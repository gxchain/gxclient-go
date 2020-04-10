package login

import (
	"gxclient-go/rpc"
	"strconv"
)

const APIID = "1"

type API struct {
	caller rpc.Caller
}

func NewAPI(caller rpc.Caller) *API {
	return &API{caller}
}

func (api *API) call(method string, args []interface{}, reply interface{}) error {
	err := api.caller.Connect()
	if err != nil {
		return err
	}
	return api.caller.Call(rpc.APIID(APIID), method, args, reply)
}

func (api *API) GetApiByName(name string) (*uint8, error) {
	var id uint8
	err := api.call("get_api_by_name", []interface{}{name}, &id)
	return &id, err
}

func (api *API) Login(username, password string) (bool, error) {
	var resp bool
	err := api.call("login", []interface{}{username, password}, &resp)
	return resp, err
}

func (api *API) Database() (rpc.APIID, error) {
	var id rpc.APIID
	var idUint uint64
	err := api.call("database", rpc.EmptyParams, &idUint)
	id = rpc.APIID(strconv.FormatUint(idUint, 10))
	return id, err
}

func (api *API) History() (rpc.APIID, error) {
	var id rpc.APIID
	var idUint uint64
	err := api.call("history", rpc.EmptyParams, &idUint)
	id = rpc.APIID(strconv.FormatUint(idUint, 10))
	return id, err
}

func (api *API) NetworkBroadcast() (rpc.APIID, error) {
	var id rpc.APIID
	var idUint uint64
	err := api.call("network_broadcast", rpc.EmptyParams, &idUint)
	id = rpc.APIID(strconv.FormatUint(idUint, 10))
	return id, err
}
