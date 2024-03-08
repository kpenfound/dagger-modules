// Run CircleCI config in Dagger
//
// A module that attempts to parse and execute a CircleCI configuration
// using Dagger primitives

package main

import (
	"context"
	"strings"
)

const CONFIG = "./.circleci/config.yml"

type Encircle struct{}

// Execute a job in a workflow
func (e *Encircle) EncircleJob(ctx context.Context, d *Directory, job string) (string, error) {
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

// Execute an entire workflow
func (e *Encircle) EncircleWorkflow(ctx context.Context, d *Directory, workflow string) (string, error) {

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
