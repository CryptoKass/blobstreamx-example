package main

import (
	"blobstreamx-example/client"
	shareloader "blobstreamx-example/contracts/ShareLoader.sol"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	// step 0. get eth-rpc and trpc endpoint from args
	ethEndpoint := os.Args[1]
	trpcEndpoint := os.Args[2]
	shareloaderAddress := common.HexToAddress(os.Args[3]) // e.g. 0xf2787995D9eb43b57eAdB361227Ddf4FEC99b5Df

	// step 1: connect to eth and trpc endpoints
	eth, err := ethclient.Dial(ethEndpoint)
	if err != nil {
		panic(fmt.Errorf("failed to connect to the Ethereum node: %w", err))
	}
	trpc, err := http.New(trpcEndpoint, "/websocket")
	if err != nil {
		panic(fmt.Errorf("failed to connect to the Tendermint RPC: %w", err))
	}
	fmt.Println("Successfully connected to the Ethereum node and Tendermint RPC")

	// step 2: generate share proof
	sp := &client.SharePointer{
		Height: 1490181,
		Start:  1,
		End:    14,
	}

	proof, root, err := client.GetShareProof(eth, trpc, sp)
	if err != nil {
		panic(fmt.Errorf("failed to get share proof: %w", err))
	}
	fmt.Println("Successfully generated share proof")

	// step 3: verify the share proof
	loader, err := shareloader.NewShareloader(shareloaderAddress, eth)
	if err != nil {
		panic(fmt.Errorf("failed to instantiate shareloader contract: %w", err))
	}

	valid, errCodes, err := loader.VerifyShares(nil, *proof, root)
	if err != nil {
		panic(fmt.Errorf("failed to verify share proof: %w", err))
	}

	fmt.Println("Proof valid:", valid, "Error codes:", errCodes)
}
