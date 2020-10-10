go_library(
    name = "aws",
    srcs = [
        "certificates.go",
        "instanceidentitydocument-provider.go",
        "instanceidentitydocument-verifier.go",
        "metadata.go",
    ],
    visibility = [
        "//cmd/credential-providers/...",
        "//cmd/credential-verifiers/...",
    ],
)

go_test(
    name = "aws_test",
    srcs = [
        "instanceidentitydocument-provider_test.go",
    ],
    deps = [
        ":aws",
        "//third_party/go:stretchr_testify",
    ],
    external = True,
)
