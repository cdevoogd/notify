package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
)

const (
	notifyHost  = "ntfy.sh"
	titleHeader = "Title"
)

var args struct {
	CustomMessage string   `arg:"-m,--message" help:"Override the message sent in the notification"`
	Topic         string   `arg:"-t,--topic,env:NOTIFY_TOPIC" help:"Topic name on ntfy.sh"`
	Command       string   `arg:"positional,required" help:"The command to execute"`
	Args          []string `arg:"positional" help:"Arguments to pass to the command. Use -- to prevent notify from parsing these arguments."`
}

func toStderr(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func sendNotification(title, message string) error {
	topicURL := url.URL{Scheme: "https", Host: notifyHost, Path: args.Topic}

	req, err := http.NewRequest(http.MethodPost, topicURL.String(), strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("error creating POST request: %w", err)
	}

	req.Header.Set(titleHeader, title)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func main() {
	arg.MustParse(&args)

	if args.Topic == "" {
		toStderr("A topic from ntfy.sh is required but was not set")
		toStderr("Set a topic by setting the NOTIFY_TOPIC environment variable or by using the -t/--topic flags")
		os.Exit(1)
	}

	cmd := exec.Command(args.Command, args.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err := cmd.Run()
	if err != nil {
		toStderr("Error starting command:", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	elapsed = elapsed.Round(time.Millisecond)

	title := fmt.Sprintf("%q has finished", args.Command)
	message := args.CustomMessage
	if message == "" {
		message = fmt.Sprintf("Exit Code: %d, Elapsed: %s", cmd.ProcessState.ExitCode(), elapsed)
	}

	err = sendNotification(title, message)
	if err != nil {
		toStderr("Error sending notification:", err)
		os.Exit(1)
	}
}
