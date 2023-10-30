import dagger
from dagger.mod import function

CLI = "netlify-cli@16.9.3"

@function
def deploy(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
    return (
            netlify_base()
            .with_mounted_secret("NETLIFY_AUTH_TOKEN", token)
            .with_env_variable("NETLIFY_SITE_ID", site)
            .with_mounted_directory("/src", dir)
            .with_exec(["netlify", "deploy", "--dir", "/src"])
            .stdout()
    )


def netlify_base() -> dagger.Container:
    return (
        dagger.container().from_("node:21-alpine")
        .with_exec(["npm", "install", "-g", CLI])
    )
