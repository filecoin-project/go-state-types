module github.com/filecoin-project/go-state-types

go 1.23.0

toolchain go1.24.2

retract v0.12.7 // wrongfully skipped a patch version, use v0.12.6 or v0.12.8&^

require (
	github.com/filecoin-project/go-address v1.2.0
	github.com/filecoin-project/go-amt-ipld/v4 v4.4.0
	github.com/filecoin-project/go-bitfield v0.2.4
	github.com/filecoin-project/go-commp-utils/v2 v2.1.0
	github.com/filecoin-project/go-hamt-ipld/v3 v3.4.1
	github.com/ipfs/go-block-format v0.2.0
	github.com/ipfs/go-cid v0.5.0
	github.com/ipfs/go-ipld-cbor v0.2.0
	github.com/ipld/go-ipld-prime v0.21.0
	github.com/minio/sha256-simd v1.0.1
	github.com/multiformats/go-multibase v0.2.0
	github.com/multiformats/go-multihash v0.2.3
	github.com/multiformats/go-varint v0.0.7
	github.com/stretchr/testify v1.10.0
	github.com/whyrusleeping/cbor-gen v0.3.1
	golang.org/x/crypto v0.38.0
	golang.org/x/sync v0.14.0
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/filecoin-project/go-fil-commcid v0.2.0 // indirect
	github.com/filecoin-project/go-fil-commp-hashhash v0.2.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/ipfs/go-ipfs-util v0.0.3 // indirect
	github.com/ipfs/go-ipld-format v0.6.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.8 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.1.0 // indirect
	github.com/multiformats/go-base36 v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.89.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.3.0 // indirect
)
