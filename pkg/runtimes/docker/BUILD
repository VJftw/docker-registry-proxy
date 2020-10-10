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
