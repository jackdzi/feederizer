# feederizer

Installation: Add a database feederizer.db (or any name) and update the config.toml.

TODO: add functionality for locally hosted (not docker) such that the server gets ran as a subproccess

To deploy locally, simply edit the config.toml and set docker to false.

If you want to host the server and allow external IP calls to it, navigate to `server/internal/api/NewRouter.go` and add the desired external IP to the whitelist.
