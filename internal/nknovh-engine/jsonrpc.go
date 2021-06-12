package nknovh_engine

import (
	"encoding/json"
	"io/ioutil"
	"strings"
	)

func (o *NKNOVH) jrpc_get(obj *JsonRPCConf) ([]byte, error) {
	r := RPCRequest{Jsonrpc: "2.0", Method: obj.Method, Id: 1, Params:obj.Params}
	var answer []byte
	client := obj.Client
	if b, err := json.Marshal(r); err != nil {
		return answer, err
	} else {
		if res, err := client.Post("http://" + obj.Ip + ":30003/", "application/json", strings.NewReader(string(b))); err != nil {
			return answer, err
		} else {
			defer res.Body.Close()
			answer, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return answer, err
			}
			//o.log.Syslog("ANSWER RAW: " +string(answer), "nodes")
			return answer, nil
		}
	}
}
