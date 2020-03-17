// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/globocom/huskyCI/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println("error found: ", err)
		os.Exit(1)
	}
}
