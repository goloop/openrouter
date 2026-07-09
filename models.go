package openrouter

import (
	"context"
	"net/http"

	"github.com/goloop/ai"
)

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

// GetModel returns the model with the given ID. OpenRouter has no per-model
// endpoint, so it looks the model up in the full list and reports a 404
// [ai.APIError] when the ID is not routed.
func (c *Client) GetModel(ctx context.Context, id string) (*Model, error) {
	models, err := c.Models(ctx)
	if err != nil {
		return nil, err
	}
	for i := range models {
		if models[i].ID == id {
			return &models[i], nil
		}
	}
	return nil, &ai.APIError{
		Status:  http.StatusNotFound,
		Message: "model not found: " + id,
	}
}
