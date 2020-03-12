package types

import (
	"gxclient-go/transaction"
	"gxclient-go/util"
)

type TypeDef struct {
	NewTypeName string `json:"new_type_name"`
	Type        string `json:"type"`
}

func (o TypeDef) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.NewTypeName)
	encoder.Encode(o.Type)
	return nil
}

type ErrorMessage struct {
	ErrorCode uint64 `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

func (o ErrorMessage) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.ErrorCode)
	encoder.Encode(o.ErrorMsg)
	return nil
}

type Struct struct {
	Name   string  `json:"name"`
	Base   string  `json:"base"`
	Fields []Field `json:"fields"`
}

func (o Struct) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.Name)
	encoder.Encode(o.Base)
	encoder.EncodeUVarint(uint64(len(o.Fields)))
	for _, f := range o.Fields {
		encoder.Encode(f)
	}
	return nil
}

type Table struct {
	Name      string   `json:"name"`
	IndexType string   `json:"index_type"`
	KeyNames  []string `json:"key_names"`
	KeyTypes  []string `json:"key_types"`
	Type      string   `json:"type"`
}

func (o Table) MarshalTransaction(encoder *transaction.Encoder) error {
	uname := util.StringToName(o.Name)
	encoder.Encode(uname)
	encoder.Encode(o.IndexType)

	encoder.EncodeUVarint(uint64(len(o.KeyNames)))
	for _, keyName := range o.KeyNames {
		encoder.Encode(keyName)
	}

	encoder.EncodeUVarint(uint64(len(o.KeyTypes)))
	for _, keyType := range o.KeyTypes {
		encoder.Encode(keyType)
	}
	encoder.Encode(o.Type)
	return nil
}

type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (o Field) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.Name)
	encoder.Encode(o.Type)
	return nil
}

type Action struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Payable bool   `json:"payable"`
}

func (o Action) MarshalTransaction(encoder *transaction.Encoder) error {
	uname := util.StringToName(o.Name)
	encoder.Encode(uname)
	encoder.Encode(o.Type)
	encoder.Encode(o.Payable)
	return nil
}

type Abi struct {
	Version       string         `json:"version"`
	Types         []TypeDef      `json:"types"`
	Structs       []Struct       `json:"structs"`
	Actions       []Action       `json:"actions"`
	Tables        []Table        `json:"tables"`
	ErrorMessages []ErrorMessage `json:"error_messages"`
	AbiExtensions []interface{}  `json:"abi_extensions"`
}

func (o Abi) MarshalTransaction(encoder *transaction.Encoder) error {
	encoder.Encode(o.Version)
	encoder.EncodeUVarint(uint64(len(o.Types)))
	for _, t := range o.Types {
		encoder.Encode(t)
	}

	encoder.EncodeUVarint(uint64(len(o.Structs)))
	for _, s := range o.Structs {
		encoder.Encode(s)
	}
	encoder.EncodeUVarint(uint64(len(o.Actions)))
	for _, a := range o.Actions {
		encoder.Encode(a)
	}
	encoder.EncodeUVarint(uint64(len(o.Tables)))
	for _, t := range o.Tables {
		encoder.Encode(t)
	}
	encoder.EncodeUVarint(uint64(len(o.ErrorMessages)))
	for _, m := range o.ErrorMessages {
		encoder.Encode(m)
	}
	//encoder.Encode(o.AbiExtensions)
	encoder.EncodeUVarint(0)
	return nil
}

type ContractAccountProperties struct {
	ID                            ObjectID `json:"id"`
	MembershipExpirationDate      string   `json:"membership_expiration_date"`
	Registrar                     string   `json:"registrar"`
	Referrer                      string   `json:"referrer"`
	LifetimeReferrer              string   `json:"lifetime_referrer"`
	NetworkFeePercentage          int64    `json:"network_fee_percentage"`
	LifetimeReferrerFeePercentage int64    `json:"lifetime_referrer_fee_percentage"`
	ReferrerRewardsPercentage     int64    `json:"referrer_rewards_percentage"`
	Name                          string   `json:"name"`
	Statistics                    string   `json:"statistics"`
	WhitelistingAccounts          []string `json:"whitelisting_accounts"`
	BlacklistingAccounts          []string `json:"blacklisting_accounts"`
	WhitelistedAccounts           []string `json:"whitelisted_accounts"`
	BlacklistedAccounts           []string `json:"blacklisted_accounts"`
	XAbi                          Abi      `json:"abi"`
	VmType                        string   `json:"vm_type"`
	VmVersion                     string   `json:"vm_version"`
	Code                          string   `json:"code"`
	CodeVersion                   string   `json:"code_version"`
}
