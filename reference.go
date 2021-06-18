package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type functionTags struct {
	function string
	line     int
	lineEnd  int
	content  string
	rType    string
	file     string
}

type functionReference struct {
	uuid     string
	function string
	line     int
	file     string
	content  string
	folder   bool
	indent   int
}

type functionReferences []functionReference

func getFunctionNameGlobal(filename string, linenum int) (string, error) {
	globalOut, err := execute("global", "-f", filename)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`(\S+)\s+(\S+)\s+(\S+)\s+(.*)\n`)
	all := re.FindAllStringSubmatch(globalOut, -1)
	var functionSlice []functionTags
	for _, item := range all {
		functionSlice = append(functionSlice, functionTags{
			function: item[1],
			line:     str2int(item[2]),
			content:  item[4],
		})
	}
	function := ""
	for i, functionItem := range functionSlice {
		if i == len(functionSlice)-1 {
			function = functionItem.function
			break
		} else if i == 0 && linenum < functionItem.line {
			function = ""
			break
		} else if linenum >= functionItem.line && linenum < functionSlice[i+1].line {
			function = functionItem.function
			break
		}
	}
	return function, nil
}

func getFunctionNameCtags(filename string, linenum int) (string, error) {
	ctagsOut, err := execute("ctags", "--fields=+ne", "-o", "-", "--sort=no", filename)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`(\S+)\s+\S+\s+/\^(.*)/;"\s+(\S+)\s+.*line:(\d+).*end:(\d+)`)
	all := re.FindAllStringSubmatch(ctagsOut, -1)
	var functionSlice []functionTags
	for _, item := range all {
		functionSlice = append(functionSlice, functionTags{
			function: item[1],
			line:     str2int(item[4]),
			lineEnd:  str2int(item[5]),
			content:  item[3],
			rType:    item[2],
		})
	}
	function := ""
	for _, functionItem := range functionSlice {
		if linenum >= functionItem.line && linenum <= functionItem.lineEnd {
			function = functionItem.function
			break
		}
	}
	return function, nil
}
func getFunctionName(filename string, linenum int) (string, error) {
	if isCtags {
		return getFunctionNameCtags(filename, linenum)
	} else if isGlobal {
		return getFunctionNameGlobal(filename, linenum)
	} else {
		return "", fmt.Errorf("Neither Ctags nor GNU Global exist")
	}
}

func getDefine(function string) ([]functionTags, error) {
	if function == "" {
		return []functionTags{}, fmt.Errorf("Empty function name")
	}
	globalArgs := []string{
		"-x",
		function,
	}
	if extendTag {
		globalArgs = []string{
			"-ix",
			".*" + function + ".*",
		}
	}
	globalOut, err := execute("global", globalArgs...)
	if err != nil {
		return []functionTags{}, err
	}
	re := regexp.MustCompile(`(\S+)\s+(\S+)\s+(\S+)\s+(.*)\n`)
	all := re.FindAllStringSubmatch(globalOut, -1)
	var functionSlice []functionTags
	for _, item := range all {
		filename := item[3]
		linenum := str2int(item[2])
		function, err := getFunctionName(filename, linenum)
		if err != nil || function == "" {
			continue
		}
		functionSlice = append(functionSlice, functionTags{
			file:     filename,
			line:     linenum,
			function: function,
			content:  item[4],
		})
	}
	if len(functionSlice) == 0 {
		return functionSlice, fmt.Errorf("No Define find")
	}
	return functionSlice, nil
}

func getReference(function string) ([]functionTags, error) {
	globalOut, err := execute("global", "-rx", function)
	if err != nil {
		return []functionTags{}, err
	}
	re := regexp.MustCompile(`(\S+)\s+(\S+)\s+(\S+)\s+(.*)\n`)
	all := re.FindAllStringSubmatch(globalOut, -1)
	var functionSlice []functionTags
	for _, item := range all {
		filename := item[3]
		linenum := str2int(item[2])
		function, err := getFunctionName(filename, linenum)
		if err != nil || function == "" {
			continue
		}
		functionSlice = append(functionSlice, functionTags{
			file:     filename,
			line:     linenum,
			function: function,
			content:  item[4],
		})
	}
	return functionSlice, nil
}

