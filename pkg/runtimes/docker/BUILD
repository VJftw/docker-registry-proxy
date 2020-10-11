go_library(
    name = "docker",
    srcs = [
        "image-service.go",
        "image-util.go",
    ],
    deps = [
        "//api/proto/v1",
        "//third_party/go:docker",
        "//third_party/go:google_cadvisor",
        "//third_party/go:k8s_cri-api",
        "//third_party/go:docker_distribution",
    ],
    visibility = [
        "//cmd/clients/...",
    ]
)

go_test(
    name = "docker_test",
    srcs = [
        "image-util_test.go",
    ],
    deps = [
        ":docker",
        "//api/proto/v1",
        "//third_party/go:stretchr_testify",
        "//third_party/go:grpc",
    ],
    external = True,
)
