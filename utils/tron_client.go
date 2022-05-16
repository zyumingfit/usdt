package utils

import (
	"os"

	tronPb "github.com/TRON-US/USDT/tron/pb/api"

	"google.golang.org/grpc"
)

var client tronPb.WalletClient
var solidityClient tronPb.WalletSolidityClient

func init() {
	conn, err := TronConnection()
	if err != nil {
		os.Exit(0)
	}
	client = tronPb.NewWalletClient(conn)

	solidityConn, err := TronSolidityConnection()
	if err != nil {
		os.Exit(0)
	}
	solidityClient = tronPb.NewWalletSolidityClient(solidityConn)
}

// Structure method to get gRPC connection to tron server.
func TronConnection() (*grpc.ClientConn, error) {
	//conn, err := grpc.Dial("grpc.trongrid.io:50051", grpc.WithInsecure())
	conn, err := grpc.Dial("34.220.77.106:50051", grpc.WithInsecure())
	//conn, err := grpc.Dial("grpc.shasta.trongrid.io:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Structure method to get gRPC connection to tron server.
func TronSolidityConnection() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial("34.220.77.106:50061", grpc.WithInsecure())
	//conn, err := grpc.Dial("grpc.shasta.trongrid.io:50052", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Structure method to get tron full node gRPC client.
func TronClient() tronPb.WalletClient {
	return client
}

// Structure method to get tron solidity node gRPC client.
func TronSolidityClient() tronPb.WalletSolidityClient {
	return solidityClient
}
