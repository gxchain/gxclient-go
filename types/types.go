package types

import (
	"fmt"
	"gxclient-go/transaction"
)

var (
	ErrRPCClientNotInitialized      = fmt.Errorf("RPC client is not initialized")
	ErrNotImplemented               = fmt.Errorf("not implemented")
	ErrInvalidInputType             = fmt.Errorf("invalid input type")
	ErrInvalidInputLength           = fmt.Errorf("invalid input length")
	ErrInvalidPublicKey             = fmt.Errorf("invalid PublicKey")
	ErrInvalidAddress               = fmt.Errorf("invalid Address")
	ErrPublicKeyChainPrefixMismatch = fmt.Errorf("PublicKey database prefix mismatch")
	ErrAddressChainPrefixMismatch   = fmt.Errorf("Address database prefix mismatch")
	ErrInvalidChecksum              = fmt.Errorf("invalid checksum")
	ErrNoSigningKeyFound            = fmt.Errorf("no signing key found")
	ErrNoVerifyingKeyFound          = fmt.Errorf("no verifying key found")
	ErrInvalidDigestLength          = fmt.Errorf("invalid digest length")
	ErrInvalidPrivateKeyCurve       = fmt.Errorf("invalid PrivateKey curve")
	ErrCurrentChainConfigIsNotSet   = fmt.Errorf("current database config is not set")
)

type Int8 int8
type UInt8 uint8
type UInt16 uint16
type UInt32 uint32
type UInt64 uint64
type Int64 int64

type AssetType Int8

func (num UInt8) MarshalTransaction(enc *transaction.Encoder) error {
	return enc.EncodeNumber(uint8(num))
}

func (num UInt16) MarshalTransaction(enc *transaction.Encoder) error {
	return enc.EncodeNumber(uint16(num))
}

func (num UInt32) MarshalTransaction(enc *transaction.Encoder) error {
	return enc.EncodeNumber(uint32(num))
}

func (num UInt64) MarshalTransaction(enc *transaction.Encoder) error {
	return enc.EncodeNumber(uint64(num))
}

type WorkerInitializerType UInt8

const (
	WorkerInitializerTypeRefund WorkerInitializerType = iota
	WorkerInitializerTypeVestingBalance
	WorkerInitializerTypeBurn
)

const (
	AssetTypeUndefined AssetType = -1
	AssetTypeCoreAsset AssetType = iota
	AssetTypeUIA
	AssetTypeSmartCoin
	AssetTypePredictionMarket
)

type SpaceType Int8

const (
	SpaceTypeUndefined SpaceType = -1
	SpaceTypeProtocol  SpaceType = iota
	SpaceTypeImplementation
)

type PredicateType UInt8

const (
	PredicateAccountNameEqLit PredicateType = iota
	PredicateAssetSymbolEqLit
	PredicateBlockId
)

type ObjectType Int8

const (
	ObjectTypeUndefined ObjectType = -1
)

//for SpaceTypeProtocol
const (
	ObjectTypeBase ObjectType = iota + 1
	ObjectTypeAccount
	ObjectTypeAsset
	ObjectTypeForceSettlement
	ObjectTypeCommiteeMember
	ObjectTypeWitness
	ObjectTypeLimitOrder
	ObjectTypeCallOrder
	ObjectTypeCustom
	ObjectTypeProposal
	ObjectTypeOperationHistory
	ObjectTypeWithdrawPermission
	ObjectTypeVestingBalance
	ObjectTypeWorker
	ObjectTypeBalance
)

// for SpaceTypeImplementation
const (
	ObjectTypeGlobalProperty ObjectType = iota + 1
	ObjectTypeDynamicGlobalProperty
	ObjectTypeAssetDynamicData
	ObjectTypeAssetBitAssetData
	ObjectTypeAccountBalance
	ObjectTypeAccountStatistics
	ObjectTypeTransaction
	ObjectTypeBlockSummary
	ObjectTypeAccountTransactionHistory
	ObjectTypeBlindedBalance
	ObjectTypeChainProperty
	ObjectTypeWitnessSchedule
	ObjectTypeBudgetRecord
	ObjectTypeSpecialAuthority
)

type AccountCreateExtensionsType UInt8

const (
	AccountCreateExtensionsNullExt AccountCreateExtensionsType = iota
	AccountCreateExtensionsOwnerSpecial
	AccountCreateExtensionsActiveSpecial
	AccountCreateExtensionsBuyback
)

type SpecialAuthorityType UInt8

const (
	SpecialAuthorityTypeNoSpecial SpecialAuthorityType = iota
	SpecialAuthorityTypeTopHolders
)

type VestingPolicyType UInt8

const (
	VestingPolicyTypeLinear VestingPolicyType = iota
	VestingPolicyTypeCCD
)

type FeeParametersType UInt8

const (
	FeeParametersTypeTransfer FeeParametersType = iota
	FeeParametersTypeLimitOrderCreate
	FeeParametersTypeLimitOrderCancel
	FeeParametersTypeCallOrderUpdate
	FeeParametersTypeFillOrder
	FeeParametersTypeAccountCreate
	FeeParametersTypeAccountUpdate
	FeeParametersTypeAccountWhitelist
	FeeParametersTypeAccountUpgrade
	FeeParametersTypeAccountTransfer
	FeeParametersTypeAssetCreate
	FeeParametersTypeAssetUpdate
	FeeParametersTypeAssetUpdateBitasset
	FeeParametersTypeAssetUpdateFeedProducers
	FeeParametersTypeAssetIssue
	FeeParametersTypeAssetReserve
	FeeParametersTypeAssetFundFeePool
	FeeParametersTypeAssetSettle
	FeeParametersTypeAssetGlobalSettle
	FeeParametersTypeAssetPublishFeed
	FeeParametersTypeWitnessCreate
	FeeParametersTypeWitnessUpdate
	FeeParametersTypeProposalCreate
	FeeParametersTypeProposalUpdate
	FeeParametersTypeProposalDelete
	FeeParametersTypeWithdrawPermissionCreate
	FeeParametersTypeWithdrawPermissionUpdate
	FeeParametersTypeWithdrawPermissionClaim
	FeeParametersTypeWithdrawPermissionDelete
	FeeParametersTypeCommitteeMemberCreate
	FeeParametersTypeCommitteeMemberUpdate
	FeeParametersTypeCommitteeMemberUpdateGlobalParameters
	FeeParametersTypeVestingBalanceCreate
	FeeParametersTypeVestingBalanceWithdraw
	FeeParametersTypeWorkerCreate
	FeeParametersTypeCustom
	FeeParametersTypeAssert
	FeeParametersTypeBalanceClaim
	FeeParametersTypeOverrideTransfer
	FeeParametersTypeTransferToBlind
	FeeParametersTypeBlindTransfer
	FeeParametersTypeTransferFromBlind
	FeeParametersTypeAssetSettleCancel
	FeeParametersTypeAssetClaimFees
)
