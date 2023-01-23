package teaconditions

// Condition on some file exists or is absent
type TeaCondPerm struct {
	/* Usage:

	gid:
	 - kvm
	 - adm

	uid: root

	*/
	BaseTeaCondition
}

// NewTeaCondPerm constructor to a file-related conditions
func NewTeaCondPerm(clause, target string) (*TeaCondPerm, error) {
	tcf := new(TeaCondPerm)

	return tcf, nil
}

// IsSatisfied returns if a specific condition is met for a specified user permissions
func (tcf *TeaCondPerm) IsSatisfied() bool {

	return true
}
