package builtin

import (
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateMethodNum(t *testing.T) {

	methodNum, err := GenerateMethodNum("Receive")
	require.NoError(t, err)
	require.Equal(t, methodNum, abi.MethodNum(3726118371))
}
