package nknovh_engine

import (
		"net/http"
		"templater"
		"github.com/gobwas/ws"
		"github.com/gobwas/ws/wsutil"
		"github.com/julienschmidt/httprouter"
		"fmt"
		"sync"
		"net"
		"log"
		"encoding/json"
		"strconv"
		"errors"
		"time"
		)


type WebHelper struct {
	temp *templater.Templater
}

type Web struct {
	Response map[int]WSReply
	Helper *WebHelper
	Methods map[string]func(*WSQuery,*CLIENT) (error, WSReply)
	MethodsReqAuth []string
	MethodsReadOnly []string
	MethodsToAll []string
	WsPool *WsPool
}

type WsPool struct {
	Clients map[int]*WsClients
	ActiveIps map[string]int
	i uint64
	mu sync.RWMutex
	mu_ips sync.RWMutex
}

type WsClients struct {
	list map[uint64]*CLIENT
	mu sync.RWMutex
}

type WSQuery struct {
	Method string `json:"Method"`
	Value map[string]interface{} `json:"Value,omitempty"`
}
		
type WSReply struct {
	Method string `json:"Method"`
	Code int `json:"Code"`
	Error bool `json:"Error,omitempty`
	ErrMessage string `json:"ErrMessage,omitempty"`
	Value interface{}
}

type Netstatus struct {
	Relays              int64   `json:"relays"`
 	AverageUptime       int     `json:"average_uptime"`
	AverageRelays       uint64  `json:"average_relays"`
	RelaysPerHour       uint64  `json:"relays_per_hour"`
	ProposalSubmitted   int     `json:"proposalSubmitted"`
	PersistNodesCount   int     `json:"persist_nodes_count"`
	NodesCount          int     `json:"nodes_count"`
 	LastHeight          int     `json:"last_height"`
 	LastTimestamp       uint64  `json:"last_timestamp"`
 	AverageBlockTime    float64 `json:"average_blockTime"`
	AverageBlocksPerDay float64 `json:"average_blocksPerDay"`
	LatestUpdate        string  `json:"latest_update"`
}


type CLIENT struct {
	HashId int
	Ip string
	ReadOnly bool
	ConnId uint64
	NotWs bool
	WsConnection net.Conn
}

func (o *NKNOVH) RegisterMethods() {
	o.Web.Methods = map[string]func(*WSQuery, *CLIENT) (error, WSReply){}
	o.Web.MethodsReqAuth = []string{"getfullstack", "addnodes", "rmnodes", "getmynodes", "getnetstatus", "getmywallets", "getprices", "getnodedetails", "savemysettings", "logout"}
	o.Web.MethodsReadOnly = []string{"getfullstack", "getmynodes", "getnetstatus", "getmywallets", "getprices", "getnodedetails"}
	o.Web.MethodsToAll = []string{"addnodes", "rmnodes", "savemysettings"}
	o.Web.Methods["auth"] = o.apiAuth
	o.Web.Methods["logout"] = o.apiLogout
	o.Web.Methods["genid"] = o.apiGenId
	o.Web.Methods["getfullstack"] = o.apiFullstack 
	o.Web.Methods["addnodes"] = o.apiAddNodes
	o.Web.Methods["getmywallets"] = o.apiMyWallets
	o.Web.Methods["getprices"] = o.apiPrices
	o.Web.Methods["getmynodes"] = o.apiMyNodes
	o.Web.Methods["getnetstatus"] = o.apiNetstatus
	o.Web.Methods["rmnodes"] = o.apiRmNodes
	o.Web.Methods["getdaemon"] = o.apiDaemon
	o.Web.Methods["getlanguage"] = o .apiLanguage
	o.Web.Methods["savemysettings"] = o.apiSaveSettings
	o.Web.Methods["getnodedetails"] = o.apiGetNodeDetails
	o.Web.Methods["getnodeipbypublickey"] = o.apiGetNodeIpByPublicKey
	return
}

