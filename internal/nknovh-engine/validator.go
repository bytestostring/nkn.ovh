package nknovh_engine

import (
		"regexp"
)

type Validator struct {
	Expr map[string]*regexp.Regexp
}

func buildValidator() *Validator {
	v := new(Validator)
	v.Expr = map[string]*regexp.Regexp{}
	v.Expr["Addr"] = regexp.MustCompile(`^tcp://(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:([0-9]){0,5}$`)
	v.Expr["Ipv4"] = regexp.MustCompile(`^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)
	v.Expr["Id"] = regexp.MustCompile(`^([A-Za-z0-9]{64})$`)
	v.Expr["PublicKey"] = regexp.MustCompile(`^([A-Za-z0-9]{64})$`)
	v.Expr["SyncState"] = regexp.MustCompile(`^WAIT_FOR_SYNCING|SYNC_STARTED|SYNC_FINISHED|PERSIST_FINISHED$`)
	v.Expr["Tlsjsonrpcdomain"] = regexp.MustCompile(`^[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}\.ipv4\.(?:nknlabs|staticdns([0-9][0-9]{0,2}|1000))\.io$`)
	v.Expr["Tlswebsocketdomain"] = regexp.MustCompile(`^[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}\.ipv4\.(?:nknlabs|staticdns([0-9][0-9]{0,2}|1000))\.io$`)
	v.Expr["Version"] = regexp.MustCompile(`^([0-9\.A-Za-z\-]*)$`)
	return v
}

func (v *Validator) IsNodeStateValid(s *NodeState) bool {
	if s.Error != nil {
		return false
	}
	var b bool
	if b = v.Expr["Addr"].MatchString(s.Result.Addr); !b {
		return false
	}
	if len(s.Result.ID) != 64 || len(s.Result.PublicKey) != 64 {
		return false
	}
	if b = v.Expr["Id"].MatchString(s.Result.ID); !b {
		return false
	}
	if b = v.Expr["PublicKey"].MatchString(s.Result.PublicKey); !b {
		return false
	}
	if b = v.Expr["SyncState"].MatchString(s.Result.SyncState); !b {
		return false
	}
	if b = v.Expr["Tlsjsonrpcdomain"].MatchString(s.Result.Tlsjsonrpcdomain); !b {
		return false
	}
	if b = v.Expr["Tlswebsocketdomain"].MatchString(s.Result.Tlswebsocketdomain); !b {
		return false
	}
	if len(s.Result.Version) > 64 {
		return false
	}
	if b = v.Expr["Version"].MatchString(s.Result.Version); !b {
		return false
	}
	return true
}

func (v *Validator) IsIPv4Valid(s string) bool {
	if b := v.Expr["Ipv4"].MatchString(s); !b {
		return false
	}
	return true
}


func (v *Validator) IsNodeNeighborValid(s *NodeNeighbor) bool {
	if s.Error != nil {
		return false
	}
	var b bool
	l := len(s.Result)
	for i := 0; i < l; i++ {
		if b = v.Expr["Addr"].MatchString(s.Result[i].Addr); !b {
			return false
		}
		if len(s.Result[i].ID) != 64 || len(s.Result[i].PublicKey) != 64 {
			return false
		}
		if b = v.Expr["Id"].MatchString(s.Result[i].ID); !b {
			return false
		}
		if b = v.Expr["PublicKey"].MatchString(s.Result[i].PublicKey); !b {
			return false
		}
		if b = v.Expr["SyncState"].MatchString(s.Result[i].SyncState); !b {
			return false
		}
		if b = v.Expr["Tlsjsonrpcdomain"].MatchString(s.Result[i].Tlsjsonrpcdomain); !b {
			return false
		}
		if b = v.Expr["Tlswebsocketdomain"].MatchString(s.Result[i].Tlswebsocketdomain); !b {
			return false
		}
	}
	return true
}
