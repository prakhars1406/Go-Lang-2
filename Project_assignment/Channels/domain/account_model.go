package Domain

type AccountsSubsciptionDetail struct {
	AccountBalance      float32
	CurrentSubscription string
}

type Accounts struct {
	Accounts []AccountsDetail
}
type AccountsDetail struct {
	ID               string
	PackId           int
	Balance          float32
	SubscriptionDate string
	RechargeDate     string
	ExtraChannels    []int
	Email            string
	Phone            string
}
