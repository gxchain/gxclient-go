package http

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/pquerna/ffjson/ffjson"
	"gxclient-go/rpc"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

type rpcRequest struct {
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      uint64        `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
}

func (p *rpcRequest) reset() {
	p.ID = 0
	p.Params = nil
	p.Method = ""
}

type ResponseErrorContext struct {
	Level      string `json:"level"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Method     string `json:"method"`
	Hostname   string `json:"hostname"`
	ThreadName string `json:"thread_name"`
	Timestamp  string `json:"timestamp"`
}

type ResponseErrorStack struct {
	Context ResponseErrorContext `json:"context"`
	Format  string               `json:"format"`
	Data    interface{}          `json:"data"`
}

type ResponseErrorData struct {
	Code    int                  `json:"code"`
	Name    string               `json:"name"`
	Message string               `json:"message"`
	Stack   []ResponseErrorStack `json:"stack"`
}

type ResponseError struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    ResponseErrorData `json:"data"`
}

func (p ResponseError) Error() string {
	return p.Message
}

type rpcResponseString struct {
	ID      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   ResponseError `json:"error"`
}

func (p rpcResponseString) HasError() bool {
	return p.Error.Code != 0
}

type HttpTransport struct {
	*http.Client
	endpointURL string
	*ffjson.Encoder
	req    rpcRequest
	decBuf *bytes.Buffer
}

func NewHttpTransport(url string) (*HttpTransport, error) {
	cli := HttpTransport{
		endpointURL: url,
	}

	return &cli, nil
}

func (p *HttpTransport) Connect() error {
	p.Client = &http.Client{
		Timeout: 20 * time.Second,
	}

	p.decBuf = new(bytes.Buffer)
	p.Encoder = ffjson.NewEncoder(p.decBuf)

	return nil
}

func (caller *HttpTransport) Call(api rpc.APIID, method string, args []interface{}, reply interface{}) error {
	// increase request id
	caller.req.JsonRpc = "2.0"
	caller.req.Method = method
	if caller.req.ID == math.MaxUint64 {
		caller.req.ID = 0
	}
	caller.req.ID++
	caller.req.Params = args

	if err := caller.Encode(&caller.req); err != nil {
		return errors.Wrap(err, "Encode")
	}

	data := caller.decBuf.Bytes()
	s := string(data[:])
	println(s)

	req, err := http.NewRequest("POST", caller.endpointURL, caller.decBuf)
	if err != nil {
		return errors.Wrap(err, "NewRequest")
	}

	req.Close = true
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := caller.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	s1 := string(body[:])
	println(s1)

	var res rpcResponseString
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	var d []byte
	d, err = json.Marshal(res.Result)
	if d == nil || err != nil {
		return errors.Wrap(err, "marshal error")
	}

	if err := json.Unmarshal(d, &reply); err != nil {
		return err
	}

	return nil
}

func (caller *HttpTransport) SetCallback(api rpc.APIID, method string, notice func(args json.RawMessage)) error {
	return nil
}

func (caller *HttpTransport) Close() error {
	return nil
}
