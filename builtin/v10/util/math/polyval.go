package math

import (
	"math/big"

	math9 "github.com/filecoin-project/go-state-types/builtin/v9/util/math"
)

const Precision128 = math9.Precision128

func Polyval(p []*big.Int, x *big.Int) *big.Int {
	return math9.Polyval(p, x)
}
