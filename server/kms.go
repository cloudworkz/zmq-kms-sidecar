package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	cloudkms "cloud.google.com/go/kms/apiv1"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

func getKeyName(projectID string, keyRingID string, locationID string, cryptoKeyID string, cryptoKeyVersion string) string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s/cryptoKeyVersions/%s", projectID, locationID, keyRingID, cryptoKeyID, cryptoKeyVersion)
}

func getAsymmetricPublicKey(keyName string) (interface{}, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &kmspb.GetPublicKeyRequest{
		Name: keyName,
	}
	response, err := client.GetPublicKey(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public key: %+v", err)
	}

	keyBytes := []byte(response.Pem)
	block, _ := pem.Decode(keyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %+v", err)
	}
	return publicKey, nil
}

func encryptRSA(abstractKey interface{}, plaintext []byte) ([]byte, error) {
	rsaKey, ok := abstractKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA")
	}

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaKey, plaintext, nil)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %+v", err)
	}
	return ciphertext, nil
}

func decryptRSA(keyName string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &kmspb.AsymmetricDecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}

	response, err := client.AsymmetricDecrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decryption request failed: %+v", err)
	}
	return response.Plaintext, nil
}
