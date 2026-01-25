// internal/llm/llm.go
package llm

import "context"

type Message struct {
	Role    string // "system" | "user" | "assistant"
	Content string
}

type Request struct {
	Model    ModelID
	Messages []Message
}

type StreamChunk struct {
	Text string
	Done bool
}

type Client interface {
	Name() string
	DefaultModel() ModelID
	Complete(ctx context.Context, req Request) (string, error)
	Stream(ctx context.Context, req Request) (<-chan StreamChunk, error)
}

type ModelID string
