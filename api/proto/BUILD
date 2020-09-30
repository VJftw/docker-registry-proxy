sh_cmd(
    name = "genproto",
    cmd = "set -x; rsync -avhL plz-out/go/src/$PKG_DIR/ $PKG_DIR/",
)

sh_cmd(
    name = "lint",
    cmd = "set -x; $(out_location //third_party/lang:prototool) lint $PKG_DIR",
    data = ["//third_party/lang:prototool"],
)

sh_cmd(
    name = "format",
    cmd = "set -x; $(out_location //third_party/lang:prototool) format -fw $PKG_DIR",
    data = ["//third_party/lang:prototool"],
)

sh_cmd(
    name = "break",
    cmd = "set -x; $(out_location //third_party/lang:prototool) break check $PKG_DIR",
    data = ["//third_party/lang:prototool"],
)
