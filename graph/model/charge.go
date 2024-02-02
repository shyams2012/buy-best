package model

type Charge struct {
	Amount       int64  `json:"amount"`
	ReceiptEmail string `json:"receiptMail"`
	ProductName  string `json:"productName"`
	Customer     string `json:"customer"`
}
