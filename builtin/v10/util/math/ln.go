package math

import (
	math9 "github.com/filecoin-project/go-state-types/builtin/v9/util/math"

	"github.com/filecoin-project/go-state-types/big"
)

func Ln(z big.Int) big.Int {
	return math9.Ln(z)
}
