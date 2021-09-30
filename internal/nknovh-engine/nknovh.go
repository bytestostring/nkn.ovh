package nknovh_engine

import (
	 "database/sql"
	 "time"
	 "strconv"
	 "encoding/json"
	 "sync"
	 "regexp"
	 "strings"
	 "reflect"
	 "net/http"
	 "runtime"
	 "sort"
	 "errors"
	 // "fmt"
	 )

type NKNOVH struct {
	conf *configuration
	log *logger
	sql *Mysql
	NodeInfo *NodeInfo
	threads *Threads
	http *Http
	Nknsdk *Nknsdk
	Web *Web
	Reporter *Reporter
	Validator *Validator
}


type Http struct {
	NeighborClient *http.Client
	MainClient *http.Client
	DirtyClient *http.Client
}

type Threads struct {
	Counter int
	Neighbors chan struct{}
	Main chan struct{}
	Dirty chan struct{}
	Flush sync.Mutex
}

type Reporter struct {
	dirty map[string]*DBNode
	main map[string]*DBNode
	mu_dirty sync.RWMutex
	mu_main sync.RWMutex
	dirtyReady chan bool
	mainReady chan bool
}


type RPCError struct {
	Code int
	Status string
	Description string
	WalletAddress string
	PublicKey string
}

type NodeInfo struct {
	ips []string
	m_nodes map[string][]uint64
	d_nodes map[string][]uint64
	dirty_keys []string
	States []*NodeState
	Neighbors []*NodeNeighbor
	CounterFinish int
	mux sync.RWMutex
	ANLast map[string]float64
	ANLastMux sync.RWMutex
	ANArray map[int]map[int][]int
	ANArrayMux []sync.RWMutex
}

type DBNode struct {
	Ip string
	Ids []uint64
	Dirty bool
	LastStatus string
}

func (o *NKNOVH) Build() error {

	o.log = &logger{dir: "logs"}
	o.log.Init()
	// Get configuration json file
	var config *configuration
	conf, err := config.configure()
	if err != nil {
		o.log.Syslog(err.Error(), "main")
		return err
	}
	o.sql = &Mysql{log: o.log}
	o.sql.build()
	o.conf = conf
	var conslen int = len(o.conf.Db)

	//Creating connections mysql/pgsql
	for i := 0; i < conslen; i++ {
		if err := o.sql.createConnect(o.conf.Db[i].Host, o.conf.Db[i].DbType, o.conf.Db[i].Login, o.conf.Db[i].Password, o.conf.Db[i].Database, o.conf.Db[i].MaxOpenConns, o.conf.Db[i].MaxIdleConns, o.conf.Db[i].InsideName); err != nil {
			return err
		}
	}

	//Prepare mysql queries
	if err := o.sql.prepare(); err != nil {
		o.log.Syslog("o.sql.prepare() has returned error: ("+err.Error()+")", "sql")
		return err
	}

	if err := o.updateConfig("Version", o.conf.Version); err != nil {
		o.log.Syslog("o.updateConfig has returned error: ("+err.Error()+")", "main")
		return err
	}

	timezone := time.Now().Format("Z0700")
	if err := o.updateConfig("Timezone", timezone); err != nil {
	o.log.Syslog("o.updateConfig has returned error: ("+err.Error()+")", "main")
		return err
	}

	//Create NKN wallet for NKN SDK API
	if err := o.walletCreate(); err != nil {
		return err
	}
	if err := o.nknConnect(); err != nil {
		return err
	}

	o.Validator = buildValidator()

	o.threads = &Threads{
							Neighbors: make(chan struct{}, o.conf.Threads.Neighbors),
							Main: make(chan struct{}, o.conf.Threads.Main),
							Dirty: make(chan struct{}, o.conf.Threads.Dirty),
						}
	o.threads.Counter = 0
	o.NodeInfo = &NodeInfo{
							mux: sync.RWMutex{},
							ips: make([]string, 0),
							m_nodes: map[string][]uint64{},
							d_nodes: map[string][]uint64{},
							ANLast: map[string]float64{},
						}

	o.Reporter = &Reporter{
					dirty: map[string]*DBNode{},
					main: map[string]*DBNode{},
					dirtyReady: make(chan bool),
					mainReady: make(chan bool),
				}

	o.NodeInfo.ANArray = map[int]map[int][]int{}
	o.NodeInfo.ANArrayMux = make([]sync.RWMutex, 0)
	for i := 0; i < 256; i++ {
		o.NodeInfo.ANArrayMux = append(o.NodeInfo.ANArrayMux, sync.RWMutex{})
		o.NodeInfo.ANArray[i] = map[int][]int{}
		for n := 0; n < 256; n++ {
			o.NodeInfo.ANArray[i][n] = make([]int, 0)
		}
	}
	var netTransport = &http.Transport{DisableKeepAlives: true}
	o.http = &Http{
					MainClient: &http.Client{Timeout: time.Duration(o.conf.MainPoll.ConnTimeout)*time.Second, Transport: netTransport,},
					DirtyClient: &http.Client{Timeout: time.Duration(o.conf.DirtyPoll.ConnTimeout)*time.Second, Transport: netTransport,},
					NeighborClient: &http.Client{Timeout: time.Duration(o.conf.NeighborPoll.ConnTimeout)*time.Second, Transport: netTransport,},
				}
	return nil
}

