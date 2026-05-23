package hooks

import (
	"fmt"

	"github.com/browningluke/mangathr/v2/internal/logging"
)

// HooksConfig is the top-level config struct for all hook types, embedded
// in the global Config and marshalled from the "hooks" YAML key.
type HooksConfig struct {
	Discord    []DiscordHookConfig    `yaml:"discord"`
	Webhook    []WebhookHookConfig    `yaml:"webhook"`
	Subcommand []SubcommandHookConfig `yaml:"subcommand"`
}

// controller is the package-level singleton, initialised via SetConfig.
var controller *Controller

// SetConfig initialises the package-level controller from cfg.
func SetConfig(cfg HooksConfig) {
	controller = newController(cfg)
}

// Fire dispatches event to all registered hooks. When abortOnError is set on a
// failing hook, the error is returned and the caller should abort the operation.
func Fire(event string, ctx HookContext) error {
	if controller == nil {
		return nil
	}
	return controller.Fire(event, ctx)
}

// FinalizeAggregates flushes all pending aggregate buffers, firing each hook
// once with its collected contexts. Call at the end of each command run.
func FinalizeAggregates() error {
	if controller == nil {
		return nil
	}
	return controller.FinalizeAggregates()
}

// Controller manages hook dispatch and aggregate buffering.
type Controller struct {
	hooks   []Hook
	buffers map[string][]HookContext // keyed by hook name
}

func newController(cfg HooksConfig) *Controller {
	var hooks []Hook
	for _, d := range cfg.Discord {
		hooks = append(hooks, newDiscordHook(d))
	}
	for _, w := range cfg.Webhook {
		hooks = append(hooks, newWebhookHook(w))
	}
	for _, s := range cfg.Subcommand {
		hooks = append(hooks, newSubcommandHook(s))
	}
	return &Controller{
		hooks:   hooks,
		buffers: make(map[string][]HookContext),
	}
}

func (c *Controller) Fire(event string, ctx HookContext) error {
	ctx.Event = event
	isError := ctx.Error != nil
	_ = isError // event name already encodes success/error

	for _, h := range c.hooks {
		if !h.ShouldFire(event) {
			continue
		}
		if h.IsAggregate() {
			c.buffers[h.Name()] = append(c.buffers[h.Name()], ctx)
			continue
		}
		if !h.FireIfEmpty() && ctx.Chapter.Count == 0 {
			continue
		}
		if err := h.Fire(ctx); err != nil {
			if h.AbortOnError() {
				return fmt.Errorf("hook %q: %w", h.Name(), err)
			}
			logging.Warningln(fmt.Sprintf("hook %q failed (abortOnError=false): %v", h.Name(), err))
		}
	}
	return nil
}

func (c *Controller) FinalizeAggregates() error {
	for _, h := range c.hooks {
		if !h.IsAggregate() {
			continue
		}

		buf := c.buffers[h.Name()] // nil/empty if no events were collected

		// If nothing was collected and fireIfEmpty is false, skip entirely.
		if len(buf) == 0 && !h.FireIfEmpty() {
			continue
		}

		aggCtx := buildAggregateContext(buf)

		// If events were collected but total chapters is zero, apply the same gate.
		if !h.FireIfEmpty() && aggCtx.ChapterCount == 0 {
			delete(c.buffers, h.Name())
			continue
		}

		if err := h.FireAggregate(aggCtx); err != nil {
			delete(c.buffers, h.Name())
			if h.AbortOnError() {
				return fmt.Errorf("hook %q: %w", h.Name(), err)
			}
			logging.Warningln(fmt.Sprintf("hook %q failed (abortOnError=false): %v", h.Name(), err))
			continue
		}
		delete(c.buffers, h.Name())
	}
	return nil
}

func buildAggregateContext(items []HookContext) AggregateHookContext {
	var chapterCount, errorCount int
	seen := make(map[string]struct{})
	for _, ctx := range items {
		chapterCount += ctx.Chapter.Count
		if ctx.Error != nil {
			errorCount++
		}
		seen[ctx.Manga.Title] = struct{}{}
	}
	return AggregateHookContext{
		Items:        items,
		ChapterCount: chapterCount,
		ErrorCount:   errorCount,
		MangaCount:   len(seen),
	}
}
