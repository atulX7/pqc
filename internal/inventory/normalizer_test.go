package inventory

import "testing"

func TestInferCryptoUsageTypeForPyCryptodomePatterns(t *testing.T) {
	tests := []struct {
		name        string
		matchedText string
		ruleName    string
		want        string
	}{
		{
			name:        "rsa import",
			matchedText: "from Crypto.PublicKey import RSA",
			ruleName:    "PyCryptodome RSA import",
			want:        "crypto library import",
		},
		{
			name:        "rsa import key",
			matchedText: "cert = RSA.import_key(get_pub_cert())",
			ruleName:    "PyCryptodome RSA key import",
			want:        "key import/loading",
		},
		{
			name:        "private key reference",
			matchedText: `privateKey = RSA.import_key(open("rsa_private_key.pem", "rb").read())`,
			ruleName:    "PyCryptodome RSA key import",
			want:        "key file reference",
		},
		{
			name:        "public key reference",
			matchedText: `publicKey = RSA.import_key(open("rsa_public_key.pem", "rb").read())`,
			ruleName:    "PyCryptodome RSA key import",
			want:        "key file reference",
		},
		{
			name:        "oaep",
			matchedText: "cipher = PKCS1_OAEP.new(publicKey)",
			ruleName:    "PyCryptodome RSA OAEP encryption",
			want:        "encryption/decryption",
		},
		{
			name:        "pkcs signature",
			matchedText: "signature = pkcs1_15.new(privateKey).sign(hashValue)",
			ruleName:    "PyCryptodome RSA PKCS#1 signature",
			want:        "digital signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := InferCryptoUsageType(tt.matchedText, tt.ruleName)
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
