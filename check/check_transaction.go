package check

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/TRON-US/USDT/log"
	tronPb "github.com/TRON-US/USDT/tron/pb/api"
	tronCore "github.com/TRON-US/USDT/tron/pb/core"
	"github.com/TRON-US/USDT/utils"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// function to check transaction.
func TransactionCheck(txSignedFile string) error {
	log.Logger().Info("Check transaction begin.")

	fileName := fmt.Sprint(txSignedFile)
	// Read transactionSigned file.
	transactions, err := os.Open(fileName)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Open file failed, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := transactions.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()

	success := 0
	failed := 0

	wg := sync.WaitGroup{}

	// New reader of transaction file.
	br := bufio.NewReader(transactions)
	for {
		// Read line.
		num, _, c := br.ReadLine()
		if c == io.EOF {
			log.Logger().Info("Read line over.")
			break
		}

		// Convert string to int.
		intNum, err := strconv.Atoi(string(num))
		if err != nil {
			log.Logger().Info("Convert string to int failed.")
			failed++
			break
		}

		// Read line.
		transactionHex, _, f := br.ReadLine()
		if f == io.EOF {
			log.Logger().Info("Read line over.")
			break
		}

		// Get raw data.
		transactionByte, err := hex.DecodeString(string(transactionHex))
		if err != nil {
			log.Logger().Error("Hex decode error.")
			failed++
			continue
		}

		transaction := &tronCore.Transaction{}
		// Unmarshal transaction.
		err = proto.Unmarshal(transactionByte, transaction)
		if err != nil {
			log.Logger().Error("Unmarshal transaction message failed.")
			failed++
			continue
		}

		// Marshal raw data.
		raw, err := proto.Marshal(transaction.GetRawData())
		if err != nil {
			log.Logger().Error("Marshal raw data failed.")
			failed++
			continue
		}
		txraw := &tronCore.TransactionRaw{}
		txraw = transaction.GetRawData()
		txcontract := &tronCore.Transaction_Contract{}
		txcontract = txraw.Contract[0]
		value := txcontract.Parameter.Value

		trigger := &tronCore.TriggerSmartContract{}
		err = proto.Unmarshal(value,trigger)
		if err != nil {
			   log.Logger().Error("Unmarshal transaction message failed.")
			   failed++
			   continue
		}
		data := hex.EncodeToString(trigger.Data)
		to_addr := "41" + data[8:72][24:]
		hexAddress, err := hex.DecodeString(to_addr)
		if err != nil {
			log.Logger().Error("Hex decode string error.")
			continue
		}

		base58Address, err := utils.Encode58Check(hexAddress)
		if err != nil {
			log.Logger().Error("Base58 encode error.")
			continue
		}


		transactionHash, err := utils.Hash(raw)

		go func(hash []byte) {
			wg.Add(1)
			defer wg.Done()
			// get transaction.
			result, err := utils.TronSolidityClient().GetTransactionById(context.Background(), &tronPb.BytesMessage{
				Value: hash,
			})
			if err != nil {
				log.Logger().Info(fmt.Sprintf("Line [%d], id [%s] receiver [%s] get transaction error, reasons: [%v]", intNum, hex.EncodeToString(transactionHash), base58Address, err))
				return
			}

			if result.GetRet() == nil || len(result.GetRet()) != 1 {
				log.Logger().Info(fmt.Sprintf("Line [%d], id [%s] receiver [%s] get transaction ret format error.", intNum, hex.EncodeToString(transactionHash), base58Address))
				return
			}

			if result.GetRet()[0].GetContractRet().String() != "SUCCESS" {
				log.Logger().Info(fmt.Sprintf("Line [%d], id [%s] receiver [%s] transaction result is not success, result: [%s].", intNum, hex.EncodeToString(transactionHash), base58Address, result.GetRet()[0].GetContractRet().String()))
				return
			}

			success++
			//successLog := fmt.Sprintf("Line [%d], id [%s] receiver [%s] transaction SUCCESS.", intNum, hex.EncodeToString(transactionHash), base58Address)
			//log.Logger().Info(successLog)

		}(transactionHash)
	}

	wg.Wait()

	if failed != 0 {
		log.Logger().Info(fmt.Sprintf("Check transaction failed, number: [%d] is unconfirmed, please retry.", failed))
		return nil
	}

	log.Logger().Info(fmt.Sprintf("Check transaction success, number: [%d]", success))
	log.Logger().Info("Check transaction end.")
	return nil

}
