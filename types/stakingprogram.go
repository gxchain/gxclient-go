package types

type StakiongPrograms []StakingProgram

type StakingProgram struct {
	ProgramId   string `json:"program_id"`
	Weight      uint32 `json:"weight"`
	StakingDays uint32 `json:"staking_days"`
}
