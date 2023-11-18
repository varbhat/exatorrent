target "docker-metadata-action" {}

target "default" {
  inherits = ["docker-metadata-action"]
}

target "artifact" {
  inherits = ["docker-metadata-action"]
  target = "artifact"
  output = ["type=local,dest=./artifact"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
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
