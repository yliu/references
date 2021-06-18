package main

import (
	"log"
	"os"
	"regexp"

	ui "github.com/gizak/termui/v3"
)

type Mode int

const (
	ModeNormal Mode = iota
	ModeFunction
	ModeHelp
)

func newUI(args []string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	reflist := createRefList()
	funcbox := createFuncBox()
	helpbox := createHelpBox()

	mode := ModeFunction
	prevMode := ModeNormal
	newFunc := ""

	if len(args) > 0 {
		for _, funcname := range args {
			reflist.AddFunction(funcname)
		}
		mode = ModeNormal
	}
	if mode == ModeFunction {
		funcbox.Render()
	}

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		if mode == ModeFunction {
			switch e.ID {
			case "<Enter>":
				newFunc = funcbox.Text
				reflist.AddFunction(newFunc)
				mode = ModeNormal
			case "<Escape>", "<C-c>":
				reflist.Render()
				mode = ModeNormal
			case "<Backspace>":
				if len(funcbox.Text) > 0 {
					funcbox.Update(funcbox.Text[:len(funcbox.Text)-1])
				}
			case "<C-w>":
				funcbox.Update("")
			case "<C-x>":
				extendTag = !extendTag
				funcbox.Render()
			case "<C-d>":
				if funcbox.Text == "" {
					ui.Close()
					os.Exit(0)
				}
			case "?":
				helpbox.Render()
				prevMode = mode
				mode = ModeHelp

			}
			match, err := regexp.MatchString(`^[A-Za-z0-9_]$`, e.ID)
			if err != nil {
				continue
			}
			if match {
				funcbox.Update(funcbox.Text + e.ID)
			}
		} else if mode == ModeHelp {
			if e.Type != ui.MouseEvent {
				mode = prevMode
				reflist.Render()
				if mode == ModeFunction {
					funcbox.Render()
				}
			}
		} else {
			switch e.ID {
			case "q", "<C-c>", "Q":
				ui.Close()
				os.Exit(0)
			case "j", "<Down>":
				reflist.ScrollDown()
			case "k", "<Up>":
				reflist.ScrollUp()
			case "<C-d>":
				reflist.ScrollHalfPageDown()
			case "<C-u>":
				reflist.ScrollHalfPageUp()
			case "<C-f>", "<PageDown>":
				reflist.ScrollPageDown()
			case "<C-b>", "<PageUp>":
				reflist.ScrollPageUp()
			case "<C-x>":
				extendTag = !extendTag
			case "<Home>":
				reflist.ScrollTop()
			case "G", "<End>":
				reflist.ScrollBottom()
			case "h", "H", "?":
				helpbox.Render()
				prevMode = mode
				mode = ModeHelp
			case "<C-l>":
				reflist.Data = functionReferences{}
				reflist.UpdateRows()
				fallthrough
			case "n":
				funcbox.Update("")
				mode = ModeFunction
			case "D":
				reflist.Data.RemoveFunctionByIndex(reflist.List.SelectedRow)
				reflist.UpdateRows()
			case "<Resize>":
				x, y := ui.TerminalDimensions()
				reflist.ReRect(x, y)
			case "<Enter>":
				if reflist.Data.Size() == 0 {
					funcbox.Update("")
					mode = ModeFunction
					continue
				}
				reflist.Data.ReferenceByIndex(reflist.List.SelectedRow)
				reflist.UpdateRows()
			case "o", "O":
				reflist.Data.OpenFile(reflist.List.SelectedRow)
				ui.Close()
				if err := ui.Init(); err != nil {
					log.Fatalf("failed to recreate termui: %v", err)
				}
				defer ui.Close()
				reflist.Render()
			case "<MouseLeft>":
				payload := e.Payload.(ui.Mouse)
				x, y := payload.X, payload.Y
				reflist.MouseClick(x, y)
			case "<MouseWheelDown>":
				reflist.ScrollAmount(3)
			case "<MouseWheelUp>":
				reflist.ScrollAmount(-3)
			}
		}
	}
}
