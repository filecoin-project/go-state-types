package math

import (
	"math/big"

	math9 "github.com/filecoin-project/go-state-types/builtin/v9/util/math"
)

func Parse(coefs []string) []*big.Int {
	return math9.Parse(coefs)
}
