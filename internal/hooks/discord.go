package hooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// DiscordHookConfig is the config for a Discord webhook hook.
type DiscordHookConfig struct {
	Name         string   `yaml:"name"`
	WebhookURL   string   `yaml:"webhookURL"`
	AbortOnError bool     `yaml:"abortOnError"`
	FireIfEmpty  bool     `yaml:"fireIfEmpty"`
	On           []string `yaml:"on"`
	Aggregate    bool     `yaml:"aggregate"`
	Embed        bool     `yaml:"embed"`
	Template     struct {
		// Used when embed: true
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Color       int    `yaml:"color"`
		Footer      string `yaml:"footer"`
		// Used when embed: false
		Message string `yaml:"message"`
	} `yaml:"template"`
}

type DiscordHook struct {
	baseHook
	cfg DiscordHookConfig
}

func newDiscordHook(cfg DiscordHookConfig) *DiscordHook {
	return &DiscordHook{
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

func (h *DiscordHook) Fire(ctx HookContext) error {
	return h.send(ctx)
}

func (h *DiscordHook) FireAggregate(ctx AggregateHookContext) error {
	return h.send(ctx)
}

func (h *DiscordHook) send(data any) error {
	payload, err := h.buildPayload(data)
	if err != nil {
		return fmt.Errorf("hooks: discord %q: build payload: %w", h.cfg.Name, err)
	}
	return h.post(payload)
}

// discordEmbedPayload is the JSON envelope for embed-style messages.
type discordEmbedPayload struct {
	Embeds []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string         `json:"title,omitempty"`
	Description string         `json:"description,omitempty"`
	Color       int            `json:"color,omitempty"`
	Footer      *discordFooter `json:"footer,omitempty"`
}

type discordFooter struct {
	Text string `json:"text"`
}

type discordMessagePayload struct {
	Content string `json:"content"`
}

func (h *DiscordHook) buildPayload(data any) ([]byte, error) {
	if h.cfg.Embed {
		title, err := renderTemplate(h.cfg.Template.Title, data)
		if err != nil {
			return nil, err
		}
		desc, err := renderTemplate(h.cfg.Template.Description, data)
		if err != nil {
			return nil, err
		}
		footer, err := renderTemplate(h.cfg.Template.Footer, data)
		if err != nil {
			return nil, err
		}

		embed := discordEmbed{
			Title:       title,
			Description: desc,
			Color:       h.cfg.Template.Color,
		}
		if footer != "" {
			embed.Footer = &discordFooter{Text: footer}
		}
		return json.Marshal(discordEmbedPayload{Embeds: []discordEmbed{embed}})
	}

	msg, err := renderTemplate(h.cfg.Template.Message, data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(discordMessagePayload{Content: msg})
}

func (h *DiscordHook) post(payload []byte) error {
	resp, err := http.Post(h.cfg.WebhookURL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("hooks: discord %q: request failed: %w", h.cfg.Name, err)
	}
	defer resp.Body.Close()

	// Discord returns 204 No Content on success.
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("hooks: discord %q: unexpected status %d", h.cfg.Name, resp.StatusCode)
	}
	return nil
}
