package teaboxlib

import "github.com/gdamore/tcell/v2"

// EGA colors :-)
var EGAColorBlack = tcell.NewHexColor(0x000000).TrueColor()
var EGAColorBlue = tcell.NewHexColor(0x000088).TrueColor()
var EGAColorGreen = tcell.NewHexColor(0x008800).TrueColor()
var EGAColorCyan = tcell.NewHexColor(0x008888).TrueColor()
var EGAColorRed = tcell.NewHexColor(0x880000).TrueColor()
var EGAColorMagenta = tcell.NewHexColor(0x880088).TrueColor()
var EGAColorYellow = tcell.NewHexColor(0xAA5500).TrueColor()
var EGAColorBrown = EGAColorYellow.TrueColor()
var EGAColorWhite = tcell.NewHexColor(0xAAAAAA).TrueColor() // In EGA the white color darkens YOU!
var EGAColorLightGray = tcell.ColorLightGrey.TrueColor()
var EGAColorDarkGray = tcell.NewHexColor(0x555555).TrueColor()
var EGAColorBrightBlack = EGAColorDarkGray.TrueColor() // In EGA we have 256 shades of black!
var EGAColorBrightBlue = tcell.NewHexColor(0x5555FF).TrueColor()
var EGAColorBrightGreen = tcell.NewHexColor(0x33DD33).TrueColor()
var EGAColorBrightCyan = tcell.NewHexColor(0x55FFFF).TrueColor()
var EGAColorBrightRed = tcell.NewHexColor(0xFF5555).TrueColor()
var EGAColorBrightMagenta = tcell.NewHexColor(0xFF55FF).TrueColor()
var EGAColorBrightYellow = tcell.NewHexColor(0xFFFF55).TrueColor()
var EGAColorBrightWhite = tcell.NewHexColor(0xFFFFFF).TrueColor()
