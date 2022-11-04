package teaboxlib

import "github.com/gdamore/tcell/v2"

// EGA colors :-)
var EGAColorBlack = tcell.NewHexColor(0x000000)
var EGAColorBlue = tcell.NewHexColor(0x000088)
var EGAColorGreen = tcell.NewHexColor(0x008800)
var EGAColorCyan = tcell.NewHexColor(0x008888)
var EGAColorRed = tcell.NewHexColor(0x880000)
var EGAColorMagenta = tcell.NewHexColor(0x880088)
var EGAColorYellow = tcell.NewHexColor(0xAA5500)
var EGAColorBrown = EGAColorYellow
var EGAColorWhite = tcell.NewHexColor(0xAAAAAA) // In EGA the white color darkens YOU!
var EGAColorLightGray = EGAColorWhite
var EGAColorDarkGray = tcell.NewHexColor(0x555555)
var EGAColorBrightBlack = EGAColorDarkGray // In EGA we have 256 shades of black!
var EGAColorBrightBlue = tcell.NewHexColor(0x5555FF)
var EGAColorBrightGreen = tcell.NewHexColor(0x33DD33)
var EGAColorBrightCyan = tcell.NewHexColor(0x55FFFF)
var EGAColorBrightRed = tcell.NewHexColor(0xFF5555)
var EGAColorBrightMagenta = tcell.NewHexColor(0xFF55FF)
var EGAColorBrightYellow = tcell.NewHexColor(0xFFFF55)
var EGAColorBrightWhite = tcell.NewHexColor(0xFFFFFF)
