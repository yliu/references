package main

import (
	"log"
	"reflect"
	"regexp"
	"strings"
	"unsafe"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/mattn/go-runewidth"
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

type refList struct {
	List        *widgets.List
	CurrentLine int
	Dx          int
	Dy          int
}

type functionBox struct {
	Text      string
	Prefix    string
	Suffix    string
	Paragraph *widgets.Paragraph
}
type helpBox struct {
	Text      string
	Paragraph *widgets.Paragraph
}

func createRefList(txt string) *refList {
	var refList refList
	refList.List = widgets.NewList()
	width, height := ui.TerminalDimensions()
	refList.Dx = -1
	refList.Dy = -1
	refList.List.SetRect(refList.Dx, refList.Dy, width+1, height+1)
	refList.List.WrapText = false
	refList.List.TextStyle = ui.NewStyle(ui.ColorClear)
	refList.List.SelectedRowStyle = ui.NewStyle(ui.ColorWhite, ui.ColorBlue)
	refList.List.Border = true
	refList.UpdateContent(txt)
	return &refList
}
func (reflist *refList) UpdateContent(txt string) {
	rows := strings.Split(txt, "\n")
	reflist.UpdateRows(rows)
}
func (reflist *refList) UpdateRows(rows []string) (enlarge int) {
	if len(reflist.List.Rows) == 1 {
		row0 := strings.TrimSpace(reflist.List.Rows[0])
		if row0 == simpleHelp || row0 == "" {
			enlarge = len(rows)
		}
	} else if len(reflist.List.Rows) < len(rows) {
		enlarge = len(rows) - len(reflist.List.Rows)
	}
	reflist.List.Rows = rows
	if reflist.List.SelectedRow >= len(reflist.List.Rows) {
		reflist.List.SelectedRow = len(reflist.List.Rows) - 1
	}
	width := reflist.List.GetRect().Dx()
	for i, row := range reflist.List.Rows {
		re := regexp.MustCompile(`\[(.*)\]\(.*\)`)
		rowTrimSpace := strings.TrimRight(row, " ")
		rowReal := re.ReplaceAllString(rowTrimSpace, "$1")
		rowLen := runewidth.StringWidth(rowReal)
		if rowLen+2 < width {
			reflist.List.Rows[i] = rowTrimSpace + strings.Repeat(" ", width-rowLen-2)
		}
	}
	ui.Render(reflist.List)
	return enlarge
}
func (reflist *refList) Render() {
	ui.Render(reflist.List)
}
func (reflist *refList) ReRect(x, y int) {
	reflist.List.SetRect(reflist.Dx, reflist.Dy, x+1, y+1)
	reflist.UpdateRows(reflist.List.Rows)
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollDown() {
	reflist.List.ScrollDown()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollUp() {
	reflist.List.ScrollUp()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollPageDown() {
	reflist.List.ScrollPageDown()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollPageUp() {
	reflist.List.ScrollPageUp()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollHalfPageDown() {
	reflist.List.ScrollHalfPageDown()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollHalfPageUp() {
	reflist.List.ScrollHalfPageUp()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollBottom() {
	reflist.List.ScrollBottom()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollTop() {
	reflist.List.ScrollTop()
	ui.Render(reflist.List)
}
func (reflist *refList) ScrollAmount(amount int) {
	reflist.List.ScrollAmount(amount)
	ui.Render(reflist.List)
}
func (reflist *refList) SelectLastRow(last int) {
	reflist.List.SelectedRow = len(reflist.List.Rows) - last
	ui.Render(reflist.List)
}
func (reflist *refList) MouseClick(x, y int) {
	if x > reflist.List.Min.X && x < reflist.List.Max.X-1 && y > reflist.List.Min.Y && y < reflist.List.Max.Y-1 {
		reflist.List.SelectedRow = hackGetTopRow(reflist.List) + y + reflist.List.Min.X + 1
	}
	ui.Render(reflist.List)
}

func hackGetTopRow(l *widgets.List) int {
	pointerVal := reflect.ValueOf(l)
	val := reflect.Indirect(pointerVal)
	member := val.FieldByName("topRow")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*int)(ptrToY)
	return *realPtrToY
}

func createFuncBox() *functionBox {
	var funcBox functionBox
	funcBox.Prefix = "> "
	funcBox.Suffix = "[ ](bg:blue)[ ](bg:clear)"

	t := widgets.NewParagraph()

	t.Title = " Input Function Name "
	t.Text = funcBox.Prefix + funcBox.Text + funcBox.Suffix
	funcBox.Paragraph = t
	funcBox.Render()
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
	ui.Render(funcbox.Paragraph)
}
func (funcbox *functionBox) Update(txt string) {
	if len(txt) > 50 {
		return
	}
	funcbox.Text = txt
	funcbox.Paragraph.Text = funcbox.Prefix + funcbox.Text + funcbox.Suffix
	ui.Render(funcbox.Paragraph)
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

func newUI() {
	var data functionReferences
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	reflist := createRefList("")
	funcbox := createFuncBox()
	helpbox := createHelpBox()

	mode := "function"
	newFunc := ""

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents

		if mode == "function" {
			match, err := regexp.MatchString(`^[A-Za-z0-9_]$`, e.ID)
			if err != nil {
				continue
			}

			switch {
			case e.ID == "<Enter>":
				newFunc = funcbox.Text
				data.AddFunction(newFunc)
				enlarge := reflist.UpdateRows(data.List())
				if enlarge > 0 {
					reflist.SelectLastRow(enlarge)
				}
				fallthrough
			case e.ID == "<Escape>" || e.ID == "<C-c>":
				reflist.Render()
				mode = ""
			case e.ID == "<Backspace>":
				if len(funcbox.Text) > 0 {
					funcbox.Update(funcbox.Text[:len(funcbox.Text)-1])
				}
			case e.ID == "<C-w>":
				funcbox.Update("")
			case match:
				funcbox.Update(funcbox.Text + e.ID)
			}
			continue
		} else if mode == "help" {
			if e.Type != ui.MouseEvent {
				mode = ""
				reflist.Render()
			}
		}

		switch e.ID {
		case "q", "<C-c>", "Q":
			return
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
		case "<Home>":
			reflist.ScrollTop()
		case "G", "<End>":
			reflist.ScrollBottom()
		case "h", "H":
			helpbox.Render()
			mode = "help"
		case "<C-l>":
			data = functionReferences{}
			reflist.UpdateRows(data.List())
			fallthrough
		case "n":
			funcbox.Update("")
			mode = "function"
		case "D":
			data.RemoveFunctionByIndex(reflist.List.SelectedRow)
			reflist.UpdateRows(data.List())
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			reflist.ReRect(x, y)
		case "<Enter>":
			data.ReferenceByIndex(reflist.List.SelectedRow)
			reflist.UpdateRows(data.List())
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
