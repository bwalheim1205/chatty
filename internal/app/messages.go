package app

type LLMCompleteMsg struct {
	Text string
	Err  error
}

type LLMStreamChunk struct {
	Text string
	Done bool
	Err  error
}
