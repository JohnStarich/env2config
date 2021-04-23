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
ENV E2C_CONFIGS=myconf
ENV MYCONF_OPTS_FILE=/output/my-config.yaml
ENV MYCONF_OPTS_FORMAT=yaml

ENV MYCONF_bind_url=http://example.com
ENV MYCONF_db.address=db.example.com
```

At runtime, `env2config` will generate a new yaml file from environment variables and output it to `/output/my-config.yaml`.
The above configuration will generate:
```yaml
bind_url: http://example.com
db:
  address: db.example.com
```

## Projects using env2config

* [JohnStarich/docker-matrix-appservice-slack](https://github.com/JohnStarich/docker-matrix-appservice-slack)
* [JohnStarich/docker-matrix-pantalaimon](https://github.com/JohnStarich/docker-matrix-pantalaimon)
* [JohnStarich/docker-synapse](https://github.com/JohnStarich/docker-synapse)
