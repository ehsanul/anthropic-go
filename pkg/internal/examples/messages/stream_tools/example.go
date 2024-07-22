package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ehsanul/anthropic-go/v3/pkg/anthropic"
	"github.com/ehsanul/anthropic-go/v3/pkg/anthropic/client/native"
)

type WeatherRequest struct {
	City string `json:"city" jsonschema:"required,description=city to get the weather for"`
	Unit string `json:"unit" jsonschema:"enum=celsius,enum=fahrenheit,description=temperature unit to return"`
}

func main() {
	ctx := context.Background()
	client, err := native.MakeClient(native.Config{
		APIKey: "your-api-key",
	})
	if err != nil {
		panic(err)
	}

	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude35Sonnet,
		MaxTokensToSample: 1024,
		ToolChoice: &anthropic.ToolChoice{
			Type: "auto", // let the model choose the tool, to demonstrate mixed text and tool use content blocks in the stream
		},
		Tools: []anthropic.Tool{
			{
				Name:        "get_weather",
				Description: "Get the weather",
				InputSchema: anthropic.GenerateInputSchema(&WeatherRequest{}),
			},
		},
		Messages: []anthropic.MessagePartRequest{
			{
				Role: "user",
				Content: []anthropic.ContentBlock{
					anthropic.NewTextContentBlock("How are you? Also, what's the weather in Charleston in Fahrenheit?"),
				},
			},
		},
		Stream: true,
	}

	rCh, errCh := client.MessageStream(ctx, request)

	chunk := &anthropic.MessageStreamResponse{}
	done := false
	messageResponse := anthropic.MessageResponse{}
	var currentContentBlock anthropic.ContentBlock
	currentBuilder := strings.Builder{}

	for {
		select {
		case chunk = <-rCh:
			var out string
			switch chunk.Type {
			case string(anthropic.MessageEventTypePing):
				continue
			case string(anthropic.MessageEventTypeMessageStart):
				messageResponse = chunk.Message
			case string(anthropic.MessageEventTypeMessageDelta):
				if chunk.Delta.StopReason != "" {
					messageResponse.StopReason = chunk.Delta.StopReason
				}
				if chunk.Delta.StopSequence != "" {
					messageResponse.StopSequence = chunk.Delta.StopSequence
				}
				if chunk.Usage.OutputTokens != 0 {
					messageResponse.Usage.OutputTokens = chunk.Usage.OutputTokens
				}
			case string(anthropic.MessageEventTypeMessageStop):
				done = true
			case string(anthropic.MessageEventTypeContentBlockStart):
				currentContentBlock = chunk.ContentBlock
				if chunk.ContentBlock.ContentBlockType() == "tool_use" {
					out = fmt.Sprintf("\nTool Use: %s\n", chunk.ContentBlock.(anthropic.ToolUseContentBlock).Name)
				} else {
					out = "\n"
				}
			case string(anthropic.MessageEventTypeContentBlockStop):
				partResponse := toPartResponse(currentContentBlock, currentBuilder.String())
				messageResponse.Content = append(messageResponse.Content, partResponse)
				currentContentBlock = nil
				currentBuilder.Reset()
			case string(anthropic.MessageEventTypeContentBlockDelta):
				switch chunk.Delta.Type {
				case "text_delta":
					out = chunk.Delta.Text
					currentBuilder.WriteString(chunk.Delta.Text)
				case "input_json_delta":
					out = chunk.Delta.PartialJson
					currentBuilder.WriteString(chunk.Delta.PartialJson)
				}
			default:
				panic("unexpected event type: " + chunk.Type)
			}
			fmt.Print(out)
		case err := <-errCh:
			fmt.Printf("\n\nError: %s\n\n", err)
			done = true
		}

		if chunk.Type == "message_stop" || done {
			break
		}
	}

	fmt.Println()
	fmt.Println("-------------------FINAL RESULT----------------------")
	jsonResponse, err := json.MarshalIndent(messageResponse, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonResponse))
	fmt.Println("-----------------------------------------------------")
}

func toPartResponse(contentBlock anthropic.ContentBlock, content string) anthropic.MessagePartResponse {
	switch contentBlock.ContentBlockType() {
	case "text":
		return anthropic.MessagePartResponse{
			Type: contentBlock.ContentBlockType(),
			Text: content,
		}
	case "tool_use":
		var input map[string]interface{}
		if err := json.Unmarshal([]byte(content), &input); err != nil {
			panic(err)
		}

		toolUseContentBlock := contentBlock.(anthropic.ToolUseContentBlock)
		return anthropic.MessagePartResponse{
			Type:  toolUseContentBlock.Type,
			ID:    toolUseContentBlock.ID,
			Name:  toolUseContentBlock.Name,
			Input: input,
		}
	default:
		panic("unexpected content block type: " + contentBlock.ContentBlockType())
	}
}
