package interest

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/TRON-US/USDT/log"
	"github.com/TRON-US/USDT/utils"
)

// Generate audit_account.txt and base_account.txt
func GenerateInterest(date string, limit int64, airdropAddress []byte) error {
	log.Logger().Info("Generate interest file begin.")

	var totalNumber = 0
	var totalBalance int64 = 0
	var totalInterestAmount int64 = 0
	var baseAmount int64 = 0
	var baseInterestAmount int64 = 0
	var baseNumber = 0
	var auditAmount int64 = 0
	var auditInterestAmount int64 = 0
	var auditNumber = 0
	var whitelistAmount int64 = 0
	var whitelistInterestAmount int64 = 0
	var whitelistNumber = 0
	var blacklistAmount int64 = 0
	var blacklistInterestAmount int64 = 0
	var blacklistNumber = 0
	var zeroAmountNumber = 0

	var accountMap map[string]string
	accountMap = make(map[string]string)

	// Transfer string date to int date.
	intDate, err := strconv.Atoi(date)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Convert date from string to int failed, reasons: [%v]", err))
		return err
	}

	// Get rate by date
	rate, err := getRate(intDate)
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Get interest rate by date failed, reasons: [%v]", err))
		return err
	}

	// Read accounts file by date.
	fileName := fmt.Sprint("./", date, ".txt")
	accounts, err := os.Open(fileName)
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
	baseAccounts, err := os.Create(fmt.Sprint("base_account_", date, ".txt"))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create base account file error, reasons: [%v]", err))
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

	// Create base accounts file.
	auditAccounts, err := os.Create(fmt.Sprint("audit_account_", date, ".txt"))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create audit account file error, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := auditAccounts.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()

	// New writer of audit account.
	auditAccountsWriter := bufio.NewWriter(auditAccounts)
	defer func() {
		_ = auditAccountsWriter.Flush()
	}()

	// Create blacklist accounts file.
	blacklistAccounts, err := os.Create(fmt.Sprint("blacklist_account_", date, ".txt"))
	if err != nil {
		log.Logger().Error(fmt.Sprintf("Create blacklist account file error, reasons: [%v]", err))
		return err
	}
	defer func() {
		if err := blacklistAccounts.Close(); err != nil {
			log.Logger().Error(fmt.Sprintf("File close failed, reasons: [%v]", err))
		}
	}()

	// New writer of blacklist account.
	blacklistAccountsWriter := bufio.NewWriter(blacklistAccounts)
	defer func() {
		_ = blacklistAccountsWriter.Flush()
	}()

	// New reader of accounts file.
	br := bufio.NewReader(accounts)
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
			continue
		}

		// Check line format.
		if len(s) != 2 {
			log.Logger().Error("Incorrect line format.")
			continue
		}

		if accountMap[s[0]] != "" {
			log.Logger().Info(fmt.Sprintf("Account [%s] is repeated!", s[0]))
		}
		accountMap[s[0]] = s[1]

		if strings.EqualFold(hex.EncodeToString(airdropAddress), s[0]) {
			log.Logger().Info(fmt.Sprintf("Address [%s] is airdrop address.", s[0]))
			continue
		}

		hexAddress, err := hex.DecodeString(s[0])
		if err != nil {
			log.Logger().Error("Hex decode string error.")
			continue
		}

		base58Address, err := utils.Encode58Check(hexAddress)
		if err != nil {
			log.Logger().Error("Base58 encode error.")
			continue
		}

		// Transfer account balance to int64.
		int64Balance, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			log.Logger().Error("Convert balance from string to int64 error.")
			continue
		}

		totalBalance = totalBalance + int64Balance

		// Calculating interest and write file.
		amount := float64(int64Balance) * rate
		if int64(amount) != 0 {
			line := fmt.Sprint(s[0], "\t", base58Address, "\t", s[1], "\t", int64(amount))

			// Get blacklist.
			blacklist, err := GetBlacklist()
			if err != nil {
				log.Logger().Error(fmt.Sprintf("Get blacklist failed, reasons: [%v]", err))
				return err
			}

			// Check account is in blacklist.
			blacklistFlg := false
			for i := 0; i < len(blacklist); i++ {
				if strings.EqualFold(blacklist[i].Address, s[0]) {
					log.Logger().Info(fmt.Sprintf("Account[%s] is in blacklist.", s[0]))
					blacklistFlg = true
					break
				}
			}
			if blacklistFlg {
				_, err = fmt.Fprintln(blacklistAccounts, line)
				if err != nil {
					log.Logger().Error("Write line to blacklist account file error.")
					continue
				}
				blacklistNumber++
				blacklistAmount = blacklistAmount + int64Balance
				blacklistInterestAmount = blacklistInterestAmount + int64(amount)
				totalNumber++
				totalInterestAmount = totalInterestAmount + int64(amount)
				continue
			}

			// Get whitelist.
			whitelist, err := GetWhitelist()
			if err != nil {
				log.Logger().Error(fmt.Sprintf("Get whitelist failed, reasons: [%v]", err))
				return err
			}

			// Check account is in whitelist.
			whitelistFlg := false
			for i := 0; i < len(whitelist); i++ {
				if strings.EqualFold(whitelist[i].Address, s[0]) {
					log.Logger().Info(fmt.Sprintf("Account[%s] is in whitelist.", s[0]))
					whitelistFlg = true
					break
				}
			}
			if whitelistFlg {
				_, err = fmt.Fprintln(baseAccountsWriter, line)
				if err != nil {
					log.Logger().Error("Write line to base account file error.")
				}
				whitelistNumber++
				whitelistAmount = whitelistAmount + int64Balance
				whitelistInterestAmount = whitelistInterestAmount + int64(amount)
				baseNumber++
				baseAmount = baseAmount + int64Balance
				baseInterestAmount = baseInterestAmount + int64(amount)
				totalNumber++
				totalInterestAmount = totalInterestAmount + int64(amount)
				continue
			}

			if int64(amount) <= limit {
				_, err = fmt.Fprintln(baseAccountsWriter, line)
				if err != nil {
					log.Logger().Error("Write line to base account file error.")
					continue
				}
				baseNumber++
				baseAmount = baseAmount + int64Balance
				baseInterestAmount = baseInterestAmount + int64(amount)
				totalNumber++
				totalInterestAmount = totalInterestAmount + int64(amount)
			} else {
				_, err = fmt.Fprintln(auditAccountsWriter, line)
				if err != nil {
					log.Logger().Error("Write line to audit account file error.")
					continue
				}
				auditNumber++
				auditAmount = auditAmount + int64Balance
				auditInterestAmount = auditInterestAmount + int64(amount)
				totalNumber++
				totalInterestAmount = totalInterestAmount + int64(amount)
			}
		} else {
			log.Logger().Info(fmt.Sprintf("Account[%s] Interest is 0.", s[0]))
			zeroAmountNumber++
			totalNumber++
		}

	}

	log.Logger().Info("Generate interest file end.")
	log.Logger().Info(fmt.Sprintf("\nReport:\n"+
		"\tAnnual percentage rate:                      [%d]\n"+
		"\tDaily percentage rate:                       [%f]\n"+
		"\tAll USDT accounts  on TRON:                  [%d]\n"+
		"\tAll USDT balance on TRON:                    [%s]\n"+
		"\tTotal USDT if airdropping for all:           [%s]\n"+
		"\tZero USDT balance accounts  on TRON:         [%d]\n"+
		"\tBlacklisted accounts  on TRON:               [%d]\n"+
		"\tBlacklisted USDT amount on TRON:             [%s]\n"+
		"\tBlacklisted airdrop USDT amount:             [%s]\n"+
		"\tTo be audited accounts  on TRON:             [%d]\n"+
		"\tTo be audited USDT amount on TRON:           [%s]\n"+
		"\tTo be audited airdrop USDT amount:           [%s]\n"+
		"\tWhitelist accounts  on TRON:                 [%d]\n"+
		"\tWhitelist USDT amount on TRON:               [%s]\n"+
		"\tWhitelist airdrop USDT amount:               [%s]\n"+
		"\tActual USDT accounts  on TRON:               [%d]\n"+
		"\tActual USDT amount on TRON:                  [%s]\n"+
		"\tActual airdrop USDT amount:                  [%s]\n", int(rate*365*100), rate, totalNumber, intAmountFormat(totalBalance),
		intAmountFormat(totalInterestAmount), zeroAmountNumber, blacklistNumber, intAmountFormat(blacklistAmount),
		intAmountFormat(blacklistInterestAmount), auditNumber, intAmountFormat(auditAmount),
		intAmountFormat(auditInterestAmount), whitelistNumber, intAmountFormat(whitelistAmount),
		intAmountFormat(whitelistInterestAmount), baseNumber, intAmountFormat(baseAmount),
		intAmountFormat(baseInterestAmount)))
	return nil
}

