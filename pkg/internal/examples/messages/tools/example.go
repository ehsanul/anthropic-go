package main

import (
	"context"
	"encoding/json"
	"fmt"

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

	// Prepare a message request
	request := &anthropic.MessageRequest{
		Model:             anthropic.Claude35Sonnet,
		MaxTokensToSample: 1024,
		ToolChoice: &anthropic.ToolChoice{
			// force the use of the tool
			Type: "tool",
			Name: "get_weather",
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
					anthropic.NewTextContentBlock("what's the weather in Charleston?"),
				},
			},
		},
	}

	// Call the Message method
	response, err := client.Message(ctx, request)
	if err != nil {
		panic(err)
	}

	if response.StopReason == "tool_use" {
		// Do something with the tool use response
		responseJson, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Tool use response: %s\n", responseJson)
	}
}
