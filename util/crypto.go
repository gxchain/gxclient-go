package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"github.com/juju/errors"
	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/crypto/ripemd160"
)

func Ripemd160(in []byte) ([]byte, error) {
	h := ripemd160.New()

	if _, err := h.Write(in); err != nil {
		return nil, errors.Annotate(err, "Write")
	}

	sum := h.Sum(nil)
	return sum, nil
}

func Ripemd160Checksum(in []byte) ([]byte, error) {
	buf, err := Ripemd160(in)
	if err != nil {
		return nil, errors.Annotate(err, "Ripemd160")
	}

	return buf[:4], nil
}
func Sha512Checksum(in []byte) ([]byte, error) {
	buf := sha512.Sum512(in)
	return buf[:4], nil
}

func Sha256(in []byte) []byte {
	buf := sha256.Sum256(in)
	return buf[:]
}

func CharToSymbol(c byte) uint64 {
	if c >= 'a' && c <= 'z' {
		return uint64((c - 'a') + 6)
	}
	if c >= '1' && c <= '5' {
		return uint64((c - '1') + 1)
	}
	return 0
}

func StringToName(str string) uint64 {
	bs := []byte(str)
	var name uint64 = 0

	len := len(bs)
	if len <= 0 {
		return 0
	}

	i := 0
	for ; i < len && i < 12; i++ {
		name |= (CharToSymbol(bs[i]) & 0x1f) << uint(64-5*(i+1))
	}

	if i == 12 && len > 12 {
		name |= CharToSymbol(bs[12]) & 0x0F
	}

	return name
}

func ToBytes(in interface{}) []byte {
	b, err := ffjson.Marshal(in)
	if err != nil {
		panic("ToBytes: unable to marshal input")
	}
	return b
}

func ToMap(in interface{}) map[string]interface{} {
	b, err := ffjson.Marshal(in)
	if err != nil {
	}

	m := make(map[string]interface{})
	if err := ffjson.Unmarshal(b, &m); err != nil {
		panic("ToMap: unable to unmarshal input")
	}

	return nil
}
