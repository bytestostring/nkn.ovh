package nknovh_engine

import (
		"database/sql"
		_ "github.com/go-sql-driver/mysql"
		"errors"
		"math"
		"regexp"
		"net"
		"strings"
		"strconv"
		"io/ioutil"
		"fmt"
		"time"
		"encoding/json"
)

func (o *NKNOVH) updateUniqWatch(c *CLIENT) error {
	if c.HashId == -1 {
		return nil
	}
	_, err := o.sql.stmt["main"]["WebUpdateUniqWatch"].Exec(c.HashId)
	if err != nil {
		o.log.Syslog("Cannot execute WebUpdateUniqWatch", "sql")
		return err
	}
	return nil
}

func (o *NKNOVH) WsError(q *WSQuery, code int) (err error, r WSReply) {
	var ok bool
	if r, ok = o.Web.Response[code]; ok {
		r.Method = q.Method
		err = errors.New(r.ErrMessage)
		return
	}
	err = errors.New("Response key is not found")
	return nil, WSReply{Method: q.Method, Code: -1, Error: true, ErrMessage: "Response key is not found"}
}

func (o *NKNOVH) apiGetNodeDetails(q *WSQuery, c *CLIENT) (err error, r WSReply) {

	t0 := time.Now()
	var node_id int
	if raw_node_id, ok := q.Value["NodeId"].(float64); !ok {
		if raw_node_id_s, ok := q.Value["NodeId"].(string); !ok {
			return o.WsError(q, 19)
		} else {
			x, err := strconv.Atoi(raw_node_id_s)
			if err != nil {
				return o.WsError(q, 19)
			}
			node_id = x
		}
	} else {
		node_id = int(raw_node_id)
	}

	var node_ip string
	var node_name string
	row := o.sql.stmt["main"]["WebSelectNodeInfoById+HashId"].QueryRow(node_id, c.HashId)
	err = row.Scan(&node_name, &node_ip)
	switch {
		case err == sql.ErrNoRows:
			return o.WsError(q, 18)
		case err != nil:
			return o.WsError(q, 1)
	}
	var data NodeSt
	state := &JsonRPCConf{Ip:node_ip, Method:"getnodestate", Params: &json.RawMessage{'{','}'}, Client: o.http.MainClient, UnmarshalData: &data.State}

	t1 := time.Now()
	res, err := o.jrpc_get(state)
	t1_time := time.Now().Sub(t1)
	if err != nil && len(res) == 0 {
		return o.WsError(q, 20)
	}
	if err != nil && len(res) > 0 {
		return o.WsError(q, 21)
	}
	if b := o.isNodeStateValid(&data.State); !b {
		return o.WsError(q, 25)
	}
	m := map[string]interface{}{}

	if data.State.Error != nil {
		r := o.respErrorHandling(data.State.Error)
		m["Data"] = r
		return nil, WSReply{Method: q.Method, Code: 3, Value: m,}
	}

	neighbor := &JsonRPCConf{Ip:node_ip, Method:"getneighbor", Params: &json.RawMessage{'{','}'}, Client: o.http.MainClient, UnmarshalData: &data.Neighbor}

	t2 := time.Now()
	res, err = o.jrpc_get(neighbor)
	t2_time := time.Now().Sub(t2)

	if err != nil && len(res) == 0 {
		return o.WsError(q, 22)
	}
	if err != nil && len(res) > 0 {
		return o.WsError(q, 23)
	}

	if data.Neighbor.Error != nil {
		return o.WsError(q, 24)
	}

	type NodeStats struct {
		MinPing int
		AvgPing int
		MaxPing int
		NeighborCount int
		NeighborPersist int
		RelaysPerHour uint64
		NodeState *NodeState
	}

	ns := new(NodeStats)
	if data.State.Result.Uptime > 0 {
		ns.RelaysPerHour = uint64(math.Floor(float64(data.State.Result.RelayMessageCount)/float64(data.State.Result.Uptime)*3600))
	} else {
		ns.RelaysPerHour = 0
	}
	ns.NodeState = &data.State

	//Get the neighbors info
	ncount := len(data.Neighbor.Result)
	if ncount != 0 {
		var min int = -1
		var max int = -1
		var sumping int
		var sumpersist int
		for i := 0; i < ncount; i++ {
			if min == -1 || max == -1 {
				min = data.Neighbor.Result[i].RoundTripTime
				max = data.Neighbor.Result[i].RoundTripTime
			}
			if data.Neighbor.Result[i].RoundTripTime > max {
				max = data.Neighbor.Result[i].RoundTripTime
			}
			if data.Neighbor.Result[i].RoundTripTime < min {
				min = data.Neighbor.Result[i].RoundTripTime
			}

			sumping += data.Neighbor.Result[i].RoundTripTime
			if data.Neighbor.Result[i].SyncState == "PERSIST_FINISHED" {
				sumpersist++
			}
		}
		ns.AvgPing = int(math.Round(float64(sumping)/float64(ncount)))
		ns.MaxPing = max
		ns.MinPing = min
		ns.NeighborCount = ncount
		ns.NeighborPersist = sumpersist
	}

	m["NodeStats"] = ns
	t0_time := time.Now().Sub(t0)
	m["DebugInfo"] = map[string]interface{}{
			"GetnodestateTime": t1_time.String(),
			"GetneighborTime": t2_time.String(),
			"HandlingTime": t0_time.String(),
	}

	return nil, WSReply{Method: q.Method, Code: 0, Value: m,}
}

