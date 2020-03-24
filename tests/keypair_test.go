package tests

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"gxclient-go/keypair"
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
