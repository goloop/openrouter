package openrouter_test

import (
	"encoding/json"
	"fmt"

	"github.com/goloop/ai"
	"github.com/goloop/openrouter"
)

func ExampleNew() {
	c := openrouter.New("sk-or-...")
	_ = c // use c.Generate, c.Stream, c.ChatCompletion, ...
	fmt.Println(openrouter.ModelClaudeSonnet)
	// Output: anthropic/claude-sonnet-4
}

// ExampleClient_Generate builds a request. Sending it needs a real API key, so
// this example only shows the shape.
func ExampleClient_Generate() {
	req := &ai.Request{
		Model: openrouter.ModelGPT4o,
		Messages: []ai.Message{
			ai.UserText("Name the capital of France."),
		},
	}
	fmt.Println(req.Model, len(req.Messages))
	// Output: openai/gpt-4o 1
}

// ExampleClient_GetModel notes single-model lookup. OpenRouter has no per-model
// endpoint, so GetModel finds the model in the full list and reports a 404
// ai.APIError when the ID is not routed.
func ExampleClient_GetModel() {
	c := openrouter.New("sk-or-...")
	_ = c // m, _ := c.GetModel(ctx, openrouter.ModelClaudeSonnet)
	fmt.Println(openrouter.ModelClaudeSonnet)
	// Output: anthropic/claude-sonnet-4
}

// ExampleTool shows a tool definition passed with a request.
func ExampleTool() {
	tool := ai.Tool{
		Name:        "get_weather",
		Description: "Get the current weather for a city.",
		Schema:      json.RawMessage(`{"type":"object","properties":{"city":{"type":"string"}}}`),
	}
	fmt.Println(tool.Name)
	// Output: get_weather
}
