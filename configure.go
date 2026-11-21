package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/rabidaudio/inshrimpentation/database"
	"github.com/rabidaudio/inshrimpentation/tasks"
	"go.yaml.in/yaml/v3"
)

type ConfigFile struct {
	Tasks map[tasks.Task]TaskConfig `yaml:"tasks"`
}

type TaskConfig struct {
	Interval time.Duration  `yaml:"interval"`
	Enabled  *bool          `yaml:"enabled,omitempty"`
	Options  map[string]any `yaml:"options,omitempty"`
}

func configure(ctx context.Context, confpath string) error {
	db := database.GetDB(ctx)

	f, err := os.Open(confpath)
	if err != nil {
		return err
	}
	defer f.Close()

	var config ConfigFile
	d := yaml.NewDecoder(f)
	if err = d.Decode(&config); err != nil {
		return err
	}

	return db.InTx(ctx, nil, func(db *database.DB) error {
		for _, t := range tasks.AllTasks() {
			tc, ok := config.Tasks[t]
			if !ok {
				// if it wasn't defined in the configuration file, disable it
				if err = db.SetConfigEnabled(ctx, string(t), false); err != nil {
					return err
				}
				continue
			}
			if err = db.SetInterval(ctx, string(t), tc.Interval); err != nil {
				return err
			}
			enabled := tc.Enabled == nil || *tc.Enabled
			if err = db.SetConfigEnabled(ctx, string(t), enabled); err != nil {
				return err
			}
			if err = db.SetConfigOptions(ctx, string(t), tc.Options); err != nil {
				return err
			}
		}

		confs, err := db.GetConfigs(ctx)
		if err != nil {
			return err
		}
		for _, c := range confs {
			log.Printf("[%v] every:%v enabled:%v opts:%+v\n", c.Type, c.Interval.String(), c.Enabled, c.Options)
		}

		return nil
	})
}
