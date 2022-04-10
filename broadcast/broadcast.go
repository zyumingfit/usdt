package broadcast

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/TRON-US/USDT/log"
	tronCore "github.com/TRON-US/USDT/tron/pb/core"
	"github.com/TRON-US/USDT/utils"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

// function to sign and broadcast transaction.
func TransactionBroadcast(batch uint32, txSignedFile string) error {
	log.Logger().Info("Broadcast transaction begin.")

	// Read transactionSigned file.
	fileName := fmt.Sprint(txSignedFile)
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

	// Create failed transactions file.
	failedTransactions, err := os.Create(fmt.Sprint("failedTransaction.txt"))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create transaction file error, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := failedTransactions.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()

	// New writer of transactions.
	failedTransactionsWriter := bufio.NewWriter(failedTransactions)
	defer func() {
		_ = failedTransactionsWriter.Flush()
	}()

	success := 0
	failed := 0

	broadcastBatchhCnt := uint32(0)
	batchIntervalTime := 3

	// New reader of transaction file.
	br := bufio.NewReader(transactions)
	for {
		if batchIntervalTime > 0 && broadcastBatchhCnt == batch{
			log.Logger().Info(fmt.Sprintf("Wait for %d second.....", batchIntervalTime))
			timeAfter := time.After(time.Second * time.Duration(batchIntervalTime))
			 <-timeAfter
			 broadcastBatchhCnt = 0
		}
		broadcastBatchhCnt +=1


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
			continue
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
		// Unmarshal raw data.
		err = proto.Unmarshal(transactionByte, transaction)
		if err != nil {
			log.Logger().Error("Hex decode error.")
			failed++
			continue
		}

		// broadcast transaction
		result, err := utils.TronClient().BroadcastTransaction(context.Background(), transaction)
		if err != nil {
			log.Logger().Info(fmt.Sprintf("Line [%d] broadcast contract error, reasons: [%v]", intNum, err))
			// Write number.
			_, err = fmt.Fprintln(failedTransactionsWriter, strconv.Itoa(intNum))
			if err != nil {
				log.Logger().Error("Write num to failed transaction file error.")
				failed++
				continue
			}

			// Write failed transaction.
			_, err = fmt.Fprintln(failedTransactionsWriter, string(transactionHex))
			if err != nil {
				log.Logger().Error("Write transaction to failed transaction file error.")
				failed++
				continue
			}
			failed++
			continue
		}
		if !result.GetResult() {
			log.Logger().Info(fmt.Sprintf("Line [%d] broadcast contract failed, reasons: [%s].", intNum, string(result.Message)))
			// Write number.
			_, err = fmt.Fprintln(failedTransactionsWriter, strconv.Itoa(intNum))
			if err != nil {
				log.Logger().Error("Write num to failed transaction file error.")
				failed++
				continue
			}

			// Write failed transaction.
			_, err = fmt.Fprintln(failedTransactionsWriter, string(transactionHex))
			if err != nil {
				log.Logger().Error("Write transaction to failed transaction file error.")
				failed++
				continue
			}
			failed++
			continue
		}
		success++
		log.Logger().Info(fmt.Sprintf("Line [%d] broadcast SUCCESS.", intNum))
	}

	if failed != 0 {
		log.Logger().Info(fmt.Sprintf("Broadcast transaction failed, number: [%d], please retry.", failed))
		return nil
	}

	log.Logger().Info(fmt.Sprintf("Broadcast transaction success, number: [%d]", success))
	log.Logger().Info("Broadcast transaction end.")
	return nil

}
