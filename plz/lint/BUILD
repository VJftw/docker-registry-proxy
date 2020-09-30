sh_cmd(
    name = "proto",
    cmd = """
    set -e
    dirs=$(plz query alltargets --include=grpc | cut -f1 -d":" | cut -c 3-)
    for dir in \\\$dirs; do
        set -x; 
        $(out_location //third_party/lang:prototool) lint \\\$dirs
        set +x;
    done
    """,
    data = ["//third_party/lang:prototool"],
)
