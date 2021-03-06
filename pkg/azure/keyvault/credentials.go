package keyvault

import (
	"os"
	"strings"

	"get.porter.sh/plugin/azure/pkg/azure/azureconfig"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest"
	azureauth "github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
)

// GetCredentials gets an authorizer for Azure
func GetCredentials(cfg azureconfig.Config, l hclog.Logger) (autorest.Authorizer, error) {

	azureAuthEnvVarNames := []string{
		azureauth.TenantID,
		azureauth.ClientID,
		azureauth.ClientSecret,
		azureauth.CertificatePath,
		azureauth.CertificatePassword,
		azureauth.Username,
		azureauth.Password,
	}

	prefix := cfg.EnvAzurePrefix
	if prefix != "" && prefix != "AZURE_" {
		for _, v := range azureAuthEnvVarNames {
			env := prefix + strings.TrimPrefix(v, "AZURE_")
			val := os.Getenv(env)
			os.Setenv(v, val)
		}
	}

	var authorizer autorest.Authorizer
	var err error

	// Attempt to login with az cli if no vars are set.

	if noAzureAuthEnvVarsAreSet(azureAuthEnvVarNames) {
		authorizer, err = auth.NewAuthorizerFromCLI()
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create an azure authorizer from azure cli")
		}

		return authorizer, nil
	}

	// NewAuthorizierFromEnvironment attempts to authenticate using credentials, certicates, user name and password and MSI however if we get here MSI login wll be skipped as env vars are set so one of the other methods will be attempted

	authorizer, err = auth.NewAuthorizerFromEnvironment()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create an azure authorizer from environment")
	}

	return authorizer, nil
}

func noAzureAuthEnvVarsAreSet(azureAuthEnvVarNames []string) bool {
	for _, v := range azureAuthEnvVarNames {
		val := os.Getenv(v)
		if len(val) > 0 {
			return false
		}
	}
	return true
}
