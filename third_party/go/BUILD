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
        "ptypes/any",
        "ptypes/struct",
        "protoc-gen-go/descriptor",
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
        "http2/hpack",
        "html",
        "trace",
        "context",
        "internal/timeseries",
        "http/httpguts",
        "idna",
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
        "transform",
    ],
)

go_get(
    name = "x_sys",
    get = "golang.org/x/sys",
    revision = "release-branch.go1.15",
    install = [
        "unix",
        "internal/unsafeheader",
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
        "attributes",
        "codes",
        "status",
        "backoff",
        "balancer",
        "balancer/base",
        "balancer/grpclb/state",
        "connectivity",
        "internal/...",
        "resolver",
        "keepalive",
        "credentials",
        "serviceconfig",
        "balancer/roundrobin",
        "encoding",
        "encoding/proto",
        "metadata",
        "peer",
        "stats",
        "health",
        "reflection",
        "health/grpc_health_v1",
        "grpclog",
        "tap",
        "binarylog/grpc_binarylog_v1",
        "reflection/grpc_reflection_v1alpha",
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
    install = [
        "",
        "internal/plugin",
    ],
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
    install = [
        "",
        "vendor/github.com/mattn/go-colorable",
        "vendor/github.com/mattn/go-isatty",
        "vendor/golang.org/x/sys/unix",
    ],
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
    install = [
        "",
        "zapcore",
        "buffer",
        "internal/...",
    ],
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
    install = [
        "",
        "hcl/printer",
        "hcl/ast",
        "hcl/parser",
        "hcl/token",
        "hcl/scanner",
        "hcl/strconv",
        "json/scanner",
        "json/token",
        "json/parser",
    ],
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
    install = [
        "",
        "mem",
    ],
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
        "logging",
        "logging/zap/...",
        "recovery",
        "tags",
    ],
    revision = "v1.2.2",
    deps = [
        ":protobuf",
        ":uber_zap",
        ":x_net",
        ":grpc",
    ],
    licences = ["Apache-2.0"],
)

go_get(
    name = "docker",
    get = "github.com/docker/docker",
    repo = "github.com/moby/moby",
    install = [
        "api",
        "api/types/...",
        "client",
        "vendor/github.com/pkg/errors",
        "vendor/github.com/docker/go-connections/sockets",
        "vendor/github.com/docker/distribution/reference",
        "errdefs",
        "vendor/github.com/docker/go-connections/tlsconfig",
        "vendor/github.com/opencontainers/go-digest",
        "vendor/github.com/docker/go-units",
        "vendor/github.com/docker/go-connections/nat",
        "vendor/github.com/opencontainers/image-spec/specs-go/v1",
        "vendor/golang.org/x/net/proxy",
        "vendor/github.com/docker/distribution/digestset",
        "vendor/github.com/containerd/containerd/errdefs",
        "vendor/github.com/docker/distribution/registry/api/errcode",
        "vendor/github.com/sirupsen/logrus",
        "vendor/google.golang.org/grpc/codes",
        "vendor/google.golang.org/grpc/status",
        "vendor/github.com/opencontainers/image-spec/specs-go",
        "vendor/github.com/gogo/protobuf/proto",
        "vendor/golang.org/x/net/internal/socks",
        "vendor/golang.org/x/sys/unix",
        "vendor/github.com/golang/protobuf/proto",
        "vendor/github.com/golang/protobuf/ptypes",
        "vendor/google.golang.org/genproto/googleapis/rpc/status",
        "vendor/google.golang.org/grpc/internal",
        "vendor/github.com/golang/protobuf/ptypes/any",
        "vendor/github.com/golang/protobuf/ptypes/duration",
        "vendor/github.com/golang/protobuf/ptypes/timestamp",
        "vendor/google.golang.org/grpc/connectivity",
        "vendor/google.golang.org/grpc/grpclog",
    ],
    revision = "v19.03.13",
    licences = ["Apache-2.0"],
)

go_get(
    name = "google_cadvisor",
    get = "github.com/google/cadvisor",
    install = [
        "fs",
        "devicemapper",
        "utils",
    ],
    revision = "v0.37.0",
    licences = ["Apache-2.0"],
    deps = [
        ":mistifyio_zfs",
        ":k8s_klog",
        ":k8s_utils",
    ]
)

