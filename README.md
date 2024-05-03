# BlobstreamX Example

A simple example attempt to verify a share proof on chain using blobstreams [DAVerifier](https://github.com/celestiaorg/blobstream-contracts/blob/master/src/lib/verifier/DAVerifier.sol) contract library, as per the examples from the official [documentation](https://docs.celestia.org/developers/blobstream-proof-queries#example-rollup-that-uses-the-daverifier).

## Usage

```
go run main.go
```

This example will fetch data inclusion and share proofs from using Celestia json rpc and then attempts to verify the generated proof on chain using the deployed ShareLoader contract below.

### ShareLoader Contract

A simple contract that wraps the [DAVerifier](https://github.com/celestiaorg/blobstream-contracts/blob/master/src/lib/verifier/DAVerifier.sol) library's
`verifySharesToDataRootTupleRoot` method is deployed on Arb-Sepolia
at `0xf2787995D9eb43b57eAdB361227Ddf4FEC99b5Df`.

[View on Arbiscan](https://sepolia.arbiscan.io/address/0xf2787995D9eb43b57eAdB361227Ddf4FEC99b5Df#code)
