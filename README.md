# env2config
Convert environment variables to configuration files. Great for containerizing existing projects that rely on config files.

Some projects just don't support environment variables out of the box. `env2config` can easily be copied in and inserted into the container's `ENTRYPOINT` for simple conversion of environment variables to various config file formats.

## Getting Started
To get started, just add these lines to your Dockerfile:
```Dockerfile
COPY --from=johnstarich/env2config:latest /env2config /
ENTRYPOINT ["/env2config"]
```
_If you already have an entrypoint, inject env2config just before the first arg:_ `ENTRYPOINT ["/env2config", "/myentrypoint.sh"]`

Then add your configs, for example:
```Dockerfile
# Comma separated config names
ENV E2C_CONFIGS=myconf,other

# <name>_OPTS_<setting> are generation settings for this config.
# The FILE and FORMAT opts are required, TEMPLATE is optional.
# Supported formats: yaml, json, toml, ini
ENV MYCONF_OPTS_FILE=/output/my-config.yaml
ENV MYCONF_OPTS_FORMAT=yaml
# <name>_<key> are mappings from config file keys to environment variables.
ENV MYCONF_bind_url=http://example.com
ENV MYCONF_db.address=db.example.com

ENV OTHER_OPTS_FILE=/output/other.yaml
ENV OTHER_OPTS_FORMAT=yaml
ENV OTHER_addresses.0=http://replica0.example.com
ENV OTHER_addresses.1=http://replica1.example.com
```

At runtime, `env2config` will generate the above configuration with 2 new yaml files from environment variables and output it to `/output/my-config.yaml` and `/output/other.yaml`.

`/output/my-config.yaml`:
```yaml
bind_url: http://example.com
db:
  address: db.example.com
```

`/output/other.yaml`:
```yaml
addresses:
    - http://replica0.example.com
    - http://replica1.example.com
```

To require an environment variable with a custom source, use the pattern `<name>_OPTS_IN_<key>=<env>`.
For example, `MYCONF_OPTS_IN_url=BIND_URL` will require the `$BIND_URL` variable, then set it in the myconf config as `url`.

## Projects using env2config

* [JohnStarich/docker-matrix-appservice-slack](https://github.com/JohnStarich/docker-matrix-appservice-slack)
* [JohnStarich/docker-matrix-pantalaimon](https://github.com/JohnStarich/docker-matrix-pantalaimon)
* [JohnStarich/docker-synapse](https://github.com/JohnStarich/docker-synapse)
* [JohnStarich/docker-matrix-turn](https://github.com/JohnStarich/docker-matrix-turn)
