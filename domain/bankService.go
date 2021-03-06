package domain

import (
	"fmt"
	"sort"
	"time"
)

type BankService struct {
	Repository AccountRepository
	Clock      Clock
}

func (bankService *BankService) Deposit(deposit Deposit) {
	bankService.Repository.Save(Transaction{Amount: deposit.Amount, Date: bankService.Clock.Now()})
}

func (bankService *BankService) Withdrawal(withdrawal Withdrawal) {
	bankService.Repository.Save(Transaction{Amount: -withdrawal.Amount, Date: bankService.Clock.Now()})
}

func (bankService *BankService) Report() string {
	var transactions = bankService.Repository.GetTransactions()
	var reportTransactions = bankService.createReportTransactionsFrom(transactions)
	sort.Sort(sort.Reverse(ByDate(reportTransactions)))
	return bankService.createReportFrom(reportTransactions)
}

func (bankService *BankService) createReportTransactionsFrom(transactions []Transaction) []ReportTransaction {
	var total Amount
	var reportTransactions []ReportTransaction
	for _, transaction := range transactions {
		total += transaction.Amount
		reportTransactions = append(reportTransactions, ReportTransaction{Transaction: transaction, Total: total})
	}
	return reportTransactions
}

func (BankService *BankService) createReportFrom(reportTransactions []ReportTransaction) string {
	var report = "date || transaction || balance\n"
	for _, reportTransaction := range reportTransactions {
		report = report + fmt.Sprintf("%s || %v || %v\n", reportTransaction.Date.Format("2006-01-02"), reportTransaction.Amount, reportTransaction.Total)
	}
	return report
}

type ByDate []ReportTransaction

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

type Amount float32

type Deposit struct {
	Amount
}

type Withdrawal struct {
	Amount
}

type Transaction struct {
	Amount
	Date time.Time
}

type ReportTransaction struct {
	Transaction
	Total Amount
}

type Clock interface {
	Now() time.Time
}
