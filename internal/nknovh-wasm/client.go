package nknovh_wasm

import (
	"fmt"
	"syscall/js"
	"xwasmapi"
	"sort"
	"bytes"
	"time"
	"net"
	"errors"
	"encoding/json"
	"text/template"
	"sync"
	"math"
	"regexp"
	"strings"
	"strconv"
	"github.com/fvbommel/sortorder"
)

	//Some html templates
	var wallabel string = `<div style="margin: 20px 0 0 0;"><p>%[1]s %[2]d:</p><input id="setwal-%[2]d" type="text" class="inputtext" value="%[4]s" placeholder="%[3]s"></div>`
	var link_reference string = `<a href="javascript:void(0);" onclick="showModal('reference')">[?]</a>`
	var wallet_div string = `<div class="wallet %[1]s" id="wallet-%[2]d">
	<p style="font-weight: bold">%[3]v %[4]v</p>
	<p><a href="https://explorer.nkn.org/detail/address/%[5]v" rel="noreferrer" target="_blank" title="Explorer">%[5]v</a></p>
	<p>%[6]s</p>
	</div>`
	var node_template string = `<div class="td"><input type="checkbox" id="controlNode-%[1]v" name="controlNode-%[1]v" value="%[1]v"></div><div class="td nodeName">%[2]v</div><div class="td nodeIP">%[3]v</div><div class="td nodeSyncState">%[4]v</div><div class="td nodeProposal">%[5]v</div><div class="td nodeHeight">%[6]v</div><div class="td nodeUptime">%[7]v</div><div class="td nodeRelays">%[8]v</div><div class="td nodeRelays10">%[9]v</div><div class="td nodeRelays60">%[10]v</div><div class="td nodeVersion">%[11]v</div><div class="td nodeUpdated">%[12]v</div>`
	var arrow_asc string = "<span>&#9660;</span>"
	var arrow_desc string = "<span>&#9650;</span>"
	var nodes_sumstat string = `<div class="tr" id="sum_tr"><div class="td"></div><div class="td">%s</div><div class="td">*</div><div class="td">*</div><div class="td">%.f</div><div class="td">*</div><div class="td">%.2f %s</div><div class="td">%s</div><div class="td">%s</div><div class="td">%s</div><div class="td">*</div><div class="td">*</div></div>`

func (c *CLIENT) CheckVersion(actual string) bool {
	doc := js.Global().Get("document")
	if c.Version == "" {
		c.Version = actual
		div := doc.Call("getElementById", "site_version")
		if !div.Truthy() {
			fmt.Println("div#site_version is not Truthy")
			return true
		}
		div.Set("textContent", "Version: " + actual)
		div.Get("style").Set("display", "inline-block")
		return true
	}
	if c.Version != actual {
		html := js.Global().Get("version_src").String()
		data := map[string]interface{}{
			"LANG": c.LANG,
			"CurVersion": c.Version,
			"LastVersion": actual,
		}
		if err, s := c.handlingTemplate(&html, &data); err != nil {
			fmt.Println(err.Error())
			return true
		} else {
			doc.Get("body").Set("innerHTML", s)
			return false
		}
	}

	return true
}

func (c *CLIENT) ToggleCheckBox() {
	doc := js.Global().Get("document")
	ca := doc.Call("getElementById", "control-all")
	if !ca.Truthy() {
		return
	}
	nt := doc.Call("getElementById", "nodes_table")
	nodes_divs := nt.Call("querySelectorAll", "div[id*='Node-']")
	l := nodes_divs.Get("length").Int()
	var b bool = ca.Get("checked").Bool()

	for i := 0; i < l; i++ {
		div := nodes_divs.Index(i)
		if !div.Truthy() {
			continue
		}
		if x := div.Get("style").Get("display").String(); x != "none" {
			check := div.Call("querySelector", "input[type=checkbox]")
			if !check.Truthy() {
				continue
			}
			if b {
				check.Set("checked", true)
			} else {
				check.Set("checked", false)
			}
		}
	}
	return
}


func (c *CLIENT) AddWalletLabels(add_field bool) {
	var n int = 1
	var w string
	var l int

	doc := js.Global().Get("document")
	wloading := doc.Call("getElementById", "settings_wallets_loading")
	if !wloading.Truthy() && add_field == false {
		return
	}
	modal := doc.Call("getElementById", "settingsModal")
	wlist := doc.Call("getElementById", "settingsWalletsList")
	if !modal.Truthy() || !wlist.Truthy() {
		fmt.Println("modal/wlist is not Truthy")
		return
	}

	if add_field != false {
		n = wlist.Call("querySelectorAll", "[id*=setwal-]").Get("length").Int()
		if n >= 3 {
			return
		}
		n++
		s := fmt.Sprintf(wallabel, c.LANG.Modal["settings"]["label_wal"], n, c.LANG.Modal["settings"]["wal_placeholder"], "")
		wlist.Call("insertAdjacentHTML", "beforeend", s)
		if n >= 3 {
			c.W.HideById("set_addmorewallets")
			return
		} else {
			c.W.ShowById("set_addmorewallets")
		}
		return
	}
	if c.Wallets != nil {
		l = len(c.Wallets.Value.Wallets)
	} else {
		l = 0
	}
	if l > 0 {
		for i := 0; i < l; i++ {
			if n == 3 {
				c.W.HideById("set_addmorewallets")
			}
			s := "wal" + strconv.Itoa(n)
			check := doc.Call("getElementById", s)
			if !check.Truthy() {
				s := fmt.Sprintf(wallabel, c.LANG.Modal["settings"]["label_wal"], n, c.LANG.Modal["settings"]["wal_placeholder"], c.Wallets.Value.Wallets[i].NknWallet)
				w += s
				n++
			}
		}
	} else {
		s := "wal" + strconv.Itoa(n)
		check := doc.Call("getElementById", s)
		if !check.Truthy() {
			s := fmt.Sprintf(wallabel, c.LANG.Modal["settings"]["label_wal"], n, c.LANG.Modal["settings"]["wal_placeholder"], "")
			w += s
		} else {
			return
		}
	}
	wlist.Set("innerHTML", w)
	return
}

func (c *CLIENT) SwitchTab(showTab string) {
	doc := js.Global().Get("document")

	colorTab := func(sel string) {
		concat := sel + showTab
		sw := doc.Call("querySelectorAll", "span[id^='"+ sel +"']")
		if !sw.Truthy() {
			return
		}
		l := sw.Get("length").Int()
		for i := 0; i < l; i++ {
			tabElem := sw.Index(i)
			id := tabElem.Get("id")
			if id.String() == concat {
				tabElem.Get("classList").Call("add", "active")
			} else {
				tabElem.Get("classList").Call("remove", "active")
			}
		}
	}

	ShowHideContent := func(m map[string]string) {
		for key, val := range m {
				if key == showTab {
					c.W.ShowById(val)
				} else {
					c.W.HideById(val)
				}
			}
	}

	switch x := showTab; x {
		case "single","multiple":
			m := map[string]string{}
			m["single"] = "addNodesSingle"
			m["multiple"] = "addNodesMultiple"

			colorTab("switch-nodes-")
			ShowHideContent(m)
		break
		case "wallets", "notifications":
			m := map[string]string{}
			m["wallets"] = "settingsWallets"
			m["notifications"] = "settingsNotifications"

			colorTab("switch-settings-")
			ShowHideContent(m)
		break
	}
}

