package alphavantage

import (
	"context"

	"github.com/mkorenkov/alphavantage/formtype"
	"github.com/pkg/errors"
)

// CompanyProfile makes API request and returns parsed response
func CompanyProfile(ctx context.Context, httpClient HTTPClient, apiKey string, symbol string) (CompanyProfileInfo, error) {
	url := buildURL(apiKey, "OVERVIEW", symbol)
	res := CompanyProfileInfo{}
	if err := makeRequest(ctx, httpClient, url, &res); err != nil {
		return res, errors.Wrap(err, "CompanyProfile error")
	}
	return res, nil
}

// BalanceSheets makes API request and returns parsed response
func BalanceSheets(ctx context.Context, httpClient HTTPClient, apiKey string, symbol string) ([]BalanceSheetStatement, error) {
	url := buildURL(apiKey, "BALANCE_SHEET", symbol)
	response := rawBalanceSheetResponse{}
	if err := makeRequest(ctx, httpClient, url, &response); err != nil {
		return nil, errors.Wrap(err, "BalanceSheets error")
	}
	res := make([]BalanceSheetStatement, 0, len(response.AnnualReports)+len(response.QuarterlyReports))
	for _, raw := range response.AnnualReports {
		b, err := fromBalanceSheet(raw, formtype.Form10K)
		if err != nil {
			return nil, errors.Wrap(err, "BalanceSheets parsing error")
		}
		res = append(res, b)
	}
	for _, raw := range response.QuarterlyReports {
		b, err := fromBalanceSheet(raw, formtype.Form10Q)
		if err != nil {
			return nil, errors.Wrap(err, "BalanceSheets parsing error")
		}
		res = append(res, b)
	}
	return res, nil
}

// CashFlows makes API request and returns parsed response
func CashFlows(ctx context.Context, httpClient HTTPClient, apiKey string, symbol string) ([]CashFlowStatement, error) {
	url := buildURL(apiKey, "CASH_FLOW", symbol)
	response := rawCashFlowResponse{}
	if err := makeRequest(ctx, httpClient, url, &response); err != nil {
		return nil, errors.Wrap(err, "CashFlows error")
	}
	res := make([]CashFlowStatement, 0, len(response.AnnualReports)+len(response.QuarterlyReports))
	for _, raw := range response.AnnualReports {
		b, err := fromCashFlow(raw, formtype.Form10K)
		if err != nil {
			return nil, errors.Wrap(err, "CashFlows parsing error")
		}
		res = append(res, b)
	}
	for _, raw := range response.QuarterlyReports {
		b, err := fromCashFlow(raw, formtype.Form10Q)
		if err != nil {
			return nil, errors.Wrap(err, "CashFlows parsing error")
		}
		res = append(res, b)
	}
	return res, nil
}

// IncomeStatements makes API request and returns parsed response
func IncomeStatements(ctx context.Context, httpClient HTTPClient, apiKey string, symbol string) ([]IncomeStatement, error) {
	url := buildURL(apiKey, "INCOME_STATEMENT", symbol)
	response := rawIncomeStatementResponse{}
	if err := makeRequest(ctx, httpClient, url, &response); err != nil {
		return nil, errors.Wrap(err, "IncomeStatements error")
	}
	res := make([]IncomeStatement, 0, len(response.AnnualReports)+len(response.QuarterlyReports))
	for _, raw := range response.AnnualReports {
		b, err := fromIncomeStatement(raw, formtype.Form10K)
		if err != nil {
			return nil, errors.Wrap(err, "IncomeStatements parsing error")
		}
		res = append(res, b)
	}
	for _, raw := range response.QuarterlyReports {
		b, err := fromIncomeStatement(raw, formtype.Form10Q)
		if err != nil {
			return nil, errors.Wrap(err, "IncomeStatements parsing error")
		}
		res = append(res, b)
	}
	return res, nil
}

