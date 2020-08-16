package model

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Transaction struct {
	gorm.Model
	Txid         uuid.UUID `gorm:"primary_key"` //(primary_key=True, default=uuid.uuid4, editable=False)
	Token        uuid.UUID //(default=uuid.uuid4, editable=False)
	Createdat    time.Time
	Status       string `gorm:"type:varchar(50);default:'DRAFT'"`
	Hash         string `gorm:"type:varchar(128);default:'0000'"`
	Action       string `gorm:"type:varchar(50)"`
	FromUser     string `gorm:"type:varchar(128);default:'0000'"`
	ToUser       string `gorm:"type:varchar(128);default:'0000'"`
	Remarks      string
	Authtokentbr string
	Quantity     string `gorm:"type:varchar(2)"`
	Jwt          string `gorm:"type:varchar(2048)"` //models.BinaryField(max_length=2048, null=True)
}

type Usernotification struct {
	gorm.Model
	Userid    string `gorm:"primary_key;type:varchar(1024)"`
	Deviceid  string `gorm:"type:varchar(1024);default:'0000'"`
	Token     string `gorm:"type:varchar(1024);default:'0000'"`
	Cnt       int    `gorm:"default:0;AUTO_INCREMENT"`
	Createdat time.Time
}
type Devicedata struct {
	Brand         string
	BuildNumber   int
	Carrier       string
	DeviceCountry string
	DeviceId      string
	DeviceLocale  string
	Manufacturer  string
	Timezone      string
	UniqueId      string
}
type Callbackdata struct {
	Device    Devicedata
	DeviceId  string
	Email     string
	Id        string
	Name      string
	PublicKey string
	Token     string
	Txid      string
}
type CollaborationData struct {
	OrganizationDid  string
	OrganizationName string
}
type CollaborationCompaniesData struct {
	OrganizationDid  string
	OrganizationName string
}
type ParamMasterTable struct {
	ParamType string
	ParamCode int `gorm:"primary_key"`
	ParamName string
}
type ParamMasterTableResponse struct {
	ParamCode int
	ParamName string
}
type GetparamlistRes struct {
	ParamName []string `json:"params"`
}
type VouchersJsonTemplate struct {
	VoucherUuid      string
	Did              string
	VoucherTitle     string
	TypeOfVoucher    string
	EffectiveDate    string
	ExpiryDate       string
	Agent            []string
	Background       string
	Transferable     bool
	UsageNumber      int
	OrganizationDid  string
	Latitude         float64
	Longitude        float64
	Description      string
	Logo             string
	Tags             []string
	Collaboration    []CollaborationData
	Visibilty        []int //endorsement types
	Redeemability    []int //endorsement types
	UsableDays       []string
	Price            float64
	StartDisplay     string
	EndDisplay       string
	SpecialDate      string
	PrereserveCount  int
	OrganizationName string
	Status           string
	Maximumclaims    int
}
type VouchersDataTemplate struct {
	gorm.Model
	VoucherId        uuid.UUID `gorm:"primary_key"`
	VoucherTitle     string    `gorm:"type:varchar(30)"`
	Createdat        time.Time
	Createdby        string
	TypeOfVoucher    string
	EffectiveDate    string `gorm:"default:'nill'"`
	ExpiryDate       string
	Background       string
	Logo             string
	Description      string
	Transferable     bool
	ForNonMemberUse  bool
	Visibility       string `gorm:"default:'non_members'"` //members/non_members
	UsageNumber      int
	OrganizationDid  string
	OrganizationName string
	Latitude         float64
	Longitude        float64
	Status           string `gorm:"default:'draft'"` //draft,inactive,active,cancelled
	Expired          bool
	UsableDays       string
	Price            float64
	StartDisplay     string `gorm:"default:'nill'"`
	EndDisplay       string `gorm:"default:'nill'"`
	SpecialDate      string `gorm:"default:'nill'"`
	ClaimCount       int    `gorm:"default:0"`
	RedeemCount      int    `gorm:"default:0"`
	Maximumclaims    int
}
type SearchVouchersDataByUuid struct {
	VoucherId        uuid.UUID
	VoucherTitle     string
	Createdat        time.Time
	Createdby        string
	OrganizationDid  string
	OrganizationName string
	TypeOfVoucher    string
	EffectiveDate    string
	ExpiryDate       string
	Background       string
	Logo             string
	Description      string
	Latitude         float64
	Longitude        float64
	Transferable     bool
	ForNonMemberUse  bool
	UsageNumber      int
	Tags             []string
	Status           string
	Expired          bool
	UsableDays       string
	UsableOn         []string
	Price            float64
	StartDisplay     string
	EndDisplay       string
	SpecialDate      string
	PrereserveCount  int
	Redeemability    []int //endorsement types
	Collaboration    []CollaborationCompaniesData
	Maximumclaims    int
}
type SearchVouchersDataByUuidResponse struct {
	VoucherId        uuid.UUID
	VoucherTitle     string
	Createdat        time.Time
	Createdby        string
	OrganizationDid  string
	OrganizationName string
	TypeOfVoucher    string
	EffectiveDate    string
	ExpiryDate       string
	Background       string
	Logo             string
	Description      string
	Latitude         float64
	Longitude        float64
	Transferable     bool
	ForNonMemberUse  bool
	UsageNumber      int
	Tags             []string
	Status           string
	Expired          bool
	UsableDays       []string
	Price            float64
	StartDisplay     string
	EndDisplay       string
	SpecialDate      string
	PrereserveCount  int
	Redeemability    []int //endorsement types
	Visibility       []int //endorsement types
	Collaboration    []CollaborationCompaniesData
	Maximumclaims    int
}

