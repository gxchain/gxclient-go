package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"gxclient-go/rpc"
	"gxclient-go/types"
)

type API struct {
	caller rpc.Caller
	id     rpc.APIID
}

func NewAPI(id rpc.APIID, caller rpc.Caller) *API {
	return &API{id: id, caller: caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	err := api.caller.Connect()
	if err != nil {
		return err
	}
	return api.caller.Call(api.id, method, args, reply)
}

func (api *API) setCallback(method string, callback func(raw json.RawMessage)) error {
	return api.caller.SetCallback(api.id, method, callback)
}

// GET ChainId of entry point
func (api *API) GetChainId() (string, error) {
	var resp string
	err := api.call("get_chain_id", rpc.EmptyParams, &resp)
	return resp, err
}

// Gets dynamic global properties of current blockchain
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error) {
	var resp DynamicGlobalProperties
	err := api.call("get_dynamic_global_properties", rpc.EmptyParams, &resp)
	return &resp, err
}

// Get block by block height
func (api *API) GetBlock(blockNum uint32) (*Block, error) {
	var resp Block
	err := api.call("get_block", []interface{}{blockNum}, &resp)
	return &resp, err
}

func (api *API) GetObjects(objectIds ...string) ([]json.RawMessage, error) {
	var resp []json.RawMessage
	err := api.call("get_objects", []interface{}{objectIds}, &resp)
	return resp, err
}

func (api *API) GetObject(objectId string) (json.RawMessage, error) {
	var resp []json.RawMessage
	err := api.call("get_objects", []interface{}{[]string{objectId}}, &resp)
	return resp[0], err
}

//get_account_by_name
func (api *API) GetAccount(account string) (*types.Account, error) {
	var resp *types.Account
	if err := api.call("get_account_by_name", []interface{}{account}, &resp); err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.Errorf("account %s not exist", account)
	}
	return resp, nil
}

func (api *API) GetAccounts(accounts ...string) ([]*types.Account, error) {
	ret := make([]*types.Account, len(accounts))
	for i, account := range accounts {
		a, err := api.GetAccount(account)
		if err != nil {
			return nil, err
		}
		ret[i] = a
	}
	return ret, nil
}

func (api *API) GetAccountsByIds(ids ...string) ([]*types.Account, error) {
	var resp []*types.Account
	if err := api.call("get_accounts", []interface{}{ids}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// GetAccountBalances
// Get an account’s balances in various assets.
func (api *API) GetAccountBalances(accountID string, assets ...string) ([]*types.AssetAmount, error) {
	var resp []*types.AssetAmount
	array := []string{}
	err := api.call("get_account_balances", []interface{}{accountID, append(array, assets...)}, &resp)
	return resp, err
}

//get_key_references
func (api *API) GetAccountsByPublicKeys(publicKeys ...string) (*[][]string, error) {
	var resp *[][]string
	err := api.call("get_key_references", []interface{}{publicKeys}, &resp)
	return resp, err
}

//get_key_references
func (api *API) GetAccountsByPublicKey(publicKeys string) ([]string, error) {
	var resp [][]string
	if err := api.call("get_key_references", []interface{}{[]string{publicKeys}}, &resp); err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}
	//Deduplication
	var result []string
	temp := map[string]struct{}{}
	for _, item := range resp[0] {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result, nil
}

//lookup_asset_symbols
func (api *API) GetAssets(symbols ...string) ([]*Asset, error) {
	var resp []*Asset
	err := api.call("lookup_asset_symbols", []interface{}{symbols}, &resp)
	return resp, err
}

// LookupAssetSymbols get assets corresponding to the provided symbol or IDs
func (api *API) GetAsset(symbol string) (*Asset, error) {
	var resp []*Asset
	if err := api.call("lookup_asset_symbols", []interface{}{[]string{symbol}}, &resp); err != nil {
		return nil, err
	}
	if resp[0] == nil {
		return nil, errors.Errorf("assets %s not exist", symbol)
	}
	return resp[0], nil
}

//get_contract_account_by_name
func (api *API) GetWitnessByAccount(accountId string) (*Witness, error) {
	var resp *Witness
	if err := api.call("get_witness_by_account", []interface{}{accountId}, &resp); err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, errors.Errorf("%s is not witness", accountId)
	}
	return resp, nil
}

// Semantically equivalent to get_account_balances, but takes a name instead of an ID.
func (api *API) GetNamedAccountBalances(account string, assets ...string) ([]*types.AssetAmount, error) {
	var resp []*types.AssetAmount
	array := []string{}
	err := api.call("get_named_account_balances", []interface{}{account, append(array, assets...)}, &resp)
	return resp, err
}

// LookupAccounts gets names and IDs for registered accounts
// lower_bound_name: Lower bound of the first name to return
// limit: Maximum number of results to return must not exceed 1000
func (api *API) LookupAccounts(lowerBoundName string, limit uint16) (AccountsMap, error) {
	var resp AccountsMap
	err := api.call("lookup_accounts", []interface{}{lowerBoundName, limit}, &resp)
	return resp, err
}

// GetRequiredFee fetchs fee for operations
func (api *API) GetRequiredFee(ops []types.Operation, assetID string) ([]types.AssetAmount, error) {
	var resp []types.AssetAmount

	opsJSON := []interface{}{}
	for _, o := range ops {
		_, err := json.Marshal(o)
		if err != nil {
			return []types.AssetAmount{}, err
		}

		opArr := []interface{}{o.Type(), o}

		opsJSON = append(opsJSON, opArr)
	}
	if err := api.call("get_required_fees", []interface{}{opsJSON, assetID}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

// get_staking_objects
func (api *API) GetStakingObjects(accountID string) ([]*types.StakingObject, error) {
	var resp []*types.StakingObject
	err := api.call("get_staking_objects", []interface{}{accountID}, &resp)
	return resp, err
}

func (api *API) GetStakingPrograms() ([]*types.StakingProgram, error) {
	var programs []*types.StakingProgram
	globalProperties, err := api.getGlobalProperties()
	if err != nil {
		return nil, err
	}
	properties := gjson.Parse(globalProperties.Properties)
	extensions := properties.Get("parameters").Get("extensions").Array()
	for _, ex := range extensions {
		if ex.Array()[0].String() == "11" {
			for _, param := range ex.Array()[1].Get("params").Array() {
				program := types.StakingProgram{
					ProgramId:   param.Array()[0].String(),
					Weight:      uint32(param.Array()[1].Get("weight").Uint()),
					StakingDays: uint32(param.Array()[1].Get("staking_days").Uint()),
				}
				programs = append(programs, &program)
			}
		}
	}
	return programs, err
}

func (api *API) getGlobalProperties() (*GlobalProperties, error) {
	var resp *GlobalProperties
	err := api.call("get_global_properties", rpc.EmptyParams, &resp)
	return resp, err
}

// GetBlockHeader returns block header by the given block number
func (api *API) GetBlockHeader(blockNum uint32) (*BlockHeader, error) {
	var resp BlockHeader
	err := api.call("get_block_header", []interface{}{blockNum}, &resp)
	return &resp, err
}

func (api *API) GetTransactionByTxid(txid string) (*types.Transaction, error) {
	var resp *types.Transaction
	if err := api.call("get_transaction_rows", []interface{}{txid}, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (api *API) GeTransactionExtByTxid(txid string) (*types.TransactionExt, error) {
	var resp *types.TransactionExt
	if err := api.call("get_transaction_by_txid", []interface{}{txid}, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}
