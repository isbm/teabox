package teaconditions

import (
	"fmt"
	"os/user"
	"strings"
)

// Condition on some file exists or is absent
type TeaCondPerm struct {
	currUsr *user.User
	/* Usage:

	- gid:
	  - kvm
	  - adm
	  mesage: you should be in the valid groups

	- uid:
	  - root
	  message: You are not Groot!

	*/
	BaseTeaCondition
}

// NewTeaCondPerm constructor to a file-related conditions
func NewTeaCondPerm(message, clause string, targets []string) (*TeaCondPerm, error) {
	tcf := new(TeaCondPerm)
	tcf.message = message

	switch clause {
	case "gid", "uid":
		tcf.clause = clause
	default:
		return nil, fmt.Errorf("clause should be either 'gid' or 'uid': '%s'", clause)
	}

	for _, target := range targets {
		target = strings.TrimSpace(target)
		if target != "" {
			tcf.targets = append(tcf.targets, target)
		}
	}

	if len(tcf.targets) == 0 {
		return nil, fmt.Errorf("permission conditions require at lease one GID or UID specified")
	}

	var err error
	tcf.currUsr, err = user.Current()

	if err != nil {
		return nil, err
	}

	return tcf, nil
}

// IsSatisfied returns if a specific condition is met for a specified user permissions
func (tcf *TeaCondPerm) IsSatisfied() bool {
	switch tcf.clause {
	case "gid":
		for _, gid := range tcf.targets {
			if tcf.currUsr.Gid == gid {
				return true
			}
		}
	case "uid":
		for _, uid := range tcf.targets {
			if tcf.currUsr.Uid == uid {
				return true
			}
		}
	}

	return false
}
