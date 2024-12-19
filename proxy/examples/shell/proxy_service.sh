#!/usr/bin/env dagger shell

github.com/kpenfound/dagger-modules/proxy | with-service "postgres" $(container | from "postgres" | as-service) 5433 5432 | service

