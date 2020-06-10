package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/go-cleanhttp"
	"golang.org/x/oauth2"
)

const (
	TenantID     = "53109908-9db2-4dbd-ab3e-9c40ab19bac9"
	ClientID     = "453b11f3-cc26-422c-b466-2c6ddd60f6eb"
	ClientSecret = ""
)

type azureProvider struct {
	oidcVerifier *oidc.IDTokenVerifier
	httpClient   *http.Client
}

type oidcDiscoveryInfo struct {
	Issuer  string `json:"issuer"`
	JWKSURL string `json:"jwks_uri"`
}

func NewAzureProvider() (*azureProvider, error) {
	httpClient := cleanhttp.DefaultClient()
	discoveryURL := "https://login.microsoftonline.com/localgotodo.onmicrosoft.com/v2.0/.well-known/openid-configuration"
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

	return &azureProvider{
		oidcVerifier: oidcVerifier,
		httpClient:   httpClient,
	}, nil
}
