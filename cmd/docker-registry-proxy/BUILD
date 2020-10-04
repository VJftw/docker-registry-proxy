go_binary(
    name = "docker-registry-proxy",
    srcs = ["main.go"],
    deps = [
        "//api/proto/v1",
        "//pkg/cmd",
        "//pkg/plugin",
        "//pkg/runtimes/docker/registryproxy",
        "//third_party/go:spf13_viper",
    ],
)
