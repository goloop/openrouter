package openrouter

import "context"

// Model describes a model listed by OpenRouter.
type Model struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	ContextLength int    `json:"context_length,omitempty"`
	Pricing       struct {
		Prompt     string `json:"prompt,omitempty"`
		Completion string `json:"completion,omitempty"`
	} `json:"pricing,omitempty"`
}

// Models lists the models OpenRouter can route to.
func (c *Client) Models(ctx context.Context) ([]Model, error) {
	var out struct {
		Data []Model `json:"data"`
	}
	if err := c.getJSON(ctx, "/models", &out); err != nil {
		return nil, err
	}
	return out.Data, nil
}
