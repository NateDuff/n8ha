package ha

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Credentials is the chained token credential used for authentication
var authRecordPath = os.TempDir() + "/auth.json"

// contextKey is a custom type for context keys
type contextKey string

const (
	recordKey contextKey = "record"
	tokenKey  contextKey = "token"
)

var (
	GetChainedCredentialFunc         = GetChainedCredential
	GetBlobCredentialFunc            = GetBlobCredential
	o                                sync.Once
	GetEnvFunc                       = os.Getenv
	NewTableClientFunc               = aztables.NewClient
	NewBlobClientFunc                = azblob.NewClient
	NewManagedIdentityCredentialFunc = azidentity.NewManagedIdentityCredential
	NewEnvironmentCredentialFunc     = azidentity.NewEnvironmentCredential
	NewChainedTokenCredentialFunc    = azidentity.NewChainedTokenCredential
	Credentials                      *azidentity.ChainedTokenCredential
	BlobClient                       *azblob.Client
	InternalTableClient              *aztables.Client
)

// GetChainedCredential creates a chained token credential based on the value of isManagedIdentityEnabled
func GetChainedCredential() {
	o.Do(func() {
		creds := &[]azcore.TokenCredential{}

		isManagedIdentityEnabled := GetEnvFunc("MSI_CLIENT_ID") != ""

		if isManagedIdentityEnabled {
			clientId := azidentity.ClientID(GetEnvFunc("MSI_CLIENT_ID"))
			opts := &azidentity.ManagedIdentityCredentialOptions{ID: clientId}

			managed, err := NewManagedIdentityCredentialFunc(opts)

			if err != nil {
				fmt.Println("Error creating managed identity credential:", err)

			}

			*creds = []azcore.TokenCredential{
				managed,
			}
		} else {
			envVarCreds, err := NewEnvironmentCredentialFunc(nil)

			if err != nil {
				fmt.Println("Error creating environment credential:", err)
			}

			*creds = []azcore.TokenCredential{
				envVarCreds,
			}
		}

		var err error
		Credentials, err = NewChainedTokenCredentialFunc(*creds, nil)
		if err != nil {
			fmt.Println("Error creating chained token credential:", err)
		}
	})
}

func GetBlobCredential() (*azblob.Client, error) {
	if BlobClient != nil {
		return BlobClient, nil
	}

	GetChainedCredentialFunc()

	accountURL := "https://" + os.Getenv("STORAGE_ACCOUNT_NAME") + ".blob.core.windows.net"

	BlobClient, err := NewBlobClientFunc(accountURL, Credentials, nil)
	if err != nil {
		fmt.Println("Error creating blob client:", err)
		return nil, err
	}

	if BlobClient == nil {
		fmt.Println("Blob client is nil")
		return nil, fmt.Errorf("Blob client is nil")
	}

	return BlobClient, nil
}

func GetInternalTableCredential(tableName string) *aztables.Client {
	if InternalTableClient != nil {
		return InternalTableClient
	}

	GetChainedCredentialFunc()

	InternalTableClient, err := aztables.NewClient("https://stdcsinternal.table.core.windows.net/"+tableName, Credentials, nil)
	if err != nil {
		fmt.Println("Error creating table client:", err)
	}

	return InternalTableClient
}

// getAzAccessToken gets an access token for the API using the Credentials variable
func getAzAccessToken(ctx context.Context, creds azcore.TokenCredential) (azcore.AccessToken, error) {
	t, err := creds.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{"https://management.azure.com/.default"},
	})
	if err != nil {
		return azcore.AccessToken{}, err
	}

	// Return the token
	return t, nil
}

// GetUserName gets the username of the currently logged in user
func GetUserName(ctx context.Context) (string, error) {
	record := ctx.Value(recordKey)

	if record != nil {
		return record.(azidentity.AuthenticationRecord).Username, nil
	}

	token, ok := ctx.Value(tokenKey).(string)
	if !ok {
		return "", fmt.Errorf("token not found in context")
	}

	// Convert the string token JWT to a map
	var tokenMap map[string]interface{}
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) < 2 {
		return "", fmt.Errorf("invalid token format")
	}

	payload, err := base64.RawURLEncoding.DecodeString(tokenParts[1])
	if err != nil {
		return "", fmt.Errorf("error decoding token payload: %v", err)
	}

	err = json.Unmarshal(payload, &tokenMap)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling token payload: %v", err)
	}

	userID, ok := tokenMap["unique_name"].(string)
	if !ok {
		userID, _ = tokenMap["appid"].(string)
	}

	return userID, nil
}

