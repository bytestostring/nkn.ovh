package nknovh_engine

type configuration struct {
	Version string
	Db []struct {
		Host string `json:"Host"`
		Login string `json:"Login"`
		Database string `json:"Database"`
		Password string `json:"Password"`
		DbType string `json:"DbType"`
		MaxOpenConns int `json:"MaxOpenConns"`
		MaxIdleConns int `json:"MaxIdleConns"`
		InsideName string `json:"InsideName"`
	} `json:"Db"`
	NeighborPoll struct {
		ConnTimeout int `json:"ConnTimeout"`
		Interval int `json:"Interval"`
		RemoveInterval int `json:"RemoveInterval"`
	} `json:"NeighborPoll"`
	MainPoll struct {
		ConnTimeout int `json:"ConnTimeout"`
		Interval int `json:"Interval"`
		EntriesPerNode int `json:"EntriesPerNode"`
	} `json:"MainPoll"`
	Threads struct {
		Neighbors int `json:"Neighbors"`
		Main int `json:"Main"`
		Dirty int `json:"Dirty"`
	} `json:"Threads"`
	Wallets struct {
		Interval int `json:"Interval"`
	} `json:"Wallets"`
	SeedList []string `json:"SeedList"`
}
