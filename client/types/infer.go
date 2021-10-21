// Copyright 2021 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// InferOutput holds all data from Infer output.
type InferOutput []InferResult

// InferResult holds the detailed information from Infer output.
type InferResult struct {
	Type     string `json:"bug_type"`
	Message  string `json:"qualifier"`
	File     string `json:"file"`
	Line     string `json:"line"`
	Severity string `json:"severity"`
	Title    string `json:"bug_type_hum"`
}