func (o *NKNOVH) Run() error {

	var ch []chan bool = make([]chan bool, 4)
	
	//Run polls
	go o.reporterManager()
	go o.createPoll("neighborPoll", o.conf.NeighborPoll.Interval, ch[0], false, o.neighborPoll)
	go o.createPoll("mainPoll", o.conf.MainPoll.Interval, ch[1], true, o.mainPoll)
	go o.createPoll("walletPoll", o.conf.Wallets.Interval, ch[2], false, o.walletPoll)
	go o.createPoll("dirtyPoll", o.conf.DirtyPoll.Interval, ch[3], true, o.dirtyPoll)

	for i := 0; i < len(ch); i++ {
		select {
		case msg1 := <-ch[0]:
			if msg1 == false {
				return errors.New("Neighbors Poll has ended")
			}
		break
		case msg1 := <-ch[1]:
			if msg1 == false {
				return errors.New("Main Poll has ended")
			}
		break
		case msg1 := <-ch[2]:
			if msg1 == false {
				return errors.New("Wallet Poll has ended")
			}
		break
		case msg1 := <-ch[3]:
			if msg1 == false {
				return errors.New("Dirty Poll has ended")
			}
		break
		}
	}
	return errors.New("Any poll has exited")
}

func (o *NKNOVH) reporterManager() {

	dbnode_handling := func(nodes *map[string]*DBNode, dirty bool) error {
		if len(*nodes) < 1 {
			return nil
		}
		if dirty {
			o.Reporter.mu_dirty.Lock()
			defer o.Reporter.mu_dirty.Unlock()
		} else {
			o.Reporter.mu_main.Lock()
			defer o.Reporter.mu_main.Unlock()
		}

		defer func() {
			*nodes = map[string]*DBNode{}
		}()

		var hash_id int
		var name string
		for _, i := range *nodes {

			loop2:
			for _, node_id := range i.Ids {
				row := o.sql.stmt["main"]["selectNodeHashNameById"].QueryRow(node_id)
				err := row.Scan(&hash_id, &name)
				switch {
					case err == sql.ErrNoRows:
						continue loop2
					break
					case err != nil:
						continue loop2
					break
				}
				//row = o.sql.stmt["main"]
			}
		}
		return nil
	}

	for {
		select {
			case <- o.Reporter.dirtyReady:
				go dbnode_handling(&o.Reporter.dirty, true)
			break
			case <- o.Reporter.mainReady:
				go dbnode_handling(&o.Reporter.main, false)
			break
		}
	}
}

func (o *NKNOVH) updateConfig(name string, value string) error {

	var id uint
	row := o.sql.stmt["main"]["selectDaemonIdByName"].QueryRow(&name)
	err := row.Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		if _, err1 := o.sql.stmt["main"]["insertDaemon"].Exec(&name, &value); err1 != nil {
				o.log.Syslog("Stmt insertDaemon has returned an error: ("+err1.Error()+")", "sql")
				return err1
			}
	break
	case err != nil:
		o.log.Syslog("Can't execute row.Scan(): "+err.Error(), "sql")
		return err
	break
	default:
		if _, err1 := o.sql.stmt["main"]["updateDaemonById"].Exec(&value, &id); err1 != nil {
			o.log.Syslog("Stmt updateDaemonById has returned an error: ("+err1.Error()+")", "sql")
			return err1
		}
	}
	return nil
}


func (o *NKNOVH) getANFromDB() error {
	o.NodeInfo.ips = make([]string, 0)
	rows, err := o.sql.stmt["main"]["selectAllIpsAN"].Query()
	if err != nil { 
		return err
	}
	defer rows.Close()

	var db_ip string
	for rows.Next() {
		if err := rows.Scan(&db_ip); err != nil {
			return err
		}
		o.NodeInfo.ips = append(o.NodeInfo.ips, db_ip)
	}
	return nil
}	


func (o *NKNOVH) createPoll(pollName string, interval int, ch chan bool, even bool, f func() error) error {
	o.log.Syslog(pollName + " is starting...", "main")
	o.log.Syslog("[" + pollName + "] Waiting for syncing", "main")

	//Prepare variables
	var dtime time.Duration = time.Duration(interval)
	var lap uint64 = 1
	var iteration_start time.Time
	var iteration_end time.Time
	var iteration_time time.Duration
	var iteration_all time.Duration = 0
	var iteration_average time.Duration = 0
	var inside_error bool

	// Sync
	sleeping(time.Duration(time.Now().Second())*time.Second + time.Duration(time.Now().Nanosecond())*time.Nanosecond, 60, even)
	o.log.Syslog(pollName + " loop is starting!", "main")
	for {
		iteration_start = time.Now()
		if inside_error == true {
			inside_error = false
			sleeping(iteration_time, dtime, even)
		}

		//Run something
		if err := f(); err != nil {
			inside_error = true
			iteration_time = time.Now().Sub(iteration_start)
			o.log.Syslog("["+pollName+"] Got an error: " + err.Error(), "main")
			continue
		}

		iteration_end = time.Now()
		iteration_time = iteration_end.Sub(iteration_start)
		iteration_all = iteration_time + iteration_all
		iteration_average = iteration_all / time.Duration(lap)
		o.log.Syslog("["+strconv.FormatUint(lap, 10)+"] ["+pollName+"] Cycle iteration time: "+iteration_time.String()+" (Average: "+iteration_average.String()+"); Sleeping", "cycles")
		lap++
		sleeping(iteration_time, dtime, even)
	}
	ch <- false
	return nil
}


