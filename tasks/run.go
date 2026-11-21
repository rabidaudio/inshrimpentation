package tasks

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
	"golang.org/x/sync/errgroup"
)

func RunAllTasks(ctx context.Context, force bool) ([]*database.Measurement, []*database.Alert, error) {
	return RunTasks(ctx, force, AllTasks()...)
}

func RunTasks(ctx context.Context, force bool, tasks ...Task) ([]*database.Measurement, []*database.Alert, error) {
	db := database.GetDB(ctx)

	configs, err := db.GetConfigs(ctx)
	if err != nil {
		return nil, nil, err
	}
	tc := map[Task]*database.Config{}
	for _, c := range configs {
		tc[Task(c.Type)] = c
	}

	// by using buffered channels we reduce the amount of routine switching
	mmc := make(chan *database.Measurement, len(tasks))
	aac := make(chan *database.Alert, len(tasks))

	tasksRun := map[Task]time.Time{}

	cwg, ectx := errgroup.WithContext(ctx)

	// runners
	rwg := sync.WaitGroup{}
	for _, t := range tasks {
		rwg.Go(func() {
			c := tc[t]
			if c == nil {
				log.Printf("tasks: skipping %v, reason: not configured\n", t)
				return
			}
			if !c.Enabled && !force {
				log.Printf("tasks: skipping %v, reason: disabled\n", t)
				return
			}
			if !c.IsDue() {
				if force {
					log.Printf("tasks: running %v, reason: forced\n", t)
				} else {
					log.Printf("tasks: skipping %v, reason: due in: %v\n", t, c.DueIn())
					return
				}
			} else {
				log.Printf("tasks: running %v, reason: due at: %v\n", t, c.DueAt())
			}

			runner := taskRunners[t]
			m, a, err := runner(ectx, c.Options)
			for _, mm := range m {
				mmc <- mm
			}
			if a != nil {
				aac <- a
			}
			// instead of crashing, add an alert for any errors
			if err != nil {
				aac <- &database.Alert{
					Type:    fmt.Sprintf("error:%s", t),
					Time:    time.Now(),
					Message: fmt.Errorf("task failed [%v]: %w", t, err).Error(),
				}
			}
			tasksRun[t] = time.Now()
		})
	}

	var mm []*database.Measurement
	var aa []*database.Alert

	cwg.Go(func() error {
		for m := range mmc {
			mm = append(mm, m)
			err := db.Record(ectx, m)
			if err != nil {
				return err
			}
		}
		return nil
	})
	cwg.Go(func() error {
		for a := range aac {
			aa = append(aa, a)
			err := db.AddAlert(ectx, a)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// when the runners finish, close the channels and record runs
	cwg.Go(func() error {
		rwg.Wait()
		close(mmc)
		close(aac)

		for t, at := range tasksRun {
			if err = db.TrackRun(ectx, string(t), at); err != nil {
				return err
			}
		}
		return nil
	})

	// wait for collectors and runners to complete
	err = cwg.Wait()
	return mm, aa, err
}
