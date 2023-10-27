import dagger
from dagger.mod import function

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"

@function
async def proxy(svc: dagger.Service, name: str, frontend: int) -> dagger.Container:
    ctr = dagger.container().from_("nginx")
    ports = await svc.ports()
    cfg = get_config(ports[0], name, frontend)

    ctr = ctr.with_service_binding(name, svc).with_exposed_port(frontend)
    ctr = ctr.with_exec(['sh', '-c', f'echo {cfg} > {NGINX_CONFIG}'])

    return ctr

@function
async def additional_proxy(ctr: dagger.Container, svc: dagger.Service, name: str, frontend: int) -> dagger.Container:
    ports = await svc.ports()
    cfg = get_config(ports[0], name, frontend)

    ctr = ctr.with_service_binding(name, svc).with_exposed_port(frontend)
    ctr = ctr.with_exec(['sh', '-c', f'echo {cfg} >> {NGINX_CONFIG}'])

    return ctr

def get_config(port: dagger.Port, name: str, frontend: int) -> str:
    backend = port.port()
    protocol = port.protocol()
    return f'''
server {{
    listen {frontend};
    listen [::]:{frontend};

    server_name {name};

    location / {{
        proxy_pass {protocol}://{name}:{backend};
        include proxy_params;
    }}
}}
    '''
