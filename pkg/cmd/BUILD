go_library(
    name = "cmd",
    srcs = [
        "cmd.go",
        "grpc.go",
        "http.go",
        "log.go",
        "plugin.go",
    ],
    deps = [
        "//pkg/plugin",
        "//third_party/go:grpc",
        "//third_party/go:spf13_cobra",
        "//third_party/go:spf13_viper",
        "//third_party/go:uber_zap",
        "//third_party/go:grpc-ecosystem_grpc-middleware",
        "//api/proto/v1",
    ],
    visibility = ["PUBLIC"],
)