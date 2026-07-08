# openrouter - reference

The full reference for the `openrouter` package: the client, the shared
`goloop/ai` model, chat completions (interface and native), streaming, app
attribution and models.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Creating a client](#creating-a-client)
- [Generate and Stream](#generate-and-stream)
- [Native chat completions](#native-chat-completions)
- [App attribution](#app-attribution)
- [Tools, images and system prompts](#tools-images-and-system-prompts)
- [Models](#models)
- [Options and errors](#options-and-errors)

## Mental model

`openrouter.Client` implements `ai.Client`, the provider-agnostic contract from
`github.com/goloop/ai`. OpenRouter is a routing gateway - one API and key reach
many model providers, with models namespaced as `provider/model`. The shared
`Generate` and `Stream` cover the common ground - chat with tools, images and
streaming - so code written against the interface runs on any provider.

The wire format is chat-completions compatible; native power lives in
`ChatCompletion` and model listing.

```go
import (
	"github.com/goloop/ai"
	"github.com/goloop/openrouter"
)
```

## Creating a client

```go
c := openrouter.New(os.Getenv("OPENROUTER_API_KEY"))

c = openrouter.New(apiKey,
	openrouter.WithReferer("https://myapp.example"),
	openrouter.WithTitle("My App"),
)
```

The base URL defaults to `https://openrouter.ai/api/v1`.

## Generate and Stream

```go
resp, err := c.Generate(ctx, &ai.Request{
	Model:    openrouter.ModelClaudeSonnet,
	System:   "You are concise.",
	Messages: []ai.Message{ai.UserText("Name three primary colors.")},
})
resp.Text()
resp.ToolCalls()
resp.Usage
```

`Stream` returns `iter.Seq2[ai.Chunk, error]`: text deltas as chunks with
`Text`, a finished tool call as a chunk with `ToolCall`, and a final chunk with
`Done` and `Usage`.

## Native chat completions

For provider-only options build a `ChatRequest` and call `ChatCompletion` or
`ChatCompletionStream`:

```go
resp, err := c.ChatCompletion(ctx, &openrouter.ChatRequest{
	Model:          openrouter.ModelGPT4o,
	Messages:       []openrouter.ChatMessage{{Role: "user", Content: "as JSON"}},
	ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
})
```

## App attribution

`WithReferer` and `WithTitle` set the `HTTP-Referer` and `X-Title` headers that
OpenRouter uses to attribute and rank traffic per application. Both are
optional.

## Tools, images and system prompts

Tool use, images and system prompts use the shared `ai` types: `ai.Tool`,
`ai.Image`, `ai.ToolResult` and a `RoleSystem` message or the `System` field.
Tool results are sent back as `RoleTool` messages whose `ai.ToolResult.ID`
matches the `ai.ToolUse.ID`. Inline image bytes are sent as a base64 data URI.

## Models

Any `provider/model` string works; constants such as `ModelGPT4o`,
`ModelClaudeSonnet` and `ModelGeminiFlash` are provided for convenience.

```go
models, err := c.Models(ctx)
models[0].ID            // "openai/gpt-4o"
models[0].ContextLength
```

## Options and errors

Options: `WithBaseURL`, `WithHTTPClient`, `WithTimeout`, `WithMaxRetries`,
`WithHeader`, `WithReferer`, `WithTitle`.

A non-success response becomes an `*ai.APIError` with `Status`, `Type`, `Code`,
`Message` and the raw body:

```go
var apiErr *ai.APIError
if errors.As(err, &apiErr) && apiErr.Status == http.StatusTooManyRequests {
	// back off
}
```

Requests missing a model or messages fail before the network with
`ai.ErrNoModel` or `ai.ErrNoMessages`.
