package builtin

import (
	"encoding/binary"
	"unicode"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/minio/blake2b-simd"
	"golang.org/x/xerrors"
)

// Generates a standard FRC-42 compliant method number
// Reference: https://github.com/filecoin-project/FIPs/blob/master/FRCs/frc-0042.md
func GenerateFRCMethodNum(name string) (abi.MethodNum, error) {
	err := validateMethodName(name)
	if err != nil {
		return 0, err
	}

	if name == "Constructor" {
		return MethodConstructor, nil
	}

	digest := blake2b.Sum512([]byte("1|" + name))

	for i := 0; i < 64; i += 4 {
		methodId := binary.BigEndian.Uint32(digest[i : i+4])
		if methodId >= (1 << 24) {
			return abi.MethodNum(methodId), nil
		}
	}

	return abi.MethodNum(0), xerrors.Errorf("Could not generate method num from method name %s:", name)
}

func validateMethodName(name string) error {
	if name == "" {
		return xerrors.Errorf("empty name string")
	}

	if !(unicode.IsUpper(rune(name[0])) || name[0] == "_"[0]) {
		return xerrors.Errorf("Method name first letter must be uppercase or underscore, method name: %s", name)
	}

	for _, c := range name {
		if !(unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_') {
			return xerrors.Errorf("method name has illegal characters, method name: %s", name)
		}
	}

	return nil
}

func MustGenerateFRCMethodNum(name string) abi.MethodNum {
	methodNum, err := GenerateFRCMethodNum(name)
	if err != nil {
		panic(err)
	}
	return methodNum
}
