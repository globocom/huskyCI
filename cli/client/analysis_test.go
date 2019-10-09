package client

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/globocom/huskyCI/client/types"
	"github.com/stretchr/testify/assert"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

func TestClientStart(t *testing.T) {
	t.Run(
		"Test Start() with error 5XX",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte(""))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			_, err := hcli.Start("repo_url", "repo_branch")
			if err == nil {
				t.Fatalf("CLIENT START: fail to validate Husky API response (%v)", err)
			}
		})
	t.Run(
		"Test Start() without RID",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte(""))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			_, err := hcli.Start("repo_url", "repo_branch")
			if err == nil {
				t.Fatalf("CLIENT START: fail to validate Husky API RID (%v)", err)
			}
		})
	t.Run(
		"Test Start() without token",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Husky-Token") == "Husky-Token-Value" {
					w.Header().Set("X-Request-Id", "token")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(""))
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(""))
				}
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			_, err := hcli.Start("repo_url", "repo_branch")
			if err == nil {
				t.Fatalf("CLIENT START: fail to validate call without Husky-Token (%v)", err)
			}
		})
	t.Run(
		"Test Start() with token",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("Husky-Token") == "Husky-Token-Value" {
					w.Header().Set("X-Request-Id", "token")
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(""))
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(""))
				}
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com", Token: "Husky-Token-Value"})
			hcli.httpCli = httpClient

			result, err := hcli.Start("repo_url", "repo_branch")
			if err != nil {
				t.Fatalf("CLIENT START: fail to validate call with Husky-Token (%v)", err)
			}

			assert.Equal(t, types.Analysis{
				URL:    "repo_url",
				Branch: "repo_branch",
				RID:    "token",
			}, result)
		})
}

func TestClientGet(t *testing.T) {
	t.Run(
		"Test Get() with error 5XX",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusGatewayTimeout)
				w.Write([]byte(""))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			_, err := hcli.Get("request-id-example-test")
			if err == nil {
				t.Fatalf("CLIENT GET: fail to validate Husky API status code (%v)", err)
			}
		})
	t.Run(
		"Test Get() with 200 OK, but malformated JSON",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}}"))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			_, err := hcli.Get("request-id-example-test")
			if err == nil {
				t.Fatalf("CLIENT GET: fail to validate Husky API response (%v)", err)
			}
		})
	t.Run(
		"Test Get() with 200 OK and wellformated JSON",
		func(t *testing.T) {
			var validOutput = `{
				"RID": "request-id-example-test",
				"repositoryURL": "git@github.com:example/example-rails.git",
				"repositoryBranch": "vuln-branch",
				"commitAuthors": [
				  "rafaveira3@gmail.com"
				],
				"status": "finished",
				"result": "failed",
				"errorFound": "",
				"containers": [
				],
				"startedAt": "2019-10-01T11:58:49.159-03:00",
				"finishedAt": "2019-10-01T11:59:40.064-03:00",
				"codes": [
				  {
					"language": "JavaScript",
					"files": [
					  "app/assets/config/manifest.js",
					  "app/assets/javascripts/application.js",
					  "app/assets/javascripts/cable.js",
					  "app/assets/javascripts/realize.js",
					  "public/javascripts/i18n.js",
					  "public/javascripts/translations.js"
					]
				  },
				  {
					"language": "Ruby",
					"files": [
					  "app/controllers/users_controller.rb"
					]
				  }
				],
				"huskyciresults": {
				  "rubyresults": {
					"brakemanoutput": {
					  "mediumvulns": [
						{
						  "language": "Ruby",
						  "securitytool": "Brakeman",
						  "confidence": "Medium",
						  "file": "app/code/app/controllers/users_controller.rb",
						  "line": "48",
						  "code": "params.require(:user).permit(:name, :email, :password, :password_confirmation, :role)",
						  "details": "https://brakemanscanner.org/docs/warning_types/mass_assignment/Potentially dangerous key allowed for mass assignment",
						  "type": "Mass Assignment"
						}
					  ]
					}
				  }
				}
			  }`
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(validOutput))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			result, err := hcli.Get("request-id-example-test")
			if err != nil {
				t.Fatalf("CLIENT GET: fail to validate Husky API response (%v)", err)
			}

			assert.IsType(t, types.Analysis{}, result)
		})
}

func TestClientMonitor(t *testing.T) {
	t.Run(
		"Test Monitor() with timeout",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Second)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}"))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			timeoutMonitor, _ := time.ParseDuration("1s")
			retryMonitor, _ := time.ParseDuration("100ms")
			_, err := hcli.Monitor("request-id-example-test", timeoutMonitor, retryMonitor)
			if err == nil {
				t.Fatalf("CLIENT MONITOR: fail to timeout Husky API (%v)", err)
			}
		})
	t.Run(
		"Test Monitor() with Get() fail",
		func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("{}}"))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			timeoutMonitor, _ := time.ParseDuration("1s")
			retryMonitor, _ := time.ParseDuration("100ms")
			_, err := hcli.Monitor("request-id-example-test", timeoutMonitor, retryMonitor)
			if err == nil {
				t.Fatalf("CLIENT MONITOR: fail to timeout Husky API (%v)", err)
			}
		})
	t.Run(
		"Test Monitor() with Get() finished",
		func(t *testing.T) {
			var validOutput = `{
					"RID": "request-id-example-test",
					"repositoryURL": "git@github.com:example/example-rails.git",
					"repositoryBranch": "vuln-branch",
					"commitAuthors": [
					  "rafaveira3@gmail.com"
					],
					"status": "finished",
					"result": "failed",
					"errorFound": "",
					"containers": [
					],
					"startedAt": "2019-10-01T11:58:49.159-03:00",
					"finishedAt": "2019-10-01T11:59:40.064-03:00",
					"codes": [
					  {
						"language": "JavaScript",
						"files": [
						  "app/assets/config/manifest.js",
						  "app/assets/javascripts/application.js",
						  "app/assets/javascripts/cable.js",
						  "app/assets/javascripts/realize.js",
						  "public/javascripts/i18n.js",
						  "public/javascripts/translations.js"
						]
					  },
					  {
						"language": "Ruby",
						"files": [
						  "app/controllers/users_controller.rb"
						]
					  }
					],
					"huskyciresults": {
					  "rubyresults": {
						"brakemanoutput": {
						  "mediumvulns": [
							{
							  "language": "Ruby",
							  "securitytool": "Brakeman",
							  "confidence": "Medium",
							  "file": "app/code/app/controllers/users_controller.rb",
							  "line": "48",
							  "code": "params.require(:user).permit(:name, :email, :password, :password_confirmation, :role)",
							  "details": "https://brakemanscanner.org/docs/warning_types/mass_assignment/Potentially dangerous key allowed for mass assignment",
							  "type": "Mass Assignment"
							}
						  ]
						}
					  }
					}
				  }`
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(validOutput))
			})
			httpClient, teardown := testingHTTPClient(h)
			defer teardown()

			hcli := NewClient(types.Target{Endpoint: "http://example.com"})
			hcli.httpCli = httpClient

			timeoutMonitor, _ := time.ParseDuration("1s")
			retryMonitor, _ := time.ParseDuration("100ms")
			result, err := hcli.Monitor("request-id-example-test", timeoutMonitor, retryMonitor)
			if err != nil {
				t.Fatalf("CLIENT MONITOR: fail to timeout Husky API (%v)", err)
			}
			assert.IsType(t, types.Analysis{}, result)
		})
}
