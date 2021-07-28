# STEP 1 - Build the binary

# use the golang image as build image
FROM golang

# build arguments in order to reuse it for arm containers
ARG GOOS=linux
ARG GOARCH=amd64
ARG GOARM

# activate go modules
ENV GO111MODULE on

# copy the local package files to the container's workspace.
COPY . /go/src/github.com/dj95/huekit

# set the working directory to build the application
WORKDIR /go/src/github.com/dj95/huekit

# compile the program
RUN CGO_ENABLED=0 GOOS="${GOOS}" GOARCH="${GOARCH}" go build \
    -ldflags="-s -w" \
    -a \
    -installsuffix cgo \
    -o /go/bin/huekit \
    /go/src/github.com/dj95/huekit/cmd/huekit/main.go


# STEP 2 - Build a minimal container

# start from scratch
FROM scratch

# copy the static executable
COPY --from=0 /go/bin/huekit /huekit

# define the entrypoint
ENTRYPOINT ["/huekit"]
