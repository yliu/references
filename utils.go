package main

import (
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func execute(name string, arg ...string) (string, error) {
	out, err := exec.Command(name, arg...).CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}

func checkCtags() error {
	out, err := execute("ctags", "--version")
	if err != nil {
		return fmt.Errorf("%s", out)
	}
	re := regexp.MustCompile(`Universal Ctags`)
	match := re.Match([]byte(out))
	if match {
		return nil
	}
	return fmt.Errorf("It uses Universal Ctags instead of Exuberant Ctags")
}

func checkGlobal() error {
	xrand := fmt.Sprintf("Rand%d", rand.Uint64())
	out, err := execute("global", "-rx", xrand)
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(out))
	}
	return nil
}

func genUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return "00000000-0000-0000-0000-000000000000"
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func str2int(str string) int {
	num, _ := strconv.ParseInt(str, 10, 64)
	return int(num)
}
