package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gxclient-go/faucet"
	"testing"
)

func TestApi_GetRegister(t *testing.T) {
	transaction, err := faucet.Register(testFaucet, "cli-wallet-test-9", testPub, testPub, testPub)
	require.Nil(t, err)
	str, _ := json.Marshal(transaction)
	fmt.Println(string(str))
}
