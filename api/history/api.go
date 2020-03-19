package history

import (
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

// GetMarketHistory returns market history base/quote (candlesticks) for the given period
func (api *API) GetMarketHistory(base, quote types.ObjectID, bucketSeconds uint32, start, end types.Time) ([]*Bucket, error) {
	var resp []*Bucket
	err := api.call("get_market_history", []interface{}{base.String(), quote.String(), bucketSeconds, start, end}, &resp)
	return resp, err
}

// GetMarketHistoryBuckets returns a list of buckets that can be passed to
// `GetMarketHistory` as the `bucketSeconds` argument
func (api *API) GetMarketHistoryBuckets() ([]uint32, error) {
	var resp []uint32
	err := api.call("get_market_history_buckets", rpc.EmptyParams, &resp)
	return resp, err
}

// GetFillOrderHistory returns filled orders
func (api *API) GetFillOrderHistory(base, quote types.ObjectID, limit uint32) ([]*OrderHistory, error) {
	var resp []*OrderHistory
	err := api.call("get_fill_order_history", []interface{}{base.String(), quote.String(), limit}, &resp)
	return resp, err
}

// GetAccountHistory gets operations relevant to the specified account
// account: The account whose history should be queried
// stop: ID of the earliest operation to retrieve
// limit: Maximum number of operations to retrieve (must not exceed 100)
// start: ID of the most recent operation to retrieve
func (api *API) GetAccountHistory(account, stop string, limit int, start string) ([]*OperationHistory, error) {
	var history []*OperationHistory
	err := api.call("get_account_history", []interface{}{account, stop, limit, start}, &history)
	return history, err
}
