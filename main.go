package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/cdevoogd/notify/internal/log"
	"github.com/cdevoogd/notify/internal/ntfy"
)

var args struct {
	ntfy.Config
	CustomMessage string   `arg:"-m,--message" help:"Override the message sent in the notification"`
	Command       string   `arg:"positional,required" help:"The command to execute"`
	Args          []string `arg:"positional" help:"Arguments to pass to the command. Use -- to prevent notify from parsing these arguments."`
}

func main() {
	arg.MustParse(&args)

	notifier, err := ntfy.NewNotifier(&args.Config)
	if err != nil {
		log.Fatal("Error initializing notifier:", err)
	}

	cmd := exec.Command(args.Command, args.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err = cmd.Start()
	if err != nil {
		log.Fatal("Error starting command:", err)
	}

	// Ignoring the error here since non-zero exits will return an error.
	_ = cmd.Wait()

	title := fmt.Sprintf("%q has finished", args.Command)
	message := args.CustomMessage
	if message == "" {
		ec := cmd.ProcessState.ExitCode()
		elapsed := time.Since(start).Round(time.Millisecond)
		message = fmt.Sprintf("Exit Code: %d, Elapsed: %s", ec, elapsed)
	}

	err = notifier.SendNotification(title, message)
	if err != nil {
		log.Fatal("Error sending notification:", err)
	}
}