func (c *CLIENT) walletsInfoUpdate() {
	if c.Wallets == nil {
		return
	}
	var n = len(c.Wallets.Value.Wallets)
	if n < 1 {
		c.W.ShowById("wallets_nf")
		return
	}
	var wallets string
	var wclass string
	var usd_val float64
	for i := 0; i < n; i++ {
		wclass = ""
		if c.Wallets.Value.Wallets[i].Balance < 0 {
			wclass = "waiting"
			wallets += fmt.Sprintf(wallet_div, wclass, c.Wallets.Value.Wallets[i].Id, c.LANG.WalletTracker["walletname_label"], i+1, c.Wallets.Value.Wallets[i].NknWallet, fmt.Sprintf("%s: %s", c.LANG.WalletTracker["balance_label"], c.LANG.WalletTracker["wait_for_update"]))
		} else {
			if c.Prices != nil {
				usd_val = c.Wallets.Value.Wallets[i].Balance * c.Prices.Value.Usd
			} else {
				usd_val = 0
			}
			s := fmt.Sprintf("%[1]s: %2.5f NKN [ %3.2f$ ]", c.LANG.WalletTracker["balance_label"], c.Wallets.Value.Wallets[i].Balance, usd_val)
			wallets += fmt.Sprintf(wallet_div, wclass, c.Wallets.Value.Wallets[i].Id, c.LANG.WalletTracker["walletname_label"], i+1, c.Wallets.Value.Wallets[i].NknWallet, s)
		}
	}
	doc := js.Global().Get("document")
	user_wallets := doc.Call("getElementById", "user_wallets")
	if user_wallets.Truthy() {
		user_wallets.Set("innerHTML", wallets)
	}
	return
}

func (c *CLIENT) ParseAll() (resp bool) {
	if c.Nodes == nil || c.Netstatus == nil || c.Daemon == nil {
		return
	}

	resp = true
	if c.Daemon != nil {
		if c.Daemon.Error {
			x := c.Daemon.Code
			c.GenErr(c.Daemon.ErrMessage, "default", x)
		} else {
			if b := c.CheckVersion(c.Daemon.Value.Version); b == false {
				c.mux.AutoUpdater.Lock()
				if c.AutoUpdaterIsStarted {
					c.AutoUpdaterStopCh <- true
				}
				c.mux.AutoUpdater.Unlock()
				return
			}
		}
	}
	if c.Netstatus != nil {
		if c.Netstatus.Error {
			x := c.Netstatus.Code
			c.GenErr(c.Netstatus.ErrMessage, "default", x)
		} else {
			c.parseNetstatus()
		}
	}

	if c.Wallets != nil {
		if c.Wallets.Error {
			x := c.Wallets.Code
			c.GenErr(c.Wallets.ErrMessage, "getmywallets", x)
		} else {
			c.walletsInfoUpdate()
			c.AddWalletLabels(false)
		}
	}

	if c.Nodes != nil {
		if c.Nodes.Error {
			x := c.Nodes.Code
			c.GenErr(c.Nodes.ErrMessage, "default", x)
		} else {
			switch x := c.Nodes.Code; x {
				case 0:
					c.W.HideById("nodes_nf")
					c.PreSortNodes()
					c.SortAndParseNodes()
				break
				case 3:
					c.W.ShowById("nodes_nf")
				break
			}
		}
	}
	c.calcNodesSummary()
	return
}

func (c *CLIENT) ParseFullstack(f *GetFullstack) {
	c.mux.StartView.Lock()
	defer c.mux.StartView.Unlock()
	//Parse main response
	if f.Error {
		x := f.Code
		c.GenErr(f.ErrMessage, "default", x)
		return
	}

	//Parse daemon response
	if f.Value.Daemon.Error {
		x := f.Value.Daemon.Code
		c.GenErr(f.Value.Daemon.ErrMessage, "default", x)
	} else {
		c.Daemon = &f.Value.Daemon
		if b := c.CheckVersion(c.Daemon.Value.Version); b == false {
			c.mux.AutoUpdater.Lock()
			if !c.AutoUpdaterIsStarted {
				c.AutoUpdaterStopCh <- true
			}
			c.mux.AutoUpdater.Unlock()
			return
		}
	}

	//Parse prices
	if f.Value.Prices.Error {
		x := f.Value.Prices.Code
		c.GenErr(f.Value.Prices.ErrMessage, "default", x)
	} else {
		c.Prices = &f.Value.Prices
	}

	//Parse Netstatus
	if f.Value.Netstatus.Error {
		x := f.Value.Netstatus.Code
		c.GenErr(f.Value.Netstatus.ErrMessage, "default", x)
	} else {
		c.Netstatus = &f.Value.Netstatus
		c.parseNetstatus()
	}
	//Parse Wallets
	if f.Value.Wallets.Error {
		x := f.Value.Wallets.Code
		c.GenErr(f.Value.Wallets.ErrMessage, "getmywallets", x)
	} else {
		c.Wallets = &f.Value.Wallets
		c.walletsInfoUpdate()
		c.AddWalletLabels(false)
	}
	//Parse user's nodes
	if f.Value.Nodes.Error {
		x := f.Value.Nodes.Code
		c.GenErr(f.Value.Nodes.ErrMessage, "default", x)
	} else {
		switch x := f.Value.Nodes.Code; x {
			case 0:
				c.W.HideById("nodes_nf")
				c.Nodes = &f.Value.Nodes
				c.PreSortNodes()
				c.SortAndParseNodes()
			break
			case 3:
				c.W.ShowById("nodes_nf")
				c.Nodes = &Nodes{}
			break
		}
	}
	c.calcNodesSummary()
	c.switchLoading(false)
	return	
}

