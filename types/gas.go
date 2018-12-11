package types

// GasOutput is the struct that holds all data from Gas output.
type GasOutput struct {
	Issues []GasIssue `json:"Issues"`
	Stats  GasStats   `json:"Stats"`
}

// GasIssue is the struct that holds all issues from Gas output.
type GasIssue struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	RuleID     string `json:"rule_id"`
	Details    string `json:"details"`
	File       string `json:"file"`
	Code       string `json:"code"`
	Line       string `json:"line"`
}

// GasStats is the struct that holds all stats from Gas output.
type GasStats struct {
	Files int `json:"files"`
	Lines int `json:"lines"`
	Nosec int `json:"nosec"`
	Found int `json:"found"`
}
