package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/shopspring/decimal"
)

type apiServer struct {
	// db     *someDatabase

	router *httprouter.Router
}

type StorageType int

const (
	InMemory StorageType = iota
	Mongo
)

type Payment struct {
	Amount               decimal.Decimal    `json:"amount"`
	BeneficiaryParty     BeneficiaryParty   `json:"beneficiary_party"`
	ChargesInformation   ChargesInformation `json:"charges_information"`
	Currency             string             `json:"currency"`
	DebtorParty          DebtorParty        `json:"debtor_party"`
	EndToEndReference    string             `json:"end_to_end_reference"`
	Fx                   Fx                 `json:"fx"`
	NumericReference     string             `json:"numeric_reference"`
	PaymentID            string             `json:"payment_id"`
	PaymentPurpose       string             `json:"payment_purpose"`
	PaymentScheme        string             `json:"payment_scheme"`
	PaymentType          string             `json:"payment_type"`
	ProcessingDate       string             `json:"processing_date"`
	Reference            string             `json:"reference"`
	SchemePaymentSubType string             `json:"scheme_payment_sub_type"`
	SchemePaymentType    string             `json:"scheme_payment_type"`
	SponsorParty         SponsorParty       `json:"sponsor_party"`
}
type BeneficiaryParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}
type SenderCharges struct {
	Amount   decimal.Decimal `json:"amount"`
	Currency string          `json:"currency"`
}
type ChargesInformation struct {
	BearerCode              string          `json:"bearer_code"`
	ReceiverChargesAmount   decimal.Decimal `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string          `json:"receiver_charges_currency"`
	SenderCharges           []SenderCharges `json:"sender_charges"`
}
type DebtorParty struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}
type Fx struct {
	ContractReference string          `json:"contract_reference"`
	ExchangeRate      string          `json:"exchange_rate"`
	OriginalAmount    decimal.Decimal `json:"original_amount"`
	OriginalCurrency  string          `json:"original_currency"`
}
type SponsorParty struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}
