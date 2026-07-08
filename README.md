[![deps.dev](https://img.shields.io/badge/deps.dev-insights-4c8dbc)](https://deps.dev/go/github.com%2Fgoloop%2Fopenrouter) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/openrouter/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/openrouter) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# openrouter

`openrouter` is a Go client for the OpenRouter API. It implements the
`github.com/goloop/ai` interface, so it looks and works like every other goloop
AI provider. OpenRouter is a routing gateway: one API and key reach many model
providers, with models namespaced as `provider/model`.

## Features

- Chat completions: `Generate` for a single response, `Stream` for
  token-by-token output through `iter.Seq2`.
- Tool use (function calling), multimodal image input and system prompts.
- Native `ChatCompletion` and `ChatCompletionStream` with the full option set.
- Model listing across all routed providers.
- App-attribution headers via `WithReferer` and `WithTitle`.
- Retries on 429 and 5xx with backoff; normalized, typed API errors.
- Depends only on `github.com/goloop/ai` and the standard library.

## Installation

```sh
go get github.com/goloop/openrouter
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/goloop/ai"
	"github.com/goloop/openrouter"
)

func main() {
	c := openrouter.New(os.Getenv("OPENROUTER_API_KEY"),
		openrouter.WithReferer("https://myapp.example"),
		openrouter.WithTitle("My App"),
	)

	resp, err := c.Generate(context.Background(), &ai.Request{
		Model:    openrouter.ModelClaudeSonnet,
		Messages: []ai.Message{ai.UserText("Say hello in one word.")},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Text())
}
```

## Streaming

```go
for chunk, err := range c.Stream(ctx, req) {
	if err != nil {
		break
	}
	fmt.Print(chunk.Text)
}
```

## Choosing a model

Any `provider/model` string works; a few are provided as constants
(`ModelGPT4o`, `ModelClaudeSonnet`, `ModelGeminiFlash`, ...). List everything
available with `c.Models(ctx)`.

## Documentation

Full reference: **[DOC.md](DOC.md)** (Ukrainian: **[DOC.UK.md](DOC.UK.md)**).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT - see [LICENSE](LICENSE).
