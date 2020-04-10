package rpc

import (
	"encoding/json"
	"io"
)

var EmptyParams = []interface{}{}

type APIID string

type Caller interface {
	Call(api APIID, method string, args []interface{}, reply interface{}) error
	SetCallback(api APIID, method string, callback func(raw json.RawMessage)) error
	Connect() error
}

type CallCloser interface {
	Caller
	io.Closer
}
