package main

import (
	"reflect"
	"regexp"
	"strings"
	"unsafe"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/mattn/go-runewidth"
)

type refList struct {
	List        *widgets.List
	CurrentLine int
	Dx          int
	Dy          int
	Data        functionReferences
}

func createRefList() *refList {
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
	return &refList
}
func (reflist *refList) UpdateRowsImpl(rows []string) (enlarge int) {
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
func (reflist *refList) UpdateRows() (enlarge int) {
	return reflist.UpdateRowsImpl(reflist.Data.List())
}
func (reflist *refList) Render() {
	ui.Render(reflist.List)
}
func (reflist *refList) ReRect(x, y int) {
	reflist.List.SetRect(reflist.Dx, reflist.Dy, x+1, y+1)
	reflist.UpdateRowsImpl(reflist.List.Rows)
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
func (reflist *refList) AddFunction(funcname string) {
	n, err := reflist.Data.AddFunction(funcname)
	if err != nil {
		reflist.UpdateRows()
		return
	}
	enlarge := reflist.UpdateRows()
	if enlarge > 0 {
		reflist.SelectLastRow(enlarge)
	}
	if n == 1 {
		reflist.Data.ReferenceByIndex(reflist.List.SelectedRow)
		reflist.UpdateRows()
	}
}

func hackGetTopRow(l *widgets.List) int {
	pointerVal := reflect.ValueOf(l)
	val := reflect.Indirect(pointerVal)
	member := val.FieldByName("topRow")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*int)(ptrToY)
	return *realPtrToY
}
