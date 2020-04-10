package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"gxclient-go/rpc"
	"io/ioutil"
	"math"
	"net/http"
	"sync"
	"time"
)

type Transport struct {
	Url    string
	client http.Client

	requestID uint64
	reqMutex  sync.Mutex
}

func NewTransport(url string) *Transport {
	timeout := time.Duration(20 * time.Second)

	return &Transport{
		client: http.Client{
			Timeout: timeout,
		},
		Url: url,
	}
}

func (p *Transport) Connect() error {
	return nil
}

func (caller *Transport) Call(api rpc.APIID, method string, args []interface{}, reply interface{}) error {
	caller.reqMutex.Lock()
	defer caller.reqMutex.Unlock()

	// increase request id
	if caller.requestID == math.MaxUint64 {
		caller.requestID = 0
	}
	caller.requestID++

	request := rpc.RPCRequest{
		Method: "call",
		ID:     caller.requestID,
		Params: []interface{}{api, method, args},
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := caller.client.Post(caller.Url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read body")
	}

	var rpcResponse rpc.RPCResponse
	if err = json.Unmarshal(respBody, &rpcResponse); err != nil {
		return errors.Wrapf(err, "failed to unmarshal response: %+v", string(respBody))
	}

	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	if rpcResponse.Result != nil {
		if err := json.Unmarshal(*rpcResponse.Result, reply); err != nil {
			return errors.Wrapf(err, "failed to unmarshal rpc result: %+v", string(*rpcResponse.Result))
		}
	}

	return nil
}

func (caller *Transport) SetCallback(api rpc.APIID, method string, notice func(args json.RawMessage)) error {
	return nil
}

func (caller *Transport) Close() error {
	return nil
}
