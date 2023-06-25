package resource

import (
	"github.com/go-ai-agent/core/runtime"
)

// ContentMap - slice of any content to be included in a message
type ContentMap map[string][]any

// Credentials - type for a credentials function
type Credentials func() (username string, password string, err error)

// DatabaseUrl - struct for a database connection string
type DatabaseUrl struct {
	Url string
}

func NewStatusCode(status **runtime.Status) func() int {
	return func() int {
		return int((*(status)).Code())
	}
}

// AccessCredentials - access function for Credentials in a message
func AccessCredentials(msg *Message) Credentials {
	if msg == nil || msg.Content == nil {
		return nil
	}
	for _, c := range msg.Content {
		if fn, ok := c.(Credentials); ok {
			return fn
		}
	}
	return nil
}

// AccessDatabaseUrl - access function for a DatabaseUrl in a message
func AccessDatabaseUrl(msg *Message) DatabaseUrl {
	if msg == nil || msg.Content == nil {
		return DatabaseUrl{}
	}
	for _, c := range msg.Content {
		if url, ok := c.(DatabaseUrl); ok {
			return url
		}
	}
	return DatabaseUrl{}
}
