# Waypoint Dobi Builder Plugin

Builder plugin for [HashiCorp Waypoint](https://www.waypointproject.io/).

## Usage

1. Setup a [dobi](https://github.com/dnephin/dobi) configuration for your build:

``` yaml
meta:
  project: project

image=split:
  image: registry.gitlab.com/user/app
  pull-base-image-on-build: true
  cache-from:
    - registry.gitlab.com/user/app:latest
    - "registry.gitlab.com/user/app:{env.TAG}"
  context: .
  tags:
    - "{env.TAG}"
    - latest

```

2. Use Waypoint to build your images with a `dobi` builder:

``` hcl
project = "project"

app "app" {
    build {
        use "dobi" {
            image = "app"

            ## Optional environment for Dobi.  It automatically
            ## inherits the OS environment.
            env = {
                TAG = "develop"
            }
        }
    }
}
```

``` shell
waypoint build
```

```
✓ Initializing dobi build context
✓ Executing command: dobi app:build
 │ .....
 │ .....
 │ Step 7/7 : ENTRYPOINT ["/app/app"]
 │  ---> Using cache
 │  ---> 7c055d74d442
 │ Successfully built 7c055d74d442
 │ Successfully tagged registry.gitlab.com/user/app:develop
 │ [image:build split] registry.gitlab.com/user/app Create
```

## Build

Install all the [required
software](https://www.waypointproject.io/docs/extending-waypoint/creating-plugins#requirements)
by following the official plugin development instructions.

Run the Makefile to compile the plugin, the `Makefile` will build the
plugin for all architectures.  You can comment out those which you are
not interested in.

```shell
make
```

```shell
Build Protos
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./builder/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./registry/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./platform/output.proto
protoc -I . --go_out=plugins=grpc:. --go_opt=paths=source_relative ./release/output.proto

Compile Plugin
# Clear the output
rm -rf ./bin
GOOS=linux GOARCH=amd64 go build -o ./bin/linux_amd64/waypoint-plugin-dobi ./main.go
GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_amd64/waypoint-plugin-dobi ./main.go
GOOS=windows GOARCH=amd64 go build -o ./bin/windows_amd64/waypoint-plugin-dobi.exe ./main.go
GOOS=windows GOARCH=386 go build -o ./bin/windows_386/waypoint-plugin-dobi.exe ./main.go
```

After the build is finished, run

``` shell
make install
```

to make the plugin available for Waypoint.

## Building with Docker

To build plugins for release you can use the `build-docker` Makefile
target, this will build your plugin for all architectures and create
zipped artifacts which can be uploaded to an artifact manager such as
GitHub releases.

The built artifacts will be output in the `./releases` folder.

```shell
make build-docker

rm -rf ./releases
DOCKER_BUILDKIT=1 docker build --output releases --progress=plain .
#1 [internal] load .dockerignore
#1 transferring context: 2B done
#1 DONE 0.0s

#...

#14 [export_stage 1/1] COPY --from=build /go/plugin/bin/*.zip .
#14 DONE 0.1s

#15 exporting to client
#15 copying files 36.45MB 0.1s done
#15 DONE 0.1s
```

## Building and releasing with GitHub Actions

When cloning the template a default GitHub Action is created at the
path `.github/workflows/build-plugin.yaml`. You can use this action to
automatically build and release your plugin.

The action has two main phases:
1. **Build** - This phase builds the plugin binaries for all the supported architectures. It is triggered when pushing
   to a branch or on pull requests.
1. **Release** - This phase creates a new GitHub release containing the built plugin. It is triggered when pushing tags
   which starting with `v`, for example `v0.1.0`.

You can enable this action by clicking on the `Actions` tab in your GitHub repository and enabling GitHub Actions.
