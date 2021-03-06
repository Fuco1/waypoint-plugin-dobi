module github.com/Fuco1/waypoint-dobi

go 1.14

require (
	github.com/golang/protobuf v1.4.3
	github.com/hashicorp/go-hclog v0.14.1
	github.com/hashicorp/waypoint-plugin-sdk v0.0.0-20201021094150-1b1044b1478e
	github.com/mitchellh/go-glint v0.0.0-20201015034436-f80573c636de
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.0.1 // indirect
	google.golang.org/protobuf v1.25.0
)

// replace github.com/hashicorp/waypoint-plugin-sdk => ../../waypoint-plugin-sdk
