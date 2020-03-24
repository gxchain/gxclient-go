package faucet

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"gxclient-go/types"
	"io/ioutil"
	oldHttp "net/http"
	"time"
)

type rpcError struct {
	Error *rpcErrorStuuct `json:"error,omitempty"`
}

type rpcErrorStuuct struct {
	Base []string `json:"base,omitempty"`
}

func Register(faucet, account, activeKey, ownerKey, memoKey string) (*types.Transaction, error) {
	accountInfo := types.RegisterAccountInfo{}
	accountInfo.ActiveKey = activeKey
	if len(ownerKey) > 0 {
		accountInfo.OwnerKey = ownerKey
	} else {
		accountInfo.OwnerKey = activeKey
	}
	if len(memoKey) > 0 {
		accountInfo.MemoKey = memoKey
	} else {
		accountInfo.MemoKey = activeKey
	}
	accountInfo.Name = account

	accountReg := types.RegisterAccount{}
	accountReg.Account = accountInfo

	decBuf := new(bytes.Buffer)
	enc := ffjson.NewEncoder(decBuf)

	if err := enc.Encode(&accountReg); err != nil {
		return nil, errors.Wrap(err, "Encode")
	}

	c := &oldHttp.Client{
		Timeout: 10 * time.Second,
	}
	req, err := oldHttp.NewRequest("POST", faucet, decBuf)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest")
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s1 := string(body[:])
	println(s1)

	var rpcErr *rpcError
	if err = json.Unmarshal(body, &rpcErr); err != nil {
		return nil, err
	}
	if rpcErr.Error != nil {
		return nil, errors.Errorf("%s", string(body[:]))
	}

	var res *types.Transaction
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}

	return res, nil
}
