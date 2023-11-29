import dagger
from dagger.mod import field, function, object_type

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"

def init() -> dagger.Container:
    return (
        dagger.container().from_("nginx:1.25.3")
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
        self.ctr = (
            self.ctr
            .with_new_file(f"/etc/nginx/conf.d/{name}.conf", contents=cfg)
            .with_service_binding(name, service)
            .with_exposed_port(frontend)
        )
        return self

    @function
    def service(self) -> dagger.Service:
        """Get the proxy Service"""
        return (
            self.ctr
            .with_exec(["nginx", "-g", "daemon off;"])
            .as_service()
        )

def get_config(port: int, name: str, frontend: int) -> str:
    return f'''
server {{
    listen {frontend};
    listen [::]:{frontend};

    server_name {name};

    location / {{
        proxy_pass http://{name}:{port};
        proxy_set_header Host $http_host;
        proxy_set_header X-Forwarded-Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }}
}}'''
