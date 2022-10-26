package teaboxlib

import "github.com/gdamore/tcell/v2"

// This is a configuration of the UI, filled with the defaults.
// Use this to patch for own settings/branding

var MAIN_MENU_WIDTH int = 30

// Labels
var LABEL_BACK = "◀◃◃ Back"
var LABEL_EXIT = "Exit ▹▹▶"
var LABEL_SEP = "─"

// Colors theme
var WORKSPACE_BACKGROUND = tcell.ColorDarkGreen
var WORKSPACE_HEADER = tcell.ColorGreen
var WORKSPACE_HEADER_TEXT = tcell.ColorBlack