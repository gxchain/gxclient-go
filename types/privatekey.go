package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"gxclient-go/transaction"
	"io"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/juju/errors"
)

var (
	ErrInvalidCurve               = fmt.Errorf("invalid elliptic curve")
	ErrSharedKeyTooBig            = fmt.Errorf("shared key params are too big")
	ErrSharedKeyIsPointAtInfinity = fmt.Errorf("shared key is point at infinity")
)

type PrivateKeys []PrivateKey

type PrivateKey struct {
	priv *btcec.PrivateKey
	pub  *PublicKey
	raw  []byte
}

func NewPrivateKeyFromBrainKey(brainKey, seq string) (*PrivateKey, error) {
	sha512 := sha512.Sum512([]byte(brainKey + " " + seq))
	hashByte := sha256.Sum256(sha512[:])
	privateKey, err := NewDeterministicPrivateKey(bytes.NewBuffer(hashByte[:]))
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func NewDeterministicPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	return newRandomPrivateKey(randSource)
}

func newRandomPrivateKey(randSource io.Reader) (*PrivateKey, error) {
	rawPrivKey := make([]byte, 32)
	written, err := io.ReadFull(randSource, rawPrivKey)
	if err != nil {
		return nil, fmt.Errorf("error feeding crypto-rand numbers to seed ephemeral private key: %s", err)
	}
	if written != 32 {
		return nil, fmt.Errorf("couldn't write 32 bytes of randomness to seed ephemeral private key")
	}

	privKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), rawPrivKey)

	pub, err := NewPublicKey(privKey.PubKey())
	if err != nil {
		return nil, errors.Annotate(err, "NewPublicKey")
	}

	raw := append([]byte{128}, privKey.D.Bytes()...)
	raw = append(raw, checksum(raw)...)

	return &PrivateKey{
		priv: privKey,
		pub:  pub,
		raw:  raw,
	}, nil
}

func checksum(data []byte) []byte {
	c1 := sha256.Sum256(data)
	c2 := sha256.Sum256(c1[:])
	return c2[0:4]
}

func (p PrivateKey) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.EncodeUVarint(uint64(len(p.raw))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	if err := enc.Encode(p.raw); err != nil {
		return errors.Annotate(err, "encode raw")
	}

	return nil
}

func NewPrivateKeyFromWif(wifPrivateKey string) (*PrivateKey, error) {
	w, err := btcutil.DecodeWIF(wifPrivateKey)
	if err != nil {
		return nil, errors.Annotate(err, "DecodeWIF")
	}

	priv := w.PrivKey
	raw := base58.Decode(wifPrivateKey)
	pub, err := NewPublicKey(priv.PubKey())
	if err != nil {
		return nil, errors.Annotate(err, "NewPublicKey")
	}

	k := PrivateKey{
		priv: priv,
		raw:  raw,
		pub:  pub,
	}

	return &k, nil
}

func (p PrivateKey) PublicKey() *PublicKey {
	return p.pub
}

func (p PrivateKey) ECPrivateKey() *btcec.PrivateKey {
	return p.priv
}

func (p PrivateKey) ToECDSA() *ecdsa.PrivateKey {
	return p.priv.ToECDSA()
}

func (p PrivateKey) Bytes() []byte {
	return p.priv.Serialize()
}

func (p PrivateKey) ToHex() string {
	return hex.EncodeToString(p.Bytes())
}

func (p PrivateKey) ToWIF() string {
	return base58.Encode(p.raw)
}

func (p PrivateKey) SignCompact(hash []byte) (sig []byte, err error) {
	sig, err = btcec.SignCompact(btcec.S256(), p.ECPrivateKey(), hash, true)
	return
}

func (p PrivateKey) SharedSecret(pub *PublicKey, skLen, macLen int) (sk []byte, err error) {
	puk := pub.ToECDSA()
	pvk := p.priv

	if pvk.PublicKey.Curve != puk.Curve {
		return nil, ErrInvalidCurve
	}

	if skLen+macLen > pub.MaxSharedKeyLength() {
		return nil, ErrSharedKeyTooBig
	}

	x, _ := puk.Curve.ScalarMult(puk.X, puk.Y, pvk.D.Bytes())
	if x == nil {
		return nil, ErrSharedKeyIsPointAtInfinity
	}

	sk = make([]byte, skLen+macLen)
	skBytes := x.Bytes()
	copy(sk[len(sk)-len(skBytes):], skBytes)
	return sk, nil
}
