import dagger
from dagger import dag, function

@function
async def check(directory: dagger.Directory) -> str:
    return await (
        dag.container()
        .from_("ghcr.io/astral-sh/ruff")
        .with_mounted_directory("/src", directory)
        .with_workdir("/src")
        .with_exec(["ruff", "check", "."])
        .stdout()
    )
