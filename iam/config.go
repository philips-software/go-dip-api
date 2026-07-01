package iam

import (
	"io"

	hsdpsigner "github.com/philips-software/go-nih-signer"
)

// Config contains the configuration of a client
type Config struct {
	Region           string
	Environment      string
	OAuth2ClientID   string
	OAuth2Secret     string
	SharedKey        string
	SecretKey        string
	BaseIAMURL       string
	BaseIDMURL       string
	OrgAdminUsername string
	OrgAdminPassword string
	IAMURL           string
	IDMURL           string
	// TokenAudience overrides the "aud" claim of the JWT used for service
	// login. When empty the audience defaults to the IAM access token
	// endpoint. Some IAM deployments (e.g. Keycloak-backed realms) expect the
	// realm issuer as the audience instead.
	TokenAudience string
	Scopes        []string
	RootOrgID        string
	DebugLog         io.Writer
	Signer           *hsdpsigner.Signer
}
