grpc_library(
    name = "v1",
    srcs = [
        "v1/authentication.proto",
    ],
    languages = ["go"],
    visibility = ["PUBLIC"],
)

sh_cmd(
    name = "go_proto_pkg",
    srcs = [":v1"],
    cmd = """
cp -r plz-out/go/src/$(out_location $SRCS) api/
rm -rf $(out_location $SRCS)/v1/v1
""",
)
