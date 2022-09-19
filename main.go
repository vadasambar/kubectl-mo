/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var namespaces = []string{}

// https://github.com/ahmetb/kubectx/blob/master/kubens#L26
var kubensDir string

func main() {
	h, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubensDir = filepath.Join(h, ".kube", "kubens")

	args := []string{}
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, "-ns") {
			a := strings.Replace(arg, "-ns=", "", 1)
			for _, n := range strings.Split(a, ",") {
				namespaces = append(namespaces, strings.TrimSpace(n))
			}
			continue
		}
		args = append(args, strings.TrimSpace(arg))
	}

	// If no `-ns` is specified, fall back to current and last switched-to namespace
	if len(namespaces) == 0 {
		namespaces = append(namespaces, getCurretnNs())

		prevNsFile := filepath.Join(kubensDir, getCurrentCtx())
		o, err := os.ReadFile(prevNsFile)
		if err != nil {
			panic(err)
		}

		namespaces = append(namespaces, string(o))
	}

	for _, n := range namespaces {

		a := []string{}
		a = append(a, "-n", n)

		execKubectlCommand(append(a, args...), true, false)
	}

}

func execKubectlCommand(a []string, printOutput, returnOutput bool) string {
	var stdout bytes.Buffer
	c := exec.Command("kubectl", a...)
	if printOutput {
		fmt.Println("executing command", c.String())
	}
	c.Stderr = &stdout
	c.Stdout = &stdout

	if err := c.Run(); err != nil {
		if printOutput {
			fmt.Println(stdout.String())
		}

		if returnOutput {
			return strings.TrimSpace(stdout.String())
		}
		return ""
	}

	if printOutput {
		fmt.Println(stdout.String())

	}
	return strings.TrimSpace(stdout.String())
}

func getCurretnNs() string {
	// kubectl config view --minify -o jsonpath='{..namespace}'
	cmd := []string{
		"config", "view", "--minify", "-o", "jsonpath={..namespace}",
	}
	return execKubectlCommand(cmd, false, true)
}

func getCurrentCtx() string {
	// kubectl config current-context
	cmd := []string{
		"config", "current-context",
	}
	return execKubectlCommand(cmd, false, true)
}