// CompanyProfileInfo parsed version of CompanyProfileInfo data received from alphavantage
type CompanyProfileInfo struct {
	Symbol                     string `json:"Symbol"`
	AssetType                  string `json:"AssetType"`
	Name                       string `json:"Name"`
	Description                string `json:"Description"`
	Exchange                   string `json:"Exchange"`
	Currency                   string `json:"Currency"`
	Country                    string `json:"Country"`
	Sector                     string `json:"Sector"`
	Industry                   string `json:"Industry"`
	Address                    string `json:"Address"`
	FullTimeEmployees          int64  `json:"FullTimeEmployees,string"`
	FiscalYearEnd              string `json:"FiscalYearEnd"`
	LatestQuarter              Date   `json:"LatestQuarter"`
	MarketCapitalization       int64  `json:"MarketCapitalization,string"`
	EBITDA                     int64  `json:"EBITDA,string"`
	PERatio                    Money  `json:"PERatio"`
	PEGRatio                   Money  `json:"PEGRatio"`
	BookValue                  Money  `json:"BookValue"`
	DividendPerShare           Money  `json:"DividendPerShare"`
	DividendYield              Money  `json:"DividendYield"`
	EPS                        Money  `json:"EPS"`
	RevenuePerShareTTM         Money  `json:"RevenuePerShareTTM"`
	ProfitMargin               Money  `json:"ProfitMargin"`
	OperatingMarginTTM         Money  `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM          Money  `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM          Money  `json:"ReturnOnEquityTTM"`
	RevenueTTM                 int64  `json:"RevenueTTM,string"`
	GrossProfitTTM             int64  `json:"GrossProfitTTM,string"`
	DilutedEPSTTM              Money  `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY Money  `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  Money  `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         Money  `json:"AnalystTargetPrice"`
	TrailingPE                 Money  `json:"TrailingPE"`
	ForwardPE                  Money  `json:"ForwardPE"`
	PriceToSalesRatioTTM       Money  `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           Money  `json:"PriceToBookRatio"`
	EVToRevenue                Money  `json:"EVToRevenue"`
	EVToEBITDA                 Money  `json:"EVToEBITDA"`
	Beta                       Money  `json:"Beta"`
	High52Week                 Money  `json:"52WeekHigh"`
	Low52Week                  Money  `json:"52WeekLow"`
	SMA50                      Money  `json:"50DayMovingAverage"`
	SMA200                     Money  `json:"200DayMovingAverage"`
	SharesOutstanding          int64  `json:"SharesOutstanding,string"`
	SharesFloat                int64  `json:"SharesFloat,string"`
	SharesShort                int64  `json:"SharesShort,string"`
	SharesShortPriorMonth      int64  `json:"SharesShortPriorMonth,string"`
	ShortRatio                 Money  `json:"ShortRatio"`
	ShortPercentOutstanding    Money  `json:"ShortPercentOutstanding"`
	ShortPercentFloat          Money  `json:"ShortPercentFloat"`
	PercentInsiders            Money  `json:"PercentInsiders"`
	PercentInstitutions        Money  `json:"PercentInstitutions"`
	ForwardAnnualDividendRate  Money  `json:"ForwardAnnualDividendRate"`
	ForwardAnnualDividendYield Money  `json:"ForwardAnnualDividendYield"`
	PayoutRatio                Money  `json:"PayoutRatio"`
	DividendDate               Date   `json:"DividendDate"`
	ExDividendDate             Date   `json:"ExDividendDate"`
	LastSplitFactor            string `json:"LastSplitFactor"`
	LastSplitDate              Date   `json:"LastSplitDate"`
}

type rawBalanceSheetResponse struct {
	Symbol           string                `json:"symbol"`
	AnnualReports    []rawBalanceSheetItem `json:"annualReports"`
	QuarterlyReports []rawBalanceSheetItem `json:"quarterlyReports"`
}

type rawIncomeStatementResponse struct {
	Symbol           string                   `json:"symbol"`
	AnnualReports    []rawIncomeStatementItem `json:"annualReports"`
	QuarterlyReports []rawIncomeStatementItem `json:"quarterlyReports"`
}

type rawCashFlowResponse struct {
	Symbol           string            `json:"symbol"`
	AnnualReports    []rawCashFlowItem `json:"annualReports"`
	QuarterlyReports []rawCashFlowItem `json:"quarterlyReports"`
}

