package openrouter

import "github.com/goloop/ai"

// DefaultBaseURL is the OpenRouter API base URL, including the version segment.
const DefaultBaseURL = "https://openrouter.ai/api/v1"

// Convenience model identifiers. OpenRouter routes to many providers; a model
// is namespaced as "provider/model". Any model string is accepted; use Models
// to discover what is available.
const (
	ModelGPT4o        = "openai/gpt-4o"
	ModelClaudeSonnet = "anthropic/claude-sonnet-4"
	ModelGeminiFlash  = "google/gemini-2.5-flash"
	ModelLlama3170B   = "meta-llama/llama-3.1-70b-instruct"
	ModelDeepSeekChat = "deepseek/deepseek-chat"
	ModelMistralLarge = "mistralai/mistral-large"
)

// Client is an OpenRouter API client. It implements [ai.Client] and adds the
// provider's native endpoints. The wire format is chat-completions compatible.
type Client struct {
	opts    ai.Options
	referer string
	title   string
}

var _ ai.Client = (*Client)(nil)

// New returns a Client for the given API key. Shared options (WithBaseURL,
// WithHTTPClient, WithTimeout, WithMaxRetries, WithHeader) and OpenRouter
// options (WithReferer, WithTitle) configure it.
func New(apiKey string, opts ...Option) *Client {
	s := settings{}
	for _, o := range opts {
		o(&s)
	}

	o := ai.NewOptions(apiKey, s.aiOpts...)
	if o.BaseURL == "" {
		o.BaseURL = DefaultBaseURL
	}

	return &Client{opts: o, referer: s.referer, title: s.title}
}
