package documents

import (
	"account-service/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

var DebuggerStatus bool = false

const blockchainApiBasePath = "/api/blockchain/"
const cryptoApiBasePath = "/api/crypto/"

var DLAccessUrl2 string = os.Getenv("BLOCKCHAIN_API_ADDR")
var BLOCKCHAIN_API_URL string = os.Getenv("API_BLOCKCHAIN_URL")
var CRYPTO_API_URL string = os.Getenv("API_CRYPTO_URL")
var VOUCHER_SERVICE_API_URL string = os.Getenv("API_VOUCHER_SERVICE_URL")
var EVENT_API_URL string = os.Getenv("API_EVENT_URL")

type Service struct {
	client                *http.Client
	BlockchainApiBasePath string // API endpoint base URL
	CryptoApiBasePath     string
	UserAgent             string // optional additional User-Agent fragment
}

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			// this function is required when you want to introduce your custom format.
			// In my case I wanted file and line to look like this `file="engine.go:141`
			// but f.File provides a full path along with the file name.
			// So in `formatFilePath()` function I just trimmet everything before the file name
			// and added a line number in the end
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}
func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}
func New(client *http.Client) (*Service, error) {
	//fmt.Println("new is called", client)
	if client == nil {
		logrus.Error("client is nil")
		return nil, errors.New("client is nil")
	}
	URL := os.Getenv("URL")
	s := &Service{client: client, BlockchainApiBasePath: URL + blockchainApiBasePath, CryptoApiBasePath: URL + cryptoApiBasePath}
	//fmt.Println("returning new service", s)
	return s, nil
}

