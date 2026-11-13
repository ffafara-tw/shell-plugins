package yarn

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/importer"
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
				MarkdownDescription: "Access token used to authenticate to an npm-compatible registry.",
				Secret:              true,
			},
			{
				Name:                fieldname.Organization,
				MarkdownDescription: "The package scope this registry should be used for, without the leading @.",
				Optional:            true,
			},
			{
				Name:                fieldname.Host,
				MarkdownDescription: "The npm-compatible registry host or URL that accepts this access token.",
				Optional:            true,
			},
		},
		DefaultProvisioner: scopeEnvProvisioner{},
		Importer: importer.TryAll(
			importer.TryAllEnvVars(fieldname.Token, "NPM_TOKEN", "NODE_AUTH_TOKEN", "YARN_AUTH_TOKEN", "YARN_TOKEN"),
			TryyarnConfigFile(),
		)}
}

var nonEnvVarChars = regexp.MustCompile(`[^A-Z0-9]+`)

// envVarName returns the environment variable referenced by npmAuthToken in .yarnrc.yml for the given scope:
// {PREFIX}_NPM_AUTH_TOKEN. Without a scope it falls back to the variable Yarn reads for the default registry.
func envVarName(scope string) string {
	prefix := strings.TrimPrefix(strings.TrimSpace(scope), "@")
	prefix = strings.Trim(nonEnvVarChars.ReplaceAllString(strings.ToUpper(prefix), "_"), "_")
	if prefix == "" {
		prefix = "YARN"
	}
	return prefix + "_NPM_AUTH_TOKEN"
}

type scopeEnvProvisioner struct{}

func (scopeEnvProvisioner) Provision(_ context.Context, in sdk.ProvisionInput, out *sdk.ProvisionOutput) {
	out.AddEnvVar(envVarName(in.ItemFields[fieldname.Organization]), in.ItemFields[fieldname.Token])
}

func (scopeEnvProvisioner) Deprovision(_ context.Context, _ sdk.DeprovisionInput, _ *sdk.DeprovisionOutput) {
	// Nothing to do here: environment variables get wiped automatically when the process exits.
}

func (scopeEnvProvisioner) Description() string {
	return "Provision {SCOPE}_NPM_AUTH_TOKEN environment variable"
}

// TryyarnConfigFile imports an auth token for every npmScope in ~/.yarnrc.yml that defines npmAuthToken.
func TryyarnConfigFile() sdk.Importer {
	return importer.TryFile("~/.yarnrc.yml", func(ctx context.Context, contents importer.FileContents, in sdk.ImportInput, out *sdk.ImportAttempt) {
		var config Config
		if err := contents.ToYAML(&config); err != nil {
			out.AddError(err)
			return
		}

		if candidate, ok := newCandidate("", config.NPMAuthToken, config.NPMRegistryServer); ok {
			out.AddCandidate(candidate)
		}

		scopes := make([]string, 0, len(config.NPMScopes))
		for scope := range config.NPMScopes {
			scopes = append(scopes, scope)
		}
		sort.Strings(scopes)

		for _, scope := range scopes {
			settings := config.NPMScopes[scope]
			if candidate, ok := newCandidate(scope, settings.NPMAuthToken, settings.NPMRegistryServer); ok {
				out.AddCandidate(candidate)
			}
		}
	})
}

func newCandidate(scope, token, registry string) (sdk.ImportCandidate, bool) {
	token = strings.TrimSpace(token)
	// Tokens that reference an environment variable, e.g. "${SOFTWARE_NPM_AUTH_TOKEN:-}", hold no secret.
	if token == "" || strings.Contains(token, "${") {
		return sdk.ImportCandidate{}, false
	}

	fields := map[sdk.FieldName]string{fieldname.Token: token}
	nameHint := "yarn"
	if scope = strings.TrimPrefix(strings.TrimSpace(scope), "@"); scope != "" {
		fields[fieldname.Organization] = scope
		nameHint = fmt.Sprintf("yarn %s", scope)
	}
	if registry = strings.TrimSpace(registry); registry != "" {
		fields[fieldname.Host] = registry
	}

	return sdk.ImportCandidate{Fields: fields, NameHint: importer.SanitizeNameHint(nameHint)}, true
}

// Config is the subset of the .yarnrc.yml schema used to import auth tokens.
type Config struct {
	NPMAuthToken      string                 `yaml:"npmAuthToken"`
	NPMRegistryServer string                 `yaml:"npmRegistryServer"`
	NPMScopes         map[string]ScopeConfig `yaml:"npmScopes"`
}

type ScopeConfig struct {
	NPMAuthToken      string `yaml:"npmAuthToken"`
	NPMRegistryServer string `yaml:"npmRegistryServer"`
}
