package types

import (
	"gxclient-go/transaction"
)

func NewAccountCreateOperation(registrar, referrer GrapheneID, referrerPercent UInt16, owner, active Authority, name string, options AccountOptions) *AccountCreateOperation {
	op := &AccountCreateOperation{
		Registrar:       registrar,
		Referrer:        referrer,
		ReferrerPercent: referrerPercent,
		Owner:           owner,
		Active:          active,
		Name:            name,
		Options:         options,
		Extensions:      AccountCreateExtensions{},
	}

	return op
}

type AccountCreateOperation struct {
	OperationFee
	Registrar       GrapheneID              `json:"registrar"`
	Referrer        GrapheneID              `json:"referrer"`
	ReferrerPercent UInt16                  `json:"referrer_percent"`
	Owner           Authority               `json:"owner"`
	Active          Authority               `json:"active"`
	Name            string                  `json:"name"`
	Options         AccountOptions          `json:"options"`
	Extensions      AccountCreateExtensions `json:"extensions"`
}

func (p AccountCreateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.Encode(int8(p.Type()))
	enc.Encode(p.Fee)
	enc.Encode(p.Registrar)
	enc.Encode(p.Referrer)
	enc.Encode(p.ReferrerPercent)
	enc.Encode(p.Name)
	enc.Encode(p.Owner)
	enc.Encode(p.Active)
	enc.Encode(p.Options)
	enc.EncodeUVarint(0)

	return enc.Err()
}

func (op *AccountCreateOperation) Type() OpType { return AccountCreateOpType }
