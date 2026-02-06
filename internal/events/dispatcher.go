package events

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/CABGenOrg/cabgen_backend/internal/models"
	"github.com/CABGenOrg/cabgen_backend/internal/repositories"
)

type Dispatcher interface {
	Run(ctx context.Context)
}

type result struct {
	eventID uint
	err     error
}

type dispatcher struct {
	eventRepo repositories.EventRepository
	registry  Registry
	interval  time.Duration
	nWorkers  int
	jobCh     chan models.Event
	resultCh  chan result
}

func NewDispatcher(eventRepo repositories.EventRepository,
	registry Registry, interval time.Duration, nWorkers int) Dispatcher {
	return &dispatcher{
		eventRepo: eventRepo,
		registry:  registry,
		interval:  interval,
		nWorkers:  nWorkers,
		jobCh:     make(chan models.Event, 100),
		resultCh:  make(chan result, 100),
	}
}

func (d *dispatcher) poll(ctx context.Context) {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			events, err := d.eventRepo.GetEvents(ctx, 10)
			if err != nil {
				continue
			}

			for _, ev := range events {
				d.jobCh <- ev
			}
		}
	}
}

func (d *dispatcher) worker(ctx context.Context) {
	for ev := range d.jobCh {
		func() {
			defer func() {
				if r := recover(); r != nil {
					d.resultCh <- result{
						eventID: ev.ID,
						err:     fmt.Errorf("panic in handler: %v", r),
					}
				}
			}()

			handler, ok := d.registry.Get(ev.Name)
			if !ok {
				d.resultCh <- result{
					eventID: ev.ID,
					err:     fmt.Errorf("handler not found"),
				}
				return
			}

			err := handler(ctx, ev.Payload)
			d.resultCh <- result{eventID: ev.ID, err: err}
		}()
	}
}

func (d *dispatcher) listenResults(ctx context.Context) {
	for res := range d.resultCh {
		if res.err != nil {
			d.eventRepo.MarkFailed(ctx, res.eventID, res.err.Error())
			continue
		}
		d.eventRepo.MarkDone(ctx, res.eventID)
	}
}

func (d *dispatcher) Run(ctx context.Context) {
	var wgWriter, wgWorkers sync.WaitGroup

	// DB Writer
	wgWriter.Add(1)
	go func() {
		defer wgWriter.Done()
		d.listenResults(ctx)
	}()

	// Worker Pool
	for range d.nWorkers {
		wgWorkers.Add(1)
		go func() {
			defer wgWorkers.Done()
			d.worker(ctx)
		}()
	}

	// Poller
	go d.poll(ctx)

	// Wait context to be canceled
	<-ctx.Done()

	close(d.jobCh)
	wgWorkers.Wait()
	close(d.resultCh)
	wgWriter.Wait()
}
