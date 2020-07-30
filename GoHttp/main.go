package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	documents "account-service/document"
	"account-service/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	logrus "github.com/sirupsen/logrus"
)

var db *gorm.DB
var err error
var service *documents.Service
var BLOCKCHAIN_API_URL string = ""
var CRYPTO_API_URL string = ""
var VOUCHER_SERVICE_API_URL string = ""
var EVENT_API_URL string = ""
var DebuggerStatus bool = false

type Credentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

func addevents(source string, evt string, actor string, text string) {
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	issuedtime := strconv.FormatInt(umillisec, 10)
	eventargs := []string{source, evt, issuedtime, actor, text}
	service.Addevent(eventargs)
}
func createtempaccount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	InvitationId := params["InvitationId"]
	var CommonRes model.CommonResponse
	if InvitationId == "" || InvitationId == "undefined" {
		CommonRes.Code = "ERR_INVALID_ID"
		CommonRes.Message = "InvitationId cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var status string
	if db.Where("invitation_id = ?", InvitationId).First(&model.TempAccountsDetail{}).Pluck("Status", status).RecordNotFound() {
		logrus.Info("No record found")
	} else {
		logrus.Error("InvitationId already used")
		CommonRes.Code = "ERR_INVALID_ID"
		CommonRes.Message = "InvitationId already used"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	u1 := uuid.NewV4()
	now := time.Now()
	after := now.AddDate(0, 0, +14)
	words := strings.Fields(after.String())
	words[1] = "23:59:59"
	if err := db.Create(&model.TempAccountsDetail{TemporaryId: u1.String(), InvitationId: InvitationId, ExpiryDate: words[0] + " " + words[1] + " " + words[2] + " " + words[3], Status: "initiated"}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB_CONNECTION"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var Responsedata model.TempAccountsResponse
	Responsedata.TemporaryAccountId = u1.String()
	Responsedata.ExpiryDate = words[0] + " " + words[1] + " " + words[2] + " " + words[3]
	// var res1 string
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Responsedata)
}
func enrolltoEmtrust(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	memberdid := params["member"]
	var CommonRes model.CommonResponse
	if memberdid == "" || memberdid == "undefined" {
		CommonRes.Code = "ERR_INVALID_DID"
		CommonRes.Message = "memberdid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var accountid []string
	if err := db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").Pluck("account_id", &accountid).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if accountid == nil {
		CommonRes.Code = "ERR_ACCOUNT_ID_INVALID"
		CommonRes.Message = "System account not present"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var accountDid []string
	if err := db.Model(&model.AccountsDetail{}).Where("account_id = ?", accountid[0]).Pluck("account_did", &accountDid).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var c model.EnrolltoEmtrustReq
	err = json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_PARSE"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.PublicKey == "" {
		CommonRes.Code = "ERR_PUB_KEY"
		CommonRes.Message = "PublicKey cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

	u1 := uuid.NewV4()
	MemberId := u1.String()
	endorsement, err := populateEndorsement(accountid[0], memberdid, MemberId, "member", "Member", "Member for organization")
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_ENDORSEMENT"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

	// get identity of user if already exist
	identity, err := fetchDidDocumentForUser(memberdid)
	if "" != identity.Other { // append new endorsements
		for _, v := range identity.DID.Endorsements {
			if strings.Contains(v.Issuer, accountid[0]) && strings.ToLower(v.CredentialSubject.AssociatedWith.Type) == "member" {
				logrus.Error("Already a member")
				CommonRes.Code = "208"
				CommonRes.Message = "Already a member"
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAlreadyReported)
				json.NewEncoder(w).Encode(CommonRes)
				return
			}
		}
		identity.DID.Endorsements = append(identity.DID.Endorsements, endorsement)
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, string(didAsBytes)}, true)
		addevents("account-service", "ENROLL_MEMBER", memberdid, strings.Join([]string{identity.Id, string(didAsBytes)}, ","))

	} else { // generate new endorsement
		identity, err = generateIdentityStructure(memberdid, c.PublicKey, accountDid[0], endorsement, "member", "Y")
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, false)
		addevents("account-service", "ENROLL_MEMBER", memberdid, strings.Join([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, ","))
	}

	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_BC"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersDetail{AccountID: accountid[0], MemberID: MemberId, Name: "", Managed: "Internal", ExternalMemberId: "", ExternalMembershipUrl: ""}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersRole{MemberID: MemberId, TypeOfRole: "Member"}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersRequiredParams{AccountID: accountid[0], MemberID: MemberId, Name: "", Email: "", Address: "", Phone: "", Sex: "", Age: "", DOB: ""}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	// on success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.EnrolltoEmtrustRes{IssuerDid: accountDid[0], Diddocument: identity.DID})
}
func enrollmembers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationdid := params["account"]
	var CommonRes model.EnrollCommonResponse
	if organizationdid == "nil" || organizationdid == "undefined" {
		var accountdid []string
		if err := db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").Pluck("account_did", &accountdid).Error; err != nil {
			logrus.Error(err.Error())
			CommonRes.Status = http.StatusBadRequest
			CommonRes.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
		if accountdid == nil {
			logrus.Error("System account not present")
			CommonRes.Status = http.StatusBadRequest
			CommonRes.Message = "System account not present"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
		organizationdid = accountdid[0]
	} else {
		var enroll []bool
		if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", organizationdid).Pluck("enroll", &enroll).Error; err != nil {
			logrus.Error(err.Error())
			CommonRes.Status = http.StatusBadRequest
			CommonRes.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
		if len(enroll) == 0 {
			logrus.Error("Organization not found")
			CommonRes.Status = http.StatusBadRequest
			CommonRes.Message = "Organization not found"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
		if !enroll[0] {
			CommonRes.Status = http.StatusBadRequest
			CommonRes.Message = "Cannot enroll to this organization"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
	}
	var c model.EnrollMembersReq
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.MemberDid == "" {
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = "MemberDid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.PublicKey == "" {
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = "PublicKey cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	u1 := uuid.NewV4()
	MemberId := u1.String()
	var accountid []string

	if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", organizationdid).Pluck("account_id", &accountid).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if len(accountid) == 0 {
		logrus.Error("Organization not found")
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = "Organization not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	// populate new endorsement for admin user
	endorsement, err := populateEndorsement(accountid[0], c.MemberDid, MemberId, "member", "Member", "Member for organization")

	// get identity of user if already exist
	identity, err := fetchDidDocumentForUser(c.MemberDid)
	if "" != identity.Other { // append new endorsements
		for _, v := range identity.DID.Endorsements {
			if strings.Contains(v.Issuer, accountid[0]) && strings.ToLower(v.CredentialSubject.AssociatedWith.Type) == "member" {
				CommonRes.Status = http.StatusAlreadyReported
				CommonRes.Message = "Already a member"
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAlreadyReported)
				json.NewEncoder(w).Encode(CommonRes)
				return
			}
			if strings.Contains(v.Issuer, accountid[0]) && strings.ToLower(v.CredentialSubject.AssociatedWith.Type) == "admin" {
				// populate new endorsement for admin user
				endorsement, err = populateEndorsement(accountid[0], c.MemberDid, v.ID, "member", "Member", "Member for organization")

			}
		}
		identity.DID.Endorsements = append(identity.DID.Endorsements, endorsement)
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, string(didAsBytes)}, true)
		addevents("account-service", "ENROLL_MEMBER", c.MemberDid, strings.Join([]string{identity.Id, string(didAsBytes)}, ","))
	} else { // generate new endorsement
		identity, err = generateIdentityStructure(c.MemberDid, c.PublicKey, organizationdid, endorsement, "member", "Y")
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, false)
		addevents("account-service", "ENROLL_MEMBER", c.MemberDid, strings.Join([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, ","))
	}

	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersDetail{AccountID: accountid[0], MemberID: MemberId, Name: c.Name, Managed: c.Managed, ExternalMemberId: c.ExternalMemberId, ExternalMembershipUrl: c.ExternalMembershipUrl}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersRole{MemberID: MemberId, TypeOfRole: "Member"}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Create(&model.MembersRequiredParams{AccountID: accountid[0], MemberID: MemberId, Name: c.RequiredParams.Name, Email: c.RequiredParams.Email, Address: c.RequiredParams.Address, Phone: c.RequiredParams.Phone, Sex: c.RequiredParams.Sex, Age: c.RequiredParams.Age, DOB: c.RequiredParams.DOB}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Status = http.StatusBadRequest
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	// on success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.EnrolltoEmtrustRes{IssuerDid: organizationdid, Diddocument: identity.DID})
}

/*
Creating account for the organization and assign admin
Parameters Required
	- DID : user's decentralized ID
	- PublicKey : user's public KEY
	- OtherDetails : extra details if any
	- DeviceID : user's device id (//TODO: why its even required?
	- TemporaryAccountId : temporary account id, linked with invite
	- BusinessName : Name of the business
	- PreferredSite : creating unique subdomain on emtrust
Steps for Account Flow
	1. Create Account keys and did for account.
	2. Generate Endorsements for admin role for users did.
	3. Prepare or Get existing DID document for user.
	4. Append new endorsement to did document.
	5. Prepare updated DID document and get it signed.
	6. Update did document on blockchain.
	7. Return with updated did.
*/
func createAccount(w http.ResponseWriter, r *http.Request) {

	bodyDecoder := json.NewDecoder(r.Body)
	var accountsDetailReq model.AccountsDetailReq
	err := bodyDecoder.Decode(&accountsDetailReq)
	var CommonRes model.CommonResponse
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DECODE"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	//TODO: validate all mandatory input

	// validate the invitation of account creation
	err = validateTemporaryAccount(accountsDetailReq.TemporaryAccountId)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_INVALID_ID"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

	// setup the account
	businessDid, businessPublicKey, didDocument, _int_biz_id, _intr_admin_id, err := setupAccountKeyPairs(accountsDetailReq.BusinessName)
	if businessDid == "" {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_ID_CREATION"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	didAsBytes, _ := json.Marshal(didDocument)
	// persist organisation did on the chain
	err = persistIdentityOnChain([]string{businessDid, businessPublicKey, "", "", "", string(didAsBytes)}, false)
	addevents("account-service", "CREATE_ORGANIZATION", accountsDetailReq.DID, strings.Join([]string{businessDid, businessPublicKey, "", "", "", string(didAsBytes)}, ","))
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	// populate new endorsement for admin user
	endorsement, err := populateEndorsement(_int_biz_id, accountsDetailReq.DID, _intr_admin_id, "admin", "Admin", "Admin for organization")

	// get identity of user if already exist
	identity, err := fetchDidDocumentForUser(accountsDetailReq.DID)

	if "" != identity.Other { // append new endorsements
		identity.DID.Endorsements = append(identity.DID.Endorsements, endorsement)
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, string(didAsBytes)}, true)
		addevents("account-service", "CREATE_ORGANIZATION", accountsDetailReq.DID, strings.Join([]string{identity.Id, string(didAsBytes)}, ","))
	} else { // generate new endorsement
		identity, err = generateIdentityStructure(accountsDetailReq.DID, accountsDetailReq.PublicKey, businessDid, endorsement, "admin", "Y")
		didAsBytes, _ := json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, false)
		addevents("account-service", "CREATE_ORGANIZATION", accountsDetailReq.DID, strings.Join([]string{identity.Id, identity.PublicKey, "", "", "", string(didAsBytes)}, ","))
	}
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	//as last update lets update table records
	sysAccountFlag := strings.Contains(strings.ToLower(accountsDetailReq.BusinessName), "emtrust")
	linkAccountWithDid(_int_biz_id, businessDid, identity.Id, accountsDetailReq.BusinessName, accountsDetailReq.PreferredSite, sysAccountFlag, accountsDetailReq.TemporaryAccountId, accountsDetailReq.Background)

	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// on success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.EnrolltoEmtrustRes{IssuerDid: businessDid, Diddocument: identity.DID})
}

