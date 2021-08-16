package nknovh_wasm

type LANG struct {
	LangValues map[string]string `json:"lang_values"`
	CurrentLang           string `json:"current_lang"`
	SenseTimezone         string `json:"sense_timezone"`
	SenseOf               string `json:"sense_of"`
	SenseDays             string `json:"sense_days"`
	SenseHours            string `json:"sense_hours"`
	SenseSeconds          string `json:"sense_seconds"`
	SenseMinutes          string `json:"sense_minutes"`
	SenseMonth            string `json:"sense_month"`
	SenseYears            string `json:"sense_years"`
	SenseEvery            string `json:"sense_every"`
	SenseRelayh           string `json:"sense_relayh"`
	WarnNotMining         string `json:"warn_not_mining"`
	WarnNotMiningMultiple string `json:"warn_not_mining_multiple"`
	Refresher			  map[string]string `json:"refresher"`
	Nav                   map[string]string `json:"nav"`
	Donate		      map[string]string`json:"donate"`
	Attention             map[string]string `json:"attention"`
	Index                 map[string]string`json:"index"`
	WalletTracker         map[string]string`json:"wallet_tracker"`
	NodesTables map[string]map[string]string `json:"nodes_tables"`
	Modal map[string]map[string]string `json:"modal"`
	Outdated map[string]string`json:"outdated"`
	Answers map[string]map[int]string `json:"answers"`
}