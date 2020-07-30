package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AccountsDetail struct {
	gorm.Model
	AccountID                string `gorm:"primary_key"`
	AccountDid               string
	AccountName              string
	BusinessName             string
	PreferredSite            string
	Address                  string
	Latitude                 string
	Longitude                string
	Logo                     string
	Status                   string
	CreatedBy                string
	Email                    string
	MembershipURL            string
	Enroll                   bool `gorm:"type:boolean"`
	ExternalMember           bool `gorm:"type:boolean"`
	ExternalMembershipUrl    string
	MembershipRequiredFields string
	Systemaccount            string
	Background               string `gorm:"default:'null'"`
}
type GetAccountsDetailResponse struct {
	gorm.Model
	AccountID                string
	AccountDid               string
	AccountName              string
	BusinessName             string
	PreferredSite            string
	Address                  string
	Latitude                 string
	Longitude                string
	Logo                     string
	Status                   string
	CreatedBy                string
	Email                    string
	MembershipURL            string
	Enroll                   bool
	ExternalMember           bool
	ExternalMembershipUrl    string
	MembershipRequiredFields []string
	AccountFeatures          []int
	Background               string
}
type Setaccountsdetailreq struct {
	BusinessName             string
	Address                  string
	Latitude                 string
	Longitude                string
	Logo                     string
	Email                    string
	MembershipURL            string
	ExternalMember           bool
	ExternalMembershipUrl    string
	MembershipRequiredFields string
	Enroll                   bool
	Background               string
}
type AccountFeaturesTable struct {
	gorm.Model
	AccountDid         string
	Paramcode          int
	AssociatedBy       string
	UnitPrice          int    `gorm:"default:0"`
	EffectiveStartDate string `gorm:"default:'nill'"`
	Status             string `gorm:"default:'draft'"` //active//inactive//draft
	Recurring          string `gorm:"default:'draft'"` //true//false//draft
}
type AccountFeaturesTableRes struct {
	AssociatedAt       time.Time `json:"associatedAt"`
	AccountDid         string    `json:"accountDid"`
	Paramcode          int       `json:"paramCode"`
	Paramname          string    `json:"paramName"`
	AssociatedBy       string    `json:"associatedBy"`
	UnitPrice          int       `json:"unitPrice"`
	EffectiveStartDate string    `json:"effectiveStartDate"`
	Status             string    `json:"status"`
	Recurring          string    `json:"recurring"`
}
type UpdateFeaturesTableReq struct {
	UnitPrice          int    `json:"unitPrice"`
	EffectiveStartDate string `json:"effectiveStartDate"`
	Status             string `json:"status"`
	Recurring          string `json:"recurring"`
}
type AccountFeaturesTableReq struct {
	Paramcode []int
	Did       string
}
type AccountsDetailReq struct {
	DID                string
	PublicKey          string
	OtherDetails       string
	DeviceID           string
	TemporaryAccountId string
	BusinessName       string
	PreferredSite      string
	Background         string
}
type TempAccountsDetail struct {
	gorm.Model
	TemporaryId        string
	InvitationId       string
	PermanentAccountID string
	ExpiryDate         string
	Status             string
}
type TempAccountsResponse struct {
	TemporaryAccountId string
	ExpiryDate         string
}
type MembersDetail struct {
	gorm.Model
	MemberID              string `gorm:"primary_key"`
	AccountID             string
	Name                  string
	Managed               string
	ExternalMemberId      string
	ExternalMembershipUrl string
}
type MembersRequiredParams struct {
	MemberID  string
	AccountID string
	Name      string
	Email     string
	Phone     string
	Address   string
	Sex       string
	DOB       string
	Age       string
}
type GetallmembersRes struct {
	MemberID        string
	MemberDID       string
	AccountID       string
	Name            string
	Email           string
	Phone           string
	Address         string
	Sex             string
	DOB             string
	Age             string
	EndorsementType []int
}
type EnrollMembersReq struct {
	MemberDid             string
	PublicKey             string
	Name                  string
	Managed               string
	ExternalMemberId      string
	ExternalMembershipUrl string
	RequiredParams        MembersRequiredParams //0-Name,1-Email,2-Phone,3-Address,4-Sex,5-DOB,6-Age
}
type EnrolltoEmtrustReq struct {
	PublicKey string `json:"publickey"`
}
type EnrolltoEmtrustRes struct {
	IssuerDid   string      `json:"issuer"`
	Diddocument DIDDocument `json:"did"`
}
type MembersRole struct {
	gorm.Model
	MemberID   string
	TypeOfRole string
}
type AccountsRequest struct {
	BusinessName  string
	PreferredSite string
}
type Credentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}
type GenerateKeyPairResponse struct {
	Did         string
	PublicKey   string
	DidDocument DIDDocument `json:"didDocument"`
}
type EndorsementResponse struct {
	TempId string
	Status string
}
type TempaccessResponse struct {
	Created     bool
	Updateddate string
	AccountDid  string
}
type TempaacountResponse struct {
	gorm.Model
	Status             string
	PermanentAccountID string
}
type PermanentacountResponse struct {
	AccountDid string
	CreatedBy  string
}
type CreateaccountResponse struct {
	Created     bool
	Updateddate string
	AccountDid  string
	UserDid     string
}
type GetenrollparamsResposne struct {
	RequiredParam []string
}
type SetaccountdetailsResponse struct {
	Message string
}
type IdentResponse struct {
	Status   int
	Enrolled bool
	Message  string
}
type CommonResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}
type IdentData struct {
	Status     int         `json:"status"`
	DeviceInfo string      `json:"deviceInfo"`
	Email      string      `json:"email"`
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Other      string      `json:"other"`
	PublicKey  string      `json:"publicKey"`
	DID        DIDDocument `json:"didDocument"`
}
type AssociatedWith struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Desc string `json:"desc"`
}
type Endorsement struct {
	Context           []string  `json:"@context"`
	ID                string    `json:"id"`
	Type              []string  `json:"type"`
	Issuer            string    `json:"issuer"`
	IssuanceDate      time.Time `json:"issuanceDate"`
	CredentialSubject struct {
		ID             string         `json:"id"`
		AssociatedWith AssociatedWith `json:"associatedWith"`
	} `json:"credentialSubject"`
	Proof struct {
		Type               string    `json:"type"`
		Created            time.Time `json:"created"`
		ProofPurpose       string    `json:"proofPurpose"`
		VerificationMethod string    `json:"verificationMethod"`
		Jws                string    `json:"jws"`
	} `json:"proof"`
}

