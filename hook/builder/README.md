# Captain Hook

Captain Hook is a Docker build system that uses [Captain](https://github.com/harbur/captain) to build, tag and push images to the docker registry.

It provides an http server that will handle `push` events from GitHub webhooks to trigger a build. When a `push` event is received, it runs the [build]() docker image, which:

1. Clones your code.
2. Pulls the last built docker image for the branch.
3. Builds the image with `captain build`.
4. Pushes the image with `captain push`.

## Scale out

Captain Hook only needs to talk to the Docker daemon API to run the builder image. The best way to scale out your build cluster is to use [Docker Swarm](https://github.com/docker/swarm).