func (o *NKNOVH) dirtyPoll() error {
	defer func() {
		o.Reporter.dirtyReady <- true
	}()

	if err := o.getNodesFromDB(true); err != nil {
		return err
	}
	if err := o.fetchNodesInfo(true); err != nil {
		return err
	}
	o.rmNodesByFcnt(180, 1)
	o.rmNodesByFcnt(5040, 0)
	return nil
}

func (o *NKNOVH) mainPoll() error {
	defer func() {
		o.Reporter.mainReady <- true
	}()
	if err := o.getNodesFromDB(false); err != nil {
		return err
	}
	if err := o.fetchNodesInfo(false); err != nil {
		return err
	}
	return nil
}

func (o *NKNOVH) neighborPoll() error {

	if _, err := o.sql.stmt["main"]["clearAN"].Exec(); err != nil {
		return err
	}
	if err := o.dbIpsToArray(); err != nil {
		return err
	}
	if err := o.updateAN(); err != nil {
		return err
	}
	if err := o.swapAndClearAN(); err != nil {
		return err
	} 
	if err := o.saveANStatus(); err != nil {
		return err
	}
	return nil
}

func (o *NKNOVH) swapAndClearAN() error {
	tx, err := o.sql.db["main"].Begin()

	if err != nil {
		o.log.Syslog("Cannot create tx: " + err.Error(), "sql")
		return err
	}

	o.NodeInfo.ANLastMux.Lock()
	defer o.NodeInfo.ANLastMux.Unlock()
	defer tx.Rollback()

	if _, err = tx.Stmt(o.sql.stmt["main"]["clearANStats"]).Exec(); err != nil {
		o.log.Syslog("Cannot execute clearANStats: " + err.Error(), "sql")
		return err
	}
	if _, err = tx.Stmt(o.sql.stmt["main"]["copyANtoStats"]).Exec(); err != nil {
		o.log.Syslog("Cannot execute copyANStats: " + err.Error(), "sql")
		return err
	}

	if _, err = tx.Stmt(o.sql.stmt["main"]["clearAN"]).Exec(); err != nil {
		o.log.Syslog("Cannot execute clearAN: " + err.Error(), "sql")
		return err
	}
	if err = tx.Commit(); err != nil {
		o.log.Syslog("Cannot commit transaction (swapAndClearAN): " + err.Error(), "sql")
		return err
	}
	return nil
}

func (o *NKNOVH) fetchNodesInfo(dirty bool) error {
	var wg sync.WaitGroup
	var nodes_list *map[string][]uint64
	var http_client *http.Client
	var threads *chan struct {}

	if dirty {
		nodes_list = &o.NodeInfo.d_nodes
		http_client = o.http.DirtyClient
		threads = &o.threads.Dirty
		l := len(o.NodeInfo.dirty_keys)
		for i := 0; i < l; i++ {
			dbnode := new(DBNode)
			dbnode.Ip = o.NodeInfo.dirty_keys[i]
			dbnode.Ids = (*nodes_list)[dbnode.Ip]
			dbnode.Dirty = dirty
			r := &JsonRPCConf{Ip:dbnode.Ip, Method:"getnodestate", Params: &json.RawMessage{'{','}'}, Client: http_client,}
			wg.Add(1)
			*threads <- struct{}{}
			go o.getInfo(&wg, r, "UpdateNode", threads, dbnode)	
		}
	} else {
		nodes_list = &o.NodeInfo.m_nodes
		http_client = o.http.MainClient
		threads = &o.threads.Main
		for k, v := range *nodes_list {
			dbnode := new(DBNode)
			dbnode.Ip = k
			dbnode.Ids = v
			dbnode.Dirty = dirty
			r := &JsonRPCConf{Ip:k, Method:"getnodestate", Params: &json.RawMessage{'{','}'}, Client: http_client,}
			wg.Add(1)
			*threads <- struct{}{}
			go o.getInfo(&wg, r, "UpdateNode", threads, dbnode)
		}
	}
	wg.Wait()
	num_routines := runtime.NumGoroutine()
	o.log.Syslog("func fetchNodesInfo() is finished, active numbers of goroutines: " + strconv.Itoa(num_routines), "main")
	return nil
}

