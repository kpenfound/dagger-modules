# Zenith

A module to build and update your Dagger Zenith CLI and Engine based on the latest from the Zenith branch

## Usage

Authenticate with Dockerhub:
`docker login`

Run this dagger query, replacing `$DOCKERHUB_USERNAME` with your Dockerhub username and the `operatingSystem` with your OS (`darwin` or `linux`):
`echo '{ zenith { update(repo: "$DOCKERHUB_USERNAME", operatingSystem: "darwin") { export(path:"./daggerz")} }}' | dagger query -m github.com/kpenfound/dagger-modules/zenith`

This will build a new engine image at `$DOCKERHUB_USERNAME/dagger-zenith-engine:v0.0.{epoch}`, which the resulting CLI at `./daggerz` will use as it's engine. Replace your existing Dagger Zenith binary with this new binary, and you're updated!