func (c *CLIENT) parseNetstatus() {
	var netstatus *Netstatus
	if c.Netstatus == nil {
		return
	}
	netstatus = c.Netstatus
	var au_sense string
	var au float64
	x := float64(netstatus.Value.AverageUptime)
	if n := x; n <= 3600/60 {
		au = x
		au_sense = c.LANG.SenseSeconds
	} else if n := x/60; n <= 59 {
		au = n
		au_sense = c.LANG.SenseMinutes
	} else if n := x/3600; n <= 24 {
		au = n
		au_sense = c.LANG.SenseHours
	} else if n := x/3600/24; n <= 365 {
		au = n
		au_sense = c.LANG.SenseDays
	} else if n := x/3600/24; n > 365 {
		au = n
		au_sense = c.LANG.SenseYears
	}
	averageUptime := fmt.Sprintf("%.2f", au)
	relaysPerHour := NumSeparate(netstatus.Value.RelaysPerHour)
	averageRelays := NumSeparate(netstatus.Value.AverageRelays)
	nodesCount := NumSeparate(netstatus.Value.NodesCount)
	persistNodesCount := NumSeparate(netstatus.Value.PersistNodesCount)

	lu := netstatus.Value.LatestUpdate + "+03:00"
	//Fix timezone
	time_layout := "2006-01-02 15:04:05Z07:00"
	time_layout2 := "2006-01-02 / 15:04:05" 
	t_now :=  time.Now()
	t_zone, t_offset := t_now.Zone()
	t_loc := time.FixedZone(t_zone, t_offset)
	t, err := time.Parse(time_layout, lu)
	if err != nil {
		fmt.Println(err)
	}
	t = t.In(t_loc)
	latestUpdate := t.Format(time_layout2)

	c.mux.NodesSummary.Lock()
	c.NodesSummary["all"] = map[string]float64{}
	c.NodesSummary["all"]["AverageBlockTime"] = netstatus.Value.AverageBlockTime
	c.NodesSummary["all"]["AverageBlocksPerDay"] = netstatus.Value.AverageBlocksPerDay
	c.NodesSummary["all"]["LastTimestamp"] = float64(netstatus.Value.LastTimestamp)
	c.NodesSummary["all"]["LastHeight"] = float64(netstatus.Value.LastHeight)
	c.NodesSummary["all"]["RelaysPerHour"] = float64(netstatus.Value.RelaysPerHour)
	c.NodesSummary["all"]["NodesCount"] = float64(netstatus.Value.NodesCount)
	c.NodesSummary["all"]["PersistNodesCount"] = float64(netstatus.Value.PersistNodesCount)
	c.mux.NodesSummary.Unlock()

	doc := js.Global().Get("document")
	doc.Call("getElementById", "ns-average_uptime").Set("textContent", averageUptime)
	doc.Call("getElementById", "ns-average_uptime-sense").Set("textContent", au_sense)
	doc.Call("getElementById", "ns-average_relays").Set("textContent", averageRelays)
	doc.Call("getElementById", "ns-average_relays-sense").Set("textContent",c.LANG.SenseRelayh)
	doc.Call("getElementById", "ns-relays_per_hour").Set("textContent", relaysPerHour)
	doc.Call("getElementById", "ns-relays_per_hour-sense").Set("textContent", c.LANG.SenseRelayh)
	doc.Call("getElementById", "ns-latest_update").Set("textContent", latestUpdate)
	doc.Call("getElementById", "ns-latest_update-sense").Set("textContent", c.LANG.SenseTimezone)
	doc.Call("getElementById", "ns-nodes_count").Set("textContent", nodesCount)
	doc.Call("getElementById", "ns-persist_nodes_count").Set("textContent", persistNodesCount)
	c.W.HideById("jNST_loading")
	c.W.ShowById("jNST")
	return
}


func (c *CLIENT) GenErr(errtext string, section string, code int, id ...string) {
	var endtext string
	var elemid string = "error"
	if val, ok := c.LANG.Answers[section][code]; ok {
		endtext = val
	} else {
		endtext = errtext
	}
	if len(id) > 0 {
		elemid = id[0]
	}
	dom := js.Global().Get("document").Call("getElementById", elemid)
		if !dom.Truthy() {
			fmt.Println("No access to div#error")
			return
		}
	dom.Set("innerText", endtext)
	dom.Get("style").Set("display", "block")
	return
}

func (c *CLIENT) PreSortNodes(dom ...*js.Value) {
	doc := js.Global().Get("document")
	if dom != nil {
		id := dom[0].Get("id").String()
		if c.Sort != id {
			c.Sort_type = "ASC"
		}
		if c.Sort_type == "ASC" && c.Sort == id {
			c.Sort_type = "DESC"
		} else {
			c.Sort_type = "ASC"
		}
		c.Sort = id
	}
	tr_top := doc.Call("getElementById", "tr_top").Get("children")
	arr := js.Global().Get("Object").Call("entries", tr_top)
	arr.Call("forEach", js.FuncOf(func(_ js.Value, elems []js.Value) interface{} {
		span := elems[0].Index(1).Get("children").Call("item", 0)
		if span.Truthy() {
			if span.Get("id").String() != "control-all" {
				c.W.Remove(&span)
			}
		}
		return nil
	}))
	parent := doc.Call("getElementById", c.Sort)
	new_span := doc.Call("createElement", "span")
	if c.Sort_type == "ASC" {
		new_span.Set("innerHTML", arrow_asc)
	} else {
		new_span.Set("innerHTML", arrow_desc)
	}
	parent.Call("append", new_span)
	c.W.LocalStorage("set", "sort", c.Sort)
	c.W.LocalStorage("set", "sort_type", c.Sort_type)
	return
}


