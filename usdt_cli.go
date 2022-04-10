package main

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/TRON-US/USDT/utils"
	"github.com/TRON-US/USDT/wink"
	"os"
	"strconv"
	"strings"

	"github.com/TRON-US/USDT/broadcast"
	"github.com/TRON-US/USDT/check"
	"github.com/TRON-US/USDT/interest"
	"github.com/TRON-US/USDT/log"
	"github.com/TRON-US/USDT/sign"
	"github.com/TRON-US/USDT/transaction"
	eth "github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type UsdtConfig struct {
	pk              *ecdsa.PrivateKey
	limit           int64
	rate            float64
	airdropAddress  []byte
	contractAddress []byte
	precision int64
	batch 			int
	assetName []byte
}

// Init signed private key by user input.
func (usdt *UsdtConfig) InitPrivateKey() {
	// Get hex private key by keyboard input.
	fmt.Println("Please enter private key:")
	f := bufio.NewReader(os.Stdin)
	hexPrivateKey, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Convert private key from hex to struct.
	ecdsaPK, err := eth.HexToECDSA(strings.Replace(hexPrivateKey, "\n", "", -1))
	if err != nil {
		fmt.Printf("Convert private key to ECDSA format failed, reasons: [%v]\n", err)
		return
	}

	usdt.pk = ecdsaPK
}

// Init usdt airdrop params.
func (usdt *UsdtConfig) Init() {
	// Get airdrop address by keyboard input.
	fmt.Println("Please enter airdrop address:")
	f := bufio.NewReader(os.Stdin)
	hexAirdropAddress, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Decode hex airdrop address.
	usdt.airdropAddress, err = hex.DecodeString(strings.Replace(hexAirdropAddress, "\n", "", -1))
	if err != nil {
		fmt.Printf("Decode hex airdrop address failed, reasons: [%v]\n", err)
		return
	}

	// Get smart contract address by keyboard input.
	fmt.Println("Please enter token contract address:")
	f = bufio.NewReader(os.Stdin)
	hexContractAddress, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Decode hex contract address.
	usdt.contractAddress, err = hex.DecodeString(strings.Replace(hexContractAddress, "\n", "", -1))
	if err != nil {
		fmt.Printf("Decode hex smart contract address failed, reasons: [%v]\n", err)
		return
	}

	// Get precision
	fmt.Println("Please enter tonken precision:")
	f = bufio.NewReader(os.Stdin)
	precisionString, err := f.ReadString('\n')
	if err != nil {
			fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
			return
		}
	usdt.precision, err = strconv.ParseInt(strings.Replace(precisionString, "\n", "", -1), 10, 64)
	if err != nil {
			fmt.Printf("Convert precision to int64, reasons: [%v]\n", err)
			return
		}

	fmt.Println("Please enter rate:")
	f = bufio.NewReader(os.Stdin)
	rateString, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Convert interest rate to int64.
	usdt.rate, err = strconv.ParseFloat(strings.Replace(rateString, "\n", "", -1), 64)
	if err != nil {
		fmt.Printf("Convert interest rate to int64, reasons: [%v]\n", err)
		return
	}
}

// Init usdt airdrop params.
func (usdt *UsdtConfig) InitTrxAirdrop() {
	// Get airdrop address by keyboard input.
	fmt.Println("Please enter airdrop address:")
	f := bufio.NewReader(os.Stdin)
	hexAirdropAddress, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Decode hex airdrop address.
	usdt.airdropAddress, err = hex.DecodeString(strings.Replace(hexAirdropAddress, "\n", "", -1))
	if err != nil {
		fmt.Printf("Decode hex airdrop address failed, reasons: [%v]\n", err)
		return
	}
}

// Init usdt airdrop params.
func (usdt *UsdtConfig) InitTrc10Airdrop() {
	// Get airdrop address by keyboard input.
	fmt.Println("Please enter airdrop address:")
	f := bufio.NewReader(os.Stdin)
	hexAirdropAddress, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Decode hex airdrop address.
	usdt.airdropAddress, err = hex.DecodeString(strings.Replace(hexAirdropAddress, "\n", "", -1))
	if err != nil {
		fmt.Printf("Decode hex airdrop address failed, reasons: [%v]\n", err)
		return
	}

	fmt.Println("Please enter asset name:")
	f = bufio.NewReader(os.Stdin)
	assetName, err := f.ReadString('\n')
	if err != nil{
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}
	assetName = strings.Replace(assetName, "\n", "", -1)
	usdt.assetName = []byte(assetName)
}


