package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	version = "0.3.1"
)

var (
	interval int
	timeout  int
	quiet    bool
	halt     bool
	help     bool
)

func main() {
	flag.IntVar(&interval, "interval", 2, "interval in seconds defaulting to 1")
	flag.IntVar(&timeout, "timeout", -1, "duration until the program ends")
	flag.BoolVar(&quiet, "quiet", false, "only output stderr")
	flag.BoolVar(&halt, "halt", false, "halt on failure")
	flag.BoolVar(&help, "help", false, "output this help information")

	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()

	val := strings.Join(flag.Args(), " ")
	commands := []string{"-c", val}

	fmt.Print("\033[2J")   // clear screen
	fmt.Print("\033[1;1H") // move to position 1,1
	fmt.Print("\033[s")    // cursor mark

	done := make(chan bool)

	go func() {
		finishAt := time.Now().Add(time.Duration(timeout) * time.Second)
		for {

			if timeout > 0 && time.Now().After(finishAt) {
				close(done)
			}

			cmd := exec.Command("sh", commands...)
			out, err := cmd.CombinedOutput()

			fmt.Print("\033[u\033[K") // reset cursor

			if !quiet || err != nil {
				fmt.Printf("%s", out)
			}

			if halt && err != nil {
				exitCode := 124
				if xerr, ok := err.(*exec.ExitError); ok {
					exitCode = xerr.ExitCode()
				}
				os.Exit(exitCode)
			}

			time.Sleep(time.Duration(interval) * time.Second)
		}
	}()

	<-done
}
