# drone-rsync

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-rsync/status.svg)](http://beta.drone.io/drone-plugins/drone-rsync)
[![Coverage Status](https://aircover.co/badges/drone-plugins/drone-rsync/coverage.svg)](https://aircover.co/drone-plugins/drone-rsync)
[![](https://badge.imagelayers.io/plugins/drone-rsync:latest.svg)](https://imagelayers.io/?images=plugins/drone-rsync:latest 'Get your own badge on imagelayers.io')

Drone plugin to deploy or update a project via Rsync. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Binary

Build the binary using `make`:

```
make deps build
```

## Docker

Build the container using `make`:

```
make deps docker
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_HOST=foo.com \
  -e PLUGIN_USER=root \
  -e PLUGIN_KEY="$(cat ${HOME}/.ssh/id_rsa)" \
  -e PLUGIN_SOURCE=dist/ \
  -e PLUGIN_TARGET=/path/on/server \
  -e PLUGIN_DELETE=true \
  -e PLUGIN_RECURSIVE=true \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/drone-rsync
```
