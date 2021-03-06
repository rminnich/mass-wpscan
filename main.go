// Copyright 2017 Jayson Grace. All rights reserved
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"log"
	"os"
	"sync"
)

var (
	inputFile string
	wpParams  string
	outfile   string
	errmsg = color.Red
	warn = color.Yellow
	msg = color.Green
)

func fatal(exitval int, fmt string, args ...interface{}) {
	errmsg(fmt, args...)
	os.Exit(exitval)
}

func scanTargets(targets []string, wpParams string, cmdOutput []string, wg *sync.WaitGroup) []string {
	var output string
	for _, target := range targets {
		msg("Scanning %s with wpscan, please wait...", target)
		cmd := "wpscan" + " --url " + target + " " + wpParams
		wg.Add(1)
		output = exeCmd(string(cmd), wg)
		cmdOutput = append(cmdOutput, output)
		wg.Wait()
	}
	return cmdOutput
}

func main() {
	var cmdOutput []string

	flag.Parse()

	// if there's no input, print usage
	if flag.NFlag() == 0 || validateInput() == false {
		usage()
	}

	paramSlice := splitStringSpaceSlice(wpParams)

	validateWpParams(paramSlice)

	wg := new(sync.WaitGroup)
	msg("Updating wpscan, please wait...")
	wg.Add(1)
	output := exeCmd("wpscan --update", wg)
	wg.Wait()
	cmdOutput = append(cmdOutput, output)

	// Get targets
	targets, err := readLines(inputFile)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	cmdOutput = scanTargets(targets, wpParams, cmdOutput, wg)

	if outfile != "" {
		if err := writeLines(cmdOutput, outfile); err != nil {
			log.Fatalf("writeLines: %s", err)
		}
	} else {
		// No output file has been specified - print output from the command
		fmt.Println(cmdOutput)
	}
}

func init() {
	flag.StringVar(&inputFile, "i", "", "Input file with targets.")
	flag.StringVar(&wpParams, "p", "", "Arguments to run with wpscan.")
	flag.StringVar(&outfile, "o", "", "File to output information to.")
}

func usage() {
	os.Args[0] = os.Args[0] + " [options]"
	flag.Usage()
	os.Exit(1)
}

func validateInput() bool {
	if inputFile == "" || wpParams == "" {
		errmsg("You must specify an input file with targets and parameters for wpscan!")
		errmsg("Example: mass-wpscan -i vuln_targets.txt -p \"-r --batch -e vt,tt,u,vp\"")
		return false
	}
	return true
}

func validateWpParams(parameters []string) {
	for _, p := range parameters {
		if p == "--url" {
			fatal(1, "You can not include the --url parameter, all targets should be in your input file!")
		}
	}
}
