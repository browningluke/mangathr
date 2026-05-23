package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// WebhookHookConfig is the config for a generic HTTP webhook hook.
type WebhookHookConfig struct {
	Name         string            `yaml:"name"`
	WebhookURL   string            `yaml:"webhookURL"`
	RequestType  string            `yaml:"requestType"`  // GET/POST/PUT/PATCH; default POST
	SuccessCode  int               `yaml:"successCode"`  // default 200
	AbortOnError bool              `yaml:"abortOnError"`
	FireIfEmpty  bool              `yaml:"fireIfEmpty"`
	On           []string          `yaml:"on"`
	Aggregate    bool              `yaml:"aggregate"`
	Body         string            `yaml:"body"`    // text/template → JSON body
	Headers      string            `yaml:"headers"` // text/template → JSON map of headers
}

type WebhookHook struct {
	baseHook
	cfg WebhookHookConfig
}

func newWebhookHook(cfg WebhookHookConfig) *WebhookHook {
	if cfg.RequestType == "" {
		cfg.RequestType = "POST"
	}
	if cfg.SuccessCode == 0 {
		cfg.SuccessCode = 200
	}
	return &WebhookHook{
		baseHook: baseHook{
			name:         cfg.Name,
			events:       cfg.On,
			abortOnError: cfg.AbortOnError,
			fireIfEmpty:  cfg.FireIfEmpty,
			aggregate:    cfg.Aggregate,
		},
		cfg: cfg,
	}
}

func (h *WebhookHook) Fire(ctx HookContext) error {
	return h.send(ctx)
}

func (h *WebhookHook) FireAggregate(ctx AggregateHookContext) error {
	return h.send(ctx)
}

func (h *WebhookHook) send(data any) error {
	body, err := renderTemplate(h.cfg.Body, data)
	if err != nil {
		return fmt.Errorf("hooks: webhook %q: render body: %w", h.cfg.Name, err)
	}

	headers, err := h.renderHeaders(data)
	if err != nil {
		return fmt.Errorf("hooks: webhook %q: render headers: %w", h.cfg.Name, err)
	}

	method := strings.ToUpper(h.cfg.RequestType)
	req, err := http.NewRequest(method, h.cfg.WebhookURL, bytes.NewBufferString(body))
	if err != nil {
		return fmt.Errorf("hooks: webhook %q: build request: %w", h.cfg.Name, err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("hooks: webhook %q: request failed: %w", h.cfg.Name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != h.cfg.SuccessCode {
		return fmt.Errorf("hooks: webhook %q: unexpected status %d (want %d)",
			h.cfg.Name, resp.StatusCode, h.cfg.SuccessCode)
	}
	return nil
}

func (h *WebhookHook) renderHeaders(data any) (map[string]string, error) {
	if h.cfg.Headers == "" {
		return nil, nil
	}
	rendered, err := renderTemplate(h.cfg.Headers, data)
	if err != nil {
		return nil, err
	}
	var result map[string]string
	if err := json.Unmarshal([]byte(rendered), &result); err != nil {
		return nil, fmt.Errorf("parse rendered headers as JSON: %w", err)
	}
	return result, nil
}