// Generate account interest file by account snapshot.
func (usdt *UsdtConfig) GenerateInterestFile() {
	if usdt.limit == 0 {
		fmt.Println("You should init first!")
		return
	}

	// Get date string by keyboard input.
	fmt.Println("Please enter date you want to deal with:")
	f := bufio.NewReader(os.Stdin)
	date, err := f.ReadString('\n')
	if err != nil {
		fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
		return
	}

	// Generate interest file.
	err = interest.GenerateInterest(strings.Replace(date, "\n", "", -1), usdt.limit, usdt.airdropAddress)
	if err != nil {
		fmt.Printf("Generate interest file error, reasons: [%v]\n", err)
	}
}

// Generate account interest file by account snapshot.
func (usdt *UsdtConfig) GenerateWinkInterestFile() {
	if usdt.rate == 0 || usdt.airdropAddress == nil {
		fmt.Println("You should init first!")
		return
	}

	// Generate interest file.
	err := wink.GenerateInterest(usdt.airdropAddress, usdt.rate, usdt.precision)
	if err != nil {
		fmt.Printf("Generate interest file error, reasons: [%v]\n", err)
	}
}

// Generate transaction file by audit file.
func (usdt *UsdtConfig) GenerateTransactionFile() {
	if usdt.airdropAddress == nil || len(usdt.airdropAddress) == 0 || usdt.contractAddress == nil || len(usdt.contractAddress) == 0 {
		fmt.Println("You should init first!")
		return
	}

	// Generate transaction file.
	err := transaction.CreateSmartContract(usdt.airdropAddress, usdt.contractAddress, "audit.txt", "transaction.txt", "")
	if err != nil {
		fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
	}
}

func (usdt *UsdtConfig) GenerateTrxTransactionFile() {
	if usdt.airdropAddress == nil || len(usdt.airdropAddress) == 0{
		fmt.Println("You should init first!")
		return
	}

	// Generate transaction file.
	err := transaction.TransferContract(usdt.airdropAddress, "audit.txt", "transaction.txt")
	if err != nil {
		fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
	}
}
func (usdt *UsdtConfig) GenerateTrc10TransactionFile() {
	if usdt.airdropAddress == nil || len(usdt.airdropAddress) == 0 || usdt.assetName == nil || len(usdt.assetName) == 0{
		fmt.Println("You should init first!")
		return
	}

	// Generate transaction file.
	err := transaction.TransferAssetContract(usdt.airdropAddress, usdt.assetName, "audit.txt", "transaction.txt")
	if err != nil {
		fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
	}
}

// Generate transactionSigned file by transaction file.
func (usdt *UsdtConfig) Sign() {
	if usdt.pk == nil {
		fmt.Println("You should init private key first!")
		return
	}

	// Sign transaction
	err := sign.TransactionSigned(usdt.pk, "transaction.txt", "transactionSigned.txt")
	if err != nil {
		fmt.Printf("Sign transactions error, reasons: [%v]\n", err)
	}
}

// Broadcast signed transaction by file.
func (usdt *UsdtConfig) Broadcast() {
	err := broadcast.TransactionBroadcast(uint32(usdt.batch), "transactionSigned.txt")
	if err != nil {
		fmt.Printf("Broadcast transactions error, reasons: [%v]\n", err)
	}
}

// Transaction check.
func (usdt *UsdtConfig) Check() {
	err := check.TransactionCheck("transactionSigned.txt")
	if err != nil {
		fmt.Printf("Transactions check error, reasons: [%v]\n", err)
	}
}

