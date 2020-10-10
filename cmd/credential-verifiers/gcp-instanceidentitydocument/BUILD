go_binary(
    name = "gcp-instanceidentitydocument",
    srcs = ["main.go"],
    deps = [
        "//api/proto/v1",
        "//pkg/plugin",
        "//pkg/auth/gcp",
        "//third_party/go:protobuf",
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
        "//pkg/auth/gcp",
        "//pkg/plugin",
        "//third_party/go:dgrijalva_jwt",
    ],
)
