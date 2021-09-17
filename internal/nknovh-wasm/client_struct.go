package nknovh_wasm

import (
		"xwasmapi"
		"sync"
		"syscall/js"
)

type CLIENT struct {
	Hash string
	Lang string
	Sort string
	Sort_type string
	Hide_attention bool
	Version string
	EntriesPerPage int
	CurrentPage int
	Daemon *Daemon
	Netstatus *Netstatus
	Nodes *Nodes
	Wallets *Wallets
	Prices *Prices
	Debug bool
	ShowOfflineOnly bool
	NodesSummary map[string]map[string]float64
	Conf *Conf
	Cached *Cached
	LANG *LANG
	Objects *Objects
	W *xwasmapi.Xwasmapi
	ws *js.Value
	mux *Mutexes
	AutoUpdaterIsStarted bool
	AutoUpdaterStartCh chan bool
	AutoUpdaterStopCh chan bool
	apiMethods map[string]func(*WSReply) interface{}
}

type Mutexes struct {
	Websocket *sync.Mutex
	StartView *sync.Mutex
	Nodes *sync.Mutex
	NodesSummary *sync.Mutex
	AutoUpdater *sync.RWMutex
}

type Objects struct {
	Channels map[string]*chan struct{}
	Listeners map[string]*js.Func
}

type Conf struct {
	DefaultLanguage string
	DefaultEntriesPerPage int
}

type Cached struct {
	Pages map[string]string
	Lang map[string]*LANG
}

type Wallets struct {
	Code  int `json:"Code"`
	Error bool `json:"Error, omitempty"`
	ErrMessage	string	`json:"ErrMessage, omitempty"`
	Value struct {
		Wallets []struct {
			Id int `json:"Id"`
			NknWallet string `json:"NknWallet"`
			Balance float64 `json:"Balance`
		} `json:"Wallets, omitempty`
	} `json:"Value, omitempty"`
}

type Prices struct {
	Code  int `json:"Code"`
	Error bool `json:"Error, omitempty"`
	ErrMessage string `json:"ErrMessage, omitempty"`
	Value struct {
		Usd float64 `json:"usd"`
	} `json:"Value, omitempty"`
}

type Netstatus struct {
	Method string `json: "Method"`
	Code  int `json:"Code"`
	Error bool `json:"Error, omitempty"`
	ErrMessage	string	`json:"ErrMessage, omitempty"`
	Value struct {
		Relays              int64   `json:"relays"`
		AverageUptime       int     `json:"average_uptime"`
		AverageRelays       int     `json:"average_relays"`
		RelaysPerHour       int     `json:"relays_per_hour"`
		ProposalSubmitted   int     `json:"proposalSubmitted"`
		PersistNodesCount   int     `json:"persist_nodes_count"`
		NodesCount          int     `json:"nodes_count"`
		LastHeight          int     `json:"last_height"`
		LastTimestamp       int     `json:"last_timestamp"`
		AverageBlockTime    float64 `json:"average_blockTime"`
		AverageBlocksPerDay float64 `json:"average_blocksPerDay"`
		LatestUpdate        string  `json:"latest_update"`
	} `json:"Value, omitempty"`
}

type Nodes struct {
	Method string `json: "Method"`
	Code  int `json:"Code"`
	Error bool `json:"Error, omitempty"`
	ErrMessage	string `json:"ErrMessage, omitempty"`
	Value struct {
		List []struct {
			NodeID            int    `json:"NodeId"`
			Err               int   `json:"Err,omitempty"`
			IP                string `json:"Ip"`
			Name              string `json:"Name"`
			Uptime            int    `json:"Uptime"`
			SyncState         string `json:"SyncState"`
			RelayMessageCount int    `json:"RelayMessageCount"`
			Currtimestamp     int    `json:"Currtimestamp"`
			ProposalSubmitted int    `json:"ProposalSubmitted"`
			RelaysPerHour     int    `json:"RelaysPerHour"`
			RelaysPerHour10   int    `json:"RelaysPerHour10"`
			RelaysPerHour60   int    `json:"RelaysPerHour60"`
			Height            int    `json:"Height"`
			Version           string `json:"Version"`
			LatestUpdate      string `json:"LatestUpdate"`
		} `json:"List"`
	} `json:"Value, omitempty"`
}

type Daemon struct {
	Code  int `json:"Code"`
	Error bool `json:"Error, omitempty"`
	ErrMessage string `json:"ErrMessage, omitempty"`
	Value struct {
		Version string `json:"Version"`
		Timezone string `json:"Timezone"`
	} `json:"Value, omitempty"`
}

type GetFullstack struct {
	Method string `json: "Method"`
	Code  int  `json:"Code"`
	Error bool `json:"Error"`
	ErrMessage string `json: "ErrMessage, omitempty"`
	Value struct {
		Netstatus Netstatus `json: "Netstatus, omitempty"`
		Nodes Nodes `json: "Nodes, omitempty"`
		Wallets Wallets `json: "Wallets, omitempty"`
		Prices Prices `json: "Prices, omitempty"`
		Daemon Daemon `json: "Daemon, omitempty"`
	} `json:"Value, omitempty"`
}

type NodeLookup struct {
	MinPing         int `json:"MinPing"`
	AvgPing         int `json:"AvgPing"`
	MaxPing         int `json:"MaxPing"`
	NeighborCount   int `json:"NeighborCount"`
	NeighborPersist int `json:"NeighborPersist"`
	RelaysPerHour   int `json:"RelaysPerHour"`
	NodeState NodeState `json:"NodeState"`
}

type RPCErrorState struct {
	Code int `json:"code,omitempty"`
	Data string `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	WalletAddress string `json:"walletAddress,omitempty"`
}


type NodeState struct {
	Id      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error *RPCErrorState `json:"error,omitempty"`
	Result  struct {
		Addr               string `json:"addr"`
		Currtimestamp      int    `json:"currTimeStamp"`
		Height             int    `json:"height"`
		ID                 string `json:"id"`
		Jsonrpcport        int    `json:"jsonRpcPort"`
		ProposalSubmitted  int    `json:"proposalSubmitted"`
		ProtocolVersion    int    `json:"protocolVersion"`
		PublicKey          string `json:"publicKey"`
		RelayMessageCount  uint64 `json:"relayMessageCount"`
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