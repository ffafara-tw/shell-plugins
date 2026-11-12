package yarn

import (
	"testing"
	
	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/plugintest"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)
	
func TestAuthTokenProvisioner(t *testing.T) {
	plugintest.TestProvisioner(t, AuthToken().DefaultProvisioner, map[string]plugintest.ProvisionCase{
		"default": {
			ItemFields: map[sdk.FieldName]string{ // TODO: Check if this is correct
				fieldname.Token: "ghp_Qpz9WNHxLIFSzAn029MeXOSOrwhnJEXAMPLE",
			},
			ExpectedOutput: sdk.ProvisionOutput{
				Environment: map[string]string{
					"YARN_TOKEN": "ghp_Qpz9WNHxLIFSzAn029MeXOSOrwhnJEXAMPLE",
				},
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
		// TODO: If you implemented a config file importer, add a test file example in yarn/test-fixtures
		// and fill the necessary details in the test template below.
		"config file": {
			Files: map[string]string{
				// "~/path/to/config.yml": plugintest.LoadFixture(t, "config.yml"),
			},
			ExpectedCandidates: []sdk.ImportCandidate{
			// 	{
			// 		Fields: map[sdk.FieldName]string{
			// 			fieldname.Token: "ghp_Qpz9WNHxLIFSzAn029MeXOSOrwhnJEXAMPLE",
			// 		},
			// 	},
			},
		},
	})
}
