sh_cmd(
    name = "proto",
    cmd = """
    set -e
    dirs=$(plz query alltargets --include=grpc | cut -f1 -d":" | cut -c 3-)
    for dir in \\\$dirs; do
        echo "-> Linting proto files in \\\$dir"
        $(out_location //third_party/lang:prototool) lint \\\$dir
    done
    """,
    data = ["//third_party/lang:prototool"],
)

sh_cmd(
    name = "go",
    cmd = """
    set -e
    dirs="$(plz query alltargets --include=go | grep -v third_party | cut -f1 -d":" | cut -c 3- | xargs)"
    echo "-> Linting go files"
    $(out_location //third_party/lang:golang_ci_lint) run \\\$dirs
    """,
    data = ["//third_party/lang:golang_ci_lint"],
)