func (o *NKNOVH) apiSaveSettings(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	var ok bool
	wallets_limit := 3
	var wallets_passed []string
	var raw_wallets []interface{}

	re_wallets := regexp.MustCompile(`^NKN([A-Za-z0-9]{33})$`)
	tx, err := o.sql.db["main"].Begin()
	if err != nil {
		return o.WsError(q, 1)
	}

	defer tx.Rollback()

	addWallet := func(val string) error {
		_, err := tx.Stmt(o.sql.stmt["main"]["WebAddWallet"]).Exec(c.HashId, val)
		if err != nil {
			return err
		}
		return nil
	}
	rmWallet := func(id int) error {
		_, err := tx.Stmt(o.sql.stmt["main"]["WebRmWalletById"]).Exec(id)
		if err != nil {
			return err
		}
		return nil
	}
	//Wallets handling
	if raw_wallets, ok = q.Value["Wallets"].([]interface{}); ok {
		if len(raw_wallets) > wallets_limit {
			return o.WsError(q, 17)
		}
		for i, _ := range raw_wallets {
			val, ok := raw_wallets[i].(string)
			if !ok {
				return o.WsError(q, 17)
			}
			if ok = re_wallets.MatchString(val); !ok {
				return o.WsError(q, 17)
			}
			wallets_passed = append(wallets_passed, val)
		}
		if len(wallets_passed) < 1 {
			//Remove all wallets
			_, err := tx.Stmt(o.sql.stmt["main"]["WebRmAllWalletsByHash"]).Exec(c.HashId)
			if err != nil {
				o.log.Syslog("Cannot execute Tx stmt query: " + err.Error(), "sql")
				return o.WsError(q, 1)
			}
		} else {
			db_wallets := make([]string, 0, 3)
			db_wallets_id := make([]int, 0, 3)
			rows, err := o.sql.stmt["main"]["WebGetMyWallets"].Query(c.HashId)
			if err != nil {
				return o.WsError(q, 1)
			}
			defer rows.Close()
			var (
				id int
				nkn_wallet string
				balance float64
			)
			for rows.Next() {
				if err = rows.Scan(&id, &nkn_wallet, &balance); err != nil {
					return o.WsError(q, 1)
				}
				db_wallets = append(db_wallets, nkn_wallet)
				db_wallets_id = append(db_wallets_id, id)
			}

			dont_remove_ids := make([]int, 0)
			db_len := len(db_wallets)

			for i, _ := range wallets_passed {
				if x := FindStringInSlice(db_wallets, wallets_passed[i]); x == db_len {
					if err = addWallet(wallets_passed[i]); err != nil {
						return o.WsError(q, 1)
					}
				} else {
					dont_remove_ids = append(dont_remove_ids, db_wallets_id[x])
				}
			}
	
			y := len(dont_remove_ids)
			for i, _ := range db_wallets_id {
				if x := FindIntInSlice(dont_remove_ids, db_wallets_id[i]); x == y {
					if err = rmWallet(db_wallets_id[i]); err != nil {
						return o.WsError(q, 1)
					}
				}
			}
		}
	}

	if err = tx.Commit(); err != nil {
		o.log.Syslog("Cannot Commit Tx Query: " + err.Error(), "sql")
		return o.WsError(q, 1)
	}
	m := map[string]interface{}{}
	m["Data"] = "All settings saved"
	return nil, WSReply{Method: q.Method, Code: 0, Value: m,}
}


