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
	data []byte
	err  error
)

func usage(s int) {
	fmt.Printf("\n%s\n", `Usage:
	envsubsty [file]
Or
	cat file.txt | envsubsty	
	`)
	os.Exit(s)
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
		usage(1)
	}
}

func envsubsty(s []byte) []byte {
	str := string(s)
	extVar := regexp.MustCompile(`\$(\{[0-9A-Za-z_-]+([:?=+-]{1,2}([^{}]*(\$\{[^{}]+\})*[^{}]*)?)?\}|[0-9A-Za-z_-]+)`)
	for _, i := range extVar.FindAllString(str, -1) {
		out, _ := exec.Command("sh", "-c", `eval printf '%s'`+i).Output()
		if string(out) != "" {
			str = strings.ReplaceAll(str, i, string(out))
		}
	}
	return []byte(str)
}

func main() {
	flag.Parse()
	switch flag.NArg() {
	case 0:
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			data, err = ioutil.ReadAll(os.Stdin)
			check(err)
			data = envsubsty(data)
			fmt.Println(string(data))
		} else {
			usage(1)
		}
		break
	case 1:
		data, err = ioutil.ReadFile(flag.Arg(0))
		check(err)
		data = envsubsty(data)
		ioutil.WriteFile(flag.Arg(0), data, 0644)
		break
	default:
		usage(1)
	}
}