func (o *NKNOVH) RegisterResponse() {
	o.Web.Response[1] = WSReply{Code: 1, Error: true, ErrMessage: "Cannot execute SQL query"}
	o.Web.Response[2] = WSReply{Code: 2, Error: false, ErrMessage: "Nodes added partially"}
	o.Web.Response[3] = WSReply{Code: 3, Error: false, ErrMessage: "No info/entries in a database"}
	o.Web.Response[4] = WSReply{Code: 4, Error: true, ErrMessage: "An error occured"}
	o.Web.Response[5] = WSReply{Code: 5, Error: true, ErrMessage: "Incorrect query"}
	o.Web.Response[6] = WSReply{Code: 6, Error: true, ErrMessage: "Name of node(-s) too long or incorrect format."}
	o.Web.Response[7] = WSReply{Code: 7, Error: true, ErrMessage: "Multiple variable must be boolean"}
	o.Web.Response[8] = WSReply{Code: 8, Error: true, ErrMessage: "Ip variable must be valid string"}
	o.Web.Response[9] = WSReply{Code: 9, Error: true, ErrMessage: "Incorrect ip address(-es)"}
	o.Web.Response[10] = WSReply{Code: 10, Error: true, ErrMessage: "Passed ip address(-es) not IPv4"}
	o.Web.Response[11] = WSReply{Code: 11, Error: true, ErrMessage: "Passed ip address(-es) not in public network"}
	o.Web.Response[12] = WSReply{Code: 12, Error: true, ErrMessage: "Nodes limit exceeded"}
	o.Web.Response[13] = WSReply{Code: 13, Error: true, ErrMessage: "Wrong delimiter passed"}
	o.Web.Response[14] = WSReply{Code: 14, Error: true, ErrMessage: "The nodes weren't added since they had been already created."}
	o.Web.Response[15] = WSReply{Code: 15, Error: true, ErrMessage: "Wrong data of NodesId passed"}
	o.Web.Response[16] = WSReply{Code: 16, Error: true, ErrMessage: "Wallets overflow"}
	o.Web.Response[17] = WSReply{Code: 17, Error: true, ErrMessage: "One or more of the passed wallets are not in the correct format"}
	o.Web.Response[18] = WSReply{Code: 18, Error: true, ErrMessage: "One or more Id of passed nodes are not found. No changes."}

	//Link to apiGetNodeDetails
	o.Web.Response[19] = WSReply{Code: 19, Error: true, ErrMessage: "Wrong data of NodeId passed"}
	o.Web.Response[20] = WSReply{Code: 20, Error: true, ErrMessage: "The node is offline / No reply from the node before a timeout"}
	o.Web.Response[21] = WSReply{Code: 21, Error: true, ErrMessage: "Cannot decode json of the node response (getnodestate)"}
	o.Web.Response[22] = WSReply{Code: 22, Error: true, ErrMessage: "The node is online, but information about neighbors has not been received before a timeout"}
	o.Web.Response[23] = WSReply{Code: 23, Error: true, ErrMessage: "Cannot decode json of the node response (getneighbor)"}
	o.Web.Response[24] = WSReply{Code: 24, Error: true, ErrMessage: "Query returned an error (getneighbor)"}

	//Link to apiGetNodeIpByPublicKey
	o.Web.Response[25] = WSReply{Code: 25, Error: true, ErrMessage: "PublicKey is not set"}
	o.Web.Response[26] = WSReply{Code: 26, Error: true, ErrMessage: "Wrong PublicKey passed"}

	o.Web.Response[230] = WSReply{Code: 230, Error: true, ErrMessage: "No view variable passed, the variable must be string"}
	o.Web.Response[231] = WSReply{Code: 231, Error: true, ErrMessage: "No Locale variable passed, the variable must be string"}
	o.Web.Response[232] = WSReply{Code: 232, Error: true, ErrMessage: "Locale or View passed variables are overflow"}
	o.Web.Response[233] = WSReply{Code: 233, Error: true, ErrMessage: "Passed lang package is not found in package slice"}
	o.Web.Response[234] = WSReply{Code: 234, Error: true, ErrMessage: "Passed language pack is not found as JSON file"}

	o.Web.Response[240] = WSReply{Code: 240, Error: true, ErrMessage: "GenRandomSHA256 returned error"}
	o.Web.Response[252] = WSReply{Code: 252, Error: true, ErrMessage: "You have created at least 3 ID for the latest 30 minutes"}
	o.Web.Response[253] = WSReply{Code: 253, Error: true, ErrMessage: "You have no authorization"}
	o.Web.Response[254] = WSReply{Code: 254, Error: true, ErrMessage: "Incorrect ID length"}
	o.Web.Response[255] = WSReply{Code: 255, Error: true, ErrMessage: "Passed ID is not found"}
	o.Web.Response[500] = WSReply{Code: 500, Error: true, ErrMessage: "Internal server error"}

	//Main errors
	o.Web.Response[1000] = WSReply{Code: 1000, Error: true, ErrMessage: "Method variable is not passed or it has wrong format"}
	o.Web.Response[1001] = WSReply{Code: 1001, Error: true, ErrMessage: "The passed Method is not found"}
	o.Web.Response[1002] = WSReply{Code: 1002, Error: true, ErrMessage: "Connections limit is reached"}
	o.Web.Response[1003] = WSReply{Code: 1003, Error: true, ErrMessage: "Passed JSON is incorrect"}
	return
}

