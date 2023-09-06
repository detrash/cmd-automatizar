package store

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"recy_network/util"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/shopspring/decimal"
)

type DatabaseStore struct {
	PoolDB *pgxpool.Pool
	Config util.Config
}

func NewDatabaseStore(config util.Config) *DatabaseStore {

	user := config.DBUser
	password := config.DBPass
	host := config.DBHost
	dbname := config.DBNAME

	poolDB := CreatePGXPool(user, password, host, "5432", dbname, "certificateDB")

	if poolDB == nil {
		err := errors.New("DatabaseStore:DATABASE NIL")
		log.Println(err)
		os.Exit(1)
	}

	return &DatabaseStore{PoolDB: poolDB, Config: config}
}

func (dbStore *DatabaseStore) GetInfoForRegistry(start string) ([]InfoRegistry, error) {

	layoutISO := "2006-01-02"
	date_start, err_start := time.Parse(layoutISO, start)

	if err_start != nil {
		return nil, err_start
	}

	//'2022-01-01'::date
	strQuery :=
		`select 
		f.id, 
		f."walletAddress", 
		f."formMetadataUrl" 
		from  "User" u  
		inner join "Form" f  on (f."userId"=u.id)  
		left  join "Certificate" c on (f.id=c.formid)
		where   1=1 
		and u."profileType"='RECYCLER'
		and f."createdAt"  > $1 
		and f."isFormAuthorizedByAdmin"=true  
		and f."walletAddress"  is not null 
		and c.id  is null`

	rows, err := dbStore.PoolDB.Query(context.Background(), strQuery, date_start)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	infos := make([]InfoRegistry, 0)
	for rows.Next() {
		info := InfoRegistry{}
		err := rows.Scan(
			&info.FormID,
			&info.FormWallet,
			&info.URLCertificate)

		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil

}

func (dbStore *DatabaseStore) SaveCertificate(info InfoRegistry) error {

	strSQL := `INSERT INTO public."Certificate" (id, txhash, certid, createdat, formid, walletperform,towallet) 
	          VALUES($1, $2 ,  $3 , CURRENT_TIMESTAMP  at time zone 'utc'  , $4 , $5, $6)`

	_, err := dbStore.PoolDB.Exec(context.Background(), strSQL,
		info.ID,
		info.TxHash,
		info.IDCertificate,
		info.FormID,
		info.FromWallet,
		info.ToWallet)

	return err

}

func (dbStore *DatabaseStore) GetInfoIssue() ([]IssueReport, error) {

	strQuery :=
		`select u.id, f."walletAddress", sum(d.amount) as total 
		from "User" u  inner join "Form" f  on (u.id=f."userId")
		inner join "Document" d on (f.id=d."formId")
		inner join "Certificate" c  on (f.id=c.formid)
		left join issue_certificate ic on (c.id=ic.certificateid)
		where 1=1 and ic.certificateid is null  
		group by u.id , f."walletAddress"`

	rows, err := dbStore.PoolDB.Query(context.Background(), strQuery)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	infos := make([]IssueReport, 0)
	for rows.Next() {
		info := IssueReport{}
		err := rows.Scan(
			&info.UserID,
			&info.Wallet,
			&info.TotalIssue)
		if err != nil {
			return nil, err
		}

		info.Allocations = dbStore.GetAllocation(info.TotalIssue)
		info.Report, err = dbStore.GetFormReport(info.UserID, info.Wallet)
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (dbStore *DatabaseStore) SaveIssueReport(issueReport IssueReport) error {

	strSQL := `
	INSERT INTO public."Issue"
	(id, userid, total_issuance, report, wallet, createdat)
	VALUES($1, $2, $3, $4, $5, CURRENT_TIMESTAMP  at time zone 'utc' )`

	strReport, err := json.Marshal(issueReport)
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbStore.PoolDB.Exec(context.Background(), strSQL,
		issueReport.ID,
		issueReport.UserID,
		issueReport.TotalIssue,
		string(strReport),
		issueReport.Wallet)

	if err != nil {
		log.Fatal(err)
	}

	strSQLIssueCertificate := `
		INSERT INTO public.issue_certificate
		(issued, issueid, certificateid)
		VALUES(TRUE, $1 , $2 )`

	var arrReports = issueReport.Report
	for _, r := range arrReports {
		_, err = dbStore.PoolDB.Exec(context.Background(), strSQLIssueCertificate,
			&issueReport.ID,
			&r.CertID)

		if err != nil {
			log.Fatal(err)
		}
	}

	//insert allocation
	strSQLAlocation := `
	INSERT INTO public."Allocation" 
	(txhash, "percent", wallet, issueid, total)
	VALUES( $1 , $2, $3, $4, $5);`

	arrAllocation := issueReport.Allocations
	for _, allocation := range arrAllocation {
		_, err = dbStore.PoolDB.Exec(context.Background(), strSQLAlocation,
			&allocation.TxHash,
			&allocation.Percent,
			&allocation.Wallet,
			&issueReport.ID,
			&allocation.Total)

		if err != nil {
			log.Fatal(err)
		}
	}

	return err

}

func (dbStore *DatabaseStore) GetFormReport(userID string, walletAddr string) ([]FormReport, error) {

	strQuery :=
		`select f.id, f."walletAddress", c.id, c.certid, 
		 f."formMetadataUrl", f."createdAt", d."residueType", d.amount 
		 from "User" u  inner join "Form" f  on (u.id=f."userId")
		 inner join "Document" d on (f.id=d."formId")
		 inner join "Certificate" c  on (f.id=c.formid)
		 where 1=1 
		 and u.id=$1
		 and f."walletAddress"=$2`

	rows, err := dbStore.PoolDB.Query(context.Background(), strQuery, userID, walletAddr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	reports := make(map[string]FormReport, 0)
	//f.id, f."walletAddress", c.id, c.certid, f."formMetadataUrl", d."residueType", d.amount
	for rows.Next() {
		//reports
		var FormID string
		var FormWallet string
		var NftID int64
		var CertID string
		var FormURL string
		var FormCreated time.Time
		var ResidueType string
		var Amount decimal.Decimal

		err := rows.Scan(
			&FormID,
			&FormWallet,
			&CertID,
			&NftID,
			&FormURL,
			&FormCreated,
			&ResidueType,
			&Amount,
		)

		layoutISO := "2006-01-02"
		dateForm := FormCreated.Format(layoutISO)

		val, ok := reports[FormID]
		if ok {
			val.Residues = append(val.Residues, Residue{
				Type:   ResidueType,
				Amount: Amount,
			})
			reports[FormID] = val

		} else {
			reports[FormID] = FormReport{
				FormID:      FormID,
				FormWallet:  FormWallet,
				FormCreated: dateForm,
				CertID:      CertID,
				NftID:       NftID,
				Url:         FormURL,
				Residues: []Residue{
					{
						Type:   ResidueType,
						Amount: Amount,
					}},
			}
		}

		if err != nil {
			return nil, err
		}
	}

	arrReports := make([]FormReport, 0)
	for _, r := range reports {
		arrReports = append(arrReports, r)

	}
	return arrReports, nil
}

func (dbStore *DatabaseStore) GetAllocation(totalIssue decimal.Decimal) []Allocation {

	allocations := make([]Allocation, 4)

	var CEM = decimal.NewFromInt32(100)

	var percent = decimal.NewFromFloat(dbStore.Config.WalletRecyclePerc)
	var total = totalIssue.Mul(percent).Div(CEM)

	allocations[0] = Allocation{
		TxHash:  "",
		Percent: percent,
		Wallet:  "",
		Total:   total,
	}

	percent = decimal.NewFromFloat(dbStore.Config.WalletDTrashPerc)
	total = totalIssue.Mul(percent).Div(CEM)
	allocations[1] = Allocation{
		TxHash:  "",
		Percent: percent,
		Wallet:  dbStore.Config.WalletDTrash,
		Total:   total,
	}

	percent = decimal.NewFromFloat(dbStore.Config.WalletLiquidezPerc)
	total = totalIssue.Mul(percent).Div(CEM)
	allocations[2] = Allocation{
		TxHash:  "",
		Percent: percent,
		Wallet:  dbStore.Config.WalletLiquidez,
		Total:   total,
	}

	percent = decimal.NewFromFloat(dbStore.Config.WalletUsuariosPerc)
	total = totalIssue.Mul(percent).Div(CEM)
	allocations[3] = Allocation{
		TxHash:  "",
		Percent: percent,
		Wallet:  dbStore.Config.WalletUsuarios,
		Total:   total,
	}

	return allocations

}

func (store *DatabaseStore) Close() {
	store.PoolDB.Close()

}