func (c *CLIENT) calcNodesSummary() {
	c.mux.NodesSummary.Lock()
	defer c.mux.NodesSummary.Unlock()
	var (
		r string
		r10 string
		r60 string
		au float64
		au_sense string
		wop_sense string = c.LANG.SenseSeconds
		mining_sense string
		averageBlocksPerDay float64
		controlPercentage float64
		waitRewardMonth float64
		waitRewardOne float64
		wait_per_month float64
		wait_per_month_usd float64
		inactiveNodes int
	)
	if _, ok := c.NodesSummary["all"]; !ok {
		return
	}
	if _, ok := c.NodesSummary["client"]; !ok {
		return
	}

	//Parse Relays
	if c.NodesSummary["client"]["RelaysPerHour"] == 0 {
		r = "N/A"
	} else {
		r = fmt.Sprintf("%.2fk", c.NodesSummary["client"]["RelaysPerHour"]/1000)
	}
	if c.NodesSummary["client"]["RelaysPerHour10"] == 0 {
		r10 = "N/A"
	} else {
		r10 = fmt.Sprintf("%.2fk", c.NodesSummary["client"]["RelaysPerHour10"]/1000)
	}
	if c.NodesSummary["client"]["RelaysPerHour60"] == 0 {
		r60 = "N/A"
	} else {
		r60 = fmt.Sprintf("%.2fk", c.NodesSummary["client"]["RelaysPerHour60"]/1000)
	}

	//Parse average uptime for client nodes
	x := c.NodesSummary["client"]["AverageUptime"]
	if n := x; n <= 3600/60 {
		au = x
		au_sense = c.LANG.SenseSeconds
	} else if n := x/60; n <= 59 {
		au = n
		au_sense = c.LANG.SenseMinutes
	} else if n := x/3600; n <= 24 {
		au = n
		au_sense = c.LANG.SenseHours
	} else if n := x/3600/24; n <= 365 {
		au = n
		au_sense = c.LANG.SenseDays
	} else if n := x/3600/24; n > 365 {
		au = n
		au_sense = c.LANG.SenseYears
	}

	//Parse rewards and network control
	if x := c.NodesSummary["all"]["AverageBlocksPerDay"]; x == 0 {
		averageBlocksPerDay = 3850
	} else {
		averageBlocksPerDay = x
	}
	controlPercentage = c.NodesSummary["client"]["RelaysPerHour"]/c.NodesSummary["all"]["RelaysPerHour"]*100
	waitRewardMonth = averageBlocksPerDay*30/100*controlPercentage
	inactiveNodes = int(c.NodesSummary["client"]["Nodes"] - c.NodesSummary["client"]["ActiveNodes"])
	if controlPercentage != 0 {
		waitRewardOne = 1440/(averageBlocksPerDay/100*controlPercentage)/60
		if waitRewardOne > 24 {
			waitRewardOne = waitRewardOne/24
			wop_sense = c.LANG.SenseDays
		} else if waitRewardOne < 1 {
			waitRewardOne = waitRewardOne*60
			wop_sense = c.LANG.SenseMinutes
		} else {
			wop_sense = c.LANG.SenseHours
		}
	}
	if x := c.NodesSummary["client"]["ActiveNodes"]; x > 1 {
		mining_sense = c.LANG.NodesTables["Other"]["aremining_label"]
	} else {
		mining_sense = c.LANG.NodesTables["Other"]["ismining_label"]
	}

	doc := js.Global().Get("document")
	wait_per_month = waitRewardMonth*11.09
	if c.Prices != nil {
		wait_per_month_usd = wait_per_month * c.Prices.Value.Usd

		//Add prices to title
		title := doc.Get("title")
		title_slice := strings.Split(title.String(), "|")
		var title_end string
		if len(title_slice) >= 2 {
			title_end = strings.TrimSpace(title_slice[1])
		} else {
			title_end = title.String()
		}
		doc.Set("title", fmt.Sprintf("%.5f $ | %s ", c.Prices.Value.Usd, title_end))

	}

	//Final, printing result

	elSumstat := doc.Call("getElementById", "sum_tr")
	elNodesCount := doc.Call("getElementById", "sum-NodesCount")
	elNetworkControl := doc.Call("getElementById", "sum-NetworkControl")
	elAllRelays := doc.Call("getElementById", "sum-AllRelays")
	elAverageRelays := doc.Call("getElementById", "sum-AverageRelays")
	elWaitProposalMonth := doc.Call("getElementById", "sum-waitProposalMonth")
	elWaitNKNMonth := doc.Call("getElementById", "sum-waitNKNMonth")
	elWaitOneProposal := doc.Call("getElementById", "sum-waitOneProposal")

	if elSumstat.Truthy() {
		c.W.Remove(&elSumstat)
	}
	sumstat := fmt.Sprintf(nodes_sumstat, c.LANG.NodesTables["Other"]["b_sum_label"], c.NodesSummary["client"]["Proposal"], au, au_sense, r, r10, r60)
	div_sumstat := doc.Call("getElementById", "tr_top")
	div_sumstat.Call("insertAdjacentHTML", "afterend", sumstat)

	elAverageRelays.Set("textContent", fmt.Sprintf("%s %s", NumSeparate(int(c.NodesSummary["client"]["AverageRelays"])), c.LANG.SenseRelayh))
	elNodesCount.Set("textContent", fmt.Sprintf("%d %s %d %s", int(c.NodesSummary["client"]["ActiveNodes"]), c.LANG.SenseOf, int(c.NodesSummary["client"]["Nodes"]), mining_sense))
	elNetworkControl.Set("textContent", fmt.Sprintf("%.5f%%", controlPercentage))
	elAllRelays.Set("textContent", fmt.Sprintf("%s %s", NumSeparate(int(c.NodesSummary["client"]["RelaysPerHour"])), c.LANG.SenseRelayh))
	elWaitProposalMonth.Set("textContent", fmt.Sprintf("≈ %.2f / %s", waitRewardMonth, c.LANG.SenseMonth))
	elWaitNKNMonth.Set("textContent", fmt.Sprintf("≈ %.2f NKN / %s [ %.2f$ ]", wait_per_month, c.LANG.SenseMonth, wait_per_month_usd))
	elWaitOneProposal.Set("textContent", fmt.Sprintf("%s ≈ %.2f %s", c.LANG.SenseEvery, waitRewardOne, wop_sense))

	c.W.ShowById("jNST_client")
	if x := inactiveNodes; x >= 1 {
		if x > 1 {
			c.GenErr(fmt.Sprintf("%d %s", x, c.LANG.WarnNotMiningMultiple), "default", -1)
		} else {
			c.GenErr(fmt.Sprintf("%d %s", x, c.LANG.WarnNotMining), "default", -1)
		}
	} else {
		c.W.HideById("error")
	}

	return

}

