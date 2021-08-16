package xwasmapi

import (
	"fmt"
	"syscall/js"
	//"encoding/json"
	//"io/ioutil"
	//"time"
	//"bytes"
	"errors"
)

type Xwasmapi struct {
	//http_client *http.Client
}

/*
type AjaxAnswer struct {
	Code int
	Error bool `json:"Error, omitempty"`
	ErrMessage string `json:"ErrMessage, omitempty"`
	Value json.RawMessage `json: "Value, omitempty"`
}

type ApiRequest struct {
	Method string
	TypeQuery string
	DataType string
	Data json.RawMessage
}
*/

func New() *Xwasmapi {
	var x Xwasmapi
	//var tp = &http.Transport{DisableKeepAlives: true}
	//x.http_client = &http.Client{Timeout: 3*time.Second, Transport: tp,}
	return &x
}
/*
func (o *Xwasmapi) Ajax(url string, type_query string, datatype string, cb func([]byte, error), raw ...json.RawMessage) (answer []byte, err error) {
	var res *http.Response
	if type_query == "GET" {
		res, err = o.http_client.Get(url)
		cb(answer, err)
		return
	}
	if type_query == "POST" {
		b, err := json.Marshal(raw)
		if err != nil {
			cb(answer, err)
			return answer, err
		}
		if res, err = o.http_client.Post(url, datatype, bytes.NewReader(b)); err != nil {
			cb(answer, err)
			return answer, err
		}
	}
	if res.StatusCode != 200 {
		err = errors.New("StatusCode of the response is not 200")
		cb(answer, err)
		return
	}
	defer res.Body.Close()
	answer, err = ioutil.ReadAll(res.Body)
	if err != nil {
		cb(answer, err)
		return
	}
	cb(answer, err)
	return
}


func (o *Xwasmapi) apiQuery(r *ApiRequest, cb func([]byte, error)) (a *AjaxAnswer, err error) {
	a = new(AjaxAnswer)
	err = nil
	bytes, err := o.Ajax("/api/" +r.Method, r.TypeQuery, r.DataType, cb, r.Data)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, a)
	if err != nil {
		return
	}
	return
}

*/
func (o *Xwasmapi) HideById(id string) error {
	d := js.Global().Get("document").Call("getElementById", id)
	if !d.Truthy() {
		fmt.Println("No access to div#" + id)
		return errors.New("No access to div#" + id)
	}
	d.Get("style").Set("display", "none")
	return nil
}


func (o *Xwasmapi) ShowById(id string) error {
	d := js.Global().Get("document").Call("getElementById", id)
	if !d.Truthy() {
		fmt.Println("No access to div#" + id)
		return errors.New("No access to div#" + id)
	}
	d.Get("style").Set("display", "block")
	return nil
}

func(o *Xwasmapi) Redirect(url string) {
	js.Global().Get("location").Set("href", url)
	return
}

func (o *Xwasmapi) Detach(id string) (error, js.Value) {
	d := js.Global().Get("document").Call("getElementById", id)
	if !d.Truthy() {
		fmt.Println("No access to div#" + id)
		return errors.New("No access to div#" + id), js.Value{}
	}
	r := d.Get("parentElement").Call("removeChild", d)
	return nil, r
}

func (o *Xwasmapi) Remove(dom *js.Value) error {
	if !dom.Truthy() {
		fmt.Println("No access to DOM Element")
		return errors.New("No access to DOM Element")
	}
	dom.Call("remove")
	return nil
}

func (o *Xwasmapi) Get(a ...interface{}) (err error, elem *js.Value) {
	l := len(a)
	if len(a) < 1 {
		return errors.New("Please provide a path to DOM Element"), &js.Value{}
	}
	dom := js.Global()
	for i := 0; i < l; i++ {
		dom = dom.Get(a[i].(string))
	}
	if !dom.Truthy() {
		return errors.New("DOM Element not Truthy"), &js.Value{}
	}
	return nil, &dom
}

func (o *Xwasmapi) LocalStorage(type_call string, params ...string) (err error, s string) {
	l := len(params)
	switch x := type_call; x {
		case "get":
			if l < 1 {
				return errors.New("You must set a variable name"), s
			}
			s = js.Global().Get("localStorage").Call("getItem", params[0]).String()
			if s == "<null>" {
				return errors.New("No variable found"), s
			} else {
				return
			}
		break
		case "set": 
			if l != 2 {
				return errors.New("You must set two parameters"), s
			}
			js.Global().Get("localStorage").Call("setItem", params[0], params[1])
			return
		break
		case "remove":
			if l > 0 {
				for i := 0; i < l; i++ {
					js.Global().Get("localStorage").Call("removeItem", params[i])
				}
				return
			}
			return errors.New("You must set the variable(s) for remove Method"), s
		break
		case "clear":
			js.Global().Get("localStorage").Call("clear")
			return
		default:
		return errors.New("You must set the type_call variable to get, set, remove or clear"), s
	}

	return
}

func (o *Xwasmapi) WsCreate(path string) (error, *js.Value) {
	host := js.Global().Get("location").Get("host")
	protocol := js.Global().Get("location").Get("protocol")
	if !host.Truthy() || !protocol.Truthy() {
		return errors.New("hostname or protocol is not Truthy"), &js.Value{}
	}
	var proto_ws string
	if protocol.String() == "http:" {
		proto_ws = "ws"
	} else {
		proto_ws = "wss"
	}
	ws := js.Global().Get("WebSocket").New(fmt.Sprintf("%s://%s/%s", proto_ws, host.String(), path))
	return nil, &ws
}

/*
func (o *Xwasmapi) ApiQuery(r *ApiRequest, cb func([]byte, error)) error {
	o.apiQuery(r,cb)
	return nil
}
*/