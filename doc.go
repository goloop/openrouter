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
// # Structured output
//
// ai.Request.Format maps onto the provider's response_format, so a request for
// JSON is enforced rather than merely asked for, and ai.Response.JSON decodes
// the reply. Plain JSON mode also puts ai.Format.Instruction into the system
// prompt, because this wire format rejects json_object unless the word "json"
// appears in the messages.
//
// # Hosted capabilities
//
// This provider routes to models rather than serving one, so what a request
// can do depends on the route behind it, and a hosted search is a routing
// feature rather than a tool the model is offered. ai.Hosted is refused with
// ai.ErrNoHosted rather than answered differently per route.
//
// The refusal is the documented behavior, not a gap waiting to be filled
// in silence: an answer produced without the search that was asked for
// looks exactly like one produced with it. A caller who would rather have
// the answer anyway asks again without ai.Request.Hosted.
//
// # Asking what this driver can do
//
// Capabilities describes this driver for the decision taken before a call:
// whether to offer a feature at all, and whether it needs one request or two.
//
//	if ai.SupportsHosted(c, ai.Hosted{Kind: ai.HostedWebSearch}) { ... }
//
// It is a hint and not a permission - support also depends on the model, the
// account and the region - so ai.ErrNoHosted and ai.ErrNoFormat remain the
// source of truth and a caller still handles them. What changes is that a
// refusal the provider only reports as a 400 now arrives as those same
// sentinels, wrapped around the original ai.APIError, so one errors.Is covers
// a limitation this driver knew in advance and one it learned over the wire.
//
// WithReferer and WithTitle set the app-attribution headers OpenRouter uses for
// ranking. It depends only on goloop/ai and the standard library.
package openrouter