func (s *Service) GenerateKeyPair(accountId string) (int, []byte, string) {
	url := CRYPTO_API_URL + "/certs/generateKeyPair"
	//	fmt.Println("service is", s)
	//url := s.CryptoApiBasePath + "/certs/generateKeyPair"
	arg := "{\"id\":\"" + accountId + "\"}"
	//	logrus.Info("args::" + arg)
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return 500, nil, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	ResponseString := string(body)
	if DebuggerStatus {
		logrus.Info("GenerateKeyPair response body:" + ResponseString)
	}
	return res.StatusCode, body, ResponseString
}
func (s *Service) QueryIdentity(did string) (int, []byte, string) {
	url := BLOCKCHAIN_API_URL + "/ident/" + did
	//	fmt.Println("service is", s)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err.Error())
		return 500, nil, ""
	}
	//	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return 500, nil, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return 500, nil, ""
	}
	var result model.Identreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.Contains(result.Result, "Error") {
		status = 500
	} else {
		status = 200
	}
	ResponseStringg := string(body)
	if DebuggerStatus {
		logrus.Info("QueryIdentity response body:" + ResponseStringg)
	}
	return status, []byte(result.Result), ResponseStringg
}
func (s *Service) GenerateEndorsement(args []string) (int, string) {
	url := CRYPTO_API_URL + "/endorsements/create"
	arg := "{\"id\":\"" + args[0] + "\",\"endorsementType\":\"" + args[3] + "\",\"beneficiaryId\":\"" + args[1] + "\",\"selfAssociationId\":\"" + args[2] + "\",\"associationType\":\"" + args[4] + "\",\"associationDesc\":\"" + args[5] + "\"}"
	if DebuggerStatus {
		logrus.Info("GenerateEndorsement arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//	res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return 500, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	ResponseString := string(body)
	if DebuggerStatus {
		logrus.Info("GenerateEndorsement response body:" + ResponseString)
	}
	return res.StatusCode, ResponseString
}
func (s *Service) GenerateDidDocument(args []string) (int, string) {
	if len(args) == 0 || len(args) != 5 {
		return 500, ""
	}
	if args[0] == "" || args[1] == "" || args[2] == "" || args[3] == "" || args[4] == "" {
		return 500, ""
	}
	url := CRYPTO_API_URL + "/certs/generateDidDocument"
	arg := "{\"did\":\"" + args[0] + "\",\"publickey\":\"" + args[1] + "\",\"controllerdid\":\"" + args[2] + "\",\"endorsements\":" + strings.Trim(args[3], "\"") + ",\"entitytype\":\"" + args[4] + "\",\"pubkeyasauth\":\"Y\"}"
	if DebuggerStatus {
		logrus.Info("GenerateDidDocument arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	ResponseString := string(body)
	if DebuggerStatus {
		logrus.Info("GenerateDidDocument response body:" + ResponseString)
	}
	return res.StatusCode, ResponseString
}
func (s *Service) RegisterDidDocument(args []string) (int, string) {
	if len(args) == 0 || len(args) != 6 {
		return 500, ""
	}
	if args[0] == "" {
		return 500, ""
	}
	if args[1] == "" {
		return 500, ""
	}
	url := BLOCKCHAIN_API_URL + "/ident/register"
	//arg := []string{}
	//arg := "{\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"," + "\"" + args[2] + "\"," + "\"" + args[3] + "\"," + "\"" + args[4] + "\","  + "\"" + args[5] + "\"" +"]}"
	param := map[string]interface{}{"fcn": "register", "args": args, "channel": "identchannel", "contract": "ident"}
	arg, err := json.Marshal(param)
	if DebuggerStatus {
		logrus.Info("RegisterDidDocument arg:" + string(arg))
	}
	payload := bytes.NewBuffer(arg) //strings.NewReader(string(arg))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	var result model.Registerreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.ToLower(result.Result) == "successs" {
		status = 200
	} else {
		status = 500
	}
	ResponseString := string(body)
	if DebuggerStatus {
		logrus.Info("RegisterDidDocument response body:" + ResponseString)
	}
	return status, ResponseString
}
func (s *Service) UpdateDidDocument(args []string) (int, string) {
	if len(args) == 0 {
		return 500, ""
	}
	if args[0] == "" && args[1] == "" {
		return 500, ""
	}
	if args[0] == "" {
		return 500, ""
	}
	url := BLOCKCHAIN_API_URL + "/ident"
	//arg := "{\"fcn\":\"update\",\"peers\":\"peer0.trade\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"]}"
	param := map[string]interface{}{"fcn": "update", "args": args, "channel": "identchannel", "contract": "ident"}
	arg, err := json.Marshal(param)
	if DebuggerStatus {
		logrus.Info("UpdateDidDocument args:" + string(arg))
	}
	payload := bytes.NewBuffer(arg) //strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	var result model.Registerreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.ToLower(result.Result) == "successs" {
		status = 200
	} else {
		status = 500
	}
	ResponseString := string(body)
	if DebuggerStatus {
		logrus.Info("UpdateDidDocument response body:" + ResponseString)
	}
	return status, ResponseString
}

func (s *Service) QueryMembersdetail(memberid string) (int, []byte, string) {
	query := "{\"selector\":{\"other\":{\"$regex\":\"MembershipCredential\",\"$regex\":\"" + memberid + "\"}}}"
	url := BLOCKCHAIN_API_URL + "/ident/fetch"
	param := map[string]interface{}{"fcn": "queryIdentities", "args": [1]string{query}, "channel": "identchannel", "contract": "ident"}
	arg, err := json.Marshal(param)
	if DebuggerStatus {
		logrus.Info("QueryMembersdetail args:" + string(arg))
	}
	payload := bytes.NewBuffer(arg) //strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	var result model.Identreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.Contains(result.Result, "Error") {
		status = 500
	} else {
		status = 200
	}
	ResponseStringg := string(body)
	if DebuggerStatus {
		logrus.Info("QueryMembersdetail response body:" + ResponseStringg)
	}
	return status, []byte(result.Result), ResponseStringg
}
func (s *Service) Addevent(args []string) (int, string) {
	if len(args) == 0 || len(args) != 5 {
		return 200, ""
	}

	url := EVENT_API_URL + "/publish"
	arg := "{\"source\":\"" + args[0] + "\",\"evt\":\"" + args[1] + "\",\"time\":\"" + args[2] + "\",\"actor\": \"" + args[3] + "\",\"text\":\"" + args[4] + "\"}"
	//arg := "{\"source\":\""+args[0]+"\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"],\"peers\":[\"peer0.trade\"]}"
	fmt.Println(arg)
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return res.StatusCode, ""
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return res.StatusCode, ""
	}
	ResponseStringg := string(body)
	return res.StatusCode, ResponseStringg
}
