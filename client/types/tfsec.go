// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "encoding/json"

// TFSecOutput is the struct that holds all data from TFSec output.
type TFSecOutput struct {
	Warnings json.RawMessage `json:"warnings"`
	Results  []TFSecResult   `json:"results"`
}

// TFSecResult is the struct that holds detailed information of results from TFSec output.
type TFSecResult struct {
	RuleID      string   `json:"rule_id"`
	Link        string   `json:"link"`
	Location    Location `json:"location"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
}

// Location is the struct that holds detailed information of location from each result
type Location struct {
	Filename  string `json:"filename"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}
