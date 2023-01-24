package teaconditions

import (
	"fmt"
)

// Processor of conditions.
// It will decide if a condition meets its requirements

// TeaConditionsProcessor implement a conditions processor per a module,
// which will check each condition and will determine its state to further
// load a module or not.
type TeaConditionsProcessor struct {
	conditions  []TeaCondition
	buffResult  int // -1 to Unknown, 0 to false, 1 to true
	buffMessage string
}

// NewTeaConditionsProcessor constructor.
func NewTeaConditionsProcessor(rules []map[string][]string) (*TeaConditionsProcessor, error) {
	cond := new(TeaConditionsProcessor)
	cond.conditions = []TeaCondition{}
	cond.buffResult = -1

	if err := cond.loadConditions(rules); err != nil {
		return nil, err
	}

	return cond, nil
}

// Satisfied conditions returns true, otherwise false, if at least one is not satisfied.
func (cnd *TeaConditionsProcessor) Satisfied() bool {
	// Called already?
	if cnd.buffResult > 0 {
		return true
	} else if cnd.buffResult == 0 {
		return false
	}

	// Process and pre-buff
	for _, condition := range cnd.conditions {
		if !condition.IsSatisfied() {
			cnd.buffMessage = condition.GetInfoMessage()
			cnd.buffResult = 0
			return false
		}
	}

	cnd.buffResult = 1
	return true
}

// GetInfoMessage returns a message for a failed condition. If message is requested before
// Satisfied() is called, theen it will call it first to pre-buffer the results.
func (cnd *TeaConditionsProcessor) GetInfoMessage() string {
	if cnd.buffResult < 0 {
		cnd.Satisfied()
	}

	return cnd.buffMessage
}

func (cnd *TeaConditionsProcessor) load(condition map[string][]string) (TeaCondition, error) {
	var message string
	if _, ok := condition["message"]; ok {
		message = condition["message"][0]
		delete(condition, "message")
	} else {
		return nil, fmt.Errorf("condition has no message defined")
	}

	for rule := range condition {
		switch rule {
		case "all-absent", "absent", "all-present", "present":
			return NewTeaCondFile(message, rule, condition[rule])
		case "gid", "uid":
			return NewTeaCondPerm(message, rule, condition[rule])
		default:
			return nil, fmt.Errorf("unknown condition")
		}
	}

	return nil, fmt.Errorf("condition '%v' was not recognised", condition)
}

// Load conditions
func (cnd *TeaConditionsProcessor) loadConditions(rules []map[string][]string) error {
	for _, def := range rules {
		condition, err := cnd.load(def)
		if err != nil {
			return err
		}

		cnd.conditions = append(cnd.conditions, condition)
	}

	return nil
}
