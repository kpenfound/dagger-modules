package main

import (
	"context"
	"fmt"

	"github.com/kpenfound/dagger-modules/encircle/internal/circle"
)

type Executor struct {
	ctx    context.Context
	client *Client
}

func newExecutor(ctx context.Context) *Executor {
	return &Executor{
		ctx:    ctx,
		client: dag,
	}
}

func setup(ctx context.Context) (*circle.Config, *Executor, error) {
	executor := newExecutor(ctx)
	cfg, err := circle.ReadConfig(CONFIG)

	return cfg, executor, err
}

func (e *Executor) executeJob(name string, job *circle.Job) error {
	fmt.Printf("running job %s\n", name)
	src := e.client.Host().Directory(".")
	runner := e.client.Container().
		Pipeline(name).
		From(job.docker[0].image).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithNewFile("/envfile", ContainerWithNewFileOpts{
			Permissions: 0777,
			Contents:    "",
		}).
		WithEnvVariable("BASH_ENV", "/envfile")

	for _, s := range job.steps {
		runner = s.toDagger(runner, map[string]string{})
	}
	_, err := runner.Sync(e.ctx)
	return err
}

func (e *Executor) executeWorkflow(name string, workflow *circle.Workflow, jobs map[string]*circle.Job) error {
	fmt.Printf("running workflow %s\n", name)
	for _, jobName := range workflow.jobs {
		err := e.executeJob(jobName, jobs[jobName])
		if err != nil {
			return err
		}
	}
	return nil
}
