package transaction

import (
	"bufio"
	"encoding/hex"
	"fmt"
	tronCore "github.com/TRON-US/USDT/tron/pb/core"
	"github.com/TRON-US/USDT/utils"
	"github.com/TRON-US/USDT/log"
	"github.com/golang/protobuf/proto"
	"context"
	"io"
	"os"
	"strconv"
	"strings"
)

func TransferAssetContract(airDropAddress []byte, assetName []byte, fileAudit string, fileTx string) error {
	log.Logger().Info("Create Transfer Asset Contract begin.")
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
	success := 0
	failed := 0
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
		interest, err := strconv.ParseInt(s[3], 10, 64)
		if err != nil {
			log.Logger().Error("Convert balance from string to int64 error.")
			failed++
			continue
		}

		to := fmt.Sprintf("41%s", s[0][2:])
		to_address, err := hex.DecodeString(to)
		if err != nil {
			log.Logger().Error(fmt.Sprintf("Decode hex receive address failed, reasons: [%v]\n", err))
			failed++
			continue
		}


		transferContract := &tronCore.TransferAssetContract{
			AssetName: assetName,
			OwnerAddress:    airDropAddress,
			ToAddress: to_address,
			Amount: interest,
		}
		// Trigger smart contract.
		transaction, err :=  utils.TronClient().TransferAsset2(context.Background(), transferContract)
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
	log.Logger().Info(fmt.Sprintf("Create smart contract success, number: [%d]", success))
	log.Logger().Info("Create smart contract end.")
	return nil
}

