package documents

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"voucher-service/model"

	logrus "github.com/sirupsen/logrus"
)

const blockchainApiBasePath = "/api/blockchain/"
const cryptoApiBasePath = "/api/crypto/"

var BLOCKCHAIN_API_URL string = ""
var CRYPTO_API_URL string = ""
var ACCOUNT_SERVICE_API_URL string = ""
var EVENT_API_URL string = ""
var DebuggerStatus bool = false

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
func (s *Service) ClaimVouchers(args []string) (int, string) {
	url := BLOCKCHAIN_API_URL + "/dtc/issue"
	arg := "{\"fcn\":\"issue\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"," + "\"" + args[2] + "\"," + "\"" + args[3] + "\"," + "\"" + args[4] + "\"," + "\"" + args[5] + "\"," + "\"" + args[6] + "\"," + "\"" + args[7] + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("ClaimVouchers arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
		return 500, ""
	}
	req.Header.Add("content-type", "application/json")
	//	res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err)
		return 500, ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return 500, ""
	}
	if DebuggerStatus {
		logrus.Info("ClaimVouchers response body:" + string(body))
	}
	ResponseString := string(body)
	return res.StatusCode, ResponseString
}
func (s *Service) QueryTokensByOwner(did string) (int, []byte) {
	url := BLOCKCHAIN_API_URL + "/dtc/fetch"
	arg := "{\"fcn\":\"queryTokensByOwner\",\"args\":[\"" + did + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("QueryTokensByOwner arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)

	}
	if DebuggerStatus {
		logrus.Info("QueryTokensByOwner response body:" + string(body))
	}
	return res.StatusCode, body
}
func (s *Service) GetHistory(voucheruuid string) (int, []byte) {
	url := BLOCKCHAIN_API_URL + "/dtc/fetch"
	arg := "{\"fcn\":\"getHistory\",\"args\":[\"" + voucheruuid + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("GetHistory arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}
	if DebuggerStatus {
		logrus.Info("GetHistory response body:" + string(body))
	}
	//	Response, _ := json.Marshal(MockJson)
	// Convert bytes to string.
	//ResponseString := string(body)
	return res.StatusCode, body
}
func (s *Service) GetVoucherDetails(voucheruuid string) (int, []byte) {
	url := BLOCKCHAIN_API_URL + "/dtc/fetch"
	arg := "{\"fcn\":\"readEntry\",\"args\":[\"" + voucheruuid + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("GetVoucherDetails arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}
	if DebuggerStatus {
		logrus.Info("GetVoucherDetails response body:" + string(body))
	}
	return res.StatusCode, body
}
func (s *Service) QueryTokensByIssuer(did string) (int, []byte) {
	url := BLOCKCHAIN_API_URL + "/dtc/fetch"
	arg := "{\"fcn\":\"queryTokensByIssuer\",\"args\":[\"" + did + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("QueryTokensByIssuer arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	//res, err := http.DefaultClient.Do(req)
	res, err := s.client.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}
	if DebuggerStatus {
		logrus.Info("QueryTokensByIssuer response body:" + string(body))
	}
	//	Response, _ := json.Marshal(MockJson)
	// Convert bytes to string.
	//ResponseString := string(body)
	return res.StatusCode, body
}
func TransferVoucher(args []string) int {
	url := BLOCKCHAIN_API_URL + "/dtc/transfer"
	arg := "{\"fcn\":\"transfer\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("TransferVoucher arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	if DebuggerStatus {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("TransferVoucher response body:" + string(body))
	}
	return res.StatusCode
}
func UpdateVoucher(args []string) int {
	url := BLOCKCHAIN_API_URL + "/dtc/update"
	arg := "{\"fcn\":\"update\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"],\"channel\": \"tradechannel\", \"contract\": \"dtc\"}"
	if DebuggerStatus {
		logrus.Info("UpdateVoucher arg:" + arg)
	}
	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	if DebuggerStatus {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			logrus.Error(err)
		}
		logrus.Info("UpdateVoucher response body:" + string(body))
	}
	return res.StatusCode
}
func Addevent(args []string) int {
	url := EVENT_API_URL + "/publish"
	arg := "{\"source\":\"" + args[0] + "\",\"evt\":\"" + args[1] + "\",\"time\":\"" + args[2] + "\",\"actor\": \"" + args[3] + "\",\"text\":\"" + args[4] + "\"}"
	//arg := "{\"source\":\""+args[0]+"\",\"args\":[" + "\"" + args[0] + "\"," + "\"" + args[1] + "\"],\"peers\":[\"peer0.trade\"]}"

	payload := strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)

	//req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	return res.StatusCode
}
func QueryIdentity(did string) (int, []byte) {
	url := BLOCKCHAIN_API_URL + "/ident/" + did
	//	fmt.Println("service is", s)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err.Error())
		return 500, nil
	}
	//	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return 500, nil
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return 500, nil
	}
	if DebuggerStatus {
		logrus.Info("QueryIdentity response body:" + string(body))
	}
	var result model.Identreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.Contains(result.Result, "Error") {
		status = 500
	} else {
		status = 200
	}
	return status, []byte(result.Result)
}
func Getaccountdetails(uniqueList string) (error, []byte, int) {
	url := ACCOUNT_SERVICE_API_URL + "/accounts/" + uniqueList
	//url := "http://localhost:8080/accounts/" + uniqueList
	//fmt.Println("Getaccountdetails:: URL:" + url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err)
	}
	// Convert bytes to string.
	//ResponseString := string(body)
	return err, body, res.StatusCode
}
func QueryMembersdetail(memberid string) (int, []byte) {
	query := "{%22selector%22:{%22other%22:{%22$regex%22:%22MembershipCredential%22,%22$regex%22:%22" + memberid + "%22}}}"
	url := BLOCKCHAIN_API_URL + "/ident/fetch"
	param := map[string]interface{}{"fcn": "queryIdentities", "args": query, "channel": "identchannel", "contract": "ident"}
	arg, err := json.Marshal(param)
	if DebuggerStatus {
		logrus.Info("QueryMembersdetail args::" + string(arg))
	}
	payload := bytes.NewBuffer(arg) //strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err.Error())
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
	}
	if DebuggerStatus {
		logrus.Info("QueryMembersdetail response body:" + string(body))
	}
	var result model.Identreponse
	json.Unmarshal([]byte(body), &result)
	status := 200
	if strings.Contains(result.Result, "Error") {
		status = 500
	} else {
		status = 200
	}
	return status, []byte(result.Result)
}
