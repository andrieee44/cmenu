package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"github.com/adrg/xdg"
)

func help() {
	fmt.Fprintln(os.Stderr, "usage: cmenu <MENU> [<FILE>]")
	os.Exit(1)
}

func exit(err error) {
	fmt.Fprintln(os.Stderr, "cmenu:", err)
	os.Exit(1)
}

func exitIf(err error) {
	if err != nil {
		exit(err)
	}
}

func exists(path string) bool {
	var err error

	_, err = os.Stat(path)
	if err == nil {
		return true
	}

	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	exit(err)

	return false
}

func jsonFile() *os.File {
	var (
		file           *os.File
		path, dataFile string
		err            error
	)

	if len(os.Args) == 2 {
		return os.Stdin
	}

	path = os.Args[2]
	if exists(path) {
		file, err = os.Open(path)
		exitIf(err)

		return file
	}

	dataFile = filepath.Join(xdg.DataHome, "cmenu", path)
	if exists(dataFile) {
		file, err = os.Open(dataFile)
		exitIf(err)

		return file
	}

	exit(fmt.Errorf("%s: file does not exist", path))

	return nil
}

func jsonCmds() map[string]string {
	var (
		cmds map[string]string
		file *os.File
	)

	file = jsonFile()
	exitIf(json.NewDecoder(file).Decode(&cmds))
	exitIf(file.Close())

	return cmds
}

func cmdsKeys(cmds map[string]string) []string {
	var (
		keys []string
		key  string
	)

	for key = range cmds {
		keys = append(keys, key)
	}

	slices.Sort(keys)

	return keys
}

func keyMenu(keys []string) string {
	var (
		cmd *exec.Cmd
		buf bytes.Buffer
	)

	cmd = exec.Command("/bin/sh", "-c", "--", os.Args[1])
	cmd.Stdin = strings.NewReader(strings.Join(keys, "\n"))
	cmd.Stdout = &buf
	exitIf(cmd.Err)
	exitIf(cmd.Run())

	return strings.TrimSuffix(buf.String(), "\n")
}

func main() {
	var (
		cmds     map[string]string
		key, val string
		ok       bool
	)

	if len(os.Args) == 1 {
		help()
	}

	cmds = jsonCmds()
	key = keyMenu(cmdsKeys(cmds))

	val, ok = cmds[key]
	if !ok {
		exit(fmt.Errorf("%s is not a valid key", key))
	}

	fmt.Println(val)
}
