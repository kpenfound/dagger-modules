import dagger
from dagger.mod import function

CLI = "netlify-cli@16.9.3"

@function
def deploy(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
    return (
            netlify_base(dir, token, site)
            .with_exec(["netlify", "deploy", "--dir", "/src", "--prod"])
            .stdout()
    )

@function
def preview(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
    return (
            netlify_base(dir, token, site)
            .with_exec(["netlify", "deploy", "--dir", "/src"])
            .stdout()
    )


def netlify_base(dir: dagger.Directory, token: dagger.Secret, site: str) -> dagger.Container:
    return (
        dagger.container().from_("node:21-alpine")
        .with_exec(["npm", "install", "-g", CLI])
        .with_mounted_secret("NETLIFY_AUTH_TOKEN", token)
        .with_env_variable("NETLIFY_SITE_ID", site)
        .with_mounted_directory("/src", dir)
    )
