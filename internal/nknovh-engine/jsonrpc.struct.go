package nknovh_engine

import (
		"encoding/json"
		"time"
		"net/http"
)
type RPCResponse struct {
	Error struct {
		Code int `json:"code"`
		Data string `json:"data"`
		Message string `json:"message"`
	} `json:"error"`
	Jsonrpc string `json:"jsonrpc"`
	Id string `json:"id"`
}

type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method string      `json:"method"`
	Params *json.RawMessage `json:"params,ommitempty"`
	Id     int      `json:"id"`
}

type JsonRPCConf struct {
	Timeout time.Duration
	Ip string
	Method string
	Params *json.RawMessage
	Client *http.Client
}

type NodeSt struct {
	State NodeState
	Neighbor NodeNeighbor
}

type NodeState struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Addr               string `json:"addr"`
		Currtimestamp      int    `json:"currTimeStamp"`
		Height             int    `json:"height"`
		ID                 string `json:"id"`
		Jsonrpcport        int    `json:"jsonRpcPort"`
		ProposalSubmitted  int    `json:"proposalSubmitted"`
		ProtocolVersion    int    `json:"protocolVersion"`
		Publickey          string `json:"publicKey"`
		RelayMessageCount  int    `json:"relayMessageCount"`
		SyncState          string `json:"syncState"`
		Tlsjsonrpcdomain   string `json:"tlsJsonRpcDomain"`
		Tlsjsonrpcport     int    `json:"tlsJsonRpcPort"`
		Tlswebsocketdomain string `json:"tlsWebsocketDomain"`
		Tlswebsocketport   int    `json:"tlsWebsocketPort"`
		Uptime             int    `json:"uptime"`
		Version            string `json:"version"`
		Websocketport      int    `json:"websocketPort"`
	} `json:"result"`
}


type NodeNeighbor struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result []struct {
		Addr               string `json:"addr"`
		Height             int    `json:"height"`
		ID                 string `json:"id"`
	//	Isoutbound         bool   `json:"isOutbound"`
	//	Jsonrpcport        int    `json:"jsonRpcPort"`
	//	Protocolversion    int    `json:"protocolVersion"`
	//	PublicKey          string `json:"publicKey"`
	//	RoundTripTime      int    `json:"roundTripTime"`
		SyncState          string `json:"syncState"`
	//	Tlsjsonrpcdomain   string `json:"tlsJsonRpcDomain"`
	//	Tlsjsonrpcport     int    `json:"tlsJsonRpcPort"`
	//	Tlswebsocketdomain string `json:"tlsWebsocketDomain"`
	//	Tlswebsocketport   int    `json:"tlsWebsocketPort"`
	//	Websocketport      int    `json:"websocketPort"`
	} `json: "result"`
}
