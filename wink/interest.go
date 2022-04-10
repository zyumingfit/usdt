package wink
import (
	"bufio"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"os"
	"strings"
	"github.com/TRON-US/USDT/log"
	"github.com/TRON-US/USDT/utils"
)
// Generate audit.txt by account.csv.
func GenerateInterest(airdropAddress []byte, rate float64, precision int64) error {
	log.Logger().Info("Generate interest file begin.")
	var accountMap map[string]string
	accountMap = make(map[string]string)
	// Read accounts file.
	accounts, err := os.Open("./account.csv")
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Open file failed, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := accounts.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()
	// Create base accounts file.
	baseAccounts, err := os.Create("audit.txt")
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create audit file error, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := baseAccounts.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()
	// New writer of base account.
	baseAccountsWriter := bufio.NewWriter(baseAccounts)
	defer func() {
		_ = baseAccountsWriter.Flush()
	}()
	//totalAirdropAmount := int64(0)
	totalAirdropAmount := big.NewFloat(0)
	failedNum := 0
	// New reader of accounts file.
	reader := csv.NewReader(accounts)
	_, err = reader.Read()
	if err == io.EOF {
		log.Logger().Info("Read line over.")
		return err
	}
	for {
		// Read line.
		record, err := reader.Read()
		if err == io.EOF {
			log.Logger().Info("Read line over.")
			break
		}
		if record == nil || len(record) != 2 {
			log.Logger().Info("Column format error.")
			return errors.New("column format error")
		}
		if accountMap[record[0]] != "" {
			log.Logger().Info(fmt.Sprintf("Account [%s] is repeated!", record[0]))
		}
		accountMap[record[0]] = record[1]
		// Decode address to byte.
		addressByte, err := utils.Decode58Check(record[0])
		if err != nil {
			log.Logger().Error(fmt.Sprintf("Decode address [%s] error, reasons: [%v]", record[0], err))
			return err
		}
		hexAddress := hex.EncodeToString(addressByte)
		if strings.EqualFold(hex.EncodeToString(airdropAddress), hexAddress) {
			log.Logger().Info(fmt.Sprintf("Address [%s] is airdrop address.", hexAddress))
			continue
		}
		// Transfer account balance to int64.
		balance, ok := big.NewFloat(0).SetString(record[1])
		if !ok {
			log.Logger().Error("Convert balance from string to int64 error.")
			failedNum++
			continue
		}

		/*
			balance, err := strconv.ParseInt(record[1], 10, 64)
			if err != nil {
				log.Logger().Error("Convert balance from string to int64 error.")
				failedNum++
				continue
			}
		*/
		// Calculating interest and write file.
		//amount := float64(balance) / rate
		_rate := big.NewFloat(rate * 1000000)

		_precision := math.Pow10(int(precision))
		amount := balance.Quo(balance, _rate)
		amount = amount.Mul(amount, big.NewFloat(_precision))
		if amount.Cmp(big.NewFloat(0)) > 0 {
			line := fmt.Sprint(hexAddress, "\t", record[0], "\t", record[1], "\t", amount.Text('f', 0))
			_, err = fmt.Fprintln(baseAccountsWriter, line)
			if err != nil {
				log.Logger().Error("Write line to base account file error.")
				failedNum++
				continue
			}
			totalAirdropAmount = totalAirdropAmount.Add(totalAirdropAmount, amount)
		} else {
			log.Logger().Info(fmt.Sprintf("Account[%s] Interest is 0.", record[0]))
		}
	}
	log.Logger().Info(fmt.Sprintf("Generate interest file end, airdrop amount: %s.", totalAirdropAmount.Text('f', 0)))
	return nil
}
