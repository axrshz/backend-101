### Image Processing API

Create `.env` from `.env.example` and provide the S3 credentials.

Build and start the API, worker, PostgreSQL, Redis, and database migration:

```sh
docker compose up --build
```

Open `http://localhost:3000`.

Stop the services:

```sh
docker compose down
```

Remove the persisted PostgreSQL and Redis data as well:

```sh
docker compose down -v
```
