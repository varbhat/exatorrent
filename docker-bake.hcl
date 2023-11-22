target "docker-metadata-action" {
  tags = ["exatorrent:local"]
}

target "default" {
  inherits = ["docker-metadata-action"]
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/arm/v7",
  ]
}

target "artifact" {
  inherits = ["docker-metadata-action"]
  target = "artifact"
  output = ["type=local,dest=./artifact"]
}

target "artifact-darwin" {
  inherits = ["docker-metadata-action", "artifact"]
  target = "artifact-darwin"
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
  ]
}

target "artifact-all" {
  inherits = ["artifact"]
  target = "artifact"
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/arm/v7",
    "linux/ppc64le",
  ]
}

target "release" {
  target = "release"
  output = ["type=local,dest=./release"]
  contexts = {
    artifacts = "./artifact"
  }
}
