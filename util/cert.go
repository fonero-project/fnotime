// Copyright (c) 2017 The Fonero developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package util

import (
	"crypto/elliptic"
	"io/ioutil"
	"os"
	"time"

	"github.com/fonero-project/fnod/certgen"
)

// GenCertPair generates a key/cert pair to the paths provided.
func GenCertPair(org, certFile, keyFile string) error {
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	cert, key, err := certgen.NewTLSCertPair(elliptic.P521(), org,
		validUntil, nil)
	if err != nil {
		return err
	}

	// Write cert and key files.
	if err = ioutil.WriteFile(certFile, cert, 0666); err != nil {
		return err
	}
	if err = ioutil.WriteFile(keyFile, key, 0600); err != nil {
		os.Remove(certFile)
		return err
	}

	return nil
}
