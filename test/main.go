package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"os"

	"software.sslmate.com/src/go-pkcs12"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func main() {
	// Load environment variables
	certPath := os.Getenv("TEMP") + "\\msiCert.pfx"
	certPassword := "@Nduff1284" //os.Getenv("CERT_PASSWORD")
	tenantID := "1d122f81-68ca-4aeb-91cd-0d1a6029e07b"
	clientID := "a048a641-afb4-4666-831c-907a4f612331"
	resourceURL := "499b84ac-1321-427f-aa17-267ca6975798"

	// Read certificate file
	certData, err := os.ReadFile(certPath)
	if err != nil {
		fmt.Printf("Failed to read certificate file: %v\n", err)
		return
	}

	// Decode certificate
	privateKey, certificate, err := pkcs12.Decode(certData, certPassword)
	if err != nil {
		fmt.Printf("Failed to parse PFX: %v\n", err)
		return
	}

	// Create credential
	cred, err := azidentity.NewClientCertificateCredential(tenantID, clientID, []*x509.Certificate{certificate}, privateKey, nil)
	if err != nil {
		fmt.Printf("Failed to create credential: %v\n", err)
		return
	}

	// Get access token
	ctx := context.Background()
	token, err := cred.GetToken(ctx, policy.TokenRequestOptions{
		Scopes: []string{resourceURL + "/.default"},
	})
	if err != nil {
		fmt.Printf("Failed to get access token: %v\n", err)
		return
	}

	// Print access token
	fmt.Println(token.Token)
}
