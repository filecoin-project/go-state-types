package abi_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
)

func newIDAddr(t testing.TB, id uint64) address.Address {
	address, err := address.NewIDAddress(id)
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func newActorAddr(t testing.TB, data string) address.Address {
	address, err := address.NewActorAddress([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	return address
}

func TestAddrKey(t *testing.T) {
	id_address_1 := newIDAddr(t, 101)
	id_address_2 := newIDAddr(t, 102)
	actor_address_1 := newActorAddr(t, "actor1")
	actor_address_2 := newActorAddr(t, "222")

	t.Run("address to key string conversion", func(t *testing.T) {
		assert.Equal(t, "\x00\x65", abi.AddrKey(id_address_1).Key())
		assert.Equal(t, "\x00\x66", abi.AddrKey(id_address_2).Key())
		assert.Equal(t, "\x02\x58\xbe\x4f\xd7\x75\xa0\xc8\xcd\x9a\xed\x86\x4e\x73\xab\xb1\x86\x46\x5f\xef\xe1", abi.AddrKey(actor_address_1).Key())
		assert.Equal(t, "\x02\xaa\xd0\xb2\x98\xa9\xde\xab\xbb\xb6\u007f\x80\x5f\x66\xaa\x68\x8c\xdd\x89\xad\xf5", abi.AddrKey(actor_address_2).Key())
	})
}
