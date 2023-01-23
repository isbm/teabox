package teaconditions

// TeaCondition interface for each specific condition
type TeaCondition interface {
	IsSatisfied() bool
	GetInfoMessage() string
}