func (o *NKNOVH) InternalErrorJson(w http.ResponseWriter, errx error) {
	w.WriteHeader(500)
	res := o.Web.Response[500]
	if errx != nil {
		o.log.Syslog("Internal server error: " + errx.Error(), "http")
	}
	if b, err := json.Marshal(res); err == nil {
		w.Write(b)
	}
	return
}


func (o *NKNOVH) WsClientCreate(conn net.Conn) *CLIENT {
	t := time.Now()

	o.Web.WsPool.mu.Lock()
	o.Web.WsPool.i++
	c := &CLIENT{HashId: -1, WsConnection: conn, ConnId: o.Web.WsPool.i}
	if v, ok := o.Web.WsPool.Clients[c.HashId]; !ok {
		o.Web.WsPool.Clients[c.HashId] = new(WsClients)
		o.Web.WsPool.Clients[c.HashId].list = map[uint64]*CLIENT{}
		o.Web.WsPool.Clients[c.HashId].list[o.Web.WsPool.i] = c
		o.Web.WsPool.mu.Unlock()
	} else {
		o.Web.WsPool.mu.Unlock()
		v.mu.Lock()
		v.list[o.Web.WsPool.i] = c
		v.mu.Unlock()
	}
	t_x := time.Now().Sub(t).String()
	o.log.Syslog("WsClientCreate time: " + t_x, "debug")

	//ONLY FOR DEBUG MODE
	cnt := 0
	o.Web.WsPool.mu.RLock()
	for i, _ := range o.Web.WsPool.Clients {
		o.Web.WsPool.Clients[i].mu.RLock()
		cnt += len(o.Web.WsPool.Clients[i].list)
		o.Web.WsPool.Clients[i].mu.RUnlock()
	}
	o.Web.WsPool.mu.RUnlock()
	o.log.Syslog("Active ws connections: " + strconv.Itoa(cnt), "debug")
	return c
}

func (o *NKNOVH) WsClientClose(c *CLIENT) {
	t := time.Now()
	o.Web.WsPool.Clients[c.HashId].mu.Lock()
	c.WsConnection.Close()
	o.Web.WsPool.Clients[c.HashId].mu.Unlock()
	o.WsMultiConnectDecrease(c.Ip)
	o.WsClientGC(c)
	t_x := time.Now().Sub(t).String()
	o.log.Syslog("WsClientClose time: " + t_x, "debug")
	return
}

func (o *NKNOVH) WsClientGC(c *CLIENT) {
	o.Web.WsPool.Clients[c.HashId].mu.Lock()
	delete(o.Web.WsPool.Clients[c.HashId].list, c.ConnId)
	if len(o.Web.WsPool.Clients[c.HashId].list) == 0 {
		o.Web.WsPool.mu.Lock()
		delete(o.Web.WsPool.Clients, c.HashId)
		o.Web.WsPool.mu.Unlock()
	} else {
		o.Web.WsPool.Clients[c.HashId].mu.Unlock()
	}
	return
}

func (o *NKNOVH) WsClientUpdate(c *CLIENT, hashId int) {
	t := time.Now()
	if c.NotWs {
		c.HashId = hashId
		return
	}
	o.WsClientGC(c)

	debugf := func() {
			for x, _ := range o.Web.WsPool.Clients {
				for i, _ := range o.Web.WsPool.Clients[x].list {
					s := fmt.Sprintf("HashId: %v >> Context: %v", x, o.Web.WsPool.Clients[x].list[i])
					o.log.Syslog(s, "debug")
				}
			}
	}
	c.HashId = hashId
	o.Web.WsPool.mu.Lock()
	if v, ok := o.Web.WsPool.Clients[c.HashId]; !ok {
		o.Web.WsPool.Clients[c.HashId] = new(WsClients)
		o.Web.WsPool.Clients[c.HashId].list = map[uint64]*CLIENT{}
		o.Web.WsPool.Clients[c.HashId].list[c.ConnId] = c
		debugf()
		o.Web.WsPool.mu.Unlock()
	} else {
		debugf()
		o.Web.WsPool.mu.Unlock()
		v.mu.Lock()
		v.list[c.ConnId] = c
		v.mu.Unlock()
	}
	t_x := time.Now().Sub(t).String()
	o.log.Syslog("WsClientUpdate time: " + t_x, "debug")

	return
}

