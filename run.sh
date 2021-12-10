#!/usr/bin/env bash

set -eo pipefail

DC="${DC:-exec}"

# If we're running in CI we need to disable TTY allocation for docker-compose
# commands that enable it by default, such as exec and run.
TTY=""
if [[ ! -t 1 ]]; then
    TTY="-T"
fi

# -----------------------------------------------------------------------------
# Helper functions start with _ and aren't listed in this script's help menu.
# -----------------------------------------------------------------------------

function _dc {
    export DOCKER_BUILDKIT=1
    docker-compose ${TTY} "${@}"
}

function _use_local_env {
    sort -u environment/local.env | grep -v '^$\|^\s*\#' > './environment/local.env.tempfile'
    export "$(< environment/local.env.tempfile xargs)"
    rm environment/local.env.tempfile
}

# ----------------------------------------------------------------------------

function up {
    docker image rm pg-vakt:local || true
    docker build -t pg-vakt:local .

    docker run --rm -it \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v pg-vakt-media:/app/media \
    -v "$(pwd)"/environment:/app/config \
    --env-file environment/pg_vakt.env \
    -e HOST_HOSTNAME="$(hostname)" \
    pg-vakt:local
}

function push {
    docker image rm barklan/pg-vakt:rolling || true
    docker build -t barklan/pg-vakt:rolling .
    docker image push barklan/pg-vakt:rolling
}

function direct {
    scp -r environment "helper:/home/docker"

    ssh -tt -o StrictHostKeyChecking=no "helper" \
    'docker stop pg-vakt || true && \
    docker rm pg-vakt || true && \
    cd /home/docker && \
    docker volume create pg-vakt-media && \
    docker image rm barklan/pg-vakt:rolling || true && \
    docker run --rm --name pg-vakt -it -d \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v pg-vakt-media:/app/media \
    -v "$(pwd)"/environment:/app/config \
    --env-file environment/pg_vakt.env \
    -e HOST_HOSTNAME=ctopanel.com \
    barklan/pg-vakt:rolling'
}

# -----------------------------------------------------------------------------

function help {
    printf "%s <task> [args]\n\nTasks:\n" "${0}"

    compgen -A function | grep -v "^_" | cat -n
}

TIMEFORMAT=$'\nTask completed in %3lR'
time "${@:-help}"
