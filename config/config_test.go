package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/philips-software/go-dip-api/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	iamService := c.
		Region("us-east").
		Env("client-test").
		Service("iam")
	if !assert.NotNil(t, iamService) {
		return
	}
	assert.Equal(t, "https://iam-client-test.us-east.philips-healthsuite.com", iamService.URL)
}

func TestPreview(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	iamService := c.
		Region("us-east").
		Env("preview").
		Service("iam")
	if !assert.NotNil(t, iamService) {
		return
	}
	assert.Equal(t, "https://iam.us-east-1.stg-foundation-security.hsp.philips.com", iamService.URL)

	idmService := c.
		Region("us-east").
		Env("preview").
		Service("idm")
	if !assert.NotNil(t, idmService) {
		return
	}
	assert.Equal(t, "https://idm.us-east-1.stg-foundation-security.hsp.philips.com", idmService.URL)
}

func TestSTL(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}

	stlService := c.
		Region("us-east").
		Service("stl")
	if !assert.NotNil(t, stlService) {
		return
	}
	assert.Equal(t, "na1.vpn.hsdp.io", stlService.Domain)
}

func TestOpts(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !assert.True(t, ok) {
		return
	}
	basePath := filepath.Dir(filename)
	hsdpJsonFile := filepath.Join(basePath, "hsdp.json")
	data, err := os.ReadFile(hsdpJsonFile)
	if !assert.Nil(t, err) {
		return
	}
	configReader := bytes.NewReader(data)
	c, err := config.New(
		config.WithEnv("client-test"),
		config.WithRegion("us-east"),
		config.FromReader(configReader))
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	iam := c.Service("iam")
	assert.Equal(t, "https://iam-client-test.us-east.philips-healthsuite.com", iam.URL)
}

func TestMissing(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	missingService := c.
		Region("us-east").
		Service("bogus")
	assert.Equal(t, "", missingService.URL)
}

func TestRegions(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	regions := c.Regions()
	assert.Less(t, 0, len(regions))
	assert.Contains(t, regions, "eu-west")
}

func TestServices(t *testing.T) {
	c, err := config.New()
	if !assert.Nil(t, err) {
		return
	}
	if !assert.NotNil(t, c) {
		return
	}
	services := c.Region("us-east").Env("client-test").Services()
	assert.Less(t, 0, len(services))
	assert.Contains(t, services, "idm")
	assert.Contains(t, services, "iam")
}
