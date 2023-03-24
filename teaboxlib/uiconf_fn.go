package teaboxlib

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type UiConfig struct {
	tc *TeaConf
}

func NewUiConfig() *UiConfig {
	uic := new(UiConfig)
	return uic
}

// Setup configuration
func (uic *UiConfig) Setup(conf *TeaConf) *UiConfig {
	uic.tc = conf
	return uic.setLabels().setWorkspace().setMenu().setForms().setCommon()
}

func (uic *UiConfig) setCommon() *UiConfig {
	s := uic.tc.GetRootConfig().Find("ui:widgets")
	w, e := s.Int("menu-width", "")
	if e == nil {
		MAIN_MENU_WIDTH = w
	}

	return uic
}

// Set labels
func (uic *UiConfig) setLabels() *UiConfig {
	s := uic.tc.GetRootConfig().Find("ui:widgets")
	for _, k := range []string{"label-back", "label-exit", "label-sep", "label-more"} {
		l := s.String(k, "")
		if l == "" {
			continue
		}

		switch k {
		case "label-back":
			LABEL_BACK = l
		case "label-exit":
			LABEL_EXIT = l
		case "label-sep":
			LABEL_SEP = l
		case "label-more":
			LABEL_MORE = l
		}
	}

	return uic
}

// Set workspace
func (uic *UiConfig) setWorkspace() *UiConfig {
	s := uic.tc.GetRootConfig().Find("ui:colors")
	for _, k := range []string{"background", "header-background", "header-foreground"} {
		if strings.ToLower(s.String(k, "")) != "default" && s.String(k, "") != "" {
			c := uic.getColor(s.Raw()[k])
			if c == nil {
				continue
			}

			switch k {
			case "background":
				WORKSPACE_BACKGROUND = *c
			case "header-background":
				WORKSPACE_HEADER = *c
			case "header-foreground":
				WORKSPACE_HEADER_TEXT = *c
			}
		}
	}
	return uic
}

// Set menu
func (uic *UiConfig) setMenu() *UiConfig {
	s := uic.tc.GetRootConfig().Find("ui:colors")
	for _, k := range []string{"menu-border-selected", "menu-border", "menu-item", "menu-item-selected"} {
		if strings.ToLower(s.String(k, "")) != "default" && s.String(k, "") != "" {
			c := uic.getColor(s.Raw()[k])
			if c == nil {
				continue
			}

			switch k {
			case "menu-border-selected":
				MENU_BORDER_SELECTED = *c
			case "menu-border":
				MENU_BORDER = *c
			case "menu-item":
				MENU_ITEM = *c
			case "menu-item-selected":
				MENU_ITEM_SELECTED = *c
			}
		}
	}

	return uic
}

// Set forms
func (uic *UiConfig) setForms() *UiConfig {
	s := uic.tc.GetRootConfig().Find("ui:colors")
	for _, k := range []string{
		"button-background", "button-background-selected", "button-foreground",
		"button-foreground-selected", "border", "border-selected"} {
		if strings.ToLower(s.String(k, "")) != "default" && s.String(k, "") != "" {
			c := uic.getColor(s.Raw()[k])
			if c == nil {
				continue
			}

			switch k {
			case "button-background":
				FORM_BUTTON_BACKGROUND = *c
			case "button-background-selected":
				FORM_BUTTON_BACKGROUND_SELECTED = *c
			case "button-foreground":
				FORM_BUTTON_TEXT = *c
			case "button-foreground-selected":
				FORM_BUTTON_TEXT_SELECTED = *c
			case "border":
				FORM_BORDER = *c
			case "border-selected":
				FORM_BORDER_SELECTED = *c
			}
		}
	}

	return uic
}

// Get hex color. If "data" is "default" or none, nil returned. Format: "0xrrggbb"
//
// Note:
//
//	  This method also understands decimal numbers, those are coming from YAML
//		 If a user did not used quotes. In that case 0xff will be as 255.
//		 All YAML notations are accepted:
//		 - "0xaabbcc"
//		 - 0xaabbcc
//		 - 0xAABBCC
func (uic *UiConfig) getColor(in interface{}) *tcell.Color {
	var data string
	_, ok := in.(int)
	if ok {
		data = "0x" + strconv.FormatInt(int64(in.(int)), 16)
	} else {
		data = fmt.Sprintf("%v", in)
	}

	if strings.ToLower(data) != "default" && len(data) == 8 {
		r, e := strconv.ParseInt(data[2:4], 16, 64)
		if e != nil {
			return nil
		}

		g, e := strconv.ParseInt(data[4:6], 16, 64)
		if e != nil {
			return nil
		}

		b, e := strconv.ParseInt(data[6:8], 16, 64)
		if e != nil {
			return nil
		}

		c := tcell.NewRGBColor(int32(r), int32(g), int32(b))
		return &c
	}

	return nil
}
