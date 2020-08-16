package documents_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"voucher-service/documents"
	"voucher-service/model"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

var _ = Describe("Documents", func() {
	var server *ghttp.Server
	var Statuscode int

	BeforeEach(func() {
		server = ghttp.NewServer()
		os.Setenv("BLOCKCHAIN_API_ADDR", server.URL())
		os.Setenv("URL", "http://emtrust.io")

	})

	AfterEach(func() {
		server.Close()
	})

	Describe("ClaimVouchers", func() {

		Context("When api calls voucher is getting claimed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"status": "SUCCESS",
						"info": "ce0e4b69a971bde5756afc9d2556c34feec299fe95b3eb678358d477a783a1ad"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			Statuscode, _ = service.ClaimVouchers([]string{"cfd3298b-de30-42f8-97fd-fe7ae35c122e", "923fa3ed-0c11-48c8-a245-68f273c6df40", "Rs 121 OFF", "fixed value", "{'ExpiryDate':'1582309799999','Description':'get Rs 100 off','Logo':'','Background':'','Transferable':'true','OrganizationDid':'did:emtrust:0xa429540b824f5b63140d'}", "0", "did:emtrust:0x983c7c977bdb55003a5b", "did:emtrust:0x6767b1a114e8041012b6"})
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
		})
	})

	Describe("QueryTokensByOwner", func() {

		Context("When api calls QueryTokensByOwner", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			Statuscode, _ = service.QueryTokensByOwner("did:emtrust:0x6738d1c089e7ad3c981d")
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
		})
		Context("When api calls QueryTokensByOwner passing correct did", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[
						{
							"Key": "c772b2d4-f383-4ba6-97e5-4d26b534fea8",
							"Record": {
								"Template": "200b6ce6-6732-4d66-987c-4d8cf1b9079f",
								"body": "{\"ExpiryDate\":\"1585247399999\",\"Description\":\"save 15$ on next order\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x50890a0b60a6de170e0a\",\"OrganizationName\":\"Tollgate\",\"CreateBy\":\"did:emtrust:0x3fd1ab3d93310d8f0e5d\",\"IssuedAt\":\"1584600463671\",\"Status\":\"Claimed\"}",
								"issuer": "did:emtrust:0x50890a0b60a6de170e0a",
								"name": "Save 15$ off",
								"owner": "did:emtrust:0x6596eb9e8a7a1a264e20",
								"type": "Discount",
								"uuid": "c772b2d4-f383-4ba6-97e5-4d26b534fea8",
								"value": "0"
							}
						}
					]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})
			var response []byte
			service, _ := documents.New(client)
			Statuscode, response = service.QueryTokensByOwner("did:emtrust:0x6596eb9e8a7a1a264e20")
			var data []model.QueryTokensByOwnerResponse
			json.Unmarshal([]byte(response), &data)
			It("should return key similar key", func() {
				Expect(data[0].Key).To(Equal("c772b2d4-f383-4ba6-97e5-4d26b534fea8"))
			})
		})
		Context("When api calls QueryTokensByOwner, passing wrong did", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})
			var response []byte
			service, _ := documents.New(client)
			Statuscode, response = service.QueryTokensByOwner("sdfsd")
			var data []model.QueryTokensByOwnerResponse
			json.Unmarshal([]byte(response), &data)
			It("should return empty array", func() {
				Expect(len(data)).To(Equal(0))
			})
		})
	})

	Describe("GetHistory", func() {

		Context("When api calls GetHistory when I pass wrong voucher id", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})
			var response []byte
			service, _ := documents.New(client)
			Statuscode, response = service.GetHistory("asd")
			var data []model.GetHistoryResponse
			json.Unmarshal([]byte(response), &data)
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
			It("should return Empty array", func() {
				Expect(len(data)).To(Equal(0))
			})
		})
		Context("When api calls GetHistory when I pass Correct voucher id", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[
						{
							"TxId": "ff655031725e3ef216ff02b87e4173edc991f9b847b4d6f31d538d6ff616ec64",
							"Value": {
								"uuid": "c772b2d4-f383-4ba6-97e5-4d26b534fea8",
								"Template": "200b6ce6-6732-4d66-987c-4d8cf1b9079f",
								"name": "Save 15$ off",
								"type": "Discount",
								"body": "{\"ExpiryDate\":\"1585247399999\",\"Description\":\"save 15$ on next order\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x50890a0b60a6de170e0a\",\"OrganizationName\":\"Tollgate\",\"CreateBy\":\"did:emtrust:0x3fd1ab3d93310d8f0e5d\",\"IssuedAt\":\"1584600463671\",\"Status\":\"Claimed\"}",
								"value": "0",
								"owner": "did:emtrust:0x6596eb9e8a7a1a264e20",
								"issuer": "did:emtrust:0x50890a0b60a6de170e0a"
							},
							"Timestamp": "2020-03-19 06:47:43.945 +0000 UTC",
							"IsDelete": "false"
						},
						{
							"TxId": "897b508ae507b994e42c0b41706ba36ff047163cfa0f8cda5c7982ce2f07462b",
							"Value": {
								"uuid": "c772b2d4-f383-4ba6-97e5-4d26b534fea8",
								"Template": "200b6ce6-6732-4d66-987c-4d8cf1b9079f",
								"name": "Save 15$ off",
								"type": "Discount",
								"body": "{\"ExpiryDate\":\"1585247399999\",\"Description\":\"save 15$ on next order\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x50890a0b60a6de170e0a\",\"OrganizationName\":\"Tollgate\",\"CreateBy\":\"did:emtrust:0x3fd1ab3d93310d8f0e5d\",\"IssuedAt\":\"1584602360042\",\"Status\":\"Redeemed\"}",
								"value": "0",
								"owner": "did:emtrust:0x6596eb9e8a7a1a264e20",
								"issuer": "did:emtrust:0x50890a0b60a6de170e0a"
							},
							"Timestamp": "2020-03-19 07:19:20.305 +0000 UTC",
							"IsDelete": "false"
						},
						{
							"TxId": "166974dd0c5bf8e5cdeed7c883f415f0c2b95ef4a4ea51767ef49a296d6f0d68",
							"Value": {
								"uuid": "c772b2d4-f383-4ba6-97e5-4d26b534fea8",
								"Template": "200b6ce6-6732-4d66-987c-4d8cf1b9079f",
								"name": "Save 15$ off",
								"type": "Discount",
								"body": "{\"ExpiryDate\":\"1585247399999\",\"Description\":\"save 15$ on next order\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x50890a0b60a6de170e0a\",\"OrganizationName\":\"Tollgate\",\"CreateBy\":\"did:emtrust:0x3fd1ab3d93310d8f0e5d\",\"IssuedAt\":\"1584602360042\",\"Status\":\"Redeemed\"}",
								"value": "0",
								"owner": "did:emtrust:0x50890a0b60a6de170e0a",
								"issuer": "did:emtrust:0x50890a0b60a6de170e0a"
							},
							"Timestamp": "2020-03-19 07:19:22.744 +0000 UTC",
							"IsDelete": "false"
						}
					]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})
			var response []byte
			service, _ := documents.New(client)
			Statuscode, response = service.GetHistory("c772b2d4-f383-4ba6-97e5-4d26b534fea8")
			var data []model.GetHistoryResponse
			json.Unmarshal([]byte(response), &data)
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
			It("should return Empty array", func() {
				Expect(data[0].Value.Uuid).To(Equal("c772b2d4-f383-4ba6-97e5-4d26b534fea8"))
			})
		})
	})
	Describe("GetVoucherDetails", func() {

		Context("When api calls GetVoucherDetails", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			Statuscode, _ = service.GetVoucherDetails("045cf5eb-9906-407b-8e50-116e3619843e")
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
		})
	})
	Describe("QueryTokensByIssuer", func() {

		Context("When api calls QueryTokensByIssuer", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			Statuscode, _ = service.QueryTokensByIssuer("did:emtrust:0x6738d1c089e7ad3c981d")
			It("should return status code 200", func() {
				Expect(Statuscode).To(Equal(200))
			})
		})
	})
})
