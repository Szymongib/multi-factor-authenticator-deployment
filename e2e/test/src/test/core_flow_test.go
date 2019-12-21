package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"testing"

	"github.com/szymongib/multi-factor-authenticator-core/pkg/api/contract"

	"github.com/stretchr/testify/assert"

	"github.com/szymongib/multi-factor-authenticator-core/pkg/client"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/pretty"
)

type authRequestConfig struct {
	authMethodName string
	credentials    interface{}
}

func TestFullFlow(t *testing.T) {

	userCreds := testSuite.GenerateCredentials()
	fmt.Println("Test credentials: ")
	_, err := spew.Println(userCreds)
	require.NoError(t, err)

	defer func() {
		// TODO - delete test user
	}()

	fmt.Printf("\nRegistering %s user...\n", userCreds.Email)
	idResponse, err := testSuite.CoreClient.RegisterUser(userCreds.Email, userCreds.Password)
	require.NoError(t, err)
	printDataAsJSONOrFail(t, idResponse)

	t.Run("should return 409 when email already in use", func(t *testing.T) {
		response, err := testSuite.CoreClient.RegisterUserRaw(userCreds.Email, userCreds.Password)
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, response.StatusCode)
		dumpResponse(t, response)
	})

	// TODO - cleanup after test or generate random creds (best both)

	_, idTokenResponse := generateIdToken(t, userCreds)

	fmt.Println("Testing Authentication Methods...")

	t.Run("Test Password authentication method", func(t *testing.T) {
		securedClient := client.NewIdTokenSecuredAPIClient(config.MultiFactorAuthenticatorCoreURL, idTokenResponse.Token.Token, config.SkipTLSVerify)

		passwordCredentials := testSuite.GeneratePasswordAuthMethodCredentials()
		fmt.Println("Password credentials: ")
		_, err := spew.Println(userCreds)
		require.NoError(t, err)

		fmt.Println("Enabling Password Method...")
		tokenResponse, err := securedClient.EnableAuthenticationMethod(config.PasswordMethod.Name, passwordCredentials)
		require.NoError(t, err)
		printDataAsJSONOrFail(t, tokenResponse)

		t.Run("should fail to authenticate without Password auth method token", func(t *testing.T) {
			fmt.Println("Logging in user...")
			loginResponse, err := testSuite.CoreClient.LoginUser(userCreds.Email, userCreds.Password)
			require.NoError(t, err)
			printDataAsJSONOrFail(t, loginResponse)

			fmt.Println("Trying to generate Id Token without required Password auth token...")
			response, err := testSuite.CoreClient.GenerateIdTokenRaw(loginResponse.Token.Token, nil)
			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, response.StatusCode)
			responseBytes, err := httputil.DumpResponse(response, true)
			require.NoError(t, err)
			fmt.Println("Failed response: ", string(responseBytes))
		})

		fmt.Println("Generating Id Token...")
		passwordAuthRequestConfig := authRequestConfig{
			authMethodName: config.PasswordMethod.Name,
			credentials:    passwordCredentials,
		}
		_, idTokenResponse = generateIdToken(t, userCreds, passwordAuthRequestConfig)
		require.NoError(t, err)
		printDataAsJSONOrFail(t, idTokenResponse)

		fmt.Println("Disabling Password Method...")
		disabled, err := securedClient.DisableAuthenticationMethod(config.PasswordMethod.Name)
		require.NoError(t, err)
		printDataAsJSONOrFail(t, disabled)

		fmt.Println("Assert descriptions present")
		authMethods, err := securedClient.GetUserAuthenticationMethods()
		require.NoError(t, err)
		printDataAsJSONOrFail(t, authMethods)

		for _, method := range authMethods.AuthenticationMethods {
			assert.NotEmpty(t, method.Description)
		}

	})
}

func Test_Errors(t *testing.T) {

	t.Run("should return 401 if user not found", func(t *testing.T) {
		userCreds := testSuite.GenerateCredentials()
		fmt.Println("Test credentials: ")
		_, err := spew.Println(userCreds)
		require.NoError(t, err)

		fmt.Println("Logging in user...")
		loginResponse, err := testSuite.CoreClient.LoginUserRaw(userCreds.Email, userCreds.Password)
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, loginResponse.StatusCode)
	})
}

func generateIdToken(t *testing.T, userCreds contract.UserCredentials, authConfigs ...authRequestConfig) (contract.LoginResponse, contract.IdTokenResponse) {
	fmt.Println("Logging in user...")
	loginResponse, err := testSuite.CoreClient.LoginUser(userCreds.Email, userCreds.Password)
	require.NoError(t, err)
	printDataAsJSONOrFail(t, loginResponse)

	tokens := make(contract.Tokens, len(authConfigs))
	for i, authCfg := range authConfigs {
		fmt.Printf("Logging in to %s Authentication Method...\n", authCfg.authMethodName)
		authTokenResponse, err := testSuite.CoreClient.LoginToAuthenticationService(t, authCfg.authMethodName, loginResponse.Token.Token, authCfg.credentials)
		require.NoError(t, err)
		printDataAsJSONOrFail(t, authTokenResponse)

		tokens[i] = contract.Token{
			AuthenticationMethod: authCfg.authMethodName,
			Token:                authTokenResponse.Token,
		}
	}

	fmt.Println("Generating Id Token...")
	idTokenResponse, err := testSuite.CoreClient.GenerateIdToken(loginResponse.Token.Token, tokens)
	require.NoError(t, err)
	printDataAsJSONOrFail(t, idTokenResponse)

	return loginResponse, idTokenResponse
}

func printDataAsJSONOrFail(t *testing.T, data interface{}) {
	err := printDataAsJSON(data)
	assert.NoError(t, err)
}

func printDataAsJSON(data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	prettyResponse := pretty.Pretty(bytes)

	fmt.Println(string(pretty.Color(prettyResponse, nil)))

	return nil
}

func dumpResponse(t *testing.T, response *http.Response) {
	bytes, err := httputil.DumpResponse(response, true)
	require.NoError(t, err)
	fmt.Println(string(bytes))
}
