package batch

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/filecoin-project/go-state-types/exitcode"
	"github.com/stretchr/testify/require"
)

// Tests to match with Rust fil_actors_runtime::serialization
func TestSerializationBatchReturn(t *testing.T) {
	testCases := []struct {
		params BatchReturn
		hex    string
	}{
		{
			params: BatchReturn{},
			// [0,[]]
			hex: "820080",
		},
		{
			params: BatchReturn{SuccessCount: 1},
			// [1,[]]
			hex: "820180",
		},
		{
			params: BatchReturn{FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrIllegalArgument}}},
			// [0,[[0,16]]]
			hex: "820081820010",
		},
		{
			params: BatchReturn{SuccessCount: 2, FailCodes: []FailCode{
				{Idx: 1, Code: exitcode.SysErrOutOfGas},
				{Idx: 2, Code: exitcode.ErrIllegalState},
				{Idx: 4, Code: exitcode.ErrIllegalArgument},
			}},
			// [2,[[1,7],[2,20],[4,16]]]
			hex: "820283820107820214820410",
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var rt BatchReturn
			req.NoError(rt.UnmarshalCBOR(&buf))
			req.Equal(tc.params, rt)
		})
	}
}

func TestBatchReturn(t *testing.T) {
	req := require.New(t)

	br := BatchReturn{}
	req.Equal(0, br.Size())
	req.True(br.AllOk())
	req.Equal([]exitcode.ExitCode{}, br.Codes())

	br = BatchReturn{SuccessCount: 1}
	req.Equal(1, br.Size())
	req.True(br.AllOk())
	req.Equal([]exitcode.ExitCode{exitcode.Ok}, br.Codes())
	req.Equal(exitcode.Ok, br.CodeAt(0))

	br = BatchReturn{FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrIllegalArgument}}}
	req.Equal(1, br.Size())
	req.False(br.AllOk())
	req.Equal([]exitcode.ExitCode{exitcode.ErrIllegalArgument}, br.Codes())
	req.Equal(exitcode.ErrIllegalArgument, br.CodeAt(0))

	br = BatchReturn{SuccessCount: 1, FailCodes: []FailCode{{Idx: 1, Code: exitcode.ErrIllegalArgument}}}
	req.Equal(2, br.Size())
	req.False(br.AllOk())
	req.Equal([]exitcode.ExitCode{exitcode.Ok, exitcode.ErrIllegalArgument}, br.Codes())
	req.Equal(exitcode.Ok, br.CodeAt(0))
	req.Equal(exitcode.ErrIllegalArgument, br.CodeAt(1))

	br = BatchReturn{SuccessCount: 1, FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrForbidden}}}
	req.Equal(2, br.Size())
	req.False(br.AllOk())
	req.Equal([]exitcode.ExitCode{exitcode.ErrForbidden, exitcode.Ok}, br.Codes())
	req.Equal(exitcode.ErrForbidden, br.CodeAt(0))
	req.Equal(exitcode.Ok, br.CodeAt(1))

	br = BatchReturn{SuccessCount: 2, FailCodes: []FailCode{
		{Idx: 1, Code: exitcode.SysErrOutOfGas},
		{Idx: 2, Code: exitcode.ErrIllegalState},
		{Idx: 4, Code: exitcode.ErrIllegalArgument},
	}}
	req.Equal(5, br.Size())
	req.False(br.AllOk())
	req.Equal([]exitcode.ExitCode{exitcode.Ok, exitcode.SysErrOutOfGas, exitcode.ErrIllegalState, exitcode.Ok, exitcode.ErrIllegalArgument}, br.Codes())
	req.Equal(exitcode.Ok, br.CodeAt(0))
	req.Equal(exitcode.SysErrOutOfGas, br.CodeAt(1))
	req.Equal(exitcode.ErrIllegalState, br.CodeAt(2))
	req.Equal(exitcode.Ok, br.CodeAt(3))
	req.Equal(exitcode.ErrIllegalArgument, br.CodeAt(4))
}
