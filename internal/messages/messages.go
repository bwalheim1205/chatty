package messages

type ChatMessage struct {
	Role    string
	Content string
}

type NewLLMResponse struct {
	Text string
}

type ErrLLM struct {
	Err error
}
