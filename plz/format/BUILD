sh_cmd(
    name = "proto",
    cmd = """
dirs=$(plz query alltargets --include=grpc | cut -f1 -d":" | cut -c 3-)
for dir in \\\$dirs; do
    $(out_location //third_party/lang:prototool) format -fw \\\$dirs
done
    """,
    data = ["//third_party/lang:prototool"],
)

sh_cmd(
    name = "go",
    cmd = """
set -e
dirs="$(plz query alltargets --include=go | grep -v third_party | cut -f1 -d":" | cut -c 3- | xargs)"
set -x;
$(out_location //third_party/lang:gofmt) -s -w \\\$dirs
    """,
    data = ["//third_party/lang:gofmt"],
)