type rawBalanceSheetItem struct {
	FiscalDateEnding                string `json:"fiscalDateEnding"`
	ReportedCurrency                string `json:"reportedCurrency"`
	TotalAssets                     string `json:"totalAssets"`
	IntangibleAssets                string `json:"intangibleAssets"`
	EarningAssets                   string `json:"earningAssets"`
	OtherCurrentAssets              string `json:"otherCurrentAssets"`
	TotalLiabilities                string `json:"totalLiabilities"`
	TotalShareholderEquity          string `json:"totalShareholderEquity"`
	DeferredLongTermLiabilities     string `json:"deferredLongTermLiabilities"`
	OtherCurrentLiabilities         string `json:"otherCurrentLiabilities"`
	CommonStock                     string `json:"commonStock"`
	RetainedEarnings                string `json:"retainedEarnings"`
	OtherLiabilities                string `json:"otherLiabilities"`
	Goodwill                        string `json:"goodwill"`
	OtherAssets                     string `json:"otherAssets"`
	Cash                            string `json:"cash"`
	TotalCurrentLiabilities         string `json:"totalCurrentLiabilities"`
	ShortTermDebt                   string `json:"shortTermDebt"`
	CurrentLongTermDebt             string `json:"currentLongTermDebt"`
	OtherShareholderEquity          string `json:"otherShareholderEquity"`
	PropertyPlantEquipment          string `json:"propertyPlantEquipment"`
	TotalCurrentAssets              string `json:"totalCurrentAssets"`
	LongTermInvestments             string `json:"longTermInvestments"`
	NetTangibleAssets               string `json:"netTangibleAssets"`
	ShortTermInvestments            string `json:"shortTermInvestments"`
	NetReceivables                  string `json:"netReceivables"`
	LongTermDebt                    string `json:"longTermDebt"`
	Inventory                       string `json:"inventory"`
	AccountsPayable                 string `json:"accountsPayable"`
	TotalPermanentEquity            string `json:"totalPermanentEquity"`
	AdditionalPaidInCapital         string `json:"additionalPaidInCapital"`
	CommonStockTotalEquity          string `json:"commonStockTotalEquity"`
	PreferredStockTotalEquity       string `json:"preferredStockTotalEquity"`
	RetainedEarningsTotalEquity     string `json:"retainedEarningsTotalEquity"`
	TreasuryStock                   string `json:"treasuryStock"`
	AccumulatedAmortization         string `json:"accumulatedAmortization"`
	OtherNonCurrrentAssets          string `json:"otherNonCurrrentAssets"`
	DeferredLongTermAssetCharges    string `json:"deferredLongTermAssetCharges"`
	TotalNonCurrentAssets           string `json:"totalNonCurrentAssets"`
	CapitalLeaseObligations         string `json:"capitalLeaseObligations"`
	TotalLongTermDebt               string `json:"totalLongTermDebt"`
	OtherNonCurrentLiabilities      string `json:"otherNonCurrentLiabilities"`
	TotalNonCurrentLiabilities      string `json:"totalNonCurrentLiabilities"`
	NegativeGoodwill                string `json:"negativeGoodwill"`
	Warrants                        string `json:"warrants"`
	PreferredStockRedeemable        string `json:"preferredStockRedeemable"`
	CapitalSurplus                  string `json:"capitalSurplus"`
	LiabilitiesAndShareholderEquity string `json:"liabilitiesAndShareholderEquity"`
	CashAndShortTermInvestments     string `json:"cashAndShortTermInvestments"`
	AccumulatedDepreciation         string `json:"accumulatedDepreciation"`
	CommonStockSharesOutstanding    string `json:"commonStockSharesOutstanding"`
}

