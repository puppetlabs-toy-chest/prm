module github.com/puppetlabs/prm

go 1.16

replace github.com/puppetlabs/prm/docs/md => ./docs/md

require (
	github.com/Masterminds/semver v1.5.0
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/Microsoft/hcsshim v0.9.2 // indirect
	github.com/containerd/cgroups v1.0.2 // indirect
	github.com/containerd/containerd v1.5.9 // indirect
	github.com/docker/docker v20.10.12+incompatible
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/go-version v1.4.0
	github.com/json-iterator/go v1.1.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.3
	github.com/moby/sys/mount v0.3.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/muesli/termenv v0.10.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/opencontainers/image-spec v1.0.2
	github.com/opencontainers/runc v1.1.0 // indirect
	github.com/otiai10/copy v1.7.0
	github.com/puppetlabs/pdkgo v0.0.0-20220126155800-c78bfaf4fb87
	github.com/puppetlabs/prm/docs/md v0.0.0-20220131145518-8c57e81fbea0
	github.com/rs/zerolog v1.26.1
	github.com/spf13/afero v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/proto/otlp v0.12.0 // indirect
	golang.org/x/net v0.0.0-20220127200216-cd36cc0744dd // indirect
	golang.org/x/sys v0.0.0-20220128215802-99c3d69c2c27 // indirect
	google.golang.org/genproto v0.0.0-20220126215142-9970aeb2e350 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/ini.v1 v1.66.3 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)
