package openrouter

import (
	"fmt"

	"github.com/goloop/ai"
)

// checkHosted refuses a request that asks the provider to run a capability on
// its own side.
//
// This provider is a router: what a request can do depends on the model behind
// the route, and its search is a routing feature rather than a tool the model
// is offered. There is no answer this driver can give for the provider as a
// whole, and answering per route would be a capability table that is wrong the
// week a route changes. A later version adds it as an explicit opt-in.
//
// Returning [ai.ErrNoHosted] before the request leaves is the documented
// behavior, not a placeholder: an answer produced without the search that was
// asked for looks exactly like one produced with it, so failing loudly is the
// only way a caller can tell the difference.
func checkHosted(req *ai.Request) error {
	if len(req.Hosted) == 0 {
		return nil
	}
	return fmt.Errorf("%w: a hosted search depends on the route and is not offered here", ai.ErrNoHosted)
}
