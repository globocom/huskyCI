package github

import (
	"net/http"
	"net/url"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("NewDeviceFlow", func() {
	Context("Normal", func() {
		It("should return a new GitHub device flow", func() {
			want := DeviceFlow{
				baseURI: DefaultBaseURI,
				client:  http.DefaultClient,
			}
			got := NewDeviceFlow(DefaultBaseURI, http.DefaultClient)
			Expect(got).To(Equal(want))
		})
	})
})

var _ = Describe("DeviceFlow", func() {
	var (
		server *ghttp.Server
		df     DeviceFlow
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		baseURI, _ := url.Parse(server.URL())
		df = DeviceFlow{baseURI: baseURI, client: http.DefaultClient}
	})

	AfterEach(func() { server.Close() })

	Describe("GetCodes", func() {
		var (
			req = &GetCodesRequest{
				ClientID: "client_id",
				Scope:    "scope",
			}
			statusCode int
			resp       interface{}
		)

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/login/device/code"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
						"Accept":       []string{"application/json"},
					}),
					ghttp.VerifyJSONRepresenting(req),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resp),
				),
			)
		})

		When("The request succeeds", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				resp = &GetCodesResponse{
					DeviceCode:      "device_code",
					UserCode:        "user_code",
					VerificationURI: "https://www.example.com/verify",
					ExpiresIn:       600,
					Interval:        60,
				}
			})

			It("should return device codes without errors", func() {
				got, err := df.GetCodes(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(resp))
			})
		})

		When("The response is not found", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
				resp = &ErrResponse{Err: "Not Found"}
			})

			It("should return the not found error", func() {
				_, err := df.GetCodes(req)
				Expect(err).To(MatchError("Not Found"))
			})
		})
	})

	Describe("GetAccessToken", func() {
		var (
			req = &GetAccessTokenRequest{
				ClientID:   "client_id",
				DeviceCode: "device_code",
				GrantType:  "grant_type",
			}
			statusCode int
			resp       interface{}
		)

		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/login/oauth/access_token"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
						"Accept":       []string{"application/json"},
					}),
					ghttp.VerifyJSONRepresenting(req),
					ghttp.RespondWithJSONEncodedPtr(&statusCode, &resp),
				),
			)
		})

		When("The request succeeds", func() {
			BeforeEach(func() {
				statusCode = http.StatusOK
				resp = &GetAccessTokenResponse{
					AccessToken: "access_token",
					TokenType:   "bearer",
					Scope:       "scope",
				}
			})

			It("should return access token without errors", func() {
				got, err := df.GetAccessToken(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(got).To(Equal(resp))
			})
		})

		When("The response is not found", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
				resp = &ErrResponse{Err: "Not Found"}
			})

			It("should return the not found error", func() {
				_, err := df.GetAccessToken(req)
				Expect(err).To(MatchError("Not Found"))
			})
		})
	})
})
