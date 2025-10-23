package teaboxlib

import (
	"github.com/gdamore/tcell/v2"
)

// This is a configuration of the UI, filled with the defaults.
// Use this to patch for own settings/branding

var MAIN_MENU_WIDTH int = 30

// Labels
var LABEL_BACK = "◀ Back"
var LABEL_EXIT = "Exit ▶"
var LABEL_SEP = "─"
var LABEL_MORE = "…"

// Colors theme
var WORKSPACE_BACKGROUND = EGAColorGreen
var WORKSPACE_HEADER = EGAColorBrightGreen
var WORKSPACE_HEADER_TEXT = tcell.ColorBlack

var FORM_BUTTON_BACKGROUND = EGAColorLightGray
var FORM_BUTTON_BACKGROUND_SELECTED = EGAColorBrightWhite
var FORM_BUTTON_TEXT = EGAColorBlack
var FORM_BUTTON_TEXT_SELECTED = EGAColorBlack

var MENU_BORDER_SELECTED = EGAColorBrightWhite
var MENU_BORDER = EGAColorBrightGreen
var MENU_ITEM_SELECTED = EGAColorLightGray
var MENU_ITEM = EGAColorDarkGray

var FORM_BORDER_SELECTED = EGAColorBrightWhite
var FORM_BORDER = EGAColorBrightGreen
var FORM_BACKGROUND_COLOR_FOCUSED = EGAColorBrightGreen
var FORM_FIELD_TEXT = EGAColorBlack
var FORM_FIELD_BACKGROUND = EGAColorLightGray
var FORM_FIELD_BACKGROUND_DARKER = EGAColorDarkGray
var FORM_FIELD_BACKGROUND_FOCUSED = EGAColorBrightWhite

// Logs
var LOG_FILENAME = "/var/log/teabox.log"