/*
Validate temporary account, return error is not exist
*/
func validateTemporaryAccount(tempAccountId string) error {
	var count int64
	db.Where("temporary_id = ?", tempAccountId).First(&model.TempAccountsDetail{}).Count(&count)
	if count == 0 {
		logrus.Error("invalid invitation")
		return errors.New("invalid invitation")
	}
	db.Model(&model.TempAccountsDetail{}).Where("temporary_id = ? AND status = ?", tempAccountId, "completed").Count(&count)
	if count == 1 {
		logrus.Error("temporary_id already used")
		return errors.New("temporary_id already used")
	}
	if db.Where(&model.TempAccountsDetail{TemporaryId: tempAccountId, Status: "initiated"}).RecordNotFound() {
		logrus.Error("invalid invitation")
		return errors.New("invalid invitation")
	}
	return nil
}

/*
Setup Accounts Private and Public Keys
*/
func setupAccountKeyPairs(name string) (string, string, model.DIDDocument, string, string, error) {

	businessName := strings.ToLower(name)
	// check for system account or not
	if strings.Contains(businessName, "emtrust") {
		var count int64
		db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").Count(&count) // check if any record exist
		if count >= 1 {
			logrus.Error("system account already exist")
			return "", "", model.DIDDocument{}, "", "", errors.New("system account already exist")
		}
	}

	// assign new random uuid for account id
	_db_acc_id := uuid.NewV4().String()
	db.Create(&model.AccountsDetail{AccountID: _db_acc_id, Status: "Draft"})

	// assign new random for member id and link account id
	_db_member_id := uuid.NewV4().String()
	db.Create(&model.MembersDetail{AccountID: _db_acc_id, MemberID: _db_member_id})

	//role is required?
	// db.Create(&model.MembersRole{MemberID: _db_member_id, TypeOfRole: "admin"})

	_, _api_resp, _api_resp_string := service.GenerateKeyPair(_db_acc_id)
	if DebuggerStatus {
		logrus.Info("Response received for keypair::" + _api_resp_string)
	}
	var business_did model.GenerateKeyPairResponse
	err := json.Unmarshal(_api_resp, &business_did)

	if err != nil {
		logrus.Error("unable to persist the record on the chain")
		return "", "", model.DIDDocument{}, "", "", errors.New("unable to generate keys for new account")
	}

	return business_did.Did, business_did.PublicKey, business_did.DidDocument, _db_acc_id, _db_member_id, nil
}

func persistIdentityOnChain(args []string, updateFlag bool) error {
	var statusCode int

	if updateFlag {
		statusCode, _ = service.UpdateDidDocument(args)
	} else {
		log.Println("Registering new identity...", args)
		statusCode, statusStr := service.RegisterDidDocument(args)
		log.Println("Registration response is now... ", statusCode, statusStr)
	}

	if statusCode == 500 {
		logrus.Error("unable to persist the record on the chain")
		return errors.New("unable to persist the record on the chain")
	}

	return nil
}

/*
 Fetching and preparing endorsements
*/
func populateEndorsement(issuerDid string, beneficiaryDid string, membershipId string, endorsementType string, associationType string, associationDesc string) (model.Endorsement, error) {

	_, endorsementJson := service.GenerateEndorsement([]string{issuerDid, beneficiaryDid, membershipId, endorsementType, associationType, associationDesc})
	//log.Println("received endorsements from server", endorsementJson)
	var endorsement model.Endorsement
	err := json.Unmarshal([]byte(endorsementJson), &endorsement)

	if err != nil {
		logrus.Error(err.Error())
		return model.Endorsement{}, errors.New("unable to populate endorsement")
	}

	return endorsement, nil
}

/*
 Fetch member details from the memberid
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchmemberdetailsfromId(id string) (model.DIDDocument, error) {

	_, identJson, _ := service.QueryMembersdetail(id)
	var details []model.MemberData
	var did model.DIDDocument
	err := json.Unmarshal(identJson, &details)
	if err != nil {
		logrus.Error(err.Error())
		return model.DIDDocument{}, err
	}
	if details[0].Record.Other != "" {

		err = json.Unmarshal([]byte(details[0].Record.Other), &did)
	}
	if err != nil {
		logrus.Error(err.Error())
		return model.DIDDocument{}, err
	}
	return did, nil
}

/*
 Fetch member enrolled from the account id
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchmemberlist(id string) ([]string, error) {

	status, identJson, _ := service.QueryMembersdetail(id)
	if status != 200 {
		logrus.Error("QueryMembersdetail error")
		return nil, err
	}
	var details []model.MemberData
	err := json.Unmarshal(identJson, &details)
	if err != nil {
		logrus.Error(err.Error())
		return nil, err
	}
	membersDid := []string{}
	for _, v := range details {
		membersDid = append(membersDid, v.Key)
	}
	return membersDid, nil
}

/*
 Fetch did document from the chain
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchDidDocumentForUser(did string) (model.IdentData, error) {

	_, identJson, _api_resp_string := service.QueryIdentity(did)
	if DebuggerStatus {
		logrus.Info("Response received for QueryIdentity::" + _api_resp_string)
	}

	var identity model.IdentData = model.IdentData{}
	err := json.Unmarshal(identJson, &identity)
	if err != nil {
		return model.IdentData{}, err
	}

	if identity.Other != "" {
		var did model.DIDDocument
		err = json.Unmarshal([]byte(identity.Other), &did)
		identity.DID = did
	}

	return identity, nil
}

/*
 Create the Identity Structure with just did and public key with new endorsement
*/
func generateIdentityStructure(did string, publicKey string, controllerDid string, endorsement model.Endorsement, entityType string, keyAsAuth string) (model.IdentData, error) {

	//log.Println("generating identity structure with endorsement", endorsement)
	endorsementsAsString, err := json.Marshal(endorsement)
	if err != nil {
		logrus.Error("unable to marshal the endorsements into json")
		return model.IdentData{}, err
	}
	args := []string{did, publicKey, controllerDid, string(endorsementsAsString), entityType}
	_, generatedDidString := service.GenerateDidDocument(args) // TODO: it need to support passing entity type and keyAsAuth

	var didDocument model.DIDDocument
	err = json.Unmarshal([]byte(generatedDidString), &didDocument)

	if err != nil {
		logrus.Error("unable to unmarshal the did string into endorsements::" + generatedDidString)
		return model.IdentData{}, err
	}

	identity := model.IdentData{Id: did, PublicKey: publicKey, Other: generatedDidString, DID: didDocument}

	return identity, nil

}

