import dagger
from dagger.mod import field, function, object_type

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"

def init() -> dagger.Container:
    return (
        dagger.container().from_("nginx:1.25.3")
        .with_entrypoint([])
        .with_exec(['sh', '-c', f'echo "" > {NGINX_CONFIG}'])
    )

@object_type
class Proxy:
    """Forwards multiple services into a single service with multiple ports"""

    ctr: dagger.Container = field(default=init)

    @function
    def with_service(
        self,
        service: dagger.Service,
        name: str,
        frontend: int,
        backend: int
    ) -> "Proxy":
        """Add a service to proxy"""
        cfg = get_config(backend, name, frontend)

        ctr = self.ctr.with_service_binding(name, service).with_exposed_port(frontend)
        self.ctr = ctr.with_exec(['sh', '-c', f'echo "{cfg}" >> {NGINX_CONFIG}'])

        return self

    @function
    def service(self) -> dagger.Service:
        """Get the proxy Service"""
        ctr = self.ctr.with_exec(["/docker-entrypoint.sh", "nginx", "-g", "daemon off;"])
        return ctr.as_service()

def get_config(port: int, name: str, frontend: int) -> str:
    return f'''
server {{
    listen {frontend};
    listen [::]:{frontend};

    server_name {name};

    location / {{
        proxy_pass http://{name}:{port};
    }}
}}'''
