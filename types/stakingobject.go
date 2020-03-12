package types

type StakingObjects []StakingObject

type StakingObject struct {
	ID             ObjectID    `json:"id"`
	Owner          ObjectID    `json:owner`
	TrustNode      ObjectID    `json:"trust_node"`
	Amount         AssetAmount `json:"amount"`
	CreateDateTime Time        `json:"create_date_time"`
	ProgramId      string      `json:"program_id"`
	StakingDays    uint32      `json:"staking_days"`
	Weight         uint32      `json:"weight"`
	IsValid        bool        `json:"is_valid"`
}