func (c *CLIENT) SortAndParseNodes() {
	c.mux.Nodes.Lock()
	switch x := c.Sort; x {
		case "t_name":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return sortorder.NaturalLess(c.Nodes.Value.List[i].Name, c.Nodes.Value.List[j].Name)
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return !sortorder.NaturalLess(c.Nodes.Value.List[i].Name, c.Nodes.Value.List[j].Name)
				})
			}
		break
		case "t_status":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].SyncState > c.Nodes.Value.List[j].SyncState
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
				return c.Nodes.Value.List[i].SyncState < c.Nodes.Value.List[j].SyncState 
				})
			}
		break
		case "t_uptime":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].Uptime > c.Nodes.Value.List[j].Uptime 
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
				return c.Nodes.Value.List[i].Uptime < c.Nodes.Value.List[j].Uptime 
				})
			}
		break
		case "t_ip":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return bytes.Compare(net.ParseIP(c.Nodes.Value.List[i].IP), net.ParseIP(c.Nodes.Value.List[j].IP)) > 0
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return bytes.Compare(net.ParseIP(c.Nodes.Value.List[j].IP), net.ParseIP(c.Nodes.Value.List[i].IP)) > 0
				})
			}
		break
		case "t_proposal":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].ProposalSubmitted > c.Nodes.Value.List[j].ProposalSubmitted
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].ProposalSubmitted < c.Nodes.Value.List[j].ProposalSubmitted
				})
			}
		break	
		case "t_height":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].Height > c.Nodes.Value.List[j].Height
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].Height < c.Nodes.Value.List[j].Height
				})
			}
		break	
		case "t_relay":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour > c.Nodes.Value.List[j].RelaysPerHour
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour < c.Nodes.Value.List[j].RelaysPerHour
				})
			}
		break
		case "t_relay10":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour10 > c.Nodes.Value.List[j].RelaysPerHour10
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour10 < c.Nodes.Value.List[j].RelaysPerHour10
				})
			}
		break
		case "t_relay60":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour60 > c.Nodes.Value.List[j].RelaysPerHour60
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].RelaysPerHour60 < c.Nodes.Value.List[j].RelaysPerHour60
				})
			}
		break
		case "t_version":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].Version > c.Nodes.Value.List[j].Version
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					return c.Nodes.Value.List[i].Version < c.Nodes.Value.List[j].Version
				})
			}
		break
		case "t_latestup":
			if c.Sort_type == "ASC" {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					lay := "2006-01-02 15:04:05"
					t1, err := time.Parse(lay, c.Nodes.Value.List[i].LatestUpdate)
					if err != nil {
						fmt.Println(err)
					}
					u1 := t1.Unix()
					t2, err := time.Parse(lay, c.Nodes.Value.List[j].LatestUpdate)
					if err != nil {
						fmt.Println(err)
					}
					u2 := t2.Unix()
					return u1 < u2
				})
			} else {
				sort.SliceStable(c.Nodes.Value.List, func(i, j int) bool {
					lay := "2006-01-02 15:04:05"
					t1, err := time.Parse(lay, c.Nodes.Value.List[i].LatestUpdate)
					if err != nil {
						fmt.Println(err)
					}
					u1 := t1.Unix()
					t2, err := time.Parse(lay, c.Nodes.Value.List[j].LatestUpdate)
					if err != nil {
						fmt.Println(err)
					}
					u2 := t2.Unix()
					return u1 > u2
				})
			}
		break
	}
	var (
		sumUptime int
		sumRelaysPerHour int
		sumRelaysPerHour10 int
		sumRelaysPerHour60 int
		averageRelays int
		averageUptime int
		sumOffline int
		sumProposal int
		sumActiveNodes int
		sumNodes int
		RelaysViewK string
		RelaysViewK10 string
		RelaysViewK60 string
		UptimeView string
		VersionView string
		UpdateView string

		r_uptime float64
		r_relays float64
		r_relays10 float64
		r_relays60 float64
		r_version string
		r_syncstate string
		r_nodeid int
		r_name string
		r_ip string
		r_proposal int
		r_height int
		r_update string
		r_err int
		status string
		node_class string
		waiting_status string
	)

	var nodes []map[string]interface{}
	m, err := json.Marshal(c.Nodes.Value.List)
	c.mux.Nodes.Unlock()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := json.Unmarshal(m, &nodes); err != nil {
		fmt.Println(err.Error())
		return
	}
	var div string

	doc := js.Global().Get("document")
	dom_nodesTable := doc.Call("getElementById", "nodes_table")

	//Fix timezone
	time_layout := "2006-01-02 15:04:05Z0700"
	time_layout2 := "2006-01-02 / 15:04:05" 
	t_now :=  time.Now()
	t_zone, t_offset := t_now.Zone()
	t_loc := time.FixedZone(t_zone, t_offset)

	waiting_status = c.LANG.NodesTables["Other"]["b_node_waiting"]
	var iteration_start time.Time
	var iteration_end time.Time
	var iteration_time time.Duration
	iteration_start = time.Now()

	startShow := c.EntriesPerPage*c.CurrentPage-c.EntriesPerPage+1
	stopShow := startShow+c.EntriesPerPage
	for _, val := range nodes {	
		sumNodes++

		r_nodeid = int(val["NodeId"].(float64))
		r_name = val["Name"].(string)
		r_ip = val["Ip"].(string)
		r_syncstate = val["SyncState"].(string)
		r_update = val["LatestUpdate"].(string)
		r_update += c.Daemon.Value.Timezone

		if x := val["Err"]; x != nil {
			r_err = int(val["Err"].(float64))
		} else {
			r_err = 0
		}

		t, err := time.Parse(time_layout, r_update)
		if err != nil {
			fmt.Println(err)
		}
		t = t.In(t_loc)

		UpdateView = t.Format(time_layout2)

		switch x := r_syncstate; x {
		case "PERSIST_FINISHED":
			node_class = "mining"
			sumActiveNodes++
		break
		case "_OUT_":
			status = "Out of NKN Network"
			status = fmt.Sprintf("%s %s", status, link_reference)
			node_class = "warning_out"
			sumOffline++
		case "PRUNING DB", "GENERATION ID", "WAIT_FOR_SYNCING", "SYNC_STARTED", "SYNC_FINISHED":
			node_class = "warning"
			sumOffline++
		default:
			node_class = "warning"
		}

		switch x := r_err; x {
		case 0,3:
			r_uptime = val["Uptime"].(float64)
			r_relays = val["RelaysPerHour"].(float64)
			r_relays10 = val["RelaysPerHour10"].(float64)
			r_relays60 = val["RelaysPerHour60"].(float64)
			r_version = val["Version"].(string)
			r_height = int(val["Height"].(float64))
			r_proposal = int(val["ProposalSubmitted"].(float64))

			if x == 0 {
				status = r_syncstate
				sumUptime += int(r_uptime)
				if status == "PERSIST_FINISHED" {
					sumRelaysPerHour += int(r_relays)
				}
				sumProposal += r_proposal
			}

			RelaysViewK = fmt.Sprintf("%.2fk", r_relays/1000)
			if r_relays10 > 0 {
				RelaysViewK10 = fmt.Sprintf("%.2fk", r_relays10/1000)
				if x == 0 && status == "PERSIST_FINISHED" {
					sumRelaysPerHour10 += int(r_relays10)
				}
			} else {
				RelaysViewK10 = "N/A"
			}
			if r_relays60 > 0 {
				RelaysViewK60 = fmt.Sprintf("%.2fk", r_relays60/1000)
				if x == 0 && status == "PERSIST_FINISHED" {
					sumRelaysPerHour60 += int(r_relays60)
				}
			} else {
				RelaysViewK60 = "N/A"
			}
			if r_version == "" {
				VersionView = "N/A"
			} else {
				VersionView = r_version
			}
			if r_uptime >= 3600*24 {
				UptimeView = fmt.Sprintf("%.2fd", r_uptime/(3600*24))
			} else if r_uptime >= 3600 {
				UptimeView = fmt.Sprintf("%.2fh", r_uptime/(3600))
			} else {
				UptimeView = fmt.Sprintf("%ds", int(r_uptime))
			}
		break
		case 1:
			status = "OFFLINE"
			node_class = "error"
			sumOffline++
		break
		case 2:
			status = waiting_status
			node_class = "waiting"
			sumOffline++
		break
		}
		div_id := fmt.Sprintf("Node-%v", r_nodeid)

		nodediv := js.Global().Get("document").Call("getElementById", div_id)
		if nodediv.Truthy() {
			classList := nodediv.Get("classList")
			if x := classList.Call("contains", node_class); !x.Truthy() {
				classList.Call("remove", classList.Index(1).String())
				classList.Call("add", node_class)
			}
			divs := nodediv.Call("getElementsByTagName", "div")

			divs.Index(1).Set("textContent", r_name)
			divs.Index(2).Set("textContent", r_ip)
			divs.Index(3).Set("innerHTML", status)
			divs.Index(11).Set("textContent", UpdateView)

			if r_err == 0 || r_err == 3 {
				divs.Index(4).Set("textContent", r_proposal)
				divs.Index(5).Set("textContent", r_height)
				divs.Index(6).Set("textContent", UptimeView)
				divs.Index(7).Set("textContent", RelaysViewK)
				divs.Index(8).Set("textContent", RelaysViewK10)
				divs.Index(9).Set("textContent", RelaysViewK60)
				divs.Index(10).Set("textContent", VersionView)
			} else {
				divs.Index(4).Set("textContent", "N/A")
				divs.Index(5).Set("textContent", "N/A")
				divs.Index(6).Set("textContent", "N/A")
				divs.Index(7).Set("textContent", "N/A")
				divs.Index(8).Set("textContent", "N/A")
				divs.Index(9).Set("textContent", "N/A")
				divs.Index(10).Set("textContent", "N/A")
			}
			if err, div_save := c.W.Detach(div_id); err == nil {
				dom_nodesTable.Call("append", div_save)
				if !(sumNodes >= startShow && sumNodes < stopShow) {
					div_save.Get("style").Set("display", "none")
				} else {
					div_save.Get("style").Set("display", "table-row")
				}
			} else {
				fmt.Println("Cannot detach by ID: ", div_id)
			}
			continue
		}
		if r_err == 0 || r_err == 3 {
			div = fmt.Sprintf(node_template, r_nodeid, r_name, r_ip, status, r_proposal, r_height, UptimeView, RelaysViewK, RelaysViewK10, RelaysViewK60, VersionView, UpdateView)
		} else {
			div = fmt.Sprintf(node_template, r_nodeid, r_name, r_ip, status, "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", "N/A", UpdateView)
		}
		dom_div := doc.Call("createElement", "div")
		dom_div.Set("id", div_id) 
		dom_div.Set("className", "tr " + node_class)
		dom_div.Set("innerHTML", div)
		if !(sumNodes >= startShow && sumNodes < stopShow) {
			dom_div.Get("style").Set("display", "none")
		}
		dom_nodesTable.Call("append", dom_div)
	}
	c.ReloadTableSettings()
	//cut that
	iteration_end = time.Now()
	iteration_time = iteration_end.Sub(iteration_start)
	fmt.Println(iteration_time)

	if sumActiveNodes > 0 {
		activeNodes := sumNodes - sumOffline
		if (sumRelaysPerHour > 0) && activeNodes > 0 {
			averageRelays = sumRelaysPerHour/activeNodes
			averageUptime = sumUptime/activeNodes
		}
	}
	c.mux.NodesSummary.Lock()
	c.NodesSummary["client"] = map[string]float64{}
	c.NodesSummary["client"]["Nodes"] = float64(sumNodes)
	c.NodesSummary["client"]["ActiveNodes"] = float64(sumActiveNodes)
	c.NodesSummary["client"]["RelaysPerHour10"] = float64(sumRelaysPerHour10)
	c.NodesSummary["client"]["RelaysPerHour"] = float64(sumRelaysPerHour)
	c.NodesSummary["client"]["RelaysPerHour10"] = float64(sumRelaysPerHour10)
	c.NodesSummary["client"]["RelaysPerHour60"] = float64(sumRelaysPerHour60)
	c.NodesSummary["client"]["AverageRelays"] = float64(averageRelays)
	c.NodesSummary["client"]["AverageUptime"] = float64(averageUptime)
	c.NodesSummary["client"]["Proposal"] = float64(sumProposal)
	c.mux.NodesSummary.Unlock()
	return
}