func (o *NKNOVH) isOutOfNetwork(dbnode *DBNode, node *NodeState) (error, bool) {
	var id uint64
	var timestamp uint64
	var last_timestamp float64
	var last_height float64
	var average_blockTime float64
	var ok bool
	var correction uint64 = 300
	var diff_timestamp uint64
	var min_block_difference float64 = 7

	if node.Result.Uptime < 1200 {
		return nil, false
	}
	
	timestamp = uint64(time.Now().Unix())
	diff := timestamp - uint64(node.Result.Uptime) + correction
	
	o.NodeInfo.ANLastMux.RLock()
	defer o.NodeInfo.ANLastMux.RUnlock()
	if last_height, ok = o.NodeInfo.ANLast["Height"]; !ok {
		return nil, false
	}
	if last_timestamp, ok = o.NodeInfo.ANLast["Timestamp"]; !ok {
		return nil, false
	}
	if average_blockTime, ok = o.NodeInfo.ANLast["averageBlockTime"]; !ok {
		return nil, false
	}
	row := o.sql.stmt["main"]["selectIdByIpANLast"].QueryRow(dbnode.Ip)
	err := row.Scan(&id)
	switch {
		case err == sql.ErrNoRows:
			var diff_height uint64
			var node_height uint64 = uint64(node.Result.Height)
			if node_height > uint64(last_height) {
				diff_height = node_height - uint64(last_height)
			} else {
				diff_height = uint64(last_height) - node_height
			}
			diff_timestamp = timestamp - uint64(last_timestamp)
			block_difference := float64(diff_height)  - float64(diff_timestamp) / average_blockTime 
			if diff < uint64(last_timestamp) && (block_difference > min_block_difference || block_difference < min_block_difference - min_block_difference * 2) {
				return nil, true
			}
		case err != nil:
			return err, false
	}
	return nil, false
}


func (o *NKNOVH) UpdateNode(node *NodeState, params interface{}) {

	if b := o.Validator.IsNodeStateValid(node); !b {
		o.log.Syslog("isNodeStateValid has returned false", "nodes")
		return
	}

	dbnode := params.(*DBNode)
	minute := time.Now().Minute()
	var id uint64
	var failcnt int64
	var ftf uint8
	for _, node_id := range dbnode.Ids {
		//Add historical data once per 10 minutes
		if d := minute % 10; d == 0 {
			if _, err := o.sql.stmt["main"]["insertNodeStats"].Exec(&node_id, &node.Result.ID, &node.Result.Currtimestamp, &node.Result.Height, &node.Result.ProposalSubmitted, &node.Result.ProtocolVersion, &node.Result.RelayMessageCount, &node.Result.SyncState, &node.Result.Uptime, &node.Result.Version); err != nil {
					o.log.Syslog("Stmt insertNodeStats has returned an error: ("+err.Error()+")", "sql")
				}
			o.rmOldHistory(node_id)
		}


		// Exclude a node from dirtyPoll if it is in
		if _, err := o.sql.stmt["main"]["updateNodeToMain"].Exec(&node_id); err != nil {
			o.log.Syslog("Stmt updateNodeToMain has returned an error: ("+err.Error()+")", "sql")
		}

		//Detect out of NKN Network
		var isNodeOuted bool
		if err, b := o.isOutOfNetwork(dbnode, node); err == nil && b == true {
			node.Result.SyncState = "OUT"
			isNodeOuted = true
		}

		if dbnode.Dirty {
			o.Reporter.mu_dirty.Lock()
			dbnode.LastStatus = node.Result.SyncState
			if _, b := o.Reporter.dirty[dbnode.Ip]; !b {
				o.Reporter.dirty[dbnode.Ip] = dbnode
			}
			o.Reporter.mu_dirty.Unlock()
		} else if isNodeOuted {
			o.Reporter.mu_main.Lock()
			dbnode.LastStatus = node.Result.SyncState
			if _, b := o.Reporter.main[dbnode.Ip]; !b {
				o.Reporter.main[dbnode.Ip] = dbnode
			}
			o.Reporter.mu_main.Unlock()
		}

		//Add the last data
		row := o.sql.stmt["main"]["selectNodeLastIdByNodeId"].QueryRow(&node_id)
		err := row.Scan(&id, &failcnt, &ftf)
		switch {
		case err == sql.ErrNoRows:
			if _, err1 := o.sql.stmt["main"]["insertNodeLast"].Exec(node_id, &node.Result.ID, &node.Result.Currtimestamp, &node.Result.Height, &node.Result.ProposalSubmitted, &node.Result.ProtocolVersion, &node.Result.RelayMessageCount, &node.Result.SyncState, &node.Result.Uptime, &node.Result.Version, 0, 0); err1 != nil {
				o.log.Syslog("Stmt insertNodeLast has returned an error: ("+err1.Error()+")", "sql")
			}
		break
		case err != nil:
			o.log.Syslog("Can't execute row.Scan(): "+err.Error(), "sql")
		break
		default:
			if _, err1 := o.sql.stmt["main"]["updateNodeLastById"].Exec(&node.Result.ID, &node.Result.Currtimestamp, &node.Result.Height, &node.Result.ProposalSubmitted, &node.Result.ProtocolVersion, &node.Result.RelayMessageCount, &node.Result.SyncState, &node.Result.Uptime, &node.Result.Version, 0, 0, &id); err1 != nil {
				o.log.Syslog("Stmt updateNodeLastById has returned an error: ("+err1.Error()+")", "sql")
			}
		}

	}
	return
}

