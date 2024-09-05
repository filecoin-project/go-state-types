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
		name   string
		params BatchReturn
		hex    string
	}{
		{
			name:   "empty",
			params: BatchReturn{},
			// [0,[]]
			hex: "820080",
		},
		{
			name:   "single success",
			params: BatchReturn{SuccessCount: 1},
			// [1,[]]
			hex: "820180",
		},
		{
			name:   "single failure",
			params: BatchReturn{FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrIllegalArgument}}},
			// [0,[[0,16]]]
			hex: "820081820010",
		},
		{
			name: "multiple success",
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
		t.Run(tc.name, func(t *testing.T) {
			req := require.New(t)

			var buf bytes.Buffer
			req.NoError(tc.params.MarshalCBOR(&buf))
			req.Equal(tc.hex, hex.EncodeToString(buf.Bytes()))
			var br BatchReturn
			req.NoError(br.UnmarshalCBOR(&buf))
			req.Equal(tc.params, br)
			req.NoError(br.Validate())
		})
	}
}

func TestBatchReturn(t *testing.T) {
	req := require.New(t)

	t.Run("empty", func(t *testing.T) {
		br := BatchReturn{}
		req.Equal(0, br.Size())
		req.True(br.AllOk())
		req.Equal([]exitcode.ExitCode{}, br.Codes())
		_, err := br.CodeAt(0)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})

	t.Run("single success", func(t *testing.T) {
		br := BatchReturn{SuccessCount: 1}
		req.Equal(1, br.Size())
		req.True(br.AllOk())
		req.Equal([]exitcode.ExitCode{exitcode.Ok}, br.Codes())
		ec, err := br.CodeAt(0)
		req.NoError(err)
		req.Equal(exitcode.Ok, ec)
		_, err = br.CodeAt(1)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})

	t.Run("single failure", func(t *testing.T) {
		br := BatchReturn{FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrIllegalArgument}}}
		req.Equal(1, br.Size())
		req.False(br.AllOk())
		req.Equal([]exitcode.ExitCode{exitcode.ErrIllegalArgument}, br.Codes())
		ec, err := br.CodeAt(0)
		req.NoError(err)
		req.Equal(exitcode.ErrIllegalArgument, ec)
		_, err = br.CodeAt(1)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})

	t.Run("multiple success", func(t *testing.T) {
		br := BatchReturn{SuccessCount: 1, FailCodes: []FailCode{{Idx: 1, Code: exitcode.ErrIllegalArgument}}}
		req.Equal(2, br.Size())
		req.False(br.AllOk())
		req.Equal([]exitcode.ExitCode{exitcode.Ok, exitcode.ErrIllegalArgument}, br.Codes())
		ec, err := br.CodeAt(0)
		req.NoError(err)
		req.Equal(exitcode.Ok, ec)
		ec, err = br.CodeAt(1)
		req.NoError(err)
		req.Equal(exitcode.ErrIllegalArgument, ec)
		req.Equal(exitcode.Ok, br.Codes()[0])
		_, err = br.CodeAt(2)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})

	t.Run("multiple failure", func(t *testing.T) {
		br := BatchReturn{SuccessCount: 1, FailCodes: []FailCode{{Idx: 0, Code: exitcode.ErrForbidden}}}
		req.Equal(2, br.Size())
		req.False(br.AllOk())
		req.Equal([]exitcode.ExitCode{exitcode.ErrForbidden, exitcode.Ok}, br.Codes())
		ec, err := br.CodeAt(0)
		req.NoError(err)
		req.Equal(exitcode.ErrForbidden, ec)
		ec, err = br.CodeAt(1)
		req.NoError(err)
		req.Equal(exitcode.Ok, ec)
		_, err = br.CodeAt(2)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})

	t.Run("mixed", func(t *testing.T) {
		br := BatchReturn{SuccessCount: 2, FailCodes: []FailCode{
			{Idx: 1, Code: exitcode.SysErrOutOfGas},
			{Idx: 2, Code: exitcode.ErrIllegalState},
			{Idx: 4, Code: exitcode.ErrIllegalArgument},
		}}
		req.Equal(5, br.Size())
		req.False(br.AllOk())
		req.Equal([]exitcode.ExitCode{exitcode.Ok, exitcode.SysErrOutOfGas, exitcode.ErrIllegalState, exitcode.Ok, exitcode.ErrIllegalArgument}, br.Codes())
		ec, err := br.CodeAt(0)
		req.NoError(err)
		req.Equal(exitcode.Ok, ec)
		ec, err = br.CodeAt(1)
		req.NoError(err)
		req.Equal(exitcode.SysErrOutOfGas, ec)
		ec, err = br.CodeAt(2)
		req.NoError(err)
		req.Equal(exitcode.ErrIllegalState, ec)
		ec, err = br.CodeAt(3)
		req.NoError(err)
		req.Equal(exitcode.Ok, ec)
		ec, err = br.CodeAt(4)
		req.NoError(err)
		req.Equal(exitcode.ErrIllegalArgument, ec)
		_, err = br.CodeAt(5)
		req.Error(err, "index out of bounds")
		req.NoError(br.Validate())
	})
}

func TestBatchReturn_Validate(t *testing.T) {
	tests := []struct {
		name        string
		batchReturn BatchReturn
		errorMsg    string
	}{
		{
			name: "valid batchreturn",
			batchReturn: BatchReturn{
				SuccessCount: 5,
				FailCodes: []FailCode{
					{Idx: 1, Code: exitcode.ErrIllegalArgument},
					{Idx: 3, Code: exitcode.ErrIllegalState},
					{Idx: 6, Code: exitcode.ErrNotPayable},
				},
			},
		},
		{
			name: "failcodes not in strictly increasing order",
			batchReturn: BatchReturn{
				SuccessCount: 5,
				FailCodes: []FailCode{
					{Idx: 1, Code: exitcode.ErrIllegalArgument},
					{Idx: 3, Code: exitcode.ErrIllegalState},
					{Idx: 3, Code: exitcode.ErrNotPayable},
				},
			},
			errorMsg: "fail codes are not in strictly increasing order",
		},
		{
			name: "failcodes contain index out of bounds",
			batchReturn: BatchReturn{
				SuccessCount: 5,
				FailCodes: []FailCode{
					{Idx: 1, Code: exitcode.ErrIllegalArgument},
					{Idx: 3, Code: exitcode.ErrIllegalState},
					{Idx: 10, Code: exitcode.ErrNotPayable},
				},
			},
			errorMsg: "index out of bounds",
		},
		{
			name: "gaps between failures exceed successcount",
			batchReturn: BatchReturn{
				SuccessCount: 2,
				FailCodes: []FailCode{
					{Idx: 1, Code: exitcode.ErrIllegalArgument},
					{Idx: 4, Code: exitcode.ErrIllegalState},
					{Idx: 7, Code: exitcode.ErrNotPayable},
				},
			},
			errorMsg: "index out of bounds",
		},
		{
			name: "initial gap exceeds successcount",
			batchReturn: BatchReturn{
				SuccessCount: 1,
				FailCodes: []FailCode{
					{Idx: 2, Code: exitcode.ErrIllegalArgument},
				},
			},
			errorMsg: "index out of bounds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := require.New(t)
			err := tt.batchReturn.Validate()
			// req.NoError(err)
			if tt.errorMsg != "" {
				req.ErrorContains(err, tt.errorMsg)
			} else {
				req.NoError(err)
			}
		})
	}
}
