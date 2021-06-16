package nknovh_engine

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Mysql struct {
	db   map[string]*sql.DB
	log  *logger
	stmt map[string]map[string]*sql.Stmt
}

func (o *Mysql) build() {
	//	o.log = log
	o.db = map[string]*sql.DB{}
	o.stmt = map[string]map[string]*sql.Stmt{"main": map[string]*sql.Stmt{}}
}

func (o *Mysql) prepare() error {
	queries := map[string]map[string]string{
		"main": map[string]string{
			"fetchUniqs": "SELECT * FROM uniq",
			"insertAN": "INSERT IGNORE INTO all_nodes(ip,addr,node_id,syncState,height) VALUES(?,?,?,?,?)",
			"selectLastHeightAN": "SELECT height, currtimestamp FROM all_nodes WHERE height = (SELECT MAX(height) FROM all_nodes) LIMIT 1",
			"selectIdByAddrAN": "SELECT id FROM all_nodes WHERE addr = ?",
			"selectIdByIpAN": "SELECT id FROM all_nodes WHERE ip = ?",
			"selectAllIpsAN": "SELECT ip FROM all_nodes",
			"rmNodesByInterval": "DELETE FROM all_nodes WHERE latest_update < NOW() - INTERVAL ? MINUTE",
			"rmNodesByInactive": "DELETE FROM all_nodes WHERE uptime IS NULL",
			"updateNodeByIpAN": "UPDATE all_nodes SET node_id = ?, syncState = ?, uptime = ?, proposalSubmitted = ?, relayMessageCount = ?, height = ?, version = ?, currtimestamp = ?, latest_update = CURRENT_TIMESTAMP() WHERE ip = ?",
			"selectAllAN": "SELECT syncState, uptime, proposalSubmitted, relayMessageCount FROM all_nodes",
			"insertANStats": "INSERT INTO all_nodes_stats(relays, average_uptime, average_relays, relays_per_hour, proposalSubmitted, persist_nodes_count, nodes_count, last_height, last_timestamp, average_blockTime, average_blocksPerDay) VALUES(?,?,?,?,?,?,?,?,?,?,?)",
			"selectAllNodesDirty": "SELECT id, ip FROM nodes WHERE dirty = 1 ORDER BY dirty_fcnt ASC",
			"selectAllNodesNotDirty": "SELECT id, ip FROM nodes WHERE dirty = 0",
			"selectNodeIpById": "SELECT ip FROM nodes WHERE id = ?",
			"insertNodeStats": "INSERT INTO nodes_history(node_id,NID,Currtimestamp,Height,ProposalSubmitted,ProtocolVersion,RelayMessageCount,SyncState,Uptime,Version) VALUES(?,?,?,?,?,?,?,?,?,?)",
			"countNodeHistory": "SELECT count(id) as cnt FROM nodes_history WHERE node_id = ?",
			"rmOldHistory": "DELETE FROM nodes_history WHERE node_id = ? ORDER BY id ASC LIMIT ?",
			"selectNodeHistLastIdByNodeId": "SELECT id FROM nodes_history WHERE node_id = ? ORDER BY id DESC LIMIT 1",
			"updateNodeHistById": "UPDATE nodes_history SET SyncState = ?, latest_update = CURRENT_TIMESTAMP() WHERE id = ?",
			"updateNodeToDirty": "UPDATE nodes SET dirty = 1, dirty_fcnt = dirty_fcnt+1 WHERE id = ?",
			"updateNodeToMain": "UPDATE nodes SET dirty = 0, dirty_fcnt = 0 WHERE id = ? AND dirty = 1",
			"selectWallets": "SELECT id, nkn_wallet, balance FROM wallets",
			"updateWalletBalanceById": "UPDATE wallets SET balance = ? WHERE id = ?",
			"getPriceByName": "SELECT id FROM prices WHERE name = ?",
			"insertPrice": "INSERT INTO prices(name,price) VALUES(?,?)",
			"updatePriceById": "UPDATE prices SET price = ?, last_update = CURRENT_TIMESTAMP() WHERE id = ?",
			"selectDaemonIdByName": "SELECT id FROM daemon WHERE name = ?",
			"updateDaemonById": "UPDATE daemon SET value = ? WHERE id = ?",
			"insertDaemon": "INSERT INTO daemon(name,value) VALUES(?,?)",
			"rmNodesByFcnt": "DELETE FROM nodes WHERE id IN (SELECT node_id FROM nodes_last WHERE failcnt > ? AND firsttime_failed = ?)",
			"selectNodeLastIdByNodeId": "SELECT id,failcnt,firsttime_failed FROM nodes_last WHERE node_id = ?",
			"insertNodeLast": "INSERT INTO nodes_last(node_id,NID,Currtimestamp,Height,ProposalSubmitted,ProtocolVersion,RelayMessageCount,SyncState,Uptime,Version,failcnt,firsttime_failed) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)",
			"updateNodeLastById": "UPDATE nodes_last SET NID = ?, Currtimestamp = ?, Height = ?, ProposalSubmitted = ?, ProtocolVersion = ?, RelayMessageCount = ?, SyncState = ?, Uptime = ?, Version = ?, failcnt = ?, firsttime_failed = ?, latest_update = CURRENT_TIMESTAMP() WHERE id = ?",
		},
	}
	var stmt *sql.Stmt
	var err error
	for key, val := range queries {
		for key2, value := range val {
			stmt, err = o.db[key].Prepare(queries[key][key2])
			if err != nil {
				o.log.Syslog("Can't prepare an query: "+value, "sql")
				return err
			}
			o.stmt[key][key2] = stmt

		}
	}
	return nil
}

func (o *Mysql) createConnect(host string, dbtype string, login string, password string, database string, moc int, mic int, inside string) error {
	var sqlinfo string
	if dbtype == "mysql" {
		sqlinfo = fmt.Sprintf("%s:%s@/%s", login, password, database)
	} else if dbtype == "postgres" {
		sqlinfo = fmt.Sprintf("host=%s user=%s password=%s sslmode=disable", host, login, password)
	}
	db, err := sql.Open(dbtype, sqlinfo)
	if err != nil {
		o.log.Syslog("Cannot create connect to database: "+err.Error(), "sql")
		return err
	}
	err = db.Ping()
	if err != nil {
		o.log.Syslog("Cannot create connect to database: "+err.Error(), "sql")
		return err
	}
	o.log.Syslog("["+inside+"] Connection to DB \""+database+"\" has successfully created", "sql")
	o.db[inside] = db
	db.SetMaxOpenConns(moc)
	db.SetMaxIdleConns(mic)
	return nil
}
