package nknovh_wasm

import (
		"syscall/js"
		"encoding/json"
		"fmt"
		"time"
		"sync"
)

type WSQuery struct {
	Method string `json:"Method"`
	Value map[string]interface{} `json:"Value, omitempty"`
}
		
type WSReply struct {
	Method string `json:"Method"`
	Code int `json:"Code"`
	Error bool `json:"Error, omitempty`
	ErrMessage string `json:"ErrMessage, omitempty"`
	Value map[string]interface{} `json: "Value, omitempty"`
}

func (c *CLIENT) WsOnOpen() {
	f := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		if c.Hash != "" {
			c.SetLanguage("view_src", "")
			c.WsAuth()
			c.mux.AutoUpdater.Lock()
			if b := c.AutoUpdaterIsStarted; !b {
				fmt.Println("AutoUpdater is stopped. Sending a message to run")

				c.AutoUpdaterStartCh <- true
			}
			c.mux.AutoUpdater.Unlock()
		} else {
			c.SetLanguage("index_src", "")
		}
		return nil
	})
	c.ws.Set("onopen", f)
	fmt.Println("WebSocket onopen registered")
	return
}

func (c *CLIENT) SetLanguage(view string, locale string) {
	c.mux.StartView.Lock()
	c.switchLoading(true)
	if locale != "" {
		c.W.LocalStorage("set", "lang", locale)
	} else {
		if err, storage_lang := c.W.LocalStorage("get", "lang"); err == nil {
			locale = storage_lang 
		} else {
			locale = c.Conf.DefaultLanguage
		}
	}

	c.Lang = locale
	if _, ok := c.Cached.Lang[locale]; !ok {
		m := map[string]interface{}{}
		m["Locale"] = locale
		m["View"] = view
		c.WsSend("getlanguage", m)
	} else {
		c.LANG = c.Cached.Lang[locale]
		c.handlingLangPages(view, locale)
		c.mux.StartView.Unlock()
	}
}

func (c *CLIENT) WsOnClose() {
	check := make(chan int)
	ch := make(chan int)

	reconnect_notification := func(rchan chan int) {
		doc := js.Global().Get("document")
		elem := doc.Call("getElementById", "disconnected")
		if !elem.Truthy() {
			close(rchan)
			return
		}
		sec_elem := doc.Call("getElementById", "disconnected-seconds")
		i := 0
		var last int
		for {
			select {
			case x := <-rchan:
				if last < x {
					c.W.HideById("disconnected-process")
					c.W.ShowById("disconnected-failed")
				} else {
					c.W.HideById("disconnected-failed")
					c.W.ShowById("disconnected-process")
				}
				if x < -3 {
					c.W.HideById("disconnected-failed")
					c.W.HideById("disconnected-process")
					c.W.ShowById("disconnected-success")
					//time.Sleep(2 * time.Second)
					c.W.HideById("disconnected")
					c.W.HideById("disconnected-success")
					return
				}
				if i == 0 {
					c.W.ShowById("disconnected")
					i++
				}
				last = x
				sec_elem.Set("innerHTML", x)

			break
			}
		}
	}

	checkstatus := func() {
		seconds_to_reconnect := 6
		str_mux := &sync.Mutex{}
		rchan := make(chan int, 0)
		go reconnect_notification(rchan)
		rchan_tick := time.NewTicker(50 * time.Millisecond)
		tick := time.NewTicker(50 * time.Millisecond)
		tick_checker := time.NewTicker(5 * time.Millisecond)
		tick_checker.Stop()
		rchan_tick.Stop()
		defer tick.Stop()
		defer rchan_tick.Stop()
		defer tick_checker.Stop()
		for {
			select {
			case <-rchan_tick.C:
				str_mux.Lock()
				seconds_to_reconnect--
				rchan <- seconds_to_reconnect
				str_mux.Unlock()
			case <-tick.C:
				str_mux.Lock()
				seconds_to_reconnect = 6
				str_mux.Unlock()

				rchan_tick.Reset(1 * time.Second)
				tick.Reset(5 * time.Second)
				x := c.ws.Get("readyState").Int()
				ch <- x
			case y := <-check:
				if y == 0 {
					rchan_tick.Stop()
					rchan <- -10
			
					return
				}
				if y == 1 {
					x := c.ws.Get("readyState").Int()
					if x == 0 {
						c.mux.Websocket.Lock()
						c.WsOnOpen()
						c.WsOnMessage()
						c.mux.Websocket.Unlock()
						tick_checker.Reset(1 * time.Millisecond)
						looping:
						for {
							select {
								case <-tick_checker.C:
								x := c.ws.Get("readyState").Int()
								if x == 1 {
									ch <- x
									break looping
								}
								if x == 2 || x == 3 {
									break looping
								}
							}
						}
					}
					if x == 1 {
						continue
					}
				}
			}
		}
	}

	f := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		go func() {
			c.mux.AutoUpdater.Lock()
			if c.AutoUpdaterIsStarted {
				c.AutoUpdaterStopCh <- true
			}
			c.mux.AutoUpdater.Unlock()
			go checkstatus()
			for {
				select {
					case x := <-ch:
					if x == 1 {
						c.mux.Websocket.Lock()
						c.WsOnClose()
						c.mux.Websocket.Unlock()
						check <- 0
						return
					}
					if x == 3 || x == 2 {
						c.mux.Websocket.Lock()
						_, c.ws = c.W.WsCreate("polling")
						c.mux.Websocket.Unlock()
						check <- 1
						continue
					}
				}
			}	
		}()
		return nil
	})
	c.ws.Set("onclose", f)
	fmt.Println("WebSocket onclose registered")
}

