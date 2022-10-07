package math

import (
	math9 "github.com/filecoin-project/go-state-types/builtin/v9/util/math"

	"math/big"
)

func ExpNeg(x *big.Int) *big.Int {
	return math9.ExpNeg(x)
}