// retrieveRecord retrieves the authentication record from a file
func retrieveRecord() (azidentity.AuthenticationRecord, error) {
	record := azidentity.AuthenticationRecord{}
	b, err := os.ReadFile(authRecordPath)
	if err == nil {
		err = json.Unmarshal(b, &record)
	}
	return record, err
}

// storeRecord stores the authentication record to a file
func storeRecord(record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	if err == nil {
		err = os.WriteFile(authRecordPath, b, 0600)
	}
	return err
}

// WithUserAuth returns a context with user authentication
func WithUserAuth(ctx context.Context) (context.Context, error) {
	record, err := retrieveRecord()
	if err != nil {
		// If there is an error reading the record, assume it doesn't exist
		record = azidentity.AuthenticationRecord{}
	}
	c, err := cache.New(nil)
	if err != nil {
		return nil, fmt.Errorf("Error creating cache: %v", err)
	}

	creds, err := azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
		TenantID: os.Getenv("AZURE_TENANT_ID"),
		// If record is zero, the credential will start with no user logged in
		AuthenticationRecord: record,
		// Credentials cache in memory by default. Setting Cache with a
		// nonzero value from cache.New() enables persistent caching.
		Cache: c,
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating interactive browser credential: %v", err)
	}

	if record == (azidentity.AuthenticationRecord{}) {
		// No stored record; call Authenticate to acquire one.
		// This will prompt the user to authenticate interactively.
		record, err = creds.Authenticate(ctx, &policy.TokenRequestOptions{
			TenantID: os.Getenv("AZURE_TENANT_ID"),
		})
		if err != nil {
			return nil, fmt.Errorf("Error authenticating user: %v", err)
		}
		err = storeRecord(record)
		if err != nil {
			return nil, fmt.Errorf("Error storing authentication record: %v", err)
		}
	}

	log.Printf("Authenticated as %s", record.Username)

	// Add the record to context
	ctx = context.WithValue(ctx, recordKey, record)

	Credentials, err = azidentity.NewChainedTokenCredential([]azcore.TokenCredential{creds}, nil)

	return ctx, nil
}

// WithMSIAuth returns a context with MSI authentication
func WithMSIAuth(ctx context.Context) (context.Context, error) {
	isManagedIdentityEnabled := os.Getenv("MSI_CLIENT_ID") != ""

	var creds []azcore.TokenCredential

	if isManagedIdentityEnabled {
		clientID := azidentity.ClientID(os.Getenv("MSI_CLIENT_ID"))
		opts := &azidentity.ManagedIdentityCredentialOptions{ID: clientID}

		managed, err := azidentity.NewManagedIdentityCredential(opts)

		if err != nil {
			fmt.Println("Error creating managed identity credential:", err)

		}

		creds = []azcore.TokenCredential{
			managed,
		}
	} else {
		envVarCreds, err := azidentity.NewEnvironmentCredential(nil)

		if err != nil {
			fmt.Println("Error creating environment credential:", err)
		}

		creds = []azcore.TokenCredential{
			envVarCreds,
		}
	}

	Credentials, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		fmt.Println("Error creating chained token credential:", err)
	}

	token, err := getAzAccessToken(ctx, Credentials)
	if err != nil {
		return nil, fmt.Errorf("Error getting access token: %v", err)
	}

	ctx = context.WithValue(ctx, tokenKey, token.Token)

	return ctx, nil
}

// WithAuthContext returns a context with the appropriate authentication method
func WithAuthContext(ctx context.Context, daemon bool) (context.Context, error) {
	if daemon {
		return WithMSIAuth(ctx)
	}

	return WithUserAuth(ctx)
}
