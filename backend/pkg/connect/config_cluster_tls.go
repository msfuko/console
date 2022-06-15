// Copyright 2022 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package connect

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// ConfigClusterTLS is the config if you want to connect to Kafka connect REST API via (mutual) TLS
type ConfigClusterTLS struct {
	Enabled               bool   `yaml:"enabled" json:"enabled"`
	CaFilepath            string `yaml:"caFilepath" json:"caFilepath"`
	CertFilepath          string `yaml:"certFilepath" json:"certFilepath"`
	KeyFilepath           string `yaml:"keyFilepath" json:"keyFilepath"`
	InsecureSkipTLSVerify bool   `yaml:"insecureSkipTlsVerify" json:"insecureSkipTLSVerify"`
}

func (c *ConfigClusterTLS) SetDefaults() {
	c.Enabled = false
}

func (c *ConfigClusterTLS) Validate() error {
	if !c.Enabled {
		return nil
	}

	return nil
}

func (c *ConfigClusterTLS) TLSConfig() (*tls.Config, error) {
	if !c.Enabled {
		return &tls.Config{}, nil
	}

	// 1. Create CA cert pool
	caCertPool := x509.NewCertPool()
	if c.CaFilepath != "" {
		ca, err := ioutil.ReadFile(c.CaFilepath)
		if err != nil {
			return nil, err
		}
		isSuccessful := caCertPool.AppendCertsFromPEM(ca)
		if !isSuccessful {
			return nil, fmt.Errorf("failed to append ca file to cert pool, is this a valid PEM format?")
		}
	}

	// 2. If configured load TLS cert & key - Mutual TLS
	var certificates []tls.Certificate
	if c.CertFilepath != "" && c.KeyFilepath != "" {
		cert, err := ioutil.ReadFile(c.CertFilepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file for schema registry client: %w", err)
		}

		privateKey, err := ioutil.ReadFile(c.KeyFilepath)
		if err != nil {
			return nil, fmt.Errorf("failed to read key file for schema registry client: %w", err)
		}

		pemBlock, _ := pem.Decode(privateKey)
		if pemBlock == nil {
			return nil, fmt.Errorf("no valid private key found")
		}

		tlsCert, err := tls.X509KeyPair(cert, privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate pair for schema registry client: %w", err)
		}
		certificates = []tls.Certificate{tlsCert}
	}

	return &tls.Config{
		InsecureSkipVerify: c.InsecureSkipTLSVerify,
		Certificates:       certificates,
		RootCAs:            caCertPool,
	}, nil
}
