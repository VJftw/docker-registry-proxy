go_library(
    name = "registryproxy",
    srcs = [
        "auth.go",
        "grpc.go",
        "handler.go",
        "proxy.go",
        "redirect-follower.go",
        "registryproxy.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go:dgrijalva_jwt",
        "//third_party/go:uber_zap",
        "//third_party/go:google_containerregistry",
        "//pkg/cmd",
    ],
    visibility = ["//cmd/..."],
)

go_test(
    name = "registryproxy_test",
    srcs = [
        "handler_test.go",
    ],
    deps = [
        ":registryproxy",
        "//third_party/go:stretchr_testify",
    ],
    external = True,
)
