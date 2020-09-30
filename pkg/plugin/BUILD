go_lib(
    name = "plugin",
    srcs = [
        "auth-provider.go",
        "auth-verifier.go",
        "configuration.go",
        "discovery.go",
        "plugin.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go:hashicorp_plugin",
        "//third_party/go:grpc",
    ],
    visibility = ["PUBLIC"],
)