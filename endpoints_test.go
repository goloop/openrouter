package openrouter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/goloop/ai"
)

func TestChatCompletionNative(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var req ChatRequest
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &req)
		if string(req.ResponseFormat) != `{"type":"json_object"}` {
			t.Errorf("response_format = %s", req.ResponseFormat)
		}
		io.WriteString(w, `{"model":"m","choices":[{"index":0,`+
			`"message":{"role":"assistant","content":"{}"},"finish_reason":"stop"}]}`)
	})
	defer done()

	resp, err := c.ChatCompletion(context.Background(), &ChatRequest{
		Model:          "m",
		Messages:       []ChatMessage{{Role: "user", Content: "hi"}},
		ResponseFormat: json.RawMessage(`{"type":"json_object"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Choices[0].Message.Content != "{}" {
		t.Errorf("content = %v", resp.Choices[0].Message.Content)
	}
}

func TestAttributionHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HTTP-Referer") != "https://app.example" {
			t.Errorf("referer = %q", r.Header.Get("HTTP-Referer"))
		}
		if r.Header.Get("X-Title") != "My App" {
			t.Errorf("title = %q", r.Header.Get("X-Title"))
		}
		io.WriteString(w, `{"model":"m","choices":[{"index":0,`+
			`"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}]}`)
	}))
	defer srv.Close()

	c := New("key", WithBaseURL(srv.URL), WithMaxRetries(0),
		WithReferer("https://app.example"), WithTitle("My App"))
	if _, err := c.Generate(context.Background(), &ai.Request{
		Model: "m", Messages: []ai.Message{ai.UserText("hi")},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestModels(t *testing.T) {
	c, done := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, `{"data":[{"id":"openai/gpt-4o","name":"GPT-4o","context_length":128000}]}`)
	})
	defer done()

	models, err := c.Models(context.Background())
	if err != nil || len(models) != 1 || models[0].ID != "openai/gpt-4o" {
		t.Fatalf("models: %v %+v", err, models)
	}
	if models[0].ContextLength != 128000 {
		t.Errorf("context length = %d", models[0].ContextLength)
	}
}
