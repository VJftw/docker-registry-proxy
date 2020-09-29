filegroup(
    name="v1_srcs",
    srcs = [
        "v1/authentication.proto",
        "v1/configuration.proto",
    ]
)

grpc_library(
    name = "v1",
    srcs = [":v1_srcs"],
    languages = ["go"],
    visibility = ["PUBLIC"],
)

sh_cmd(
    name = "go_proto_pkg",
    srcs = [":v1"],
    cmd = """
cp -r -L plz-out/go/src/$(out_location $SRCS) api/
rm -rf $(out_location $SRCS)/v1/v1
""",
)

sh_cmd(
    name = "lint",
    srcs = [":v1_srcs"],
    cmd = """
set -x
$(out_location //third_party/lang:prototool) lint $SRCS
""",
    data = ["//third_party/lang:prototool"],
)
