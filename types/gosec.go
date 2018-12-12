package types

// GosecOutput is the struct that holds all data from Gosec output.
type GosecOutput struct {
	GosecIssues []GosecIssue `json:"Issues"`
	GosecStats  GosecStats   `json:"Stats"`
}

// GosecIssue is the struct that holds all issues from Gosec output.
type GosecIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GosecStats is the struct that holds all stats from Gosec output.
type GosecStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}
