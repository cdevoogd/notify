package ntfy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const titleHeader = "Title"

// Config stores the configuration required to push notifications to a ntfy.sh instance.
type Config struct {
	Topic string `arg:"-t,--topic,env:NOTIFY_TOPIC,required" help:"Topic name on ntfy.sh"`
}

// Notifier sends push notifications using a topic from https://ntfy.sh.
type Notifier struct {
	topicURL url.URL
}

func NewNotifier(cfg *Config) (*Notifier, error) {
	if cfg == nil {
		return nil, errors.New("config cannot be nil")
	}

	if cfg.Topic == "" {
		return nil, errors.New("a topic is required but was not set")
	}

	n := &Notifier{topicURL: url.URL{Scheme: "https", Host: "ntfy.sh", Path: cfg.Topic}}
	return n, nil
}

func (n *Notifier) SendNotification(title, message string) error {
	req, err := http.NewRequest(http.MethodPost, n.topicURL.String(), strings.NewReader(message))
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
