package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

func help() {
	fmt.Fprintln(os.Stderr, "usage: cmenu MENU [FILE]")
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

func panicIf(err error) {
	if err != nil {
		panic(fmt.Errorf("cmenu: %s", err))
	}
}

func jsonFile() *os.File {
	var (
		f   *os.File
		err error
	)

	if len(os.Args) == 2 {
		return os.Stdin
	}

	f, err = os.Open(os.Args[2])
	exitIf(err)

	return f
}

func jsonCmds() map[string]string {
	var (
		cmds map[string]string
		f    *os.File
	)

	cmds = make(map[string]string)
	f = jsonFile()
	exitIf(json.NewDecoder(f).Decode(&cmds))
	panicIf(f.Close())

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

	return keys
}

func keyReader(keys []string) io.Reader {
	slices.Sort(keys)

	return strings.NewReader(strings.Join(keys, "\n"))
}

func keyMenu(keys []string) string {
	var (
		cmd *exec.Cmd
		buf bytes.Buffer
	)

	cmd = exec.Command("/bin/sh", "-c", os.Args[1])
	cmd.Stdin = keyReader(keys)
	cmd.Stdout = &buf
	panicIf(cmd.Err)
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
