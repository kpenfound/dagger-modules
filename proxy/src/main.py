"""
Proxy multiple services though a single service

This module allows you to proxy any number of Dagger Services
through a single Dagger Service on specified ports
"""

import dagger
from dagger import dag, function, field, object_type

NGINX_CONFIG = "/etc/nginx/conf.d/default.conf"


def init() -> dagger.Container:
    return (
        dag.container()
        .from_("nginx:1.25.3")
        .with_new_file(
            "/etc/nginx/stream.conf",
            contents=f"""stream {{ include /etc/nginx/stream.d/*.conf; }}""",
        )
        .with_new_file(
            "/etc/nginx/nginx.conf",
            contents=f"""
user  nginx;
worker_processes  auto;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;

events {{
    worker_connections  1024;
}}

include /etc/nginx/stream.conf;

http {{
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;
}}
            """,
        )
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
        backend: int,
        is_tcp: bool = False,
    ) -> "Proxy":
        """Add a service to proxy"""
        cfg = get_config(backend, name, frontend, is_tcp)
        conf_path = (
            f"/etc/nginx/stream.d/{name}.conf"
            if is_tcp
            else f"/etc/nginx/conf.d/{name}.conf"
        )
        self.ctr = (
            self.ctr.with_new_file(conf_path, contents=cfg)
            .with_service_binding(name, service)
            .with_exposed_port(frontend)
        )
        return self

    @function
    def service(self) -> dagger.Service:
        """Get the proxy Service"""
        return self.ctr.as_service(args=["nginx", "-g", "daemon off;"])


def get_config(port: int, name: str, frontend: int, is_tcp: bool) -> str:
    if is_tcp:
        return f"""
    server {{
        listen {frontend};
        listen [::]:{frontend};
        proxy_pass {name}:{port};
    }}
    """
    else:
        return f"""
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
    }}"""