func (c *CLIENT) SetEntriesPerPage(num int) {
	c.EntriesPerPage = num
	c.W.LocalStorage("set", "entriesPerPage", strconv.Itoa(num))
	var nodes_count int
	c.mux.Nodes.Lock()
	if c.Nodes != nil {
		nodes_count = len(c.Nodes.Value.List)
	} else {
		nodes_count = 0
	}
	c.mux.Nodes.Unlock()
	if x := float64(nodes_count)/float64(c.EntriesPerPage); x < float64(c.CurrentPage) && x != 0 {
		c.CurrentPage = int(math.Ceil(x))
	}
	c.SortTableLite()
	c.ReloadTableSettings()
	return
}

func (c *CLIENT) SetPage(num int) {
	c.CurrentPage = num
	c.SortTableLite()
	c.ReloadTableSettings()
	return
}

func (c *CLIENT) SortTableLite() {
	doc := js.Global().Get("document")
	dom_nt := doc.Call("getElementById", "nodes_table")
	startShow := c.EntriesPerPage*c.CurrentPage-c.EntriesPerPage+1
	stopShow := startShow+c.EntriesPerPage
	childs := dom_nt.Get("children")
	l := childs.Get("length").Int()
	var nodes_count int
	for i := 0; i < l; i++ {
		div := childs.Index(i)
		id := div.Get("id").String()
		if !strings.HasPrefix(id, "Node-") {
			continue
		}
		nodes_count++
		if !(nodes_count >= startShow && nodes_count < stopShow) {
			div.Get("style").Set("display", "none")
		} else {
			div.Get("style").Set("display", "table-row")
		}
	}
	return
}

func (c *CLIENT) AddNodes() {
	doc := js.Global().Get("document")
	button := doc.Call("getElementById", "addNodeButton")
	button.Set("disabled", true)
	single := doc.Call("getElementById", "addNodesSingle")
	var multiple bool
	data := map[string]interface{}{}
	if x := single.Get("style").Get("display").String(); x == "none" {
		multiple = true
	} else {
		multiple = false
	}

	if !multiple {
		data["Ip"] = doc.Call("getElementById", "nodeIP").Get("value").String()
		data["Name"] = doc.Call("getElementById", "nodeName").Get("value").String()
		data["Multiple"] = false
		c.WsSend("addnodes", data)
		return
	}
	data["Ip"] = doc.Call("getElementById", "nodeIPList").Get("value").String()
	data["Name"] = doc.Call("getElementById", "nodeNameList").Get("value").String()
	data["Multiple"] = true
	c.WsSend("addnodes", data)
	return
}

func (c *CLIENT) RmNodes() {
	doc := js.Global().Get("document")
	nt := doc.Call("getElementById", "nodes_table")
	checkboxes := nt.Call("querySelectorAll", "input[type=checkbox]")
	l := checkboxes.Get("length").Int()
	nodes := make([]int, 0)
	for i := 0; i < l; i++ {
		if checkboxes.Index(i).Get("checked").Bool() == false {
			continue
		}
		id := checkboxes.Index(i).Get("id").String()
		s := strings.Split(id, "-")
		if node_id, err := strconv.Atoi(s[1]); err == nil {
			nodes = append(nodes, node_id)
		}
	}

	m := map[string]interface{}{}
	m["NodesId"] = nodes
	c.WsSend("rmnodes", m)
	return
}

func (c *CLIENT) ReloadTableSettings() {
	doc := js.Global().Get("document")
	sel := doc.Call("getElementById", "selEntriesPerPage")
	if sel.Truthy() {
		sel.Set("value", c.EntriesPerPage)
	}
	var nodes_count int
	var pages float64
	var toback string = `<div class="node_table_page back" onclick="prevPage()">&laquo;</div>`
	var tonext string = `<div class="node_table_page next" onclick="nextPage()">&raquo;</div>`
	var topage string = `<div class="node_table_page %[2]v" onclick="setPage(%[1]v);">%[1]v</div>`
	c.mux.Nodes.Lock()
	if c.Nodes != nil {
		nodes_count = len(c.Nodes.Value.List)
	} else {
		nodes_count = 0
	}
	c.mux.Nodes.Unlock()
	if x := float64(nodes_count)/float64(c.EntriesPerPage); x > 1 {	
		pages = math.Ceil(x)
		if pages < float64(c.CurrentPage) {
			c.CurrentPage = int(pages)
			c.SortTableLite()
		}
	} else {
		pages = 1
	}
	div_pages := doc.Call("getElementById", "nodes_pages")
	if !div_pages.Truthy() {
		fmt.Println("div#nodes_pages not found")
		return
	}
	var toback_comp bool
	var final_div string
	var addclass string
	num_pages := int(pages)
	for i := 1; i <= num_pages; i++ {
		addclass = ""
		if toback_comp == false && c.CurrentPage > i {
			final_div += toback
			toback_comp = true
		}
		if c.CurrentPage == i {
			addclass = "current"
		}
		final_div += fmt.Sprintf(topage, i, addclass) + " "
		if x := i+1; x > num_pages && c.CurrentPage < i {
			final_div += tonext
		}
	}
	div_pages.Set("innerHTML", final_div)
	return
}

func (c *CLIENT) ShowHideModal(id string, event string) {
	doc := js.Global().Get("document")
	id_mod := id + "Modal"
	modal := doc.Call("getElementById", id_mod)
	if !modal.Truthy() {
		fmt.Println("Modal is undefined")
		return
	}
	if _, ok := c.Objects.Listeners["Modal"]; !ok {
		listener := js.FuncOf(func(_ js.Value, event []js.Value) interface{} {
			target_id := event[0].Get("target").Get("id")
			if target_id.String() == id_mod {
				c.ShowHideModal(id, "hide")
			}
			return nil
		})
		c.Objects.Listeners["Modal"] = &listener
	}
	if event == "show" {
		doc.Get("body").Call("addEventListener", "click", c.Objects.Listeners["Modal"])
		modal.Get("style").Set("display", "block")
	}
	if event == "hide" {
		doc.Get("body").Call("removeEventListener", "click", c.Objects.Listeners["Modal"])
		delete(c.Objects.Listeners, "Modal")
		modal.Get("style").Set("display", "none")		
	}
	return
}

