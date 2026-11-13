package yarn

import (
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/needsauth"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
)

func yarnCLI() schema.Executable {
	return schema.Executable{
		Name:      "yarn CLI",
		Runs:      []string{"yarn"},
		DocsURL:   sdk.URL("https://yarnpkg.com/cli"),
		NeedsAuth: needsauth.IfAll(
			needsauth.NotForHelpOrVersion(),
			needsauth.NotWithoutArgs(),
			needsauth.IfAny(
				needsauth.ForCommand("install"),
			),
		),
		Uses: []schema.CredentialUsage{
			{
				Name: credname.AuthToken,
			},
		},
	}
}