go_get(
    name = "mistifyio_zfs",
    get = "github.com/mistifyio/go-zfs",
    revision = "v1.0.0",
    licences = ["Apache-2.0"],
)

go_get(
    name = "k8s_klog",
    get = "k8s.io/klog",
    repo = "github.com/kubernetes/klog",
    revision = "v2.2.0",
    licences = ["Apache-2.0"],
    deps = [":logr"],
)

go_get(
    name = "k8s_utils",
    get = "k8s.io/utils",
    repo = "github.com/kubernetes/utils",
    install = [
        "io",
        "mount",
        "exec",
    ],
    revision = "2df71ebbae66",
    licences = ["Apache-2.0"],
    deps = [":k8s_klog"],
)

go_get(
    name = "logr",
    get = "github.com/go-logr/logr",
    revision = "v0.2.1",
    licences = ["Apache-2.0"],
)

go_get(
    name = "k8s_cri-api",
    get = "k8s.io/cri-api",
    repo = "github.com/kubernetes/cri-api",
    install = [
        "pkg/apis/runtime/v1alpha2",
    ],
    revision = "v0.19.2",
    licences = ["Apache-2.0"],
    deps = [
        ":grpc",
        ":gogo_protobuf",
    ],
)

go_get(
    name = "gogo_protobuf",
    get = "github.com/gogo/protobuf",
    install = [
        "gogoproto",
        "proto",
        "sortkeys",
        "protoc-gen-gogo/descriptor",
    ],
    revision = "v1.3.1",
    licences = ["BSD-3-Clause"],
)

go_get(
    name = "docker_distribution",
    get = "github.com/docker/distribution",
    install = [
        "reference",
        "digestset",
        "vendor/github.com/opencontainers/go-digest",
    ],
    revision = "v2.7.1",
    licences = ["Apache-2.0"],
)

go_get(
    name = "dgrijalva_jwt",
    get = "github.com/dgrijalva/jwt-go",
    revision = "v3.2.0",
    licences = ["MIT"],
)

go_get(
    name = "mozilla_pkcs7",
    get = "go.mozilla.org/pkcs7",
    repo = "github.com/mozilla-services/pkcs7",
    revision = "432b2356ecb18209c1cec25680b8a23632794f21",
    licences = ["MIT"],
)

go_get(
    name = "google_containerregistry",
    get = "github.com/google/go-containerregistry",
    install = [
        "vendor/github.com/docker/cli/cli/config",
        "vendor/github.com/docker/cli/cli/config/types",
        "vendor/github.com/docker/cli/cli/config/configfile",
        "vendor/github.com/docker/cli/cli/config/credentials",
        "vendor/github.com/docker/docker/pkg/homedir",
        "vendor/github.com/docker/docker-credential-helpers/client",
        "vendor/github.com/docker/docker-credential-helpers/credentials",
        "vendor/github.com/pkg/errors",
        "pkg/authn",
        "pkg/name",
        "pkg/logs",
        "pkg/internal/retry",
        "pkg/internal/retry/wait",
        "pkg/v1/remote/transport",
    ],
    revision = "v0.1.3",
    licences = ["MIT"],
)

go_get(
    name = "stretchr_testify",
    get = "github.com/stretchr/testify",
    install = [
        "assert", 
        "require",
    ],
    revision = "v1.6.1",
    licences = ["MIT"],
    deps = [
        ":davecgh_spew",
        ":pmezard_difflib",
        ":yaml_v3",
    ],
)

go_get(
    name = "davecgh_spew",
    get = "github.com/davecgh/go-spew",
    install = ["spew"],
    revision = "v1.1.1",
    licences = ["MIT"],
)

go_get(
    name = "pmezard_difflib",
    get = "github.com/pmezard/go-difflib",
    install = ["difflib"],
    revision = "v1.0.0",
    licences = ["ISC"],
)

go_get(
    name = "yaml_v3",
    get = "gopkg.in/yaml.v3",
    revision = "v3",
    licences = ["Apache-2.0", "MIT"],
)
