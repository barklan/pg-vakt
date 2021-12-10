# pg-vakt

First, configure postgres:

```yaml
db:
  container_name: prod_db
  image: 'postgres:14'
  command: |
    postgres -c wal_level=replica
    -c archive_mode=on
    -c archive_timeout=600
    -c archive_command='cp %p /pgbackups/%f'
  volumes:
    - /home/ubuntu/prod/pgbackups:/pgbackups

  ...
```

On backup server create volume for backups:

```bash
docker volume create pg-vakt-media
```

Optionally set hostname:

```bash
hostnamectl set-hostname example.com
```

Configure pg-vakt (`./pg_vakt.env`):

```ini
PG_VAKT_TG_BOT_TOKEN=3242342fdfh323342:43gdf33223r23f22fsvjojojnft
PG_VAKT_TG_CHAT_ID=-13534534534534
PG_VAKT_INTERVAL_MINUTES=360
PG_VAKT_CONTINUOUS=true
PG_VAKT_CONTINUOUS_PATH=/home/ubuntu/prod/pgbackups  # path where WAL files are copied to
PG_VAKT_SSH_USER=root
PG_VAKT_SSH_HOSTNAME=example.com
PG_VAKT_SSH_KEY_FILENAME=example_key.pem  # this is relative to /app/config in container
PG_VAKT_CONTAINER_NAME=prod_db
PG_VAKT_DATABASE=app
```

Run it:

```bash
docker run --rm --name pg-vakt -it \
-v /var/run/docker.sock:/var/run/docker.sock \
-v pg-vakt-media:/app/media \
-v "$(pwd)"/environment:/app/config \
--env-file pg_vakt.env \
-e HOST_HOSTNAME="$(hostname)" \
barklan/pg-vakt:rolling
```