// Get interest rate by date.
func getRate(date int) (float64, error) {
	if date == 20190716 { //1 days
		return stageA, nil
	} else if date == 20190717 { //1 days
		return stageB, nil
	} else if date >= 20190718 && date <= 20190727 { //10 days
		return stageC, nil
	} else if date >= 20190728 && date <= 20190815 { //19 days
		return stageD, nil
	} else {
		return 0, errors.New("date illegal")
	}
}

func intAmountFormat(amount int64) string {
	s := strconv.FormatInt(amount, 10)
	length := len(s)
	if length <= 6 {
		return fmt.Sprint("0.", utils.GetStringZero(6-length), s)
	} else {
		decimal := s[length-6:]
		steamed := s[:length-6]

		num := 0
		if len(steamed)%3 == 0 {
			num = (len(steamed) / 3) - 1
		} else {
			num = len(steamed) / 3
		}

		s1 := strings.Split(steamed, "")
		s2 := make([]string, len(s1)+num)

		count := 0

		l := len(s1) - 1

		for i := len(s2) - 1; i >= 0; i-- {
			if count == 3 {
				count = 0
				s2[i] = ","
				continue
			}

			s2[i] = s1[l]
			l--
			count++
		}

		amountString := ""
		for i := 0; i < len(s2); i++ {
			amountString = amountString + s2[i]
		}

		return amountString + "." + decimal
	}
}