// BalanceSheetStatement parsed version of BalanceSheet data received from alphavantage
type BalanceSheetStatement struct {
	FormType                        formtype.FormType `json:"formType"`
	FiscalDateEnding                Date              `json:"fiscalDateEnding"`
	ReportedCurrency                string            `json:"reportedCurrency"`
	TotalAssets                     int64             `json:"totalAssets"`
	IntangibleAssets                int64             `json:"intangibleAssets"`
	EarningAssets                   int64             `json:"earningAssets"`
	OtherCurrentAssets              int64             `json:"otherCurrentAssets"`
	TotalLiabilities                int64             `json:"totalLiabilities"`
	TotalShareholderEquity          int64             `json:"totalShareholderEquity"`
	DeferredLongTermLiabilities     int64             `json:"deferredLongTermLiabilities"`
	OtherCurrentLiabilities         int64             `json:"otherCurrentLiabilities"`
	CommonStock                     int64             `json:"commonStock"`
	RetainedEarnings                int64             `json:"retainedEarnings"`
	OtherLiabilities                int64             `json:"otherLiabilities"`
	Goodwill                        int64             `json:"goodwill"`
	OtherAssets                     int64             `json:"otherAssets"`
	Cash                            int64             `json:"cash"`
	TotalCurrentLiabilities         int64             `json:"totalCurrentLiabilities"`
	ShortTermDebt                   int64             `json:"shortTermDebt"`
	CurrentLongTermDebt             int64             `json:"currentLongTermDebt"`
	OtherShareholderEquity          int64             `json:"otherShareholderEquity"`
	PropertyPlantEquipment          int64             `json:"propertyPlantEquipment"`
	TotalCurrentAssets              int64             `json:"totalCurrentAssets"`
	LongTermInvestments             int64             `json:"longTermInvestments"`
	NetTangibleAssets               int64             `json:"netTangibleAssets"`
	ShortTermInvestments            int64             `json:"shortTermInvestments"`
	NetReceivables                  int64             `json:"netReceivables"`
	LongTermDebt                    int64             `json:"longTermDebt"`
	Inventory                       int64             `json:"inventory"`
	AccountsPayable                 int64             `json:"accountsPayable"`
	TotalPermanentEquity            int64             `json:"totalPermanentEquity"`
	AdditionalPaidInCapital         int64             `json:"additionalPaidInCapital"`
	CommonStockTotalEquity          int64             `json:"commonStockTotalEquity"`
	PreferredStockTotalEquity       int64             `json:"preferredStockTotalEquity"`
	RetainedEarningsTotalEquity     int64             `json:"retainedEarningsTotalEquity"`
	TreasuryStock                   int64             `json:"treasuryStock"`
	AccumulatedAmortization         int64             `json:"accumulatedAmortization"`
	OtherNonCurrrentAssets          int64             `json:"otherNonCurrrentAssets"`
	DeferredLongTermAssetCharges    int64             `json:"deferredLongTermAssetCharges"`
	TotalNonCurrentAssets           int64             `json:"totalNonCurrentAssets"`
	CapitalLeaseObligations         int64             `json:"capitalLeaseObligations"`
	TotalLongTermDebt               int64             `json:"totalLongTermDebt"`
	OtherNonCurrentLiabilities      int64             `json:"otherNonCurrentLiabilities"`
	TotalNonCurrentLiabilities      int64             `json:"totalNonCurrentLiabilities"`
	NegativeGoodwill                int64             `json:"negativeGoodwill"`
	Warrants                        int64             `json:"warrants"`
	PreferredStockRedeemable        int64             `json:"preferredStockRedeemable"`
	CapitalSurplus                  int64             `json:"capitalSurplus"`
	LiabilitiesAndShareholderEquity int64             `json:"liabilitiesAndShareholderEquity"`
	CashAndShortTermInvestments     int64             `json:"cashAndShortTermInvestments"`
	AccumulatedDepreciation         int64             `json:"accumulatedDepreciation"`
	CommonStockSharesOutstanding    int64             `json:"commonStockSharesOutstanding"`
}

