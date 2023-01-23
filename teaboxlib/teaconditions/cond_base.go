package teaconditions

// BaseTeaCondition contains a mixin for all required
// common methods of a condition object.
type BaseTeaCondition struct {
	clause  string
	targets []string
	message string
}

// GetMessage returns a message content that needs to be displayed
// *otherwise* (if a condition requirements are not met)
func (btc *BaseTeaCondition) GetInfoMessage() string {
	return btc.message
}
