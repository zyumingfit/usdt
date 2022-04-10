package sign

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/TRON-US/USDT/log"
	tronCore "github.com/TRON-US/USDT/tron/pb/core"
	"github.com/TRON-US/USDT/utils"

	"github.com/golang/protobuf/proto"
)

// Loop transactions file, sign transaction and generate transactionSigned file.
func TransactionSigned(key *ecdsa.PrivateKey, fileTx string, fileTxSigned string) error {
	log.Logger().Info("Sign transaction begin.")
	// Create sign transaction file.
	transactionSigned, err := os.Create(fmt.Sprint(fileTxSigned))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create signed transaction file error, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := transactionSigned.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()

	// New writer of transaction signed.
	transactionSignedWriter := bufio.NewWriter(transactionSigned)
	defer func() {
		_ = transactionSignedWriter.Flush()
	}()

	// Read transaction file.
	fileName := fmt.Sprint(fileTx)
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

	// New reader of transaction file.
	br := bufio.NewReader(transactions)
	for {
		// Read line.
		a, _, c := br.ReadLine()
		if c == io.EOF {
			log.Logger().Info("Read line over.")
			break
		}

		num := string(a)

		// Write serial number.
		_, err = fmt.Fprintln(transactionSignedWriter, num)
		if err != nil {
			log.Logger().Error("Write num to transaction file error.")
			failed++
			continue
		}

		// Read line.
		a, _, c = br.ReadLine()
		if c == io.EOF {
			log.Logger().Info("Read line over.")
			break
		}

		// Check line is nil.
		if a == nil || len(a) == 0 {
			log.Logger().Error("Line is nil")
			failed++
			continue
		}

		// Get transaction raw data.
		transactionBytes, err := hex.DecodeString(string(a))
		if err != nil {
			log.Logger().Error("Hex decode error.")
			failed++
			continue
		}

		transaction := &tronCore.Transaction{}
		// Unmarshal raw data.
		err = proto.Unmarshal(transactionBytes, transaction)
		if err != nil {
			log.Logger().Error("Hex decode error.")
			failed++
			continue
		}

		// Marshal raw data.
		rawData, err := proto.Marshal(transaction.GetRawData())
		if err != nil {
			log.Logger().Error("Marshal raw data error.")
			failed++
			continue
		}

		// Sign transaction.
		sign, err := utils.Sign(rawData, key)
		if err != nil {
			log.Logger().Error(fmt.Sprintf("Line [%s] sign error, reasons: [%v]", num, err))
			failed++
			continue
		}

		transaction.Signature = append(transaction.GetSignature(), sign)

		// Marshal transaction.
		transactionMarshal, err := proto.Marshal(transaction)
		if err != nil {
			log.Logger().Error("Marshal transaction error.")
			failed++
			continue
		}

		log.Logger().Info(fmt.Sprintf("Line [%s] sign success.", num))

		// Write transaction and sign.
		_, err = fmt.Fprintln(transactionSignedWriter, hex.EncodeToString(transactionMarshal))
		if err != nil {
			log.Logger().Error("Write signed transaction to transactionSigned file error.")
			failed++
			continue
		}
		success++
	}

	if failed != 0 {
		log.Logger().Info(fmt.Sprintf("Sign transaction failed, number: [%d], please retry.", failed))
		return nil
	}

	log.Logger().Info(fmt.Sprintf("Sign transaction success, number: [%d]", success))
	log.Logger().Info("Sign transaction end.")
	return nil
}
