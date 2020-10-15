go_binary(
    name = "kubelet-image-service",
    srcs = ["main.go"],
    deps = [
        "//api/proto/v1",
        "//pkg/cmd",
        "//pkg/plugin",
        "//pkg/runtimes/docker",
        "//third_party/go:spf13_viper",
        "//third_party/go:k8s_cri-api",
    ],
    visibility = ["//build/clients/kubelet-image-service/..."],
)
