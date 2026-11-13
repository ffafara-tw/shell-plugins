package yarn

import (
	"testing"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"scope": {
			ItemFields: map[sdk.FieldName]string{
				fieldname.Token:        "npm_EXAMPLE",
				fieldname.Organization: "my-software",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{"MY_SOFTWARE_NPM_AUTH_TOKEN": "npm_EXAMPLE"},
			},
		},
		"no scope": {
			ItemFields: map[sdk.FieldName]string{fieldname.Token: "npm_EXAMPLE"},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{"YARN_NPM_AUTH_TOKEN": "npm_EXAMPLE"},
			},
		},
	})
}

func TestAuthTokenImporter(t *testing.T) {
	plugintest.TestImporter(t, AuthToken().Importer, map[string]plugintest.ImportCase{
		"environment": {
			Environment: map[string]string{ // TODO: Check if this is correct
				"YARN_TOKEN": "ghp_Qpz9WNHxLIFSzAn029MeXOSOrwhnJEXAMPLE",
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token: "ghp_Qpz9WNHxLIFSzAn029MeXOSOrwhnJEXAMPLE",
					},
				},
			},
		},
		"config file": {
			Files: map[string]string{
				"~/.yarnrc.yml": plugintest.LoadFixture(t, "yarnrc.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "npm_EXAMPLE_other_token",
						fieldname.Organization: "other",
					},
					NameHint: "yarn other",
				},
				{
					Fields: map[sdk.FieldName]string{
						fieldname.Token:        "npm_EXAMPLE_company_token",
						fieldname.Organization: "company",
						fieldname.Host:         "https://npm.private.registry.com",
					},
					NameHint: "yarn company",
				},
			},
		},
	})
}
