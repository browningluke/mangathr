package hooks

import (
	"fmt"
	"os"
	"os/exec"
)

// SubcommandHookConfig is the config for a subcommand hook.
type SubcommandHookConfig struct {
	Name         string            `yaml:"name"`
	Command      string            `yaml:"command"`
	AbortOnError bool              `yaml:"abortOnError"`
	FireIfEmpty  bool              `yaml:"fireIfEmpty"`
	On           []string          `yaml:"on"`
	Aggregate    bool              `yaml:"aggregate"`
	Args         []string          `yaml:"args"` // each entry is a text/template
	Env          map[string]string `yaml:"env"`  // values are text/templates
}

type SubcommandHook struct {
	baseHook
	cfg SubcommandHookConfig
}

func newSubcommandHook(cfg SubcommandHookConfig) *SubcommandHook {
	return &SubcommandHook{
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

func (h *SubcommandHook) Fire(ctx HookContext) error {
	return h.run(ctx)
}

func (h *SubcommandHook) FireAggregate(ctx AggregateHookContext) error {
	return h.run(ctx)
}

func (h *SubcommandHook) run(data any) error {
	args, err := h.renderArgs(data)
	if err != nil {
		return fmt.Errorf("hooks: subcommand %q: render args: %w", h.cfg.Name, err)
	}

	env, err := h.renderEnv(data)
	if err != nil {
		return fmt.Errorf("hooks: subcommand %q: render env: %w", h.cfg.Name, err)
	}

	cmd := exec.Command(h.cfg.Command, args...)
	cmd.Env = append(os.Environ(), env...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("hooks: subcommand %q failed: %w\nOutput: %s",
			h.cfg.Name, err, output)
	}
	return nil
}

func (h *SubcommandHook) renderArgs(data any) ([]string, error) {
	args := make([]string, 0, len(h.cfg.Args))
	for _, tmpl := range h.cfg.Args {
		rendered, err := renderTemplate(tmpl, data)
		if err != nil {
			return nil, err
		}
		args = append(args, rendered)
	}
	return args, nil
}

func (h *SubcommandHook) renderEnv(data any) ([]string, error) {
	env := make([]string, 0, len(h.cfg.Env))
	for k, tmpl := range h.cfg.Env {
		rendered, err := renderTemplate(tmpl, data)
		if err != nil {
			return nil, err
		}
		env = append(env, k+"="+rendered)
	}
	return env, nil
}
