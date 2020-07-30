package model

type SendEmailReq struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type PublicViewDetailsReq struct {
	Did         string                     `json:"did"`
	Name        string                     `json:"name"`
	Email       string                     `json:"email"`
	Phone       string                     `json:"phone"`
	About       string                     `json:"about"`
	Pic         string                     `json:"pic"`
	Endorsement []PublicViewEndorsementReq `json:"endorsement"`
}
type PublicViewDetails struct {
	Did   string `json:"did" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	About string `json:"about"`
	Pic   string `json:"pic"`
}
type PublicViewEndorsementReq struct {
	AccountDid string `json:"accountDid"`
	Type       string `json:"type"`
	Date       string `json:"date"`
	Desc       string `json:"desc"`
}
type PublicViewEndorsement struct {
	UserDid    string `json:"UserDid"`
	AccountDid string `json:"accountDid"`
	Type       string `json:"type"`
	Date       string `json:"date"`
	Desc       string `json:"desc"`
}
type Credentials struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
}
type CommonResponse struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}
type QR_Result struct {
	Txid string `json:"txid"`
	Data string `json:"data"`
}
type QR_Status struct {
	Status string `json:"status"`
	Jwt    string `json:"jwt"`
	Id     string `json:"id"`
}
type Accountdetailsreq struct {
	Did string `json:"did"`
	Jwt string `json:"jwt"`
}
type Identreponse struct {
	Result string `json:"result"`
}
type Validatereponse struct {
	Valid bool `json:"valid"`
}
type IdentResponseResult struct {
	Id        string `json:"id"`
	PublicKey string `json:"publicKey"`
}
