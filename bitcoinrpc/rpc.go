package bitcoinrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RpcClient struct {
	URL      string
	Username string
	Password string
	Client   *http.Client
}

type RpcRequest struct {
	Jsonrpc string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *RpcError       `json:"error"`
	ID     int             `json:"id"`
}

// RpcError captures an error returned by the RPC server
type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewRpcClient creates a new RpcClient
func NewRpcClient(host string, username string, password string) *RpcClient {
	return &RpcClient{
		URL:      host,
		Username: username,
		Password: password,
		Client:   &http.Client{},
	}
}

// Do makes an RPC call to the specified method with the provided parameters
func (rpc *RpcClient) Do(method string, params []interface{}) (json.RawMessage, error) {
	requestData := RpcRequest{
		Jsonrpc: "1.0",
		ID:      1,
		Method:  method,
		Params:  params,
	}
	jsonRequest, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", rpc.URL, bytes.NewBuffer(jsonRequest))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(rpc.Username, rpc.Password)

	resp, err := rpc.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResponse RpcResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return nil, err
	}

	if rpcResponse.Error != nil {
		return nil, fmt.Errorf("RPC error (%d): %s", rpcResponse.Error.Code, rpcResponse.Error.Message)
	}

	return rpcResponse.Result, nil
}
