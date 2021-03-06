package main

import (
	"github.com/Fuco1/waypoint-dobi/builder"
	// "github.com/Fuco1/waypoint-dobi/platform"
	// "github.com/Fuco1/waypoint-dobi/registry"
	// "github.com/Fuco1/waypoint-dobi/release"
	sdk "github.com/hashicorp/waypoint-plugin-sdk"
)

func main() {
	// sdk.Main allows you to register the components which should
	// be included in your plugin
	// Main sets up all the go-plugin requirements

	sdk.Main(sdk.WithComponents(
		// Comment out any components which are not
		// required for your plugin
		&builder.Builder{},
		// &registry.Registry{},
		// &platform.Platform{},
		// &release.ReleaseManager{},
	))
}
