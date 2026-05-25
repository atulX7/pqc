package inventory

type Finding struct {
	SourceType  string `json:"source_type"`
	FilePath    string `json:"file_path"`
	LineNumber  int    `json:"line_number"`
	MatchedText string `json:"matched_text"`
	RuleID      string `json:"rule_id"`
	RuleName    string `json:"rule_name"`
	Algorithm   string `json:"algorithm"`
	Severity    string `json:"severity"`
	RiskType    string `json:"risk_type"`
	Confidence  string `json:"confidence"`
	Priority    int    `json:"-"`
}

type BusinessMetadata struct {
	ApplicationName       string
	SensitivityFlags      []string
	ExposureLevel         string
	BusinessCriticality   string
	SecrecyLifetimeYears  int
	CryptoAgilityLevel    string
	VendorDependencyLevel string
	MigrationComplexity   string
}

type CryptoAsset struct {
	ApplicationName       string   `json:"application_name"`
	SourceType            string   `json:"source_type"`
	FilePath              string   `json:"file_path"`
	LineNumber            int      `json:"line_number"`
	CryptoUsageType       string   `json:"crypto_usage_type"`
	Algorithm             string   `json:"algorithm"`
	Severity              string   `json:"severity"`
	RiskType              string   `json:"risk_type"`
	Confidence            string   `json:"confidence"`
	MatchedText           string   `json:"matched_text"`
	SensitivityFlags      []string `json:"sensitivity_flags"`
	ExposureLevel         string   `json:"exposure_level"`
	BusinessCriticality   string   `json:"business_criticality"`
	SecrecyLifetimeYears  int      `json:"secrecy_lifetime_years"`
	CryptoAgilityLevel    string   `json:"crypto_agility_level"`
	VendorDependencyLevel string   `json:"vendor_dependency_level"`
	MigrationComplexity   string   `json:"migration_complexity"`
	QuantumRiskScore      float64  `json:"quantum_risk_score"`
	RiskLevel             string   `json:"risk_level"`
	Recommendation        string   `json:"recommendation"`
}