func fromBalanceSheet(balanceSheet rawBalanceSheetItem, formType formtype.FormType) (res BalanceSheetStatement, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = BalanceSheetStatement{}
			err = r.(error)
		}
	}()

	return BalanceSheetStatement{
		FormType:                        formType,
		FiscalDateEnding:                panicParseDate(balanceSheet.FiscalDateEnding),
		ReportedCurrency:                balanceSheet.ReportedCurrency,
		TotalAssets:                     panicParseInt64ish(balanceSheet.TotalAssets),
		IntangibleAssets:                panicParseInt64ish(balanceSheet.IntangibleAssets),
		EarningAssets:                   panicParseInt64ish(balanceSheet.EarningAssets),
		OtherCurrentAssets:              panicParseInt64ish(balanceSheet.OtherCurrentAssets),
		TotalLiabilities:                panicParseInt64ish(balanceSheet.TotalLiabilities),
		TotalShareholderEquity:          panicParseInt64ish(balanceSheet.TotalShareholderEquity),
		DeferredLongTermLiabilities:     panicParseInt64ish(balanceSheet.DeferredLongTermLiabilities),
		OtherCurrentLiabilities:         panicParseInt64ish(balanceSheet.OtherCurrentLiabilities),
		CommonStock:                     panicParseInt64ish(balanceSheet.CommonStock),
		RetainedEarnings:                panicParseInt64ish(balanceSheet.RetainedEarnings),
		OtherLiabilities:                panicParseInt64ish(balanceSheet.OtherLiabilities),
		Goodwill:                        panicParseInt64ish(balanceSheet.Goodwill),
		OtherAssets:                     panicParseInt64ish(balanceSheet.OtherAssets),
		Cash:                            panicParseInt64ish(balanceSheet.Cash),
		TotalCurrentLiabilities:         panicParseInt64ish(balanceSheet.TotalCurrentLiabilities),
		ShortTermDebt:                   panicParseInt64ish(balanceSheet.ShortTermDebt),
		CurrentLongTermDebt:             panicParseInt64ish(balanceSheet.CurrentLongTermDebt),
		OtherShareholderEquity:          panicParseInt64ish(balanceSheet.OtherShareholderEquity),
		PropertyPlantEquipment:          panicParseInt64ish(balanceSheet.PropertyPlantEquipment),
		TotalCurrentAssets:              panicParseInt64ish(balanceSheet.TotalCurrentAssets),
		LongTermInvestments:             panicParseInt64ish(balanceSheet.LongTermInvestments),
		NetTangibleAssets:               panicParseInt64ish(balanceSheet.NetTangibleAssets),
		ShortTermInvestments:            panicParseInt64ish(balanceSheet.ShortTermInvestments),
		NetReceivables:                  panicParseInt64ish(balanceSheet.NetReceivables),
		LongTermDebt:                    panicParseInt64ish(balanceSheet.LongTermDebt),
		Inventory:                       panicParseInt64ish(balanceSheet.Inventory),
		AccountsPayable:                 panicParseInt64ish(balanceSheet.AccountsPayable),
		TotalPermanentEquity:            panicParseInt64ish(balanceSheet.TotalPermanentEquity),
		AdditionalPaidInCapital:         panicParseInt64ish(balanceSheet.AdditionalPaidInCapital),
		CommonStockTotalEquity:          panicParseInt64ish(balanceSheet.CommonStockTotalEquity),
		PreferredStockTotalEquity:       panicParseInt64ish(balanceSheet.PreferredStockTotalEquity),
		RetainedEarningsTotalEquity:     panicParseInt64ish(balanceSheet.RetainedEarningsTotalEquity),
		TreasuryStock:                   panicParseInt64ish(balanceSheet.TreasuryStock),
		AccumulatedAmortization:         panicParseInt64ish(balanceSheet.AccumulatedAmortization),
		OtherNonCurrrentAssets:          panicParseInt64ish(balanceSheet.OtherNonCurrrentAssets),
		DeferredLongTermAssetCharges:    panicParseInt64ish(balanceSheet.DeferredLongTermAssetCharges),
		TotalNonCurrentAssets:           panicParseInt64ish(balanceSheet.TotalNonCurrentAssets),
		CapitalLeaseObligations:         panicParseInt64ish(balanceSheet.CapitalLeaseObligations),
		TotalLongTermDebt:               panicParseInt64ish(balanceSheet.TotalLongTermDebt),
		OtherNonCurrentLiabilities:      panicParseInt64ish(balanceSheet.OtherNonCurrentLiabilities),
		TotalNonCurrentLiabilities:      panicParseInt64ish(balanceSheet.TotalNonCurrentLiabilities),
		NegativeGoodwill:                panicParseInt64ish(balanceSheet.NegativeGoodwill),
		Warrants:                        panicParseInt64ish(balanceSheet.Warrants),
		PreferredStockRedeemable:        panicParseInt64ish(balanceSheet.PreferredStockRedeemable),
		CapitalSurplus:                  panicParseInt64ish(balanceSheet.CapitalSurplus),
		LiabilitiesAndShareholderEquity: panicParseInt64ish(balanceSheet.LiabilitiesAndShareholderEquity),
		CashAndShortTermInvestments:     panicParseInt64ish(balanceSheet.CashAndShortTermInvestments),
		AccumulatedDepreciation:         panicParseInt64ish(balanceSheet.AccumulatedDepreciation),
		CommonStockSharesOutstanding:    panicParseInt64ish(balanceSheet.CommonStockSharesOutstanding),
	}, nil
}

