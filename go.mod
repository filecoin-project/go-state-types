module github.com/filecoin-project/go-state-types

go 1.18

retract v0.12.7 // wrongfully skipped a patch version, use v0.12.6 or v0.12.8&^

require (
	github.com/filecoin-project/go-address v1.1.0
	github.com/filecoin-project/go-amt-ipld/v4 v4.2.0
	github.com/filecoin-project/go-bitfield v0.2.4
	github.com/filecoin-project/go-commp-utils/nonffi v0.0.0-20220905160352-62059082a837
	github.com/filecoin-project/go-hamt-ipld/v3 v3.1.0
	github.com/ipfs/go-block-format v0.0.3
	github.com/ipfs/go-cid v0.3.2
	github.com/ipfs/go-ipld-cbor v0.0.6
	github.com/ipld/go-ipld-prime v0.19.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/minio/sha256-simd v1.0.1
	github.com/multiformats/go-multibase v0.0.3
	github.com/multiformats/go-multihash v0.2.1
	github.com/multiformats/go-varint v0.0.6
	github.com/stretchr/testify v1.7.0
	github.com/whyrusleeping/cbor-gen v0.0.0-20230923211252-36a87e1ba72f
	golang.org/x/crypto v0.1.0
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/filecoin-project/go-commp-utils v0.1.3 // indirect
	github.com/filecoin-project/go-fil-commcid v0.1.0 // indirect
	github.com/ipfs/go-ipfs-util v0.0.2 // indirect
	github.com/ipfs/go-ipld-format v0.0.2 // indirect
	github.com/klauspost/cpuid/v2 v2.2.3 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-base32 v0.0.3 // indirect
	github.com/multiformats/go-base36 v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polydawn/refmt v0.0.0-20201211092308-30ac6d18308e // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	golang.org/x/sys v0.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	lukechampine.com/blake3 v1.1.6 // indirect
)
