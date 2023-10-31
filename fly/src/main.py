import dagger
from dagger.mod import function

@function
def deploy(app: str, image: str, token: dagger.Secret) -> str:
    config = f'''
app = "{app}"
primary_region = "ord"

[build]
  image = "{image}"

[http_service]
  internal_port = 8080
  force_https = true
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
    '''
    return (
        fly_base(token)
        .with_new_file("/fly.toml", config) # TODO: make more of these things options
        .with_exec(["/root/.fly/bin/flyctl", "deploy", "--config", "/fly.toml"])
    )

def fly_base(token: dagger.Secret) -> dagger.Container:
    return (
        dagger.container().from_("alpine:latest")
        .with_exec(["apk", "add", "curl"])
        .with_exec(["sh", "-c", "curl -L https://fly.io/install.sh | sh"])
        .with_secret_variable("FLY_API_TOKEN", token)
    )

