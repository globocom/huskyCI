// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/globocom/huskyCI/cli/cmd"
	"github.com/globocom/huskyCI/cli/errorcli"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		errorcli.Handle(err)
	}
}
