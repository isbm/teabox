package teaboxlib

import (
	"fmt"
	"sort"
)

// Interface
type TeaConfComponent interface {
	SetGroup(id string)
	GetGroup() string
	GetType() string
	GetTitle() string
	SetTitle(title string)
	GetId() string
	GetCommandName() string
	GetCommands() []map[string]interface{}
	Len() int
	Add(mod TeaConfComponent) *TeaConfBaseEntity
	GetChildren() []TeaConfComponent
}

// Entity
type TeaConfBaseEntity struct {
	title    string
	etype    string
	id       string // command or group or module name
	children []TeaConfComponent
}

func (tcb *TeaConfBaseEntity) getChildrenContainer() []TeaConfComponent {
	if tcb.children == nil {
		tcb.children = []TeaConfComponent{}
	}
	return tcb.children
}

func (tcb *TeaConfBaseEntity) Add(mod TeaConfComponent) *TeaConfBaseEntity {
	tcb.children = append(tcb.getChildrenContainer(), mod)
	sort.Slice(tcb.children, func(i, j int) bool {
		return tcb.children[i].GetTitle() < tcb.children[j].GetTitle()
	})
	return tcb
}

func (tcb *TeaConfBaseEntity) GetChildren() []TeaConfComponent {
	return tcb.children
}

func (tcb *TeaConfBaseEntity) Len() int {
	return len(tcb.getChildrenContainer())
}

func (tcb *TeaConfBaseEntity) GetTitle() string {
	return tcb.title
}

func (tcb *TeaConfBaseEntity) SetTitle(title string) {
	tcb.title = title
}

func (tcb *TeaConfBaseEntity) GetId() string {
	return tcb.id
}

func (tcb *TeaConfBaseEntity) GetType() string {
	return tcb.etype
}

func (tcb *TeaConfBaseEntity) GetGroup() string {
	return tcb.id
}

func (tcb *TeaConfBaseEntity) SetGroup(id string) {
	tcb.id = id
}

func (tcb *TeaConfBaseEntity) GetCommandName() string {
	// This cannot be a group, so the ID in this case is a command
	return tcb.id
}

// A single module in the menu root, opens immediately a form to run
type TeaConfModule struct {
	conditions []map[string]string
	commands   []map[string]interface{} // Further decomposition pending

	TeaConfBaseEntity
}

func NewTeaConfModule(title string) *TeaConfModule {
	tcm := new(TeaConfModule)
	tcm.SetTitle(title)
	tcm.etype = "module"

	return tcm
}

func (tcf *TeaConfModule) SetCondition(cond interface{}) *TeaConfModule {
	if cond != nil {
		tcf.conditions = []map[string]string{}
		condset, ok := cond.([]interface{})
		if !ok {
			panic("Wrong configuration of the module: " + tcf.title)
		}

		for _, icnd := range condset {
			cnd := map[string]string{}
			imcnd := icnd.(map[interface{}]interface{})
			for k, v := range imcnd {
				ks, ok := k.(string)
				if !ok {
					panic(fmt.Sprintf("Wrong configuration of the module %s: key %v is not string", tcf.title, k))
				}
				vs, ok := v.(string)
				if !ok {
					panic(fmt.Sprintf("Wrong configuration of the module %s: value %v is not string", tcf.title, v))
				}

				cnd[ks] = vs
			}
			tcf.conditions = append(tcf.conditions, cnd)
		}
	}
	return tcf
}

func (tcf *TeaConfModule) GetCondition() []map[string]string {
	return tcf.conditions
}

func (tcf *TeaConfModule) SetCommands(commands interface{}) *TeaConfModule {
	if commands != nil {
		data, ok := commands.([]map[string]interface{})
		if ok {
			tcf.commands = data
		} // XXX: else pop-up somethig or log etc...
	}
	return tcf
}

func (tcf *TeaConfModule) GetCommands() []map[string]interface{} {
	return tcf.commands
}

// Group of the modules by a topic
type TeaConfGroup struct {
	TeaConfBaseEntity
}

func NewTeaConfGroup(id string) *TeaConfGroup {
	tcg := new(TeaConfGroup)
	tcg.id = id
	tcg.etype = "group"
	tcg.SetTitle(id) // group is also a title in this case
	return tcg
}

func (tcg *TeaConfGroup) GetCommands() []map[string]interface{} {
	return nil
}

// Command
type TeaConfCmd struct {
	TeaConfBaseEntity
}

func NewTeaConfCmd(id, title string) *TeaConfCmd {
	tc := new(TeaConfCmd)
	tc.SetTitle(title)
	tc.id = id
	tc.etype = "command"
	return tc
}

func (tc *TeaConfCmd) GetCommands() []map[string]interface{} {
	return nil
}
