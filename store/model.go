package store

import (
	"time"

	"github.com/shopspring/decimal"
)

type InfoRegistry struct {
	ID             string
	IDCertificate  uint64
	TxHash         string
	URLCertificate string
	CreatedAt      time.Time
	FromWallet     string
	ToWallet       string
	FormWallet     string
	FormID         string
}

type IssueReport struct {
	ID          string
	UserID      string
	Wallet      string
	TotalIssue  decimal.Decimal
	Report      []FormReport
	Allocations []Allocation
}

type FormReport struct {
	FormID      string
	FormWallet  string
	FormCreated string
	CertID      string
	NftID       int64
	Url         string
	Residues    []Residue
}

type Residue struct {
	Type   string
	Amount decimal.Decimal
}

type Allocation struct {
	TxHash  string
	Percent decimal.Decimal
	Wallet  string
	Total   decimal.Decimal
}
