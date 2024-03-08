"""
Execute Ruff on a Python project

A utility module to run ruff check on a Python project
"""

import dagger
from dagger import dag, function

@function
async def check(directory: dagger.Directory) -> str:
    """run ruff check"""
    return await (
        dag.container()
        .from_("python:3.10-alpine")
        .with_exec(["pip", "install", "ruff"])
        .with_mounted_directory("/src", directory)
        .with_workdir("/src")
        .with_exec(["ruff", "check", "."])
        .stdout()
    )
