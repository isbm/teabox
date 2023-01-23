package teaconditions

import (
	"fmt"
	"strings"

	wzlib_utils "github.com/infra-whizz/wzlib/utils"
)

// Condition on some file exists or is absent

type TeaCondFile struct {
	/* Usage:

	absent: /path/to/file
	message: /path/to/file is present

	or:

	present: /path/to/other/file
	message: /path/to/other/file is absent
	*/
	BaseTeaCondition
}

// NewTeaCondFile constructor to a file-related conditions
func NewTeaCondFile(message, clause string, targets []string) (*TeaCondFile, error) {
	tcf := new(TeaCondFile)
	tcf.message = message

	switch clause {
	case "present", "all-present", "absent", "all-absent":
		tcf.clause = clause
	default:
		return nil, fmt.Errorf("clause should be either 'absent' or 'present', not '%s'", clause)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("each target should contain a path to a file or a directory")
	}

	for _, target := range targets {
		if !strings.HasPrefix(target, "/") {
			return nil, fmt.Errorf("target path should be always absolute")
		}
		tcf.targets = append(tcf.targets, target)
	}

	return tcf, nil
}

// IsSatisfied returns if a specific condition is met for a target file (present or absent)
// File can be any. NOTE: permissions are affecting this, if a user has no access to that file,
// then the file is absent to that user in this context.
func (tcf *TeaCondFile) IsSatisfied() bool {
	switch tcf.clause {
	case "present":
		for _, p := range tcf.targets {
			return wzlib_utils.FileExists(p)
		}
	case "all-present":
		for _, p := range tcf.targets {
			if !wzlib_utils.FileExists(p) {
				return false
			}
		}
	case "absent":
		for _, p := range tcf.targets {
			return !wzlib_utils.FileExists(p)
		}
	case "all-absent":
		for _, p := range tcf.targets {
			if wzlib_utils.FileExists(p) {
				return false
			}
		}
	}

	return true
}
