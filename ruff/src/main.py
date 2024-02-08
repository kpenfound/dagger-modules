import dagger
from dagger import dag, function

@function
async def check(directory: dagger.Directory) -> str:
    return await (
        dag.container()
        .from_("pipeline-components/ruff:latest")
        .with_mounted_directory("/src", directory)
        .with_workdir("/src")
        .with_exec(["ruff", "check", "."])
        .stdout()
    )
