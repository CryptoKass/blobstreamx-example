package client

import (
	"context"
	"fmt"

	shareloader "blobstreamx-example/contracts/ShareLoader.sol"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/tendermint/rpc/client/http"
)

type SharePointer struct {
	Height int64
	Start  int64
	End    int64
}

// GetShareProof returns the share proof for the given share pointer.
// Ready to be used with the DAVerifier library.
// RE: https://docs.celestia.org/developers/blobstream-proof-queries#example-rollup-that-uses-the-daverifier
func GetShareProof(eth *ethclient.Client, trpc *http.HTTP, sp *SharePointer) (*shareloader.SharesProof, [32]byte, error) {
	ctx := context.Background()

	// 1. Get the data commitment
	dataCommitment, err := GetDataCommitment(eth, sp.Height, 10_000_000)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to get data commitment: %w", err)
	}

	// 2. Get the block
	blockRes, err := trpc.Block(ctx, &sp.Height)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to get block: %w", err)
	}

	// 3. get data root inclusion commitment
	dcProof, err := trpc.DataRootInclusionProof(ctx, uint64(sp.Height), dataCommitment.StartBlock, dataCommitment.EndBlock)
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to get data root inclusion proof: %w", err)
	}

	// 4. get share proof
	shareProof, err := trpc.ProveShares(ctx, (uint64(sp.Height)), uint64(sp.Start), uint64(sp.End))
	if err != nil {
		return nil, [32]byte{}, fmt.Errorf("failed to get share proof: %w", err)
	}

	nonce := dataCommitment.ProofNonce.Uint64()
	height := uint64(sp.Height)

	blockDataRoot := [32]byte(blockRes.Block.DataHash)

	return &shareloader.SharesProof{
		Data:             shareProof.Data,
		ShareProofs:      toNamespaceMerkleMultiProofs(shareProof.ShareProofs),
		Namespace:        *namespace(shareProof.NamespaceID, uint8(shareProof.NamespaceVersion)),
		RowRoots:         toRowRoots(shareProof.RowProof.RowRoots),
		RowProofs:        toRowProofs(shareProof.RowProof.Proofs),
		AttestationProof: toAttestationProof(nonce, height, blockDataRoot, dcProof.Proof),
	}, blockDataRoot, nil
}
