// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

// BrakemanOutput is the struct that holds issues and stats found on a Brakeman scan.
type BrakemanOutput struct {
	Warnings []WarningItem `json:"warnings"`
}

// WarningItem is the struct that holds all detailed information of a vulnerability found.
type WarningItem struct {
	Type       string `json:"warning_type"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	File       string `json:"file"`
	Line       int    `json:"line"`
	Details    string `json:"link"`
	Confidence string `json:"confidence"`
}
