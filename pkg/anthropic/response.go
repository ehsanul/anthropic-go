package anthropic

// CompletionResponse is the response from the Anthropic API for a completion request.
type CompletionResponse struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Stop       string `json:"stop"`
}

// StreamResponse is the response from the Anthropic API for a stream of completions.
type StreamResponse struct {
	Completion string `json:"completion"`
	StopReason string `json:"stop_reason"`
	Model      string `json:"model"`
	Stop       string `json:"stop"`
	LogID      string `json:"log_id"`
}

// MessageResponse is a subset of the response from the Anthropic API for a message response.
type MessagePartResponse struct {
	Type string `json:"type"`

	/* for type = "text" */
	Text string `json:"text,omitempty"`

	/* for type = "tool_use" */
	ID    string                 `json:"id,omitempty"`
	Name  string                 `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`
}

// MessageResponse is the response from the Anthropic API for a message response.
type MessageResponse struct {
	ID           string                `json:"id"`
	Type         string                `json:"type"`
	Model        string                `json:"model"`
	Role         string                `json:"role"`
	Content      []MessagePartResponse `json:"content"`
	StopReason   string                `json:"stop_reason"`
	Stop         string                `json:"stop"`
	StopSequence string                `json:"stop_sequence"`
	Usage        MessageUsage          `json:"usage"`
}

type MessageUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type MessageStreamResponse struct {
	Type         string             `json:"type"`
	Message      MessageResponse    `json:"message,omitempty"`
	ContentBlock ContentBlock       `json:"content_block,omitempty"`
	Delta        MessageStreamDelta `json:"delta,omitempty"`
	Usage        MessageStreamUsage `json:"usage,omitempty"`

	Index int `json:"index,omitempty"`
}

type MessageStreamDelta struct {
	Type string `json:"type,omitempty"`

	// for delta type = "text_delta"
	Text string `json:"text,omitempty"`

	// for delta type = "input_json_delta"
	PartialJson string `json:"partial_json,omitempty"`

	// for chunk/event/MessageStreamResponse type = "message_delta", and delta type = "" (i.e. omitted)
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
}

type MessageStreamUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
