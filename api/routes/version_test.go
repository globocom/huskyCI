// Copyright 2019 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routes_test

import (
	apiContext "github.com/globocom/huskyCI/api/context"
	"github.com/globocom/huskyCI/api/routes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getRequestResult", func() {

	expected := map[string]string{
		"version": apiContext.DefaultConf.GetAPIVersion(),
		"date":    apiContext.DefaultConf.GetAPIReleaseDate(),
	}

	apiContext.DefaultConf.SetOnceConfig()
	config := apiContext.APIConfiguration

	Context("When version and date are requested", func() {
		It("Should return a map with API version and date", func() {
			Expect(routes.GetRequestResult(config)).To(Equal(expected))
		})
	})

})