func (usdt *UsdtConfig) RunCli() {
	for again := true; again; {
		f := bufio.NewReader(os.Stdin)
		fmt.Println("Please enter the command: ")
		cmm, err := f.ReadString('\n')
		if err != nil {
			fmt.Printf("There were errors reading, err:[%v], exiting program.\n", err)
			continue
		}
		switch {
		case strings.EqualFold(cmm, "InitTrc20Airdrop\n"):
			usdt.Init()
		case strings.EqualFold(cmm, "InitTrxAirdrop\n"):
			usdt.InitTrxAirdrop()
		case strings.EqualFold(cmm, "InitTrc10Airdrop\n"):
			usdt.InitTrc10Airdrop()
		case strings.EqualFold(cmm, "InitPrivateKey\n"):
			usdt.InitPrivateKey()
		case strings.EqualFold(cmm, "GenerateInterestFile\n"):
			usdt.GenerateInterestFile()
		case strings.EqualFold(cmm, "GenerateWinkInterestFile\n"):
			usdt.GenerateWinkInterestFile()
		case strings.EqualFold(cmm, "GenerateTransactionFile\n"):
			usdt.GenerateTransactionFile()
		case strings.EqualFold(cmm, "GenerateTrxTransactionFile\n"):
			usdt.GenerateTrxTransactionFile()
		case strings.EqualFold(cmm, "GenerateTrc10TransactionFile\n"):
			usdt.GenerateTrc10TransactionFile()
		case strings.EqualFold(cmm, "Sign\n"):
			usdt.Sign()
		case strings.EqualFold(cmm, "Broadcast\n"):
			//Get batch number
			fmt.Println("Number of transactions broadcast in a batch:")
			f = bufio.NewReader(os.Stdin)
			batchNumber, err := f.ReadString('\n')
			if err != nil{
				fmt.Printf("Read string from buffer error, reasons: [%v]\n", err)
				return
			}
			batch, err := strconv.ParseInt(strings.Replace(batchNumber, "\n", "", -1), 10, 32)
			if err != nil{
				fmt.Printf("Convert batch to int64, reasons: [%v]\n", err)
				return

			}
			usdt.batch = int(batch)
			usdt.Broadcast()
		case strings.EqualFold(cmm, "Check\n"):
			usdt.Check()
		case strings.EqualFold(cmm, "help\n"):
			helpString := fmt.Sprint(
				"[InitTrc20Airdrop]	---Init trc20 airdrop params.\n",
				"[InitTrxAirdrop]          ---Init trx airdrop params.\n",
				"[InitTrc10Airdrop]          ---Init trx airdrop params.\n",
				"[InitPrivateKey]          ---Init usdt airdrop address's private key.\n",
				"[GenerateInterestFile]    ---Generate account interest file by account snapshot.\n",
				"[GenerateWinkInterestFile]---Generate wink airdrop file by account snapshot.\n",
				"[GenerateTransactionFile] ---Generate transaction file by audit file.\n",
				"[GenerateTrxTransactionFile] ---Generate Trx airdrop transaction file by audit file.\n",
				"[GenerateTrc10TransactionFile] ---Generate Trx airdrop transaction file by audit file.\n",
				"[Sign]                    ---Sign transaction by transaction file.\n",
				"[Broadcast]               ---Broadcast transaction by transactionSigned file.\n",
				"[Check]                   ---Check smart contract is success.\n",
				"[Exit]                    ---Exit the transaction")
			fmt.Println(helpString)
		case strings.EqualFold(cmm, "Exit\n"):
			again = false
		default:
			fmt.Println("Wrong inputting command, please enter 'help' to get the commands.")
		}
	}
}
type Conf struct {
	AirdropAddress 	string `yaml:"airdrop_address"`
	PrivateKey 		string `yaml:"private_key"`
	Fullnode 		string  `yaml:"fullnode"`
	Solditynode 	string `yaml:"soliditynode"`
	Batch 			uint32 `yaml:"batch"`
	Airdrops 		map[string]AirDropConf `yaml:"airdrops"`
}
type AirDropConf struct {
	Type 			string `ymal:"type"`
	Decimals 		uint32 `yaml:"decimals"`
	AssetName       string `yaml:"assetname"`
	ContractAddress string `yaml:"contract_address"`
	Enable 			bool `yaml:"enable"`
	Audit 			string `yaml:"audit"`
	TxFile 			string `yaml:"tx_file"`
	SignedFile 		string `yaml:"signed_file"`
	BlacklistFile 		string `yaml:"blacklist_file"`
}