type RedeemVouhcerDetail struct {
	VoucherId       uuid.UUID
	UsableDays      string
	SpecialDate     string
	OrganizationDid string
	ExpiryDate      string
	Background      string
	Logo            string
	Description     string
	Transferable    bool
	Createdby       string
	PrereserveCount int
}
type CollaboratingCompanies struct {
	gorm.Model
	VoucherId       uuid.UUID
	OrganizationDid string
}
type RedeemabilityTable struct {
	gorm.Model
	VoucherId uuid.UUID
	Paramcode int
}
type DealsVisibiltyTable struct {
	gorm.Model
	VoucherId uuid.UUID
	Paramcode int
}
type OwnerVouchers struct {
	Key    string              `json:"Key"`
	Record OwnerVouchersRecord `json:"Record"`
}
type OwnerVouchersRecord struct {
	Template string `json:"Template"`
}
type SearchVouchersDataResponse struct {
	VoucherId        uuid.UUID
	VoucherTitle     string
	Createdat        time.Time
	Createdby        string
	OrganizationDid  string
	OrganizationName string
	TypeOfVoucher    string
	ExpiryDate       string
	Background       string
	Logo             string
	Description      string
	Latitude         float64
	Longitude        float64
	Transferable     bool
	ForNonMemberUse  bool
	UsageNumber      int
	Tags             []string
	Status           string
	Expired          bool
}

type SearchDealsResponse struct {
	VoucherId               uuid.UUID
	VoucherTitle            string
	Createrorganization_did string
	OrganizationName        string
	Createdby               string
	TypeOfVoucher           string
	Effective_date          string
	ExpiryDate              string
	Background              string
	Logo                    string
	Description             string
	Latitude                float64
	Longitude               float64
	Transferable            bool
	Expired                 bool
	Status                  string
	Collaborganization_did  string
	Visibility              string
	Usable_days             []string
	Price                   float64
	Paramcode               int
}
type SearchDealsResponseOutput struct {
	VoucherId        uuid.UUID
	VoucherTitle     string
	OrganizationDid  string
	OrganizationName string
	Createdby        string
	TypeOfVoucher    string
	Effective_date   string
	ExpiryDate       string
	Background       string
	Logo             string
	Description      string
	Latitude         float64
	Longitude        float64
	Transferable     bool
	Expired          bool
	Status           string
	Enrolled         bool
}

type VouchersTemplateData struct {
	VoucherId     uuid.UUID
	VoucherTitle  string
	Createdat     time.Time
	Createdby     string
	TypeOfVoucher string
	ExpiryDate    string
	Background    string
	Description   string
	Logo          string
	Tags          []string
}
type VoucherTagsTemplate struct {
	VoucherID uuid.UUID
	Tag       string
}
type VoucherAgentsTemplate struct {
	VoucherID uuid.UUID
	Agent     string
}
type Agents struct {
	gorm.Model
	Name  string
	Did   string `gorm:"primary_key"`
	Orgid string
}
type Students struct {
	Name string
}
type ClaimedVouchersTemplate struct {
	Token     uuid.UUID `gorm:"primary_key"`
	VoucherId uuid.UUID
	Did       string
	Claimedat string
	Expiresat string
}
type PublicKey struct {
	Id              string
	Type            string
	PublicKeyBase58 string
	Authorizations  []string
}

type Authentication struct {
	Type   string
	PubKey string
}
type MemberOf struct {
	Id   string
	Type string
	Name string
}

type CredentialSubject struct {
	Id     string
	Member MemberOf
}

type Proof struct {
	Type               string
	Created            string
	ProofPurpose       string
	VerificationMethod string
	Jws                string
}

