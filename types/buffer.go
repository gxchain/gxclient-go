package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	"gxclient-go/transaction"
	"io"
)

type Buffer []byte
type Buffers []Buffer

func (p *Buffer) UnmarshalJSON(data []byte) error {
	var b string
	if err := ffjson.Unmarshal(data, &b); err != nil {
		return errors.Annotate(err, "Unmarshal")
	}

	return p.FromString(b)
}

func (p Buffer) Bytes() []byte {
	return p
}

func (p Buffer) Length() int {
	return len(p)
}

func (p Buffer) String() string {
	return hex.EncodeToString(p)
}

func (p *Buffer) FromString(data string) error {
	buf, err := hex.DecodeString(data)
	if err != nil {
		return errors.Annotate(err, "DecodeString")
	}

	*p = buf
	return nil
}

func (p Buffer) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(p.String())
}

func (p Buffer) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	if err := enc.Encode(p.Bytes()); err != nil {
		return errors.Annotate(err, "encode bytes")
	}

	return nil
}

//Encrypt AES-encrypts the buffer content
func (p *Buffer) Encrypt(cipherKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, errors.Annotate(err, "NewCipher")
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+p.Length())
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, errors.Annotate(err, "ReadFull")
	}

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(
		ciphertext[aes.BlockSize:],
		p.Bytes(),
	)

	return ciphertext, nil
}

//Decrypt AES decrypts the buffer content
func (p *Buffer) Decrypt(cipherKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(cipherKey)
	if err != nil {
		return nil, errors.Annotate(err, "NewCipher")
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if byteLen := p.Length(); byteLen < aes.BlockSize {
		return nil, errors.Errorf("invalid cipher size %d", byteLen)
	}

	buf := p.Bytes()
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	// XORKeyStream can work in-place if the two arguments are the same.
	cipher.NewCFBDecrypter(block, iv).XORKeyStream(buf, buf)

	return buf, nil
}

func BufferFromString(data string) (b Buffer, err error) {
	b = Buffer{}
	err = b.FromString(data)
	return
}
