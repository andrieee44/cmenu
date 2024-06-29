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
)

func help() {
	fmt.Fprintln(os.Stderr, `usage: cmenu MENU [FILE]

cmenu wraps MENU to choose from JSON object in FILE.
if no FILE, read from STDIN.
FILE must have key value pairs of type string.

example: echo '{"key":"value"}' | cmenu fzf

cmenu pipes all keys to MENU.
MENU must output a valid key.
cmenu only outputs the chosen value.`)

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

func exists(name string) bool {
	var err error

	_, err = os.Stat(name)

	return !errors.Is(err, os.ErrNotExist)
}

func dataDir() string {
	const dirName string = "cmenu"

	var dir string

	dir = os.Getenv("XDG_DATA_HOME")
	if dir != "" {
		return filepath.Join(dir, dirName)
	}

	dir = os.Getenv("HOME")
	if dir != "" {
		return filepath.Join(dir, ".local", "share", dirName)
	}

	panic(errors.New("cmenu: $HOME is empty"))
}

func jsonFile() *os.File {
	var (
		f              *os.File
		name, dataFile string
		err            error
	)

	if len(os.Args) == 2 {
		return os.Stdin
	}

	name = os.Args[2]
	if exists(name) {
		f, err = os.Open(name)
		exitIf(err)

		return f
	}

	dataFile = filepath.Join(dataDir(), name)
	if exists(dataFile) {
		f, err = os.Open(dataFile)
		exitIf(err)

		return f
	}

	exit(fmt.Errorf("%s: file does not exist", name))

	return nil
}

func jsonCmds() map[string]string {
	var (
		cmds map[string]string
		f    *os.File
	)

	f = jsonFile()
	exitIf(json.NewDecoder(f).Decode(&cmds))
	exitIf(f.Close())

	return cmds
}

func cmdsKeys(cmds map[string]string) []string {
	var (
		keys []string
		k    string
	)

	for k = range cmds {
		keys = append(keys, k)
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
