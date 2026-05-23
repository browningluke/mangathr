package workerpool

import (
	"fmt"
	"github.com/alitto/pond"
)

type pool struct {
	tasks   []func()
	workers int
}

func New(workers int) *pool {
	return &pool{
		workers: workers,
	}
}

func (p *pool) AddTask(f func()) {
	p.tasks = append(p.tasks, f)
}

// Run submits all tasks to a worker pool and blocks until completion.
// Returns the first worker panic as an error, if any.
func (p *pool) Run() error {
	// Buffered so the panic handler never blocks even if we return early.
	wpErr := make(chan error, 1)

	panicHandler := func(v interface{}) {
		err, ok := v.(error)
		if !ok {
			err = fmt.Errorf("worker panic: %v", v)
		}
		// Non-blocking send: only the first panic is captured.
		select {
		case wpErr <- err:
		default:
		}
	}

	wp := pond.New(p.workers, 0, pond.PanicHandler(panicHandler))

	for _, task := range p.tasks {
		wp.Submit(task)
	}

	// StopAndWait blocks until all submitted tasks have finished (or
	// been dropped after a panic), eliminating the busy-wait spin loop.
	wp.StopAndWait()

	select {
	case err := <-wpErr:
		return err
	default:
		return nil
	}
}