type rawCashFlowItem struct {
	FiscalDateEnding               string `json:"fiscalDateEnding"`
	ReportedCurrency               string `json:"reportedCurrency"`
	Investments                    string `json:"investments"`
	ChangeInLiabilities            string `json:"changeInLiabilities"`
	CashflowFromInvestment         string `json:"cashflowFromInvestment"`
	OtherCashflowFromInvestment    string `json:"otherCashflowFromInvestment"`
	NetBorrowings                  string `json:"netBorrowings"`
	CashflowFromFinancing          string `json:"cashflowFromFinancing"`
	OtherCashflowFromFinancing     string `json:"otherCashflowFromFinancing"`
	ChangeInOperatingActivities    string `json:"changeInOperatingActivities"`
	NetIncome                      string `json:"netIncome"`
	ChangeInCash                   string `json:"changeInCash"`
	OperatingCashflow              string `json:"operatingCashflow"`
	OtherOperatingCashflow         string `json:"otherOperatingCashflow"`
	Depreciation                   string `json:"depreciation"`
	DividendPayout                 string `json:"dividendPayout"`
	StockSaleAndPurchase           string `json:"stockSaleAndPurchase"`
	ChangeInInventory              string `json:"changeInInventory"`
	ChangeInAccountReceivables     string `json:"changeInAccountReceivables"`
	ChangeInNetIncome              string `json:"changeInNetIncome"`
	CapitalExpenditures            string `json:"capitalExpenditures"`
	ChangeInReceivables            string `json:"changeInReceivables"`
	ChangeInExchangeRate           string `json:"changeInExchangeRate"`
	ChangeInCashAndCashEquivalents string `json:"changeInCashAndCashEquivalents"`
}

// CashFlowStatement parsed version of CashFlow data received from alphavantage
type CashFlowStatement struct {
	FormType                       formtype.FormType `json:"formType"`
	FiscalDateEnding               Date              `json:"fiscalDateEnding"`
	ReportedCurrency               string            `json:"reportedCurrency"`
	Investments                    int64             `json:"investments"`
	ChangeInLiabilities            int64             `json:"changeInLiabilities"`
	CashflowFromInvestment         int64             `json:"cashflowFromInvestment"`
	OtherCashflowFromInvestment    int64             `json:"otherCashflowFromInvestment"`
	NetBorrowings                  int64             `json:"netBorrowings"`
	CashflowFromFinancing          int64             `json:"cashflowFromFinancing"`
	OtherCashflowFromFinancing     int64             `json:"otherCashflowFromFinancing"`
	ChangeInOperatingActivities    int64             `json:"changeInOperatingActivities"`
	NetIncome                      int64             `json:"netIncome"`
	ChangeInCash                   int64             `json:"changeInCash"`
	OperatingCashflow              int64             `json:"operatingCashflow"`
	OtherOperatingCashflow         int64             `json:"otherOperatingCashflow"`
	Depreciation                   int64             `json:"depreciation"`
	DividendPayout                 int64             `json:"dividendPayout"`
	StockSaleAndPurchase           int64             `json:"stockSaleAndPurchase"`
	ChangeInInventory              int64             `json:"changeInInventory"`
	ChangeInAccountReceivables     int64             `json:"changeInAccountReceivables"`
	ChangeInNetIncome              int64             `json:"changeInNetIncome"`
	CapitalExpenditures            int64             `json:"capitalExpenditures"`
	ChangeInReceivables            int64             `json:"changeInReceivables"`
	ChangeInExchangeRate           int64             `json:"changeInExchangeRate"`
	ChangeInCashAndCashEquivalents int64             `json:"changeInCashAndCashEquivalents"`
}