type Endorsements struct {
	Context      []string
	Id           string
	Type         []string
	Issuer       string
	IssuanceDate string
	Credsubject  CredentialSubject
	Proofdata    Proof
}
type EndorsementsData struct {
	Context     string
	Id          string
	Pubkey      []PublicKey
	Authen      []Authentication
	Endorsedata []Endorsements
}
type OrganizationDetails struct {
	gorm.Model
	OrgName       string
	OrgId         string
	Email         string
	Logo          string
	CreatedBy     string
	TransactionId string
	WebsiteUrl    string
	MembershipUrl string
	Latitude      float64
	Longitude     float64
	OrdDid        string
	PubKey        string
	Status        string
}
type OrganizationDetailsResponse struct {
	gorm.Model
	OrgName string
	Email   string
	Logo    string
	OrdDid  string
	PubKey  string
}
type OrganizationDetailRequest struct {
	OrgName       string
	Email         string
	Logo          string
	WebsiteUrl    string
	MembershipUrl string
	Latitude      float64
	Longitude     float64
	Did           string
}
type MembersDetails struct {
	gorm.Model
	Name         string
	MemberDid    string
	MembershipId string
	Email        string
	PhoneNo      string
	OrgDid       string
	Status       string
}
type MembersDetailsRequest struct {
	Name      string
	MemberDid string
	Email     string
	PhoneNo   string
	OrgDid    string
}
type VoucherClaimResponse struct {
	Message     string
	VoucherUUID string
}
type VerifyEmailResponse struct {
	Message    string
	Registered bool
}
type Voucherdata struct {
	ExpiryDate       string
	Description      string
	Logo             string
	Background       string
	Transferable     string
	OrganizationDid  string
	OrganizationName string
}
type ClaimVouchersData struct {
	VoucherId   string
	VoucherName string
	VoucherType string
	Data        Voucherdata
	IssuerDid   string
	OwnerDid    string
}
type ClaimVoucherDbData struct {
	ForNonMemberUse bool
	UsageNumber     int
	OrganizationDid string
	EffectiveDate   string
	ExpiryDate      string
	Status          string
	PrereserveCount int
	Maximumclaims   int
}
type TransferVoucherdata struct {
	Purpose     int //100=share,101=redeem
	OwnerDid    string
	NewOwnerDid string
}
type VoucherTransferResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
type UpdateVoucherdata struct {
	VouhcerUUID string
	OwnerDid    string
	NewBody     Voucherdata
}
type VoucherUpdateResponse struct {
	Message string
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
type EndorsementTypeAndId struct {
	DID  string
	Type string
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
	Created        string        `json:"created"`
	Updated        string        `json:"updated"`
}
type IdentResponse struct {
	Status   int
	Enrolled bool
	Message  string
}
type GenKeyResponse struct {
	Did       string
	PublicKey string
}
type RegisterOrgResponse struct {
	TransactionId string
}
type Vouchersdata struct {
	ExpiryDate       string `json:'ExpiryDate'`
	Description      string `json:'Description'`
	Logo             string `json:'Logo'`
	Background       string `json:'Background'`
	Transferable     string `json:'transferable'`
	OrganizationDid  string `json:'organizationdid'`
	OrganizationName string `json:'OrganizationName'`
	CreateBy         string `json:'CreateBy'`
	IssuedAt         string `json:'IssuedAt'`
	Status           string `json:'Status'`
}
type QueryTokensByOwnerResponse struct {
	Key    string
	Record QueryTokensByOwnerRecord
}
type QueryTokensByOwnerRecord struct {
	Template string
	Body     string
	Issuer   string
	Name     string
	Owner    string
	Type     string
	Uuid     string
	Value    string
	Data     Vouchersdata
	Expired  bool
}
type GetHistoryResponse struct {
	TxId      string
	Value     GetHistoryValue
	Timestamp string
	IsDelete  string
}
type GetHistoryValue struct {
	Uuid     string
	Template string
	Name     string
	Type     string
	Body     string
	Value    string
	Owner    string
	Issuer   string
	Data     Vouchersdata
}
type GetVoucherDetailsResponse struct {
	Template string
	Body     string
	Issuer   string
	Name     string
	Owner    string
	Type     string
	Uuid     string
	Value    string
	Data     Vouchersdata
}
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
}
type ScanAccountsDetail struct {
	AccountName string
	Latitude    string
	Longitude   string
}
type MemberData struct {
	Key    string       `json:"Key"`
	Record MemberRecord `json:"Record"`
}
type MemberRecord struct {
	DeviceInfo string      `json:"deviceInfo"`
	Email      string      `json:"email"`
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Other      string      `json:"other"`
	DID        DIDDocument `json:"didDocument"`
}
type Identreponse struct {
	Result string `json:"result"`
}
type Issues struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Issue string `json:"issue"`
	Image string `json:"image"`
}
type Coordinates struct {
	Id          string `json:"id"`
	X           string `json:"x"`
	Y           string `json:"y"`
	RecipientId string `json:"recipientId"`
}
type Esigndata struct {
	DocId      string        `json:"docId"`
	CreateBy   string        `json:"createBy"`
	Coordinate []Coordinates `json:"coordinates"`
}
type EsigndataTable struct {
	DocId    string `gorm:"primary_key"`
	CreateBy string
	UUID     uuid.UUID
}
type CoordinatesTable struct {
	CordinateId string    `json:"cordinateId"`
	Status      string    `gorm:"default:'draft'";json:"docId"`
	UUID        uuid.UUID `json:"uuid"`
	NewUUID     uuid.UUID `json:"newUuid"`
	X           string    `json:"x"`
	Y           string    `json:"y"`
	RecipientId string    `json:"recipientId"`
}
type EsigndataResposne struct {
	DocId      string             `json:"docId"`
	CreateBy   string             `json:"createBy"`
	Coordinate []CoordinatesTable `json:"coordinates"`
}
