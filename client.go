package gxclient_go

import (
	"bytes"
	"encoding/hex"
	_ "encoding/hex"
	"fmt"
	_ "fmt"
	"github.com/pkg/errors"
	"gxclient-go/api/broadcast"
	"gxclient-go/api/database"
	"gxclient-go/api/history"
	"gxclient-go/api/login"
	"gxclient-go/rpc"
	"gxclient-go/rpc/http"
	"gxclient-go/rpc/websocket"
	"gxclient-go/sign"
	"gxclient-go/transaction"
	"gxclient-go/types"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	cc rpc.CallCloser

	// Database represents database_api
	Database *database.API

	// NetworkBroadcast represents network_broadcast_api
	Broadcast *broadcast.API
	//
	// History represents history_api
	History *history.API

	// Login represents login_api
	Login *login.API

	chainID string

	activePriKey *types.PrivateKey

	memoPriKey *types.PrivateKey

	account *types.Account
}

// NewClient creates a new RPC client
func NewClient(actPriKeyWif, memoPriKeyWif, accountName, url string) (*Client, error) {
	// transport
	var cc rpc.CallCloser
	var err error
	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		cc, err = http.NewHttpTransport(url)
	} else {
		cc, err = websocket.NewTransport(url)
	}
	if err != nil {
		return nil, err
	}

	client := &Client{cc: cc}
	activeKey, err := types.NewPrivateKeyFromWif(actPriKeyWif)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init active private key")
	}
	client.activePriKey = activeKey

	memoKey, err := types.NewPrivateKeyFromWif(memoPriKeyWif)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init memo private key")
	}
	client.memoPriKey = memoKey

	if strings.HasPrefix(url, "http") || strings.HasPrefix(url, "https") {
		client.Database = database.NewAPI(0, cc)
		chainID, err := client.Database.GetChainId()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get database ID")
		}
		client.chainID = chainID
		client.History = history.NewAPI(1, cc)
		client.Broadcast = broadcast.NewAPI(1, cc)
		return client, nil
	}

	// login
	loginAPI := login.NewAPI(cc)
	client.Login = loginAPI

	// database
	databaseAPIID, err := loginAPI.Database()
	if err != nil {
		return nil, err
	}
	client.Database = database.NewAPI(databaseAPIID, client.cc)

	// database ID
	chainID, err := client.Database.GetChainId()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get database ID")
	}
	client.chainID = chainID

	// history
	historyAPIID, err := loginAPI.History()
	if err != nil {
		return nil, err
	}
	client.History = history.NewAPI(historyAPIID, client.cc)

	// network broadcast
	networkBroadcastAPIID, err := loginAPI.NetworkBroadcast()
	if err != nil {
		return nil, err
	}
	client.Broadcast = broadcast.NewAPI(networkBroadcastAPIID, client.cc)

	account, err := client.Database.GetAccount(accountName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init account")
	}
	client.account = account

	return client, nil
}

// Close should be used to close the client when no longer needed.
// It simply calls Close() on the underlying CallCloser.
func (client *Client) Close() error {
	return client.cc.Close()
}

// Transfer a certain amount of the given asset
func (client *Client) Transfer(to, memo, amountAsset, feeSymbol string, broadcast bool) (*types.TransactionResult, error) {
	toAccount, err := client.Database.GetAccount(to)
	if err != nil {
		return nil, err
	}

	amountAndSymbol := strings.Split(amountAsset, " ")
	if len(amountAndSymbol) != 2 {
		return nil, errors.New("amountAsset incorrect format!")
	}

	amountSymbol, err := client.Database.GetAsset(amountAndSymbol[1])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseFloat(amountAndSymbol[0], 64)
	if err != nil {
		return nil, err
	}

	amountAssets := types.AssetAmount{
		AssetID: amountSymbol.ID,
		Amount:  uint64(amount * math.Pow10(int(amountSymbol.Precision))),
	}

	fee, err := client.Database.GetAsset(feeSymbol)
	if err != nil {
		return nil, err
	}
	feeAssets := types.AssetAmount{
		AssetID: fee.ID,
		Amount:  0,
	}

	var memoOb = &types.Memo{}

	var keys []*types.PublicKey
	for k := range client.account.Active.KeyAuths {
		keys = append(keys, k)
	}

	memoOb.From = *keys[0]
	memoOb.To = toAccount.Options.MemoKey
	memoOb.Nonce = types.GetNonce()
	if len(memo) > 0 {
		err := memoOb.Encrypt(client.memoPriKey, memo)
		if err != nil {
			return nil, err
		}
	} else {
		memoOb = nil
	}

	op := types.NewTransferOperation(types.MustParseObjectID(client.account.ID.String()), types.MustParseObjectID(toAccount.ID.String()), amountAssets, feeAssets, memoOb)

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, feeAssets.AssetID.String())
	if err != nil {
		return nil, err
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{client.activePriKey.ToWIF()}, op)
	if err != nil {
		return nil, err
	}

	result := new(types.TransactionResult)
	result.SignedTransaction = stx

	if broadcast {
		resp, err := client.broadcastSync(stx)
		if err != nil {
			return result, err
		}
		result.BroadcastResponse = resp
	}
	return result, err
}

