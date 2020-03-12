package types

import (
	"encoding/json"
	"gxclient-go/transaction"
)

// NewStakingClaimOperation returns a new instance of StakingClaimOperation
func NewStakingClaimOperation(owner, stakingId ObjectID, fee AssetAmount) *StakingClaimOperation {
	op := &StakingClaimOperation{
		Owner:      owner,
		StakingId:  stakingId,
		Fee:        fee,
		Extensions: []json.RawMessage{},
	}
	return op
}

// StakingClaimOperation
type StakingClaimOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Owner      ObjectID          `json:"owner"`
	StakingId  ObjectID          `json:"staking_id"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (op *StakingClaimOperation) Type() OpType { return StakingClaimOpType }

func (op *StakingClaimOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Owner)
	enc.Encode(op.StakingId)

	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}
