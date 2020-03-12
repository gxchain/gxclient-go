# gxclient-go
A client to interact with gxchain implemented in GO

<p>
 <a href='javascript:;'>
   <img width="300px" src='https://raw.githubusercontent.com/gxchain/gxips/master/assets/images/task-gxclient.png'/>
 </a>
 <a href='javascript:;'>
   <img width="300px" src='https://raw.githubusercontent.com/gxchain/gxips/master/assets/images/task-gxclient-en.png'/>
 </a>
</p> 

## Usage
```
import "github.com/gxchain/gxclient-go"
```

# APIs
- [ ] [Keypair API](#keypair-api)
- [x] [Chain API](#chain-api)
- [ ] [Faucet API](#faucet-api)
- [x] [Account API](#account-api)
- [x] [Asset API](#asset-api)
- [ ] [Contract API](#contract-api)
- [x] [Staking API](#staking-api)

## Constructors
```
//init client
func NewClient(actPriKeyWif, memoPriKeyWif, accountName, url string) (*Client, error) {
```


## Chain API
```
//broadcast transaction
func (client *Client) broadcast(stx *types.SignedTransaction) error
//broadcast transaction
func (client *Client) broadcastSync(stx *types.SignedTransaction) (*types.BroadcastResponse, error)
// GET ChainId of entry point
func (api *API) GetChainId() (string, error)
// Gets dynamic global properties of current blockchain
func (api *API) GetDynamicGlobalProperties() (*DynamicGlobalProperties, error)
// Get block by block height
func (api *API) GetBlock(blockNum uint32) (*Block, error)
//get block objects
func (api *API) GetObjects(objectIds ...string) ([]json.RawMessage, error)
//get block object
func (api *API) GetObject(objectId string) (json.RawMessage, error)
//send transfer request to entryPoint node
func (client *Client) Transfer(to, memo, amountAsset, feeSymbol string, broadcast bool) (*types.TransactionResult, error)
```

## Account API
```
// get account info by account name
func (api *API) GetAccount(account string) (*types.Account, error)
// get accounts info by account names
func (api *API) GetAccounts(accounts ...string) ([]*types.Account, error)
//get account balances by account name
func (api *API) GetAccountBalances(accountID string, assets ...string) ([]*types.AssetAmount, error)
//get account_ids by public key
func (api *API) GetAccountsByPublicKey(publicKeys string) ([]string, error)
```

## Asset API
```
// get assets corresponding to the provided symbols or IDs
func (api *API) GetAssets(symbols ...string) ([]*Asset, error)
// get assets corresponding to the provided symbol or ID
func (api *API) GetAsset(symbol string) (*Asset, error)
```

## Staking API

```
//Get staking programs
func (api *API) GetStakingPrograms() ([]*types.StakingProgram, error)
//create staking
func (client *Client) CreateStaking(to string, amount float64, programId, feeSymbol string, broadcast bool) (*types.TransactionResult, error)
//update staking by stakingId
func (client *Client) UpdateStaking(to, stakingId, feeSymbol string, broadcast bool) (*types.TransactionResult, error) 
//claim staking by stakingId
func (client *Client) ClaimStaking(stakingId, feeSymbol string, broadcast bool) (*types.TransactionResult, error)
```

