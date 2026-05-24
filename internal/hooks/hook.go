package hooks

import "strings"

// Hook is implemented by each hook type (discord, webhook, subcommand).
type Hook interface {
	Name() string
	IsAggregate() bool
	AbortOnError() bool
	FireIfEmpty() bool
	ShouldFire(event string) bool

	Fire(ctx HookContext) error
	FireAggregate(ctx AggregateHookContext) error
}

// baseHook provides the shared fields and method implementations that all
// hook types embed to satisfy the non-send parts of the Hook interface.
type baseHook struct {
	name         string
	events       []string
	abortOnError bool
	fireIfEmpty  bool
	aggregate    bool
}

func (b *baseHook) Name() string        { return b.name }
func (b *baseHook) IsAggregate() bool   { return b.aggregate }
func (b *baseHook) AbortOnError() bool  { return b.abortOnError }
func (b *baseHook) FireIfEmpty() bool   { return b.fireIfEmpty }

func (b *baseHook) ShouldFire(event string) bool {
	for _, e := range b.events {
		if strings.EqualFold(e, event) {
			return true
		}
	}
	return false
}