func (o *NKNOVH) UpdateNodeFail(answer []byte, params interface{}) error {
	dbnode := params.(*DBNode)

	var node_ip string = dbnode.Ip
	var repeatInterval time.Duration
	var repeatTimeout time.Duration
	var threads *chan struct{}

	if dbnode.Dirty == true {
		threads = &o.threads.Dirty
		repeatInterval = 0
		repeatTimeout = 2
	} else {
		threads = &o.threads.Main
		repeatInterval = 3
		repeatTimeout = 6
	}

	netTransport := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Timeout: repeatTimeout*time.Second,Transport: netTransport,}

	//dummy, no routines
	var wg sync.WaitGroup
	r := &JsonRPCConf{Ip:node_ip, Method:"getnodestate", Params: &json.RawMessage{'{','}'}, Client: client}
	for i := 1; i < 4; i++ {
		time.Sleep(repeatInterval * time.Second)
		wg.Add(1)
		if err := o.getInfo(&wg, r, "UpdateNode", threads, params, true); err != nil {
			o.log.Syslog("[Retry " + strconv.Itoa(i) + "] The node \"" + node_ip + "\" has no responded", "nodes")
			continue
		}
		o.log.Syslog("Node \"" + node_ip + "\" is working up now!", "main")
		return nil
	}
	

	// Node gonna offline and into the dirty poll
	if !dbnode.Dirty {
		o.Reporter.mu_main.Lock()
		dbnode.LastStatus = "OFFLINE"
		if _, b := o.Reporter.main[dbnode.Ip]; !b {
			o.Reporter.main[dbnode.Ip] = dbnode
		}
		o.Reporter.mu_main.Unlock()
	}

	var last_id uint64
	var failcnt int64
	var ftf uint8
	for _, node_id := range dbnode.Ids {

		// switch to dirty
		if _, err := o.sql.stmt["main"]["updateNodeToDirty"].Exec(&node_id); err != nil {
			o.log.Syslog("Stmt updateNodeToDirty has returned an error: ("+err.Error()+")", "sql")
		}
		row := o.sql.stmt["main"]["selectNodeLastIdByNodeId"].QueryRow(&node_id)
		err := row.Scan(&last_id, &failcnt, &ftf)
		switch {
		case err == sql.ErrNoRows:
			if _, err1 := o.sql.stmt["main"]["insertNodeLast"].Exec(&node_id, "", 0, 0, 0, -1, 0, "OFFLINE", 0, "", 1, 1); err1 != nil {
				o.log.Syslog("Stmt insertNodeLast has returned an error: ("+err1.Error()+")", "sql")
			}
		break
		case err != nil:
			o.log.Syslog("Can't execute row.Scan(): "+err.Error(), "sql")
		break
		default:
			if _, err1 := o.sql.stmt["main"]["updateNodeLastById"].Exec("", 0, 0, 0, -1, 0, "OFFLINE", 0, "", failcnt+1, &ftf, &last_id); err1 != nil {
				o.log.Syslog("Stmt updateNodeLastById has returned an error: ("+err1.Error()+")", "sql")
			}
		}
	}
	return nil
}

func (o *NKNOVH) getNodesFromDB(dirty bool) error {
	var switch_dirty string
	var nodes_list *map[string][]uint64
	dirty_keys := make([]string, 0)
	if dirty {
		switch_dirty = "selectAllNodesDirty"
		nodes_list = &o.NodeInfo.d_nodes
	} else {
		switch_dirty = "selectAllNodesNotDirty"
		nodes_list = &o.NodeInfo.m_nodes
	}
	rows, err := o.sql.stmt["main"][switch_dirty].Query()
	if err != nil { 
		return err
	}
	defer rows.Close()

	//clear map
	*nodes_list = map[string][]uint64{}
	var node_id uint64
	var db_ip string
	for rows.Next() {
		if err := rows.Scan(&node_id, &db_ip); err != nil {
			return err
		}

		(*nodes_list)[db_ip] = append((*nodes_list)[db_ip], node_id)
		if dirty {
			dirty_keys = append(dirty_keys, db_ip)
		}
	}
	if dirty {
		o.NodeInfo.dirty_keys = dirty_keys
	}
	return nil
}