// Create a Staking
func (client *Client) CreateStaking(to string, amount float64, programId, feeSymbol string, broadcast bool) (*types.TransactionResult, error) {
	//trustNode
	toAccount, err := client.Database.GetAccount(to)
	if err != nil {
		return nil, err
	}
	witness, err := client.Database.GetWitnessByAccount(toAccount.ID.String())
	if err != nil {
		return nil, err
	}
	trustNodeId := types.MustParseObjectID(witness.Id)

	//amount
	asset, err := client.Database.GetAsset("GXC")
	if err != nil {
		return nil, err
	}
	amountAssets := types.AssetAmount{
		AssetID: asset.ID,
		Amount:  uint64(amount * math.Pow10(int(asset.Precision))),
	}
	fee, err := client.Database.GetAsset(feeSymbol)
	if err != nil {
		return nil, err
	}
	feeAssets := types.AssetAmount{
		AssetID: fee.ID,
		Amount:  0,
	}

	programs, err := client.Database.GetStakingPrograms()
	if err != nil {
		return nil, err
	}
	var stakingProgram *types.StakingProgram
	for _, program := range programs {
		if program.ProgramId == programId {
			stakingProgram = program
		}
	}
	if stakingProgram == nil {
		return nil, errors.Errorf("programId %s illegal!", programId)
	}

	op := types.NewStakingCreateOperation(types.MustParseObjectID(client.account.ID.String()), trustNodeId, amountAssets, feeAssets, programId, stakingProgram.Weight, stakingProgram.StakingDays)

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, feeAssets.AssetID.String())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{client.activePriKey.ToWIF()}, op)
	if err != nil {
		return nil, err
	}

	result := new(types.TransactionResult)
	result.SignedTransaction = stx

	if broadcast {
		resp, err := client.broadcastSync(stx)
		if err != nil {
			return result, err
		}
		result.BroadcastResponse = resp
	}
	return result, err
}

// Update a Staking
func (client *Client) UpdateStaking(to, stakingId, feeSymbol string, broadcast bool) (*types.TransactionResult, error) {
	toAccount, err := client.Database.GetAccount(to)
	if err != nil {
		return nil, err
	}
	witness, err := client.Database.GetWitnessByAccount(toAccount.ID.String())
	if err != nil {
		return nil, err
	}
	trustNodeId := types.MustParseObjectID(witness.Id)

	fee, err := client.Database.GetAsset(feeSymbol)
	if err != nil {
		return nil, err
	}
	feeAssets := types.AssetAmount{
		AssetID: fee.ID,
		Amount:  0,
	}

	op := types.NewStakingUpdateOperation(types.MustParseObjectID(client.account.ID.String()), trustNodeId, types.MustParseObjectID(stakingId), feeAssets)

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, feeAssets.AssetID.String())
	if err != nil {
		return nil, err
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{client.activePriKey.ToWIF()}, op)
	if err != nil {
		return nil, err
	}

	result := new(types.TransactionResult)
	result.SignedTransaction = stx

	if broadcast {
		resp, err := client.broadcastSync(stx)
		if err != nil {
			return result, err
		}
		result.BroadcastResponse = resp
	}
	return result, err
}

// Claim a Staking
func (client *Client) ClaimStaking(stakingId, feeSymbol string, broadcast bool) (*types.TransactionResult, error) {
	fee, err := client.Database.GetAsset(feeSymbol)
	if err != nil {
		return nil, err
	}
	feeAssets := types.AssetAmount{
		AssetID: fee.ID,
		Amount:  0,
	}
	op := types.NewStakingClaimOperation(types.MustParseObjectID(client.account.ID.String()), types.MustParseObjectID(stakingId), feeAssets)

	fees, err := client.Database.GetRequiredFee([]types.Operation{op}, feeAssets.AssetID.String())
	if err != nil {
		return nil, err
	}
	op.Fee.Amount = fees[0].Amount

	stx, err := client.sign([]string{client.activePriKey.ToWIF()}, op)
	if err != nil {
		return nil, err
	}

	result := new(types.TransactionResult)
	result.SignedTransaction = stx

	if broadcast {
		resp, err := client.broadcastSync(stx)
		if err != nil {
			return result, err
		}
		result.BroadcastResponse = resp
	}
	return result, err
}

func (client *Client) sign(wifs []string, operations ...types.Operation) (*types.SignedTransaction, error) {
	props, err := client.Database.GetDynamicGlobalProperties()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dynamic global properties")
	}

	block, err := client.Database.GetBlock(props.LastIrreversibleBlockNum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get block")
	}

	refBlockPrefix, err := sign.RefBlockPrefix(block.Previous)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign block prefix")
	}

	expiration := props.Time.Add(10 * time.Minute)
	stx := types.NewSignedTransaction(&types.Transaction{
		RefBlockNum:    sign.RefBlockNum(props.LastIrreversibleBlockNum - 1&0xffff),
		RefBlockPrefix: refBlockPrefix,
		Expiration:     types.Time{Time: &expiration},
	})

	for _, op := range operations {
		stx.PushOperation(op)
	}

	var b bytes.Buffer
	x := transaction.NewEncoder(&b)

	if err := x.Encode(stx.Transaction); err != nil {
		return nil, nil
	}
	s := hex.EncodeToString(b.Bytes())
	fmt.Println(s)

	if err = stx.Sign(wifs, client.chainID); err != nil {
		return nil, errors.Wrap(err, "failed to sign the transaction")
	}

	return stx, nil
}

func (client *Client) broadcast(stx *types.SignedTransaction) error {
	return client.Broadcast.BroadcastTransaction(stx.Transaction)
}

func (client *Client) broadcastSync(stx *types.SignedTransaction) (*types.BroadcastResponse, error) {
	return client.Broadcast.BroadcastTransactionSynchronous(stx.Transaction)
}
