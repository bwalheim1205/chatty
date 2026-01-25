package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwalheim1205/chatty/internal/llm"
)

type Client struct {
	provider     string
	baseURL      string
	defaultModel llm.ModelID
	http         *http.Client
}

func New(baseURL string, defaultModel llm.ModelID) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	if defaultModel == "" {
		defaultModel = "llama3.2:3b"
	}

	return &Client{
		provider:     "ollama",
		baseURL:      baseURL,
		defaultModel: defaultModel,
		http:         &http.Client{},
	}
}

type chatRequest struct {
	Model    llm.ModelID   `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func (c *Client) Name() string {
	return c.provider
}

func (c *Client) DefaultModel() llm.ModelID {
	return c.defaultModel
}

func (c *Client) Complete(ctx context.Context, req llm.Request) (string, error) {
	model := c.defaultModel
	if req.Model != "" {
		model = req.Model
	}

	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: toOllamaMessages(req.Messages),
		Stream:   false,
	})
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/chat",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama error: %s", resp.Status)
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Message.Content, nil
}

func (c *Client) Stream(ctx context.Context, req llm.Request) (<-chan llm.StreamChunk, error) {
	model := c.defaultModel
	if req.Model != "" {
		model = req.Model
	}

	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: toOllamaMessages(req.Messages),
		Stream:   true,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/api/chat",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}

	ch := make(chan llm.StreamChunk)

	go func() {
		defer resp.Body.Close()
		defer close(ch)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			var msg chatResponse
			if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
				return
			}

			if msg.Message.Content != "" {
				ch <- llm.StreamChunk{Text: msg.Message.Content}
			}

			if msg.Done {
				ch <- llm.StreamChunk{Done: true}
				return
			}
		}
	}()

	return ch, nil
}

func toOllamaMessages(msgs []llm.Message) []chatMessage {
	out := make([]chatMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, chatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return out
}
