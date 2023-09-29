package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type Zenith struct{}

func (m *Zenith) MyFunction(ctx context.Context, stringArg string) (*Container, error) {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg}).Sync(ctx)
}
func (m *Zenith) Update(ctx context.Context, repo string, operatingSystem string) (*File, error) {
	ts := time.Now().Unix()
	semver := fmt.Sprintf("v0.0.%d", ts) // Wanted to do 0.0.0+{ts}, but dockerhub does not allow + semver metadata
	registry := fmt.Sprintf("%s/dagger-zenith-engine", repo)
	tagged := fmt.Sprintf("%s:%s", registry, semver)

	_, err := dag.Dagger().
		Engine().FromZenithBranch().Worker().WithVersion(semver).Publish(ctx, tagged)
	if err != nil {
		return nil, err
	}

	bin := dag.Dagger().
		Engine().
		FromZenithBranch().
		Cli(EngineCliOpts{
			WorkerRegistry:  registry,
			OperatingSystem: operatingSystem,
			Arch:            runtime.GOARCH,
			Version:         semver,
		})
	return bin, nil
}