func fromCashFlow(cashFlow rawCashFlowItem, formType formtype.FormType) (res CashFlowStatement, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = CashFlowStatement{}
			err = r.(error)
		}
	}()

	return CashFlowStatement{
		FormType:                       formType,
		FiscalDateEnding:               panicParseDate(cashFlow.FiscalDateEnding),
		ReportedCurrency:               cashFlow.ReportedCurrency,
		Investments:                    panicParseInt64ish(cashFlow.Investments),
		ChangeInLiabilities:            panicParseInt64ish(cashFlow.ChangeInLiabilities),
		CashflowFromInvestment:         panicParseInt64ish(cashFlow.CashflowFromInvestment),
		OtherCashflowFromInvestment:    panicParseInt64ish(cashFlow.OtherCashflowFromInvestment),
		NetBorrowings:                  panicParseInt64ish(cashFlow.NetBorrowings),
		CashflowFromFinancing:          panicParseInt64ish(cashFlow.CashflowFromFinancing),
		OtherCashflowFromFinancing:     panicParseInt64ish(cashFlow.OtherCashflowFromFinancing),
		ChangeInOperatingActivities:    panicParseInt64ish(cashFlow.ChangeInOperatingActivities),
		NetIncome:                      panicParseInt64ish(cashFlow.NetIncome),
		ChangeInCash:                   panicParseInt64ish(cashFlow.ChangeInCash),
		OperatingCashflow:              panicParseInt64ish(cashFlow.OperatingCashflow),
		OtherOperatingCashflow:         panicParseInt64ish(cashFlow.OtherOperatingCashflow),
		Depreciation:                   panicParseInt64ish(cashFlow.Depreciation),
		DividendPayout:                 panicParseInt64ish(cashFlow.DividendPayout),
		StockSaleAndPurchase:           panicParseInt64ish(cashFlow.StockSaleAndPurchase),
		ChangeInInventory:              panicParseInt64ish(cashFlow.ChangeInInventory),
		ChangeInAccountReceivables:     panicParseInt64ish(cashFlow.ChangeInAccountReceivables),
		ChangeInNetIncome:              panicParseInt64ish(cashFlow.ChangeInNetIncome),
		CapitalExpenditures:            panicParseInt64ish(cashFlow.CapitalExpenditures),
		ChangeInReceivables:            panicParseInt64ish(cashFlow.ChangeInReceivables),
		ChangeInExchangeRate:           panicParseInt64ish(cashFlow.ChangeInExchangeRate),
		ChangeInCashAndCashEquivalents: panicParseInt64ish(cashFlow.ChangeInCashAndCashEquivalents),
	}, nil
}

type rawIncomeStatementItem struct {
	FiscalDateEnding                  string `json:"fiscalDateEnding"`
	ReportedCurrency                  string `json:"reportedCurrency"`
	TotalRevenue                      string `json:"totalRevenue"`
	TotalOperatingExpense             string `json:"totalOperatingExpense"`
	CostOfRevenue                     string `json:"costOfRevenue"`
	GrossProfit                       string `json:"grossProfit"`
	Ebit                              string `json:"ebit"`
	NetIncome                         string `json:"netIncome"`
	ResearchAndDevelopment            string `json:"researchAndDevelopment"`
	EffectOfAccountingCharges         string `json:"effectOfAccountingCharges"`
	IncomeBeforeTax                   string `json:"incomeBeforeTax"`
	MinorityInterest                  string `json:"minorityInterest"`
	SellingGeneralAdministrative      string `json:"sellingGeneralAdministrative"`
	OtherNonOperatingIncome           string `json:"otherNonOperatingIncome"`
	OperatingIncome                   string `json:"operatingIncome"`
	OtherOperatingExpense             string `json:"otherOperatingExpense"`
	InterestExpense                   string `json:"interestExpense"`
	TaxProvision                      string `json:"taxProvision"`
	InterestIncome                    string `json:"interestIncome"`
	NetInterestIncome                 string `json:"netInterestIncome"`
	ExtraordinaryItems                string `json:"extraordinaryItems"`
	NonRecurring                      string `json:"nonRecurring"`
	OtherItems                        string `json:"otherItems"`
	IncomeTaxExpense                  string `json:"incomeTaxExpense"`
	TotalOtherIncomeExpense           string `json:"totalOtherIncomeExpense"`
	DiscontinuedOperations            string `json:"discontinuedOperations"`
	NetIncomeFromContinuingOperations string `json:"netIncomeFromContinuingOperations"`
	NetIncomeApplicableToCommonShares string `json:"netIncomeApplicableToCommonShares"`
	PreferredStockAndOtherAdjustments string `json:"preferredStockAndOtherAdjustments"`
}

