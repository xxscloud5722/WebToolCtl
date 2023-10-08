module github.com/longyuan/gitlab.v3

go 1.19

replace github.com/longyuan/lib.v3 => ../ALib

replace github.com/longyuan/storage.v3 => ../StorageCTL

require (
	github.com/fatih/color v1.15.0
	github.com/longyuan/lib.v3 v0.0.0-00010101000000-000000000000
	github.com/longyuan/storage.v3 v0.0.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cobra v1.7.0
	github.com/xanzy/go-gitlab v0.86.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mozillazg/go-httpheader v0.4.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tencentyun/cos-go-sdk-v5 v0.7.42 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/oauth2 v0.6.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.29.1 // indirect
)
