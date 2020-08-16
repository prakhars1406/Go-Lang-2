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
	"voucher-service/documents"
	"voucher-service/model"

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

type Credentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}

var BLOCKCHAIN_API_URL string = ""
var CRYPTO_API_URL string = ""
var ACCOUNT_SERVICE_API_URL string = ""
var EVENT_API_URL string = ""
var DebuggerStatus bool = false

func addevents(source string, evt string, actor string, text string) {
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	issuedtime := strconv.FormatInt(umillisec, 10)
	eventargs := []string{source, evt, issuedtime, actor, text}
	documents.Addevent(eventargs)
}
func Queryvoucherbyissuer(OwnerDid string, OrganizationDid string) ([]model.QueryTokensByOwnerResponse, error) {
	var data []model.QueryTokensByOwnerResponse
	data2 := []model.QueryTokensByOwnerResponse{}
	resId, res := service.QueryTokensByIssuer(OwnerDid)
	if resId != 200 {
		return nil, errors.New("Failed to get claimed voucher")
	}
	json.Unmarshal([]byte(res), &data)

	var dataa model.Vouchersdata
	// var res1 string
	newIndex := 0
	now := time.Now()
	for i, v := range data {
		json.Unmarshal([]byte(v.Record.Body), &dataa)
		if dataa.OrganizationDid != OrganizationDid {
			continue
		}
		data2 = append(data2, data[i])
		//res1 = ""
		//res1 = strings.ReplaceAll(v.Record.Body, "'", "\"")

		data2[newIndex].Record.Data.ExpiryDate = dataa.ExpiryDate
		data2[newIndex].Record.Data.Description = dataa.Description
		data2[newIndex].Record.Data.Logo = dataa.Logo
		data2[newIndex].Record.Data.Background = dataa.Background
		data2[newIndex].Record.Data.Transferable = dataa.Transferable
		data2[newIndex].Record.Data.OrganizationDid = dataa.OrganizationDid
		data2[newIndex].Record.Data.OrganizationName = dataa.OrganizationName
		data2[newIndex].Record.Data.CreateBy = dataa.CreateBy
		data2[newIndex].Record.Data.Status = dataa.Status
		data2[newIndex].Record.Body = ""

		expiredate, err := strconv.Atoi(data2[newIndex].Record.Data.ExpiryDate)
		if err != nil {
			expiredate = 0
		}
		ts := int64(expiredate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff < 0 {
			data2[newIndex].Record.Expired = true
		} else {
			data2[newIndex].Record.Expired = false
		}
		newIndex++
		//fmt.Println(dataa)
	}
	return data2, nil
}
func Queryvoucherbyowner(OwnerDid string, IsOrgnizationDid bool) ([]model.QueryTokensByOwnerResponse, error) {
	var data []model.QueryTokensByOwnerResponse
	var data2 []model.QueryTokensByOwnerResponse

	resId, res := service.QueryTokensByOwner(OwnerDid)
	if resId != 200 {
		return data2, errors.New("Failed to claim voucher")
	}
	json.Unmarshal([]byte(res), &data)
	data2 = data
	var dataa model.Vouchersdata
	// var res1 strin
	now := time.Now()
	for i, v := range data {
		//res1 = ""
		//res1 = strings.ReplaceAll(v.Record.Body, "'", "\"")

		json.Unmarshal([]byte(v.Record.Body), &dataa)
		data2[i].Record.Data.ExpiryDate = dataa.ExpiryDate
		data2[i].Record.Data.Description = dataa.Description
		data2[i].Record.Data.Logo = dataa.Logo
		data2[i].Record.Data.Background = dataa.Background
		data2[i].Record.Data.Transferable = dataa.Transferable
		data2[i].Record.Data.OrganizationDid = dataa.OrganizationDid
		data2[i].Record.Data.OrganizationName = dataa.OrganizationName
		data2[i].Record.Data.CreateBy = dataa.CreateBy
		data2[i].Record.Data.IssuedAt = dataa.IssuedAt
		data2[i].Record.Data.Status = dataa.Status
		data2[i].Record.Body = ""

		expiredate, err := strconv.Atoi(data2[i].Record.Data.ExpiryDate)
		if err != nil {
			expiredate = 0
		}
		ts := int64(expiredate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff < 0 {
			data2[i].Record.Expired = true
		} else {
			data2[i].Record.Expired = false
		}
		//fmt.Println(dataa)
	}
	return data2, nil
}

/*
 Fetch did document from the chain
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchDidDocumentForUser(did string) (model.IdentData, error) {

	_, identJson := documents.QueryIdentity(did)
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
func prereserverdeals(PrereserveCount int, VoucherUuid string, OrganizationName string) error {
	voucherdata := model.VouchersDataTemplate{}
	i, _ := uuid.FromString(VoucherUuid)
	if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", i).Scan(&voucherdata).Error; err != nil {
		logrus.Error(err)
		return errors.New("Failed to pre reserve voucher")
	}
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	issuedtime := strconv.FormatInt(umillisec, 10)
	for i := 1; i <= PrereserveCount; i++ {
		u1 := uuid.NewV4()
		args := []string{u1.String(), VoucherUuid, voucherdata.VoucherTitle, voucherdata.TypeOfVoucher, "{\\\"ExpiryDate\\\":\\\"" + voucherdata.ExpiryDate + "\\\",\\\"Description\\\":\\\"" + voucherdata.Description + "\\\",\\\"Logo\\\":\\\"" + voucherdata.Logo + "\\\",\\\"Background\\\":\\\"" + voucherdata.Background + "\\\",\\\"Transferable\\\":\\\"" + strconv.FormatBool(voucherdata.Transferable) + "\\\",\\\"OrganizationDid\\\":\\\"" + voucherdata.OrganizationDid + "\\\",\\\"OrganizationName\\\":\\\"" + OrganizationName + "\\\",\\\"CreateBy\\\":\\\"" + voucherdata.Createdby + "\\\",\\\"IssuedAt\\\":\\\"" + issuedtime + "\\\",\\\"Status\\\":\\\"" + "Active" + "\\\"}", "0", voucherdata.OrganizationDid, voucherdata.OrganizationDid}
		res, _ := service.ClaimVouchers(args)
		if res != 200 {
			return errors.New("Failed to pre reserve voucher")
		}
	}
	return nil
}

func createvoucher(w http.ResponseWriter, r *http.Request) {
	var c model.VouchersJsonTemplate
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var UsageNumber int
	if c.UsageNumber == 0 {
		UsageNumber = 0
	} else {
		UsageNumber = c.UsageNumber
	}
	if c.VoucherUuid != "" && c.PrereserveCount != 0 {
		var Voucherresponse model.VoucherClaimResponse
		err = prereserverdeals(c.PrereserveCount, c.VoucherUuid, c.OrganizationName)
		if err != nil {
			Voucherresponse.Message = err.Error()
			Voucherresponse.VoucherUUID = ""
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Voucherresponse)
			return
		}
		Voucherresponse.Message = "Updated"
		Voucherresponse.VoucherUUID = c.VoucherUuid
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Voucherresponse)
		return
	}
	u1 := uuid.NewV4()
	if err := db.Create(&model.VouchersDataTemplate{VoucherId: u1, VoucherTitle: c.VoucherTitle, Createdat: time.Now(), Createdby: c.Did, TypeOfVoucher: c.TypeOfVoucher, EffectiveDate: c.EffectiveDate, ExpiryDate: c.ExpiryDate, Background: c.Background, Logo: c.Logo, Description: c.Description, Transferable: c.Transferable, UsageNumber: UsageNumber, OrganizationDid: c.OrganizationDid, UsableDays: strings.Join(c.UsableDays, ","), Price: c.Price, StartDisplay: c.StartDisplay, EndDisplay: c.EndDisplay, SpecialDate: c.SpecialDate, Maximumclaims: c.Maximumclaims}).Error; err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range c.Tags {
		if err := db.Create(&model.VoucherTagsTemplate{VoucherID: u1, Tag: v}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	for _, v := range c.Agent {
		if err := db.Create(&model.VoucherAgentsTemplate{VoucherID: u1, Agent: v}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	for _, v := range c.Collaboration {
		if err := db.Create(&model.CollaboratingCompanies{VoucherId: u1, OrganizationDid: v.OrganizationDid}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if err := db.Create(&model.CollaboratingCompanies{VoucherId: u1, OrganizationDid: c.OrganizationDid}).Error; err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if c.Redeemability[0] == 1005 {
		if err := db.Create(&model.RedeemabilityTable{VoucherId: u1, Paramcode: 1005}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	for _, v := range c.Redeemability {
		if err := db.Create(&model.RedeemabilityTable{VoucherId: u1, Paramcode: v}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if c.Visibilty[0] == 1005 {
		if err := db.Create(&model.DealsVisibiltyTable{VoucherId: u1, Paramcode: 1005}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	for _, v := range c.Visibilty {
		if err := db.Create(&model.DealsVisibiltyTable{VoucherId: u1, Paramcode: v}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	var Voucherresponse model.VoucherClaimResponse
	Voucherresponse.Message = "Created"
	Voucherresponse.VoucherUUID = u1.String()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Voucherresponse)
}

func getagent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucherid := params["voucherid"]
	// i, err := strconv.ParseUint(voucherid, 10, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	//	var Agents []model.Voucheragents
	var agent []string
	i, err := uuid.FromString(voucherid)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := db.Model(&model.VoucherAgentsTemplate{}).Where("voucher_id = ?", i).Pluck("agent", &agent).Error; err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DataJson, err := json.Marshal(agent)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DataJson)
}
func createagent(w http.ResponseWriter, r *http.Request) {
	var c model.Agents
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if c.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Name should not be null")
		return
	}
	if db.Where("did = ?", c.Did).First(&model.Agents{}).RecordNotFound() {
		if err := db.Create(&model.Agents{Name: c.Name, Did: c.Did, Orgid: c.Orgid}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Did already in used")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Record Added")
}

type Agentlist struct {
	Name  string
	Did   string
	Orgid string
}

func getagentslist(w http.ResponseWriter, r *http.Request) {
	//	var Agents []model.Voucheragents
	var agent []Agentlist
	if err := db.Model(&model.Agents{}).Scan(&agent).Error; err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DataJson, err := json.Marshal(agent)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DataJson)
}
func findassociatedOrgnziation(identity model.IdentData) []model.EndorsementTypeAndId {
	response := []model.EndorsementTypeAndId{}
	for _, v := range identity.DID.Endorsements {
		response = append(response, model.EndorsementTypeAndId{v.CredentialSubject.AssociatedWith.ID, v.CredentialSubject.AssociatedWith.Type})
	}
	return response
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
func uniquedeals(deals []model.SearchDealsResponseOutput) []model.SearchDealsResponseOutput {
	keys := make(map[model.SearchDealsResponseOutput]bool)
	list := []model.SearchDealsResponseOutput{}
	for _, entry := range deals {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func uniquedealsforissuer(deals []model.QueryTokensByOwnerResponse) []model.QueryTokensByOwnerResponse {
	keys := make(map[model.QueryTokensByOwnerResponse]bool)
	list := []model.QueryTokensByOwnerResponse{}
	for _, entry := range deals {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
func getAccountDetails() (map[string]string, map[string]string, map[string]string) {
	var organizationList []string
	db.Model(&model.VouchersDataTemplate{}).Select("organization_did").Pluck("organization_did", &organizationList)
	uniqueList := unique(organizationList)
	AcctDetails := make(map[string]string)
	AcctLat := make(map[string]string)
	AcctLong := make(map[string]string)
	for _, v := range uniqueList {
		err, data, _ := documents.Getaccountdetails(v)
		if err != nil {
			continue
		}
		var details model.ScanAccountsDetail
		err = json.Unmarshal([]byte(data), &details)
		AcctDetails[v] = details.AccountName
		AcctLat[v] = details.Latitude
		AcctLong[v] = details.Longitude
		if DebuggerStatus {
			logrus.Info("searchvouchers:: AcctDid:" + v + " AcctName:" + AcctDetails[v] + " AcctLat:" + AcctLat[v] + " AcctLong:" + AcctLong[v])
		}
	}
	return AcctDetails, AcctLat, AcctLong
}
func FindExpiredOrNot(AcctDetails map[string]string, AcctLat map[string]string, AcctLong map[string]string, Data []model.SearchVouchersDataResponse) []model.SearchVouchersDataResponse {
	Data2 := []model.SearchVouchersDataResponse{}
	now := time.Now()
	for i, _ := range Data {
		expiredate, err := strconv.Atoi(Data[i].ExpiryDate)
		if err != nil {
			expiredate = 0
		}
		ts := int64(expiredate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff < 0 {
			Data[i].Expired = true
			continue
		} else {
			Data[i].Expired = false
			Data[i].OrganizationName = AcctDetails[Data[i].OrganizationDid]
			Lat, _ := strconv.ParseFloat(AcctLat[Data[i].OrganizationDid], 64)
			Lon, _ := strconv.ParseFloat(AcctLong[Data[i].OrganizationDid], 64)
			Data[i].Latitude = Lat
			Data[i].Longitude = Lon
		}
		Data2 = append(Data2, Data[i])
	}
	return Data2
}
func FindExpiredOrNotForIndexResponse(AcctDetails map[string]string, AcctLat map[string]string, AcctLong map[string]string, Data []model.SearchDealsResponseOutput, Start int, cnt int, flag bool) []model.SearchDealsResponseOutput {
	Data2 := []model.SearchDealsResponseOutput{}
	now := time.Now()
	for i, _ := range Data {
		if flag {
			if i >= Start && i <= (Start+cnt) {
			} else {
				continue
			}
		}
		expiredate, err := strconv.Atoi(Data[i].Effective_date)
		if err != nil {
			expiredate = 0
		}
		ts := int64(expiredate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff < 0 {
			Enddate, err := strconv.Atoi(Data[i].ExpiryDate)
			if err != nil {
				Enddate = 0
			}
			ts := int64(Enddate)
			ts = ts * 1000000
			timeFromTS := time.Unix(0, ts)
			diff := timeFromTS.Sub(now)
			if diff > 0 {
				Data[i].Expired = false
				Data[i].Status = "active"
				Data[i].OrganizationName = AcctDetails[Data[i].OrganizationDid]
				Lat, _ := strconv.ParseFloat(AcctLat[Data[i].OrganizationDid], 64)
				Lon, _ := strconv.ParseFloat(AcctLong[Data[i].OrganizationDid], 64)
				Data[i].Latitude = Lat
				Data[i].Longitude = Lon
			} else {
				continue
			}
		} else {
			Data[i].Expired = false
			Data[i].Status = "non_active"
			Data[i].OrganizationName = AcctDetails[Data[i].OrganizationDid]
			Lat, _ := strconv.ParseFloat(AcctLat[Data[i].OrganizationDid], 64)
			Lon, _ := strconv.ParseFloat(AcctLong[Data[i].OrganizationDid], 64)
			Data[i].Latitude = Lat
			Data[i].Longitude = Lon
		}
		Data2 = append(Data2, Data[i])
	}
	return Data2
}

func GetTempIdCount(data []model.QueryTokensByOwnerResponse) map[string]int {
	uniquedata := uniquedealsforissuer(data)
	dict := make(map[string]int)
	for _, v := range uniquedata {
		dict[v.Record.Template] = dict[v.Record.Template] + 1
	}
	return dict
}
func FilterPrereservedVoucehers(data []model.QueryTokensByOwnerResponse) map[string]int {
	dict := make(map[string]int)
	for _, v := range data {
		if v.Record.Data.Status == "Redeemed" || v.Record.Data.Status == "Settled" {
			dict[v.Record.Template] = dict[v.Record.Template] + 1
		}
	}
	return dict
}
func getvoucher(did string, OrganizationDid string) ([]model.VouchersDataTemplate, string) {
	var data []model.QueryTokensByOwnerResponse
	var Data []model.VouchersDataTemplate
	resId, res := service.QueryTokensByIssuer(OrganizationDid)
	if resId != 200 {
		logrus.Error(err)
		return Data, "Failed to get deals"
	}
	json.Unmarshal([]byte(res), &data)
	tempid := GetTempIdCount(data)
	var data2 []model.QueryTokensByOwnerResponse
	data2, err = Queryvoucherbyissuer(OrganizationDid, OrganizationDid)
	if err != nil {
		logrus.Error(err)
		return Data, "Failed to get deals"
	}
	Ownertempid := FilterPrereservedVoucehers(data2)
	if err := db.Model(&model.VouchersDataTemplate{}).Order("created_at desc").Where("status != 'cancelled' AND Createdby = ? AND organization_did=?", did, OrganizationDid).Find(&Data).Error; err != nil {
		logrus.Error(err)
		return Data, err.Error()
	}
	now := time.Now()
	for i, _ := range Data {
		expiredate, err := strconv.Atoi(Data[i].ExpiryDate)
		if err != nil {
			logrus.Error(err)
			expiredate = 0
		}
		ts := int64(expiredate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff < 0 {
			Data[i].Expired = true
		} else {
			Data[i].Expired = false
		}
		Data[i].ClaimCount = tempid[Data[i].VoucherId.String()]
		Data[i].RedeemCount = Ownertempid[Data[i].VoucherId.String()]
	}
	return Data, ""
}
func getvoucherbyOrganization(OrganizationDid string) ([]model.SearchVouchersDataResponse, string) {
	var Data []model.SearchVouchersDataResponse
	Data2 := []model.SearchVouchersDataResponse{}
	if err := db.Model(&model.VouchersDataTemplate{}).Order("createdat desc").Where("organization_did=? AND status=?", OrganizationDid, "active").Scan(&Data).Error; err != nil {
		log.Println(err)
		return Data, err.Error()
	}
	AcctDetails, AcctLat, AcctLong := getAccountDetails()
	Data2 = FindExpiredOrNot(AcctDetails, AcctLat, AcctLong, Data)
	return Data2, ""
}
func getvoucherbyOrganizationDid(OrganizationDid string, Did string) ([]model.SearchDealsResponseOutput, error) {
	Data2 := []model.SearchDealsResponseOutput{}
	var Response model.IdentResponse
	dealsdetails := []model.SearchDealsResponse{}
	dealsdetails2 := []model.SearchDealsResponse{}
	dealsdetailsOutput := []model.SearchDealsResponseOutput{}
	var Organizations []model.EndorsementTypeAndId
	// get identity of user if already exist
	identity, _ := fetchDidDocumentForUser(Did)
	if "" == identity.Other {
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Not Enrolled"
	} else {
		Response.Status = 200
		Response.Enrolled = true
		Response.Message = "Enrolled"
		Organizations = findassociatedOrgnziation(identity)
	}
	var (
		Voucher_id              uuid.UUID
		Voucher_title           string
		Createdby               string
		Type_of_voucher         string
		Effective_date          string
		Expiry_date             string
		Background              string
		Logo                    string
		Description             string
		Transferable            bool
		Status                  string
		Createrorganization_did string
		Usable_days             string
		Price                   float64
		Start_display           string
		End_display             string
		Collaborganization_did  string
		Paramcode               int
	)
	if !Response.Enrolled {
		// sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id WHERE collaborating_companies.visibility = 'non_members' AND status = 'active' ORDER BY vouchers_data_templates.created_at desc"
		// fmt.Printf(sqlStatement)
		// rows, _ := db.Raw(sqlStatement).Rows() // (*sql.Rows, error)
		// defer rows.Close()
		// for rows.Next() {
		// 	rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display)
		// 	if DealDisplayable(Start_display, End_display) {
		// 		dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status})
		// 	} else {
		// 		continue
		// 	}
		// }
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, false, Did, Organizations))
	} else {

		sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display,collaborating_companies.organization_did,deals_visibilty_tables.paramcode FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id JOIN deals_visibilty_tables ON vouchers_data_templates.voucher_id = deals_visibilty_tables.voucher_id WHERE vouchers_data_templates.status='active' AND vouchers_data_templates.organization_did= ? ORDER BY vouchers_data_templates.created_at desc"
		rows, _ := db.Raw(sqlStatement, OrganizationDid).Rows() // (*sql.Rows, error)
		defer rows.Close()
		if rows == nil {
			return []model.SearchDealsResponseOutput{}, nil
		}
		for rows.Next() {
			rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display, &Collaborganization_did, &Paramcode)
			if DealDisplayable(Start_display, End_display) {
				dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Collaborganization_did: Collaborganization_did, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status, Paramcode: Paramcode})
			} else {
				continue
			}
		}
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, true, Did, Organizations))
	}
	AcctDetails, AcctLat, AcctLong := getAccountDetails()
	Data2 = FindExpiredOrNotForIndexResponse(AcctDetails, AcctLat, AcctLong, dealsdetailsOutput, 0, 0, false)
	return Data2, nil
}
func addenrolldetails(data []model.SearchDealsResponse, enrolled bool, Did string, Organizations []model.EndorsementTypeAndId) []model.SearchDealsResponseOutput {
	data2 := []model.SearchDealsResponseOutput{}
	if !enrolled {
		for _, v := range data {
			data2 = append(data2, model.SearchDealsResponseOutput{VoucherId: v.VoucherId, VoucherTitle: v.VoucherTitle, OrganizationDid: v.Createrorganization_did, OrganizationName: v.OrganizationName, Createdby: v.Createdby, TypeOfVoucher: v.TypeOfVoucher, Effective_date: v.Effective_date, ExpiryDate: v.ExpiryDate, Background: v.Background, Logo: v.Logo, Description: v.Description, Latitude: v.Latitude, Longitude: v.Longitude, Transferable: v.Transferable, Status: v.Status, Enrolled: false})
		}
		return data2
	} else {
		organizationLis := []string{}
		for _, v := range Organizations { //find all organization associate with Did as member
			if strings.ToLower(v.Type) == "member" {
				organizationLis = append(organizationLis, v.DID)
			}
		}
		var details []model.ParamMasterTableResponse
		db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
		ParaMap := make(map[int]string)
		for _, v := range details {
			ParaMap[v.ParamCode] = v.ParamName
		}
		for _, v := range data {
			endorsementpresent := containsEndorsement(Organizations, ParaMap[v.Paramcode], v.Collaborganization_did, v.Paramcode)
			if endorsementpresent {
				present := contains(organizationLis, v.Createrorganization_did)
				if present {
					data2 = append(data2, model.SearchDealsResponseOutput{VoucherId: v.VoucherId, VoucherTitle: v.VoucherTitle, OrganizationDid: v.Createrorganization_did, OrganizationName: v.OrganizationName, Createdby: v.Createdby, TypeOfVoucher: v.TypeOfVoucher, Effective_date: v.Effective_date, ExpiryDate: v.ExpiryDate, Background: v.Background, Logo: v.Logo, Description: v.Description, Latitude: v.Latitude, Longitude: v.Longitude, Transferable: v.Transferable, Status: v.Status, Enrolled: true})
				} else {
					data2 = append(data2, model.SearchDealsResponseOutput{VoucherId: v.VoucherId, VoucherTitle: v.VoucherTitle, OrganizationDid: v.Createrorganization_did, OrganizationName: v.OrganizationName, Createdby: v.Createdby, TypeOfVoucher: v.TypeOfVoucher, Effective_date: v.Effective_date, ExpiryDate: v.ExpiryDate, Background: v.Background, Logo: v.Logo, Description: v.Description, Latitude: v.Latitude, Longitude: v.Longitude, Transferable: v.Transferable, Status: v.Status, Enrolled: false})
				}
			}
		}
	}
	return data2
}
func ContainType(Organizations []model.EndorsementTypeAndId, OrganizationDid string, Endorsementtype int) bool {
	var details []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
	ParaMap := make(map[int]string)
	for _, v := range details {
		ParaMap[v.ParamCode] = v.ParamName
	}
	for _, v := range Organizations {
		if v.DID == OrganizationDid && strings.ToLower(v.Type) == strings.ToLower(ParaMap[Endorsementtype]) {
			return true
		}
	}
	return false
}
func DealDisplayable(Start_display string, End_display string) bool {
	now := time.Now()
	Startdate, err := strconv.Atoi(Start_display)
	if err != nil {
		Startdate = 0
	}
	ts := int64(Startdate)
	ts = ts * 1000000
	timeFromTS := time.Unix(0, ts)
	diff := timeFromTS.Sub(now)
	if diff < 0 {
		Enddate, err := strconv.Atoi(End_display)
		if err != nil {
			Enddate = 0
		}
		ts := int64(Enddate)
		ts = ts * 1000000
		timeFromTS := time.Unix(0, ts)
		diff := timeFromTS.Sub(now)
		if diff > 0 {
			return true
		}
	}
	return false
}
func getvoucherbyIndex(Startindex string, Count string, Did string) ([]model.SearchDealsResponseOutput, error) {
	Data2 := []model.SearchDealsResponseOutput{}
	var Response model.IdentResponse
	dealsdetails := []model.SearchDealsResponse{}
	dealsdetails2 := []model.SearchDealsResponse{}
	dealsdetailsOutput := []model.SearchDealsResponseOutput{}
	Start, err := strconv.Atoi(Startindex)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	cnt, err := strconv.Atoi(Count)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var Organizations []model.EndorsementTypeAndId
	// get identity of user if already exist
	identity, err := fetchDidDocumentForUser(Did)
	if "" == identity.Other {
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Not Enrolled"
	} else {
		Response.Status = 200
		Response.Enrolled = true
		Response.Message = "Enrolled"
		Organizations = findassociatedOrgnziation(identity)
	}
	var (
		Voucher_id              uuid.UUID
		Voucher_title           string
		Createdby               string
		Type_of_voucher         string
		Effective_date          string
		Expiry_date             string
		Background              string
		Logo                    string
		Description             string
		Transferable            bool
		Status                  string
		Createrorganization_did string
		Usable_days             string
		Price                   float64
		Start_display           string
		End_display             string
		Collaborganization_did  string
		Paramcode               int
	)
	if !Response.Enrolled {
		// sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id WHERE AND status = 'active' ORDER BY vouchers_data_templates.created_at desc"
		// fmt.Printf(sqlStatement)
		// rows, _ := db.Raw(sqlStatement).Rows() // (*sql.Rows, error)
		// defer rows.Close()
		// for rows.Next() {
		// 	rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display)
		// 	if DealDisplayable(Start_display, End_display) {
		// 		dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status})
		// 	} else {
		// 		continue
		// 	}
		// }
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, false, Did, Organizations))
	} else {

		sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display,collaborating_companies.organization_did,deals_visibilty_tables.paramcode FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id JOIN deals_visibilty_tables ON vouchers_data_templates.voucher_id = deals_visibilty_tables.voucher_id WHERE vouchers_data_templates.status='active' ORDER BY vouchers_data_templates.created_at desc"
		rows, _ := db.Raw(sqlStatement).Rows() // (*sql.Rows, error)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display, &Collaborganization_did, &Paramcode)
			if DealDisplayable(Start_display, End_display) {
				dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Collaborganization_did: Collaborganization_did, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status, Paramcode: Paramcode})
			} else {
				continue
			}
		}
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, true, Did, Organizations))
	}
	AcctDetails, AcctLat, AcctLong := getAccountDetails()
	Data2 = FindExpiredOrNotForIndexResponse(AcctDetails, AcctLat, AcctLong, dealsdetailsOutput, Start-1, cnt-1, true)
	return Data2, nil
}
func getvoucherbystring(Searchstring string, Did string) ([]model.SearchDealsResponseOutput, error) {
	Data2 := []model.SearchDealsResponseOutput{}
	var Response model.IdentResponse
	dealsdetails := []model.SearchDealsResponse{}
	dealsdetails2 := []model.SearchDealsResponse{}
	dealsdetailsOutput := []model.SearchDealsResponseOutput{}
	var Organizations []model.EndorsementTypeAndId
	// get identity of user if already exist
	identity, _ := fetchDidDocumentForUser(Did)
	if "" == identity.Other {
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Not Enrolled"
	} else {
		Response.Status = 200
		Response.Enrolled = true
		Response.Message = "Enrolled"
		Organizations = findassociatedOrgnziation(identity)
	}
	var (
		Voucher_id              uuid.UUID
		Voucher_title           string
		Createdby               string
		Type_of_voucher         string
		Effective_date          string
		Expiry_date             string
		Background              string
		Logo                    string
		Description             string
		Transferable            bool
		Status                  string
		Createrorganization_did string
		Usable_days             string
		Price                   float64
		Start_display           string
		End_display             string
		Collaborganization_did  string
		Paramcode               int
	)
	if !Response.Enrolled {
		// sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id WHERE collaborating_companies.visibility = 'non_members' AND status = 'active' ORDER BY vouchers_data_templates.created_at desc"
		// fmt.Printf(sqlStatement)
		// rows, _ := db.Raw(sqlStatement).Rows() // (*sql.Rows, error)
		// defer rows.Close()
		// for rows.Next() {
		// 	rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display)
		// 	if DealDisplayable(Start_display, End_display) {
		// 		dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status})
		// 	} else {
		// 		continue
		// 	}
		// }
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, false, Did, Organizations))
	} else {

		sqlStatement := "SELECT vouchers_data_templates.voucher_id,vouchers_data_templates.voucher_title,vouchers_data_templates.createdby,vouchers_data_templates.type_of_voucher,vouchers_data_templates.effective_date,vouchers_data_templates.expiry_date,vouchers_data_templates.background,vouchers_data_templates.logo,vouchers_data_templates.description,vouchers_data_templates.transferable,vouchers_data_templates.status,vouchers_data_templates.organization_did,vouchers_data_templates.usable_days,vouchers_data_templates.price,vouchers_data_templates.start_display,vouchers_data_templates.end_display,collaborating_companies.organization_did,deals_visibilty_tables.paramcode FROM vouchers_data_templates JOIN collaborating_companies ON vouchers_data_templates.voucher_id = collaborating_companies.voucher_id JOIN deals_visibilty_tables ON vouchers_data_templates.voucher_id = deals_visibilty_tables.voucher_id WHERE vouchers_data_templates.status='active' AND (LOWER(voucher_title) LIKE ? OR LOWER(description) LIKE ?) ORDER BY vouchers_data_templates.created_at desc"
		rows, _ := db.Raw(sqlStatement, "%"+strings.ToLower(Searchstring)+"%", "%"+strings.ToLower(Searchstring)+"%").Rows() // (*sql.Rows, error)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&Voucher_id, &Voucher_title, &Createdby, &Type_of_voucher, &Effective_date, &Expiry_date, &Background, &Logo, &Description, &Transferable, &Status, &Createrorganization_did, &Usable_days, &Price, &Start_display, &End_display, &Collaborganization_did, &Paramcode)
			if DealDisplayable(Start_display, End_display) {
				dealsdetails = append(dealsdetails, model.SearchDealsResponse{VoucherId: Voucher_id, VoucherTitle: Voucher_title, Createrorganization_did: Createrorganization_did, Createdby: Createdby, TypeOfVoucher: Type_of_voucher, Effective_date: Effective_date, ExpiryDate: Expiry_date, Background: Background, Logo: Logo, Description: Description, Transferable: Transferable, Collaborganization_did: Collaborganization_did, Usable_days: strings.Split(Usable_days, ","), Price: Price, Status: Status, Paramcode: Paramcode})
			} else {
				continue
			}
		}
		dealsdetails2 = dealsdetails
		dealsdetailsOutput = uniquedeals(addenrolldetails(dealsdetails2, true, Did, Organizations))
	}
	AcctDetails, AcctLat, AcctLong := getAccountDetails()
	Data2 = FindExpiredOrNotForIndexResponse(AcctDetails, AcctLat, AcctLong, dealsdetailsOutput, 0, 0, false)
	return Data2, nil
}
func UpdateVoucherUuidData(Data2 model.SearchVouchersDataByUuid, Tags []string, Paramcode []int, CollaboratingCompanies []model.CollaborationCompaniesData, Visibility []int) model.SearchVouchersDataByUuidResponse {
	var Data model.SearchVouchersDataByUuidResponse
	Data.VoucherId = Data2.VoucherId
	Data.VoucherTitle = Data2.VoucherTitle
	Data.Createdat = Data2.Createdat
	Data.Createdby = Data2.Createdby
	Data.OrganizationDid = Data2.OrganizationDid
	Data.TypeOfVoucher = Data2.TypeOfVoucher
	Data.ExpiryDate = Data2.ExpiryDate
	Data.Background = Data2.Background
	Data.Description = Data2.Description
	Data.Latitude = Data2.Latitude
	Data.Longitude = Data2.Longitude
	Data.Transferable = Data2.Transferable
	Data.ForNonMemberUse = Data2.ForNonMemberUse
	Data.UsageNumber = Data2.UsageNumber
	Data.Logo = Data2.Logo
	Data.Tags = Tags
	Data.Status = Data2.Status
	Data.EffectiveDate = Data2.EffectiveDate
	Data.UsableDays = strings.Split(Data2.UsableDays, ",")
	Data.Price = Data2.Price
	Data.StartDisplay = Data2.StartDisplay
	Data.EndDisplay = Data2.EndDisplay
	Data.SpecialDate = Data2.SpecialDate
	Data.PrereserveCount = Data2.PrereserveCount
	Data.Maximumclaims = Data2.Maximumclaims
	if len(Paramcode) == 0 || Paramcode[0] == 1005 {
		Data.Redeemability = []int{1005}
	} else {
		Data.Redeemability = Paramcode
	}
	if len(Visibility) == 0 || Visibility[0] == 1005 {
		Data.Visibility = []int{1005}
	} else {
		Data.Visibility = Visibility
	}
	Data.Collaboration = CollaboratingCompanies
	return Data
}
func getvoucherbyVoucherUuid(voucherId uuid.UUID, OrganizationDid string) (model.SearchVouchersDataByUuidResponse, error) {
	var Data2 model.SearchVouchersDataByUuid
	var Data3 model.SearchVouchersDataByUuidResponse
	var Tags []string
	if err := db.Model(&model.VouchersDataTemplate{}).Order("createdat desc").Where("voucher_id = ?", voucherId).Scan(&Data2).Error; err != nil {
		logrus.Error(err)
		return model.SearchVouchersDataByUuidResponse{}, err
	}
	if err := db.Model(&model.VoucherTagsTemplate{}).Where("voucher_id = ?", voucherId).Pluck("tag", &Tags).Error; err != nil {
		logrus.Error(err)
		return model.SearchVouchersDataByUuidResponse{}, err
	}
	VisibilityParamcode := []int{}
	if err := db.Model(&model.DealsVisibiltyTable{}).Where("voucher_id = ?", voucherId).Pluck("paramcode", &VisibilityParamcode).Error; err != nil {
		logrus.Error(err)
		return model.SearchVouchersDataByUuidResponse{}, err
	}
	Paramcode := []int{}
	if err := db.Model(&model.RedeemabilityTable{}).Where("voucher_id = ?", voucherId).Pluck("paramcode", &Paramcode).Error; err != nil {
		logrus.Error(err)
		return model.SearchVouchersDataByUuidResponse{}, err
	}
	CollaboratingCompanies := []model.CollaborationCompaniesData{}
	if err := db.Model(&model.CollaboratingCompanies{}).Where("voucher_id = ? AND organization_did !=? ", voucherId, OrganizationDid).Scan(&CollaboratingCompanies).Error; err != nil {
		logrus.Error(err)
		return model.SearchVouchersDataByUuidResponse{}, err
	}
	for i, v := range CollaboratingCompanies {
		err, data, _ := documents.Getaccountdetails(v.OrganizationDid)
		if err != nil {
			return model.SearchVouchersDataByUuidResponse{}, err
		}
		var details model.ScanAccountsDetail
		err = json.Unmarshal([]byte(data), &details)
		CollaboratingCompanies[i].OrganizationName = details.AccountName
	}
	Data3 = UpdateVoucherUuidData(Data2, Tags, Paramcode, CollaboratingCompanies, VisibilityParamcode)
	now := time.Now()
	expiredate, err := strconv.Atoi(Data3.ExpiryDate)
	if err != nil {
		expiredate = 0
	}
	ts := int64(expiredate)
	ts = ts * 1000000
	timeFromTS := time.Unix(0, ts)
	diff := timeFromTS.Sub(now)
	if diff < 0 {
		Data3.Expired = true
	} else {
		Data3.Expired = false
	}
	return Data3, nil
}

/*
 Fetch orgnization list from the chain
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchOrgForUser(did string) ([]string, error) {
	organizationList := []string{}
	status, identJson := documents.QueryIdentity(did)
	if status != 200 {
		return nil, errors.New("QueryIdentity error")
	}
	var identity model.IdentData = model.IdentData{}
	err := json.Unmarshal(identJson, &identity)
	if err != nil {
		return nil, err
	}

	if identity.Other != "" {
		var did model.DIDDocument
		err = json.Unmarshal([]byte(identity.Other), &did)
		identity.DID = did
	}
	if err == nil {
		for _, v := range identity.DID.Endorsements {
			if strings.ToLower(v.CredentialSubject.AssociatedWith.Type) == "member" {
				organizationList = append(organizationList, v.CredentialSubject.AssociatedWith.ID)
			}
		}
		return organizationList, nil
	} else {
		return nil, err
	}
}
func searchvouchers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	SearchString := queryValues.Get("SearchString")
	voucherId := queryValues.Get("VoucherId")
	Startindex := queryValues.Get("Startindex")
	Count := queryValues.Get("Count")
	Did := queryValues.Get("Did")
	Createdby := queryValues.Get("Createdby")
	OrganizationDid := queryValues.Get("OrganizationDid")

	var Data []model.SearchVouchersDataResponse
	Data2 := []model.SearchVouchersDataResponse{}
	if Createdby != "" && OrganizationDid != "" {
		if DebuggerStatus {
			logrus.Info("searchvouchers::Createdby:" + Createdby + " OrganizationDid:" + OrganizationDid)
		}
		Response, err := getvoucher(Createdby, OrganizationDid)
		if err != "" {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response)
		return
	} else if voucherId != "" && OrganizationDid != "" {
		if DebuggerStatus {
			logrus.Info("searchvouchers:: voucherId:" + voucherId)
		}
		//i, err := strconv.ParseUint(voucherId, 10, 64)
		i, _ := uuid.FromString(voucherId)
		res, err := getvoucherbyVoucherUuid(i, OrganizationDid)
		DataJsonn, err := json.Marshal(res)
		if err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(DataJsonn)
		return
	} else if OrganizationDid != "" && Did != "" {
		if DebuggerStatus {
			logrus.Info("searchvouchers::OrganizationDid:" + OrganizationDid + " ,Did::" + Did)
		}
		res1 := []model.SearchDealsResponseOutput{}
		res1, err := getvoucherbyOrganizationDid(OrganizationDid, Did)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res1)
			return
		}
	} else if OrganizationDid != "" && Createdby == "" {
		fmt.Println("searchvouchers::OrganizationDid:" + OrganizationDid)
		Response, err := getvoucherbyOrganization(OrganizationDid)
		if err != "" {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response)
		return
	} else if SearchString != "" && Did != "" {
		if DebuggerStatus {
			logrus.Info("searchvouchers:: SearchString:" + SearchString)
		}
		res1 := []model.SearchDealsResponseOutput{}
		res1, err = getvoucherbystring(SearchString, Did)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res1)
			return
		}
	} else if Startindex != "" && Count != "" && Did != "" {
		if DebuggerStatus {
			logrus.Info("searchvouchers::Startindex:" + Startindex + " Count:" + Count + " Did:" + Did)
		}

		res1 := []model.SearchDealsResponseOutput{}
		res1, err = getvoucherbyIndex(Startindex, Count, Did)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(res1)
			return
		}
	} else {
		if err := db.Model(&model.VouchersDataTemplate{}).Order("createdat desc").Where("status=?", "active").Scan(&Data).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	AcctDetails, AcctLat, AcctLong := getAccountDetails()
	Data2 = FindExpiredOrNot(AcctDetails, AcctLat, AcctLong, Data)
	DataJsonn, err := json.Marshal(Data2)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DataJsonn)
}
func checkendorsementtype(id uuid.UUID, Did string, OrganizationDid string, ident model.IdentData) (bool, error, string) {
	var Paramcode []int
	Paramname := []string{}

	if err := db.Model(&model.RedeemabilityTable{}).Where("voucher_id = ?", id).Pluck("paramcode", &Paramcode).Error; err != nil {
		logrus.Error(err)
		return false, err, ""
	}
	var details []model.ParamMasterTableResponse
	db.Model(&model.ParamMasterTable{}).Where("param_type = ?", "Endorsement_Type").Scan(&details)
	for _, v := range details {
		if containsint(Paramcode, v.ParamCode) {
			Paramname = append(Paramname, strings.ToLower(v.ParamName))
		}
	}
	var Organizations []model.EndorsementTypeAndId
	identity := ident
	Organizations = findassociatedOrgnziation(identity)
	member := false
	for _, v := range Organizations {
		if v.DID == OrganizationDid {
			if "member" == strings.ToLower(v.Type) {
				member = true
			}
			if contains(Paramname, strings.ToLower(v.Type)) {
				return true, nil, ""
			}
		} else {
			continue
		}
	}
	if member {
		return false, nil, "Cannot claim deal,endorsement type not present"
	} else {
		return false, nil, "Cannot claim deal,not enrolled to organization"
	}

}
func validateclaimreq(data model.ClaimVouchersData) string {
	if data.OwnerDid == "" {
		return "OwnerDid cannot be null"
	}
	if data.VoucherId == "" {
		return "VoucherId cannot be null"
	}
	if data.Data.OrganizationDid == "" {
		return "VoucherId cannot be null"
	}
	return ""
}
func CountTempId(data []model.QueryTokensByOwnerResponse, tempid string) int {
	uniquedata := uniquedealsforissuer(data)
	count := 0
	for _, v := range uniquedata {
		if v.Record.Template == tempid {
			count = count + 1
		}
	}
	return count
}
func claimednumber(OrganizationDid string, TempId string) (int, error) {
	var data []model.QueryTokensByOwnerResponse
	resId, res := service.QueryTokensByIssuer(OrganizationDid)
	if resId != 200 {
		logrus.Error("QueryTokensByIssuer failed")
		return 0, errors.New("Failed to count claimed deals")
	}
	json.Unmarshal([]byte(res), &data)
	count := CountTempId(data, TempId)
	return count, nil
}
func claimvouchers(w http.ResponseWriter, r *http.Request) {
	var data model.ClaimVouchersData
	var Data model.ClaimVoucherDbData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	message := validateclaimreq(data)
	if message != "" {
		logrus.Error(err)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	redeemability := false
	id, _ := uuid.FromString(data.VoucherId)
	var details []model.RedeemabilityTable
	db.Model(&model.RedeemabilityTable{}).Where("voucher_id = ?", id).Scan(&details)
	if details[0].Paramcode == 1005 {
		redeemability = false
	} else {
		redeemability = true
	}
	var Claimresponse model.VoucherClaimResponse
	var Response model.IdentResponse
	// get identity of user if already exist
	identity, err := fetchDidDocumentForUser(data.OwnerDid)
	if "" == identity.Other {
		Response.Status = 200
		Response.Enrolled = false
		Response.Message = "Not Enrolled"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response)
		return
	} else {

		if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", data.VoucherId).Scan(&Data).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if Data.Status != "active" {
			w.WriteHeader(http.StatusAlreadyReported)
			Claimresponse.Message = "Voucher not active"
			Claimresponse.VoucherUUID = ""
			DataJson, _ := json.Marshal(Claimresponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(DataJson)
			return
		}
		typesuccess := false
		mesg := ""
		if redeemability {
			typesuccess, err, mesg = checkendorsementtype(id, data.OwnerDid, Data.OrganizationDid, identity)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Failed to claim deals")
				return
			}
			if !typesuccess {
				w.WriteHeader(http.StatusAlreadyReported)
				Claimresponse.Message = mesg
				Claimresponse.VoucherUUID = ""
				DataJson, _ := json.Marshal(Claimresponse)
				w.Header().Set("Content-Type", "application/json")
				w.Write(DataJson)
				return
			}
		}
		if !DealDisplayable(Data.EffectiveDate, Data.ExpiryDate) {
			w.WriteHeader(http.StatusAlreadyReported)
			Claimresponse.Message = "Voucher not active"
			Claimresponse.VoucherUUID = ""
			DataJson, _ := json.Marshal(Claimresponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(DataJson)
			return
		}
		IssueVouhcer, err := Queryvoucherbyowner(data.OwnerDid, false)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		claimcount := 0
		for _, v := range IssueVouhcer {
			if v.Record.Template == data.VoucherId {
				claimcount++
			}
		}
		if claimcount >= Data.UsageNumber {
			w.WriteHeader(http.StatusAlreadyReported)
			Claimresponse.Message = "Voucher claimed maximum number of times"
			Claimresponse.VoucherUUID = ""
			DataJson, _ := json.Marshal(Claimresponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(DataJson)
			return
		}
		totalclaimed := 0
		totalclaimed, err = claimednumber(data.Data.OrganizationDid, data.VoucherId)
		if err != nil {
			logrus.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		Maxclaimcount := 0
		if Data.Maximumclaims == 0 {
			Maxclaimcount = 100
		} else {
			Maxclaimcount = Data.Maximumclaims
		}
		if totalclaimed >= Maxclaimcount {
			w.WriteHeader(http.StatusAlreadyReported)
			Claimresponse.Message = "Voucher claimed maximum number of times"
			Claimresponse.VoucherUUID = ""
			DataJson, _ := json.Marshal(Claimresponse)
			w.Header().Set("Content-Type", "application/json")
			w.Write(DataJson)
			return
		}
	}
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	u1 := uuid.NewV4()
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	issuedtime := strconv.FormatInt(umillisec, 10)
	args := []string{u1.String(), data.VoucherId, data.VoucherName, data.VoucherType, "{\\\"ExpiryDate\\\":\\\"" + data.Data.ExpiryDate + "\\\",\\\"Description\\\":\\\"" + data.Data.Description + "\\\",\\\"Logo\\\":\\\"" + data.Data.Logo + "\\\",\\\"Background\\\":\\\"" + data.Data.Background + "\\\",\\\"Transferable\\\":\\\"" + data.Data.Transferable + "\\\",\\\"OrganizationDid\\\":\\\"" + data.Data.OrganizationDid + "\\\",\\\"OrganizationName\\\":\\\"" + data.Data.OrganizationName + "\\\",\\\"CreateBy\\\":\\\"" + data.IssuerDid + "\\\",\\\"IssuedAt\\\":\\\"" + issuedtime + "\\\",\\\"Status\\\":\\\"" + "Claimed" + "\\\"}", "0", data.Data.OrganizationDid, data.OwnerDid}
	res, _ := service.ClaimVouchers(args)
	if res != 200 {
		w.WriteHeader(http.StatusBadRequest)
		Claimresponse.Message = "Failed to claim voucher"
		Claimresponse.VoucherUUID = ""
		DataJson, _ := json.Marshal(Claimresponse)
		w.Header().Set("Content-Type", "application/json")
		w.Write(DataJson)
		return
	}
	addevents("voucher-service", "DEALS_CLAIM", data.IssuerDid, strings.Join(args, ","))

	Claimresponse.Message = "Voucher claim success"
	Claimresponse.VoucherUUID = u1.String()
	DataJson, err := json.Marshal(Claimresponse)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(DataJson)
}

/*
 Fetch member details from the memberid
 chain responds with 500 if not found.. weired isn't it?
*/
func fetchmemberdetailsfromDid(Did string) (model.DIDDocument, error) {

	_, identJson := documents.QueryMembersdetail(Did)
	var details []model.MemberData
	var did model.DIDDocument
	err := json.Unmarshal(identJson, &details)
	if err != nil {
		return model.DIDDocument{}, err
	}
	if details[0].Record.Other != "" {

		err = json.Unmarshal([]byte(details[0].Record.Other), &did)
	}
	if err != nil {
		return model.DIDDocument{}, err
	}
	return did, nil
}
func Findreedemabilty(organizationDid string, voucheruuid string, OwnerDid string, owneridentity model.IdentData, voucherdetail []model.RedeemVouhcerDetail, vouchertempid string) (bool, error) {
	dt := time.Now()
	date := strings.Split(dt.Format("01-02-2006 15:04:05 Monday"), " ")
	usableDays := []string{}
	for _, v := range voucherdetail {
		usable := false
		if v.SpecialDate == "nill" || v.SpecialDate == "" {
			usableDays = strings.Split(v.UsableDays, ",")
			for _, v := range usableDays {
				if strings.Contains(strings.ToLower(date[2]), strings.ToLower(v)) {
					usable = true
					break
				}
			}
			if !usable {
				return false, errors.New("Redeem is available on " + strings.Join(usableDays, ",") + " only")
			}
		} else {
			i, err := strconv.ParseInt(v.SpecialDate, 10, 64)
			if err != nil {
				logrus.Error(err)
				return false, err
			}
			SpecialDate := strings.Split(time.Unix(0, i*int64(time.Millisecond)).String(), " ")
			TodaysDate := strings.Split(dt.String(), " ")
			if SpecialDate[0] != TodaysDate[0] {
				return false, errors.New("Redeem is available on " + SpecialDate[0] + " only")
			}
		}
	}
	return true, nil
}
func fetchVoucherDetails(vouchertempid uuid.UUID) ([]model.RedeemVouhcerDetail, error) {

	VoucherDetails := []model.RedeemVouhcerDetail{}
	if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", vouchertempid).Scan(&VoucherDetails).Error; err != nil {
		logrus.Error(err)
		return VoucherDetails, err
	}
	return VoucherDetails, nil
}
func fetchVoucherTempid(OwnerDid string, voucheruuid string) (string, error) {
	IssueVouhcer, err := Queryvoucherbyowner(OwnerDid, false)
	if err != nil {
		return "", err
	}
	vouchertempid := ""
	for _, v := range IssueVouhcer {
		if v.Key == voucheruuid {
			vouchertempid = v.Record.Template
			break
		}
	}
	return vouchertempid, nil
}
func UpdatVoucher(voucheruuid string, purpose int, VoucherDetails []model.RedeemVouhcerDetail) bool {
	now := time.Now()
	unixNano := now.UnixNano()
	umillisec := unixNano / 1000000
	issuedtime := strconv.FormatInt(umillisec, 10)
	err, data, _ := documents.Getaccountdetails(VoucherDetails[0].OrganizationDid)
	if err != nil {
		return false
	}
	Status := ""
	if purpose == 100 {
		Status = "Claimed"
	} else if purpose == 101 {
		Status = "Redeemed"
	} else if purpose == 102 {
		Status = "Settled"
	}
	var details model.ScanAccountsDetail
	err = json.Unmarshal([]byte(data), &details)
	args := []string{voucheruuid, "{\\\"ExpiryDate\\\":\\\"" + VoucherDetails[0].ExpiryDate + "\\\",\\\"Description\\\":\\\"" + VoucherDetails[0].Description + "\\\",\\\"Logo\\\":\\\"" + VoucherDetails[0].Logo + "\\\",\\\"Background\\\":\\\"" + VoucherDetails[0].Background + "\\\",\\\"Transferable\\\":\\\"" + strconv.FormatBool(VoucherDetails[0].Transferable) + "\\\",\\\"OrganizationDid\\\":\\\"" + VoucherDetails[0].OrganizationDid + "\\\",\\\"OrganizationName\\\":\\\"" + details.AccountName + "\\\",\\\"CreateBy\\\":\\\"" + VoucherDetails[0].Createdby + "\\\",\\\"IssuedAt\\\":\\\"" + issuedtime + "\\\",\\\"Status\\\":\\\"" + Status + "\\\"}"}
	res := documents.UpdateVoucher(args)
	if res != 200 {
		return false
	}
	return true
}
func CheckReemebility(identity model.IdentData, id uuid.UUID, NewOwnerDid string) (bool, string) {
	redeemability := false
	var details []model.RedeemabilityTable
	db.Model(&model.RedeemabilityTable{}).Where("voucher_id = ?", id).Scan(&details)
	if details[0].Paramcode == 1005 {
		redeemability = false
		return true, ""
	} else {
		redeemability = true
	}
	// get identity of user if already exist
	if "" == identity.Other {
		return false, "Not Enrolled"
	} else {
		var Data model.ClaimVoucherDbData
		if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", id).Scan(&Data).Error; err != nil {
			logrus.Error(err)
			return false, "Cannot transfer voucher"
		}
		if Data.Status != "active" {
			return false, "Voucher not active"
		}
		typesuccess := false
		mesg := ""
		if redeemability {
			typesuccess, err, mesg = checkendorsementtype(id, NewOwnerDid, Data.OrganizationDid, identity)
			if err != nil {
				return false, "Cannot transfer voucher"
			}
			if !typesuccess {
				return false, mesg
			} else {
				return true, ""
			}
		} else {
			return false, "Cannot transfer voucher"
		}
	}
}
func transferVouchers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucheruuid := params["voucher"]
	var Transferresponse model.VoucherTransferResponse
	VoucherDetails := []model.RedeemVouhcerDetail{}
	if voucheruuid == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "voucheruuid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var data model.TransferVoucherdata
	message := ""
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	if data.Purpose == 100 { //check identity of did
		if data.OwnerDid == data.NewOwnerDid {
			Transferresponse.Status = 208
			Transferresponse.Message = "Cannot share voucher to own did"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		// get identity of user if already exist
		message = "Voucher transfer success"
		identity, _ := fetchDidDocumentForUser(data.NewOwnerDid)
		if "" == identity.Other {
			Transferresponse.Status = 204
			Transferresponse.Message = "Not Enrolled"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		vouchertempid, err := fetchVoucherTempid(data.OwnerDid, voucheruuid)
		//get voucher details
		id, _ := uuid.FromString(vouchertempid)
		redeemable, message := CheckReemebility(identity, id, data.NewOwnerDid)
		if !redeemable {
			logrus.Error(err)
			Transferresponse.Status = 208
			Transferresponse.Message = message
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		VoucherDetails, err = fetchVoucherDetails(id)
		if err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
	} else if data.Purpose == 101 {
		if data.OwnerDid == data.NewOwnerDid {
			Transferresponse.Status = 208
			Transferresponse.Message = "Cannot redeem voucher to own did"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		vouchertempid, err := fetchVoucherTempid(data.OwnerDid, voucheruuid)
		//get voucher details
		id, _ := uuid.FromString(vouchertempid)
		VoucherDetails, err = fetchVoucherDetails(id)
		if err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		OrganizationDid := []string{}
		if err := db.Model(&model.CollaboratingCompanies{}).Where("voucher_id = ?", id).Pluck("organization_did", &OrganizationDid).Error; err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		if !contains(OrganizationDid, data.NewOwnerDid) {
			Transferresponse.Status = 204
			Transferresponse.Message = "Cannot Redeem voucher to this organization"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		//check valid organization.
		message = "Voucher redeem success"
		identity, _ := fetchDidDocumentForUser(data.NewOwnerDid)
		if "" == identity.Id {
			Transferresponse.Status = 204
			Transferresponse.Message = "Organization not registered"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		Owneridentity, _ := fetchDidDocumentForUser(data.OwnerDid)
		if "" == Owneridentity.Other {
			Transferresponse.Status = 204
			Transferresponse.Message = "Onwer Not Enrolled"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		redeemable, err := Findreedemabilty(data.NewOwnerDid, voucheruuid, data.OwnerDid, Owneridentity, VoucherDetails, vouchertempid)
		if err != nil {
			Transferresponse.Status = 200
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		if !redeemable {
			Transferresponse.Status = 208
			Transferresponse.Message = "voucher not redeemable"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}

	} else if data.Purpose == 102 {
		// get identity of user if already exist
		message = "Voucher Settlement success"
		identity, _ := fetchDidDocumentForUser(data.NewOwnerDid)
		if "" == identity.Other {
			Transferresponse.Status = 204
			Transferresponse.Message = "Not Enrolled"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		vouchertempid, err := fetchVoucherTempid(data.OwnerDid, voucheruuid)
		//get voucher details
		id, _ := uuid.FromString(vouchertempid)
		VoucherDetails, err = fetchVoucherDetails(id)
		if err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		if data.NewOwnerDid != VoucherDetails[0].OrganizationDid {
			Transferresponse.Status = 208
			Transferresponse.Message = "Cannot transfer voucher to this organization"
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
	} else {
		//not valid purpose of transfer
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "Purpose should be 0 or 1"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	errormesg := ""
	EventMesg := ""
	if data.Purpose == 100 {
		errormesg = "Failed to share voucher"
		EventMesg = "VOUCHER_SHARE"
	} else if data.Purpose == 101 {
		errormesg = "Failed to redeem voucher"
		EventMesg = "VOUCHER_REDEEM"
	} else if data.Purpose == 102 {
		errormesg = "Failed to settle voucher"
		EventMesg = "VOUCHER_SETTLE"
	}
	if !UpdatVoucher(voucheruuid, data.Purpose, VoucherDetails) {
		Transferresponse.Status = 208
		Transferresponse.Message = errormesg
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}

	args := []string{voucheruuid, data.NewOwnerDid}
	res := documents.TransferVoucher(args)
	if res != 200 {
		Transferresponse.Status = res
		Transferresponse.Message = errormesg
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	addevents("voucher-service", EventMesg, data.NewOwnerDid, strings.Join(args, ","))
	Transferresponse.Status = res
	Transferresponse.Message = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transferresponse)
}
func updateVouchers(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucheruuid := params["voucher"]
	var Transferresponse model.VoucherTransferResponse
	if voucheruuid == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "voucheruuid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}

	var data model.UpdateVoucherdata
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	args := []string{data.VouhcerUUID, "{\\\"ExpiryDate\\\":\\\"" + data.NewBody.ExpiryDate + "\\\",\\\"Description\\\":\\\"" + data.NewBody.Description + "\\\",\\\"Logo\\\":\\\"" + data.NewBody.Logo + "\\\",\\\"Background\\\":\\\"" + data.NewBody.Background + "\\\"}"}
	res := documents.UpdateVoucher(args)
	if res != 200 {
		Transferresponse.Status = res
		Transferresponse.Message = "Failed to update voucher"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	addevents("voucher-service", "VOUCHER_UPDATE", data.OwnerDid, strings.Join(args, ","))
	Transferresponse.Status = res
	Transferresponse.Message = "voucher update success"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transferresponse)
}
func getHistory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucheruuid := params["voucher"]
	if voucheruuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "voucheruuid cannot be null")
		return
	}
	resId, res := service.GetHistory(voucheruuid)
	if resId != 200 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Failed to claim voucher")
		return
	}
	var data []model.GetHistoryResponse
	var data2 []model.GetHistoryResponse
	json.Unmarshal([]byte(res), &data)
	data2 = data
	var dataa model.Vouchersdata
	// var res1 string
	for i, v := range data {
		//res1 = ""
		//res1 = strings.ReplaceAll(v.Record.Body, "'", "\"")
		json.Unmarshal([]byte(v.Value.Body), &dataa)
		data2[i].Value.Data.ExpiryDate = dataa.ExpiryDate
		data2[i].Value.Data.Description = dataa.Description
		data2[i].Value.Data.Logo = dataa.Logo
		data2[i].Value.Data.Background = dataa.Background
		data2[i].Value.Data.Transferable = dataa.Transferable
		data2[i].Value.Data.OrganizationDid = dataa.OrganizationDid
		data2[i].Value.Data.OrganizationName = dataa.OrganizationName
		data2[i].Value.Data.CreateBy = dataa.CreateBy
		//	data2[i].Value.Body = ""
		//fmt.Println(dataa)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data2)
}
func getvoucherdetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucheruuid := params["voucher"]
	var Transferresponse model.VoucherTransferResponse
	if voucheruuid == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "voucheruuid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	resId, res := service.GetVoucherDetails(voucheruuid)
	if resId != 200 {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "Failed to get voucher details"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var data model.GetVoucherDetailsResponse
	var data2 model.GetVoucherDetailsResponse
	json.Unmarshal([]byte(res), &data)
	if data.Template == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "Voucher Uuid invalid"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	data2 = data
	var dataa model.Vouchersdata
	json.Unmarshal([]byte(data.Body), &dataa)
	data2.Data.ExpiryDate = dataa.ExpiryDate
	data2.Data.Description = dataa.Description
	data2.Data.Logo = dataa.Logo
	data2.Data.Background = dataa.Background
	data2.Data.Transferable = dataa.Transferable
	data2.Data.OrganizationDid = dataa.OrganizationDid
	data2.Data.OrganizationName = dataa.OrganizationName
	data2.Data.CreateBy = dataa.CreateBy
	data2.Body = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data2)
}
func getclaimvouchers(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	Issuer := queryValues.Get("Issuer")
	Owner := queryValues.Get("Owner")
	OrganizationDid := queryValues.Get("OrganizationDid")
	IsOrganization := queryValues.Get("IsOrganization")
	var Transferresponse model.VoucherTransferResponse
	if Issuer != "" && OrganizationDid != "" {
		IssueVouhcer, err := Queryvoucherbyissuer(Issuer, OrganizationDid)
		if err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(IssueVouhcer)
		return
	} else if Owner != "" {
		IsDidOrganization := false
		if strings.ToLower(IsOrganization) == "true" {
			IsDidOrganization = true
		}
		IssueVouhcer, err := Queryvoucherbyowner(Owner, IsDidOrganization)
		if err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(IssueVouhcer)
		return

	}
}
func getdealsdetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	voucheruuid := params["deals"]
	var Transferresponse model.VoucherTransferResponse
	if voucheruuid == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "voucheruuid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var Data2 model.SearchVouchersDataResponse
	var Data3 model.SearchVouchersDataResponse
	var Tags []string
	i, err := uuid.FromString(voucheruuid)
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}
	if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", i).Scan(&Data2).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	if err := db.Model(&model.VoucherTagsTemplate{}).Where("voucher_id = ?", i).Pluck("tag", &Tags).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	Data3.VoucherId = Data2.VoucherId
	Data3.VoucherTitle = Data2.VoucherTitle
	Data3.Createdat = Data2.Createdat
	Data3.Createdby = Data2.Createdby
	Data3.TypeOfVoucher = Data2.TypeOfVoucher
	Data3.ExpiryDate = Data2.ExpiryDate
	Data3.Background = Data2.Background
	Data3.Description = Data2.Description
	Data3.Latitude = Data2.Latitude
	Data3.Longitude = Data2.Longitude
	Data3.Transferable = Data2.Transferable
	Data3.ForNonMemberUse = Data2.ForNonMemberUse
	Data3.UsageNumber = Data2.UsageNumber
	Data3.Logo = Data2.Logo
	Data3.Tags = Tags
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Data3)
	return
}
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func containsEndorsement(arr []model.EndorsementTypeAndId, paramcode string, OrgDid string, param int) bool {
	for _, a := range arr {
		if param == 1005 {
			return true
		}
		if a.DID == OrgDid && strings.ToLower(a.Type) == strings.ToLower(paramcode) {
			return true
		}
	}
	return false
}
func containsint(arr []int, str int) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
func containsType(arr []model.EndorsementTypeAndId, str model.EndorsementTypeAndId) bool {
	for _, a := range arr {
		if a.DID == str.DID && a.Type == str.Type {
			return true
		}
	}
	return false
}

func setdealsdetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var Transferresponse model.VoucherTransferResponse
	voucheruuid := params["deals"]
	if voucheruuid == "" {
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = "voucheruuid cannot be null"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var data model.VouchersJsonTemplate
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	i, err := uuid.FromString(voucheruuid)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	if data.Status != "" {
		if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", i).Updates(model.VouchersDataTemplate{Status: data.Status}).Error; err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
	} else {
		if err := db.Model(&model.VouchersDataTemplate{}).Where("voucher_id = ?", i).Update(model.VouchersDataTemplate{VoucherTitle: data.VoucherTitle, Createdby: data.Did, TypeOfVoucher: data.TypeOfVoucher, EffectiveDate: data.EffectiveDate, ExpiryDate: data.ExpiryDate, Background: data.Background, Logo: data.Logo, Description: data.Description, Transferable: data.Transferable, UsageNumber: data.UsageNumber, OrganizationDid: data.OrganizationDid, UsableDays: strings.Join(data.UsableDays, ","), Price: data.Price, StartDisplay: data.StartDisplay, EndDisplay: data.EndDisplay, SpecialDate: data.SpecialDate, Maximumclaims: data.Maximumclaims}).Error; err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
		db.Model(&model.VoucherTagsTemplate{}).Where("voucher_id = ?", i).Delete(&model.VoucherTagsTemplate{})
		for _, v := range data.Tags {
			if err := db.Create(&model.VoucherTagsTemplate{VoucherID: i, Tag: v}).Error; err != nil {
				logrus.Error(err)
				Transferresponse.Status = http.StatusBadRequest
				Transferresponse.Message = err.Error()
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(Transferresponse)
				return
			}

		}
		db.Model(&model.VoucherAgentsTemplate{}).Where("voucher_id = ?", i).Delete(&model.VoucherAgentsTemplate{})
		for _, v := range data.Agent {
			if err := db.Create(&model.VoucherAgentsTemplate{VoucherID: i, Agent: v}).Error; err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		db.Model(&model.CollaboratingCompanies{}).Where("voucher_id = ?", i).Delete(&model.CollaboratingCompanies{})
		for _, v := range data.Collaboration {
			if err := db.Create(&model.CollaboratingCompanies{VoucherId: i, OrganizationDid: v.OrganizationDid}).Error; err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if err := db.Create(&model.CollaboratingCompanies{VoucherId: i, OrganizationDid: data.OrganizationDid}).Error; err != nil {
			logrus.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		db.Model(&model.RedeemabilityTable{}).Where("voucher_id = ?", i).Delete(&model.RedeemabilityTable{})
		for _, v := range data.Redeemability {
			if err := db.Create(&model.RedeemabilityTable{VoucherId: i, Paramcode: v}).Error; err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		db.Model(&model.DealsVisibiltyTable{}).Where("voucher_id = ?", i).Delete(&model.DealsVisibiltyTable{})
		for _, v := range data.Visibilty {
			if err := db.Create(&model.DealsVisibiltyTable{VoucherId: i, Paramcode: v}).Error; err != nil {
				logrus.Error(err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
	}

	Transferresponse.Status = http.StatusOK
	Transferresponse.Message = "Updated"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transferresponse)
	return
}
func getreportedissue(w http.ResponseWriter, r *http.Request) {
	var Transferresponse model.VoucherTransferResponse
	Data := []model.Issues{}
	if err := db.Model(&model.Issues{}).Scan(&Data).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	DataJson, err := json.Marshal(Data)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(DataJson)
}
func getcoordinates(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	docid := params["docid"]
	var Transferresponse model.VoucherTransferResponse
	var EsigndataTable model.EsigndataTable
	if err := db.Model(&model.EsigndataTable{}).Where("doc_id = ?", docid).Scan(&EsigndataTable).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	var Coordinatesdata []model.CoordinatesTable
	if err := db.Model(&model.CoordinatesTable{}).Where("uuid = ?", EsigndataTable.UUID).Scan(&Coordinatesdata).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	EsigndataResposne := model.EsigndataResposne{DocId: EsigndataTable.DocId, CreateBy: EsigndataTable.CreateBy, Coordinate: Coordinatesdata}
	DataJson, err := json.Marshal(EsigndataResposne)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(DataJson)
}
func savecoordinates(w http.ResponseWriter, r *http.Request) {
	var c model.Esigndata
	var Transferresponse model.VoucherTransferResponse
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	u1 := uuid.NewV4()
	if err := db.Create(&model.EsigndataTable{DocId: c.DocId, CreateBy: c.CreateBy, UUID: u1}).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	u2 := uuid.NewV4()
	for _, v := range c.Coordinate {
		if err := db.Create(&model.CoordinatesTable{CordinateId: v.Id, UUID: u1, NewUUID: u2, X: v.X, Y: v.Y, RecipientId: v.RecipientId}).Error; err != nil {
			logrus.Error(err)
			Transferresponse.Status = http.StatusBadRequest
			Transferresponse.Message = err.Error()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Transferresponse)
			return
		}
	}
	Transferresponse.Status = http.StatusOK
	Transferresponse.Message = "Coordinates Saved"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transferresponse)
	return
}
func reportissue(w http.ResponseWriter, r *http.Request) {
	var c model.Issues
	var Transferresponse model.VoucherTransferResponse
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	if err := db.Create(&model.Issues{Name: c.Name, Email: c.Email, Issue: c.Issue, Image: c.Image}).Error; err != nil {
		logrus.Error(err)
		Transferresponse.Status = http.StatusBadRequest
		Transferresponse.Message = err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Transferresponse)
		return
	}
	Transferresponse.Status = http.StatusOK
	Transferresponse.Message = "Issue Reported"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Transferresponse)
	return
}
func getparamlist(w http.ResponseWriter, r *http.Request) {
	var Transferresponse model.VoucherTransferResponse
	queryValues := r.URL.Query()
	Paramtype := queryValues.Get("type")
	if Paramtype == "" {
		Transferresponse.Status = http.StatusBadRequest
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
func AddDataToParamMasterTable() {
	count := 0
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
func main() {
	logrus.Info("In main function")
	//ServicePort, _ := os.LookupEnv("PORT")

	ServicePort := os.Getenv("SERVICE_PORT")
	API_BLOCKCHAIN_URL := os.Getenv("API_BLOCKCHAIN_URL")
	API_CRYPTO_URL := os.Getenv("API_CRYPTO_URL")
	API_EVENT_URL := os.Getenv("API_EVENT_URL")
	ACCOUNT_SERVICE_API_URL := os.Getenv("API_ACCOUNT_SERVICE_URL")
	Debugger := os.Getenv("DEBUG")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

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
			logrus.Error("No .env file found")
		}
		ServicePort, _ = os.LookupEnv("PORT")
		API_BLOCKCHAIN_URL, _ = os.LookupEnv("API_BLOCKCHAIN_URL")
		API_CRYPTO_URL, _ = os.LookupEnv("API_CRYPTO_URL")
		ACCOUNT_SERVICE_API_URL, _ = os.LookupEnv("API_ACCOUNT_SERVICE_URL")
		API_EVENT_URL, _ = os.LookupEnv("API_EVENT_URL")
		Debugger, _ = os.LookupEnv("DEBUG")
	}
	BLOCKCHAIN_API_URL = API_BLOCKCHAIN_URL
	CRYPTO_API_URL = API_CRYPTO_URL
	ACCOUNT_SERVICE_API_URL = ACCOUNT_SERVICE_API_URL
	EVENT_API_URL = API_EVENT_URL
	documents.BLOCKCHAIN_API_URL = API_BLOCKCHAIN_URL
	documents.CRYPTO_API_URL = API_CRYPTO_URL
	documents.ACCOUNT_SERVICE_API_URL = ACCOUNT_SERVICE_API_URL
	documents.EVENT_API_URL = API_EVENT_URL
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
		logrus.Error(err)
		return
	}
	logrus.Info("Connection Established")
	if err := db.AutoMigrate(&model.Agents{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.VoucherTagsTemplate{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.VoucherAgentsTemplate{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.VouchersDataTemplate{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.OrganizationDetails{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.MembersDetails{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.CollaboratingCompanies{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.ParamMasterTable{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.DealsVisibiltyTable{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.Issues{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.EsigndataTable{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	if err := db.AutoMigrate(&model.CoordinatesTable{}).Error; err != nil {
		logrus.Error(err)
		return
	}
	AddDataToParamMasterTable()
	if err := db.AutoMigrate(&model.RedeemabilityTable{}).Error; err != nil {
		logrus.Error(err)
		return
	}

	router := mux.NewRouter().StrictSlash(true)
	router.Methods("POST").Path("/deals").HandlerFunc(createvoucher)
	router.Methods("GET").Path("/wai/api/getagent/{voucherid:[-A-Za-z0-9]+}").HandlerFunc(getagent)
	router.Methods("POST").Path("/wai/api/createagent").HandlerFunc(createagent)
	router.Methods("GET").Path("/wai/api/getagentslist").HandlerFunc(getagentslist)
	router.Methods("GET").Path("/deals").HandlerFunc(searchvouchers)
	router.Methods("GET").Path("/deals/{deals:[:-A-Za-z0-9]+}").HandlerFunc(getdealsdetail)
	router.Methods("PUT").Path("/deals/{deals:[:-A-Za-z0-9]+}").HandlerFunc(setdealsdetail)
	router.Methods("POST").Path("/vouchers").HandlerFunc(claimvouchers)
	router.Methods("GET").Path("/vouchers").HandlerFunc(getclaimvouchers)
	router.Methods("POST").Path("/vouchers/{voucher:[:-A-Za-z0-9]+}/transfer").HandlerFunc(transferVouchers)
	router.Methods("PUT").Path("/vouchers/{voucher:[:-A-Za-z0-9]+}").HandlerFunc(updateVouchers)
	router.Methods("GET").Path("/vouchers/{voucher:[:-A-Za-z0-9]+}/history").HandlerFunc(getHistory)
	router.Methods("GET").Path("/vouchers/{voucher:[:-A-Za-z0-9]+}").HandlerFunc(getvoucherdetails)
	router.Methods("GET").Path("/params").HandlerFunc(getparamlist)
	router.Methods("POST").Path("/report-issue").HandlerFunc(reportissue)
	router.Methods("GET").Path("/report-issue").HandlerFunc(getreportedissue)
	router.Methods("POST").Path("/esign/documents").HandlerFunc(savecoordinates)
	router.Methods("GET").Path("/esign/documents/{docid:[:-A-Za-z0-9]+}").HandlerFunc(getcoordinates)
	//log.Fatal(http.ListenAndServe(":8080", router))
	log.Fatal(http.ListenAndServe(ServicePort, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
