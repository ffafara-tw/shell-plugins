package yarn

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema"
)

func New() schema.Plugin {
	return schema.Plugin{
		Name: "yarn",
		Platform: schema.PlatformInfo{
			Name:     "yarn",
			Homepage: sdk.URL("https://yarn.com"), // TODO: Check if this is correct
		},
		Credentials: []schema.CredentialType{
			AuthToken(),
		},
		Executables: []schema.Executable{
			yarnCLI(),
		},
	}
}
