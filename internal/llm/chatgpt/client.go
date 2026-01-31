package chatgpt

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/bwalheim1205/chatty/internal/llm"
)

type Client struct {
	provider     string
	apiKey       string
	baseURL      string
	defaultModel llm.ModelID
	http         *http.Client
}

func New(apiKey string, defaultModel llm.ModelID) *Client {
	if defaultModel == "" {
		defaultModel = "gpt-5-mini"
	}

	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	return &Client{
		provider:     "chatgpt",
		apiKey:       apiKey,
		baseURL:      "https://api.openai.com/v1",
		defaultModel: defaultModel,
		http:         &http.Client{},
	}
}

// request/response types for ChatGPT
type chatRequest struct {
	Model    llm.ModelID   `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		Delta        chatMessage `json:"delta"` // used for streaming
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

// Name returns the provider name
func (c *Client) Name() string {
	return c.provider
}

// DefaultModel returns the default model
func (c *Client) DefaultModel() llm.ModelID {
	return c.defaultModel
}

// Complete performs a non-streaming request
func (c *Client) Complete(ctx context.Context, req llm.Request) (string, error) {
	model := c.defaultModel
	if req.Model != "" {
		model = req.Model
	}

	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: toChatGPTMessages(req.Messages),
	})
	if err != nil {
		return "", err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("chatgpt error: %s", resp.Status)
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	if len(out.Choices) == 0 {
		return "", fmt.Errorf("chatgpt returned no choices")
	}

	return out.Choices[0].Message.Content, nil
}

// Stream performs a streaming request
func (c *Client) Stream(ctx context.Context, req llm.Request) (<-chan llm.StreamChunk, error) {
	model := c.defaultModel
	if req.Model != "" {
		model = req.Model
	}

	body, err := json.Marshal(chatRequest{
		Model:    model,
		Messages: toChatGPTMessages(req.Messages),
		Stream:   true,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

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
			line := scanner.Text()
			// OpenAI sends lines like "data: {...}" for streaming
			if len(line) < 6 || line[:5] != "data:" {
				continue
			}
			data := line[5:]
			if data == "[DONE]" {
				ch <- llm.StreamChunk{Done: true}
				return
			}

			var msg chatResponse
			if err := json.Unmarshal([]byte(data), &msg); err != nil {
				continue
			}

			for _, choice := range msg.Choices {
				if choice.Delta.Content != "" {
					ch <- llm.StreamChunk{Text: choice.Delta.Content}
				}
			}
		}
	}()

	return ch, nil
}

// Helper to convert llm.Messages to OpenAI format
func toChatGPTMessages(msgs []llm.Message) []chatMessage {
	out := make([]chatMessage, 0, len(msgs))
	for _, m := range msgs {
		out = append(out, chatMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	return out
}
