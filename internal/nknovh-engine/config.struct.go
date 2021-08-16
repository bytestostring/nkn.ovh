package nknovh_engine

type configuration struct {
	Version string
	HttpServer struct {
		Port int `json:"Port"`
	} `json:"HttpServer"`
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
	DirtyPoll struct {
		ConnTimeout int `json:"ConnTimeout"`
		Interval int `json:"Interval"`
	} `json:"DirtyPoll"`
	Threads struct {
		Neighbors int `json:"Neighbors"`
		Main int `json:"Main"`
		Dirty int `json:"Dirty"`
	} `json:"Threads"`
	Wallets struct {
		Interval int `json:"Interval"`
	} `json:"Wallets"`
	Messengers struct {
		Telegram struct {
			Use bool `json:"Use"`
			Token string `json:"Token"`
		}
	} `json:"Messengers"`
	TrustedProxies []string `json:"TrustedProxies"`
	SeedList []string `json:"SeedList"`
}
