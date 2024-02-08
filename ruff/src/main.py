import dagger
from dagger import dag, function

@function
async def check(directory: dagger.Directory) -> str:
    return await (
        dag.container()
        .from_("python:3.10-alpine")
        .with_exec(["pip", "install", "ruff"])
        .with_mounted_directory("/src", directory)
        .with_workdir("/src")
        .with_exec(["ruff", "check", "."])
        .stdout()
    )