func (o *NKNOVH) WsRestrictMultiConnect(ip string) (error, WSReply) {
	limit := 10
	var ok bool
	o.Web.WsPool.mu_ips.Lock()
	defer o.Web.WsPool.mu_ips.Unlock()
	if _, ok = o.Web.WsPool.ActiveIps[ip]; !ok {
		o.Web.WsPool.ActiveIps[ip] = 1
		return nil, WSReply{}
	} else {
		o.Web.WsPool.ActiveIps[ip] = o.Web.WsPool.ActiveIps[ip]+1
		if o.Web.WsPool.ActiveIps[ip] > limit {
			o.log.Syslog("Connections limit is reached from IP " + ip, "debug")
			q := new(WSQuery)
			q.Method = "other"
			_, wsreply := o.WsError(q, 1002)
			return errors.New("Connections limit is reached"), wsreply
		} else {
			return nil, WSReply{}
		}
	}
}

func (o *NKNOVH) WsMultiConnectDecrease(ip string) {
	var ok bool
	o.Web.WsPool.mu_ips.Lock()
	defer o.Web.WsPool.mu_ips.Unlock()
	if _, ok = o.Web.WsPool.ActiveIps[ip]; !ok {
		return
	}
	x := o.Web.WsPool.ActiveIps[ip] - 1
	if x == 0 {
		delete(o.Web.WsPool.ActiveIps, ip)
		return
	}
	o.Web.WsPool.ActiveIps[ip] = x
}

func (o *NKNOVH) WsPolling(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		o.log.Syslog("bad \"Upgrade\" header", "wshttp")
		return
	}
	go func() {
		c := o.WsClientCreate(conn)
		defer o.WsClientClose(c)
		ip, err := o.getIp(o.conf.TrustedProxies, r)
		if err != nil {
			o.log.Syslog("getIp returned an error: " + err.Error(), "wshttp")
			return
		}
		c.Ip = ip
		if err, wsreply := o.WsRestrictMultiConnect(ip); err != nil {
			o.WriteJsonWs(&wsreply, c)
			return
		}
		for {
			msg, _, err := wsutil.ReadClientData(conn)
			if err != nil {
				o.log.Syslog(err.Error(), "wshttp")
				return
			}
			o.log.Syslog("WS Request from " + c.Ip + "; Message: " + string(msg), "wshttp")

			q := new(WSQuery)
			if err := json.Unmarshal(msg, q); err != nil {
				o.log.Syslog("Cannot unmarshal json to WSQuery: " + err.Error(), "errors")
				return
			}
			if _, ok := o.Web.Methods[q.Method]; !ok {
				return
			}
			if i := FindStringInSlice(o.Web.MethodsReqAuth, q.Method); i != len(o.Web.MethodsReqAuth) {
				if c.HashId == -1 && c.ReadOnly == false {
					//Authorization needed
					res := o.Web.Response[253]
					res.Method = q.Method
					if err := o.WriteJsonWs(&res, c); err == nil {
						continue
					}
					return
				}
			}

			o.updateUniqWatch(c)
			_, res := o.Web.Methods[q.Method](q, c)
			if i := FindStringInSlice(o.Web.MethodsToAll, q.Method); i != len(o.Web.MethodsToAll) {
				if err = o.WsSendByHashId(&res, c.HashId); err == nil {
					continue
				}
			} else {
				if err := o.WriteJsonWs(&res, c); err == nil {
					continue
				}
			}
			return
		}
	}()
}

func (o *NKNOVH) CreateIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t := o.Web.Helper.temp.New("index")
	if !t.IsComplete() {
		if err := t.GetPage("header", "main"); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		if err := t.GetPage("index", "pages"); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		if err := t.GetPage("footer", "main"); err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		t.Complete()
	}
	t.Flush()
	if err, x := getEtag("web/static/css/nknc.css"); err == nil {
		t.Set("style_etag", strconv.FormatInt(x, 10))
	}
	if err, x := getEtag("web/static/js/wasm_exec.js"); err == nil {
		t.Set("wexec_etag", strconv.FormatInt(x, 10))
	}
	if err, x := getEtag("web/static/lib.wasm"); err == nil {
		t.Set("wasm_etag", strconv.FormatInt(x, 10))
	}
    w.Write(t.View())
}

