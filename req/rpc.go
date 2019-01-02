package req

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hzxiao/goutil"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	TestNetNode = `http://seed4.neo.org:20332`
	MainNetNode = `http://seed4.neo.org:10332`
	PriNetNode  = ``

	TestNet = "testnet"
	MainNet = "mainnet"
	PriNet  = "prinet"
)

var Node string

func SetNetwork(network string, nodeAddr string) error {
	switch network {
	case TestNet:
		Node = TestNetNode
	case MainNet:
		Node = MainNetNode
	case PriNet:
		if nodeAddr == "" {
			return fmt.Errorf("you send rpc request to private network, please provider a Node address")
		}
	default:
		return fmt.Errorf("unknown network")
	}

	if nodeAddr != "" {
		Node = nodeAddr
	}

	return nil
}

func GetUtxoOutputValue(txid string, index int) (uint64, error) {
	if Node == "" {
		return 0, fmt.Errorf("Node is empth, you haven't call rpc.SetNetwork func to set network and Node")
	}

	msg := NewRpcMsg("gettxout", txid, index)
	result, err := SendRpcReq(msg)
	if err != nil {
		return 0, fmt.Errorf("send rpc msg err: %v", err)
	}

	if len(result.GetMap("error")) > 0 {
		return 0, fmt.Errorf("err: code = %v, message = %v",
			result.GetInt64P("error/code"), result.GetStringP("error/message"))
	}
	return uint64(result.GetInt64P("result/value")), nil
}

func SendRawTransaction(raw string) error {
	if Node == "" {
		return fmt.Errorf("node is empth, you haven't call rpc.SetNetwork func to set network and node")
	}

	msg := NewRpcMsg("sendrawtransaction", raw)
	result, err := SendRpcReq(msg)
	if err != nil {
		return fmt.Errorf("send rpc msg err: %v", err)
	}

	if result.GetBool("result") {
		return nil
	}

	return fmt.Errorf("err: code = %v, message = %v",
		result.GetInt64P("error/code"), result.GetStringP("error/message"))
}

type RpcMsg struct {
	JsonRpc string        `json:"jsonrpc"`
	Id      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

func NewRpcMsg(method string, params ...interface{}) *RpcMsg {
	return &RpcMsg{
		JsonRpc: "2.0",
		Id:      1,
		Method:  method,
		Params:  params,
	}
}

func SendRpcReq(msg *RpcMsg) (goutil.Map, error) {
	if Node == "" {
		return nil, fmt.Errorf("node is empth, you haven't call rpc.SetNetwork func to set network and node")
	}
	if msg == nil {
		return nil, fmt.Errorf("msg is nil")
	}

	//send request to node
	req, err := http.NewRequest("POST", Node, bytes.NewBufferString(goutil.Struct2Json(msg)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	client.Timeout = 20 * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//read response
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := jsonDecode(buf)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func jsonDecode(buf []byte) (goutil.Map, error) {
	var data = goutil.Map{}
	err := json.Unmarshal(buf, &data)
	return data, err
}
