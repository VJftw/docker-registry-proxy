package(default_visibility = ["PUBLIC"])

go_get(
    name = "protobuf",
    get = "github.com/golang/protobuf",
    install = [
        "proto",
        "ptypes",
        "ptypes/duration",
        "ptypes/timestamp",
        "ptypes/empty",
        "jsonpb",
    ],
    revision = "v1.3.5",
    licences = ["BSD-3-Clause"],
)

go_get(
    name = "protoc-gen-go",
    binary = True,
    get = [],
    install = ["github.com/golang/protobuf/protoc-gen-go"],
    deps = [":protobuf"],
    licences = ["BSD-3-Clause"],
)

go_get(
    name = "x_net",
    get = "golang.org/x/net",
    revision = "release-branch.go1.15",
    install = [
        "http2"
    ],
    deps = [":x_text"],
)

go_get(
    name = "x_text",
    get = "golang.org/x/text",
    revision = "release-branch.go1.15",
    install = [
        "secure/...",
        "unicode/...",
    ],
)

go_get(
    name = "x_sys",
    get = "golang.org/x/sys",
    revision = "release-branch.go1.15",
    install = [
        "unix",
    ],
)

go_get(
    name = "googleapis",
    get = "google.golang.org/genproto/...",
    deps = [":protobuf"],
    install = ["googleapis/rpc/status"],
    revision = "24fa4b261c55da65468f2abfdae2b024eef27dfb",
)


go_get(
    name = "grpc",
    get = "google.golang.org/grpc",
    install = [
        "",
        "codes",
        "status",
    ],
    revision = "v1.32.0",
    licences = ["Apache-2.0"],
    deps = [":protobuf", ":x_net", ":x_sys", ":googleapis"],
)

go_get(
    name = "hashicorp_plugin",
    get = "github.com/hashicorp/go-plugin",
    revision = "v1.3.0",
    licences = ["MPL-2.0"],
    deps = [
        ":protobuf", 
        ":grpc", 
        ":hashicorp_hclog", 
        ":hashicorp_yamux", 
        ":mitchellh_testing-interface",
        ":oklog_run",
    ],
)

go_get(
    name = "hashicorp_hclog",
    get = "github.com/hashicorp/go-hclog",
    revision = "v0.14.1",
    licences = ["MIT"],
    deps = [":fatih_color", ":mattn_isatty"],
)

go_get(
    name = "hashicorp_yamux",
    get = "github.com/hashicorp/yamux",
    revision = "aecfd211c9ce5303cd19780c8fc96a01b42fa298",
    licences = ["MPL-2.0"],
)

go_get(
    name = "fatih_color",
    get = "github.com/fatih/color",
    revision = "v1.9.0",
    licences = ["MIT"],
)

go_get(
    name = "mattn_isatty",
    get = "github.com/mattn/go-isatty",
    revision = "v0.0.12",
    licences = ["MIT"],
    deps = [":x_sys"],
)

go_get(
    name = "mitchellh_testing-interface",
    get = "github.com/mitchellh/go-testing-interface",
    revision = "v1.14.1",
    licences = ["MIT"],
)

go_get(
    name = "oklog_run",
    get = "github.com/oklog/run",
    revision = "v1.1.0",
    licences = ["Apache-2.0"],
)

go_get(
    name = "uber_zap",
    get = "go.uber.org/zap",
    revision = "v1.16.0",
    deps = [":uber_atomic", ":uber_multierr"],
)

go_get(
    name = "uber_atomic",
    get = "go.uber.org/atomic",
    revision = "v1.7.0",
)

go_get(
    name = "uber_multierr",
    get = "go.uber.org/multierr",
    revision = "v1.6.0",
    deps = [":uber_atomic"],
)

go_get(
    name = "spf13_cobra",
    get = "github.com/spf13/cobra",
    revision = "v1.0.0",
    licences = ["Apache-2.0"],
    deps = [":spf13_pflag"],
)

go_get(
    name = "spf13_pflag",
    get = "github.com/spf13/pflag",
    revision = "v1.0.5",
    licences = ["BSD-3-Clause"],
)

go_get(
    name = "spf13_viper",
    get = "github.com/spf13/viper",
    revision = "v1.7.1",
    licences = ["MIT"],
    deps = [":spf13_pflag"],
)