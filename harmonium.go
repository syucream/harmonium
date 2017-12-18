package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const format = `(sh|bash|csh|ksh|tcsh|zsh|)`
const tempFileName = "harmonium.sh"

func getScripts(raw []byte) ([]string, error) {
	var scripts []string

	// regexp for start sh block
	startRe, err := regexp.Compile(`(?m)^` + "```" + format + `\s*$`)
	if err != nil {
		return scripts, err
	}

	// find all index started by sh block
	matched := startRe.FindAllSubmatchIndex(raw, -1)
	text := string(raw)

	for _, m := range matched {
		// find end of supported block
		endRe, err := regexp.Compile(`(?m)^` + "```" + `$`)
		if err != nil {
			return scripts, err
		}

		offsets := endRe.FindIndex(raw[m[1]:])
		if len(offsets) > 0 {
			// m[1] ~ m[1] + offsets[0] might be a sh script
			scripts = append(scripts, text[m[1]:m[1]+offsets[0]])
		}
	}

	return scripts, nil
}

func runScript(script string) error {
	tmpFile, err := ioutil.TempFile("", tempFileName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(script); err != nil {
		return err
	}

	return exec.Command("sh", tmpFile.Name()).Run()
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: harmonium (run|extract) <file>")
		os.Exit(1)
	}

	subCommand := os.Args[1]
	filepath := os.Args[2]

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "given filepath is invalid:", filepath)
		os.Exit(1)
	}

	scripts, err := getScripts(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse error at", filepath)
		os.Exit(1)
	}

	joinedScript := strings.Join(scripts, "\n")

	if subCommand == "run" {
		if err := runScript(joinedScript); err != nil {
			fmt.Fprintln(os.Stderr, "execution failed at", filepath)
			os.Exit(1)
		}

		fmt.Println("execution successed!")
	} else if subCommand == "extract" {
		fmt.Println(joinedScript)
	} else {
		fmt.Fprintln(os.Stderr, "subcommand is invalid:", subCommand)
		os.Exit(1)
	}
}
