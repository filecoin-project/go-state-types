package verifreg

import (
	"github.com/multiformats/go-varint"
)

type AllocationId uint64

func (a AllocationId) Key() string {
	return string(varint.ToUvarint(uint64(a)))
}

type ClaimId uint64

func (a ClaimId) Key() string {
	return string(varint.ToUvarint(uint64(a)))
}