func (c *CLIENT) SaveSettings() {
	doc := js.Global().Get("document")
	data := map[string]interface{}{}

	//Handling wallets
	wallets := make([]string, 0, 3)
	wlist_div := doc.Call("getElementById", "settingsWalletsList")
	wlist_input := wlist_div.Call("querySelectorAll", "[id*=setwal-]")
	lenw := wlist_input.Get("length").Int()
	for i := 0; i < lenw; i++ {
		v := wlist_input.Index(i)
		value := strings.TrimSpace(v.Get("value").String())
		if len(value) == 0 {
			continue
		}
		wallets = append(wallets, value)
	}
	data["Wallets"] = wallets
	c.WsSend("savemysettings", data)
	return
}

func (c *CLIENT) HideAttention() {
	c.Hide_attention = true
	c.W.LocalStorage("set", "hide_attention", "yes")
	c.W.HideById("attention")
	return
}

func (c *CLIENT) handlingIndex() error {
	doc := js.Global().Get("document")
	wrapper := doc.Call("getElementById", "wrapper")
	wrapper.Set("innerHTML", "")
	index := js.Global().Get("index_src")
	if !index.Truthy() {
		fmt.Println("No index_src")
		return errors.New("no index_src")
	}
	old_content_div := doc.Call("getElementById", "content")
	if t := old_content_div.Type(); t == js.TypeObject {
		c.W.Remove(&old_content_div)
	}
	content := doc.Call("createElement", "div")
	content.Set("id", "content")
	content.Set("className", "content index")
	wrapper.Call("append", content)
	c.switchLoading(false)
	return nil
}

func (c *CLIENT) handlingViewStart() error {
	err, doc := c.W.Get("document")
	if err != nil {
		fmt.Println("Cannot fetch the document DOM Element")
		return err
	}
	old_content_div := doc.Call("getElementById", "content")
	if t := old_content_div.Type(); t == js.TypeObject {
		c.W.Remove(&old_content_div)
	}
	wrapper := doc.Call("getElementById", "wrapper")
	content := doc.Call("createElement", "div")
	content.Set("id", "content")
	content.Set("className", "content in")
	wrapper.Call("append", content)
	sl := doc.Call("getElementById", "switchLang")
	if !sl.Truthy() {
		_, sw := c.W.Get("view_switchLang")
		_, logo := c.W.Get("view_logo")
		wrapper.Call("insertAdjacentHTML", "afterbegin", logo.String())
		wrapper.Call("insertAdjacentHTML", "afterbegin", sw.String())
	}
	return nil
}

func (c *CLIENT) handlingViewEnd() error {
	c.ReloadTableSettings()
	c.switchLangActive()
	c.PreSortNodes()
	if !c.Hide_attention {
		at := js.Global().Get("attention")
		if !at.Truthy() {
			s := "Attention string is not Truthy"
			fmt.Println(s)
			return errors.New(s)
		}
		doc := js.Global().Get("document")
		at_div := doc.Call("getElementById", "attention")
		if at_div.Truthy() {		
			data := map[string]interface{}{
											"LANG": c.LANG,
											"Hash": c.Hash,
			}
			html := at.String()
			if err, s := c.handlingTemplate(&html, &data); err == nil {
				at_div.Set("innerHTML", s)
				c.W.ShowById("attention")
			} else {	
				fmt.Println(err.Error())
				return err
			}
		}
	}
	if b := c.ParseAll(); b {
		c.switchLoading(false)
	}
	return nil
}

func (c *CLIENT) handlingLangPages(view, locale string)  error {
	if view == "view_src" {
		if err := c.handlingViewStart(); err != nil {
			fmt.Println(err.Error())
			return err
		}
	} else if view == "index_src" {
		if err := c.handlingIndex(); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	dom := js.Global().Get(view)
	if  !dom.Truthy() {
		return errors.New("Element is not Truthy")
	}
	html := dom.String()
	data := map[string]interface{}{
		"LANG": c.LANG,
	}
	if err, s := c.handlingTemplate(&html, &data); err != nil {
		fmt.Println(err.Error())
		return err
	} else {
		doc := js.Global().Get("document")
		content := doc.Call("getElementById", "content")
		content.Set("innerHTML", s)
	}
	if view == "view_src" {
		if err := c.handlingViewEnd(); err != nil {
			fmt.Println(err.Error())
			return err
		}
	}
	return nil
}


func (c *CLIENT) handlingTemplate(html *string, data *map[string]interface{}) (error, string) {
	var buf bytes.Buffer

	t, err := template.New("").Parse(*html)
	if err != nil {
		return err, ""
	}
	if err := t.Execute(&buf, *data); err != nil {
		return err, ""
	}
	return nil, buf.String()
}

func (c *CLIENT) switchLoading(b bool) {
	doc := js.Global().Get("document")
	wrapper := doc.Call("getElementById", "wrapper")
	bl := doc.Call("getElementById", "baseloading")
	if b {
		c.W.HideById("content")
		if !bl.Truthy() {
			bl_object := js.Global().Get("baseloading")
			wrapper.Call("insertAdjacentHTML", "beforeend", bl_object.String())
		} else {
			bl.Get("classList").Call("remove", "hidden")
			c.W.ShowById("baseloading")
		}
		return
	}
	bl.Get("classList").Call("add", "hidden")
	return
}

func (c *CLIENT) switchLangActive() error {
	doc := js.Global().Get("document")
	sLang := doc.Call("getElementById", "switchLang")
	if !sLang.Truthy() {
		s := "switchLang div is not Truthy"
		fmt.Println(s)
		return errors.New(s)			
	}
	n := sLang.Get("childElementCount").Int()
	for i := 0; i < n; i++ {
		id := sLang.Get("children").Index(i).Get("id").String()
		if id == "" {
			continue
		}
		l := doc.Call("getElementById", id)
		if l.Get("classList").Call("contains", "active").Truthy() {
			l.Get("classList").Call("remove", "active")
		}
		if id == "lang_" + c.Lang {
			l.Get("classList").Call("add", "active")
		}
	}
	return nil
}

func (c *CLIENT) Init() {
	c.Conf = new(Conf)
	c.W = xwasmapi.New()
	c.Conf.DefaultLanguage = "en_US"
	c.Conf.DefaultEntriesPerPage = 50
	c.CurrentPage = 1
	c.Cached = &Cached{Pages: map[string]string{}, Lang: map[string]*LANG{},}
	c.Debug = true
	c.Objects = &Objects{Listeners:map[string]*js.Func{},}
	c.mux = &Mutexes{AutoUpdater: &sync.RWMutex{}, Nodes: &sync.Mutex{}, NodesSummary: &sync.Mutex{}, StartView: &sync.Mutex{}, Websocket: &sync.Mutex{},}
	c.NodesSummary = map[string]map[string]float64{}
	c.AutoUpdaterStartCh = make(chan bool)
	c.AutoUpdaterStopCh = make(chan bool)

	c.apiMethods = map[string]func(*WSReply) interface{}{}
	c.apiMethods["auth"] = c.apiAuth
	c.apiMethods["genid"] = c.apiGenId
	c.apiMethods["getfullstack"] = c.apiFullstack
	c.apiMethods["getmynodes"] = c.apiMyNodes
	c.apiMethods["getmywallets"] = c.apiMyWallets
	c.apiMethods["getnetstatus"] = c.apiNetstatus
	c.apiMethods["getprices"] = c.apiPrices
	c.apiMethods["getdaemon"] = c.apiDaemon
	c.apiMethods["addnodes"] = c.apiAddNodes
	c.apiMethods["rmnodes"] = c.apiRmNodes
	c.apiMethods["getlanguage"] = c.apiLanguage
	c.apiMethods["savemysettings"] = c.apiSaveSettings
	c.RegisterEvents()
	go c.AutoUpdater()
}

func (c *CLIENT) GetHashFromPath() (error, string) {
	re_login := regexp.MustCompile(`^/login/auth/([A-Za-z0-9]{64})(/?)$`)
	//DEPRECATED
	re_id := regexp.MustCompile(`^/id/([A-Za-z0-9]{64})(/?)$`)
	err, path := c.W.Get("document", "location", "pathname")
	if err != nil {
		return errors.New("Cannot get pathname"), ""
	}
	path_s := path.String()
	if s := re_login.FindStringSubmatch(path_s); s != nil {
		return nil, s[1]
	}
	if s := re_id.FindStringSubmatch(path_s); s != nil {
		return nil, s[1]
	}
	return nil, ""
}

func (c *CLIENT) RegisterEvents() {
	//Loading event
	baseloading_cb := js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
		c.W.HideById("baseloading")
		return nil
	})

	baseloading := js.Global().Get("document").Call("getElementById", "baseloading")
	if baseloading.Truthy() {
		baseloading.Call("addEventListener", "webkitAnimationEnd", baseloading_cb)
		baseloading.Call("addEventListener", "animationend", baseloading_cb)
		baseloading.Call("addEventListener", "oanimationend", baseloading_cb)
	}
	fmt.Println("Events registered")
}

