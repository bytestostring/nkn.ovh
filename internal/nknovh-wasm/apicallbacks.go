package nknovh_wasm

import (
		"syscall/js"
		"encoding/json"
		"fmt"
		"strconv"
)

func (c *CLIENT) apiLanguage(data *WSReply) interface{} {
	defer c.mux.StartView.Unlock()

	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		js.Global().Call("alert", data.ErrMessage)
		return false
	}
	locale, ok := data.Value["Locale"].(string)
	if !ok {
		js.Global().Call("alert","No Locale string passed in callback func")
		return false
	}
	view, ok := data.Value["View"].(string)
	if !ok {
		js.Global().Call("alert","No View string passed in callback func")
		return false
	}
	value, ok := data.Value["Data"].(string)
	if !ok {
		js.Global().Call("alert","No Data string passed in callback func")
		return false
	}

	LANG := new(LANG)
	if err := json.Unmarshal([]byte(value), LANG); err != nil {
		fmt.Println("Cannot unmarshal language pack")
		fmt.Println(err.Error())
		js.Global().Call("alert", "Cannot unmarshal language package, please reload the page later")
		return false
	}
	c.LANG = LANG
	c.Cached.Lang[locale] = LANG
	
	//parsing
	c.handlingLangPages(view, locale)
	return true

}

func (c *CLIENT) apiGenId(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	hash, ok := data.Value["Hash"].(string)
	if !ok {
		return false
	}
	c.WsAuth(hash)
	return true
}

func (c *CLIENT) apiAuth(data *WSReply) interface{} {
	if data.Error  {
		if c.Hash != "" {
			c.Hash = ""
			c.W.LocalStorage("remove", "hash")
			c.mux.AutoUpdater.Lock()
			if c.AutoUpdaterIsStarted {
				c.AutoUpdaterStopCh <- true
			}
			c.mux.AutoUpdater.Unlock()
			history := js.Global().Get("history")
			history.Call("pushState", nil, nil, "/")
			go c.Run()
		}

		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}

	if c.Hash == "" {
		hash, ok := data.Value["Hash"].(string)
		if !ok {
			return false
		}
		c.W.LocalStorage("set", "hash", hash)
		go c.Run()
		return true
	}
	return true
}

func (c *CLIENT) apiFullstack(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	if data.Code != 0 {
		return false
	}
	fstack := new(GetFullstack)
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Canot marshal data O.o How?")
		return false
	}
	if err := json.Unmarshal(b, fstack); err != nil {
		fmt.Println("Cannot unmarshal fullstack")
		return false
	}
	c.ParseFullstack(fstack)
	return true
}

func (c *CLIENT) apiAddNodes(data *WSReply) interface{} {
	doc := js.Global().Get("document")
	button := doc.Call("getElementById", "addNodeButton")
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code, "addNodesErr")
		button.Set("disabled", false)
		return false
	}
	c.W.HideById("addNodesErr")
	doc.Call("getElementById", "nodeIP").Set("value", "")
	doc.Call("getElementById", "nodeIPList").Set("value", "")
	doc.Call("getElementById", "nodeName").Set("value", "")
	doc.Call("getElementById", "nodeNameList").Set("value", "")
	button.Set("disabled", false)

	c.ShowHideModal("addNodes", "hide")
	/*
	c.W.ShowById("completedQuery")
	js.Global().Call("setTimeout", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.W.HideById("completedQuery")
		return nil
	}), 1500)
	*/
	//partial
	c.WsSend("getmynodes")
	return true
}

func (c *CLIENT) apiRmNodes(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}

	var raw_nodes []interface{}
	var nodes []int
	var ok bool
	var node_id float64
	if raw_nodes, ok = data.Value["NodesId"].([]interface{}); !ok {
		return false
	}

	for i, _ := range raw_nodes {
		if node_id, ok = raw_nodes[i].(float64); !ok {
			return false
		}
		nodes = append(nodes, int(node_id))
	}
	doc := js.Global().Get("document")
	for i, _ := range nodes {
		rmnode := doc.Call("getElementById", "Node-" + strconv.Itoa(nodes[i]))
		c.W.Remove(&rmnode)
	}
	ca := doc.Call("getElementById", "control-all")
	if ca.Truthy() {
		if ca.Get("checked").Bool() == true {
			ca.Set("checked", false)
		}
	}
	c.WsSend("getmynodes")
	return nil
}


func (c *CLIENT) apiMyNodes(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	nodes := new(Nodes)
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, nodes); err != nil {
		fmt.Println("Cannot unmarshal to nodes")
		return false
	}
	switch x := data.Code; x {
		case 0:
			c.W.HideById("nodes_nf")
			c.Nodes = nodes
			c.PreSortNodes()
			c.SortAndParseNodes()
			c.calcNodesSummary()

		break
		case 3:
			c.W.ShowById("nodes_nf")
			c.Nodes = &Nodes{}
			c.PreSortNodes()
			c.SortAndParseNodes()
			c.calcNodesSummary()
		break
	}
	return nil
}

func (c *CLIENT) apiSaveSettings(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code, "settingsErr")
		return false
	}

	c.WsSend("getmywallets")
	c.ShowHideModal("settings", "hide")
	/*
	c.W.ShowById("completedQuery")
	js.Global().Call("setTimeout", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.W.HideById("completedQuery")
		return nil
	}), 1500)
	*/
	return nil
}

func (c *CLIENT) apiMyWallets(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	wallets := new(Wallets)
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, wallets); err != nil {
		fmt.Println("Cannot unmarshal to wallets")
		return false
	}
	fmt.Println("Walletinfoupdate run")
	c.Wallets = wallets
	c.walletsInfoUpdate()
	c.AddWalletLabels(false)
	return nil
}

func (c *CLIENT) apiNetstatus(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	netstatus := new(Netstatus)
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, netstatus); err != nil {
		fmt.Println("Cannot unmarshal to netstatus")
		return false
	}
	c.Netstatus = netstatus

	c.parseNetstatus()
	return nil
}

func (c *CLIENT) apiPrices(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	prices := new(Prices)
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, prices); err != nil {
		fmt.Println("Cannot unmarshal to prices")
		return false
	}
	c.Prices = prices
	return nil
}

func (c *CLIENT) apiDaemon(data *WSReply) interface{} {
	if data.Error {
		c.GenErr(data.ErrMessage, "default", data.Code)
		return false
	}
	daemon := new(Daemon)
	b, _ := json.Marshal(data)
	if err := json.Unmarshal(b, daemon); err != nil {
		fmt.Println("Cannot unmarshal to daemon")
		return false
	}
	c.Daemon = daemon
	if x := c.CheckVersion(c.Daemon.Value.Version); x == false {
		c.mux.AutoUpdater.Lock()
		if c.AutoUpdaterIsStarted {
			c.AutoUpdaterStopCh <- true
		}
		c.mux.AutoUpdater.Unlock()
		return true
	}
	return nil
}