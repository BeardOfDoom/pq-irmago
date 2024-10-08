package myirmaserver

import (
	"path/filepath"
	"testing"

	irma "github.com/BeardOfDoom/pq-irmago"
	"github.com/BeardOfDoom/pq-irmago/internal/test"
	"github.com/BeardOfDoom/pq-irmago/server"
	"github.com/stretchr/testify/assert"
)

func validConf(t *testing.T) *Configuration {
	testdataPath := test.FindTestdataFolder(t)
	return &Configuration{
		Configuration: &server.Configuration{
			SchemesPath: filepath.Join(testdataPath, "irma_configuration"),
			Logger:      irma.Logger,
		},
		DBType:             DBTypeMemory,
		SessionLifetime:    60,
		KeyshareAttributes: []irma.AttributeTypeIdentifier{irma.NewAttributeTypeIdentifier("test.test.mijnirma.email")},
		EmailAttributes:    []irma.AttributeTypeIdentifier{irma.NewAttributeTypeIdentifier("test.test.email.email")},
	}
}

func TestConfValidation(t *testing.T) {
	_, err := New(validConf(t))
	assert.NoError(t, err)

	conf := validConf(t)
	conf.SessionLifetime = 0
	_, err = New(conf)
	assert.NoError(t, err)

	conf = validConf(t)
	conf.KeyshareAttributes = nil
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.EmailAttributes = nil
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.DBType = "UNKNOWN"
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.KeyshareAttributes = append(conf.KeyshareAttributes, irma.NewAttributeTypeIdentifier("test.test.foo.bar"))
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.KeyshareAttributes = append(conf.EmailAttributes, irma.NewAttributeTypeIdentifier("test.test.foo.bar"))
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.CORSAllowedOrigins = []string{"*"}
	_, err = New(conf)
	assert.NoError(t, err)

	conf = validConf(t)
	conf.CORSAllowedOrigins = []string{"http://example.com"}
	_, err = New(conf)
	assert.NoError(t, err)

	conf = validConf(t)
	conf.CORSAllowedOrigins = []string{"*", "http://example.com"}
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.CORSAllowedOrigins = []string{"example.com"}
	_, err = New(conf)
	assert.Error(t, err)

	conf = validConf(t)
	conf.CORSAllowedOrigins = []string{"http://example.com/foobar"}
	_, err = New(conf)
	assert.Error(t, err)
}