func (o *NKNOVH) saveANStatus() error {
	rows, err := o.sql.stmt["main"]["selectAllANLast"].Query()
	if err != nil { 
		return err
	}
	defer rows.Close()

	var (
		syncState string
		uptime sql.NullInt64
		proposalSubmitted sql.NullInt64
		relayMessageCount sql.NullInt64
		//version sql.Nullstring

		all_uptime uint64
		relays uint64
		average_uptime uint64
		average_relays uint64
		relays_per_hour float64
		persist_nodes_count int
		proposalSubmittedAll int
		nodes_count int

		uptime_ok bool
	)
	for rows.Next() {
		if err := rows.Scan(&syncState, &uptime, &proposalSubmitted, &relayMessageCount); err != nil {
			return err
		}
		nodes_count++
		if syncState == "PERSIST_FINISHED" {
			persist_nodes_count++
		}
		if uptime_ok = uptime.Valid; uptime_ok {
			all_uptime += uint64(uptime.Int64)
		}
		if ok := relayMessageCount.Valid; ok {
			relays += uint64(relayMessageCount.Int64)
			if uptime_ok == true {
				if uptime.Int64 != 0 {
					relays_per_hour += float64(relayMessageCount.Int64)/float64(uptime.Int64)*3600
				} else {
					relays_per_hour += float64(relayMessageCount.Int64)
				}
			}
		}
		if ok := proposalSubmitted.Valid; ok {
			proposalSubmittedAll += int(proposalSubmitted.Int64)
		}
	}
	average_relays = uint64(relays_per_hour/float64(nodes_count))
	average_uptime = uint64(float64(all_uptime)/float64(nodes_count))


	//Get the latest height
	var last_height uint64
	var last_timestamp uint64
	var average_blockTime float64
	var average_blocksPerDay float64
	const FirstHeightTS uint64 = 1561814790

	row := o.sql.stmt["main"]["selectLastHeightANLast"].QueryRow()
	err = row.Scan(&last_height)
	switch {
	case err == sql.ErrNoRows:
		o.log.Syslog("Cannot fetch info from db (last_height), 0 row is found", "main")
	break
	case err != nil:
		o.log.Syslog("Can't execute row.Scan(): "+err.Error(), "sql")
		return err
	default:
		o.NodeInfo.ANLastMux.Lock()
		last_timestamp = uint64(time.Now().Unix())
		if last_height > 0 {
			average_blockTime = float64(last_timestamp-FirstHeightTS)/float64(last_height)
			average_blocksPerDay = 86400/average_blockTime
			o.NodeInfo.ANLast["Timestamp"] = float64(last_timestamp)
			o.NodeInfo.ANLast["Height"] = float64(last_height)
			o.NodeInfo.ANLast["averageBlockTime"] = average_blockTime
		}
		o.NodeInfo.ANLastMux.Unlock()
	}

	if _, err := o.sql.stmt["main"]["insertANStats"].Exec(&relays, &average_uptime, &average_relays, &relays_per_hour, &proposalSubmittedAll, &persist_nodes_count, &nodes_count, &last_height, &last_timestamp, &average_blockTime, &average_blocksPerDay); err != nil {
		o.log.Syslog("Stmt insertANStats has returned an error: ("+err.Error()+")", "sql")
			return err
	}
	return nil
}

func (o *NKNOVH) getInfo(wg *sync.WaitGroup, obj *JsonRPCConf, inside_method string, threads *chan struct{}, params ...interface{}) error {
	defer wg.Done()
	var data NodeSt
	switch inside_method {
		case "UpdateNode", "UpdateNodeAN":
			obj.UnmarshalData = &data.State
		default:
			obj.UnmarshalData = &data.Neighbor
	}	

	_, err := o.jrpc_get(obj)
	if err != nil {

		// Handling UpdateNode variations
		switch inside_method {
			case "UpdateNode":
				if len(params) == 2 {
					return err
				}
				if len(params) > 0 {
					reflect.ValueOf(o).MethodByName(inside_method + "Fail").Call([]reflect.Value{reflect.ValueOf([]byte("")),reflect.ValueOf(params[0])})
				} else {
					reflect.ValueOf(o).MethodByName(inside_method + "Fail").Call([]reflect.Value{reflect.ValueOf([]byte(""))})
				}
				<-*threads
				return err
			case "AddNeighborAN":
				o.log.Syslog(err.Error(), "neighbors")
				<-*threads
				return err
			default:
				<-*threads
				return err
		}
	}
	switch method := obj.Method; method {
	case "getnodestate":
		//Handling UpdateNode variations
		if inside_method == "UpdateNode" {
			if err == nil {
				if data.State.Error != nil {
					o.UpdateNodeErr(&data.State, params[0])
					// check for no-goroutine recursive function
					if len(params) < 2 {
						<-*threads
					}
					return nil
				}
			}
		}
		if len(params) > 0 {
			reflect.ValueOf(o).MethodByName(inside_method).Call([]reflect.Value{reflect.ValueOf(&data.State), reflect.ValueOf(params[0])})
		} else {
			reflect.ValueOf(o).MethodByName(inside_method).Call([]reflect.Value{reflect.ValueOf(&data.State)})
		}
	case "getneighbor":
		reflect.ValueOf(o).MethodByName(inside_method).Call([]reflect.Value{reflect.ValueOf(&data.Neighbor)})
	}
	
	// check for no-goroutine recursive function
	if len(params) != 2 {
		<-*threads
	}
	return nil
}

func (o *NKNOVH) rmNodesByFcnt(over_failcnt int64, firsttime_failed uint8) error {
	if _, err := o.sql.stmt["main"]["rmNodesByFcnt"].Exec(&over_failcnt, &firsttime_failed); err != nil {
		o.log.Syslog("Stmt rmNodesByFcnt has returned an error: ("+err.Error()+")", "sql")
		return err
	}
	return nil
}

