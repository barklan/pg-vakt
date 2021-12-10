############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder

ARG BUILDKIT_INLINE_CACHE=1

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Get docker binary.
RUN mkdir -p /app/media
WORKDIR /app/media
RUN wget https://download.docker.com/linux/static/stable/x86_64/docker-20.10.9.tgz
RUN tar xzvf docker-20.10.9.tgz

# Create appuser
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"
WORKDIR $GOPATH/src/mypackage/myapp/

# use modules
COPY go.mod .

ENV GO111MODULE=on
RUN go mod download
RUN go mod verify

COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -o /go/bin/pgvakt ./cmd/pgvakt/.

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/media/docker/docker /usr/bin/docker

# Copy our static executable
COPY --from=builder /go/bin/pgvakt /go/bin/pgvakt

# Use an unprivileged user.
# USER appuser:appuser

# Run the pgvakt binary.
ENTRYPOINT ["/go/bin/pgvakt"]