func (c *CLIENT) WsOnMessage() {
	f := js.FuncOf(func(_ js.Value, y []js.Value) interface{} {
		inJson := []byte(y[0].Get("data").String())
		data := new(WSReply)
		if err := json.Unmarshal(inJson, data); err != nil {
			fmt.Println(err.Error())
			return false
		}
		if _, ok := c.apiMethods[data.Method]; !ok {
			fmt.Println("Method", data.Method, "not found")
			return false
		}

		c.apiMethods[data.Method](data)
		return nil
	})
	c.ws.Set("onmessage", f)
	fmt.Println("WebSocket onmessage registered")
	return
}

func (c *CLIENT) WsAuth(h ...string) {
	var hash string
	if len(h) > 0 {
		hash = h[0]
	} else {
		hash = c.Hash
	}
	data := new(WSQuery)
	data.Method = "auth"
	data.Value = map[string]interface{}{}
	data.Value["Hash"] = hash
	b, _ := json.Marshal(data)
	c.ws.Call("send", string(b))
}

func (c *CLIENT) WsGetFullstack() {
	if c.Hash == "" {
		c.GenErr("You have no authorization, please reload the page", "default", -1)
		return
	}
	data := new(WSQuery)
	data.Method = "getfullstack"
	b, _ := json.Marshal(data)
	c.ws.Call("send", string(b))
	return
}

func (c *CLIENT) WsAddNodes(datanodes map[string]interface{}) {
	data := new(WSQuery)
	data.Method = "addnodes"
	data.Value = datanodes
	b, _ := json.Marshal(data)
	c.ws.Call("send", string(b))
}

func (c *CLIENT) WsGenId() {
	data := new(WSQuery)
	data.Method = "genid"
	b, _ := json.Marshal(data)
	c.ws.Call("send", string(b)) 
	return
}

func (c *CLIENT) WsSend(method string, x ...map[string]interface{}) {
	c.mux.Websocket.Lock()
	defer c.mux.Websocket.Unlock()
	data := new(WSQuery)
	data.Method = method
	if len(x) > 0 {
		data.Value = x[0]
	}
	b, _ := json.Marshal(data)
	c.ws.Call("send", string(b))
	return
}