func (f *functionReferences) AddFunction(function string) (int, error) {
	defines, err := getDefine(function)
	if err != nil {
		return 0, err
	}
	if len(defines) == 0 {
		return 0, fmt.Errorf("Reference Empty")
	}
	for _, define := range defines {
		*f = append(*f, functionReference{
			uuid:     genUUID(),
			function: define.function,
			line:     define.line,
			file:     define.file,
			content:  define.content,
			indent:   0,
			folder:   true,
		})
	}
	return len(defines), nil
}
func (f *functionReferences) OpenFile(index int) {
	if index >= len(*f) {
		return
	}
	fref := &[]functionReference(*f)[index]
	cmd := exec.Command("vim", fref.file, fmt.Sprintf("+%d", fref.line))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func (f *functionReferences) ReferenceByIndex(index int) {
	fLen := len([]functionReference(*f))
	if index < 0 || index >= fLen {
		return
	}
	if fLen == 0 {
		return
	}

	fref := &[]functionReference(*f)[index]
	if fref.folder {
		fref.folder = false

		refs, err := getReference(fref.function)
		if err != nil {
			return
		}
		var newRefs functionReferences
		for _, ref := range refs {
			newRefs = append(newRefs, functionReference{
				uuid:     genUUID(),
				function: ref.function,
				line:     ref.line,
				file:     ref.file,
				content:  ref.content,
				indent:   fref.indent + 1,
				folder:   true,
			})
		}
		if len(newRefs) == 0 {
			return
		}

		*f = append(*f, newRefs...)
		copy((*f)[index+1+len(newRefs):], (*f)[index+1:])
		if index+1 < fLen {
			copy((*f)[index+1:index+1+len(newRefs)+1], newRefs[:])
		}
	} else {
		fref.folder = true
		var nIndex int
		for nIndex = index + 1; nIndex < fLen; nIndex++ {
			nref := &[]functionReference(*f)[nIndex]
			if fref.indent >= nref.indent {
				break
			}
		}
		copy((*f)[index+1:], (*f)[nIndex:])
		*f = (*f)[:fLen-nIndex+index+1]
	}
}

func (f *functionReferences) ReferenceByUUID(uuid string) {
	index := -1
	for i, ref := range *f {
		if ref.uuid == uuid {
			index = i
		}
	}
	if index == -1 {
		return
	}
	f.ReferenceByIndex(index)
}
func (f *functionReferences) RemoveFunctionByIndex(index int) {
	fLen := len([]functionReference(*f))
	if index < 0 || index >= fLen {
		return
	}

	fref := &[]functionReference(*f)[index]
	if fref.indent != 0 {
		return
	}
	var nIndex int
	for nIndex = index + 1; nIndex < fLen; nIndex++ {
		nref := &[]functionReference(*f)[nIndex]
		if nref.indent == 0 {
			break
		}
	}
	copy((*f)[index:], (*f)[nIndex:])
	*f = (*f)[:fLen-nIndex+index]
}
func (f *functionReferences) RemoveFunctionByUUID(uuid string) {
	index := -1
	for i, ref := range *f {
		if ref.uuid == uuid {
			index = i
		}
	}
	if index == -1 {
		return
	}
	f.RemoveFunctionByIndex(index)
}

func (f *functionReferences) Print() {
	fmt.Println("=======")
	for _, item := range []functionReference(*f) {
		f := ""
		if item.folder {
			f = "+"
		} else {
			f = "-"
		}
		fmt.Printf(
			"%s%s \x1b[31m%s\x1b[0m \x1b[36m%s:%d\x1b[0m %s\n",
			strings.Repeat(" ", 2*item.indent),
			f,
			item.function,
			item.file,
			item.line,
			"", //item.uuid,
		)
	}
}
func (f *functionReferences) List() []string {
	var funcList []string
	for _, item := range []functionReference(*f) {
		f := ""
		if item.folder {
			f = "+"
		} else {
			f = "-"
		}
		funcList = append(funcList, fmt.Sprintf(
			"%s%s [%s](fg:red) %s:%d %s",
			strings.Repeat(" ", 2*item.indent),
			f,
			item.function,
			item.file,
			item.line,
			"", //item.uuid,
		))
	}
	if len(funcList) == 0 {
		funcList = []string{simpleHelp}
	}
	return funcList
}

func (f *functionReferences) Size() int {
	return len([]functionReference(*f))
}
