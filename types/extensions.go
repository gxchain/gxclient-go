package types

import (
	"github.com/juju/errors"
	"gxclient-go/transaction"
)

type Extensions []interface{}

func (p Extensions) MarshalTransaction(enc *transaction.Encoder) error {
	if err := enc.EncodeUVarint(uint64(len(p))); err != nil {
		return errors.Annotate(err, "encode length")
	}

	for _, ex := range p {
		if err := enc.Encode(ex); err != nil {
			return errors.Annotate(err, "encode Extension")
		}
	}

	return nil
}
