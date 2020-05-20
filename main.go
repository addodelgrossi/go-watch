package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	version         = "0.3.1"
	defaultInterval = 1000
	maxArguments    = 128
	quiet           = false
	halt            = 0
)

func main() {

	val := strings.Join(os.Args[1:], " ")
	commands := []string{"-c", val}
	duration := 0 * time.Second
	// halt := false

	// fmt.Printf("commands %s, val %s, os %s\n", commands, val, runtime.GOOS)

	fmt.Print("\033[s") // cursor mark

	done := make(chan bool)

	// if quiet {
	// 	os.Stdout = os.NewFile(uintptr(syscall.Stdin), os.DevNull)
	// }

	// type IOStreams struct {
	// 	out    io.Writer
	// 	errOut io.Writer
	// }

	// streams := IOStreams{
	// 	out:    os.Stdout,
	// 	errOut: os.Stderr,
	// }

	go func() {
		finishAt := time.Now().Add(duration)
		for {

			if duration != 0 && time.Now().After(finishAt) {
				fmt.Printf("exit by duration")
				close(done)
			}

			cmd := exec.Command("sh", commands...)
			out, err := cmd.CombinedOutput()

			fmt.Print("\033[u\033[K") // reset cursor

			fmt.Println(err)

			if !quiet || err != nil {
				fmt.Printf("%s", out)
			}

			time.Sleep(2 * time.Second)
		}
	}()

	<-done

}