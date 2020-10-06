// Copyright 2020 Globo.com authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/globocom/huskyCI/cli/pkg/github"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with GitHub",
	Long: `
Log in to the GitHub.
`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{Timeout: time.Minute}
		df := github.NewDeviceFlow(github.DefaultBaseURI, client)
		getCodesResp, err := df.GetCodes(&github.GetCodesRequest{
			ClientID: github.ClientID,
		})
		if err != nil {
			return fmt.Errorf("error getting device codes: %w", err)
		}

		if err := browser.OpenURL(getCodesResp.VerificationURI); err != nil {
			return fmt.Errorf("error opening verification page: %w", err)
		}

		fmt.Printf("Please enter the user code in a browser: %s\n", getCodesResp.UserCode)
		fmt.Print("Then press Enter...")
		if _, err := bufio.NewReader(os.Stdin).ReadBytes('\n'); err != nil {
			panic(err)
		}

		resp, err := df.GetAccessToken(&github.GetAccessTokenRequest{
			ClientID:   github.ClientID,
			DeviceCode: getCodesResp.DeviceCode,
			GrantType:  github.GrantTypeDeviceCode,
		})
		if err != nil {
			return fmt.Errorf("error getting access token: %w", err)
		}

		if err := ioutil.WriteFile(".huskyci", []byte(resp.AccessToken), 0600); err != nil {
			return fmt.Errorf("error saving access token: %w", err)
		}

		fmt.Println("Login successful!🚀")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
