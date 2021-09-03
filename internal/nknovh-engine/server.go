package nknovh_engine

import (
		"net/http"
		"templater"
		"github.com/gobwas/ws"
		"github.com/gobwas/ws/wsutil"
		"github.com/julienschmidt/httprouter"
		"fmt"
		"log"
		"encoding/json"
		"strconv"
		)


type WebHelper struct {
	temp *templater.Templater
}

type Web struct {
	Response map[int]WSReply
	Helper *WebHelper
	Methods map[string]func(*WSQuery,*CLIENT) (error, WSReply)
	MethodsReqAuth []string
}

type WSQuery struct {
	Method string `json:"Method"`
	Value map[string]interface{} `json:"Value, omitempty"`
}
		
type WSReply struct {
	Method string `json:"Method"`
	Code int `json:"Code"`
	Error bool `json:"Error, omitempty`
	ErrMessage string `json:"ErrMessage, omitempty"`
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
}

func (o *NKNOVH) RegisterMethods() {
	o.Web.Methods = map[string]func(*WSQuery, *CLIENT) (error, WSReply){}
	o.Web.MethodsReqAuth = []string{"getfullstack", "addnodes", "rmnodes", "getmynodes", "getnetstatus", "getmywallets", "getprices", "savemysettings"}

	o.Web.Methods["auth"] = o.apiAuth
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

func (o *NKNOVH) WsPolling(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		o.log.Syslog("bad \"Upgrade\" header", "wshttp")
		return
	}
	go func() {
		defer conn.Close()
		ip, err := o.getIp(o.conf.TrustedProxies, r)
		if err != nil {
			o.log.Syslog("getIp returned an error: " + err.Error(), "wshttp")
			return
		}
		c := &CLIENT{HashId: -1, Ip: ip}
		for {
			msg, op, err := wsutil.ReadClientData(conn)
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
				if c.HashId == -1 {
					//Authorization needed
					res := o.Web.Response[253]
					res.Method = q.Method
					if b, err := json.Marshal(res); err == nil {
						wsutil.WriteServerMessage(conn, op, b)
						continue	
					}
					return
				}
			}

			o.updateUniqWatch(c)
			_, res := o.Web.Methods[q.Method](q, c)
			if b, err := json.Marshal(res); err == nil {
				wsutil.WriteServerMessage(conn, op, b)
				continue
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

	c := &CLIENT{HashId: -1, Ip: ip}

	var hash string
	var ok bool
	data := new(WSQuery)
	value := map[string]interface{}{}

	// application/x-www-form-urlencoded
	err = r.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	for key, val := range r.Form {
		value[key] = val[0]
	}

	if x := params.ByName("method"); x != "" {
		if _, ok = o.Web.Methods[x]; !ok {
			w.Write([]byte("No method found"))
			return
		}
		data.Method = x
	} else {
		if val, ok := value["Method"].(string); ok {
			if _, ok = o.Web.Methods[val]; !ok {
				w.Write([]byte("No method found"))
				return
			}
			data.Method = val
			delete(value, "Method")
		} else {
			w.Write([]byte("No method value"))
			return
		}
	}
	o.log.Syslog("POST Request from " + c.Ip, "http")
	//Auth
	if i := FindStringInSlice(o.Web.MethodsReqAuth, data.Method); i != len(o.Web.MethodsReqAuth) {

		if hash, ok = value["Hash"].(string); !ok {
			//DEPRECATED AND WILL BE REMOVED
			if hash, ok = value["hash"].(string); !ok {
				res := o.Web.Response[253]
				var b = make([]byte, 0)
				if b, err = json.Marshal(res); err != nil {
					o.InternalErrorJson(w, err)
					return
				}
				w.Write(b)
				return
			}
		}

		o.updateUniqWatch(c)

		mauth := map[string]interface{}{}
		mauth["Hash"] = hash
		auth := &WSQuery{Method: "auth", Value: mauth, }

		_, reply := o.apiAuth(auth, c)

		if reply.Code != 0 {
			var b = make([]byte, 0)
			if b, err = json.Marshal(reply); err != nil {
				o.InternalErrorJson(w, err)
				return
			}
			w.Write(b)
			return
		}
		delete(value, "Hash")
	}

	data.Value = value
	_, reply := o.Web.Methods[data.Method](data, c)
	var b = make([]byte, 0)
	if b, err = json.Marshal(reply); err != nil {
		o.InternalErrorJson(w, err)
		return
	}
	w.Write(b)
	return 
}

func (o *NKNOVH) Listen() {
	x := templater.NewTemplater("templates")
	wh := &WebHelper{temp: x}
	o.Web = &Web{Response: map[int]WSReply{}, Helper: wh}
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
