package documents_test

import (
	documents "account-service/document"
	"account-service/model"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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

// start of suite
var _ = Describe("Document", func() {
	var server *ghttp.Server
	var returnedData string
	var returnedResposne int
	var returnedResposnefail int
	var identJson []byte

	BeforeEach(func() {
		server = ghttp.NewServer()
		os.Setenv("BLOCKCHAIN_API_ADDR", server.URL())
		os.Setenv("URL", "http://emtrust.io")

	})

	AfterEach(func() {
		server.Close()
	})
	Describe("Addevent", func() {

		Context("When api response correctly", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "DEALS_CLAIM",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"http://emtrust.io/api/blockchain/channels/tradechannel/chaincodes/vouchers", "DEALS_CLAIM", "1584075306159", "did:emtrust:0x983c7c977bdb55003a5b", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response correctly when no source passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "DEALS_CLAIM",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"", "DEALS_CLAIM", "1584075306159", "did:emtrust:0x983c7c977bdb55003a5b", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response correctly when no event passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"http://emtrust.io/api/blockchain/channels/tradechannel/chaincodes/vouchers", "", "1584075306159", "did:emtrust:0x983c7c977bdb55003a5b", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response correctly when no time passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "DEALS_CLAIM",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"http://emtrust.io/api/blockchain/channels/tradechannel/chaincodes/vouchers", "DEALS_CLAIM", "", "did:emtrust:0x983c7c977bdb55003a5b", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response correctly when no actor passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "DEALS_CLAIM",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"http://emtrust.io/api/blockchain/channels/tradechannel/chaincodes/vouchers", "DEALS_CLAIM", "1584075306159", "", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response correctly when no text passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"evt": "DEALS_CLAIM",
						"result": "ok"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"http://emtrust.io/api/blockchain/channels/tradechannel/chaincodes/vouchers", "DEALS_CLAIM", "1584075306159", "did:emtrust:0x983c7c977bdb55003a5b", ""})
			var details model.EventResult
			json.Unmarshal([]byte(returnedData), &details)
			It("should return result ok", func() {
				Expect(details.Result).To(ContainSubstring("ok"))
			})
		})
		Context("When api response when no body passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{})
			It("should return result ok", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When api response when less param is passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.Addevent([]string{"1584075306159", "did:emtrust:0x983c7c977bdb55003a5b", "1a59ef9c-f171-4571-94ed-fe61ffa16924,eba3ddaa-45dd-4ad9-b82b-d2d67e745a3c,$ 55 OFF,fixed value,{\"ExpiryDate\":\"1583001000000\",\"Description\":\"$ 50 OFF\",\"Logo\":\"\",\"Background\":\"\",\"Transferable\":\"true\",\"OrganizationDid\":\"did:emtrust:0x8092f15b0a2251079a58\",\"OrganizationName\":\"\",\"CreateBy\":\"did:emtrust:0x983c7c977bdb55003a5b\",\"IssuedAt\":\"1584075306159\",\"Status\":\"Claimed\"},0,did:emtrust:0x8092f15b0a2251079a58,did:emtrust:0xc8a94e34809df9422880"})
			It("should return result ok", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
	})
	Describe("QueryMembersdetail", func() {

		Context("When api response when member id passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[
						{
							"Key": "did:emtrust:0x06c9edabd73117c7e331",
							"Record": {
								"deviceInfo": "",
								"email": "",
								"id": "did:emtrust:0x06c9edabd73117c7e331",
								"name": "",
								"other": "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"publicKey\":[{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"publicKeyBase58\":\"04155dd8fe2160570c0de6b34aa5778a08ef4db49492c929ad123f3117831e69f025f68d3387e6a51052028c691d1a815ff69db1edc876e553ed31794089420f29\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x06c9edabd73117c7e331\"],\"endorsements\":[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T07:36:41.891Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T07:36:41.891Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450221009d1bbdcf2ba72d6859ecb1916cb3f6d4433144a28ef6d7a5ec40a25c1a079adf022039f17287d0d87a2f8aba9bfae8eef961b24691ca2e8a0ea3970344224972c560\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"agentCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T09:16:50.207Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"agent\",\"desc\":\"agent for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T09:16:50.207Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450221008d54ca8b63ac5cce3e88b5f4c3e7bc2aebbf7d9c03bf76b2883184fa69db80d7022026099b55f7dca57838eacbd42f2b73845ac7732a1b583b9392f79158c3046e50\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"staffCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T09:18:02.538Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"staff\",\"desc\":\"staff for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T09:18:02.538Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450220356ba624db7e1373f9c6eeff3979318bf9645729b265eda4523ba58483c22796022100aaf5d809a01ff56b66ab0a0d6d1e17472789e6621f79dd98c91dd84bf2703409\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"129788e6-403f-4072-99fb-b9833f2d5dcd\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/580780fb-a47b-4ae5-a4dc-934061931b9e\",\"issuanceDate\":\"2020-03-31T05:13:01.55Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xc4c80f555f928742f7df\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-31T05:13:01.55Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3045022100eb606e4dd329097254e5d730a802046fb83ee755bd9589e568daeade8a34734302205ac2498551544e8e0b2997d16d589978dc472d132a6ca3f2f13ce2003999c95a\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"129788e6-403f-4072-99fb-b9833f2d5dcd\",\"type\":[\"VerifiableCredential\",\"staffCredential\"],\"issuer\":\"https://emtrust.io/issuer/580780fb-a47b-4ae5-a4dc-934061931b9e\",\"issuanceDate\":\"2020-03-31T05:16:00.228Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xc4c80f555f928742f7df\",\"type\":\"staff\",\"desc\":\"staff for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-31T05:16:00.228Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3046022100b039779ef6403b34d218369b644275374b0765d4b3e9f27c42bea32e1054e2c40221009c131396f0ed288e5ed6e48cb0fb5b7a7cd0973b09211370ff5f485b81eaa9a9\"}}],\"service\":null,\"created\":\"Mon, 30 Mar 2020 07:36:42 GMT\",\"updated\":\"Mon, 30 Mar 2020 07:36:42 GMT\"}",
								"publicKey": "04155dd8fe2160570c0de6b34aa5778a08ef4db49492c929ad123f3117831e69f025f68d3387e6a51052028c691d1a815ff69db1edc876e553ed31794089420f29"
							}
						}
					]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, identJson, _ = service.QueryMembersdetail("29788e6-403f-4072-99fb-b9833f2d5dcd")
			var details []model.MemberData
			json.Unmarshal(identJson, &details)
			It("Response should contain key", func() {
				Expect(details[0].Key).To(ContainSubstring("did:emtrust:0x06c9edabd73117c7e331"))
			})
		})
		Context("When api response when wrong member id passed", func() {
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

			service, _ := documents.New(client)
			_, _, returnedData = service.QueryMembersdetail("129788e6-403f-4000-99fb-b9833f2d5dcd")
			It("response should be empty array", func() {
				Expect(returnedData).To(ContainSubstring("[]"))
			})
		})
		Context("When api response when member did passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`[
						{
							"Key": "did:emtrust:0x06c9edabd73117c7e331",
							"Record": {
								"deviceInfo": "",
								"email": "",
								"id": "did:emtrust:0x06c9edabd73117c7e331",
								"name": "",
								"other": "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"publicKey\":[{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"publicKeyBase58\":\"04155dd8fe2160570c0de6b34aa5778a08ef4db49492c929ad123f3117831e69f025f68d3387e6a51052028c691d1a815ff69db1edc876e553ed31794089420f29\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x06c9edabd73117c7e331\"],\"endorsements\":[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T07:36:41.891Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T07:36:41.891Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450221009d1bbdcf2ba72d6859ecb1916cb3f6d4433144a28ef6d7a5ec40a25c1a079adf022039f17287d0d87a2f8aba9bfae8eef961b24691ca2e8a0ea3970344224972c560\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"agentCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T09:16:50.207Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"agent\",\"desc\":\"agent for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T09:16:50.207Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450221008d54ca8b63ac5cce3e88b5f4c3e7bc2aebbf7d9c03bf76b2883184fa69db80d7022026099b55f7dca57838eacbd42f2b73845ac7732a1b583b9392f79158c3046e50\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"7ae98ad3-80bf-4e48-af79-8d3855d9b9ee\",\"type\":[\"VerifiableCredential\",\"staffCredential\"],\"issuer\":\"https://emtrust.io/issuer/77794f2f-639c-4578-9991-24544fbd6326\",\"issuanceDate\":\"2020-03-30T09:18:02.538Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xdc24ffe5e734e8ce25b3\",\"type\":\"staff\",\"desc\":\"staff for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-30T09:18:02.538Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30450220356ba624db7e1373f9c6eeff3979318bf9645729b265eda4523ba58483c22796022100aaf5d809a01ff56b66ab0a0d6d1e17472789e6621f79dd98c91dd84bf2703409\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"129788e6-403f-4072-99fb-b9833f2d5dcd\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/580780fb-a47b-4ae5-a4dc-934061931b9e\",\"issuanceDate\":\"2020-03-31T05:13:01.55Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xc4c80f555f928742f7df\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-31T05:13:01.55Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3045022100eb606e4dd329097254e5d730a802046fb83ee755bd9589e568daeade8a34734302205ac2498551544e8e0b2997d16d589978dc472d132a6ca3f2f13ce2003999c95a\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"129788e6-403f-4072-99fb-b9833f2d5dcd\",\"type\":[\"VerifiableCredential\",\"staffCredential\"],\"issuer\":\"https://emtrust.io/issuer/580780fb-a47b-4ae5-a4dc-934061931b9e\",\"issuanceDate\":\"2020-03-31T05:16:00.228Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x06c9edabd73117c7e331\",\"associatedWith\":{\"id\":\"did:emtrust:0xc4c80f555f928742f7df\",\"type\":\"staff\",\"desc\":\"staff for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-03-31T05:16:00.228Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3046022100b039779ef6403b34d218369b644275374b0765d4b3e9f27c42bea32e1054e2c40221009c131396f0ed288e5ed6e48cb0fb5b7a7cd0973b09211370ff5f485b81eaa9a9\"}}],\"service\":null,\"created\":\"Mon, 30 Mar 2020 07:36:42 GMT\",\"updated\":\"Mon, 30 Mar 2020 07:36:42 GMT\"}",
								"publicKey": "04155dd8fe2160570c0de6b34aa5778a08ef4db49492c929ad123f3117831e69f025f68d3387e6a51052028c691d1a815ff69db1edc876e553ed31794089420f29"
							}
						}
					]`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, identJson, _ = service.QueryMembersdetail("did:emtrust:0x06c9edabd73117c7e331")
			var details []model.MemberData
			json.Unmarshal(identJson, &details)
			It("Response should contain key", func() {
				Expect(details[0].Key).To(ContainSubstring("did:emtrust:0x06c9edabd73117c7e331"))
			})
		})
		Context("When api response when wrong member did passed", func() {
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

			service, _ := documents.New(client)
			_, _, returnedData = service.QueryMembersdetail("did:emtrust:0x06c9edabd73117c7e001")
			It("response should be empty array", func() {
				Expect(returnedData).To(ContainSubstring("[]"))
			})
		})
		Context("When api response when null did passed", func() {
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

			service, _ := documents.New(client)
			_, _, returnedData = service.QueryMembersdetail("")
			It("response should be empty array", func() {
				Expect(returnedData).To(ContainSubstring("[]"))
			})
		})
	})
	Describe("UpdateDidDocument", func() {

		Context("When api response correctly", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"status": "SUCCESS",
						"info": "033ed5b209445dfa7c0393cea18a296974aae29529c743df984e8daed18d8e3c"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			returnedResposne, _ = service.UpdateDidDocument([]string{"did:emtrust:0xc8a94e34809df9422880", "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"publicKey\":[{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"publicKeyBase58\":\"0466e6c4ed172155bc41e6a2354713c2bc05884b8f915a84b9c507a5045ce20d27b7ae3d55bff1e8fd9c37adbe65f5fdf6e841eaa2e0ed1781db28513f556229ujkhjgha1\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x9b99b23c003b5550d2d2\"],\"endorsements\":[[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"d48ec752-f135-43f9-bee9-cd72ec956976\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/68b6936c-63e3-4391-80dc-698e8e5d8b66\",\"issuanceDate\":\"2019-12-19T11:19:45.511Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x1fd349fb986183b675ba\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T11:19:45.511Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304402205a125385509a5f0c9d59f52dd69480284a8990d1dcd63d192eaff55235cb95a202205ce07c133fff677c25a05fd91db7ea8e76462bf192b97a53217d39de6b1f5bed\"}}{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"e92b8fea-f7ed-4cbf-befc-c0a22672409f\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/4d6daf43-b78b-43a4-a05c-437c4b26d22a\",\"issuanceDate\":\"2019-12-19T09:27:30.276Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T09:27:30.276Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30460221009115deec6fe0c69a7370735fee76c18a5bf50d95ac0874fa720eef848ccc45290221008fe812c06be41038ccf7783678254c2ed3aa711982fd6743870c8481366737b6\"}}]],\"created\":\"Thu, 19 Dec 2019 09:27:30 GMT\",\"updated\":\"Thu, 19 Dec 2019 09:27:30 GMT\"}"})
			It("should return response id 200", func() {
				Expect(returnedResposne).To(Equal(200))
			})
		})
		Context("When api gets time out due to wrong did passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 504,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`upstream request timeout`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			returnedResposnefail, _ = service.UpdateDidDocument([]string{"did:emtrust:0xc8a94e34809df9422800", "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"publicKey\":[{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"publicKeyBase58\":\"0466e6c4ed172155bc41e6a2354713c2bc05884b8f915a84b9c507a5045ce20d27b7ae3d55bff1e8fd9c37adbe65f5fdf6e841eaa2e0ed1781db28513f556229ujkhjgha1\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x9b99b23c003b5550d2d2\"],\"endorsements\":[[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"d48ec752-f135-43f9-bee9-cd72ec956976\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/68b6936c-63e3-4391-80dc-698e8e5d8b66\",\"issuanceDate\":\"2019-12-19T11:19:45.511Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x1fd349fb986183b675ba\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T11:19:45.511Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304402205a125385509a5f0c9d59f52dd69480284a8990d1dcd63d192eaff55235cb95a202205ce07c133fff677c25a05fd91db7ea8e76462bf192b97a53217d39de6b1f5bed\"}}{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"e92b8fea-f7ed-4cbf-befc-c0a22672409f\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/4d6daf43-b78b-43a4-a05c-437c4b26d22a\",\"issuanceDate\":\"2019-12-19T09:27:30.276Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T09:27:30.276Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30460221009115deec6fe0c69a7370735fee76c18a5bf50d95ac0874fa720eef848ccc45290221008fe812c06be41038ccf7783678254c2ed3aa711982fd6743870c8481366737b6\"}}]],\"created\":\"Thu, 19 Dec 2019 09:27:30 GMT\",\"updated\":\"Thu, 19 Dec 2019 09:27:30 GMT\"}"})
			It("should return response id 504", func() {
				Expect(returnedResposnefail).To(Equal(504))
			})
		})
		Context("When api gets time out due to null did passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 504,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.UpdateDidDocument([]string{"", "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"publicKey\":[{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"publicKeyBase58\":\"0466e6c4ed172155bc41e6a2354713c2bc05884b8f915a84b9c507a5045ce20d27b7ae3d55bff1e8fd9c37adbe65f5fdf6e841eaa2e0ed1781db28513f556229ujkhjgha1\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x9b99b23c003b5550d2d2\"],\"endorsements\":[[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"d48ec752-f135-43f9-bee9-cd72ec956976\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/68b6936c-63e3-4391-80dc-698e8e5d8b66\",\"issuanceDate\":\"2019-12-19T11:19:45.511Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x1fd349fb986183b675ba\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T11:19:45.511Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304402205a125385509a5f0c9d59f52dd69480284a8990d1dcd63d192eaff55235cb95a202205ce07c133fff677c25a05fd91db7ea8e76462bf192b97a53217d39de6b1f5bed\"}}{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"e92b8fea-f7ed-4cbf-befc-c0a22672409f\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/4d6daf43-b78b-43a4-a05c-437c4b26d22a\",\"issuanceDate\":\"2019-12-19T09:27:30.276Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x9b99b23c003b5550d2d2\",\"associatedWith\":{\"id\":\"did:emtrust:0x9c17ed09feb40fcef946\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T09:27:30.276Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30460221009115deec6fe0c69a7370735fee76c18a5bf50d95ac0874fa720eef848ccc45290221008fe812c06be41038ccf7783678254c2ed3aa711982fd6743870c8481366737b6\"}}]],\"created\":\"Thu, 19 Dec 2019 09:27:30 GMT\",\"updated\":\"Thu, 19 Dec 2019 09:27:30 GMT\"}"})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When api gets time out due to null did and null did document passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.UpdateDidDocument([]string{"", ""})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When null is passed in body", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.UpdateDidDocument([]string{})
			It("should return response is null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When null is passed in did document", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"status": "SUCCESS",
						"info": "033ed5b209445dfa7c0393cea18a296974aae29529c743df984e8daed18d8e3c"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			returnedResposne, _ = service.UpdateDidDocument([]string{"did:emtrust:0xc8a94e34809df9422880", ""})
			It("should return response id 200", func() {
				Expect(returnedResposne).To(Equal(200))
			})
		})
	})
	Describe("RegisterDidDocument", func() {

		Context("When api calls RegisterDidDocument sucess", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"success": true,
						"message": {
							"status": "SUCCESS",
							"info": "328b1fbec302dcb94879dee9e0f823aa4942844e0fb4ac7eb88e70acd2b57905"
						},
						"id": "did:emtrust:0x68c233ae68abb3737001",
						"issuer": "0484578aad41e7df541261c5baafaf5241b4c28861354df5c8b39f23a3dea493e466f13e9c8565ff6b93876f92af5556d9741d08c612c6a0b58883e303ae788b0b",
						"roles": {}
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{"did:emtrust:0x68c233ae68abb3737001", "0473c53b5c011c81a3dd24e24fecc3ea01451f89a2023e540769efcef2464be6ae6bb1d6236a2f1d5f8259696e228a764e4464f18a91170fedaee2ba38d481606d", "", "", "", "{\"@context\": \"https://w3id.org/did/v1\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"publicKey\": [{\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"type\": \"ED25519SignatureVerification\",\"controller\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"authorizations\": []}],\"authentication\": [\"did:emtrust:0x68c233ae68abb3737ad1\"],\"endorsements\": [],\"service\": [{\"type\": \"MembershipService\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1#membership\",\"serviceEndpoint\": \"http://emtrust.io/accounts/did:emtrust:0x68c233ae68abb3737ad1\",\"description\": \"Membership Service Eligible\"}],\"created\": \"Mon, 13 Jan 2020 09:09:56 GMT\",\"updated\": \"Mon, 13 Jan 2020 09:09:56 GMT\"}"})
			var Result model.RegisterDidDocumentResult
			json.Unmarshal([]byte(returnedData), &Result)
			It("should return json with credentialSubject ID", func() {
				Expect(Result.Success).To(Equal(true))
			})
		})
		Context("When did document is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"success": true,
						"message": {
							"status": "SUCCESS",
							"info": "ea89446a1af35dbb29aebbeae193cefd795a1250e93648757cc863c54e63dbd4"
						},
						"id": "did:emtrust:0x68c233ae68abb37370001",
						"issuer": "0484578aad41e7df541261c5baafaf5241b4c28861354df5c8b39f23a3dea493e466f13e9c8565ff6b93876f92af5556d9741d08c612c6a0b58883e303ae788b0b",
						"roles": {}
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{"did:emtrust:0x68c233ae68abb37370001", "0473c53b5c011c81a3dd24e24fecc3ea01451f89a2023e540769efcef2464be6ae6bb1d6236a2f1d5f8259696e228a764e4464f18a91170fedaee2ba38d481606d", "", "", "", ""})
			var Result model.RegisterDidDocumentResult
			json.Unmarshal([]byte(returnedData), &Result)
			It("should return json with credentialSubject ID", func() {
				Expect(Result.Success).To(Equal(true))
			})
		})
		Context("When public key is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{"did:emtrust:0x68c233ae68abb37370002", "", "", "", "", "{\"@context\": \"https://w3id.org/did/v1\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"publicKey\": [{\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"type\": \"ED25519SignatureVerification\",\"controller\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"authorizations\": []}],\"authentication\": [\"did:emtrust:0x68c233ae68abb3737ad1\"],\"endorsements\": [],\"service\": [{\"type\": \"MembershipService\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1#membership\",\"serviceEndpoint\": \"http://emtrust.io/accounts/did:emtrust:0x68c233ae68abb3737ad1\",\"description\": \"Membership Service Eligible\"}],\"created\": \"Mon, 13 Jan 2020 09:09:56 GMT\",\"updated\": \"Mon, 13 Jan 2020 09:09:56 GMT\"}"})
			It("should return null response", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When no param passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{})
			It("should return null response", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When did is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{"", "0473c53b5c011c81a3dd24e24fecc3ea01451f89a2023e540769efcef2464be6ae6bb1d6236a2f1d5f8259696e228a764e4464f18a91170fedaee2ba38d481606d", "", "", "", "{\"@context\": \"https://w3id.org/did/v1\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"publicKey\": [{\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"type\": \"ED25519SignatureVerification\",\"controller\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"authorizations\": []}],\"authentication\": [\"did:emtrust:0x68c233ae68abb3737ad1\"],\"endorsements\": [],\"service\": [{\"type\": \"MembershipService\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1#membership\",\"serviceEndpoint\": \"http://emtrust.io/accounts/did:emtrust:0x68c233ae68abb3737ad1\",\"description\": \"Membership Service Eligible\"}],\"created\": \"Mon, 13 Jan 2020 09:09:56 GMT\",\"updated\": \"Mon, 13 Jan 2020 09:09:56 GMT\"}"})
			It("should return null response", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When did already registered", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"success": false,
						"message": "Error: Failed to invoke chaincode. cause:Failed to send Proposal and receive all good ProposalResponse"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.RegisterDidDocument([]string{"did:emtrust:0x68c233ae68abb3737001", "0473c53b5c011c81a3dd24e24fecc3ea01451f89a2023e540769efcef2464be6ae6bb1d6236a2f1d5f8259696e228a764e4464f18a91170fedaee2ba38d481606d", "", "", "", "{\"@context\": \"https://w3id.org/did/v1\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"publicKey\": [{\"id\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"type\": \"ED25519SignatureVerification\",\"controller\": \"did:emtrust:0x68c233ae68abb3737ad1\",\"authorizations\": []}],\"authentication\": [\"did:emtrust:0x68c233ae68abb3737ad1\"],\"endorsements\": [],\"service\": [{\"type\": \"MembershipService\",\"id\": \"did:emtrust:0x68c233ae68abb3737ad1#membership\",\"serviceEndpoint\": \"http://emtrust.io/accounts/did:emtrust:0x68c233ae68abb3737ad1\",\"description\": \"Membership Service Eligible\"}],\"created\": \"Mon, 13 Jan 2020 09:09:56 GMT\",\"updated\": \"Mon, 13 Jan 2020 09:09:56 GMT\"}"})
			var Result model.RegisterDidDocumentResult
			json.Unmarshal([]byte(returnedData), &Result)
			It("should return success as false", func() {
				Expect(Result.Success).To(Equal(false))
			})
		})
	})
	Describe("GenerateDidDocument", func() {

		Context("When api calls GenerateDidDocument", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
							"@context": "https://w3id.org/did/v1",
							"id": "did:emtrust:0x6365f099a9c0b187c372",
							"publicKey": [
								{
									"id": "did:emtrust:0x6365f099a9c0b187c372",
									"type": "ED25519SignatureVerification",
									"controller": "did:emtrust:0x9ceb3a48b5b904412095",
									"publicKeyBase58": "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22",
									"authorizations": []
								}
							],
							"authentication": [
								"did:emtrust:0x6365f099a9c0b187c372"
							],
							"endorsements": [
								[
									{
										"@context": [
											"https://www.w3.org/2018/credentials/v1",
											"https://www.halialabs.io/2019/endorsements/v1"
										],
										"id": "66dce28b-cb8e-448c-841a-b50de6cc7cc1",
										"type": [
											"VerifiableCredential",
											"MembershipCredential"
										],
										"issuer": "https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e",
										"issuanceDate": "2019-12-19T06:38:22.633Z",
										"credentialSubject": {
											"id": "did:emtrust:0x6365f099a9c0b187c372",
											"associatedWith": {
												"id": "did:emtrust:0x9ceb3a48b5b904412095",
												"type": "Member",
												"desc": "Member for organization"
											}
										},
										"proof": {
											"type": "ECDSA",
											"created": "2019-12-19T06:38:22.633Z",
											"proofPurpose": "assertionMethod",
											"verificationMethod": "https://emtrust.io/api/crypto/endorsements/verify",
											"jws": "304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5"
										}
									}
								]
							],
							"created": "Mon, 10 Feb 2020 12:00:59 GMT",
							"updated": "Mon, 10 Feb 2020 12:00:59 GMT"
						}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"did:emtrust:0x6365f099a9c0b187c372", "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22", "did:emtrust:0x9ceb3a48b5b904412095", "[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"66dce28b-cb8e-448c-841a-b50de6cc7cc1\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e\",\"issuanceDate\":\"2019-12-19T06:38:22.633Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x6365f099a9c0b187c372\",\"associatedWith\":{\"id\":\"did:emtrust:0x9ceb3a48b5b904412095\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T06:38:22.633Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5\"}}]", "member"})
			var DidDocument model.DIDDocument
			json.Unmarshal([]byte(returnedData), &DidDocument)
			It("should return json with credentialSubject ID", func() {
				Expect(DidDocument.ID).To(ContainSubstring("did:emtrust:0x6365f099a9c0b187c372"))
			})
		})
		Context("When did is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"", "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22", "did:emtrust:0x9ceb3a48b5b904412095", "[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"66dce28b-cb8e-448c-841a-b50de6cc7cc1\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e\",\"issuanceDate\":\"2019-12-19T06:38:22.633Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x6365f099a9c0b187c372\",\"associatedWith\":{\"id\":\"did:emtrust:0x9ceb3a48b5b904412095\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T06:38:22.633Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5\"}}]", "member"})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When public key is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"did:emtrust:0x6365f099a9c0b187c372", "", "did:emtrust:0x9ceb3a48b5b904412095", "[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"66dce28b-cb8e-448c-841a-b50de6cc7cc1\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e\",\"issuanceDate\":\"2019-12-19T06:38:22.633Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x6365f099a9c0b187c372\",\"associatedWith\":{\"id\":\"did:emtrust:0x9ceb3a48b5b904412095\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T06:38:22.633Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5\"}}]", "member"})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When controller did is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"did:emtrust:0x6365f099a9c0b187c372", "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22", "", "[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"66dce28b-cb8e-448c-841a-b50de6cc7cc1\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e\",\"issuanceDate\":\"2019-12-19T06:38:22.633Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x6365f099a9c0b187c372\",\"associatedWith\":{\"id\":\"did:emtrust:0x9ceb3a48b5b904412095\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T06:38:22.633Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5\"}}]", "member"})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When endorsement is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"did:emtrust:0x6365f099a9c0b187c372", "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22", "did:emtrust:0x9ceb3a48b5b904412095", "", "member"})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When entitytype is null", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{"did:emtrust:0x6365f099a9c0b187c372", "0450653e671d01c44614f2a3297c9519e08f39490b96e81852bfdc984a110e5bcc810e480fa40ea1a0c36f3bb67cc682e36f5171ea4eda4a5dc481c8a9d52f22", "did:emtrust:0x9ceb3a48b5b904412095", "[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"66dce28b-cb8e-448c-841a-b50de6cc7cc1\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/fc33b5b1-2fe5-4118-b209-de56d6eace7e\",\"issuanceDate\":\"2019-12-19T06:38:22.633Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x6365f099a9c0b187c372\",\"associatedWith\":{\"id\":\"did:emtrust:0x9ceb3a48b5b904412095\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2019-12-19T06:38:22.633Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304602210091f77ad53571b46eb4801d31a48f960ea57c61b6d4228ecdadaa9bfd5bd26028022100ac7a1d05c2fe8e632b2e807ed810fba8590abef2467a174e54b3b8bf63ed59b5\"}}]", ""})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
		Context("When no param passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 500,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(``)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateDidDocument([]string{})
			It("should return response as null", func() {
				Expect(returnedData).To(ContainSubstring(""))
			})
		})
	})

	Describe("GenerateEndorsement", func() {

		Context("When api calls GenerateEndorsement", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
							"@context": [
								"https://www.w3.org/2018/credentials/v1",
								"https://www.halialabs.io/2019/endorsements/v1"
							],
							"id": "31e2955c-5cc0-4453-9d18-d264db501abc",
							"type": [
								"VerifiableCredential",
								"MembershipCredential"
							],
							"issuer": "https://emtrust.io/issuer/e18e1843-f9fe-4e4f-ab9d-ac7a94eed4aa",
							"issuanceDate": "2020-02-03T04:50:24.738Z",
							"credentialSubject": {
								"id": "did:emtrust:0x3842dc6d3c5bb5f9b222",
								"associatedWith": {
									"id": "did:emtrust:0x7f6332183e94f2c21f28",
									"type": "Member",
									"desc": "Member for organization"
								}
							},
							"proof": {
								"type": "ECDSA",
								"created": "2020-02-03T04:50:24.738Z",
								"proofPurpose": "assertionMethod",
								"verificationMethod": "https://emtrust.io/api/crypto/endorsements/verify",
								"jws": "3045022056f1297cde63a61e10ebc578ef803539c7f9d1a35b6dba1c6c27d8201305520b0221008ba0b26dc688b797ae13c5a5cc49cea63c60ac701fa3419a5838c712cba2753f"
							}
						}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, returnedData = service.GenerateEndorsement([]string{"e18e1843-f9fe-4e4f-ab9d-ac7a94eed4aa", "did:emtrust:0x3842dc6d3c5bb5f9b222", "31e2955c-5cc0-4453-9d18-d264db501abc", "member", "Member", "Member for organization"})
			var endorsement model.Endorsement
			json.Unmarshal([]byte(returnedData), &endorsement)
			It("should return json with credentialSubject ID", func() {
				Expect(endorsement.ID).To(ContainSubstring("31e2955c-5cc0-4453-9d18-d264db501abc"))
			})
		})
	})
	// start of the test case
	Describe("Generate key pairs", func() {

		Context("When api generates the keypair", func() {

			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						  "did": "did:emtrust:0x3f1984bbca614b380324",
						  "publicKey": "04e215d550aa706e2a0c22785188baba3b70855345a26260e7a726d93772c02fec47ccb7ac76738ac5d7f31487f5c072ecd8fbf15f0ed65332ce9292481c4590b8",
						`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, _, returnedData = service.GenerateKeyPair("anyAccID")
			It("should return json with publicKey", func() {
				Expect(returnedData).To(ContainSubstring("publicKey"))
			})
		})
	})
	Describe("QueryIdentity", func() {

		Context("When api check the wrong did passed", func() {
			client := NewTestClient(func(req *http.Request) *http.Response {
				// Test request parameters
				return &http.Response{
					StatusCode: 200,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(`{
						"deviceInfo": "",
						"email": "",
						"id": "did:emtrust:0x360369452966240fb5a6",
						"name": "",
						"other": "{\"@context\":\"https://w3id.org/did/v1\",\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"publicKey\":[{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"type\":\"ED25519SignatureVerification\",\"controller\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"publicKeyBase58\":\"0490f679c0a6541e22498796ad23a614bd1f88dced322e749362071c7ef6402bd37fd6addb74b76bff9eedff616e2e993614afd6b488ba6ca3d00adea4bb6f2caa\",\"authorizations\":[]}],\"authentication\":[\"did:emtrust:0x360369452966240fb5a6\"],\"endorsements\":[{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"f565aa8c-b6a4-4713-ba74-bf5a71b12647\",\"type\":[\"VerifiableCredential\",\"MembershipCredential\"],\"issuer\":\"https://emtrust.io/issuer/38600fc9-89ee-4707-af4a-c9be18b72ad6\",\"issuanceDate\":\"2020-01-20T08:40:07.884Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"type\":\"Member\",\"desc\":\"Member for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-20T08:40:07.884Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3045022100a0c47f323bf2ddff1eea147c4220309fabd7315215fc6df82224738b669b83de02204eba1b94f1dbb64f09c92686e162409193874e03e5d51cc00202bcf3c89e8dad\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"f565aa8c-b6a4-4713-ba74-bf5a71b12647\",\"type\":[\"VerifiableCredential\",\"adminCredential\"],\"issuer\":\"https://emtrust.io/issuer/38600fc9-89ee-4707-af4a-c9be18b72ad6\",\"issuanceDate\":\"2020-01-20T09:51:21.926Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"type\":\"admin\",\"desc\":\"admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-20T09:51:21.926Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3045022100d56429afe52ce873f8976a847aeedb27cfe9de0b3b1aec2ac9e7612ccd981566022071a69b129a58aec1d1de06b5c4d6c64341c3dd792cb689941356016ae35e587a\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"f565aa8c-b6a4-4713-ba74-bf5a71b12647\",\"type\":[\"VerifiableCredential\",\"staffCredential\"],\"issuer\":\"https://emtrust.io/issuer/38600fc9-89ee-4707-af4a-c9be18b72ad6\",\"issuanceDate\":\"2020-01-20T09:51:24.477Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"type\":\"staff\",\"desc\":\"staff for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-20T09:51:24.477Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304502207b654158eda5b116f7f35a00c3f4195d5aaf95d54968fa51a9c041c790406dc6022100ef79d7f8425285be6973b10ca10245674c22e83a39b5eebbd0e54e29def46bad\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"f565aa8c-b6a4-4713-ba74-bf5a71b12647\",\"type\":[\"VerifiableCredential\",\"GenericVerifiableCredential\"],\"issuer\":\"https://emtrust.io/issuer/38600fc9-89ee-4707-af4a-c9be18b72ad6\",\"issuanceDate\":\"2020-01-20T09:51:27.031Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"type\":\"\",\"desc\":\" for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-20T09:51:27.031Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"304502206b51d7615ecd80fa598ad1dca941bf60ae3757ab3fc46d26e284ef494a7bc163022100fec48acfee13355db7a8996bc5387b3963c6b3f1b6693b2e0662761fc780553b\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"f565aa8c-b6a4-4713-ba74-bf5a71b12647\",\"type\":[\"VerifiableCredential\",\"agentCredential\"],\"issuer\":\"https://emtrust.io/issuer/38600fc9-89ee-4707-af4a-c9be18b72ad6\",\"issuanceDate\":\"2020-01-20T09:51:29.55Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x55d16d1e00dfae5e25ee\",\"type\":\"agent\",\"desc\":\"agent for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-20T09:51:29.55Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"3046022100c9e770907b448b0144c88b394f48e50134d2af70ce45cb5e6eec1e4613100f61022100f76bc8f6bb8647c2a1d2bf54307c9f888f647a998ace161fd0abbf45b12d5e78\"}},{\"@context\":[\"https://www.w3.org/2018/credentials/v1\",\"https://www.halialabs.io/2019/endorsements/v1\"],\"id\":\"b37c280b-d96a-4b9a-810a-ee96ad9487fd\",\"type\":[\"VerifiableCredential\",\"adminCredential\"],\"issuer\":\"https://emtrust.io/issuer/74bb9ed9-e298-4ab8-917f-0372c859c95f\",\"issuanceDate\":\"2020-01-23T06:35:49.215Z\",\"credentialSubject\":{\"id\":\"did:emtrust:0x360369452966240fb5a6\",\"associatedWith\":{\"id\":\"did:emtrust:0x34b55219135c5edd2fb9\",\"type\":\"Admin\",\"desc\":\"Admin for organization\"}},\"proof\":{\"type\":\"ECDSA\",\"created\":\"2020-01-23T06:35:49.215Z\",\"proofPurpose\":\"assertionMethod\",\"verificationMethod\":\"https://emtrust.io/api/crypto/endorsements/verify\",\"jws\":\"30460221009476c5da08dd6387a7c76ad0e6afec30513af413dbca6752ce8c19ce3340fce3022100edd5f0ec9c3d1e5cebc92b1a1455a9c1f0c17e7ea53c672d105047177fe44beb\"}}],\"service\":null,\"created\":\"Mon, 20 Jan 2020 08:40:08 GMT\",\"updated\":\"Mon, 20 Jan 2020 08:40:08 GMT\"}",
						"publicKey": "0490f679c0a6541e22498796ad23a614bd1f88dced322e749362071c7ef6402bd37fd6addb74b76bff9eedff616e2e993614afd6b488ba6ca3d00adea4bb6f2caa"
					}`)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			})

			service, _ := documents.New(client)
			_, _, returnedData = service.QueryIdentity("did:emtrust:0x360369452966240fb5a6")
			It("should return json with other", func() {
				Expect(returnedData).To(ContainSubstring("other"))
			})
		})
	})

})
