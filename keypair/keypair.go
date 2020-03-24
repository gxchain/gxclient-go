package keypair

import (
	"github.com/tyler-smith/go-bip39"
	"gxclient-go/types"
)

type KeyPair struct {
	BrainKey   string
	PrivateKey *types.PrivateKey
}

//Generates the key pair
func GenerateKeyPair(brainKey string) (*KeyPair, error) {
	if len(brainKey) == 0 {
		// Generate a mnemonic for memorization or user-friendly seeds
		entropy, err := bip39.NewEntropy(128)
		if err != nil {
			return nil, err
		}
		if brainKey, err = bip39.NewMnemonic(entropy); err != nil {
			return nil, err
		}
	}

	pri, err := types.NewPrivateKeyFromBrainKey(brainKey, "0")
	if err != nil {
		return nil, err
	}
	return &KeyPair{
		BrainKey:   brainKey,
		PrivateKey: pri,
	}, nil
}

//Export public key from private key
func PrivateToPublic(priWif string) (string, error) {
	pri, err := types.NewPrivateKeyFromWif(priWif)
	if err != nil {
		return "", err
	}
	return pri.PublicKey().String(), nil
}

//Check if privateKey is valid or not
func IsValidPrivate(priWif string) bool {
	_, err := types.NewPrivateKeyFromWif(priWif)
	if err != nil {
		return false
	}
	return true
}

//Check if publicKey is valid or not
func IsValidPublic(priWif string) bool {
	_, err := types.NewPublicKeyFromString(priWif)
	if err != nil {
		return false
	}
	return true
}