func (o *NKNOVH) apiRmNodes(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	var nodes []int
	var ok bool
	var node_id int
	var raw_node float64
	raw_nodes := make([]interface{}, 0)

	if raw_nodes, ok = q.Value["NodesId"].([]interface{}); !ok {
		if node_id_string, ok := q.Value["NodesId"].(string); !ok {
			return o.WsError(q, 15)
		} else {
			raw_nodes_s := strings.Split(node_id_string, ",")
			for i, _ := range raw_nodes_s {
				x, err := strconv.Atoi(strings.TrimSpace(raw_nodes_s[i]))
				if err != nil {
					return o.WsError(q, 15)
				}
				raw_nodes = append(raw_nodes, x)
			}
		}
	}
	for i, _ := range raw_nodes {
		if raw_node, ok = raw_nodes[i].(float64); !ok {
			if x, ok := raw_nodes[i].(int); !ok {
				return o.WsError(q, 15)
			} else {
				node_id = x
			}
		} else {
			node_id = int(raw_node)
		}
		nodes = append(nodes, int(node_id))
	}

	if len(nodes) < 1 {
		return o.WsError(q, 15)
	}
	tx, err := o.sql.db["main"].Begin()
	if err != nil {
		o.log.Syslog("Cannot create new Tx: " + err.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer tx.Rollback()
	for i,_ := range nodes {
		res, err := tx.Stmt(o.sql.stmt["main"]["WebRmNodes"]).Exec(c.HashId, nodes[i])
		if err != nil {
			o.log.Syslog("Cannot execute Tx stmt query: " + err.Error(), "sql")
			return o.WsError(q, 1)
		}
		if rows_affected, err := res.RowsAffected(); rows_affected == 0 && err == nil {
			o.log.Syslog("No rows affected by removing node", "sql")
			return o.WsError(q, 18)
		} else if err != nil {
			o.log.Syslog("Cannot get RowsAffected: " + err.Error(), "sql")
			return o.WsError(q, 1)
		}
	}

	err = tx.Commit()
	if err != nil {
		o.log.Syslog("Cannot Commit Tx Query: " + err.Error(), "sql")
		return o.WsError(q, 1)
	}

	m := map[string]interface{}{}
	m["Data"] = "Nodes removed successfully"
	m["NodesId"] = nodes
	return nil, WSReply{Method: q.Method, Code: 0, Value: m, }
}

func (o *NKNOVH) apiAddNodes(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	raw_ip := q.Value["Ip"]
	raw_multiple := q.Value["Multiple"]
	raw_name := q.Value["Name"]
	regexp_name := `^(\w*)$`
	nodename_max := 32
	nodes_limit := 5000

	type NodeIps struct {
		Multi []string
		MultiIP []net.IP
		One string
	}

	var node_name string
	var multiple bool
	var ok bool
	var nodes_count int
	var tmp_ip string
	ips := new(NodeIps)

	if raw_ip == nil || raw_multiple == nil || raw_name == nil {
		return o.WsError(q, 5)
	}
	if multiple, ok = raw_multiple.(bool); !ok {
		if sbool, ok := raw_multiple.(string); !ok {
			return o.WsError(q, 7)
		} else {
			x, err := strconv.ParseBool(sbool)
			if err != nil {
				return o.WsError(q, 7)
			}
			multiple = x
		}
	}
	if node_name, ok = raw_name.(string); !ok {
		return o.WsError(q, 5)
	}
	if tmp_ip, ok = q.Value["Ip"].(string); !ok {
		return o.WsError(q, 8)
	}
	re_name := regexp.MustCompile(regexp_name)

	if node_name != "" {
		if len([]rune(node_name)) > nodename_max {
			return o.WsError(q, 6)
		}
		if ok = re_name.MatchString(node_name); !ok {
			return o.WsError(q, 6)
		}
	}

	CountNodesByHash := func(x int) (error, int) {
		var cnt int
		row := o.sql.stmt["main"]["WebCountNodesByHash"].QueryRow(x)
		err := row.Scan(&cnt)
		if err == sql.ErrNoRows {
			return errors.New("No row found"), -1
		}
		if err != nil {
			return err, -1
		}
		return nil, cnt
	}

	//-1: error, 1:  inserted, 0: not inserted but err nil (IGNORE INTO)
	InsertNode := func(hashid int, name string, ip string, RawTx ...*sql.Tx) (err error, status int) {
		var tx *sql.Tx
		var res sql.Result
		var cnt int64
		status = -1
		if len(RawTx) > 0 {
			tx = RawTx[0]
			res, err = tx.Stmt(o.sql.stmt["main"]["WebInsertNode"]).Exec(hashid, name, ip)
		} else {
			res, err = o.sql.stmt["main"]["WebInsertNode"].Exec(hashid, name, ip)
		}
		if err != nil {
			return
		}

		if cnt, err = res.RowsAffected(); cnt == 0 {
			if err != nil {
				return
			}
			status = 0
			return
		}
		return nil, 1
	}

	if !multiple {
		ip := net.ParseIP(tmp_ip)
		if ip == nil {
			return o.WsError(q, 9)
		}
		if x := ip.To4(); x == nil {
			return o.WsError(q, 10)
		}
		if _, ok = IsPrivateIPv4(ip); ok {
			return o.WsError(q, 11)
		} 
		if err, nodes_count = CountNodesByHash(c.HashId); err != nil {
			return o.WsError(q, 1)
		}
		var x int
		if x = nodes_count+1; x > nodes_limit {
			return o.WsError(q, 12)
		}
		if node_name == "" {
			node_name = fmt.Sprintf("%s_%d", "Node", x)
		}
		ips.One = ip.String()
		if err, status := InsertNode(c.HashId, node_name, ips.One); err != nil {
			o.log.Syslog("InsertNode returned err:" + err.Error(), "sql")
			return o.WsError(q, 1)
		} else {
			if status == 0 {
				return o.WsError(q, 14)
			}
		}
		m := map[string]interface{}{}
		m["Info"] = "Your node added"
		return nil, WSReply{Method: q.Method, Value: m}
	}

	//Multiple
	tmp_ip = strings.TrimSpace(tmp_ip)
	if ok = strings.Contains(tmp_ip, ","); ok {
		ips.Multi = strings.Split(tmp_ip, ",")
	} else if ok = strings.Contains(tmp_ip, "\n"); ok {
		ips.Multi = strings.Split(tmp_ip, "\n")
	} else if ok = strings.Contains(tmp_ip, " "); ok {
		ips.Multi = strings.Split(tmp_ip, " ")
	} else {
		return o.WsError(q, 13)
	}

	if err, nodes_count = CountNodesByHash(c.HashId); err != nil {
		return o.WsError(q, 1)
	}

	nodeslenadd := len(ips.Multi)
	if x := nodes_count+nodeslenadd; x > nodes_limit {
		return o.WsError(q, 12)
	}
	var prefix string
	if node_name == "" {
		prefix = "Node_"
	} else {
		prefix = node_name
	}
	for i, _ := range ips.Multi {
		ip := net.ParseIP(strings.TrimSpace(ips.Multi[i]))
		if ip == nil {
			return o.WsError(q, 9)
		}
		if x := ip.To4(); x == nil {
			return o.WsError(q, 10)
		}
		if _, ok = IsPrivateIPv4(ip); ok {
			return o.WsError(q, 11)
		} 
		ips.MultiIP = append(ips.MultiIP, ip)
	}

	tx, err := o.sql.db["main"].Begin()
	if err != nil {
		o.log.Syslog("Cannot Begin Tx: " + err.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer tx.Rollback()
	var partially bool
	var cnt_nodes_added int = 0
	for i,_ := range ips.MultiIP {
		clearprefix := fmt.Sprintf("%s%d", prefix, nodes_count+i)
		if err, status := InsertNode(c.HashId, clearprefix, ips.MultiIP[i].String(), tx); err != nil {
			o.log.Syslog("InsertNode returned err:" + err.Error(), "sql")
			return o.WsError(q, 1)
		} else {
			if status == 0 {
				partially = true
				continue
			}
			cnt_nodes_added++
		}
	}

	if partially == true && cnt_nodes_added > 0 {
		if err = tx.Commit(); err != nil {
			o.log.Syslog("Cannot commit transaction: " + err.Error(), "sql")
			return o.WsError(q, 1)
		}
		return o.WsError(q, 2)
	} else if partially == true && cnt_nodes_added == 0 {
		return o.WsError(q, 14)
	}
	if err = tx.Commit(); err != nil {
		o.log.Syslog("Cannot commit transaction: " + err.Error(), "sql")
		return o.WsError(q, 1)
	}
	m := map[string]interface{}{}
	m["Info"] = "Nodes added"
	return nil, WSReply{Method: q.Method, Code: 0, Value: m }
}

func (o *NKNOVH) apiGenId(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	var cnt int
	row := o.sql.stmt["main"]["WebCheckIPCreator"].QueryRow(c.Ip)
	if errx := row.Scan(&cnt); errx != nil {
		o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	if cnt >= 3 {
		o.log.Syslog("Limit exceeded (genId) from IP: " + c.Ip, "info")
		return o.WsError(q, 252)
	}
	errx, hash := GenRandomSHA256()
	if errx != nil {
		o.log.Syslog("GenRandomSHA256 returned error: " + err.Error(), "info")
		return o.WsError(q, 240)
	}
	_, errx = o.sql.stmt["main"]["WebCreateUniq"].Exec(hash, c.Ip)
	if errx != nil {
		o.log.Syslog("Cannot exec query: "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	data := map[string]interface{}{}
	data["Hash"] = hash
	r = WSReply{Method: q.Method, Code: 0, Value: data}
	return
}

func (o *NKNOVH) apiAuth(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	var hash string
	var ok bool
	hash, ok = q.Value["Hash"].(string)
	if !ok {
		//DEPRECATED AND WILL BE REMOVED
		hash, ok = q.Value["hash"].(string)
		if !ok {
			return o.WsError(q, 5)
		}
	}

	if len(hash) != 64 {
		return o.WsError(q, 254)
	}
	var id int
	row := o.sql.stmt["main"]["WebSelectUniqByHash"].QueryRow(&hash)
	errx := row.Scan(&id)
		switch {
		case errx == sql.ErrNoRows:
			return o.WsError(q, 255)
		break
		case errx != nil:
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		break
		}
	c.HashId = id
	value := map[string]interface{}{}
	value["Hash"] = hash
	return err, WSReply{Method: q.Method, Code: 0, Value: value,}
}


func (o *NKNOVH) apiFullstack(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	if c.HashId == -1 {
		return o.WsError(q, 253)
	}
	err, netstatus := o.apiNetstatus(q, c)
	if err != nil {
		return o.WsError(q, 3)
	}
	_, daemon := o.apiDaemon(q, c)
	_, wallets := o.apiMyWallets(q, c)
	_, prices := o.apiPrices(q, c)
	_, nodes := o.apiMyNodes(q, c)

	if (nodes.Code != 0 && nodes.Code != 3) || (wallets.Code != 0 && wallets.Code != 3) || daemon.Code != 0 {
		return o.WsError(q, 4)
	}
	x := map[string]interface{}{}
	x["Netstatus"] = netstatus
	x["Daemon"] = daemon
	x["Wallets"] = wallets
	x["Prices"] = prices
	x["Nodes"] = nodes
	r = WSReply{Method: q.Method, Code: 0, Value: x}

	return
}

func (o *NKNOVH) apiLanguage(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	lang_packages := []string{"en_US", "ru_RU"}
	var locale string
	var view string
	var ok bool

	if view, ok = q.Value["View"].(string); !ok {
		return o.WsError(q, 230)
	}
	if locale, ok = q.Value["Locale"].(string); !ok {
		return o.WsError(q, 231)
	}
	if len(locale) > 10 || len(view) > 32 {
		return o.WsError(q, 232)
	}
	
	if i := FindStringInSlice(lang_packages, locale); i == len(lang_packages) {
		return o.WsError(q, 233)
	}
	read, err := ioutil.ReadFile("templates/languages/" + locale + ".json")
	if err != nil {
		o.log.Syslog("Cannot read a file: " + err.Error(), "main")
		return o.WsError(q, 234)
	}

	data := string(read)
	m := map[string]interface{}{}
	m["Locale"] = locale
	m["View"] = view
	m["Data"] = data
	r = WSReply{Method: q.Method, Code: 0, Value: m, }
	return
}

func (o *NKNOVH) apiDaemon(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	var norows bool = true
	rows, errx := o.sql.stmt["main"]["WebGetDaemon"].Query()
	if errx != nil {
		o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer rows.Close()
	m := map[string]string{}
	var name string
	var value string
	for rows.Next() {
		norows = false
		if errx := rows.Scan(&name, &value); errx != nil {
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		}
		m[name] = value
	}
	if norows {
		return o.WsError(q, 3)
	}
	r = WSReply{Method: q.Method, Code: 0, Value: m, }
	return
}

func (o *NKNOVH) apiMyNodes(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	rows, errx := o.sql.stmt["main"]["WebGetMyNodes"].Query(c.HashId)
	if errx != nil {
		o.log.Syslog("Can't execute Query WebGetMyNodes: "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer rows.Close()
	var norows bool = true
	nodes := make([]map[string]interface{}, 0)
	nodes_id := map[int]int{}
	cnt := map[int]int{}
	n := -1
	no_history_nodes := map[int]bool{}
	tmp := map[int]map[int]map[string]uint64{}

	type nodeS struct {
		Id int
		Name string
		Ip string
	}
	type nodeLastInfo struct {
		NodeId int
		NID string
		Currtimestamp uint64
		Height uint64
		ProposalSubmitted int
		ProtocolVersion int
		RelayMessageCount uint64
		SyncState string
		Uptime int
		Version string
		LatestUpdate string
	}

	loop:
	for rows.Next() {
		norows = false
		node := new(nodeS)
		if errx := rows.Scan(&node.Id, &node.Name, &node.Ip); errx != nil {
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		}
		n++
		nodes = append(nodes, map[string]interface{}{})
		nodes[n]["NodeId"] = node.Id
		nodes[n]["Ip"] = node.Ip
		nodes[n]["Name"] = node.Name
		nodes_id[node.Id] = n
		cnt[n] = 0

		nli := new(nodeLastInfo)
		rows2 := o.sql.stmt["main"]["WebGetMyNodeLastInfo"].QueryRow(node.Id)
		errx := rows2.Scan(&nli.NodeId, &nli.NID, &nli.Currtimestamp, &nli.Height, &nli.ProposalSubmitted, &nli.ProtocolVersion, &nli.RelayMessageCount, &nli.SyncState, &nli.Uptime, &nli.Version, &nli.LatestUpdate)
		switch {
			case errx == sql.ErrNoRows:
				nodes[n]["Err"] = 2
				nodes[n]["SyncState"] = "Waiting for first update"
				nodes[n]["LatestUpdate"] = time.Now().Format("2006-01-02 15:04:05")
				no_history_nodes[n] = true
				continue loop
			case errx != nil:
				o.log.Syslog("Can't execute Query WebGetMyNodeLastInfo: "+errx.Error(), "sql")
				return o.WsError(q, 1)
		}
		nodes[n]["SyncState"] = nli.SyncState
		nodes[n]["Uptime"] = nli.Uptime
		nodes[n]["Height"] = nli.Height
		nodes[n]["Version"] = nli.Version
		nodes[n]["RelayMessageCount"] = nli.RelayMessageCount
		nodes[n]["Currtimestamp"] = nli.Currtimestamp		
		nodes[n]["ProposalSubmitted"] = nli.ProposalSubmitted
		nodes[n]["LatestUpdate"] = nli.LatestUpdate

		if nli.Uptime > 0 {
			nodes[n]["RelaysPerHour"] = math.Floor(float64(nli.RelayMessageCount)/float64(nli.Uptime)*3600)
		} else {
			nodes[n]["RelaysPerHour"] = 0
		}

		nodes[n]["RelaysPerHour10"] = 0
		nodes[n]["RelaysPerHour60"] = 0
		if nli.SyncState == "OFFLINE" {
			nodes[n]["Err"] = 1
			nodes[n]["ProposalSubmitted"] = -1
			nodes[n]["SyncState"] = "_OFFLINE_"
			no_history_nodes[n] = true
		} else if nli.SyncState == "OUT" {
			nodes[n]["SyncState"] = "_OUT_"
			nodes[n]["Err"] = 3
		} else if nli.SyncState == "PRUNING DB" || nli.SyncState == "GENERATING ID" {
			no_history_nodes[n] = true
		}
	}

	if norows {
		return o.WsError(q, 3)
	}
	var sx string
	x := n+1
	inSx := make([]interface{}, x, x)
	nodesIdKeys := make([]int, 0, x)
	for key := range nodes_id {
		nodesIdKeys = append(nodesIdKeys, key)
	}
	for i := 0; i < x; i++ {
		inSx[i] = nodesIdKeys[i]
		if diff := x - i; diff == 1 {
			sx += "?"
			continue
		}
		sx += "?,"
	}
	sqlHistory := "SELECT node_id,Currtimestamp,RelayMessageCount,Uptime FROM nodes_history WHERE node_id in (" + sx + ") ORDER BY id DESC"
	stmt, errx := o.sql.db["main"].Prepare(sqlHistory)
	if errx != nil {
		o.log.Syslog("Can't Prepare sqlHistory: "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	norows = true
	rows3, errx := stmt.Query(inSx...)
	if errx != nil {
		o.log.Syslog("Can't execute query sqlhistory: "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer rows3.Close()
	for rows3.Next() {
		norows = false
		h := new(nodeLastInfo)
		if errx := rows3.Scan(&h.NodeId, &h.Currtimestamp, &h.RelayMessageCount, &h.Uptime); errx != nil {
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		}
		n := nodes_id[h.NodeId]
		if _, ok := no_history_nodes[n]; ok {
			continue
		}
		if tmp[n] == nil {
			tmp[n] = map[int]map[string]uint64{}
		}
		tmp[n][cnt[n]] = map[string]uint64{}
		tmp[n][cnt[n]]["Uptime"] = uint64(h.Uptime)
		tmp[n][cnt[n]]["Currtimestamp"] = h.Currtimestamp
		tmp[n][cnt[n]]["Relays"] = h.RelayMessageCount
		cnt[n]++
	}
	if norows {
		m := map[string]interface{}{}
		m["List"] = nodes
		r = WSReply{Method: q.Method, Code: 0, Value: m,}
		return
	}
	setRelays := func(val map[int]map[string]uint64, rh int, diff int, elem int, nodes []map[string]interface{}, key int) {
		var rtype string
		if rh == 600 {
			rtype = "RelaysPerHour10"
		} else if rh == 3600 {
			rtype = "RelaysPerHour60"
		} else {
			return
		}

		i := 0
		max_diff := float64(diff)
		ipos := i+elem
		if _, ok := val[ipos]; ok {
			diff_uptime := float64(rh) - (float64(val[i]["Uptime"]) - float64(val[ipos]["Uptime"]))
			diff_timestamp := float64(rh) - (float64(val[i]["Currtimestamp"]) - float64(val[ipos]["Currtimestamp"]))
			diff_general := diff_timestamp - diff_uptime
			if (diff_general >= 0 && diff_general <= 10) != (diff_general < 0 && diff_general >= -10) {
				if (diff_uptime <= max_diff && diff_uptime >= 0) != (diff_uptime < max_diff && diff_uptime < 0 && (max_diff+diff_uptime) > 0) {
						if (val[i]["Uptime"] > val[ipos]["Uptime"]) && (val[ipos]["Uptime"] > 0) {
							nodes[key][rtype] = math.Floor((float64(val[i]["Relays"]) - float64(val[ipos]["Relays"]))/float64(rh)*3600);
							return
						} else {
							nodes[key][rtype] = -1;
							return
						}
				} else {
					nodes[key][rtype] = -1;
					return
				}
			}
		}

	}

	opts := []map[string]int{}
	opt := map[string]int{"rh": 600, "max_diff": 60, "elem": 1}
	opt2 := map[string]int{"rh": 3600, "max_diff": 60, "elem": 6}
	opts = append(opts, opt, opt2)

	for key, _ := range tmp {
		l := len(opts)
		for z := 0; z < l; z++ {
			setRelays(tmp[key], opts[z]["rh"], opts[z]["max_diff"], opts[z]["elem"], nodes, key);
		}
	}
	m := map[string]interface{}{}
	m["List"] = nodes
	r = WSReply{Method: q.Method, Code: 0, Value: m,}
	return
}


func (o *NKNOVH) apiPrices(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	rows, errx := o.sql.stmt["main"]["WebGetPrices"].Query()
	if errx != nil {
		o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer rows.Close()
	var name string
	var price float64
	var norows bool = true
	in := map[string]interface{}{}
	for rows.Next() {
		norows = false
		if errx := rows.Scan(&name, &price); errx != nil {
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		}
		in[name] = price
	}
	if norows {
		return o.WsError(q, 3)
	}
	r = WSReply{Method: q.Method, Code: 0, Value: in, }
	return
}

func (o *NKNOVH) apiNetstatus(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	row := o.sql.stmt["main"]["WebGetNetStatus"].QueryRow()
	obj := new(Netstatus)
	errx := row.Scan(&obj.Relays, &obj.AverageUptime, &obj.AverageRelays, &obj.RelaysPerHour, &obj.ProposalSubmitted, &obj.PersistNodesCount, &obj.NodesCount, &obj.LastHeight, &obj.LastTimestamp, &obj.AverageBlockTime, &obj.AverageBlocksPerDay, &obj.LatestUpdate)
		switch {
		case errx == sql.ErrNoRows:
			return o.WsError(q, 3)
		break
		case errx != nil:
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		break
		}
	resp := WSReply{Method: q.Method, Code: 0, Value: obj, }
	return err, resp
}

func (o *NKNOVH) apiMyWallets(q *WSQuery, c *CLIENT) (err error, r WSReply) {
	rows, errx := o.sql.stmt["main"]["WebGetMyWallets"].Query(c.HashId)
	if errx != nil {
		o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
		return o.WsError(q, 1)
	}
	defer rows.Close()
	var norows bool = true
	in := []interface{}{}
	var id int64
	var wallet string
	var balance float64

	for rows.Next() {
		norows = false
		if errx := rows.Scan(&id, &wallet, &balance); errx != nil {
			o.log.Syslog("Can't execute row.Scan(): "+errx.Error(), "sql")
			return o.WsError(q, 1)
		}
		m := map[string]interface{}{}
		m["Id"] = id
		m["NknWallet"] = wallet 
		m["Balance"] = balance
		in = append(in, m)
	}
	if norows {
		return o.WsError(q, 3)
	}
	mv := map[string]interface{}{}
	mv["Wallets"] = in
	r = WSReply{Method: q.Method, Code: 0, Value: mv, }
	return
}
