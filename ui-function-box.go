package main

import (
	"fmt"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type functionBox struct {
	Text      string
	Prefix    string
	Suffix    string
	BaseTitle string
	Paragraph *widgets.Paragraph
}

func createFuncBox() *functionBox {
	var funcBox functionBox
	funcBox.Prefix = "> "
	funcBox.Suffix = "[ ](bg:blue)[ ](bg:clear)"
	funcBox.BaseTitle = " Input Function Name "

	t := widgets.NewParagraph()

	t.Title = funcBox.BaseTitle
	t.Text = funcBox.Prefix + funcBox.Text + funcBox.Suffix
	funcBox.Paragraph = t
	return &funcBox
}

func (funcbox *functionBox) Render() {
	w, h := ui.TerminalDimensions()
	helpLeft := w/2 - 60
	if helpLeft < 1 {
		helpLeft = 1
	}
	helpTop := h/2 - 2
	if helpTop < 1 {
		helpTop = 1
	}
	funcbox.Paragraph.SetRect(helpLeft, helpTop, helpLeft+62, helpTop+3)
	funcbox.Paragraph.Text = funcbox.Prefix + funcbox.Text + funcbox.Suffix
	funcbox.UpdateTitle()
	ui.Render(funcbox.Paragraph)
}

func (funcbox *functionBox) UpdateTitle() {
	tagFlag := ""
	if extendTag {
		tagFlag = "(Fuzzy) "
	}
	funcbox.Paragraph.Title = fmt.Sprintf("%s%s", funcbox.BaseTitle, tagFlag)
}

func (funcbox *functionBox) Update(txt string) {
	if len(txt) > 50 {
		return
	}
	funcbox.Text = txt
	funcbox.Paragraph.Text = funcbox.Prefix + funcbox.Text + funcbox.Suffix
	funcbox.Render()
}
