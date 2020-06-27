package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
)

const (
	// TenantID is the ID of the Azure Tenant
	TenantID = "53109908-9db2-4dbd-ab3e-9c40ab19bac9"
	// ClientID is the ID of the Azure Client
	ClientID = "1832d398-74c9-40e3-8132-d8b45b995035"
)

// AzureProvider contains information to verify tokens
type AzureProvider struct {
	oidcVerifier *oidc.IDTokenVerifier
	httpClient   *http.Client
}

type oidcDiscoveryInfo struct {
	Issuer  string `json:"issuer"`
	JWKSURL string `json:"jwks_uri"`
}

// NewAzureProvider returns struct azureProvider to verify tokens
func NewAzureProvider() (*AzureProvider, error) {
	httpClient := cleanhttp.DefaultClient()
	discoveryURL := "https://localgotodo.b2clogin.com/localgotodo.onmicrosoft.com/b2c_1_susi/v2.0/.well-known/openid-configuration"
	req, err := http.NewRequest("GET", discoveryURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: {{err}}", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s : %s", resp.Status, body)
	}

	var discoveryInfo oidcDiscoveryInfo
	if err := json.Unmarshal(body, &discoveryInfo); err != nil {
		return nil, fmt.Errorf("unable to unmarshal discovery url: {{err}}", err)
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, httpClient)
	remoteKeySet := oidc.NewRemoteKeySet(ctx, discoveryInfo.JWKSURL)

	verifierConfig := &oidc.Config{
		ClientID:             ClientID,
		SupportedSigningAlgs: []string{oidc.RS256},
	}

	oidcVerifier := oidc.NewVerifier(discoveryInfo.Issuer, remoteKeySet, verifierConfig)

	return &AzureProvider{
		oidcVerifier: oidcVerifier,
		httpClient:   httpClient,
	}, nil
}

// HandleTokenVerification verifies the authenticity of the ID Token
func HandleTokenVerification(r *http.Request) (*oidc.IDToken, error) {
	ctx := context.Background()
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		return nil, nil
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	//fmt.Printf(reqToken)
	azureP, err := NewAzureProvider()
	if err != nil {
		return nil, fmt.Errorf("Failed to get Azure Provider: " + err.Error())
	}
	idToken, err := azureP.oidcVerifier.Verify(ctx, reqToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify ID Token: " + err.Error())
	}
	return idToken, nil
}
