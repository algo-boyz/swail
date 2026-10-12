package main

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

// reference: https://v3.wails.io/guides/updater/#writing-your-own-provider

func parsePrivateKey(data []byte) (ed25519.PrivateKey, error) {
	// 1. Raw 64-byte key
	if len(data) == ed25519.PrivateKeySize {
		return ed25519.PrivateKey(data), nil
	}

	// 2. OpenSSH private key (e.g. created with ssh-keygen)
	if rawKey, err := ssh.ParseRawPrivateKey(data); err == nil {
		if key, ok := rawKey.(*ed25519.PrivateKey); ok {
			return *key, nil
		}
		if key, ok := rawKey.(ed25519.PrivateKey); ok {
			return key, nil
		}
	}

	// 3. PKCS#8 PEM block (e.g. created with openssl)
	if block, _ := pem.Decode(data); block != nil {
		if parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
			if key, ok := parsedKey.(ed25519.PrivateKey); ok {
				return key, nil
			}
		}
	}

	return nil, fmt.Errorf("unsupported key format or invalid length (%d bytes)", len(data))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run bin/sha_sign.go <path-to-file>")
		os.Exit(1)
	}

	keyData, err := os.ReadFile("updater-key")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading 'updater-key': %v\n", err)
		os.Exit(1)
	}

	privKey, err := parsePrivateKey(keyData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing private key: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening asset: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		fmt.Fprintf(os.Stderr, "Error hashing file: %v\n", err)
		os.Exit(1)
	}

	sig := ed25519.Sign(privKey, h.Sum(nil))
	fmt.Println(base64.StdEncoding.EncodeToString(sig))
}
