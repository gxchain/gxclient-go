package types

import (
	"encoding/json"
	"gxclient-go/transaction"
)

// NewTransferOperation returns a new instance of TransferOperation
func NewTransferOperation(from, to ObjectID, amount, fee AssetAmount, memo *Memo) *TransferOperation {
	op := &TransferOperation{
		From:       from,
		To:         to,
		Amount:     amount,
		Fee:        fee,
		Memo:       memo,
		Extensions: []json.RawMessage{},
	}

	return op
}

// TransferOperation
type TransferOperation struct {
	From       ObjectID          `json:"from"`
	To         ObjectID          `json:"to"`
	Amount     AssetAmount       `json:"amount"`
	Fee        AssetAmount       `json:"fee"`
	Memo       *Memo             `json:"memo,omitempty"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (op *TransferOperation) Type() OpType { return TransferOpType }

func (op *TransferOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.From)
	enc.Encode(op.To)
	enc.Encode(op.Amount)

	if op.Memo != nil && op.Memo.Message.Length() > 0 {
		enc.EncodeUVarint(1)
		enc.Encode(op.Memo)
	} else {
		//Memo?
		enc.EncodeUVarint(0)
	}
	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}
