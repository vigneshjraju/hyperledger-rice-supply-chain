package main

import (
	"fmt"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Submit a transaction to the Fabric network
func submitTxnFn(
	organization string,
	channelName string,
	chaincodeName string,
	contractName string,
	txnType string,
	privateData map[string][]byte,
	txnName string,
	args ...string,
) string {

	orgProfile := profile[organization]
	mspID := orgProfile.MSPID
	certPath := orgProfile.CertPath
	keyPath := orgProfile.KeyDirectory
	tlsCertPath := orgProfile.TLSCertPath
	gatewayPeer := orgProfile.GatewayPeer
	peerEndpoint := orgProfile.PeerEndpoint

	// gRPC connection
	clientConnection := newGrpcConnection(tlsCertPath, gatewayPeer, peerEndpoint)
	defer clientConnection.Close()

	// Identity and signer
	id := newIdentity(certPath, mspID)
	sign := newSign(keyPath)

	// Gateway connection
	gw, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(fmt.Errorf("failed to connect to gateway: %w", err))
	}
	defer gw.Close()

	network := gw.GetNetwork(channelName)
	contract := network.GetContractWithName(chaincodeName, contractName)

	fmt.Printf("\n--> Submitting transaction: %s\n", txnName)

	switch txnType {
	case "invoke":
		result, err := contract.SubmitTransaction(txnName, args...)
		if err != nil {
			panic(fmt.Errorf("failed to submit transaction: %w", err))
		}
		return fmt.Sprintf("*** Transaction submitted successfully: %s\n", string(result))

	case "query":
		evaluateResult, err := contract.EvaluateTransaction(txnName, args...)
		if err != nil {
			panic(fmt.Errorf("failed to evaluate transaction: %w", err))
		}

		var result string
		if isByteSliceEmpty(evaluateResult) {
			result = string(evaluateResult)
		} else {
			result = formatJSON(evaluateResult)
		}

		return result

	case "private":
		result, err := contract.Submit(
			txnName,
			client.WithArguments(args...),
			client.WithTransient(privateData),
		)
		if err != nil {
			panic(fmt.Errorf("failed to submit private transaction: %w", err))
		}
		return fmt.Sprintf("*** Transaction committed successfully\nresult: %s\n", result)

	default:
		return "Invalid transaction type"
	}
}

func isByteSliceEmpty(data []byte) bool {
	return len(data) == 0
}