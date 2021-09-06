package nknovh_engine

import (
		"regexp"
)

func (o *NKNOVH) isNodeStateValid(s *NodeState) bool {
	if s.Error != nil {
		return false
	}
	var b bool
	re_addr := regexp.MustCompile(`^tcp://(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:([0-9]){0,5}$`)
	if b = re_addr.MatchString(s.Result.Addr); !b {
		return false
	}

	if len(s.Result.ID) != 64 || len(s.Result.PublicKey) != 64 {
		return false
	}
	re_id := regexp.MustCompile(`^([A-Za-z0-9]{64})$`)
	if b = re_id.MatchString(s.Result.ID); !b {
		return false
	}
	if b = re_id.MatchString(s.Result.PublicKey); !b {
		return false
	}
	re_state := regexp.MustCompile(`^WAIT_FOR_SYNCING|SYNC_STARTED|SYNC_FINISHED|PERSIST_FINISHED$`)
	if b = re_state.MatchString(s.Result.SyncState); !b {
		return false
	}
	re_domain := regexp.MustCompile(`^[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}-[0-9]{1,3}\.ipv4\.nknlabs\.io$`)
	if b = re_domain.MatchString(s.Result.Tlsjsonrpcdomain); !b {
		return false
	}
	if b = re_domain.MatchString(s.Result.Tlswebsocketdomain); !b {
		return false
	}
	if len(s.Result.Version) > 64 {
		return false
	}
	re_ver := regexp.MustCompile(`^([0-9\.A-Za-z\-]*)$`)
	if b = re_ver.MatchString(s.Result.Version); !b {
		return false
	}
	return true
}

func (o *NKNOVH) isNodeNeighborValid(s *NodeNeighbor) bool {


}