type DIDDocument struct {
	Context   string `json:"@context"`
	ID        string `json:"id"`
	PublicKey []struct {
		ID              string        `json:"id"`
		Type            string        `json:"type"`
		Controller      string        `json:"controller"`
		PublicKeyBase58 string        `json:"publicKeyBase58"`
		Authorizations  []interface{} `json:"authorizations"`
	} `json:"publicKey"`
	Authentication []string      `json:"authentication"`
	Endorsements   []Endorsement `json:"endorsements"`
	Service        []Service     `json:"service"`
	Created        string        `json:"created"`
	Updated        string        `json:"updated"`
}
type GetallmembersResponse struct {
	Members []string `json:"members"`
}
type Service struct {
	Type            string `json:"type"`
	Id              string `json:"id"`
	ServiceEndpoint string `json:"serviceEndpoint"`
	Description     string `json:"description"`
}
type GetMembersDetailsRes struct {
	MemberID   string
	AccountID  string
	Name       string
	Email      string
	Phone      string
	Address    string
	Sex        string
	DOB        string
	Age        string
	TypeOfRole string
}
type UpdateMembersDetailsReq struct {
	Name               string
	Email              string
	Phone              string
	Address            string
	Sex                string
	DOB                string
	Age                string
	EndorsementType    []int
	EndorsementTypeAdd bool
}
type MemberData struct {
	Key    string       `json:"Key"`
	Record MemberRecord `json:"Record"`
}
type Registerreponse struct {
	Result string `json:"result"`
}
type Identreponse struct {
	Result string `json:"result"`
}
type MemberRecord struct {
	DeviceInfo string      `json:"deviceInfo"`
	Email      string      `json:"email"`
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Other      string      `json:"other"`
	DID        DIDDocument `json:"didDocument"`
}
type ParamMasterTableResponse struct {
	ParamCode int     `json:"paramCode"`
	ParamName string  `json:"paramName"`
	UnitPrice float64 `json:"unitPrice"`
}
type EnrollCommonResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type ParamMasterTable struct {
	ParamType string
	ParamCode int `gorm:"primary_key"`
	ParamName string
	UnitPrice float64 `gorm:"default:0"`
}
type UpdateParamMasterTable struct {
	ParamType string  `json:"paramType"`
	ParamCode int     `json:"paramCode"`
	UnitPrice float64 `json:"unitPrice"`
}
type AddParamMasterTable struct {
	ParamType string `json:"paramType"`
	ParamName string `json:"paramName"`
}
type AddFeature struct {
	FeatureType        int    `json:"featureType"`
	FeatureTypeAdd     bool   `json:"featureTypeAdd"`
	Did                string `json:"did"`
	UnitPrice          int    `json:"unitPrice"`
	EffectiveStartDate string `json:"effectiveStartDate"`
	Status             string `json:"status"`
	Recurring          string `json:"recurring"`
}
type EventResult struct {
	Evt    string
	Result string
}
type RegisterDidDocumentResult struct {
	Success bool
	Message string
	Issuer  string
	Roles   string
}
type QuickContactReq struct {
	Contact_did string `json:"contact_did"`
	Name        string `json:"name"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	Remove      bool   `json:"remove"`
}
type MemberContact struct {
	ContactId    string    `gorm:"type:varchar(37);primary_key;not null" json:"contactId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	MemberId     string    `gorm:"type:varchar(37);not null" json:"memberId"`
	ContactDid   string    `gorm:"type:varchar(37);not null" json:"contactDid"`
	Name         string    `gorm:"type:varchar" json:"name"`
	NickName     string    `gorm:"type:varchar" json:"nickName"`
	Email        string    `gorm:"type:varchar" json:"email"`
	InvitationId string    `gorm:"type:varchar(37)" json:"invitationId"`
	Active       string    `gorm:"type:varchar(10)" json:"active"`
}
type InvitationTable struct {
	InvitationId        string        `gorm:"type:varchar(37);primary_key;not null" json:"invitationId"`
	MemberContact       MemberContact `gorm:"foreignkey:InvitationId"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
	InvitaionExpiryDate time.Time     `json:"invitaionExpiryDate"`
	InvitaionAccepted   string        `gorm:"type:char(10);default:'rejected'" json:"invitaionAccepted"`
	InvitationUrl       string        `gorm:"type:varchar;not null" json:"invitationUrl"`
}
type QuickContactRes struct {
	Did         string `json:"did"`
	Name        string `json:"name"`
	NickName    string `json:"nickName"`
	Email       string `json:"email"`
	UpdatedDate string `json:"updatedDate"`
}
type SealtestReq struct {
	DID string `json:"did"`
}
type IdentResponseResult struct {
	Id        string `json:"id"`
	PublicKey string `json:"publicKey"`
}
