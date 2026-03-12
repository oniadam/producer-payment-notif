package models

type NotifPaymentWa struct {
	CustomerName      string  `json:"customerName"`
	AggrNo            string  `json:"aggrNo"`
	IdNo              string  `json:"idNo"`
	WaNo              string  `json:"waNo"`
	CustomerServiceNo string  `json:"customerServiceNo"`
	Senddtm           string  `json:"sendDtm"`
	Sendby            string  `json:"sendby"`
	Templatecode      string  `json:"templatecode"`
	LanguageCode      string  `json:"languageCode"`
	TotalPaid         float64 `json:"totalPaid"`
	TransactionSrc    string  `json:"transactionSrc"`
	Paymentmetodcode  string  `json:"paymentmetodcode"`
	Refno             string  `json:"refno"`
	RefNoWa           string  `json:"refNoWa"`
	Filepath          string  `json:"filepath"`
	Flagreversal      string  `json:"flagreversal"`
	SenderHpNo        string  `json:"senderHpNo"`
	Createdby         string  `json:"createdby"`
	Createddtm        string  `json:"createddtm"`
}
