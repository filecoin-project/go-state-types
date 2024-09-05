package batch

import (
	"errors"

	"github.com/filecoin-project/go-state-types/exitcode"
)

type BatchReturn struct {
	SuccessCount uint64
	FailCodes    []FailCode
}

type FailCode struct {
	Idx  uint64
	Code exitcode.ExitCode
}

func (b BatchReturn) Size() int {
	return int(b.SuccessCount) + len(b.FailCodes)
}

func (b BatchReturn) AllOk() bool {
	return len(b.FailCodes) == 0
}

func (b BatchReturn) Codes() []exitcode.ExitCode {
	codes := make([]exitcode.ExitCode, b.Size())
	i := 0
	for _, fc := range b.FailCodes {
		if fc.Idx > uint64(i) {
			for ; i < int(fc.Idx); i++ {
				codes[i] = exitcode.Ok
			}
		}
		codes[i] = fc.Code
		i++
	}
	for ; i < len(codes); i++ {
		codes[i] = exitcode.Ok
	}
	return codes
}

func (b BatchReturn) CodeAt(n uint64) (exitcode.ExitCode, error) {
	if n >= uint64(b.Size()) {
		return exitcode.Ok, errors.New("index out of bounds")
	}
	for _, fc := range b.FailCodes {
		if fc.Idx == n {
			return fc.Code, nil
		}
		if fc.Idx > n {
			return exitcode.Ok, nil
		}
	}
	return exitcode.Ok, nil
}

func (b BatchReturn) Validate() error {
	size := uint64(b.Size())
	var gaps uint64
	for i, fc := range b.FailCodes {
		if fc.Idx >= size {
			// will also catch the case where the gaps aren't accounted for in total size
			return errors.New("index out of bounds")
		}
		if i > 0 {
			if fc.Idx <= b.FailCodes[i-1].Idx {
				return errors.New("fail codes are not in strictly increasing order")
			}
			gaps += fc.Idx - b.FailCodes[i-1].Idx - 1
		} else {
			gaps += fc.Idx
		}
	}
	return nil
}
