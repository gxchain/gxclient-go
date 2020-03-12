package database

import (
	"encoding/json"
	"gxclient-go/types"
)

type Asset struct {
	ID                 types.ObjectID `json:"id"`
	Symbol             string         `json:"symbol"`
	Precision          uint8          `json:"precision"`
	Issuer             string         `json:"issuer"`
	DynamicAssetDataID string         `json:"dynamic_asset_data_id"`
}

type BlockHeader struct {
	TransactionMerkleRoot string            `json:"transaction_merkle_root"`
	Previous              string            `json:"previous"`
	Timestamp             types.Time        `json:"timestamp"`
	Witness               string            `json:"witness"`
	Extensions            []json.RawMessage `json:"extensions"`
}

//todo operation_results
type Block struct {
	TransactionMerkleRoot string              `json:"transaction_merkle_root"`
	Previous              string              `json:"previous"`
	Timestamp             types.Time          `json:"timestamp"`
	Witness               string              `json:"witness"`
	Extensions            []json.RawMessage   `json:"extensions"`
	WitnessSignature      string              `json:"witness_signature"`
	Transactions          []types.Transaction `json:"transactions"`
	BlockId               string              `json:"block_id"`
	SigningKey            string              `json:"signing_key"`
	TransactionIds        []string            `json:"transaction_ids"`
	RefBlockPrefix        uint32              `json:"ref_block_prefix"`
}

type MarketTicker struct {
	Time          types.Time     `json:"time"`
	Base          types.ObjectID `json:"base"`
	Quote         types.ObjectID `json:"quote"`
	Latest        string         `json:"latest"`
	LowestAsk     string         `json:"lowest_ask"`
	HighestBid    string         `json:"highest_bid"`
	PercentChange string         `json:"percent_change"`
	BaseVolume    string         `json:"base_volume"`
	QuoteVolume   string         `json:"quote_volume"`
}

type LimitOrder struct {
	ID          types.ObjectID `json:"id"`
	Expiration  types.Time     `json:"expiration"`
	Seller      types.ObjectID `json:"seller"`
	ForSale     types.Suint64  `json:"for_sale"`
	DeferredFee uint64         `json:"deferred_fee"`
	SellPrice   types.Price    `json:"sell_price"`
}

type DynamicGlobalProperties struct {
	ID                             types.ObjectID `json:"id"`
	HeadBlockNumber                uint32         `json:"head_block_number"`
	HeadBlockID                    string         `json:"head_block_id"`
	Time                           types.Time     `json:"time"`
	CurrentWitness                 types.ObjectID `json:"current_witness"`
	NextMaintenanceTime            types.Time     `json:"next_maintenance_time"`
	LastBudgetTime                 types.Time     `json:"last_budget_time"`
	AccountsRegisteredThisInterval int            `json:"accounts_registered_this_interval"`
	DynamicFlags                   int            `json:"dynamic_flags"`
	RecentSlotsFilled              string         `json:"recent_slots_filled"`
	LastIrreversibleBlockNum       uint32         `json:"last_irreversible_block_num"`
	CurrentAslot                   int64          `json:"current_aslot"`
	WitnessBudget                  int64          `json:"witness_budget"`
	RecentlyMissedCount            int64          `json:"recently_missed_count"`
	Parameters                     string         `json:"parameters"`
}

type Config struct {
	GrapheneSymbol               string `json:"GRAPHENE_SYMBOL"`
	GrapheneAddressPrefix        string `json:"GRAPHENE_ADDRESS_PREFIX"`
	GrapheneMinAccountNameLength uint8  `json:"GRAPHENE_MIN_ACCOUNT_NAME_LENGTH"`
	GrapheneMaxAccountNameLength uint8  `json:"GRAPHENE_MAX_ACCOUNT_NAME_LENGTH"`
	GrapheneMinAssetSymbolLength uint8  `json:"GRAPHENE_MIN_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxAssetSymbolLength uint8  `json:"GRAPHENE_MAX_ASSET_SYMBOL_LENGTH"`
	GrapheneMaxShareSupply       string `json:"GRAPHENE_MAX_SHARE_SUPPLY"`
}

type AccountsMap map[string]types.ObjectID

func (o *AccountsMap) UnmarshalJSON(b []byte) error {
	out := make(map[string]types.ObjectID)

	// unmarshal array
	var arr []json.RawMessage
	if err := json.Unmarshal(b, &arr); err != nil {
		return err
	}

	var (
		key string
		obj types.ObjectID
	)

	for _, item := range arr {
		account := []interface{}{&key, &obj}
		if err := json.Unmarshal(item, &account); err != nil {
			return err
		}

		out[key] = obj
	}

	*o = out
	return nil
}

type TableRowsParams struct {
	LowerBound    int64 `json:"lower_bound"`
	UpperBound    int64 `json:"upper_bound"`
	IndexPosition int64 `json:"index_position"`
	Limit         int64 `json:"limit"`
	Reverse       bool  `json:"reverse"`
}

func NewTableRowsParams() TableRowsParams {
	return TableRowsParams{0, -1, 1, 10, false}
}

type Witness struct {
	Id                    string `json:"id"`
	IsValid               bool   `json:"is_valid"`
	LastAslot             uint64 `json:"last_aslot"`
	LastConfirmedBlockNum int64  `json:"last_confirmed_block_num"`
	PayVb                 string `json:"pay_vb"`
	SigningKey            string `json:"signing_key"`
	TotalMissed           uint64 `json:"total_missed"`
	TotalVotes            string `json:"total_votes"`
	Url                   string `json:"url"`
	VoteId                string `json:"vote_id"`
	WitnessAccount        string `json:"witness_account"`
}

type Committee struct {
	Id         string `json:"id"`
	PayVb      string `json:"committee_member_account"`
	VoteId     string `json:"vote_id"`
	TotalVotes uint64 `json:"total_votes"`
	Url        string `json:"url"`
	IsValid    bool   `json:"is_valid"`
}

type GlobalProperties struct {
	Properties string
}

func (o *GlobalProperties) UnmarshalJSON(b []byte) error {
	o.Properties = string(b[:])
	return nil
}
