package nknovh_engine

import (
		"os"
		"encoding/json"
	)

func (o *configuration) configure() (*configuration, error) {
	f, err := os.Open("conf.json")
	if err != nil {
		return &configuration{}, err	
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return &configuration{}, err
	}
	return &conf,nil
}
