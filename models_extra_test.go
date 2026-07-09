package openrouter

import (
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/goloop/ai"
)

func TestGetModel(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":[`+
			`{"id":"a/one","name":"One","context_length":8192},`+
			`{"id":"b/two","name":"Two"}]}`)
	})
	defer done()

	m, err := c.GetModel(context.Background(), "b/two")
	if err != nil {
		t.Fatal(err)
	}
	if m.Name != "Two" {
		t.Errorf("name = %q", m.Name)
	}
}

func TestGetModelNotFound(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":[{"id":"a/one","name":"One"}]}`)
	})
	defer done()

	_, err := c.GetModel(context.Background(), "missing")
	var apiErr *ai.APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		t.Fatalf("err = %v", err)
	}
}
