import dagger
from dagger.mod import function

@function
def deploy(image: str, port: int, token: dagger.Secret) -> str:
    return (
        fly_base(token)
        .with_exec(["fly", "deploy", "--image", image, "--internal-port", port])
    )

def fly_base(token: dagger.Secret) -> dagger.Container:
    return (
        dagger.container().from_("alpine:latest")
        .with_exec(["apk", "add", "curl"])
        .with_exec(["sh", "-c", "curl -L https://fly.io/install.sh | sh"])
        .with_secret_variable("FLY_API_TOKEN", token)
    )

