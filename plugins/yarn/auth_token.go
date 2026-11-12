package yarn

import (
	"context"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
	"github.com/1Password/shell-plugins/sdk/provision"
	"github.com/1Password/shell-plugins/sdk/schema"
	"github.com/1Password/shell-plugins/sdk/schema/credname"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

func AuthToken() schema.CredentialType {
	return schema.CredentialType{
		Name:          credname.AuthToken,
		DocsURL:       sdk.URL("https://yarnpkg.com/configuration/yarnrc#npmAuthToken"),
		ManagementURL: sdk.URL("https://www.npmjs.com/settings/~/tokens"),
		Fields: []schema.CredentialField{
			{
				Name:                fieldname.Token,
				MarkdownDescription: "Token used to authenticate to yarn.",
				Secret:              true,
				Composition: &schema.ValueComposition{
					Length: 40,
					Prefix: "ghp_", // TODO: Check if this is correct
					Charset: schema.Charset{
						Uppercase: true,
						Lowercase: true,
						Digits:    true,
					},
				},
			},
		},
		DefaultProvisioner: provision.EnvVars(defaultEnvVarMapping),
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "NPM_TOKEN", "NODE_AUTH_TOKEN", "YARN_AUTH_TOKEN", "YARN_TOKEN"),
			TryyarnConfigFile(),
		)}
}

var defaultEnvVarMapping = map[string]sdk.FieldName{
	"npmAuthToken": fieldname.Token, // TODO: Check if this is correct
}

// TODO: Check if the platform stores the Auth Token in a local config file, and if so,
// implement the function below to add support for importing it.
func TryyarnConfigFile() sdk.Importer {
	return importer.TryFile("~/.yarnrc.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		// var config Config
		// if err := contents.ToYAML(&config); err != nil {
		// 	out.AddError(err)
		// 	return
		// }

		// if config.Token == "" {
		// 	return
		// }

		// out.AddCandidate(sdk.ImportCandidate{
		// 	Fields: map[sdk.FieldName]string{
		// 		fieldname.Token: config.Token,
		// 	},
		// })
	})
}

// TODO: Implement the config file schema
// type Config struct {
//	Token string
// }
