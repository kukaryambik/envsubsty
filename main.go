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

const version string = "0.1.3"

var (
	srcData   []byte
	err       error
	flagWrite bool
	flagHelp  bool
	flagVer   bool
	flagVars  string
)

func init() {
	// Write into the source file.
	flag.BoolVar(&flagWrite, "w", flagWrite, "Write the output to the source file.")
	// Help
	flag.BoolVar(&flagHelp, "h", flagHelp, "Show help message.")
	// Version
	flag.BoolVar(&flagVer, "V", flagVer, "Show version.")
	// List of variables
	flag.StringVar(&flagVars, "v", flagVars, "Comma or space-separated list of variables to convert.")
}

// Usage message
func Usage(s int) {
	fmt.Printf("%s\n", `Usage:
	envsubsty [-hVw] [-v 'vars'] [file|directory ...]
Or
	cat file.txt | envsubsty [-v 'vars']
Flags:`)
	flag.PrintDefaults()
	os.Exit(s)
}

// Check if error
func Check(e error) {
	if e != nil {
		fmt.Println(e)
		Usage(1)
	}
}

// Convert variables
func Convert(inData []byte) []byte {
	wrkData := string(inData)
	//varNameRegex := regexp.MustCompile(`[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*`)
	varRegex := regexp.MustCompile(
		// Creepy regular expression for finding variables.
		`\$(\{[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*([:?=+-]{1,2}([^{}]*(\$\{[^{}]+\})*[^{}]*)?)?\}|[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*)`,
	)
	varInUse := varRegex.FindAllString(wrkData, -1)
	if flagVars != "" {
		varInUse = varRegex.FindAllString(flagVars, -1)
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

// PrintVer - print version
func PrintVer() {
	fmt.Printf("Version: %v", version)
	os.Exit(0)
}

// ConvertFile - work with files or directory
func ConvertFile(file string, write bool) {
	srcData, err = ioutil.ReadFile(file)
	Check(err)
	srcData = Convert(srcData)
	if write {
		ioutil.WriteFile(file, srcData, 0644)
	} else {
		fmt.Println(string(srcData))
	}
}

// Check if it dir
func isDir(path string) bool {
	pathCheck, err := os.Stat(path)
	Check(err)
	return pathCheck.Mode().IsDir()
}

func main() {
	flag.Parse()

	if flagHelp {
		Usage(0)
	}

	if flagVer {
		PrintVer()
	}

	switch flag.NArg() {
	// If no arguments is specified, check Stdin.
	case 0:
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			srcData, err = ioutil.ReadAll(os.Stdin)
			Check(err)
			srcData = Convert(srcData)
			fmt.Println(string(srcData))
		} else {
			Usage(0)
		}
		break
	// If an argument is specified, use it as the path.
	case 1:
		var path string = flag.Arg(0)
		// If path is directory, use all files in directory
		if isDir(path) {
			files, err := ioutil.ReadDir(path)
			Check(err)
			for _, f := range files {
				if !isDir(path + f.Name()) {
					ConvertFile(path+f.Name(), flagWrite)
				}
			}
		} else {
			ConvertFile(flag.Arg(0), flagWrite)
		}
		break
	default:
		Usage(0)
	}
}