func (c *Conf) GetConf(file string) *Conf {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

func (c *Conf) Generate() {
	for name := range c.Airdrops {
		if c.Airdrops[name].Enable == false {
			continue
		}

		addressByte, err := utils.Decode58Check(c.AirdropAddress)
		if err != nil {
			addressByte, _ = hex.DecodeString(c.AirdropAddress)
		}
		auditFile := c.Airdrops[name].Audit
		txFile :=c.Airdrops[name].TxFile
		blacklistFile := c.Airdrops[name].BlacklistFile
		fmt.Println("\n\n\n-------------------------------------------------------------")
		fmt.Printf("[%s] generating....\n", name, )
		fmt.Println("-------------------------------------------------------------")

		tokenType := c.Airdrops[name].Type
		switch {
		case tokenType == "TRX":
			err := transaction.TransferContract(addressByte, auditFile, txFile)
			if err != nil {
				fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
			}
		case tokenType == "TRC10":
			assetname := []byte(c.Airdrops[name].AssetName)
			err := transaction.TransferAssetContract(addressByte, assetname, auditFile, txFile)
			if err != nil {
				fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
			}
		case tokenType == "TRC20":
			contractAddrByte, err := utils.Decode58Check(c.Airdrops[name].ContractAddress)
			if err != nil {
				contractAddrByte, _ = hex.DecodeString(c.AirdropAddress)
			}
			err = transaction.CreateSmartContract(addressByte, contractAddrByte, auditFile, txFile, blacklistFile)
			if err != nil {
				fmt.Printf("Generate transaction file error, reasons: [%v]\n", err)
			}
		default:
			log.Logger().Warn(fmt.Sprintf("unknow type [%s]", tokenType))
			return
		}
	}
}
func (c *Conf) Sign(){
	for name := range c.Airdrops {
		if c.Airdrops[name].Enable == false{
			continue
		}

		fmt.Println("\n\n\n-------------------------------------------------------------")
		fmt.Printf("[%s] signing....\n" ,name, )
		fmt.Println("-------------------------------------------------------------")
		ecdsaPK, err := eth.HexToECDSA(strings.Replace(c.PrivateKey, "\n", "", -1))
		if err != nil {
			fmt.Printf("Convert private key to ECDSA format failed, reasons: [%v]\n", err)
			return
		}
		err = sign.TransactionSigned(ecdsaPK, c.Airdrops[name].TxFile, c.Airdrops[name].SignedFile)
		if err != nil {
			fmt.Printf("Sign transactions error, reasons: [%v]\n", err)
		}
	}
}

func (c *Conf) Broadcast() {
	batch := c.Batch

	for name := range c.Airdrops {
		if c.Airdrops[name].Enable == false {
			continue
		}
		txSignedFile := c.Airdrops[name].SignedFile

		fmt.Println("\n\n\n-------------------------------------------------------------")
		fmt.Printf("[%s] broadcast....\n", name, )
		fmt.Println("-------------------------------------------------------------")
		err := broadcast.TransactionBroadcast(batch, txSignedFile)
		if err != nil {
			fmt.Printf("Broadcast transactions error, reasons: [%v]\n", err)
		}
	}
}

func (c *Conf) Check() {

	for name := range c.Airdrops {
		if c.Airdrops[name].Enable == false {
			continue
		}

		txSignedFile := c.Airdrops[name].SignedFile
		err := check.TransactionCheck(txSignedFile)
		if err != nil {
			fmt.Printf("Transactions check error, reasons: [%v]\n", err)
		}
	}
}

func main() {

	generate := flag.Bool("generate", false, "Generate transaction files.")
	sign := flag.Bool("sign", false, "Sign transaction files.")
	broadcast := flag.Bool("broadcast", false, "Broadcast signed transaction files.")
	check := flag.Bool("check", false, "Check all broadcasted transactions.")
	config := flag.String("config", "", "Specify the configuration file, the default is conf.yml.")
	interactive := flag.Bool("interactive", false, "Start the interactive command line, the default is non-interactive.")
	flag.Parse()

	if *interactive == true{
		var c UsdtConfig
		c.RunCli()
		return
	}
	funcCtr := 0
	if *generate  == true{
		funcCtr += 1
	} else if *sign == true{
		funcCtr += 1
	}else if *broadcast == true{
		funcCtr += 1
	}else if *check == true{
		funcCtr += 1
	}
	if funcCtr == 0 || funcCtr > 1{
		fmt.Println("Please select a function:\n"+
			"\t--generate: Generate transaction files.\n" +
			"\t--sign: Sign transaction files.\n" +
			"\t--broadcast: Broadcast signed transaction files.\n")
		return
	}

	configPath := "conf.yml"
	if *config != "" {
		configPath = *config
	}


	info := Conf{}
	conf := info.GetConf(configPath)
	if *generate == true{
		conf.Generate()
	}else if *sign == true{
		conf.Sign()
	}else if *broadcast == true{
		conf.Broadcast()
	}else if *check == true{
		conf.Check()
	}

}
