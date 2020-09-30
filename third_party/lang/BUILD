package(default_visibility = ["PUBLIC"])
subinclude(
    "//third_party/defs:go_tool",
    "//third_party/defs:protoc",
)

go_tool(
    name="go",
    version = "1.15.2",
    hashes = ["b49fda1ca29a1946d6bb2a5a6982cf07ccd2aba849289508ee0f9918f6bb4552"]
)

protoc_binary(
    name = "protoc",
    version = "3.13.0",
    hashes  = ["4a3b26d1ebb9c1d23e933694a6669295f6a39ddc64c3db2adf671f0a6026f82e"]
)

remote_file(
    name = "prototool_download",
    url = "https://github.com/uber/prototool/releases/download/v1.10.0/prototool-Linux-x86_64.tar.gz",
    hashes = [],
)

genrule(
    name="prototool",
    srcs = [":prototool_download"],
    cmd = ["tar -xzf $SRCS"],
    outs = ["prototool/bin/prototool"],
    binary = True,
)

GOLANG_CI_LINT_VERSION = "1.31.0"

remote_file(
    name = f"golang_ci_lint_download",
    url = "https://github.com/golangci/golangci-lint/releases/download/v%s/golangci-lint-%s-linux-amd64.tar.gz" % (GOLANG_CI_LINT_VERSION, GOLANG_CI_LINT_VERSION),
)

genrule(
    name = f"golang_ci_lint",
    srcs = [":golang_ci_lint_download"],
    cmd = "tar -xzf $SRCS",
    outs = ["golangci-lint-%s-linux-amd64/golangci-lint" % GOLANG_CI_LINT_VERSION],
    binary = True,
)