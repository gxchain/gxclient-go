package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gxclient-go/keypair"
	"gxclient-go/types"
	"testing"
)

func Test_GenerateKeyPair(t *testing.T) {
	keyPair, err := keypair.GenerateKeyPair("")
	require.Nil(t, err)
	fmt.Println(keyPair.BrainKey)
	fmt.Println(keyPair.PrivateKey.ToWIF())
	fmt.Println(keyPair.PrivateKey.PublicKey().String())
}

func Test_PrivateToPublic(t *testing.T) {
	pub, err := keypair.PrivateToPublic(testPri)
	require.Nil(t, err)
	fmt.Println(pub)
}

func Test_IsValidPrivate(t *testing.T) {
	bool := keypair.IsValidPrivate(testPri)
	fmt.Println(bool)
}

func Test_IsValidPublic(t *testing.T) {
	bool := keypair.IsValidPublic(testPub)
	fmt.Println(bool)
}

func Test_Simple(t *testing.T) {
	pub1, _ := types.NewPublicKeyFromString("GXC89s7yudddNUAYkYXADHkknJzBHJoQaJSjw8TUGKYFCeqrUFcyY")
	pub2, _ := types.NewPublicKeyFromString("GXC1111111111111111111111111111111114T1Anm")
	fmt.Println(pub1)
	fmt.Println(pub2)
}
