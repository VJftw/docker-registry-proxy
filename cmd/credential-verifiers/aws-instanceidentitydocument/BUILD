go_binary(
    name = "aws-instanceidentitydocument",
    srcs = ["main.go"],
    deps = [
        "//api/proto/v1",
        "//pkg/plugin",
        "//pkg/auth/aws",
        "//third_party/go:protobuf",
        "//third_party/go:mozilla_pkcs7",
    ],
)

go_test(
    name = "gcp-instanceidentitydocument_test",
    srcs = [
        "main.go",
        "main_test.go",
    ],
    deps = [
        "//third_party/go:stretchr_testify",
        "//api/proto/v1",
        "//pkg/auth/aws",
        "//pkg/plugin",
        "//third_party/go:mozilla_pkcs7",
    ],
)
