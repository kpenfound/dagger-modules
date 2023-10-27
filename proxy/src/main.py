import dagger
from dagger.mod import function

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"

@function
async def proxy(svc: dagger.Service, name: str, frontend: int) -> dagger.Container:
    ctr = dagger.container().from_("nginx").with_entrypoint([])
    ports = await svc.ports()
    port = await ports[0].port()
    cfg = get_config(port, name, frontend)

    ctr = ctr.with_service_binding(name, svc).with_exposed_port(frontend)
    ctr = ctr.with_exec(['sh', '-c', f'echo "{cfg}" > {NGINX_CONFIG}'])

    return ctr

@function
async def additional_proxy(ctr: dagger.Container, svc: dagger.Service, name: str, frontend: int) -> dagger.Container:
    ports = await svc.ports()
    port = await ports[0].port()
    cfg = get_config(port, name, frontend)

    ctr = ctr.with_service_binding(name, svc).with_exposed_port(frontend)
    ctr = ctr.with_exec(['sh', '-c', f'echo "{cfg}" >> {NGINX_CONFIG}'])

    return ctr

@function
async def service(ctr: dagger.Container) -> dagger.Service:
    ctr = ctr.with_exec(["/docker-entrypoint.sh", "nginx", "-g", "daemon off;"])
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
