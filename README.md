# **WARNING**

Make sure that you agree and stick to the policies
of the tile-providers before downloading!

---
<div align="center">
  <img class="logo" src="https://raw.githubusercontent.com/superboomer/map-tile-provider/master/assets/logo.png" width="128px" height="128px" alt="logo"/>
  <br>
  <br>
  <b>maptiled</b>
  <br>
  <br>

  [![build](https://github.com/superboomer/maptiled/actions/workflows/build.yml/badge.svg)](https://github.com/superboomer/maptiled/actions/workflows/build.yml)&nbsp;[![Go Report Card](https://goreportcard.com/badge/github.com/superboomer/maptiled)](https://goreportcard.com/report/github.com/superboomer/maptiled)&nbsp;[![Coverage Status](https://coveralls.io/repos/github/superboomer/maptiled/badge.svg?branch=master)](https://coveralls.io/github/superboomer/maptiled?branch=master)
</div>


 Easy CLI tool to download specified tiles by [maptile](https://github.com/superboomer/maptile/)

---
#### Options

 ***maptiled*** supports the following command-line options:

- `-s`, `--save-path`: define where maptiled save tiles.
- `-p`, `--provider-url`: url where maptile serving.
- `-z`, `--zoom`: define zoom (z) for each tile.
- `--side`: count of tile of result image square.
- `--set-max`: if provider not support specified zoom, maptiled will use maximum zoom.
- `--providers`: download specified providers. (only if maptile support them).
- `--points`: path for a points.json file.

> All environment variables are available in [source code](https://github.com/superboomer/maptiled/blob/master/internal/options/opt.go)
***
# **Install**
To start using latest released version, just run:

```
$ go install github.com/superboomer/maptiled@latest
```

# **Docker Deploy**

You can easly deploy it via docker. Basic ***docker-compose.yml*** may look like this:
```YAML
version: '3.7'

services:

  map-tile-provider:
    image: ghcr.io/superboomer/maptiled:latest
    container_name: map-tile-provider
    restart: unless-stopped
    environment:
      - PROVIDER_URL=http://example.com
      - POINTS=./example_points.json
    ports:
      - "8080:8080"
```

***

