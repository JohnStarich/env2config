# env2config
Convert environment variables to configuration files. Great for containerizing existing projects that rely on config files.

Some projects just don't support environment variables out of the box. `env2config` can easily be copied in and inserted into the container's `ENTRYPOINT` for simple conversion of environment variables to various config file formats.


## Projects using env2config

* [JohnStarich/docker-matrix-appservice-slack](https://github.com/JohnStarich/docker-matrix-appservice-slack)
* [JohnStarich/docker-matrix-pantalaimon](https://github.com/JohnStarich/docker-matrix-pantalaimon)
* [JohnStarich/docker-synapse](https://github.com/JohnStarich/docker-synapse)
