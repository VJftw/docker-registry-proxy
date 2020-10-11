sh_cmd(
    name = "genproto",
    cmd = "set -x; find $PKG_DIR/ -type f -name '*.pb.go' -delete && rsync -avhL plz-out/go/src/$PKG_DIR/ $PKG_DIR/",
)

