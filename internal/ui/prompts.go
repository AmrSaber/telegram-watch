package ui

import "github.com/fatih/color"

var boldColor = color.New(color.Bold)

var NamePrompt = boldColor.Sprint("Name:")
var HostnamePrompt = boldColor.Sprint("Hostname:")
var TelegramIdPrompt = boldColor.Sprint("Telegram ID:")
