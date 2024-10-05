//go:build !windows

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/morphy76/cgnlog/cmd/cli"
	"github.com/morphy76/cgnlog/internal/destination"
	"github.com/morphy76/cgnlog/internal/event"
	"github.com/morphy76/cgnlog/internal/source"
)

var inputFile string
var outFormat string
var keep bool
var help bool
var outFile string
var outDelay time.Duration

const steps = 6

func init() {
	flag.StringVar(&inputFile, "input", "", "Input file path")
	flag.StringVar(&outFormat, "format", "html", "Output format (html or json)")
	flag.BoolVar(&keep, "keep", false, "Keep the output file")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.DurationVar(&outDelay, "delay", 5, "Delay before exit in seconds")

	flag.Usage = func() {
		fmt.Fprint(flag.CommandLine.Output(), "\033[1;34mUsage:\033[0m\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if inputFile == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "\033[1;31mError: input file is required\033[0m\n")
		flag.Usage()
		os.Exit(1)
	}

	if outFormat != "html" && outFormat != "json" {
		fmt.Fprintf(flag.CommandLine.Output(), "\033[1;31mError: invalid format\033[0m\n")
		flag.Usage()
		os.Exit(1)
	}
}

func main() {

	doneChan := make(chan string)
	defer close(doneChan)

	exitChan := make(chan bool)
	defer close(exitChan)

	progressChan := make(chan bool)
	defer close(progressChan)

	endProgressChan := make(chan bool)
	defer close(endProgressChan)

	rowsChan := make(chan event.Event)
	defer close(rowsChan)

	go func() {

		progress := 0

		for {
			select {
			case <-progressChan:
				progress += 100 / steps
				fmt.Printf("\r\033[1;32mProgress: [%-50s] %d%%\033[0m",
					string(strings.Repeat("#", progress/2))+
						string(strings.Repeat(" ", 50-progress/2)),
					progress)
			case <-endProgressChan:
				fmt.Printf("\r\033[1;32mProgress: [%-50s] 100%%\033[0m\n",
					string(strings.Repeat("#", 50)))
				return
			}
		}
	}()

	go destination.WriteTemporaryFile(rowsChan, outFormat, doneChan, progressChan)

	if err := source.ReadInputFile(inputFile, rowsChan, progressChan); err != nil {
		fmt.Println("Failed to read input file:", err)
		os.Exit(1)
	}

	go func() {
		outFile = <-doneChan
		endProgressChan <- true

		err := cli.OpenBrowser(outFile, exitChan)
		if err != nil {
			fmt.Println("Failed to open browser:", err)
		}
	}()

	<-exitChan

	if !keep {
		<-time.After(outDelay * time.Second)
		fmt.Println("Deleting output file... " + outFile)
		if err := os.Remove(outFile); err != nil {
			fmt.Println("Failed to delete output file:", err)
		}
	} else {
		fmt.Println("Output file: " + outFile)
	}

	os.Exit(0)
}
