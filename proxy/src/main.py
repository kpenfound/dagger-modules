import dagger
from dagger.mod import Annotated, Doc, field, function, object_type

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"

@object_type
class Proxy:
    """Forwards multiple services into a single service with multiple ports"""

    ctr: Annotated[
            dagger.Container,
            Doc("Internal proxy container"),
    ] = field(default=init())

    @function
    async def with_service(
        svc: Annotated[dagger.Service, Doc("The service to proxy")],
        name: Annotated[str, Doc("The internal name of the service")],
        frontend: Annotated[int, Doc("The frontend port for the proxy")],
    ) -> Proxy:
        """Add a service to proxy"""
        ctr = dagger.container().from_("nginx").with_entrypoint([])
        ports = await svc.ports()
        port = await ports[0].port()
        cfg = get_config(port, name, frontend)

        ctr = ctr.with_service_binding(name, svc).with_exposed_port(frontend)
        ctr = ctr.with_exec(['sh', '-c', f'echo "{cfg}" >> {NGINX_CONFIG}'])

        return ctr

    @function
    async def service() -> dagger.Service:
        """Get the proxy Service"""
        ctr = self.ctr.with_exec(["/docker-entrypoint.sh", "nginx", "-g", "daemon off;"])
        return ctr.as_service()

def init() -> dagger.Container:
    return (
        dagger.container().from_("nginx")
        .with_entrypoint([])
        .with_exec(['sh', '-c', f'echo "" > {NGINX_CONFIG}'])
    )

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
