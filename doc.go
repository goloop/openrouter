// Package openrouter is a client for the OpenRouter API, built on the goloop/ai
// interface.
//
// OpenRouter is a routing gateway: one API and key reach many model providers,
// with models namespaced as "provider/model". The Client implements ai.Client,
// so Generate and Stream work the same as with any other goloop AI provider.
// On top of that it exposes the native chat completions endpoint with its full
// options, model listing and single-model lookup. The wire format is
// chat-completions compatible.
//
//	c := openrouter.New(os.Getenv("OPENROUTER_API_KEY"))
//	resp, err := c.Generate(ctx, &ai.Request{
//	    Model:    openrouter.ModelClaudeSonnet,
//	    Messages: []ai.Message{ai.UserText("Say hello in one word.")},
//	})
//
// WithReferer and WithTitle set the app-attribution headers OpenRouter uses for
// ranking. It depends only on goloop/ai and the standard library.
package openrouter
