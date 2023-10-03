package main

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kpenfound/dagger-modules/encircle/circle"
	"golang.org/x/exp/maps"
)

type Executor struct {
	ctx    context.Context
	client *Client
	dir    *Directory
	logs   []string
}

func newExecutor(ctx context.Context, dir *Directory) *Executor {
	return &Executor{
		ctx:    ctx,
		client: dag,
		dir:    dir,
	}
}

func setup(ctx context.Context, dir *Directory) (*circle.Config, *Executor, error) {
	executor := newExecutor(ctx, dir)
	cfgFile, err := dir.File(CONFIG).Contents(ctx)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := circle.ParseConfig(cfgFile)

	return cfg, executor, err
}

func (e *Executor) executeJob(name string, job *circle.Job) error {
	fmt.Printf("running job %s\n", name)
	runner := e.client.Container().
		Pipeline(name).
		From(job.Docker[0].Image).
		WithMountedDirectory("/src", e.dir).
		WithWorkdir("/src").
		WithNewFile("/envfile", ContainerWithNewFileOpts{
			Permissions: 0777,
			Contents:    "",
		}).
		WithEnvVariable("BASH_ENV", "/envfile")

	for _, s := range job.Steps {
		runner = StepToDagger(s, runner, map[string]string{})
	}
	out, err := runner.Stdout(e.ctx)
	if err != nil {
		return err
	}
	e.logs = append(e.logs, out)
	return nil
}

func (e *Executor) executeWorkflow(name string, workflow *circle.Workflow, jobs map[string]*circle.Job) error {
	fmt.Printf("running workflow %s\n", name)
	for _, jobName := range workflow.Jobs {
		err := e.executeJob(jobName, jobs[jobName])
		if err != nil {
			return err
		}
	}
	return nil
}

func OrbCommandToDagger(oc *circle.OrbCommand, c *Container, params map[string]string) *Container {
	// TODO: handle params?
	for _, s := range oc.Steps {
		c = StepToDagger(s, c, params)
	}
	return c
}

func StepToDagger(s *circle.Step, c *Container, params map[string]string) *Container {
	c = c.Pipeline(s.Name)
	if s.WorkDir != "" { // workdir relative to project root
		c = c.WithWorkdir(filepath.Join("/src", s.WorkDir))
	}
	if s.Run != nil {
		c = RunToDagger(s.Run, c, params)
	}
	if s.Command != nil {
		// Get default params
		maps.Copy(params, s.Command.OrbCommand.GetDefaultParameters())
		// Override user params
		maps.Copy(params, s.Command.Parameters)
		c = OrbCommandToDagger(s.Command.OrbCommand, c, params)
	}

	return c
}

func RunToDagger(r *circle.Run, c *Container, params map[string]string) *Container {
	c = c.Pipeline(replaceParams(r.Name, params))
	// Set env vars
	for k, v := range r.Environment {
		c = c.WithEnvVariable(k, replaceParams(v, params))
	}
	// Exec command
	command := replaceParams(r.Command, params)
	command = fmt.Sprintf("#!/bin/bash\n%s", command)
	script := fmt.Sprintf("/%s.sh", getSha(command))
	c = c.WithNewFile(script, ContainerWithNewFileOpts{
		Permissions: 0777,
		Contents:    command,
	})
	c = c.WithExec([]string{script})
	// Unset env vars
	for k := range r.Environment {
		c = c.WithoutEnvVariable(k)
	}
	return c
}

func replaceParams(target string, params map[string]string) string {
	if strings.Contains(target, "<< parameters.") || strings.Contains(target, "<<parameters.") {
		for k, v := range params {
			p1 := fmt.Sprintf("<< parameters.%s >>", k)
			p2 := fmt.Sprintf("<<parameters.%s>>", k)
			target = strings.ReplaceAll(target, p1, v)
			target = strings.ReplaceAll(target, p2, v)
		}
	}
	return target
}

func getSha(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
