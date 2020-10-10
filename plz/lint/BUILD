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

sh_cmd(
    name = "go",
    cmd = """
    set -e
    set -x;
    dirs="$(plz query alltargets --include=go | grep -v third_party | cut -f1 -d":" | cut -c 3- | xargs)"
    $(out_location //third_party/lang:golang_ci_lint) run --enable-all \\\$dirs
    set +x;
    """,
    data = ["//third_party/lang:golang_ci_lint"],
)
