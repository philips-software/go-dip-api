package iam

import "testing"

func TestAccessTokenEndpointAudienceOverride(t *testing.T) {
	// Default: audience derived from the base IAM URL.
	c, err := NewClient(nil, &Config{IAMURL: "https://iam.example.com", IDMURL: "https://idm.example.com"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got, want := c.accessTokenEndpoint(), "https://iam.example.com/oauth2/access_token"; got != want {
		t.Errorf("default audience = %q, want %q", got, want)
	}

	// Override: TokenAudience takes precedence (e.g. Keycloak realm issuer).
	aud := "https://iam.us-east-1.stg-foundation-security.hsp.philips.com/realms/healthsuite"
	co, err := NewClient(nil, &Config{
		IAMURL:        "https://iam.example.com",
		IDMURL:        "https://idm.example.com",
		TokenAudience: aud,
	})
	if err != nil {
		t.Fatalf("NewClient with override: %v", err)
	}
	if got := co.accessTokenEndpoint(); got != aud {
		t.Errorf("override audience = %q, want %q", got, aud)
	}
}
