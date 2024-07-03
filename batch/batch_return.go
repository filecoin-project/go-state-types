package batch

import "github.com/filecoin-project/go-state-types/exitcode"

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

func (b BatchReturn) CodeAt(n uint64) exitcode.ExitCode {
	for _, fc := range b.FailCodes {
		if fc.Idx == n {
			return fc.Code
		}
	}
	return exitcode.Ok
}
