package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jessevdk/go-flags"
)

const (
	notifyHost  = "ntfy.sh"
	titleHeader = "Title"
)

var opts struct {
	CustomMessage string `short:"m" long:"message" description:"Set a custom message to send"`
	Topic         string `short:"t" long:"topic" env:"NOTIFY_TOPIC" description:"Topic name to use on ntfy.sh"`
	Command       struct {
		Name      string   `positional-arg-name:"COMMAND" description:"The command for notify to execute"`
		Arguments []string `positional-arg-name:"ARGUMENTS" description:"Arguments to pass to the command. You can use -- before specifying the command to prevent notify from parsing these arguments."`
	} `positional-args:"true" required:"true"`
}

func toStderr(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func sendNotification(title, message string) error {
	topicURL := url.URL{Scheme: "https", Host: notifyHost, Path: opts.Topic}

	fmt.Printf("Sending notification to %s\n", topicURL.String())
	fmt.Printf("  Title: %s\n", title)
	fmt.Printf("  Message: %s\n", message)

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
	_, err := flags.Parse(&opts)
	if err != nil {
		if flags.WroteHelp(err) {
			os.Exit(0)
		}

		// The flags package will automatically print out errors
		os.Exit(1)
	}

	if opts.Topic == "" {
		toStderr("A topic from ntfy.sh is required but was not set")
		toStderr("Set a topic by setting the NOTIFY_TOPIC environment variable or by using the -t/--topic flags")
		os.Exit(1)
	}

	fmt.Println("Topic:", opts.Topic)
	fmt.Printf("Command: %+v\n", opts)

	cmd := exec.Command(opts.Command.Name, opts.Command.Arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	start := time.Now()
	err = cmd.Start()
	if err != nil {
		toStderr("Error starting command:", err)
		os.Exit(1)
	}

	_ = cmd.Wait()
	elapsed := time.Since(start)

	title := fmt.Sprintf("%q has finished", opts.Command.Name)
	message := opts.CustomMessage
	if message == "" {
		message = fmt.Sprintf("Exit Code: %d, Elapsed: %s", cmd.ProcessState.ExitCode(), elapsed)
	}

	err = sendNotification(title, message)
	if err != nil {
		toStderr("Error sending notification:", err)
		os.Exit(1)
	}
}
