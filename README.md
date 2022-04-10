# USDT Airdrop

## Outline
* [Background](#a)
* [Purpose of writing](#b)
* [Set up](#c)
  * [Set up for Ubuntu System](#c1)
  * [Set up for MacOS](#c2)
* [Compile](#d)
* [Run](#e)
  * [Set blacklist](#e1)
  * [Set whitelist](#e2)
  * [Download account file](#e3)
  * [Run cli](#e4)
* [Use makefile to build](#f)

<a name="a"></a>
## Background

*Tron*, one of the largest blockchain network in the world. It is faster and cheaper comparing to most other networks. <br/>
*Usdt*, introduced by [Tether](https://tether.to/), is the top 1 stablecoin in the blockchain industry. Most of current USDT is running on top of [Omni Layer](https://en.bitcoinwiki.org/wiki/Omni_Layer), while a few are running on top of [Ethereum network](https://www.ethereum.org/).

<a name="b"></a>
## Purpose of writing

This summary design document is mainly used to guide the detailed design of the USDT airdrop solution, providing a unified reference standard for detailed design, including the system's internal and external interfaces, system architecture, programming model and other major problems. <br/>
After this document has been peer reviewed, all detailed design of the system must be in accordance with the relevant standards and constraints of this document. In addition, this document serves as one of the criteria for peer review of detailed design documents.<br/>
In the detailed design process, if it is found that new summary design standards or constraints need to be added to guide the detailed design work, it must be updated and reviewed in this document to ensure the consistency and correctness of the detailed design of each module.<br/>
This document mainly describes the outline design of the airdrop solution, which includes defining the internal and external interfaces of the system, the related system architecture and design standards, and does not involve the details of the reality of the system business logic.

<a name="c"></a>
## Set up

You should install the necessary environment before you perform the airdrop.

<a name="c1"></a>
### Set up for Ubuntu System

1. Install prerequisites

```
$ sudo apt-get update
$ sudo apt upgrade -y
$ sudo apt install gcc -y
```

2. Install go

```
$ cd /tmp
$ GO_PACKAGE=go1.11.linux-amd64.tar.gz
$ wget https://dl.google.com/go/$GO_PACKAGE
$ sudo tar -xvf $GO_PACKAGE
$ sudo mv go /usr/local
$ sudo rm $GO_PACKAGE
```

3. Setting go environment

```
$ mkdir $HOME/go
$ export GOROOT=/usr/local/go
$ export GOPATH=$HOME/go
$ export PATH=$GOPATH/bin:$GOROOT/bin:$GOPATH:$PATH
$ go version
```

<a name="c2"></a>
### Set up for MacOS

1. Install Homebrew

```
$ /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

2. Install GoLang

```
$ mkdir $HOME/go
$ brew install go
$ go version
$ go version go1.11.2 darwin/amd64
```
Note: If Go has been installed and needs an update, use `$ brew update go`.

3. Set up environment in `~/.bash_profile`

```
$ export GOPATH="$HOME/go"
$ export PATH=$GOPATH:$PATH
```

<a name="d"></a>
## Compile

1. Download package

```
$ go install
```

2. Build project

```
$ go build
```
It will generate an executable file named **USDT** in the current directory.


<a name="e"></a>
## Run

<a name="e1"></a>
### Set blacklist

*blacklist.yaml* file must be in the same directory as **USDT** file.

Add address after `- address:`.

If you want to add more address to blacklist list, you should add line `- address: address`.

```
blacklist:
  - address: exampleA
  - address: exampleB
```

<a name="e2"></a>
### Set whitelist

*whitelist.yaml* file must be in the same directory as **USDT** file.

Add address after `- address:`.

If you want to add more address to whitelist list, you should add line `- address: address`.

If address is in the whitelist, it will be added in file base_account_(date you input).txt even if it exceeds limit.

```
whitelist:
  - address: exampleA
  - address: exampleB
```

<a name="e3"></a>
### Download account file

The naming convention for accounts file is date + '.txt', like `20190430.txt`.

You can download account file by shell.

```
$ bash download.sh -date yyyymmdd
```

And then you can get file named date in package.

<a name="e4"></a>
### Run cli

1. init

```
$ ./USDT
Please enter the command:
$ init
Please enter usdt airdrop address:
$ input hex address of airdrop.
Please enter usdt smart contract address:
$ input hex address of usdt smart contract
Please enter interest limit:
$ input limit of interest audit
Please enter the command:
$ help
[Init]                   ---Init usdt airdrop params.
[InitPrivateKey]         ---Init usdt airdrop address's private key.
[GenerateInterestFile]   ---Generate account interest file by account snapshot.
[GenerateTransactionFile]---Generate transaction file by audit file.
[Sign]                   ---Sign transaction by transaction file.
[Broadcast]              ---Broadcast transaction by transactionSigned file.
[Check]                  ---Check smart contract is success.
[Exit]                   ---Exit the transaction
```

2. Generate account interest file by account snapshot

```
Please enter the command:
$ GenerateInterestFile
Please enter date you want to deal with:
$ input the date you want to deal with(e.g. 20190430)
```

Now you will get `audit_account_(date you input).txt`, `base_account_(date you input).txt` and `blacklist_account_(date you input).txt`.<br/>

You must check address in `audit_account_(date you input).txt` and delete illegal address.<br/>

And then you must combine two files into one name by `audit.txt`

3. Generate transaction file by audit file

```
Please enter the command:
$ GenerateTransactionFile
```

Now you will get `transaction.txt`.

4. Sign transaction.

```
Please enter the command:
$ InitPrivateKey
Please enter private key:
$ airdrop private key
Please enter the command:
$ sign
```

Now you will get `transactionSigned.txt`.

5. Broadcast transaction.

```
Please enter the command:
$ Broadcast
```

6. Check transaction.

```
Please enter the command:
$ Check
```

<a name="f"></a>
## Use makefile to build

1. install

```
$ make install
```

2. build

```
$ make build
```

3. run

```
$ make run
```

4. clean

```
$ make clean
```
