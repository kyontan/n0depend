package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/koron/go-dproxy"
)

func DagToGraphEasyFormat(dag interface{}) (*string, error) {
	d := dproxy.New(dag)

	m, err := d.Map()
	if err != nil {
		return nil, err
	}

	var b strings.Builder

	for name, _ := range m {
		_, err := d.M(name).Map()
		if err != nil {
			return nil, err
		}

		depends, err := d.M(name).M("depends_on").Array()

		if err != nil {
			continue // task do not have `depends_on`
		}

		for _, depend := range depends {

			fmt.Fprintf(&b, "[%v]->[%v]", depend, name)
		}
	}

	s := b.String()

	return &s, nil
}

func ExecGraphEasy(input string) error {
	cmd := exec.Command("graph-easy")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		return err
	}

	// draw graph flows right as long as possible
	// http://bloodgate.com/perl/graph/manual/hinting.html#flow
	io.WriteString(stdin, "graph { flow: left; }\n")

	io.WriteString(stdin, input)
	stdin.Close()
	return cmd.Wait()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, os.Args[0], "[n0stack dag]")
		fmt.Fprintln(os.Stderr, " This program prints dependency graph of n0stack dag file(yaml).")
		fmt.Fprintln(os.Stderr, " Requires Graph::Easy perl package.")
		os.Exit(1)
	}

	filename := os.Args[1]
	b, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var dag interface{}

	err = yaml.Unmarshal(b, &dag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	graph_easy_input, err := DagToGraphEasyFormat(dag)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ExecGraphEasy(*graph_easy_input)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
