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
	v.Expr["SyncState"] = regexp.MustCompile(`^(WAIT_FOR_SYNCING|SYNC_STARTED|SYNC_FINISHED|PERSIST_FINISHED)$`)
	v.Expr["Tlsjsonrpcdomain"] = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?)*$`)
	v.Expr["Tlswebsocketdomain"] = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?(?:\.[A-Za-z0-9](?:[A-Za-z0-9\-]{0,61}[A-Za-z0-9])?)*$`)
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
	if s.Result.Tlsjsonrpcdomain != "" {
		if b = v.Expr["Tlsjsonrpcdomain"].MatchString(s.Result.Tlsjsonrpcdomain); !b {
			return false
		}
	}
	if s.Result.Tlswebsocketdomain != "" {
		if b = v.Expr["Tlswebsocketdomain"].MatchString(s.Result.Tlswebsocketdomain); !b {
			return false
		}
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
		neighbor := s.Result[i]
		if b = v.Expr["Addr"].MatchString(neighbor.Addr); !b {
			return false
		}
		if len(neighbor.ID) != 64 || len(neighbor.PublicKey) != 64 {
			return false
		}
		if b = v.Expr["Id"].MatchString(neighbor.ID); !b {
			return false
		}
		if b = v.Expr["PublicKey"].MatchString(neighbor.PublicKey); !b {
			return false
		}
		if b = v.Expr["SyncState"].MatchString(neighbor.SyncState); !b {
			return false
		}
		if neighbor.Tlsjsonrpcdomain != "" {
			if b = v.Expr["Tlsjsonrpcdomain"].MatchString(neighbor.Tlsjsonrpcdomain); !b {
				return false
			}
		}
		if neighbor.Tlswebsocketdomain != "" {
			if b = v.Expr["Tlswebsocketdomain"].MatchString(neighbor.Tlswebsocketdomain); !b {
				return false
			}
		}
	}
	return true
}
