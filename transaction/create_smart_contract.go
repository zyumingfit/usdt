package transaction
import (
    "bufio"
    "context"
    "encoding/hex"
    "fmt"
    "io"
    "math/big"
    "os"
    "strconv"
    "strings"
    "github.com/TRON-US/USDT/log"
    tronCore "github.com/TRON-US/USDT/tron/pb/core"
    "github.com/TRON-US/USDT/utils"
    "github.com/golang/protobuf/proto"
)
// Loop audit file, create transaction generate transaction file.
func CreateSmartContract(airDropAddress, contractAddress []byte, fileAudit string, fileTx string, fileBlacklist string) error {
    log.Logger().Info("Create smart contract begin.")
    // Create transaction file.
    transactions, err := os.Create(fmt.Sprint(fileTx))
    if err != nil {
        log.Logger().Error(fmt.Sprintf("Create transaction file error, reasons: [%v]", err))
        return err
    }
    defer func() {
        if err := transactions.Close(); err != nil {
            log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
        }
    }()
    // New writer of transactions.
    transactionsWriter := bufio.NewWriter(transactions)
    defer func() {
        _ = transactionsWriter.Flush()
    }()

    var blacklistAddresses = make(map[string]int)

    blacklistName := fmt.Sprint(fileBlacklist)
    if blacklistName != "" {

        blacklist, err := os.Open(blacklistName)
        if err != nil {
            log.Logger().Error(fmt.Sprintf("Open file failed, reasons: [%v]", err))
            return err

        }
        defer func() {
            if err := blacklist.Close(); err != nil {
                log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
            }
        }()

        r := bufio.NewReader(blacklist)
        for {
            // Read line.
            a, _, c := r.ReadLine()
            if c == io.EOF {
                log.Logger().Info("Read line over.")
                break
            }
            blacklistAddresses[string(a)] = 0
        }
    }

    // Read accounts file by date.
    fileName := fmt.Sprint(fileAudit)
    address, err := os.Open(fileName)
    if err != nil {
        log.Logger().Error(fmt.Sprintf("Open file failed, reasons: [%v]", err))
        return err
    }
    defer func() {
        if err := address.Close(); err != nil {
            log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
        }
    }()

    blacklistFound := make([]string, 0)

    success := 0
    failed := 0
    maxAccount := 0
    i := 0
    // New reader of accounts file.
    br := bufio.NewReader(address)
    for {
        // Read line.
        a, _, c := br.ReadLine()
        if c == io.EOF {
            log.Logger().Info("Read line over.")
            break
        }
        // Split the line and check line is nil.
        s := strings.Split(string(a), "\t")
        if s == nil || len(s) == 0 {
            log.Logger().Error("Line is nil")
            failed++
            continue
        }
        // Check line format.
        if len(s) != 4 {
            log.Logger().Error("Incorrect line format.")
            failed++
            continue
        }
        if _, ok := blacklistAddresses[s[1]];ok {
            blacklistFound = append(blacklistFound, s[1])
        	continue
        }
        interest, ok := big.NewInt(0).SetString(s[3], 10)
        if !ok{
                log.Logger().Error("Convert balance from string to int64 error.")
                failed++
                continue
        }
        // Transfer interest from decimal to hexadecimal.
        /*
            hexInterest := utils.DecHex(interest)
        */
        interest_hex_string := hex.EncodeToString(interest.Bytes())
        // Get data of smart contract.
        dataHex := utils.HashMethodName("transfer(address,uint256)") + utils.GetStringZero(24) + s[0][2:] +
            utils.GetStringZero(65-strings.Count(interest_hex_string, "")) + interest_hex_string
        data, err := hex.DecodeString(dataHex)
        if err != nil {
            log.Logger().Error("Hex decode error.")
            failed++
            continue
        }
        smartContract := &tronCore.TriggerSmartContract{
            OwnerAddress:    airDropAddress,
            ContractAddress: contractAddress,
            CallValue:       0,
            Data:            data,
            CallTokenValue:  0,
            TokenId:         0,
        }
        // Trigger smart contract.
        transaction, err := utils.TronClient().TriggerContract(context.Background(), smartContract)
        if err != nil {
            log.Logger().Error(fmt.Sprintf("Trigger contract error, reasons: [%v]", err))
            failed++
            continue
        }
        if !transaction.GetResult().GetResult() {
            log.Logger().Error("Transaction failed.")
            failed++
            continue
        }
        transaction.Transaction.RawData.FeeLimit = 40e6
        transaction.Transaction.RawData.Expiration = transaction.GetTransaction().GetRawData().GetExpiration() + 1000*60*60*23
        // Marshal transaction.
        rawData, err := proto.Marshal(transaction.GetTransaction())
        // Write number.
        _, err = fmt.Fprintln(transactionsWriter, strconv.Itoa(i))
        if err != nil {
            log.Logger().Error("Write num to transaction file error.")
            failed++
            continue
        }
        i++
        // Write transaction.
        _, err = fmt.Fprintln(transactionsWriter, hex.EncodeToString(rawData))
        if err != nil {
            log.Logger().Error("Write transaction hex to transaction file error.")
            failed++
            continue
        }
        success++
        log.Logger().Info(fmt.Sprintf("Line [%d] create smart contract success!", i-1))
    }
    if failed != 0 {
        log.Logger().Info(fmt.Sprintf("Create smart contract failed, number: [%d], please retry.", failed))
        return nil
    }
    if maxAccount != 0 {
        log.Logger().Info(fmt.Sprintf("Some accounts have abnormalities, number: [%d], please check.", maxAccount))
    }
    log.Logger().Info(fmt.Sprintf("Create smart contract success, number: [%d]", success))
    log.Logger().Info("Create smart contract end.")


    if len(blacklistFound) > 0{
        fmt.Println("\n\nThese addresses not been created transactions due to the blacklist:")
        for index, v := range(blacklistFound){
        	fmt.Println(index, ":", v)
        }

    }



    return nil
}