func (o *NKNOVH) respErrorHandling(rpcerr *RPCErrorState) *RPCError {
	obj := new(RPCError)
	obj.Code = rpcerr.Code
	if rpcerr.Message != "" {
		obj.Description = rpcerr.Message
	}
	switch obj.Code {
		case -32601:
			obj.Status = "Method not found"
			obj.Description = "The called method was not found on the server"
		case -41001:
			obj.Status = "SESSION EXPIRED"
		case -41002:
			obj.Status = "SERVICE CEILING"	
		case -41003:
			obj.Status = "ILLEGAL DATAFORMAT"
		case -42001:
			obj.Status = "INVALID METHOD"	
		case -42002:
			obj.Status = "INVALID PARAMS"
		case -42003:
			obj.Status = "VERIFY TOKEN ERROR"
		case -43001:
			obj.Status = "INVALID TRANSACTION"
		case -43002:
			obj.Status = "INVALID ASSET"	
		case -43003:
			obj.Status = "INVALID BLOCK"
		case -43004:
			obj.Status = "INVALID HASH"
		case -43005:
			obj.Status = "INVALID VERSION"
		case -44001:
			obj.Status = "UNKNOWN TRANSACTION"
		case -44002:
			obj.Status = "UNKNOWN ASSET"
		case -44003:
			obj.Status = "UNKNOWN BLOCK"
		case -44004:
			obj.Status = "UNKNOWN HASH"
		case -45001:
			obj.Status = "INTERNAL ERROR"
		case -45022:
			obj.Status = "GENERATING ID"
			if x := rpcerr.WalletAddress; x != "" {
				obj.WalletAddress = x
			}
			if x := rpcerr.PublicKey; x != "" {
				obj.PublicKey = x
			}
		case -45024:
			obj.Status = "PRUNING DB"
		case -47001:
			obj.Status = "SMARTCODE EXEC ERROR"
		default:
		scode := strconv.Itoa(obj.Code)
		obj.Status = "UNKNOWN RESPONSE CODE"
		obj.Description = "Unknown response code [" + scode + "]"
	}
	return obj
}

func (o *NKNOVH) UpdateNodeErr(resp *NodeState, params interface{}) {
	dbnode := params.(*DBNode)
	var failcnt int64
	var ftf uint8
	var last_id uint64

	resperr := o.respErrorHandling(resp.Error)

	for _, node_id := range dbnode.Ids {

		row := o.sql.stmt["main"]["selectNodeLastIdByNodeId"].QueryRow(&node_id)
		err := row.Scan(&last_id, &failcnt, &ftf)
		switch {
		case err == sql.ErrNoRows:
			if _, err1 := o.sql.stmt["main"]["insertNodeLast"].Exec(&node_id, "", 0, 0, 0, -1, 0, &resperr.Status, 0, "", 1, &ftf); err1 != nil {
				o.log.Syslog("Stmt insertNodeLast has returned an error: (" + err1.Error() + ")", "sql")
			}
		break
		case err != nil:
			o.log.Syslog("Can't execute row.Scan(): " + err.Error(), "sql")
		break
		default:
			if _, err1 := o.sql.stmt["main"]["updateNodeLastById"].Exec("", 0, 0, 0, -1, 0, &resperr.Status, 0, "", &failcnt, &ftf, &last_id); err1 != nil {
				o.log.Syslog("Stmt updateNodeLastById has returned an error: (" + err1.Error() + ")", "sql")
			}
		}
	}
	return
}

func (o *NKNOVH) UpdateNodeAN(node *NodeState) error {

	if b := o.Validator.IsNodeStateValid(node); !b {
		s := "Invalid NodeState (from UpdateNodeAN)"
		o.log.Syslog(s, "main")
		return errors.New(s)
	}

	var ip string
	re_ip := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	if tmp := re_ip.FindString(node.Result.Addr); tmp != "" {
		ip = tmp
	} else {
		return nil
	}

	if _, err := o.sql.stmt["main"]["updateNodeByIpAN"].Exec(node.Result.SyncState, node.Result.Uptime, node.Result.ProposalSubmitted, node.Result.RelayMessageCount, node.Result.Height, node.Result.Version, node.Result.Currtimestamp, ip); err != nil {
		o.log.Syslog("Can't execute updateNodeByIp: "+err.Error(), "sql")
		return err
	}
	return nil
}

func (o *NKNOVH) dbIpsToArray() error {
	//clear map
	o.NodeInfo.ANArray = map[int]map[int][]int{}
	for i := 0; i < 256; i++ {
		o.NodeInfo.ANArray[i] = map[int][]int{}
		for n := 0; n < 256; n++ {
			o.NodeInfo.ANArray[i][n] = make([]int, 0)
		}
	}

	for i := range o.NodeInfo.ips {
		if _, err := o.searchIP(o.NodeInfo.ips[i]); err != nil {
			o.log.Syslog("Get an error on searchIP func: " + err.Error(), "main")
			continue
		}
	}
	for i := 0; i < 256; i++ {
		for n := 0; n < 256; n++ {
			sort.Ints(o.NodeInfo.ANArray[i][n])
		}
	}
	return nil
}

