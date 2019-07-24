package types

import "errors"

var (
	// ErrorInternal holds the internal error string
	ErrorInternal = errors.New("internal error")
	// ErrorInvalidRID holds the invalid RID error string
	ErrorInvalidRID = errors.New("invalid rid")
	// ErrorInvalidRepository holds the invalid repository string
	ErrorInvalidRepository = errors.New("invalid repository")
	// ErrorInvalidBranch holds the invalid branch string
	ErrorInvalidBranch = errors.New("invalid branch")
	// ErrorInvalidDependencyURL holds the invalid Dependency URL string
	ErrorInvalidDependencyURL = errors.New("invalid dep url")
)
