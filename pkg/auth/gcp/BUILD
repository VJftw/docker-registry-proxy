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
    visibility = ["//plugins/..."],
)