// IncomeStatement parsed version of IncomeStatement data received from alphavantage
type IncomeStatement struct {
	FormType                          formtype.FormType `json:"formType"`
	FiscalDateEnding                  Date              `json:"fiscalDateEnding"`
	ReportedCurrency                  string            `json:"reportedCurrency"`
	TotalRevenue                      int64             `json:"totalRevenue"`
	TotalOperatingExpense             int64             `json:"totalOperatingExpense"`
	CostOfRevenue                     int64             `json:"costOfRevenue"`
	GrossProfit                       int64             `json:"grossProfit"`
	Ebit                              int64             `json:"ebit"`
	NetIncome                         int64             `json:"netIncome"`
	ResearchAndDevelopment            int64             `json:"researchAndDevelopment"`
	EffectOfAccountingCharges         int64             `json:"effectOfAccountingCharges"`
	IncomeBeforeTax                   int64             `json:"incomeBeforeTax"`
	MinorityInterest                  int64             `json:"minorityInterest"`
	SellingGeneralAdministrative      int64             `json:"sellingGeneralAdministrative"`
	OtherNonOperatingIncome           int64             `json:"otherNonOperatingIncome"`
	OperatingIncome                   int64             `json:"operatingIncome"`
	OtherOperatingExpense             int64             `json:"otherOperatingExpense"`
	InterestExpense                   int64             `json:"interestExpense"`
	TaxProvision                      int64             `json:"taxProvision"`
	InterestIncome                    int64             `json:"interestIncome"`
	NetInterestIncome                 int64             `json:"netInterestIncome"`
	ExtraordinaryItems                int64             `json:"extraordinaryItems"`
	NonRecurring                      int64             `json:"nonRecurring"`
	OtherItems                        int64             `json:"otherItems"`
	IncomeTaxExpense                  int64             `json:"incomeTaxExpense"`
	TotalOtherIncomeExpense           int64             `json:"totalOtherIncomeExpense"`
	DiscontinuedOperations            int64             `json:"discontinuedOperations"`
	NetIncomeFromContinuingOperations int64             `json:"netIncomeFromContinuingOperations"`
	NetIncomeApplicableToCommonShares int64             `json:"netIncomeApplicableToCommonShares"`
	PreferredStockAndOtherAdjustments int64             `json:"preferredStockAndOtherAdjustments"`
}

func fromIncomeStatement(income rawIncomeStatementItem, formType formtype.FormType) (res IncomeStatement, err error) {
	defer func() {
		if r := recover(); r != nil {
			res = IncomeStatement{}
			err = r.(error)
		}
	}()

	return IncomeStatement{
		FormType:                          formType,
		FiscalDateEnding:                  panicParseDate(income.FiscalDateEnding),
		ReportedCurrency:                  income.ReportedCurrency,
		TotalRevenue:                      panicParseInt64ish(income.TotalRevenue),
		TotalOperatingExpense:             panicParseInt64ish(income.TotalOperatingExpense),
		CostOfRevenue:                     panicParseInt64ish(income.CostOfRevenue),
		GrossProfit:                       panicParseInt64ish(income.GrossProfit),
		Ebit:                              panicParseInt64ish(income.Ebit),
		NetIncome:                         panicParseInt64ish(income.NetIncome),
		ResearchAndDevelopment:            panicParseInt64ish(income.ResearchAndDevelopment),
		EffectOfAccountingCharges:         panicParseInt64ish(income.EffectOfAccountingCharges),
		IncomeBeforeTax:                   panicParseInt64ish(income.IncomeBeforeTax),
		MinorityInterest:                  panicParseInt64ish(income.MinorityInterest),
		SellingGeneralAdministrative:      panicParseInt64ish(income.SellingGeneralAdministrative),
		OtherNonOperatingIncome:           panicParseInt64ish(income.OtherNonOperatingIncome),
		OperatingIncome:                   panicParseInt64ish(income.OperatingIncome),
		OtherOperatingExpense:             panicParseInt64ish(income.OtherOperatingExpense),
		InterestExpense:                   panicParseInt64ish(income.InterestExpense),
		TaxProvision:                      panicParseInt64ish(income.TaxProvision),
		InterestIncome:                    panicParseInt64ish(income.InterestIncome),
		NetInterestIncome:                 panicParseInt64ish(income.NetInterestIncome),
		ExtraordinaryItems:                panicParseInt64ish(income.ExtraordinaryItems),
		NonRecurring:                      panicParseInt64ish(income.NonRecurring),
		OtherItems:                        panicParseInt64ish(income.OtherItems),
		IncomeTaxExpense:                  panicParseInt64ish(income.IncomeTaxExpense),
		TotalOtherIncomeExpense:           panicParseInt64ish(income.TotalOtherIncomeExpense),
		DiscontinuedOperations:            panicParseInt64ish(income.DiscontinuedOperations),
		NetIncomeFromContinuingOperations: panicParseInt64ish(income.NetIncomeFromContinuingOperations),
		NetIncomeApplicableToCommonShares: panicParseInt64ish(income.NetIncomeApplicableToCommonShares),
		PreferredStockAndOtherAdjustments: panicParseInt64ish(income.PreferredStockAndOtherAdjustments),
	}, nil
}
