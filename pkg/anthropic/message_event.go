package anthropic

import (
	"encoding/json"
	"fmt"
)

type MessageEvent struct {
	Type string `json:"type"`
}

type MessageStartEvent struct {
	MessageEvent
	Message struct {
		ID           string        `json:"id"`
		Type         string        `json:"type"`
		Role         string        `json:"role"`
		Content      []interface{} `json:"content"`
		Model        string        `json:"model"`
		StopReason   string        `json:"stop_reason"`
		StopSequence string        `json:"stop_sequence"`
		Usage        struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	} `json:"message"`
}

type ContentBlockStartEvent struct {
	MessageEvent
	Index        int `json:"index"`
	ContentBlock struct {
		Type string `json:"type"`
		Text string `json:"text,omitempty"`

		/* Optional fields, only present for tools responses */
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"content_block"`
}

type PingEvent struct {
	MessageEvent
}

type ContentBlockDeltaEvent struct {
	MessageEvent
	Index int `json:"index"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`

		// Optional fields, only present for tools responses
		PartialJson string `json:"partial_json,omitempty"` // for type = "input_json_delta"
	} `json:"delta"`
}

type ContentBlockStopEvent struct {
	MessageEvent
	Index int `json:"index"`
}

type MessageDeltaEvent struct {
	MessageEvent
	Delta struct {
		Type string `json:"type"`

		StopReason   string `json:"stop_reason"`
		StopSequence string `json:"stop_sequence"`
	} `json:"delta"`
	Usage struct {
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type MessageStopEvent struct {
	MessageEvent
}

type MessageErrorEvent struct {
	MessageEvent
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

type UnsupportedEventType struct {
	Msg  string
	Code int
}

func (e UnsupportedEventType) Error() string {
	return e.Msg
}

func ParseMessageEvent(eventType MessageEventType, event string) (*MessageStreamResponse, error) {
	messageStreamResponse := &MessageStreamResponse{}
	var err error

	switch eventType {
	case MessageEventTypeMessageStart:
		messageStartEvent := &MessageStartEvent{}
		err = json.Unmarshal([]byte(event), &messageStartEvent)

		messageStreamResponse.Type = messageStartEvent.Type
		messageStreamResponse.Message.ID = messageStartEvent.Message.ID
		messageStreamResponse.Message.Type = messageStartEvent.Message.Type
		messageStreamResponse.Message.Role = messageStartEvent.Message.Role
		messageStreamResponse.Message.Model = messageStartEvent.Message.Model
		messageStreamResponse.Message.StopReason = messageStartEvent.Message.StopReason
		messageStreamResponse.Message.StopSequence = messageStartEvent.Message.StopSequence
		messageStreamResponse.Message.Usage = messageStartEvent.Message.Usage
		// set on both levels: no reason for these to ever be out of sync
		messageStreamResponse.Usage = messageStartEvent.Message.Usage

	case MessageEventTypeContentBlockStart:
		contentBlockEvent := &ContentBlockStartEvent{}
		err = json.Unmarshal([]byte(event), &contentBlockEvent)

		messageStreamResponse.Type = contentBlockEvent.Type
		switch contentBlockEvent.ContentBlock.Type {
		case "text":
			messageStreamResponse.ContentBlock = TextContentBlock{
				Type: contentBlockEvent.ContentBlock.Type,
				Text: contentBlockEvent.ContentBlock.Text,
			}
		case "tool_use":
			messageStreamResponse.ContentBlock = ToolUseContentBlock{
				ID:   contentBlockEvent.ContentBlock.ID,
				Type: contentBlockEvent.ContentBlock.Type,
				Name: contentBlockEvent.ContentBlock.Name,
			}
		}
	case MessageEventTypePing:
		pingEvent := &PingEvent{}
		err = json.Unmarshal([]byte(event), &pingEvent)

		messageStreamResponse.Type = pingEvent.Type
	case MessageEventTypeContentBlockDelta:
		contentBlockEvent := &ContentBlockDeltaEvent{}
		err = json.Unmarshal([]byte(event), &contentBlockEvent)

		messageStreamResponse.Type = contentBlockEvent.Type
		messageStreamResponse.Delta.Type = contentBlockEvent.Delta.Type
		messageStreamResponse.Delta.Text = contentBlockEvent.Delta.Text
		messageStreamResponse.Delta.PartialJson = contentBlockEvent.Delta.PartialJson
	case MessageEventTypeContentBlockStop:
		contentBlockStopEvent := &ContentBlockStopEvent{}
		err = json.Unmarshal([]byte(event), &contentBlockStopEvent)

		messageStreamResponse.Type = contentBlockStopEvent.Type
	case MessageEventTypeMessageDelta:
		messageDeltaEvent := &MessageDeltaEvent{}
		err = json.Unmarshal([]byte(event), &messageDeltaEvent)

		messageStreamResponse.Type = messageDeltaEvent.Type
		messageStreamResponse.Delta.StopReason = messageDeltaEvent.Delta.StopReason
		messageStreamResponse.Delta.StopSequence = messageDeltaEvent.Delta.StopSequence
		messageStreamResponse.Usage.OutputTokens = messageDeltaEvent.Usage.OutputTokens
		// set on both levels: no reason for these to ever be out of sync
		messageStreamResponse.Message.Usage.OutputTokens = messageDeltaEvent.Usage.OutputTokens
	case MessageEventTypeMessageStop:
		messageStopEvent := &MessageStopEvent{}
		err = json.Unmarshal([]byte(event), &messageStopEvent)

		messageStreamResponse.Type = messageStopEvent.Type
	case MessageEventTypeError:
		messageErrorEvent := &MessageErrorEvent{}
		err = json.Unmarshal([]byte(event), &messageErrorEvent)
		if err != nil {
			return messageStreamResponse, err
		}

		// error received on stream
		return messageStreamResponse, fmt.Errorf(
			"error type: %s, message: %s",
			messageErrorEvent.Error.Type,
			messageErrorEvent.Error.Message,
		)
	default:
		err = UnsupportedEventType{Msg: "unknown event type"}
	}

	return messageStreamResponse, err
}
