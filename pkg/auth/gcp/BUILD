go_library(
    name = "gcp",
    srcs = [
        "certificate-manager.go",
        "gcp.go",
        "instanceidentitydocument-verifier.go",
        "metadata.go",
    ],
    deps = [
        "//third_party/go:uber_zap",
        "//pkg/cmd",
        "//third_party/go:dgrijalva_jwt",
    ],
    visibility = [
        "//cmd/credential-providers/...",
        "//cmd/credential-verifiers/...",
    ],
)

go_test(
    name = "gcp_test",
    srcs = [
        "certificate-manager_test.go",
    ],
    deps = [
        ":gcp",
        "//third_party/go:stretchr_testify",
    ],
    external = True,
)
