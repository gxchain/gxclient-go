package types

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"github.com/juju/errors"
	"gxclient-go/transaction"
	"math/rand"
	"strconv"
	"time"
)

type Memo struct {
	From    PublicKey `json:"from"`
	To      PublicKey `json:"to"`
	Nonce   UInt64    `json:"nonce"`
	Message Buffer    `json:"message"`
}

func (p Memo) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.Encode(p.From); err != nil {
		return errors.Annotate(err, "encode from")
	}

	if err := enc.Encode(p.To); err != nil {
		return errors.Annotate(err, "encode to")
	}

	if err := enc.Encode(p.Nonce); err != nil {
		return errors.Annotate(err, "encode nonce")
	}

	if err := enc.Encode(p.Message); err != nil {
		return errors.Annotate(err, "encode Message")
	}

	return nil
}

func (m *Memo) UnmarshalJSON(b []byte) (err error) {
	stringCase := struct {
		From    PublicKey `json:"from"`
		To      PublicKey `json:"to"`
		Nonce   string    `json:"nonce"`
		Message Buffer    `json:"message"`
	}{}

	uint64Case := struct {
		From    PublicKey `json:"from"`
		To      PublicKey `json:"to"`
		Nonce   UInt64    `json:"nonce"`
		Message Buffer    `json:"message"`
	}{}

	if err = json.Unmarshal(b, &uint64Case); err == nil {
		m.From = uint64Case.From
		m.To = uint64Case.To
		m.Nonce = uint64Case.Nonce
		m.Message = uint64Case.Message
		return nil
	}

	if err = json.Unmarshal(b, &stringCase); err == nil {
		u, err := strconv.ParseUint(stringCase.Nonce, 10, 64)
		m.From = stringCase.From
		m.To = stringCase.To
		m.Nonce = UInt64(u)
		m.Message = stringCase.Message
		return err
	}

	return err
}

//Encrypt calculates a shared secret by the senders private key
//and the recipients public key, then encrypts the given memo message.
func (p *Memo) Encrypt(priv *PrivateKey, msg string) error {
	sec, err := priv.SharedSecret(&p.To, 16, 16)
	if err != nil {
		return errors.Annotate(err, "SharedSecret")
	}

	iv, blk, err := p.cypherBlock(sec)
	if err != nil {
		return errors.Annotate(err, "cypherBlock")
	}

	buf := []byte(msg)
	digest := sha256.Sum256(buf)
	mode := cipher.NewCBCEncrypter(blk, iv)

	// checksum + msg
	raw := digest[:4]
	raw = append(raw, buf...)

	if len(raw)%16 != 0 {
		raw = pad(raw, 16)
	}

	dst := make([]byte, len(raw))
	mode.CryptBlocks(dst, raw)
	p.Message = dst

	return nil
}

//Decrypt calculates a shared secret by the receivers private key
//and the senders public key, then decrypts the given memo message.
func (p Memo) Decrypt(priv *PrivateKey) (string, error) {
	var counterPartyPubKey PublicKey
	myPubKey := priv.PublicKey()

	if myPubKey.Equal(&p.To) {
		counterPartyPubKey = p.From
	} else if myPubKey.Equal(&p.From) {
		counterPartyPubKey = p.To
	} else {
		return "", errors.New("Invalid counterparty public key, cannot decrypt")
	}

	sec, err := priv.SharedSecret(&counterPartyPubKey, 16, 16)
	if err != nil {
		return "", errors.Annotate(err, "SharedSecret")
	}

	iv, blk, err := p.cypherBlock(sec)
	if err != nil {
		return "", errors.Annotate(err, "cypherBlock")
	}

	mode := cipher.NewCBCDecrypter(blk, iv)
	dst := make([]byte, len(p.Message))
	mode.CryptBlocks(dst, p.Message)

	//verify checksum
	chk1 := dst[:4]
	msg := unpad(dst[4:])
	dig := sha256.Sum256(msg)
	chk2 := dig[:4]

	if bytes.Compare(chk1, chk2) != 0 {
		return "", ErrInvalidChecksum
	}

	return string(msg), nil
}

func (p Memo) cypherBlock(sec []byte) ([]byte, cipher.Block, error) {
	ss := sha512.Sum512(sec)

	var seed []byte
	seed = append(seed, []byte(strconv.FormatUint(uint64(p.Nonce), 10))...)
	seed = append(seed, []byte(hex.EncodeToString(ss[:]))...)

	sd := sha512.Sum512(seed)
	blk, err := aes.NewCipher(sd[0:32])
	if err != nil {
		return nil, nil, errors.Annotate(err, "NewCipher")
	}

	return sd[32:48], blk, nil
}

func unpad(buf []byte) []byte {
	b := buf[len(buf)-1:][0]
	cnt := int(b)
	l := len(buf) - cnt

	a := bytes.Repeat([]byte{b}, cnt)
	if bytes.Compare(a, buf[l:]) == 0 {
		return buf[:l]
	}

	return buf
}

func pad(buf []byte, length int) []byte {
	cnt := length - len(buf)%length
	buf = append(buf, bytes.Repeat([]byte{byte(cnt)}, cnt)...)
	return buf
}

func GetNonce() UInt64 {
	rand.Seed(time.Now().Unix())
	return UInt64(rand.Uint64())
}
