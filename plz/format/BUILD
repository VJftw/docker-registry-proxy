sh_cmd(
    name = "proto",
    cmd = """
    dirs=$(plz query alltargets --include=grpc | cut -f1 -d":" | cut -c 3-)
    for dir in \\\$dirs; do
        set -x; 
        $(out_location //third_party/lang:prototool) format -fw \\\$dirs
        set +x;
    done
    """,
    data = ["//third_party/lang:prototool"],
)
