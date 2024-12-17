#!/bin/bash

# Pass a local postgres database as a service to proxy
dagger call with-service --service tcp://localhost:5432 --name postgres --frontend 5433 --backend 5432

