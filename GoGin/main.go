package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/halialabsdev/go-quaker/publicdetailsview-service/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var router *gin.Engine

func Getport() string {
	port := os.Getenv("SERVICEPORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
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
func GetQr(c *gin.Context) {
	URL := os.Getenv("URL")
	if len(URL) == 0 {
		URL = "https://dev.emtrust.io/"
	}
	url := URL + "api/v1/auth/request"
	var CommonRes model.CommonResponse
	param := map[string]interface{}{"type": "App/USER_LOGIN", "payload": "login"}
	//	logrus.Info("args::" + arg)
	arg, err := json.Marshal(param)
	payload := bytes.NewBuffer(arg) //strings.NewReader(arg)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_REQ"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_READ"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	var QR_Res model.QR_Result
	err = json.Unmarshal([]byte(body), &QR_Res)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(QR_Res)
	return
}
func Querytxid(c *gin.Context) {
	Txid := c.Param("txid")
	URL := os.Getenv("URL")
	if len(URL) == 0 {
		URL = "https://dev.emtrust.io/"
	}
	url := URL + "api/v1/auth/request?txid=" + Txid
	var CommonRes model.CommonResponse
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_REQ"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_READ"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	//var QR_Res model.QR_Status
	// err = json.Unmarshal([]byte(body), &QR_Res)
	// if err != nil {
	// 	logrus.Error(err.Error())
	// 	CommonRes.Code = "ERR_READ"
	// 	CommonRes.Message = err.Error()
	// 	c.Writer.Header().Set("Content-Type", "application/json")
	// 	c.Writer.WriteHeader(http.StatusBadRequest)
	// 	json.NewEncoder(c.Writer).Encode(CommonRes)
	// 	return
	// }
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(body)
	return
}
func VerifyEmailStatus(c *gin.Context) {
	var CommonRes model.CommonResponse
	id := c.Param("id")
	logrus.Error(id)
	error := false
	if error {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Email not verified"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	} else {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Email verified by user"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
}

func VerifyEmail(c *gin.Context) {
	var CommonRes model.CommonResponse
	id := c.Param("id")
	logrus.Error(id)
	error := false
	if error {
		CommonRes.Code = "VERIFY_ERROR"
		CommonRes.Message = "Email already verified"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	} else {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Email verified successfully"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
}
func SendEmail(c *gin.Context) {
	var CommonRes model.CommonResponse
	AuthorizationToken := c.Request.Header["Authorization"]
	if AuthorizationToken == nil {
		logrus.Error("Authorization required")
		CommonRes.Code = "AUTH_ERROR"
		CommonRes.Message = "Authorization required"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	if strings.Contains(AuthorizationToken[0], "Bearer") {
		splitToken := strings.Split(AuthorizationToken[0], "Bearer ")
		if splitToken[1] != "" {
		}

		logrus.Error(splitToken)
	}
	request := model.SendEmailReq{}
	err := c.BindJSON(&request)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_PARSE"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	error := false
	if error {
		CommonRes.Code = "SEND_ERROR"
		CommonRes.Message = "Unable to send Email! Please try again later"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	} else {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "https://dev.emtrust.io/verify/1234567890"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
}
func SaveDetails(c *gin.Context) {
	did := c.Param("did")
	var CommonRes model.CommonResponse
	AuthorizationToken := c.Request.Header["Authorization"]
	if AuthorizationToken == nil {
		logrus.Error("Authorization required")
		CommonRes.Code = "AUTH_ERROR"
		CommonRes.Message = "Authorization required"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	s := strings.Split(AuthorizationToken[0], " ")
	if len(s) < 2 {
		logrus.Error("Authorization required")
		CommonRes.Code = "AUTH_ERROR"
		CommonRes.Message = "Invalid authorization data"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	//get public key
	publicKey := QueryPublicKey(s[1], "xyz")
	if publicKey == "" {
		logrus.Error("Public Key not found")
		CommonRes.Code = "ERR_IDENTITY"
		CommonRes.Message = "Public Key not found"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	if s[0] != "Lookup" {
		if ValidateSignature(publicKey, s[2], s[1]) {
			logrus.Error("Signature validation failed")
			CommonRes.Code = "ERR_VALIDATE"
			CommonRes.Message = "Signature validation failed"
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(c.Writer).Encode(CommonRes)
			return
		}
	}
	if s[1] != did {
		logrus.Error("Cannot add/update others DID details")
		CommonRes.Code = "ERR_DATA"
		CommonRes.Message = "Cannot add/update others DID details"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	db := c.MustGet("db").(*gorm.DB)
	student := model.PublicViewDetailsReq{}
	err := c.BindJSON(&student)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_PARSE"
		CommonRes.Message = err.Error()
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	if student.Name == "" {
		logrus.Error("Name cannot be null")
		CommonRes.Code = "ERR_DATA"
		CommonRes.Message = "Name cannot be null"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	if student.Email == "" {
		logrus.Error("Email cannot be null")
		CommonRes.Code = "ERR_DATA"
		CommonRes.Message = "Email cannot be null"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	db.Model(&model.PublicViewDetails{}).Where("did = ?", did).Delete(model.PublicViewDetails{})
	if err := db.Create(&model.PublicViewDetails{Did: did, Name: student.Name, Email: student.Email, Phone: student.Phone, About: student.About, Pic: student.Pic}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB_CONNECTION"
		CommonRes.Message = "Data base connection error"
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Writer).Encode(CommonRes)
		return
	}
	db.Model(&model.PublicViewEndorsement{}).Where("user_did = ?", did).Delete(model.PublicViewEndorsement{})
	for _, v := range student.Endorsement {
		if err := db.Create(model.PublicViewEndorsement{UserDid: did, AccountDid: v.AccountDid, Type: v.Type, Date: v.Date, Desc: v.Desc}).Error; err != nil {
			logrus.Error(err.Error())
			CommonRes.Code = "ERR_DB_CONNECTION"
			CommonRes.Message = "Data base connection error"
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(c.Writer).Encode(CommonRes)
			return
		}
	}

	CommonRes.Code = "SUCCESS"
	CommonRes.Message = "Saved"
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	json.NewEncoder(c.Writer).Encode(CommonRes)
	return
}
func QueryPublicKey(did string, Jwt string) string {
	URL := os.Getenv("URL")
	if len(URL) == 0 {
		URL = "https://dev.emtrust.io/"
	}
	url := URL + "api/blockchain/ident/" + did
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err.Error())
	}
	//	req.Header.Add("Authorization", "Bearer "+Jwt)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.Error(err.Error())
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return ""
	}

	var result model.Identreponse
	var Identresult model.IdentResponseResult
	json.Unmarshal([]byte(body), &result)
	if strings.Contains(result.Result, "Error") {
		return ""
	} else {
		json.Unmarshal([]byte(result.Result), &Identresult)
		if Identresult.Id == "" {
			return ""
		}
	}

	return Identresult.PublicKey
}
func ValidateSignature(PublicKey string, Signature string, did string) bool {
	URL := os.Getenv("URL")
	if len(URL) == 0 {
		URL = "https://dev.emtrust.io/"
	}
	param := map[string]interface{}{"publicKey": PublicKey, "signature": "[" + Signature + "]", "message": did}
	arg, err := json.Marshal(param)
	payload := bytes.NewBuffer(arg) //strings.NewReader(arg)
	url := URL + "api/crypto/certs/verifySignature"
	res, err := http.NewRequest("POST", url, payload)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Error(err.Error())
		return false
	}
	var result model.Validatereponse
	json.Unmarshal([]byte(body), &result)
	return result.Valid
}

// func FindAccountdetails(c *gin.Context) {
// 	Accountdetailsreq := model.Accountdetailsreq{}
// 	var CommonRes model.CommonResponse
// 	err := c.BindJSON(&Accountdetailsreq)
// 	if err != nil {
// 		logrus.Error(err.Error())
// 		CommonRes.Code = "ERR_PARSE"
// 		CommonRes.Message = err.Error()
// 		c.Writer.Header().Set("Content-Type", "application/json")
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(c.Writer).Encode(CommonRes)
// 		return
// 	}
// 	URL := os.Getenv("URL")
// 	if len(URL) == 0 {
// 		URL = "https://dev.emtrust.io/"
// 	}
// 	url := URL + "api/account-service/accounts/" + Accountdetailsreq.Did
// 	//	fmt.Println("service is", s)
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		logrus.Error(err.Error())
// 	}
// 	req.Header.Add("Authorization", "Bearer "+Accountdetailsreq.Jwt)
// 	res, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		logrus.Error(err.Error())
// 		CommonRes.Code = "ERR_REQ"
// 		CommonRes.Message = err.Error()
// 		c.Writer.Header().Set("Content-Type", "application/json")
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(c.Writer).Encode(CommonRes)
// 		return
// 	}
// 	defer res.Body.Close()
// 	body, err := ioutil.ReadAll(res.Body)
// 	if err != nil {
// 		logrus.Error(err.Error())
// 		CommonRes.Code = "ERR_READ"
// 		CommonRes.Message = err.Error()
// 		c.Writer.Header().Set("Content-Type", "application/json")
// 		c.Writer.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(c.Writer).Encode(CommonRes)
// 		return
// 	}
// 	c.Writer.Header().Set("Content-Type", "application/json")
// 	c.Writer.WriteHeader(http.StatusOK)
// 	c.Writer.Write(body)
// 	return
// }
func FindDid(c *gin.Context) {
	Did := c.Param("did")
	db := c.MustGet("db").(*gorm.DB)
	Count := 0
	details := []model.PublicViewDetails{}
	endorsement := []model.PublicViewEndorsementReq{}
	db.Model(&model.PublicViewDetails{}).Where("did=?", Did).Count(&Count)
	if Count == 0 {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(model.PublicViewDetailsReq{})
		return
	} else {

		db.Model(&model.PublicViewDetails{}).Where("did = ?", Did).Scan(&details)
		db.Model(&model.PublicViewEndorsement{}).Where("user_did = ?", Did).Scan(&endorsement)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)
		json.NewEncoder(c.Writer).Encode(model.PublicViewDetailsReq{Did: details[0].Did, Email: details[0].Email, Name: details[0].Name, Phone: details[0].Phone, Pic: details[0].Pic, About: details[0].About, Endorsement: endorsement})
		return
	}
}
func DisplayStartPage(c *gin.Context) {
	c.HTML(http.StatusOK, "StartPage.html", nil)
}
func DisplayprofilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "profile.html", nil)
}
func DisplaySavedDetails(c *gin.Context) {
	did := c.Query("did")
	c.HTML(http.StatusOK, "Details.html", did)
}
func main() {
	logrus.Info("In main function")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	logrus.Info(host)
	logrus.Info(port)
	logrus.Info(dbname)
	logrus.Info(username)
	logrus.Info(password)
	if host == "" {
		file, _ := ioutil.ReadFile("cred.json")
		var c model.Credentials
		json.Unmarshal(file, &c)
		host = c.Host
		port = c.Port
		dbname = c.Dbname
		username = c.User
		password = c.Password
		if err := godotenv.Load(); err != nil {
			logrus.Info("No .env file found")
		}
	}
	logrus.Info("host=" + host + " port=" + port + " user=" + username + " dbname=" + dbname + " password=" + password + " sslmode=disable")
	db, err := gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+dbname+" password="+password+" sslmode=disable")

	defer db.Close()
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	logrus.Info("Connection Established")
	if err := db.AutoMigrate(&model.PublicViewDetails{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.PublicViewEndorsement{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	r := gin.Default()
	r.Static("/images", "public/images")
	r.LoadHTMLGlob("public/*.html")
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})
	r.GET("/DIDLookup", DisplayStartPage)
	r.GET("/DIDDetail", DisplaySavedDetails)
	r.GET("/MyProfile", DisplayprofilePage)
	r.GET("/did/:did", FindDid)
	//r.POST("/Accountdetails", FindAccountdetails)
	r.POST("/did/:did", SaveDetails)
	r.GET("/GetQr", GetQr)
	r.GET("/Querytxid/:txid", Querytxid)
	r.POST("/Updatedetails", SaveDetails)
	r.POST("/verify", SendEmail)
	r.POST("/verify/:id", VerifyEmail)
	r.GET("/verify/:id", VerifyEmailStatus)

	logrus.Info("Starting event app on port " + Getport())
	r.Run(Getport())
}