func (o *NKNOVH) apiPOST(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	ip, err := o.getIp(o.conf.TrustedProxies, r)
	if err != nil {
		o.InternalErrorJson(w, err)
		o.log.Syslog("getIp returned an error: " + err.Error(), "http")
		return
	}
	c := &CLIENT{HashId: -1, Ip: ip, NotWs: true,}
	o.log.Syslog("POST Request from " + c.Ip, "http")

	var hash string
	var ok bool
	data := new(WSQuery)
	value := map[string]interface{}{}

	// application/x-www-form-urlencoded
	err = r.ParseForm()
	if err != nil {
		return
	}
	for key, val := range r.Form {
		value[key] = val[0]
	}

	if len(value) == 0 {
		err = json.NewDecoder(r.Body).Decode(data)
		if err != nil {
			_, wsreply := o.WsError(data, 1003)
			o.WriteJson(&wsreply, w)
			return
		}
		if x := params.ByName("method"); x != "" {
			if _, ok = o.Web.Methods[x]; !ok {
				_, wsreply := o.WsError(data, 1001)
				o.WriteJson(&wsreply, w)
				return
			}
			data.Method = x
		}
		if _, ok = o.Web.Methods[data.Method]; !ok {
			_, wsreply := o.WsError(data, 1001)
			o.WriteJson(&wsreply, w)
			return
		}
	} else {
		if x := params.ByName("method"); x != "" {
			if _, ok = o.Web.Methods[x]; !ok {
				_, wsreply := o.WsError(data, 1001)
				o.WriteJson(&wsreply, w)
				return
			}
			data.Method = x
		} else if val, ok := value["Method"].(string); ok {
			if _, ok = o.Web.Methods[val]; !ok {
				_, wsreply := o.WsError(data, 1001)
				o.WriteJson(&wsreply, w)
				return
			}
			data.Method = val
			delete(value, "Method")
		} else {
			_, wsreply := o.WsError(data, 1000)
			o.WriteJson(&wsreply, w)
			return
		}
	}

	//Auth
	if i := FindStringInSlice(o.Web.MethodsReqAuth, data.Method); i != len(o.Web.MethodsReqAuth) {
		if hash, ok = value["Hash"].(string); !ok {
			if hash, ok = data.Value["Hash"].(string); !ok {
				res := o.Web.Response[253]
				o.WriteJson(&res, w)
				return
			}
		}

		o.updateUniqWatch(c)

		mauth := map[string]interface{}{}
		mauth["Hash"] = hash
		auth := &WSQuery{Method: "auth", Value: mauth, }

		_, reply := o.apiAuth(auth, c)

		if reply.Code != 0 {
			o.WriteJson(&reply, w)
			return
		}
		if _, ok = value["Hash"]; ok {
			delete(value, "Hash")
		}
	}
	if len(value) > 0 {
		data.Value = value
	}
	_, reply := o.Web.Methods[data.Method](data, c)
	o.WriteJson(&reply, w)
	return 
}


func (o *NKNOVH) WriteJson(data *WSReply, w http.ResponseWriter) error {
	var b = make([]byte, 0)
	var err error
	if b, err = json.Marshal(data); err != nil {
		o.InternalErrorJson(w, err)
		return err
	}
	w.Write(b)
	return nil
}


func (o *NKNOVH) WriteJsonWs(data *WSReply, c *CLIENT) error {
	var b = make([]byte, 0)
	var err error
	if b, err = json.Marshal(data); err == nil {
		if err = wsutil.WriteServerMessage(c.WsConnection, ws.OpText, b); err != nil {
			return err
		}
		return nil
	}
	return err
}

func (o *NKNOVH) Listen() {
	x := templater.NewTemplater("templates")
	wh := &WebHelper{temp: x}
	o.Web = &Web{Response: map[int]WSReply{}, Helper: wh}
	o.Web.WsPool = new(WsPool)
	o.Web.WsPool.Clients = map[int]*WsClients{}
	o.Web.WsPool.ActiveIps = map[string]int{}

	o.RegisterResponse()
	o.RegisterMethods()

	router := httprouter.New()
	router.GET("/", o.CreateIndex)
	//DEPRECATED LINE
	router.GET("/id/:hash", o.CreateIndex)
	router.GET("/login/auth/:hash", o.CreateIndex)
	router.GET("/polling", o.WsPolling)
	router.POST("/api/:method", o.apiPOST)
	router.POST("/api", o.apiPOST)

	router.NotFound = http.FileServer(http.Dir("./web/"))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(o.conf.HttpServer.Port), router))
}
