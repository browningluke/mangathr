package workerpool

import (
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

// Run a blocking abstraction over Pond's WorkerPool to download all pages
func (p *pool) Run() error {
	wpErr := make(chan error)
	panicHandler := func(p interface{}) {
		wpErr <- p.(error)
	}
	pool := pond.New(p.workers, 0, pond.PanicHandler(panicHandler))

	for _, task := range p.tasks {
		pool.Submit(task)
	}

	for {
		select {
		case err := <-wpErr:
			pool.Stop()
			return err
		default:
			if pool.SubmittedTasks() == pool.CompletedTasks() {
				return nil
			}
		}
	}
}
