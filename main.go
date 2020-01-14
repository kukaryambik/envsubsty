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
	srcData  []byte
	err      error
	writeArg bool
	helpArg  bool
	varArg   string
)

func init() {
	// Write into the source file.
	flag.BoolVar(&writeArg, "w", writeArg, "Write the output to the source file.")
	// Help
	flag.BoolVar(&helpArg, "h", helpArg, "Show help message.")
	flag.StringVar(&varArg, "v", varArg, "Comma or space-separated list of variables to convert.")
}

// Usage
func usage(s int) {
	fmt.Printf("%s\n", `Usage:
	envsubsty [-wh] [-v 'vars'] [file|directory ...]
Or
	cat file.txt | envsubsty [-v 'vars']
Flags:`)
	flag.PrintDefaults()
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
	//varNameRegex := regexp.MustCompile(`[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*`)
	varRegex := regexp.MustCompile(
		// Creepy regular expression for finding variables.
		`\$(\{[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*([:?=+-]{1,2}([^{}]*(\$\{[^{}]+\})*[^{}]*)?)?\}|[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*)`,
	)
	varInUse := varRegex.FindAllString(wrkData, -1)
	if varArg != "" {
		varInUse = varRegex.FindAllString(varArg, -1)
	}
	for _, varName := range varInUse {
		// Workaround for simply substitution complex variables (like ${VAR1:=default})
		out, _ := exec.Command("sh", "-c", `eval printf '%s'`+varName).Output()
		if string(out) != "" {
			wrkData = strings.Replace(wrkData, varName, string(out), 1)
		}
	}
	return []byte(wrkData)
}

// Work with file
func envsubstFile(file string, write bool) {
	srcData, err = ioutil.ReadFile(file)
	check(err)
	srcData = envsubsty(srcData)
	if write {
		ioutil.WriteFile(file, srcData, 0644)
	} else {
		fmt.Println(string(srcData))
	}
}

// Check if it dir
func isDir(path string) bool {
	pathCheck, err := os.Stat(path)
	check(err)
	return pathCheck.Mode().IsDir()
}

func main() {
	flag.Parse()

	if helpArg {
		usage(0)
	}

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
	// If an argument is specified, use it as the path.
	case 1:
		var path string = flag.Arg(0)
		// If path is directory, use all files in directory
		if isDir(path) {
			files, err := ioutil.ReadDir(path)
			check(err)
			for _, f := range files {
				if !isDir(path + f.Name()) {
					envsubstFile(path+f.Name(), writeArg)
				}
			}
		} else {
			envsubstFile(flag.Arg(0), writeArg)
		}
		break
	default:
		usage(0)
	}
}
