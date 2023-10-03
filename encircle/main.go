package main

import (
	"context"
	"strings"
)

const CONFIG = "./.circleci/config.yml"

type Encircle struct{}

func (m *Encircle) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}

func (d *Directory) EncircleJob(ctx context.Context, job string) (string, error) {
	cfg, executor, err := setup(ctx, d)
	if err != nil {
		return "", err
	}

	err = executor.executeJob(job, cfg.Jobs[job])
	if err != nil {
		return "", err
	}

	return strings.Join(executor.logs, "\n"), nil
}

func (d *Directory) EncircleWorkflow(ctx context.Context, workflow string) (string, error) {

	cfg, executor, err := setup(ctx, d)
	if err != nil {
		return "", err
	}

	err = executor.executeWorkflow(workflow, cfg.Workflows[workflow], cfg.Jobs)
	if err != nil {
		return "", err
	}

	return strings.Join(executor.logs, "\n"), nil
}
