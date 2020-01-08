/*
 * ----------------------------------------------------------------------------
 * "THE BEER-WARE LICENSE" (Revision 42):
 * <kukaryambik@gmail.com> wrote this file. As long as you retain this notice you
 * can do whatever you want with this stuff. If we meet some day, and you think
 * this stuff is worth it, you can buy me a beer in return.
 * ----------------------------------------------------------------------------
 */

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	srcData []byte
	err     error
)

// Usage
func usage(s int) {
	fmt.Printf("\n%s\n", `Usage:
	envsubsty [file]
Or
	cat file.txt | envsubsty	
	`)
	os.Exit(s)
}

// Check if error
func check(e error) {
	if e != nil {
		fmt.Println(e)
		usage(1)
	}
}

// Substitution of variables
func envsubsty(inData []byte) []byte {
	wrkData := string(inData)
	varRegex := regexp.MustCompile(
		// Creepy regular expression for finding variables.
		`\$(\{[0-9A-Za-z_-]+([:?=+-]{1,2}([^{}]*(\$\{[^{}]+\})*[^{}]*)?)?\}|[0-9A-Za-z_-]+)`,
	)
	for _, varName := range varRegex.FindAllString(wrkData, -1) {
		// Workaround for simply substitution complex variables (like ${VAR1:=default})
		out, _ := exec.Command("sh", "-c", `eval printf '%s'`+varName).Output()
		if string(out) != "" {
			wrkData = strings.ReplaceAll(wrkData, varName, string(out))
		}
	}
	return []byte(wrkData)
}

func main() {
	flag.Parse()
	switch flag.NArg() {
	// If no arguments is specified, check Stdin.
	case 0:
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			srcData, err = ioutil.ReadAll(os.Stdin)
			check(err)
			srcData = envsubsty(srcData)
			fmt.Println(string(srcData))
		} else {
			usage(1)
		}
		break
	// If an argument is specified, use it as the path to the file.
	case 1:
		srcData, err = ioutil.ReadFile(flag.Arg(0))
		check(err)
		srcData = envsubsty(srcData)
		ioutil.WriteFile(flag.Arg(0), srcData, 0644)
		break
	default:
		usage(1)
	}
}
