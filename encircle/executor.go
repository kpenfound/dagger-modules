package main

import (
	"context"
	"fmt"
)

type Executor struct {
	Ctx    context.Context
	Client *Client
}

func NewExecutor(ctx context.Context) *Executor {
	return &Executor{
		Ctx:    ctx,
		Client: dag,
	}
}

func (e *Executor) ExecuteJob(name string, job *Job) error {
	fmt.Printf("running job %s\n", name)
	src := e.Client.Host().Directory(".")
	runner := e.Client.Container().
		Pipeline(name).
		From(job.Docker[0].Image).
		WithMountedDirectory("/src", src).
		WithWorkdir("/src").
		WithNewFile("/envfile", ContainerWithNewFileOpts{
			Permissions: 0777,
			Contents:    "",
		}).
		WithEnvVariable("BASH_ENV", "/envfile")

	for _, s := range job.Steps {
		runner = s.ToDagger(runner, map[string]string{})
	}
	_, err := runner.Sync(e.Ctx)
	return err
}

func (e *Executor) ExecuteWorkflow(name string, workflow *Workflow, jobs map[string]*Job) error {
	fmt.Printf("running workflow %s\n", name)
	for _, jobName := range workflow.Jobs {
		err := e.ExecuteJob(jobName, jobs[jobName])
		if err != nil {
			return err
		}
	}
	return nil
}
