package models

type NotifPaymentWa struct {
	CustomerName string `json:"customerName"`
	AggrNo       string `json:"aggrNo"`
	// Amount           float64 `json:"amount"`
	WaNo             string  `json:"waNo"`
	Senddtm          string  `json:"sendDtm"`
	Sendby           string  `json:"sendby"`
	Templatecode     string  `json:"templatecode"`
	TotalPaid        float64 `json:"totalPaid"`
	TransactionSrc   string  `json:"transactionSrc"`
	Paymentmetodcode string  `json:"paymentmetodcode"`
	Refno            string  `json:"refno"`
	RefNoWa          string  `json:"refNoWa"`
	Filepath         string  `json:"filepath"`
	Flagreversal     string  `json:"flagreversal"`
	Createdby        string  `json:"createdby"`
	Createddtm       string  `json:"createddtm"`
}
