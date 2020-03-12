package types

import (
	"encoding/json"
	"gxclient-go/transaction"
)

// NewStakingCreateOperation returns a new instance of StakingCreateOperation
func NewStakingCreateOperation(owner ObjectID, trustNode ObjectID, amount, fee AssetAmount, programId string, weight, stakingDays uint32) *StakingCreateOperation {
	op := &StakingCreateOperation{
		Owner:       owner,
		TrustNode:   trustNode,
		Amount:      amount,
		Fee:         fee,
		ProgramId:   programId,
		Weight:      weight,
		StakingDays: stakingDays,
		Extensions:  []json.RawMessage{},
	}
	return op
}

// StakingCreateOperation
type StakingCreateOperation struct {
	Fee         AssetAmount       `json:"fee"`
	Owner       ObjectID          `json:"owner"`
	TrustNode   ObjectID          `json:"trust_node"`
	Amount      AssetAmount       `json:"amount"`
	ProgramId   string            `json:"program_id"`
	Weight      uint32            `json:"weight"`
	StakingDays uint32            `json:"staking_days"`
	Extensions  []json.RawMessage `json:"extensions"`
}

func (op *StakingCreateOperation) Type() OpType { return StakingCreateOpType }

func (op *StakingCreateOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(op.Type()))
	enc.Encode(op.Fee)
	enc.Encode(op.Owner)
	enc.Encode(op.TrustNode)
	enc.Encode(op.Amount)
	enc.Encode(op.ProgramId)
	enc.Encode(op.Weight)
	enc.Encode(op.StakingDays)

	//Extensions
	enc.EncodeUVarint(0)
	return enc.Err()
}
