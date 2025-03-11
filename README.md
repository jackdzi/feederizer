# feederizer

Feederizer is both a RSS aggregator and a viewer for RSS feeds, supporting usage from multiple users through its SQLite3-powered backend server, and features a fast, visually appealing terminal UI for accessing feeds.


## Installation and Configuration

Clone the repo somewhere on your computer, navigate into the directory, and run:
```sh
go build cmd/main.go
```
in both the UI and server folder. Then just run the produced binary.

### Docker

Feederizer's backend server can be deployed through Docker or can run concurrently with the terminal UI component. To run with Docker, edit the `config.toml`` file in the root of the project and change `deployment.docker` to true. To have the database be consistent between instances, create a SQLite3 database `feederizer.db` (or any name) and update `deployment.db_path``, and change the volumes option in
`docker-compose.yml` to `/path/to/feederizer.db:/data/feederizer.db:z`. If you want to host the server and allow external IP calls to it, navigate to `server/internal/api/NewRouter.go` and add the desired external IP to the whitelist.

### Local Deployment

To deploy locally, simply set `deployment.docker` to false.
