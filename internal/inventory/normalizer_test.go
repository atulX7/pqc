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
		{
			name:        "dotnet private key import",
			matchedText: "rsa.ImportRSAPrivateKey(new ReadOnlySpan<byte>(privKey), &bytesRead)",
			ruleName:    ".NET RSA private key import",
			want:        "private key import",
		},
		{
			name:        "dotnet public key import",
			matchedText: "rsa.ImportSubjectPublicKeyInfo(new ReadOnlySpan<byte>(pubKey), &bytesRead)",
			ruleName:    ".NET RSA public key import",
			want:        "public key import",
		},
		{
			name:        "dotnet jwt signing",
			matchedText: "new SigningCredentials(new RsaSecurityKey(rsa), SecurityAlgorithms.RsaSha256)",
			ruleName:    ".NET JWT SigningCredentials",
			want:        "JWT signing",
		},
		{
			name:        "dotnet jwt token handling",
			matchedText: "let handler = new JwtSecurityTokenHandler()",
			ruleName:    ".NET JWT token handler",
			want:        "JWT token handling",
		},
		{
			name:        "dotnet rsa create",
			matchedText: "use rsa = RSA.Create()",
			ruleName:    ".NET RSA key creation",
			want:        "key generation",
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
