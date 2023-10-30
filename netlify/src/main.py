import dagger
from dagger.mod import function

CLI = "netlify-cli@16.9.3"

@function
def deploy(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
    return (
        netlify_base(token)
        .with_mounted_directory("/src", dir)
        .with_exec(["netlify", "deploy", "--dir", "/src", "--site", site, "--prod"])
        .stdout()
    )

@function
def preview(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
    return (
        netlify_base(token)
        .with_mounted_directory("/src", dir)
        .with_exec(["netlify", "deploy", "--dir", "/src", "--site", site])
        .stdout()
    )

@function
def list(token: dagger.Secret) -> str:
    return (
        netlify_base(token)
        .with_exec(["netlify", "sites:list"])
        .stdout()
    )


def netlify_base(token: dagger.Secret) -> dagger.Container:
    return (
        dagger.container().from_("node:21-alpine")
        .with_exec(["npm", "install", "-g", CLI]).with_entrypoint([])
        .with_secret_variable("NETLIFY_AUTH_TOKEN", token)
    )