func (c *CLIENT) Run() {

	loaded := make(chan struct{})
	_, doc := c.W.Get("document")
	switch l := doc.Get("readyState").String(); l {
		case "loading":
			doc.Call("addEventListener", "DOMContentLoaded", js.FuncOf(func(_ js.Value, _ []js.Value) interface{} {
				close(loaded)
				return nil
			}))
		case "complete", "interactive":
			close(loaded)
		default:
			fmt.Println("Unexpected document.ReadyState: ", l)
			return
	}

	<-loaded

	var (
		hash string
		sort1 string
		sort_type string
		err error
	)

	if err, entries := c.W.LocalStorage("get", "entriesPerPage"); err == nil {
		i, err := strconv.Atoi(entries)
		if err == nil {
			c.EntriesPerPage = i
		} else {
			c.EntriesPerPage = c.Conf.DefaultEntriesPerPage
		}
	} else {
		c.EntriesPerPage = c.Conf.DefaultEntriesPerPage
	}

	err, hash = c.GetHashFromPath()
	if hash == "" {
		if err, hash = c.W.LocalStorage("get", "hash"); err == nil {
			c.Hash = hash
		} else {
			c.Hash = ""
		}
	} else {
		c.Hash = hash
		c.W.LocalStorage("set", "hash", hash)
	}
	if err, _ = c.W.LocalStorage("get", "hide_attention"); err == nil {
		c.Hide_attention = true
	} else {
		c.Hide_attention = false
	}

	err1, sort1 := c.W.LocalStorage("get", "sort")
	err2, sort_type := c.W.LocalStorage("get", "sort_type")
	if err1 == nil && err2 == nil {
		c.Sort = sort1
		c.Sort_type = sort_type
	} else {
		c.Sort = "t_name"
		c.Sort_type = "ASC"
	}

	//Create WS connection
	c.mux.Websocket.Lock()
	if c.ws == nil {
		err, ws := c.W.WsCreate("polling")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c.ws = ws
		c.WsOnOpen()
		c.WsOnClose()
		c.WsOnMessage()
	} else {
		//onopen
		if c.Hash != "" {
			c.SetLanguage("view_src", "")
			c.mux.AutoUpdater.Lock()
			if b := c.AutoUpdaterIsStarted; !b {
				c.AutoUpdaterStartCh <- true
			}
			c.mux.AutoUpdater.Unlock()
		} else {
			c.SetLanguage("index_src", "")
		}
	}
	c.mux.Websocket.Unlock()

}

func (c *CLIENT) RefreshStatus(seconds int) {
	doc := js.Global().Get("document")

	el := doc.Call("getElementById", "refreshStatus")
	if !el.Truthy() {
		fmt.Println("refreshStatus is not Truthy")
		return
	}
	el_seconds := doc.Call("getElementById", "refreshSeconds")
	cl := el.Get("classList")
	cl_length := cl.Get("length").Int()
	if cl_length == 1 && seconds == -100 {
		el_text := doc.Call("getElementById", "refreshText")
		el_sense := doc.Call("getElementById", "refreshSense")
		el_text.Set("textContent", fmt.Sprintf("%s", c.LANG.Refresher["stop"]))
		el_sense.Set("textContent", "")
		el_seconds.Set("textContent", "")
		cl.Call("add", "stopped")
		return
	} else if cl_length > 1 && seconds == -100 {
		return
	} else if cl_length > 1 && seconds != -100 {
		cl.Call("remove", "stopped")
		el_text := doc.Call("getElementById", "refreshText")
		el_sense := doc.Call("getElementById", "refreshSense")
		el_text.Set("textContent", fmt.Sprintf("%s", c.LANG.Refresher["ok"]))
		el_sense.Set("textContent", c.LANG.SenseSeconds)
	}
	data := fmt.Sprintf("%v", seconds)
	el_seconds.Set("textContent", data)
	return
}

func (c *CLIENT) AutoUpdater() {
	var timetoupdate int = 60
	reset := 1 * time.Millisecond
	normal := time.Duration(timetoupdate) * time.Second
	ticker := time.NewTicker(reset)
	countdown := time.NewTicker(1000 * time.Millisecond)
	cd_mux := &sync.Mutex{}
	
	ch_writer := make(chan int)
	pagewriter := func() {
		for {
			select {
				case x := <-ch_writer:
					c.RefreshStatus(x)
					if x == -100 {
						return
					}
			}
		}
	}
	
	defer ticker.Stop()
	defer countdown.Stop()
	for {
		<-c.AutoUpdaterStartCh
		c.mux.AutoUpdater.Lock()
		c.AutoUpdaterIsStarted = true
		c.mux.AutoUpdater.Unlock()
		fmt.Println("AutoUpdater has started")
		go func() {
			go pagewriter()
			for {
				select {
					case <-c.AutoUpdaterStopCh:
						ticker.Reset(reset)
						ch_writer <- -100
						countdown.Stop()
						c.mux.AutoUpdater.Lock()
						c.AutoUpdaterIsStarted = false
						c.mux.AutoUpdater.Unlock()
						fmt.Println("AutoUpdater has stopped")
						return
					case <-ticker.C:
						ticker.Reset(normal)
						countdown.Reset(1000 * time.Millisecond)
						cd_mux.Lock()
						timetoupdate = 61
						cd_mux.Unlock()
						c.WsGetFullstack()
					break
					case <-countdown.C:
						cd_mux.Lock()
						timetoupdate = timetoupdate - 1
						ch_writer <- timetoupdate
						cd_mux.Unlock()
					break
				}
			}
		}()
	}
	return
}