/*
Linking account with org did.. //TODO: this should also change the way its being linked
*/
func linkAccountWithDid(businessAccId string, businessDid string, createdBy string, businessName string, preferredSite string, systemAccountFlag bool, temporaryAccountId string, Background string) {

	db.Model(&model.TempAccountsDetail{}).Where("temporary_id = ?", temporaryAccountId).Update("permanent_account_id", businessAccId)
	//update account details
	if systemAccountFlag {
		db.Model(&model.AccountsDetail{}).Where("account_id = ?", businessAccId).Updates(model.AccountsDetail{Status: "active", AccountDid: businessDid, CreatedBy: createdBy, BusinessName: businessName, AccountName: businessName, PreferredSite: preferredSite, Systemaccount: "Y", Background: Background})
	} else {
		db.Model(&model.AccountsDetail{}).Where("account_id = ?", businessAccId).Updates(model.AccountsDetail{Status: "active", AccountDid: businessDid, CreatedBy: createdBy, BusinessName: businessName, AccountName: businessName, PreferredSite: preferredSite, Systemaccount: "N", Background: Background})
	}
	//update tempaccount
	db.Model(&model.TempAccountsDetail{}).Where("permanent_account_id = ?", businessAccId).Update("status", "completed")

}

func createpermanentaccount(w http.ResponseWriter, r *http.Request) {
	var c model.AccountsDetailReq
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var CommonRes model.CommonResponse
	if c.DID == "" {
		logrus.Error("DID cannot be null")
		CommonRes.Code = "ERR_INVALID_DID"
		CommonRes.Message = "DID cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.PublicKey == "" {
		logrus.Error("PublicKey cannot be null")
		CommonRes.Code = "ERR_PUBLIC_KEY"
		CommonRes.Message = "PublicKey cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.TemporaryAccountId == "" {
		logrus.Error("TemporaryAccountId cannot be null")
		CommonRes.Code = "ERR_TEMP_ACC_ID"
		CommonRes.Message = "TemporaryAccountId cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.BusinessName == "" {
		logrus.Error("BusinessName cannot be null")
		CommonRes.Code = "ERR_BUSINESS_NAME"
		CommonRes.Message = "BusinessName cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if c.PreferredSite == "" {
		logrus.Error("PreferredSite cannot be null")
		CommonRes.Code = "ERR_PRFERRED_SITE"
		CommonRes.Message = "PreferredSite cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	resId, res, _api_resp_string := service.QueryIdentity(c.DID)
	log.Println("Response received for keypair", _api_resp_string)
	if resId != 200 {
		logrus.Error("QueryIdentity error")
		CommonRes.Code = "ERR_IDENT_API"
		CommonRes.Message = "QueryIdentity api error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var dataa model.IdentData
	err = json.NewDecoder(strings.NewReader(string(res))).Decode(&dataa)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_GOLANG_LIB"
		CommonRes.Message = "NewDecoder error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	Enrolled := false
	if dataa.PublicKey == "" && dataa.DeviceInfo == "" {
		Enrolled = false
	} else {
		Enrolled = true
	}
	business := strings.ToLower(c.BusinessName)
	systemacc := "N"
	if strings.Contains(business, "emtrust") {
		systemacc = "Y"
		if !db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").First(&model.AccountsDetail{}).RecordNotFound() {
			logrus.Error("System account already present")
			CommonRes.Code = "ERR_SYS_ACC"
			CommonRes.Message = "System account already present"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
	}

	if db.Where("temporary_id = ?", c.TemporaryAccountId).First(&model.TempAccountsDetail{}).RecordNotFound() {
		logrus.Error("temporary_id not registered")
		CommonRes.Code = "ERR_TEMP_ID"
		CommonRes.Message = "temporary_id not registered"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else {
		var status model.TempaacountResponse
		if err := db.Model(&model.TempAccountsDetail{}).Where("temporary_id = ?", c.TemporaryAccountId).Scan(&status).Error; err != nil {
			logrus.Error(err.Error())
			CommonRes.Code = "ERR_DATABASE"
			CommonRes.Message = "database error"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
		if status.Status == "completed" {
			w.WriteHeader(http.StatusAlreadyReported)
			fmt.Fprintf(w, "temporary_id already used")
			return
		}
	}
	u1 := uuid.NewV4()
	if err := db.Create(&model.AccountsDetail{AccountID: u1.String(), Status: "Draft"}).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Data base connection error")
		return
	}
	if err := db.Model(&model.TempAccountsDetail{}).Where("temporary_id = ?", c.TemporaryAccountId).Update("permanent_account_id", u1.String()).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	u2 := uuid.NewV4()
	if err := db.Create(&model.MembersDetail{AccountID: u1.String(), MemberID: u2.String()}).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Data base connection error")
		return
	}
	if err := db.Create(&model.MembersRole{MemberID: u2.String(), TypeOfRole: "admin"}).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Data base connection error")
		return
	}
	resId, res11, resbody := service.GenerateKeyPair(u1.String())
	if resId != 200 {
		logrus.Error("GenerateKeyPair error")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "GenerateKeyPair error")
		return
	}
	fmt.Println("//GenerateKeyPair output//")
	fmt.Println(resbody)
	var data model.GenerateKeyPairResponse
	json.Unmarshal([]byte(res11), &data)
	args := []string{u1.String(), c.DID, u2.String()}
	resId1, res1 := service.GenerateEndorsement(args)
	if resId1 != 200 {
		logrus.Error("GenerateEndorsement error")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "GenerateEndorsement error")
		return
	}
	args2 := []string{c.DID, c.PublicKey, data.Did, res1}
	resId2, res2 := service.GenerateDidDocument(args2)
	if resId2 != 200 {
		logrus.Error("GenerateDidDocument error")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "GenerateDidDocument error")
		return
	}
	fmt.Println("//GenerateDidDocument output//")
	fmt.Println(res2)
	args3 := []string{data.Did, data.PublicKey, " ", " ", " ", strings.ReplaceAll(resbody, "\"", "\\\"")}
	for i := 0; i < 3; i++ {
		resId3, res3 := service.RegisterDidDocument(args3)
		if resId3 == 504 {
			logrus.Info("Organization register timeout")
			continue
		} else if resId3 == 200 {
			logrus.Info("Organization register sucess")
			break
		} else {
			logrus.Error("Organization register error")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Organization register error")
			return
		}
		fmt.Println("//Organization register output//")
		fmt.Println(res3)
	}
	if !Enrolled {
		args4 := []string{c.DID, c.PublicKey, " ", " ", " ", strings.ReplaceAll(res2, "\"", "\\\"")}

		for i := 0; i < 3; i++ {
			resId4, res4 := service.RegisterDidDocument(args4)
			if resId4 == 504 {
				logrus.Info("RegisterDidDocument error")
				continue
			} else if resId4 == 200 {
				logrus.Info("RegisterDidDocument error")
				break
			} else {
				logrus.Error("RegisterDidDocument error")
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "RegisterDidDocument error")
				return
			}
			fmt.Println("//individual register output//")
			fmt.Println(res4)
		}
		resId9, res9, _api_resp_string := service.QueryIdentity(c.DID)
		if DebuggerStatus {
			logrus.Info("Response received for keypair::" + _api_resp_string)
		}

		if resId9 != 200 {
			logrus.Error("QueryIdentity error")
			fmt.Fprintf(w, "QueryIdentity error")
			return
		}
		var dataa9 model.IdentData
		err = json.NewDecoder(strings.NewReader(string(res9))).Decode(&dataa9)
		if err != nil {
			logrus.Error(err.Error())
			fmt.Fprintf(w, "NewDecoder error")
			return
		}
		dataa.Other = dataa9.Other
	} else {
		fmt.Println(dataa.Other)
		index := strings.Index(dataa.Other, "endorsements")
		newOthers := dataa.Other[:(index+15)] + args2[3] + "," + dataa.Other[(index+15):]
		dataa.Other = newOthers
		fmt.Println(dataa.Other)
		args5 := []string{c.DID, strings.ReplaceAll(dataa.Other, "\"", "\\\"")}
		resId5, _ := service.UpdateDidDocument(args5)
		if resId5 != 200 {
			logrus.Error("UpdateDidDocument error")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "UpdateDidDocument error")
			return
		}
		fmt.Println("//individual update output//")
		fmt.Println(dataa.Other)
	}

	var EnrollRepsone model.EnrolltoEmtrustRes
	EnrollRepsone.IssuerDid = data.Did
	EnrollRepsone.Diddocument = dataa.DID
	if err := db.Model(&model.AccountsDetail{}).Where("account_id = ?", u1.String()).Updates(model.AccountsDetail{Status: "active", AccountDid: data.Did, CreatedBy: c.DID, BusinessName: c.BusinessName, AccountName: c.BusinessName, PreferredSite: c.PreferredSite, Systemaccount: systemacc}).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := db.Model(&model.TempAccountsDetail{}).Where("permanent_account_id = ?", u1.String()).Update("status", "completed").Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(EnrollRepsone)
}
func tempaccountaccess(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	TempaccountId := params["tempaccount"]
	if TempaccountId == "" || TempaccountId == "undefined" {
		logrus.Error("TempaccountId cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "TempaccountId cannot be null")
		return
	}
	var status model.TempaacountResponse
	if db.Where("temporary_id = ?", TempaccountId).First(&model.TempAccountsDetail{}).Scan(&status).RecordNotFound() {
		logrus.Error("No record found")
		fmt.Fprintf(w, "No record found")
		return
	}

	if err := db.Model(&model.TempAccountsDetail{}).Where("temporary_id = ?", TempaccountId).Scan(&status).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var Responsedata model.CreateaccountResponse
	if status.Status == "completed" {
		var accountdetails model.PermanentacountResponse
		if err := db.Model(&model.AccountsDetail{}).Where("account_id = ?", status.PermanentAccountID).Scan(&accountdetails).Error; err != nil {
			logrus.Error(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		Responsedata.Created = true
		Responsedata.Updateddate = status.UpdatedAt.String()
		Responsedata.AccountDid = accountdetails.AccountDid
		Responsedata.UserDid = accountdetails.CreatedBy
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Responsedata)
}
func getaccountdetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	if params == nil {
		logrus.Error("wrong did passed")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "wrong did passed")
		return
	}
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	var status model.AccountsDetail

	if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", accountDid).Scan(&status).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Error": "Specified Account not found"})
		return
	}
	AccountFeature := []int{}
	if err := db.Model(&model.AccountFeaturesTable{}).Where("account_did = ?", accountDid).Pluck("paramcode", &AccountFeature).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Error": "Specified Account not found"})
		return
	}
	MembershipRequiredFields := []string{}
	if status.MembershipRequiredFields != "" {
		MembershipRequiredFields = strings.Split(status.MembershipRequiredFields, ",")
	}
	var Response model.GetAccountsDetailResponse
	Response.ID = status.ID
	Response.CreatedAt = status.CreatedAt
	Response.UpdatedAt = status.UpdatedAt
	Response.DeletedAt = status.DeletedAt
	Response.AccountID = status.AccountID
	Response.AccountDid = status.AccountDid
	Response.AccountName = status.AccountName
	Response.BusinessName = status.BusinessName
	Response.PreferredSite = status.PreferredSite
	Response.Address = status.Address
	Response.Latitude = status.Latitude
	Response.Longitude = status.Longitude
	Response.Logo = status.Logo
	Response.Status = status.Status
	Response.CreatedBy = status.CreatedBy
	Response.Email = status.Email
	Response.MembershipURL = status.MembershipURL
	Response.ExternalMember = status.ExternalMember
	Response.ExternalMembershipUrl = status.ExternalMembershipUrl
	Response.MembershipRequiredFields = MembershipRequiredFields
	Response.Enroll = status.Enroll
	Response.Background = status.Background
	Response.AccountFeatures = AccountFeature
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response)
}
func setaccountdetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	var c model.Setaccountsdetailreq
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if c.BusinessName == "" {
		logrus.Error("BusinessName cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "BusinessName cannot be null")
		return
	}
	if c.Address == "" {
		logrus.Error("Address cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Address cannot be null")
		return
	}
	if c.Latitude == "" {
		logrus.Error("Latitude cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Latitude cannot be null")
		return
	}
	if c.Longitude == "" {
		logrus.Error("Longitude cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Longitude cannot be null")
		return
	}
	if c.Email == "" {
		logrus.Error("Email cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Email cannot be null")
		return
	}
	var account_details model.AccountsDetail
	if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", accountDid).Find(&account_details).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	account_details.AccountName = c.BusinessName
	account_details.BusinessName = c.BusinessName
	account_details.Address = c.Address
	account_details.Latitude = c.Latitude
	account_details.Longitude = c.Longitude
	account_details.Logo = c.Logo
	account_details.Email = c.Email
	account_details.MembershipURL = c.MembershipURL
	account_details.ExternalMember = c.ExternalMember
	account_details.ExternalMembershipUrl = c.ExternalMembershipUrl
	account_details.MembershipRequiredFields = c.MembershipRequiredFields
	account_details.Enroll = c.Enroll
	account_details.Background = c.Background
	if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", accountDid).Save(&account_details).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var response model.SetaccountdetailsResponse
	response.Message = "Updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func getenrollparams(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationid := params["account"]
	if organizationid == "" || organizationid == "undefined" {
		logrus.Error("organizationid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "organizationid cannot be null")
		return
	}
	var requiredfields []string
	if err := db.Model(&model.AccountsDetail{}).Where("account_id = ?", organizationid).Pluck("membership_required_fields", &requiredfields).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if len(requiredfields) == 0 {
		logrus.Error("organizationid not found")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "organizationid not found")
		return
	}
	MembershipRequiredFields := []string{}
	if requiredfields[0] != "" {
		MembershipRequiredFields = strings.Split(requiredfields[0], ",")
	}
	var resposne model.GetenrollparamsResposne
	resposne.RequiredParam = MembershipRequiredFields
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resposne)
}
func FindIdent(did string) model.IdentResponse {
	var Response model.IdentResponse
	resId, res, _api_resp_string := service.QueryIdentity(did)
	if DebuggerStatus {
		logrus.Info("Response received for keypair::" + _api_resp_string)
	}
	if resId != 200 {
		logrus.Error("QueryIdentity error")
		Response.Status = resId
		Response.Enrolled = false
		Response.Message = "Failed to get response"
		return Response
	}
	var data model.IdentData
	err := json.NewDecoder(strings.NewReader(string(res))).Decode(&data)
	if err != nil {
		logrus.Error(err.Error())
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Failed to decode json"
		return Response
	}
	if data.PublicKey == "" && data.DeviceInfo == "" {
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Not Enrolled"
	} else {
		Response.Status = 200
		Response.Enrolled = true
		Response.Message = "Enrolled"
	}
	return Response
}
func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func GetMembersIdList(MemberList []string, accountid string) ([]string, error) {
	Membersid := []string{}
	//Get all member id
	for _, v := range MemberList {
		// get identity of user if already exist
		identity, err := fetchDidDocumentForUser(v)
		if err != nil {
			return nil, err
		}
		if "" != identity.Other {
			for _, v := range identity.DID.Endorsements {
				if strings.Contains(v.Issuer, accountid) {
					Membersid = append(Membersid, v.ID)
				}
			}
		}
	}
	uniqueMembersid := []string{}
	if len(Membersid) != 0 {
		uniqueMembersid = unique(Membersid)
	}
	return uniqueMembersid, nil
}
func getallmembers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationdid := params["account"]
	var CommonRes model.CommonResponse
	if organizationdid == "" || organizationdid == "undefined" {
		logrus.Error("organizationid cannot be null")
		CommonRes.Code = "ERR_ORG_ID"
		CommonRes.Message = "organizationid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	accountId := []string{}
	if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", organizationdid).Pluck("account_id", &accountId).Error; err != nil {
		logrus.Error("Data base error")
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "Data base error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if len(accountId) == 0 {
		logrus.Error("accountId not found")
		CommonRes.Code = "ERR_DB"
		CommonRes.Message = "accountId not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	MemberList, err1 := fetchmemberlist(accountId[0])
	if err1 != nil {
		logrus.Error(err1.Error())
		CommonRes.Code = "ERR_MEMBER_LIST"
		CommonRes.Message = err1.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	membersid, err := GetMembersIdList(MemberList, accountId[0])
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_MEMBER_LIST"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	memberdetails := []model.GetallmembersRes{}
	sqlStatement := "SELECT members_required_params.member_id,members_required_params.account_id,members_required_params.name,members_required_params.phone,members_required_params.address,members_required_params.sex,members_required_params.dob,members_required_params.age,members_required_params.email FROM members_required_params FULL OUTER JOIN members_details ON members_details.member_id = members_required_params.member_id WHERE members_required_params.member_id=?"
	var (
		member_id  string
		account_id string
		name       string
		phone      string
		address    string
		sex        string
		dob        string
		age        string
		email      string
	)
	var details []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
	EndorsementTypeMap := make(map[string]int)
	for _, v := range details {
		EndorsementTypeMap[strings.ToLower(v.ParamName)] = v.ParamCode
	}
	for _, v := range membersid {
		rows, err := db.Raw(sqlStatement, v).Rows() // (*sql.Rows, error)
		defer rows.Close()
		count := 0
		for rows.Next() {
			count++
			rows.Scan(&member_id, &account_id, &name, &phone, &address, &sex, &dob, &age, &email)
			endorsementype, Did, err := GetEndorsementType(member_id)
			if err != nil {
				continue
			}
			EndorsemenTypeint := []int{}
			for _, v := range endorsementype {
				EndorsemenTypeint = append(EndorsemenTypeint, EndorsementTypeMap[strings.ToLower(v)])
			}
			memberdetails = append(memberdetails, model.GetallmembersRes{MemberID: member_id, MemberDID: Did, AccountID: account_id, Name: name, Phone: phone, Address: address, Sex: sex, DOB: dob, Age: age, Email: email, EndorsementType: EndorsemenTypeint})
		}
		if count == 0 {
			endorsementype, Did, err := GetEndorsementType(v)
			if err != nil {
				continue
			}
			EndorsemenTypeint := []int{}
			for _, v := range endorsementype {
				EndorsemenTypeint = append(EndorsemenTypeint, EndorsementTypeMap[strings.ToLower(v)])
			}
			memberdetails = append(memberdetails, model.GetallmembersRes{MemberID: v, MemberDID: Did, AccountID: accountId[0], Name: "", Phone: "", Address: "", Sex: "", DOB: "", Age: "", Email: "", EndorsementType: EndorsemenTypeint})
			continue
		}
		if err != nil {
			logrus.Error(err.Error())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberdetails)
	return
}
func getmembersdetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationid := params["account"]
	var CommonRes model.CommonResponse
	var memberdetails model.GetallmembersRes
	if organizationid == "" || organizationid == "undefined" {
		logrus.Error("organizationid cannot be null")
		CommonRes.Code = "ERR_ORG_ID"
		CommonRes.Message = "organizationid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	memberid := params["member"]
	if memberid == "" || memberid == "undefined" {
		logrus.Error("memberid cannot be null")
		CommonRes.Code = "ERR_ORG_ID"
		CommonRes.Message = "memberid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var ReqParams model.MembersRequiredParams
	if db.Model(&model.MembersRequiredParams{}).Where("member_id = ? AND account_id =  ?", memberid, organizationid).Scan(&ReqParams).RecordNotFound() {
		logrus.Error("database error")
		CommonRes.Code = "ERR_DATABASE"
		CommonRes.Message = "database error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var details []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
	EndorsementTypeMap := make(map[string]int)
	for _, v := range details {
		EndorsementTypeMap[strings.ToLower(v.ParamName)] = v.ParamCode
	}
	endorsementype, Did, err := GetEndorsementType(memberid)
	if err != nil {
		logrus.Error(err.Error())
	}
	EndorsemenTypeint := []int{}
	for _, v := range endorsementype {
		EndorsemenTypeint = append(EndorsemenTypeint, EndorsementTypeMap[strings.ToLower(v)])
	}
	memberdetails.MemberID = ReqParams.MemberID
	memberdetails.MemberDID = Did
	memberdetails.AccountID = ReqParams.AccountID
	memberdetails.Name = ReqParams.Name
	memberdetails.Phone = ReqParams.Phone
	memberdetails.Email = ReqParams.Email
	memberdetails.Address = ReqParams.Address
	memberdetails.Sex = ReqParams.Sex
	memberdetails.DOB = ReqParams.DOB
	memberdetails.Age = ReqParams.Age
	memberdetails.EndorsementType = EndorsemenTypeint
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(memberdetails)
}
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func Containsinteger(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
func updateEndorsementType(organizationid string, memberid string, EndorsementType []int, EndorsementTypeAdd bool) error {
	did, err := fetchmemberdetailsfromId(memberid)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	memberdid := did.ID
	Typepresent := []string{}
	for _, v := range did.Endorsements {
		if v.ID == memberid {
			Typepresent = append(Typepresent, strings.ToLower(v.CredentialSubject.AssociatedWith.Type))
		} else {

		}
	}
	var details []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
	EndorsemenTypestring := []string{}
	EndorsementTypeMap := make(map[int]string)
	for _, v := range details {
		EndorsementTypeMap[v.ParamCode] = v.ParamName
	}
	for _, v := range EndorsementType {
		if EndorsementTypeMap[v] == "" {
			continue
		}
		EndorsemenTypestring = append(EndorsemenTypestring, strings.ToLower(EndorsementTypeMap[v]))
	}
	TypeNotpresent := []string{}
	//Remove endorsement
	if !EndorsementTypeAdd {
		//check endorsement type present or not
		for _, EndorsemenType := range EndorsemenTypestring {
			if !Contains(Typepresent, strings.ToLower(EndorsemenType)) {
				logrus.Error("EndorsemenType not present")
				return errors.New("EndorsemenType not present")
			}
		}

	} else {
		for _, EndorsemenType := range EndorsemenTypestring {
			if !Contains(Typepresent, strings.ToLower(EndorsemenType)) {
				TypeNotpresent = append(TypeNotpresent, strings.ToLower(EndorsemenType))
			}
		}
	}
	if EndorsementTypeAdd {
		for _, EndorsemenType := range TypeNotpresent {
			endorsement, err := populateEndorsement(organizationid, memberdid, memberid, strings.ToLower(EndorsemenType), EndorsemenType, EndorsemenType+" for organization")
			if err != nil {
				logrus.Error(err.Error())
				return err
			}

			identity, err := fetchDidDocumentForUser(memberdid)
			if err != nil {
				logrus.Error(err.Error())
				return err
			}
			if "" != identity.Other { // append new endorsements
				identity.DID.Endorsements = append(identity.DID.Endorsements, endorsement)
				didAsBytes, _ := json.Marshal(identity.DID)
				// finally update it on the chain
				err = persistIdentityOnChain([]string{identity.Id, string(didAsBytes)}, true)
				addevents("account-service", "UPDATE_MEMBER_DETAIL", memberdid, strings.Join([]string{identity.Id, string(didAsBytes)}, ","))
			}
		}
	} else {
		var didAsBytes []byte
		Endorsements := []model.Endorsement{}
		identity, err := fetchDidDocumentForUser(memberdid)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
		if "" != identity.Other { // append new endorsements
			for _, v := range identity.DID.Endorsements {
				if v.ID == memberid && !Contains(EndorsemenTypestring, strings.ToLower(v.CredentialSubject.AssociatedWith.Type)) {
					Endorsements = append(Endorsements, v)
				} else if v.ID != memberid {
					Endorsements = append(Endorsements, v)
				}
			}
		}
		identity.DID.Endorsements = Endorsements
		didAsBytes, _ = json.Marshal(identity.DID)
		// finally update it on the chain
		err = persistIdentityOnChain([]string{identity.Id, string(didAsBytes)}, true)
		addevents("account-service", "UPDATE_MEMBER_DETAIL", memberdid, strings.Join([]string{identity.Id, string(didAsBytes)}, ","))
	}
	if err != nil {
		logrus.Error(err.Error())
		return err
	} else {
		return nil
	}
}
func GetEndorsementType(memberid string) ([]string, string, error) {
	did, err := fetchmemberdetailsfromId(memberid)
	if err != nil {
		logrus.Error(err.Error())
		return nil, "", err
	}
	Typepresent := []string{}
	var Did string
	for _, v := range did.Endorsements {
		if v.ID == memberid {
			Typepresent = append(Typepresent, strings.ToLower(v.CredentialSubject.AssociatedWith.Type))
			Did = v.CredentialSubject.ID
		}
	}
	return Typepresent, Did, nil
}
func updatemembersdetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationid := params["account"]
	var CommonRes model.CommonResponse
	if organizationid == "" || organizationid == "undefined" {
		logrus.Error("organizationdid cannot be null")
		CommonRes.Code = "ERR_ORG_ID"
		CommonRes.Message = "organizationid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	memberid := params["member"]
	if memberid == "" || memberid == "undefined" {
		logrus.Error("memberid cannot be null")
		CommonRes.Code = "ERR_MEMBER_ID"
		CommonRes.Message = "memberid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

	var GetDetailsRes model.UpdateMembersDetailsReq
	err := json.NewDecoder(r.Body).Decode(&GetDetailsRes)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(GetDetailsRes.EndorsementType) != 0 {
		err := updateEndorsementType(organizationid, memberid, GetDetailsRes.EndorsementType, GetDetailsRes.EndorsementTypeAdd)
		if err != nil {
			logrus.Error(err.Error())
			CommonRes.Code = "ERR_TYPE_UPDATE"
			CommonRes.Message = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		} else {
			CommonRes.Code = "SUCCESS"
			CommonRes.Message = "record updated"
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}
	}
	var count int64
	db.Model(&model.MembersRequiredParams{}).Where("member_id = ? AND account_id =  ?", memberid, organizationid).Count(&count) // check if any record exist
	if count == 0 {
		logrus.Error("memberid or organizationid not found")
		CommonRes.Code = "ERR_ID_NOT_FOUND"
		CommonRes.Message = "memberid or organizationid not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Model(&model.MembersRequiredParams{}).Where("member_id = ? AND account_id =  ?", memberid, organizationid).Updates(model.MembersRequiredParams{Name: GetDetailsRes.Name, Email: GetDetailsRes.Email, Phone: GetDetailsRes.Phone, Address: GetDetailsRes.Address, Sex: GetDetailsRes.Sex, DOB: GetDetailsRes.DOB, Age: GetDetailsRes.Age}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DATABASE"
		CommonRes.Message = "database error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	CommonRes.Code = "SUCCESS"
	CommonRes.Message = "record updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CommonRes)
	return
}
func validateaccounts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	organizationdid := params["account"]
	var CommonRes model.CommonResponse
	if organizationdid == "" || organizationdid == "undefined" {
		logrus.Error("organizationdid cannot be null")
		CommonRes.Code = "ERR_ORG_ID"
		CommonRes.Message = "organizationdid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	identity, err := fetchDidDocumentForUser(organizationdid)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DID_DOC"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if identity.Other == "" {
		logrus.Error("Orgnization did is invalid")
		CommonRes.Code = "ERR_ORG_DID_INVALID"
		CommonRes.Message = "Orgnization did is invalid"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else {
		membershipservice := false
		for _, v := range identity.DID.Service {
			if v.Type == "MembershipService" {
				membershipservice = true
				break
			}
		}
		if membershipservice {
			var accountname []string
			if err := db.Model(&model.AccountsDetail{}).Where("account_did = ?", organizationdid).Pluck("account_name", &accountname).Error; err != nil {
				logrus.Error(err.Error())
				CommonRes.Code = "ERR_DATABASE"
				CommonRes.Message = err.Error()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(CommonRes)
				return
			}
			CommonRes.Code = "SUCCESS"
			CommonRes.Message = accountname[0]
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(CommonRes)
			return
		} else {
			logrus.Error("Orgnization did is invalid")
			CommonRes.Code = "ERR_ORG_DID_INVALID"
			CommonRes.Message = "Orgnization did is invalid"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(CommonRes)
			return
		}

	}
}

// func init() {
// 	formatter := runtimee.Formatter{ChildFormatter: &logg.JSONFormatter{}}
// 	formatter.Line = true
// 	formatter.Package = true
// 	logg.SetFormatter(&formatter)
// 	logg.SetOutput(os.Stdout)
// 	logg.SetLevel(logg.InfoLevel)
// }
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
func getaccountfeatures(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	AccountFeature := []model.AccountFeaturesTable{}
	if err := db.Model(&model.AccountFeaturesTable{}).Where("account_did = ?", accountDid).Scan(&AccountFeature).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Error": "Specified Account not found"})
		return
	}
	var List []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Feature_Type").Scan(&List)
	Paramlist := make(map[int]string)
	for _, v := range List {
		Paramlist[v.ParamCode] = v.ParamName
	}
	AccountFeatureRes := []model.AccountFeaturesTableRes{}
	for _, v := range AccountFeature {
		AccountFeatureRes = append(AccountFeatureRes, model.AccountFeaturesTableRes{AssociatedAt: v.UpdatedAt, AccountDid: v.AccountDid, Paramcode: v.Paramcode, Paramname: Paramlist[v.Paramcode], AssociatedBy: v.AssociatedBy, UnitPrice: v.UnitPrice, EffectiveStartDate: v.EffectiveStartDate, Status: v.Status, Recurring: v.Recurring})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountFeatureRes)
	return
}
func getparamlist(w http.ResponseWriter, r *http.Request) {
	var Transferresponse model.CommonResponse
	queryValues := r.URL.Query()
	Paramtype := queryValues.Get("type")
	if Paramtype == "" {
		s := strconv.Itoa(http.StatusBadRequest)
		Transferresponse.Code = s
		Transferresponse.Message = "Param type should not be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var List []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", Paramtype).Scan(&List)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(List)
	return
}
func addaccountfeatures(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	params := mux.Vars(r)
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	var FeatureReq model.AddFeature
	err := json.NewDecoder(r.Body).Decode(&FeatureReq)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if FeatureReq.FeatureType != 0 {
		if !FeatureReq.FeatureTypeAdd {

			db.Model(&model.AccountFeaturesTable{}).Where("account_did = ? AND paramcode=?", accountDid, FeatureReq.FeatureType).Delete(model.AccountFeaturesTable{})

		} else {
			count := 0
			db.Model(&model.AccountFeaturesTable{}).Where("account_did = ? AND paramcode=?", accountDid, FeatureReq.FeatureType).Count(&count)
			if count == 0 {
				if err := db.Create(&model.AccountFeaturesTable{AccountDid: accountDid, Paramcode: FeatureReq.FeatureType, AssociatedBy: FeatureReq.Did, UnitPrice: FeatureReq.UnitPrice, EffectiveStartDate: FeatureReq.EffectiveStartDate, Status: FeatureReq.Status, Recurring: FeatureReq.Recurring}).Error; err != nil {
					logrus.Error(err.Error())
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Data base connection error")
					return
				}
			}
		}
	}
	CommonRes.Code = "SUCCESS"
	CommonRes.Message = "record updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CommonRes)
	return
}
func updateaccountfeatures(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	params := mux.Vars(r)
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	featurecode := params["feature"]
	if featurecode == "" || featurecode == "undefined" {
		logrus.Error("featurecode cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "featurecode cannot be null")
		return
	}
	var FeatureReq model.UpdateFeaturesTableReq
	err := json.NewDecoder(r.Body).Decode(&FeatureReq)
	if err != nil {
		logrus.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := db.Model(&model.AccountFeaturesTable{}).Where("account_did = ? AND paramcode=?", accountDid, featurecode).Updates(model.AccountFeaturesTable{UnitPrice: FeatureReq.UnitPrice, EffectiveStartDate: FeatureReq.EffectiveStartDate, Status: FeatureReq.Status, Recurring: FeatureReq.Recurring}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_DATABASE"
		CommonRes.Message = "database error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	CommonRes.Code = "SUCCESS"
	CommonRes.Message = "record updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CommonRes)
	return
}
func Getquickcontact(SearchString string, authorization string) (int, []model.MemberContact) {
	if authorization == "" {
		logrus.Error("Authorization token cannot be null")
		return 401, []model.MemberContact{}
	}
	s := strings.Split(authorization, " ")
	if len(s) < 2 {
		logrus.Error("Authorization token cannot be null")
		return 401, []model.MemberContact{}
	}
	token, _, err := new(jwt.Parser).ParseUnverified(s[1], jwt.MapClaims{})
	if err != nil {
		logrus.Error(err.Error())
		return 401, []model.MemberContact{}
	}
	did := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		did = fmt.Sprintf("%v", claims["aud"])
	} else {
		logrus.Error(err.Error())
		return 401, []model.MemberContact{}
	}
	if did == "" {
		logrus.Error("Invalid Authorization token")
		return 401, []model.MemberContact{}
	}
	identity, err := fetchDidDocumentForUser(did)
	if err != nil {
		logrus.Error(err.Error())
		return 500, []model.MemberContact{}
	}
	var accountid []string
	memberid := ""
	db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").Pluck("account_id", &accountid)
	if "" != identity.Other { // append new endorsements
		for _, v := range identity.DID.Endorsements {
			if strings.Contains(v.Issuer, "b4c846bb-0866-47a3-9be7-c553ad400aa5") { //update it when deployed
				memberid = v.ID
			}
		}
	}
	if memberid == "" {
		return 500, []model.MemberContact{}
	}
	Count := 0
	QuickContact := []model.MemberContact{}
	db.Model(&model.MemberContact{}).Where("member_id=?", memberid).Count(&Count)
	if Count == 0 {
		return 200, []model.MemberContact{}
	} else {

		db.Model(&model.MemberContact{}).Where("member_id=? AND active=?", memberid, "active").Scan(&QuickContact)
	}
	QuickContactRes := []model.MemberContact{}
	for _, v := range QuickContact {
		if strings.Contains(v.Name, strings.ToLower(SearchString)) || strings.Contains(v.Email, strings.ToLower(SearchString)) || strings.Contains(v.NickName, strings.ToLower(SearchString)) {
			QuickContactRes = append(QuickContactRes, model.MemberContact(v))
		}
	}
	if len(QuickContactRes) == 0 {
		return 400, []model.MemberContact{}
	}
	return 200, QuickContactRes
}
func sealfortest(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	sealtest := model.SealtestReq{}
	err := json.NewDecoder(r.Body).Decode(&sealtest)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_PARSE"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	errlevel := 200
	if errlevel == 400 {
		CommonRes.Code = "SEAL_ERR"
		CommonRes.Message = "Unable to seal test file for user"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Filestream xxxxx"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

}
func gettestdid(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	errlevel := 400
	if errlevel == 400 {
		CommonRes.Code = "NOT_FOUND_ERR"
		CommonRes.Message = "Test bot DID not found"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Test bot DID XXXXX"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

}

func getContact(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	queryValues := r.URL.Query()
	authorization := r.Header.Get("Authorization")

	SearchString := queryValues.Get("SearchString")
	code, QuickContactRes := Getquickcontact(SearchString, authorization)
	if code == 500 {
		logrus.Error("Interenal server error")
		CommonRes.Code = "500"
		CommonRes.Message = "Interenal server error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if code == 400 {
		logrus.Error("No contact present for search string")
		CommonRes.Code = "400"
		CommonRes.Message = "No contact present for search string"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if code == 401 {
		logrus.Error("Invalid Authorization token")
		CommonRes.Code = "401"
		CommonRes.Message = "Invalid Authorization token"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(QuickContactRes)
	return
}
func Addquickcontact(contact model.QuickContactReq, authorization string, did string) (int, model.MemberContact) {
	if authorization == "" {
		logrus.Error("Authorization token cannot be null")
		return 401, model.MemberContact{}
	}
	s := strings.Split(authorization, " ")
	if len(s) < 2 {
		logrus.Error("Authorization token cannot be null")
		return 401, model.MemberContact{}
	}
	if did == "" {
		logrus.Error("did cannot be null")
		return 400, model.MemberContact{}
	}
	token, _, err := new(jwt.Parser).ParseUnverified(s[1], jwt.MapClaims{})
	if err != nil {
		logrus.Error(err.Error())
		return 401, model.MemberContact{}
	}
	Jwtdid := ""
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		Jwtdid = fmt.Sprintf("%v", claims["aud"])
	} else {
		logrus.Error(err.Error())
		return 401, model.MemberContact{}
	}
	if Jwtdid == "" {
		logrus.Error("Invalid Authorization token")
		return 401, model.MemberContact{}
	}
	if did != Jwtdid {
		logrus.Error("Invalid Authorization token")
		return 401, model.MemberContact{}
	}
	//Get member id
	identity, err := fetchDidDocumentForUser(did)
	if err != nil {
		logrus.Error(err.Error())
		return 500, model.MemberContact{}
	}
	var accountid []string
	memberid := ""
	db.Model(&model.AccountsDetail{}).Where("systemaccount = ?", "Y").Pluck("account_id", &accountid)
	if "" != identity.Other { // append new endorsements
		for _, v := range identity.DID.Endorsements {
			if strings.Contains(v.Issuer, "b4c846bb-0866-47a3-9be7-c553ad400aa5") { //update it when deployed
				memberid = v.ID
			}
		}
	}
	if memberid == "" {
		return 400, model.MemberContact{}
	}
	if contact.Name == "" {
		return 400, model.MemberContact{}
	}
	if contact.Email == "" {
		return 400, model.MemberContact{}
	}
	if contact.Contact_did == "" {
		ContactId := uuid.NewV4()
		InvitationId := uuid.NewV4()
		loc, _ := time.LoadLocation("UTC")
		CreatedAt := time.Now().In(loc)
		//	db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", did, contact.Contact_did).Delete(model.MemberContact{})
		if err := db.Create(&model.InvitationTable{InvitationId: InvitationId.String(), CreatedAt: CreatedAt, UpdatedAt: CreatedAt, InvitationUrl: "/verification/" + InvitationId.String()}).Error; err != nil {
			logrus.Error(err.Error())
			return 500, model.MemberContact{}
		}
		if err := db.Create(&model.MemberContact{ContactId: ContactId.String(), InvitationId: InvitationId.String(), CreatedAt: CreatedAt, UpdatedAt: CreatedAt, MemberId: memberid, ContactDid: contact.Contact_did, Name: contact.Name, Email: contact.Email, NickName: contact.NickName, Active: "inactive"}).Error; err != nil {
			logrus.Error(err.Error())
			return 500, model.MemberContact{}
		}
		return 201, model.MemberContact{}
	}
	if contact.Remove {
		if contact.NickName == "Delete" {
			var InvitationId []string
			db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", memberid, contact.Contact_did).Pluck("invitation_id", &InvitationId)
			if InvitationId[0] != "" {
				db.Model(&model.InvitationTable{}).Where("invitation_id = ?", InvitationId[0]).Delete(&model.InvitationTable{})
			}
			db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", memberid, contact.Contact_did).Delete(&model.MemberContact{})
		} else {
			db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", memberid, contact.Contact_did).Update("active", "inactive")
		}
	} else {
		ContactId := uuid.NewV4()
		loc, _ := time.LoadLocation("UTC")
		CreatedAt := time.Now().In(loc)
		db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", memberid, contact.Contact_did).Delete(model.MemberContact{})
		if err := db.Create(&model.MemberContact{ContactId: ContactId.String(), CreatedAt: CreatedAt, UpdatedAt: CreatedAt, MemberId: memberid, ContactDid: contact.Contact_did, Name: contact.Name, Email: contact.Email, NickName: contact.NickName, Active: "active"}).Error; err != nil {
			logrus.Error(err.Error())
			return 500, model.MemberContact{}
		}
	}
	var Contact model.MemberContact
	db.Model(&model.MemberContact{}).Where("member_id = ? AND contact_did=?", memberid, contact.Contact_did).Scan(&Contact)

	return 200, Contact
}
func mergeContact(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	params := mux.Vars(r)
	did := params["did"]
	contact := []model.QuickContactReq{}
	err := json.NewDecoder(r.Body).Decode(&contact)
	logrus.Info(did)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "ERR_PARSE"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	errcode := 400
	if errcode == 401 {
		CommonRes.Code = "ERR_DB_CONNECTION"
		CommonRes.Message = "Data base connection error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
	} else if errcode == 200 {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Contact merged successfully"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
}
func QueryPublicKey(did string, Jwt string) string {
	url := documents.DLAccessUrl2 + "/api/v1/identities/" + did
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Error(err.Error())
	}
	req.Header.Add("Authorization", "Bearer "+Jwt)
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
func addContact(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	authorization := r.Header.Get("Authorization")
	queryValues := r.URL.Query()
	did := queryValues.Get("did")

	contact := model.QuickContactReq{}
	err = json.NewDecoder(r.Body).Decode(&contact)
	if err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "400"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

	errcode, contacts := Addquickcontact(contact, authorization, did)
	if errcode == 500 {
		CommonRes.Code = "500"
		CommonRes.Message = "Failed due to server error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(CommonRes)
	} else if errcode == 200 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contacts)
		return
	} else if errcode == 201 {
		CommonRes.Code = "SUCCESS"
		CommonRes.Message = "Invitaion sent"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else if errcode == 400 {
		CommonRes.Code = "400"
		CommonRes.Message = "Input param error."
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	} else if errcode == 401 {
		CommonRes.Code = "401"
		CommonRes.Message = "Authorization error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}

}
func getaccountfeaturedetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	accountDid := params["account"]
	if accountDid == "" || accountDid == "undefined" {
		logrus.Error("accountDid cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "accountDid cannot be null")
		return
	}
	featurecode := params["feature"]
	if featurecode == "" || featurecode == "undefined" {
		logrus.Error("featurecode cannot be null")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "featurecode cannot be null")
		return
	}
	AccountFeature := model.AccountFeaturesTable{}
	if err := db.Model(&model.AccountFeaturesTable{}).Where("account_did = ? AND paramcode=?", accountDid, featurecode).Scan(&AccountFeature).Error; err != nil {
		logrus.Error(err.Error())
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"Error": "Specified Account/paramcode not found"})
		return
	}
	var List []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Feature_Type").Scan(&List)
	Paramlist := make(map[int]string)
	for _, v := range List {
		Paramlist[v.ParamCode] = v.ParamName
	}
	AccountFeatureRes := model.AccountFeaturesTableRes{}
	AccountFeatureRes.AssociatedAt = AccountFeature.UpdatedAt
	AccountFeatureRes.AccountDid = AccountFeature.AccountDid
	AccountFeatureRes.Paramcode = AccountFeature.Paramcode
	AccountFeatureRes.Paramname = Paramlist[AccountFeature.Paramcode]
	AccountFeatureRes.AssociatedBy = AccountFeature.AssociatedBy
	AccountFeatureRes.UnitPrice = AccountFeature.UnitPrice
	AccountFeatureRes.EffectiveStartDate = AccountFeature.EffectiveStartDate
	AccountFeatureRes.Status = AccountFeature.Status
	AccountFeatureRes.Recurring = AccountFeature.Recurring
	// for _, v := range AccountFeature {
	// 	AccountFeatureRes = append(AccountFeatureRes, model.AccountFeaturesTableRes{AssociatedAt: v.UpdatedAt, AccountDid: v.AccountDid, Paramcode: v.Paramcode, Paramname: Paramlist[v.Paramcode], AssociatedBy: v.AssociatedBy, UnitPrice: v.UnitPrice, EffectiveStartDate: v.EffectiveStartDate, Status: v.Status, Recurring: v.Recurring})
	// }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(AccountFeatureRes)
	return
}
func updateparamlist(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	var UpdateReq model.UpdateParamMasterTable
	err := json.NewDecoder(r.Body).Decode(&UpdateReq)
	if err != nil {
		CommonRes.Code = "Failed"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	if err := db.Model(&model.ParamMasterTable{}).Where("param_type = ? AND param_code = ?", UpdateReq.ParamType, UpdateReq.ParamCode).Update(&model.ParamMasterTable{UnitPrice: UpdateReq.UnitPrice}).Error; err != nil {
		logrus.Error(err.Error())
		CommonRes.Code = "Failed"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	CommonRes.Code = "Success"
	CommonRes.Message = "Updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(CommonRes)
	return
}
func findMinAndMax(a []int) (max int) {
	max = a[0]
	for _, value := range a {
		if value > max {
			max = value
		}
	}
	return max
}
func addparamlist(w http.ResponseWriter, r *http.Request) {
	var CommonRes model.CommonResponse
	var UpdateReq model.AddParamMasterTable
	err := json.NewDecoder(r.Body).Decode(&UpdateReq)
	if err != nil {
		CommonRes.Code = "Failed"
		CommonRes.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CommonRes)
		return
	}
	var ParamCode []int
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", UpdateReq.ParamType).Pluck("param_code", &ParamCode)
	if len(ParamCode) == 0 {

	} else {
		max := findMinAndMax(ParamCode)
		Feature := model.ParamMasterTable{ParamType: UpdateReq.ParamType, ParamCode: max + 1, ParamName: UpdateReq.ParamName}
		db.Create(&Feature)
	}
	CommonRes.Code = "Success"
	CommonRes.Message = "Added"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(CommonRes)
	return
}
func AddDataToParamMasterTable() {
	count := 0
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Feature_Type").Count(&count)
	if count == 0 {
		Admin := model.ParamMasterTable{ParamType: "Feature_Type", ParamCode: 1011, ParamName: "GiftCard"}
		db.Create(&Admin)
		Member := model.ParamMasterTable{ParamType: "Feature_Type", ParamCode: 1012, ParamName: "EmTrust"}
		db.Create(&Member)
		Staff := model.ParamMasterTable{ParamType: "Feature_Type", ParamCode: 1013, ParamName: "EzSeal"}
		db.Create(&Staff)
		Agent := model.ParamMasterTable{ParamType: "Feature_Type", ParamCode: 1014, ParamName: "Vouchers"}
		db.Create(&Agent)

	}

	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Count(&count)
	if count == 0 {
		Admin := model.ParamMasterTable{ParamType: "Endorsement_Type", ParamCode: 1001, ParamName: "Admin"}
		db.Create(&Admin)
		Member := model.ParamMasterTable{ParamType: "Endorsement_Type", ParamCode: 1002, ParamName: "Member"}
		db.Create(&Member)
		Staff := model.ParamMasterTable{ParamType: "Endorsement_Type", ParamCode: 1003, ParamName: "Staff"}
		db.Create(&Staff)
		Agent := model.ParamMasterTable{ParamType: "Endorsement_Type", ParamCode: 1004, ParamName: "Agent"}
		db.Create(&Agent)
	}
}
func main() {
	//logg.Println("hello world")
	logrus.Info("In main function")
	ServicePort := os.Getenv("SERVICE_PORT")
	URL := os.Getenv("URL")
	BLOCKCHAIN_API_URL = os.Getenv("API_BLOCKCHAIN_URL")
	CRYPTO_API_URL = os.Getenv("API_CRYPTO_URL")
	EVENT_API_URL = os.Getenv("API_EVENT_URL")
	VOUCHER_SERVICE_API_URL = os.Getenv("API_VOUCHER_SERVICE_URL")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	Debugger := os.Getenv("DEBUG")
	if host == "" {
		file, _ := ioutil.ReadFile("cred.json")
		var c Credentials
		json.Unmarshal(file, &c)
		host = c.Host
		port = c.Port
		dbname = c.Dbname
		username = c.User
		password = c.Password
		if err := godotenv.Load(); err != nil {
			logrus.Info("No .env file found")
		}
		ServicePort, _ = os.LookupEnv("PORT")
		URL, _ = os.LookupEnv("URL")
		BLOCKCHAIN_API_URL, _ = os.LookupEnv("API_BLOCKCHAIN_URL")
		CRYPTO_API_URL, _ = os.LookupEnv("API_CRYPTO_URL")
		VOUCHER_SERVICE_API_URL, _ = os.LookupEnv("API_VOUCHER_SERVICE_URL")
		EVENT_API_URL, _ = os.LookupEnv("API_EVENT_URL")
		Debugger, _ = os.LookupEnv("DEBUG")
	}
	documents.DLAccessUrl2 = URL
	documents.BLOCKCHAIN_API_URL = BLOCKCHAIN_API_URL
	documents.CRYPTO_API_URL = CRYPTO_API_URL
	documents.VOUCHER_SERVICE_API_URL = VOUCHER_SERVICE_API_URL
	documents.EVENT_API_URL = EVENT_API_URL
	if strings.ToLower(Debugger) == "on" {
		DebuggerStatus = true
		documents.DebuggerStatus = true
	} else {
		DebuggerStatus = false
		documents.DebuggerStatus = false
	}
	if DebuggerStatus {
		logrus.Info("debugger is On")
	}
	service, err = documents.New(http.DefaultClient)
	db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+username+" dbname="+dbname+" password="+password+" sslmode=disable")
	defer db.Close()
	if err != nil {
		logrus.Error(err.Error())
	}
	logrus.Info("Connection Established")
	if err := db.AutoMigrate(&model.TempAccountsDetail{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.AccountsDetail{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.MembersDetail{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.MembersRequiredParams{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.MembersRole{}).Error; err != nil {
		logrus.Error(err.Error())
		return
	}
	if err := db.AutoMigrate(&model.ParamMasterTable{}).Error; err != nil {
		log.Println(err)
		return
	}
	if err := db.AutoMigrate(&model.InvitationTable{}).Error; err != nil {
		log.Println(err)
		return
	}
	if err := db.AutoMigrate(&model.MemberContact{}).Error; err != nil {
		log.Println(err)
		return
	}
	AddDataToParamMasterTable()
	router := mux.NewRouter().StrictSlash(true)
	if err := db.AutoMigrate(&model.AccountFeaturesTable{}).Error; err != nil {
		log.Println(err)
		return
	}
	hfkashdkas
	router.Methods("POST").Path("/tempaccounts/{InvitationId:[-A-Za-z0-9]+}").HandlerFunc(createtempaccount)
	router.Methods("POST").Path("/accounts/").HandlerFunc(createAccount)
	router.Methods("GET").Path("/tempaccounts/{tempaccount:[-A-Za-z0-9]+}").HandlerFunc(tempaccountaccess)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}").HandlerFunc(getaccountdetails)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/validate").HandlerFunc(validateaccounts)
	router.Methods("PUT").Path("/accounts/{account:[:-A-Za-z0-9]+}").HandlerFunc(setaccountdetails)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/params").HandlerFunc(getenrollparams)
	router.Methods("POST").Path("/accounts/{account:[:-A-Za-z0-9\"]+}/members").HandlerFunc(enrollmembers)
	router.Methods("POST").Path("/membership/{member:[:-A-Za-z0-9]+}/enroll/default").HandlerFunc(enrolltoEmtrust)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/members").HandlerFunc(getallmembers)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/members/{member:[:-A-Za-z0-9]+}").HandlerFunc(getmembersdetails)
	router.Methods("PUT").Path("/accounts/{account:[:-A-Za-z0-9]+}/members/{member:[:-A-Za-z0-9]+}").HandlerFunc(updatemembersdetails)
	router.Methods("GET").Path("/features").HandlerFunc(getparamlist)
	router.Methods("POST").Path("/features").HandlerFunc(addparamlist)
	router.Methods("PUT").Path("/features").HandlerFunc(updateparamlist)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/features").HandlerFunc(getaccountfeatures)
	router.Methods("POST").Path("/accounts/{account:[:-A-Za-z0-9]+}/features").HandlerFunc(addaccountfeatures)
	router.Methods("PUT").Path("/accounts/{account:[:-A-Za-z0-9]+}/features/{feature:[-A-Za-z0-9]+}").HandlerFunc(updateaccountfeatures)
	router.Methods("GET").Path("/accounts/{account:[:-A-Za-z0-9]+}/features/{feature:[-A-Za-z0-9]+}").HandlerFunc(getaccountfeaturedetail)
	router.Methods("POST").Path("/quick-contacts/").HandlerFunc(addContact)
	router.Methods("GET").Path("/quick-contacts/").HandlerFunc(getContact)
	router.Methods("POST").Path("/merge/{did:[:-A-Za-z0-9]+}").HandlerFunc(mergeContact)
	router.Methods("GET").Path("/test").HandlerFunc(gettestdid)
	router.Methods("POST").Path("/test").HandlerFunc(sealfortest)
	log.Fatal(http.ListenAndServe(ServicePort, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
