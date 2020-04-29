package offline

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcutil"
	"github.com/pkg/errors"
	"gxclient-go/sign"
)

func OfflineSign(activePriWif, unSignedHex, chainId string) (string, error) {
	var msgBuffer bytes.Buffer
	// Write the chain ID.
	rawChainID, err := hex.DecodeString(chainId)
	if err != nil {
		return "", errors.Wrapf(err, "failed to decode chain ID: %v", chainId)
	}
	if _, err := msgBuffer.Write(rawChainID); err != nil {
		return "", errors.Wrap(err, "failed to write chain ID")
	}

	rawTx, err := hex.DecodeString(unSignedHex[0 : len(unSignedHex)-2])
	if err != nil {
		return "", err
	}

	if _, err := msgBuffer.Write(rawTx); err != nil {
		return "", errors.Wrap(err, "failed to write serialized transaction")
	}

	msgBytes := msgBuffer.Bytes()
	// Compute the digest.
	digest := sha256.Sum256(msgBytes)

	w, err := btcutil.DecodeWIF(activePriWif)
	if err != nil {
		return "", err
	}
	// Set the signature array in the transaction.
	sig := sign.SignBufferSha256(digest[:], w.PrivKey.ToECDSA())
	signature := hex.EncodeToString(sig)
	return signature, nil
}
