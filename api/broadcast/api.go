package broadcast

import (
	"gxclient-go/rpc"
	"gxclient-go/types"
	"reflect"
)

type API struct {
	caller rpc.Caller
	id     rpc.APIID
}

func NewAPI(id rpc.APIID, caller rpc.Caller) *API {
	return &API{id: id, caller: caller}
}

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	err := api.caller.Connect()
	if err != nil {
		return err
	}
	return api.caller.Call(api.id, method, args, reply)
}

// BroadcastTransaction broadcast a transaction to the network.
func (api *API) BroadcastTransaction(tx *types.Transaction) error {
	var reply interface{}
	return api.call("broadcast_transaction", []interface{}{tx}, &reply)
}

func (api *API) BroadcastTransactionSynchronous(tx *types.Transaction) (*types.BroadcastResponse, error) {
	response := types.BroadcastResponse{}
	if err := api.call("broadcast_transaction_synchronous", []interface{}{tx}, &response); err != nil {
		return nil, err
	}
	return &response, nil
}
