package builtin

import (
	"github.com/filecoin-project/go-state-types/builtin/frc0042"
	"testing"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/stretchr/testify/require"
)

func TestGenerateMethodNum(t *testing.T) {

	methodNum, err := frc0042.GenerateMethodNum("Receive")
	require.NoError(t, err)
	require.Equal(t, methodNum, abi.MethodNum(3726118371))
}
