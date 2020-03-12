package types

import (
	"encoding/json"
	"gxclient-go/transaction"
)

// NewStakingUpdateOperation returns a new instance of StakingUpdateOperation
func NewStakingUpdateOperation(owner, trustNode, stakingId ObjectID, fee AssetAmount) *StakingUpdateOperation {
	op := &StakingUpdateOperation{
		Owner:      owner,
		TrustNode:  trustNode,
		StakingId:  stakingId,
		Fee:        fee,
		Extensions: []json.RawMessage{},
	}
	return op
}

// StakingUpdateOperation
type StakingUpdateOperation struct {
	Fee        AssetAmount       `json:"fee"`
	Owner      ObjectID          `json:"owner"`
	TrustNode  ObjectID          `json:"trust_node"`
	StakingId  ObjectID          `json:"staking_id"`
	Extensions []json.RawMessage `json:"extensions"`
}

func (op *StakingUpdateOperation) Type() OpType { return StakingUpdateOpType }

func (op *StakingUpdateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Owner)
	enc.Encode(op.TrustNode)
	enc.Encode(op.StakingId)

	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}
