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

const version string = "0.1.4"

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
func Convert(inData []byte, varList string) []byte {
	wrkData := string(inData)
	//varNameRegex := regexp.MustCompile(`[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*`)
	varRegex := regexp.MustCompile(
		// Creepy regular expression for finding variables.
		`\$(\{[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*([:?=+-]{1,2}([^{}]*(\$\{[^{}]+\})*[^{}]*)?)?\}|[0-9A-Za-z]([0-9A-Za-z_-]*[0-9A-Za-z])*)`,
	)
	varInUse := varRegex.FindAllString(wrkData, -1)
	if varList != "" {
		varInUse = varRegex.FindAllString(varList, -1)
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

// ConvertFile - work with file
func ConvertFile(path string, varList string, write bool) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	srcData, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	srcData = Convert(srcData, varList)
	if write {
		ioutil.WriteFile(path, srcData, 0644)
	} else {
		fmt.Println(string(srcData))
	}
	return nil
}

// ConvertDir - work with directory
func ConvertDir(path string, varList string, write bool) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !isDir(path) {
			ConvertFile(path+f.Name(), varList, write)
		}
	}
	return nil
}

// Check if it dir
func isDir(path string) bool {
	pathCheck, err := os.Stat(path)
	if err != nil {
		return false
	}
	return pathCheck.Mode().IsDir()
}

func main() {
	flag.Parse()

	if flagHelp {
		fmt.Printf(
			"%s\n\n",
			`envsubsty

Description:
	The envsubsty converts the specified environment variables in files to their value.`,
		)
		Usage(0)
	}

	if flagVer {
		fmt.Printf("Version: %v\n", version)
		os.Exit(0)
	}

	switch flag.NArg() {
	// If no arguments is specified, check Stdin.
	case 0:
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			srcData, err = ioutil.ReadAll(os.Stdin)
			Check(err)
			srcData = Convert(srcData, flagVars)
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
			err := ConvertDir(path, flagVars, flagWrite)
			Check(err)
		} else {
			err := ConvertFile(path, flagVars, flagWrite)
			Check(err)
		}
		break
	default:
		Usage(0)
	}
}