func (o *NKNOVH) updateAN() error {
	var wg sync.WaitGroup
	o.log.Syslog("[NeighborPoll] Starting get neighbors", "main")

	//Init Mainnet
	for i, _ := range o.conf.SeedList {
		o.NodeInfo.ips = append(o.NodeInfo.ips, o.conf.SeedList[i])
	}

	for x := 0; x < 2; x++ {
		for i := range o.NodeInfo.ips {
			r := &JsonRPCConf{Ip:o.NodeInfo.ips[i], Method:"getneighbor", Params: &json.RawMessage{'{','}'}, Client: o.http.NeighborClient}
			wg.Add(1)
			o.threads.Neighbors <- struct{}{}
			go o.getInfo(&wg, r, "AddNeighborAN", &o.threads.Neighbors)
		}
		 wg.Wait()
		 o.getANFromDB()
		 o.dbIpsToArray()
	}
	o.log.Syslog("[NeighborPoll] All neighbors getted", "main")
	
	for i := range o.NodeInfo.ips {
		r := &JsonRPCConf{Ip:o.NodeInfo.ips[i], Method:"getnodestate", Params: &json.RawMessage{'{','}'}, Client: o.http.NeighborClient}
		wg.Add(1)
		o.threads.Neighbors <- struct{}{}
		go o.getInfo(&wg, r, "UpdateNodeAN", &o.threads.Neighbors)
	}
	wg.Wait()
	
	o.NodeInfo.ips = make([]string, 0)

	o.log.Syslog("[NeighborPoll] Stats of all nodes saved", "main")
	return nil
}

func (o *NKNOVH) rmOldHistory(node_id uint64) error {
	var count_entries int
	row := o.sql.stmt["main"]["countNodeHistory"].QueryRow(&node_id);
	err := row.Scan(&count_entries)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		o.log.Syslog("Can't execute row.Scan(): "+err.Error(), "sql")
		return err
	default:
		diff := count_entries - o.conf.MainPoll.EntriesPerNode
		if diff <= 0 {
			return nil
		}
		if _, err := o.sql.stmt["main"]["rmOldHistory"].Exec(&node_id, &diff); err != nil {
			o.log.Syslog("Stmt rmOldHistory has returned an error: ("+err.Error()+")", "sql")
			return err
		}
	}

	return nil
}


func (o *NKNOVH) searchIP(ip string) (int, error) {
	
	var ip_split []string
	var i_first int
	var i_second int
	//var its_found bool
	var err error
	var ip2int int
	ip_split = strings.Split(ip, ".")
		if len(ip_split) < 1 {
			o.log.Syslog("IP is not splitted", "main")
			return 1, errors.New("IP is not splitted")
		}
		if i_first, err = strconv.Atoi(ip_split[0]); err != nil {
			o.log.Syslog("Cannot ParseInt ip_split[0]", "main")
			return 1, errors.New("Cannot ParseInt ip_split[0]")
		}
		if i_second, err = strconv.Atoi(ip_split[1]); err != nil {
			o.log.Syslog("Cannot ParseInt ip_split[1]", "main")
			return 1, errors.New("Cannot ParseInt ip_split[1]")
		}
	ip2int = IP4toInt(ip)
	
	o.NodeInfo.ANArrayMux[i_first].Lock()
	
	blocklen := len(o.NodeInfo.ANArray[i_first][i_second])
	
	i := sort.SearchInts(o.NodeInfo.ANArray[i_first][i_second], ip2int);
	if i < blocklen && o.NodeInfo.ANArray[i_first][i_second][i] == ip2int {
		o.NodeInfo.ANArrayMux[i_first].Unlock()
		return 0, nil
	}
	
	o.NodeInfo.ANArray[i_first][i_second] = append(o.NodeInfo.ANArray[i_first][i_second], ip2int)
	sort.Ints(o.NodeInfo.ANArray[i_first][i_second])
	o.NodeInfo.ANArrayMux[i_first].Unlock()
	return 2, nil
}



func (o *NKNOVH) AddNeighborAN(nei *NodeNeighbor) error {
	
	if b := o.Validator.IsNodeNeighborValid(nei); !b {
		s := "NodeNeighbor is not valid"
		o.log.Syslog(s, "main")
		return errors.New(s)
	}

	re_ip := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	var ip string
	var l int = len(nei.Result)
	var search int
	var err error
	for n := 0; n < l; n++ {
		if tmp := re_ip.FindString(nei.Result[n].Addr); tmp != "" {
			ip = tmp
		} else {
			o.log.Syslog("IP is nullstring", "main")
			continue
		}
		
		if search, err = o.searchIP(ip); err != nil {
			o.log.Syslog(err.Error(), "main")
			continue
		}
		if search != 2 {
			continue
		}
		
		if _, err1 := o.sql.stmt["main"]["insertAN"].Exec(ip,&nei.Result[n].Addr,&nei.Result[n].ID,&nei.Result[n].PublicKey,&nei.Result[n].SyncState,&nei.Result[n].Height); err1 != nil {
			o.log.Syslog("Stmt insertAN has returned an error: ("+err1.Error()+")", "sql")
			continue
		}
	}
	return nil
}

