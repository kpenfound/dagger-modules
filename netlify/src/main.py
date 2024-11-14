"""
Deploy apps to Netlify

A utility module for deploying apps to Netlify
"""

import dagger
from dagger import function, dag, object_type

CLI = "netlify-cli@16.9.3"

@object_type
class Netlify:
    @function
    def deploy(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
        """Deploy a site to production"""
        return (
            netlify_base(token)
            .with_mounted_directory("/src", dir)
            .with_exec(["netlify", "deploy", "--dir", "/src", "--site", site, "--prod"])
            .stdout()
        )

    @function
    def preview(dir: dagger.Directory, token: dagger.Secret, site: str) -> str:
        """Deploy a preview site"""
        return (
            netlify_base(token)
            .with_mounted_directory("/src", dir)
            .with_exec(["netlify", "deploy", "--dir", "/src", "--site", site])
            .stdout()
         )

    @function
    def list(token: dagger.Secret) -> str:
        """List sites"""
        return (
            netlify_base(token)
            .with_exec(["netlify", "sites:list"])
            .stdout()
        )


def netlify_base(token: dagger.Secret) -> dagger.Container:
    return (
        dag.container().from_("node:21-alpine")
        .with_exec(["npm", "install", "-g", CLI]).with_entrypoint([])
        .with_secret_variable("NETLIFY_AUTH_TOKEN", token)
    )
