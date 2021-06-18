package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var simpleHelp = "[n] New Function   [h] Show Help  [q] Quit"
var fullHelp = `h, H              Show help
	<Enter>           Expand/Collapse
	n                 New function
	D                 Delete function
	j, <ArrowDown>    Down
	k, <ArrowUp>      Up
	Ctrl+l            Clear
	`

type helpBox struct {
	Text      string
	Paragraph *widgets.Paragraph
}

func createHelpBox() *helpBox {
	var helpBox helpBox

	t := widgets.NewParagraph()
	t.Title = " Help "
	t.Text = fullHelp
	helpBox.Paragraph = t
	return &helpBox
}

func (helpbox *helpBox) Render() {
	width, height := ui.TerminalDimensions()
	helpLeft := width/2 - 60
	if helpLeft < 1 {
		helpLeft = 1
	}
	helpTop := height/2 - 4
	if helpTop < 1 {
		helpTop = 1
	}
	helpbox.Paragraph.SetRect(helpLeft, helpTop, helpLeft+62, helpTop+9)
	ui.Render(helpbox.Paragraph)
}
