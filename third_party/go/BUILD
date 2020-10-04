package(default_visibility = ["PUBLIC"])

sh_cmd(
    name="lint",
    cmd="""
set -x
dir=\\\$(dirname $(out_location //third_party/lang:go))
cp $(out_location //third_party/lang:go) "\\\${dir}/go" 
export PATH="$PATH:\\\${dir}"
plz run //third_party/lang:go_lint -- run pkg/... cmd/...
    """,
    data=["//third_party/lang:go"],
)

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
        "http2",
        "html",
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
    deps = [
        ":spf13_pflag",
        ":fsnotify_fsnotify",
        ":hashicorp_hcl",
        ":magiconair_properties",
        ":mitchellh_mapstructure",
        ":pelletier_toml",
        ":spf13_afero",
        ":spf13_cast",
        ":spf13_jwalterweatherman",
        ":spf13_pflag",
        ":subosito_gotenv",
        ":ini_v1",
        ":yaml_v2",
    ],
)

go_get(
    name = "fsnotify_fsnotify",
    get = "github.com/fsnotify/fsnotify",
    revision = "v1.4.9",
    licences = ["BSD-3-Clause"],
    deps = [":x_sys"],
)

go_get(
    name = "hashicorp_hcl",
    get = "github.com/hashicorp/hcl",
    revision = "v1.0.0",
    licences = ["MPL-2.0"],
)

go_get(
    name = "magiconair_properties",
    get = "github.com/magiconair/properties",
    revision = "v1.8.1",
    licences = ["BSD-2-Clause"],
)

go_get(
    name = "mitchellh_mapstructure",
    get = "github.com/mitchellh/mapstructure",
    revision = "v1.3.3",
    licences = ["MIT"],
)

go_get(
    name = "pelletier_toml",
    get = "github.com/pelletier/go-toml",
    revision = "v1.2.0",
    licences = ["MIT"],
)

go_get(
    name = "spf13_afero",
    get = "github.com/spf13/afero",
    revision = "v1.1.2",
    licences = ["Apache-2.0"],
    deps = [":x_text"],
)

go_get(
    name = "spf13_cast",
    get = "github.com/spf13/cast",
    revision = "v1.3.0",
    licences = ["MIT"],
)

go_get(
    name = "spf13_jwalterweatherman",
    get = "github.com/spf13/jwalterweatherman",
    revision = "v1.0.0",
    licences = ["MIT"],
)

go_get(
    name = "subosito_gotenv",
    get = "github.com/subosito/gotenv",
    revision = "v1.2.0",
    licences = ["MIT"],
)

go_get(
    name = "ini_v1",
    get = "gopkg.in/ini.v1",
    revision = "v1.51.0",
    licences = ["Apache-2.0"],
)


go_get(
    name = "yaml_v2",
    get = "gopkg.in/yaml.v2",
    revision = "v2.2.4",
    licences = ["Apache-2.0"],
)

go_get(
    name = "grpc-ecosystem_grpc-middleware",
    get = "github.com/grpc-ecosystem/go-grpc-middleware",
    install = [
        "",
        "logging/zap",
        "recovery",
    ],
    revision = "v1.2.2",
    deps = [
        ":protobuf",
        ":uber_zap",
        ":x_net",
        ":grpc",
    ],
)