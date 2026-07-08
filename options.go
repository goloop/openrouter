package openrouter

import (
	"net/http"
	"time"

	"github.com/goloop/ai"
)

type settings struct {
	aiOpts  []ai.Option
	referer string
	title   string
}

// Option configures a Client in New.
type Option func(*settings)

// WithBaseURL overrides the API base URL (proxies, gateways, mock servers).
func WithBaseURL(u string) Option {
	return func(s *settings) { s.aiOpts = append(s.aiOpts, ai.WithBaseURL(u)) }
}

// WithHTTPClient sets the HTTP client used for requests.
func WithHTTPClient(c *http.Client) Option {
	return func(s *settings) { s.aiOpts = append(s.aiOpts, ai.WithHTTPClient(c)) }
}

// WithTimeout sets the per-request timeout when no custom HTTP client is set.
func WithTimeout(d time.Duration) Option {
	return func(s *settings) { s.aiOpts = append(s.aiOpts, ai.WithTimeout(d)) }
}

// WithMaxRetries sets how many times a request is retried on 429 and 5xx.
func WithMaxRetries(n int) Option {
	return func(s *settings) { s.aiOpts = append(s.aiOpts, ai.WithMaxRetries(n)) }
}

// WithHeader adds a header sent with every request.
func WithHeader(key, value string) Option {
	return func(s *settings) { s.aiOpts = append(s.aiOpts, ai.WithHeader(key, value)) }
}

// WithReferer sets the HTTP-Referer header used for app attribution and
// ranking on openrouter.ai.
func WithReferer(url string) Option {
	return func(s *settings) { s.referer = url }
}

// WithTitle sets the X-Title header used for app attribution and ranking on
// openrouter.ai.
func WithTitle(title string) Option {
	return func(s *settings) { s.title = title }
}
