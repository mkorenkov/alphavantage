package alphavantage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/mkorenkov/alphavantage/formtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIBMProfileParse(t *testing.T) {
	rawTestData := []byte(`
	{
		"Symbol": "IBM",
		"AssetType": "Common Stock",
		"Name": "International Business Machines Corporation",
		"Description": "International Business Machines Corporation operates as an integrated solutions and services company worldwide. Its Cloud & Cognitive Software segment offers software for vertical and domain-specific solutions in health, financial services, and Internet of Things (IoT), weather, and security software and services application areas; and customer information control system and storage, and analytics and integration software solutions to support client mission critical on-premise workloads in banking, airline, and retail industries. It also offers middleware and data platform software, including Red Hat, which enables the operation of clients' hybrid multi-cloud environments; and Cloud Paks, WebSphere distributed, and analytics platform software, such as DB2 distributed, information integration, and enterprise content management, as well as IoT, Blockchain and AI/Watson platforms. The company's Global Business Services segment offers business consulting services; system integration, application management, maintenance, and support services for packaged software; finance, procurement, talent and engagement, and industry-specific business process outsourcing services; and IT infrastructure and platform services. Its Global Technology Services segment provides project, managed, outsourcing, and cloud-delivered services for enterprise IT infrastructure environments; and IT infrastructure support services. The company's Systems segment offers servers for businesses, cloud service providers, and scientific computing organizations; data storage products and solutions; and z/OS, an enterprise operating system, as well as Linux. Its Global Financing segment provides lease, installment payment, loan financing, short-term working capital financing, and remanufacturing and remarketing services. The company was formerly known as Computing-Tabulating-Recording Co. and changed its name to International Business Machines Corporation in 1924. The company was founded in 1911 and is headquartered in Armonk, New York.",
		"Exchange": "NYSE",
		"Currency": "USD",
		"Country": "USA",
		"Sector": "Technology",
		"Industry": "Information Technology Services",
		"Address": "One New Orchard Road, Armonk, NY, United States, 10504",
		"FullTimeEmployees": "352600",
		"FiscalYearEnd": "December",
		"LatestQuarter": "2020-06-30",
		"MarketCapitalization": "111277842432",
		"EBITDA": "15576999936",
		"PERatio": "14.0782",
		"PEGRatio": "8.7188",
		"BookValue": "23.076",
		"DividendPerShare": "6.52",
		"DividendYield": "0.0525",
		"EPS": "8.811",
		"RevenuePerShareTTM": "85.058",
		"ProfitMargin": "0.1043",
		"OperatingMarginTTM": "0.1185",
		"ReturnOnAssetsTTM": "0.0362",
		"ReturnOnEquityTTM": "0.4097",
		"RevenueTTM": "75499003904",
		"GrossProfitTTM": "36489000000",
		"DilutedEPSTTM": "8.811",
		"QuarterlyEarningsGrowthYOY": "-0.458",
		"QuarterlyRevenueGrowthYOY": "-0.054",
		"AnalystTargetPrice": "135.19",
		"TrailingPE": "14.0782",
		"ForwardPE": "11.2486",
		"PriceToSalesRatioTTM": "1.4705",
		"PriceToBookRatio": "5.3809",
		"EVToRevenue": "2.202",
		"EVToEBITDA": "11.0066",
		"Beta": "1.2071",
		"52WeekHigh": "158.75",
		"52WeekLow": "90.56",
		"50DayMovingAverage": "124.3953",
		"200DayMovingAverage": "123.3564",
		"SharesOutstanding": "890579008",
		"SharesFloat": "889189445",
		"SharesShort": "21600483",
		"SharesShortPriorMonth": "23242369",
		"ShortRatio": "4.51",
		"ShortPercentOutstanding": "0.02",
		"ShortPercentFloat": "0.0243",
		"PercentInsiders": "0.108",
		"PercentInstitutions": "58.555",
		"ForwardAnnualDividendRate": "6.52",
		"ForwardAnnualDividendYield": "0.0525",
		"PayoutRatio": "0.7358",
		"DividendDate": "2020-09-10",
		"ExDividendDate": "2020-08-07",
		"LastSplitFactor": "2:1",
		"LastSplitDate": "1999-05-27"
	}`)
	testCaseData := testCompanyProfileAPIResponse{}
	err := json.Unmarshal(rawTestData, &testCaseData)
	require.NoError(t, err)

	parsedProfile := CompanyProfileInfo{}
	err = json.Unmarshal(rawTestData, &parsedProfile)
	require.NoError(t, err)

	inputApiElements := reflect.ValueOf(&testCaseData).Elem()
	parsedResponseElements := reflect.ValueOf(&parsedProfile).Elem()
	for i := 0; i < parsedResponseElements.NumField(); i++ {
		varName := parsedResponseElements.Type().Field(i).Name

		expected := inputApiElements.FieldByName(varName).String()

		varType := parsedResponseElements.Type().Field(i).Type
		switch varType {
		case reflect.TypeOf(Money(0)):
			expectedParts := strings.Split(expected, ".")
			require.Equal(t, 2, len(expectedParts), varName)

			actualParts := strings.Split(parsedResponseElements.Field(i).Interface().(Money).String(), ".")
			require.Equal(t, 2, len(actualParts), varName)

			assert.Equal(t, expectedParts[0], actualParts[0])
			for index, value := range expectedParts[1] {
				assert.Equal(t, string(value), string(actualParts[1][index]), varName)
			}
		case reflect.TypeOf(int64(0)):
			currentResult := strconv.Itoa(int(parsedResponseElements.Field(i).Int()))
			assert.Equal(t, expected, currentResult, varName)
		case reflect.TypeOf(Date{}):
			currentResult := parsedResponseElements.Field(i).Interface().(Date).String()
			assert.Equal(t, expected, currentResult, varName)
		case reflect.TypeOf(""):
			currentResult := parsedResponseElements.Field(i).String()
			assert.Equal(t, expected, currentResult, varName)
		default:
			panic(fmt.Sprintf("unexpected type '%s'", varType))
		}
	}
}

func TestIBMIncomeStatement(t *testing.T) {
	testCases := map[string]rawIncomeStatementItem{
		"10K:2019-12-31": {
			FiscalDateEnding:                  "2019-12-31",
			ReportedCurrency:                  "USD",
			TotalRevenue:                      "77147000000",
			TotalOperatingExpense:             "25945000000",
			CostOfRevenue:                     "40659000000",
			GrossProfit:                       "36488000000",
			Ebit:                              "11511000000",
			NetIncome:                         "9431000000",
			ResearchAndDevelopment:            "5989000000",
			EffectOfAccountingCharges:         "None",
			IncomeBeforeTax:                   "10166000000",
			MinorityInterest:                  "144000000",
			SellingGeneralAdministrative:      "19956000000",
			OtherNonOperatingIncome:           "968000000",
			OperatingIncome:                   "10543000000",
			OtherOperatingExpense:             "-614000000",
			InterestExpense:                   "1344000000",
			TaxProvision:                      "731000000",
			InterestIncome:                    "349000000",
			NetInterestIncome:                 "-995000000",
			ExtraordinaryItems:                "-150000000",
			NonRecurring:                      "None",
			OtherItems:                        "None",
			IncomeTaxExpense:                  "731000000",
			TotalOtherIncomeExpense:           "529000000",
			DiscontinuedOperations:            "-4000000",
			NetIncomeFromContinuingOperations: "9435000000",
			NetIncomeApplicableToCommonShares: "9431000000",
			PreferredStockAndOtherAdjustments: "None",
		},
		"10Q:2020-06-30": {
			FiscalDateEnding:                  "2020-06-30",
			ReportedCurrency:                  "USD",
			TotalRevenue:                      "18123000000",
			TotalOperatingExpense:             "6627000000",
			CostOfRevenue:                     "9423000000",
			GrossProfit:                       "8700000000",
			Ebit:                              "1894000000",
			NetIncome:                         "1361000000",
			ResearchAndDevelopment:            "1582000000",
			EffectOfAccountingCharges:         "None",
			IncomeBeforeTax:                   "1571000000",
			MinorityInterest:                  "137000000",
			SellingGeneralAdministrative:      "5045000000",
			OtherNonOperatingIncome:           "-179000000",
			OperatingIncome:                   "2073000000",
			OtherOperatingExpense:             "-202000000",
			InterestExpense:                   "323000000",
			TaxProvision:                      "209000000",
			InterestIncome:                    "23000000",
			NetInterestIncome:                 "-300000000",
			ExtraordinaryItems:                "-1000000",
			NonRecurring:                      "None",
			OtherItems:                        "None",
			IncomeTaxExpense:                  "209000000",
			TotalOtherIncomeExpense:           "-200000000",
			DiscontinuedOperations:            "-1000000",
			NetIncomeFromContinuingOperations: "1362000000",
			NetIncomeApplicableToCommonShares: "1361000000",
			PreferredStockAndOtherAdjustments: "None",
		},
	}

	parseInt := func(v string) int64 {
		switch v {
		case "None":
			return 0
		default:
			res, _ := strconv.Atoi(v)
			return int64(res)
		}
	}

	parseDate := func(v string) Date {
		switch v {
		case "2020-06-30":
			return Date(time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC))
		case "2019-12-31":
			return Date(time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC))
		default:
			panic(fmt.Sprintf("unexpected input date '%s'", v))
		}
	}

	for form, apiResponse := range testCases {
		var formType formtype.FormType
		switch {
		case strings.HasPrefix(form, "10K:"):
			formType = formtype.Form10K
		case strings.HasPrefix(form, "10Q:"):
			formType = formtype.Form10Q
		default:
			panic(fmt.Sprintf("unexpected form type '%s'", form))
		}

		parsedResponse, err := fromIncomeStatement(apiResponse, formType)
		require.NoError(t, err)

		inputApiElements := reflect.ValueOf(&apiResponse).Elem()
		parsedResponseElements := reflect.ValueOf(&parsedResponse).Elem()

		for i := 0; i < parsedResponseElements.NumField(); i++ {
			varName := parsedResponseElements.Type().Field(i).Name
			inputValue := inputApiElements.FieldByName(varName).String()

			varType := parsedResponseElements.Type().Field(i).Type
			switch varType {
			case reflect.TypeOf(int64(0)):
				expectedResult := parseInt(inputValue)
				currentResult := int64(parsedResponseElements.Field(i).Int())
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(Date{}):
				expectedResult := parseDate(inputValue)
				currentResult := parsedResponseElements.Field(i).Interface().(Date)
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(""):
				currentResult := parsedResponseElements.Field(i).String()
				assert.Equal(t, "USD", currentResult, varName)
			case reflect.TypeOf(formtype.Form10K):
				currentResult := parsedResponseElements.Field(i).Interface().(formtype.FormType)
				assert.Equal(t, formType, currentResult, varName)
			default:
				panic(fmt.Sprintf("unexpected type '%s'", varType))
			}
		}
	}
}

func TestIBMCashFlow(t *testing.T) {
	testCases := map[string]rawCashFlowItem{
		"10K:2019-12-31": {
			FiscalDateEnding:               "2019-12-31",
			ReportedCurrency:               "USD",
			Investments:                    "6988000000",
			ChangeInLiabilities:            "-503000000",
			CashflowFromInvestment:         "-26936000000",
			OtherCashflowFromInvestment:    "-31638000000",
			NetBorrowings:                  "16284000000",
			CashflowFromFinancing:          "9042000000",
			OtherCashflowFromFinancing:     "-173000000",
			ChangeInOperatingActivities:    "1159000000",
			NetIncome:                      "9431000000",
			ChangeInCash:                   "-3124000000",
			OperatingCashflow:              "14770000000",
			OtherOperatingCashflow:         "63000000",
			Depreciation:                   "6059000000",
			DividendPayout:                 "-5707000000",
			StockSaleAndPurchase:           "-1361000000",
			ChangeInInventory:              "67000000",
			ChangeInAccountReceivables:     "491000000",
			ChangeInNetIncome:              "-848000000",
			CapitalExpenditures:            "2286000000",
			ChangeInReceivables:            "502000000",
			ChangeInExchangeRate:           "None",
			ChangeInCashAndCashEquivalents: "-3124000000",
		},
		"10Q:2020-06-30": {
			FiscalDateEnding:               "2020-06-30",
			ReportedCurrency:               "USD",
			Investments:                    "-1263000000",
			ChangeInLiabilities:            "0",
			CashflowFromInvestment:         "-1236000000",
			OtherCashflowFromInvestment:    "613000000",
			NetBorrowings:                  "-38000000",
			CashflowFromFinancing:          "-1624000000",
			OtherCashflowFromFinancing:     "-137000000",
			ChangeInOperatingActivities:    "444000000",
			NetIncome:                      "1361000000",
			ChangeInCash:                   "716000000",
			OperatingCashflow:              "3576000000",
			OtherOperatingCashflow:         "290000000",
			Depreciation:                   "1679000000",
			DividendPayout:                 "-1450000000",
			StockSaleAndPurchase:           "-167000000",
			ChangeInInventory:              "None",
			ChangeInAccountReceivables:     "589000000",
			ChangeInNetIncome:              "247000000",
			CapitalExpenditures:            "585000000",
			ChangeInReceivables:            "None",
			ChangeInExchangeRate:           "None",
			ChangeInCashAndCashEquivalents: "716000000",
		},
	}

	parseInt := func(v string) int64 {
		switch v {
		case "None":
			return 0
		default:
			res, _ := strconv.Atoi(v)
			return int64(res)
		}
	}

	parseDate := func(v string) Date {
		switch v {
		case "2020-06-30":
			return Date(time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC))
		case "2019-12-31":
			return Date(time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC))
		default:
			panic(fmt.Sprintf("unexpected input date '%s'", v))
		}
	}

	for form, apiResponse := range testCases {
		var formType formtype.FormType
		switch {
		case strings.HasPrefix(form, "10K:"):
			formType = formtype.Form10K
		case strings.HasPrefix(form, "10Q:"):
			formType = formtype.Form10Q
		default:
			panic(fmt.Sprintf("unexpected form type '%s'", form))
		}

		parsedResponse, err := fromCashFlow(apiResponse, formType)
		require.NoError(t, err)

		inputApiElements := reflect.ValueOf(&apiResponse).Elem()
		parsedResponseElements := reflect.ValueOf(&parsedResponse).Elem()

		for i := 0; i < parsedResponseElements.NumField(); i++ {
			varName := parsedResponseElements.Type().Field(i).Name
			inputValue := inputApiElements.FieldByName(varName).String()

			varType := parsedResponseElements.Type().Field(i).Type
			switch varType {
			case reflect.TypeOf(int64(0)):
				expectedResult := parseInt(inputValue)
				currentResult := int64(parsedResponseElements.Field(i).Int())
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(Date{}):
				expectedResult := parseDate(inputValue)
				currentResult := parsedResponseElements.Field(i).Interface().(Date)
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(""):
				currentResult := parsedResponseElements.Field(i).String()
				assert.Equal(t, "USD", currentResult, varName)
			case reflect.TypeOf(formtype.Form10K):
				currentResult := parsedResponseElements.Field(i).Interface().(formtype.FormType)
				assert.Equal(t, formType, currentResult, varName)
			default:
				panic(fmt.Sprintf("unexpected type '%s'", varType))
			}
		}
	}
}

func TestIBMBalanceSheet(t *testing.T) {
	testCases := map[string]rawBalanceSheetItem{
		"10K:2019-12-31": {
			FiscalDateEnding:                "2019-12-31",
			ReportedCurrency:                "USD",
			TotalAssets:                     "152186000000",
			IntangibleAssets:                "15235000000",
			EarningAssets:                   "None",
			OtherCurrentAssets:              "3997000000",
			TotalLiabilities:                "131202000000",
			TotalShareholderEquity:          "20841000000",
			DeferredLongTermLiabilities:     "3851000000",
			OtherCurrentLiabilities:         "13406000000",
			CommonStock:                     "55895000000",
			RetainedEarnings:                "162954000000",
			OtherLiabilities:                "35519000000",
			Goodwill:                        "58222000000",
			OtherAssets:                     "16369000000",
			Cash:                            "8313000000",
			TotalCurrentLiabilities:         "37701000000",
			ShortTermDebt:                   "8797000000",
			CurrentLongTermDebt:             "8797000000",
			OtherShareholderEquity:          "-198010000000",
			PropertyPlantEquipment:          "10010000000",
			TotalCurrentAssets:              "38420000000",
			LongTermInvestments:             "10786000000",
			NetTangibleAssets:               "-52617000000",
			ShortTermInvestments:            "696000000",
			NetReceivables:                  "23795000000",
			LongTermDebt:                    "54102000000",
			Inventory:                       "1619000000",
			AccountsPayable:                 "15498000000",
			TotalPermanentEquity:            "None",
			AdditionalPaidInCapital:         "None",
			CommonStockTotalEquity:          "55895000000",
			PreferredStockTotalEquity:       "None",
			RetainedEarningsTotalEquity:     "162954000000",
			TreasuryStock:                   "-169413000000",
			AccumulatedAmortization:         "None",
			OtherNonCurrrentAssets:          "14333000000",
			DeferredLongTermAssetCharges:    "None",
			TotalNonCurrentAssets:           "5182000000",
			CapitalLeaseObligations:         "5259000000",
			TotalLongTermDebt:               "54102000000",
			OtherNonCurrentLiabilities:      "35547000000",
			TotalNonCurrentLiabilities:      "93500000000",
			NegativeGoodwill:                "None",
			Warrants:                        "None",
			PreferredStockRedeemable:        "None",
			CapitalSurplus:                  "55447410000",
			LiabilitiesAndShareholderEquity: "152186000000",
			CashAndShortTermInvestments:     "9009000000",
			AccumulatedDepreciation:         "-22018000000",
			CommonStockSharesOutstanding:    "887110455",
		},
		"10Q:2020-06-30": {
			FiscalDateEnding:                "2020-06-30",
			ReportedCurrency:                "USD",
			TotalAssets:                     "154200000000",
			IntangibleAssets:                "14270000000",
			EarningAssets:                   "None",
			OtherCurrentAssets:              "4387000000",
			TotalLiabilities:                "133512000000",
			TotalShareholderEquity:          "20551000000",
			DeferredLongTermLiabilities:     "3787000000",
			OtherCurrentLiabilities:         "13812000000",
			CommonStock:                     "56135000000",
			RetainedEarnings:                "162559000000",
			OtherLiabilities:                "35937000000",
			Goodwill:                        "57833000000",
			OtherAssets:                     "18500000000",
			Cash:                            "12188000000",
			TotalCurrentLiabilities:         "38442000000",
			ShortTermDebt:                   "9289000000",
			CurrentLongTermDebt:             "9289000000",
			OtherShareholderEquity:          "-198143000000",
			PropertyPlantEquipment:          "9709000000",
			TotalCurrentAssets:              "39953000000",
			LongTermInvestments:             "9272000000",
			NetTangibleAssets:               "-51552000000",
			ShortTermInvestments:            "2063000000",
			NetReceivables:                  "19447000000",
			LongTermDebt:                    "55449000000",
			Inventory:                       "1869000000",
			AccountsPayable:                 "15341000000",
			TotalPermanentEquity:            "None",
			AdditionalPaidInCapital:         "None",
			CommonStockTotalEquity:          "56135000000",
			PreferredStockTotalEquity:       "None",
			RetainedEarningsTotalEquity:     "162559000000",
			TreasuryStock:                   "-169386000000",
			AccumulatedAmortization:         "None",
			OtherNonCurrrentAssets:          "14473000000",
			DeferredLongTermAssetCharges:    "None",
			TotalNonCurrentAssets:           "8689000000",
			CapitalLeaseObligations:         "5027000000",
			TotalLongTermDebt:               "55449000000",
			OtherNonCurrentLiabilities:      "35833000000",
			TotalNonCurrentLiabilities:      "95069000000",
			NegativeGoodwill:                "None",
			Warrants:                        "None",
			PreferredStockRedeemable:        "None",
			CapitalSurplus:                  "55686754250",
			LiabilitiesAndShareholderEquity: "154200000000",
			CashAndShortTermInvestments:     "14251000000",
			AccumulatedDepreciation:         "-21957000000",
			CommonStockSharesOutstanding:    "None",
		},
	}

	parseInt := func(v string) int64 {
		switch v {
		case "None":
			return 0
		default:
			res, _ := strconv.Atoi(v)
			return int64(res)
		}
	}

	parseDate := func(v string) Date {
		switch v {
		case "2020-06-30":
			return Date(time.Date(2020, 6, 30, 0, 0, 0, 0, time.UTC))
		case "2019-12-31":
			return Date(time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC))
		default:
			panic(fmt.Sprintf("unexpected input date '%s'", v))
		}
	}

	for form, apiResponse := range testCases {
		var formType formtype.FormType
		switch {
		case strings.HasPrefix(form, "10K:"):
			formType = formtype.Form10K
		case strings.HasPrefix(form, "10Q:"):
			formType = formtype.Form10Q
		default:
			panic(fmt.Sprintf("unexpected form type '%s'", form))
		}

		parsedResponse, err := fromBalanceSheet(apiResponse, formType)
		require.NoError(t, err)

		inputApiElements := reflect.ValueOf(&apiResponse).Elem()
		parsedResponseElements := reflect.ValueOf(&parsedResponse).Elem()

		for i := 0; i < parsedResponseElements.NumField(); i++ {
			varName := parsedResponseElements.Type().Field(i).Name
			inputValue := inputApiElements.FieldByName(varName).String()

			varType := parsedResponseElements.Type().Field(i).Type
			switch varType {
			case reflect.TypeOf(int64(0)):
				expectedResult := parseInt(inputValue)
				currentResult := int64(parsedResponseElements.Field(i).Int())
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(Date{}):
				expectedResult := parseDate(inputValue)
				currentResult := parsedResponseElements.Field(i).Interface().(Date)
				assert.Equal(t, expectedResult, currentResult, varName)
			case reflect.TypeOf(""):
				currentResult := parsedResponseElements.Field(i).String()
				assert.Equal(t, "USD", currentResult, varName)
			case reflect.TypeOf(formtype.Form10K):
				currentResult := parsedResponseElements.Field(i).Interface().(formtype.FormType)
				assert.Equal(t, formType, currentResult, varName)
			default:
				panic(fmt.Sprintf("unexpected type '%s'", varType))
			}
		}
	}
}

func TestCompanyProfile(t *testing.T) {
	httpClient := &fakeHTTPClient{
		StatusCode: http.StatusOK,
		Result: []byte(`
			{
				"Symbol": "IBM",
				"AssetType": "Common Stock",
				"Name": "International Business Machines Corporation",
				"Description": "International Business Machines Corporation operates as an integrated solutions and services company worldwide. Its Cloud & Cognitive Software segment offers software for vertical and domain-specific solutions in health, financial services, and Internet of Things (IoT), weather, and security software and services application areas; and customer information control system and storage, and analytics and integration software solutions to support client mission critical on-premise workloads in banking, airline, and retail industries. It also offers middleware and data platform software, including Red Hat, which enables the operation of clients' hybrid multi-cloud environments; and Cloud Paks, WebSphere distributed, and analytics platform software, such as DB2 distributed, information integration, and enterprise content management, as well as IoT, Blockchain and AI/Watson platforms. The company's Global Business Services segment offers business consulting services; system integration, application management, maintenance, and support services for packaged software; finance, procurement, talent and engagement, and industry-specific business process outsourcing services; and IT infrastructure and platform services. Its Global Technology Services segment provides project, managed, outsourcing, and cloud-delivered services for enterprise IT infrastructure environments; and IT infrastructure support services. The company's Systems segment offers servers for businesses, cloud service providers, and scientific computing organizations; data storage products and solutions; and z/OS, an enterprise operating system, as well as Linux. Its Global Financing segment provides lease, installment payment, loan financing, short-term working capital financing, and remanufacturing and remarketing services. The company was formerly known as Computing-Tabulating-Recording Co. and changed its name to International Business Machines Corporation in 1924. The company was founded in 1911 and is headquartered in Armonk, New York.",
				"Exchange": "NYSE",
				"Currency": "USD",
				"Country": "USA",
				"Sector": "Technology",
				"Industry": "Information Technology Services",
				"Address": "One New Orchard Road, Armonk, NY, United States, 10504",
				"FullTimeEmployees": "352600",
				"FiscalYearEnd": "December",
				"LatestQuarter": "2020-06-30",
				"MarketCapitalization": "111277842432",
				"EBITDA": "15576999936",
				"PERatio": "14.0782",
				"PEGRatio": "8.7188",
				"BookValue": "23.076",
				"DividendPerShare": "6.52",
				"DividendYield": "0.0525",
				"EPS": "8.811",
				"RevenuePerShareTTM": "85.058",
				"ProfitMargin": "0.1043",
				"OperatingMarginTTM": "0.1185",
				"ReturnOnAssetsTTM": "0.0362",
				"ReturnOnEquityTTM": "0.4097",
				"RevenueTTM": "75499003904",
				"GrossProfitTTM": "36489000000",
				"DilutedEPSTTM": "8.811",
				"QuarterlyEarningsGrowthYOY": "-0.458",
				"QuarterlyRevenueGrowthYOY": "-0.054",
				"AnalystTargetPrice": "135.19",
				"TrailingPE": "14.0782",
				"ForwardPE": "11.2486",
				"PriceToSalesRatioTTM": "1.4705",
				"PriceToBookRatio": "5.3809",
				"EVToRevenue": "2.202",
				"EVToEBITDA": "11.0066",
				"Beta": "1.2071",
				"52WeekHigh": "158.75",
				"52WeekLow": "90.56",
				"50DayMovingAverage": "124.3953",
				"200DayMovingAverage": "123.3564",
				"SharesOutstanding": "890579008",
				"SharesFloat": "889189445",
				"SharesShort": "21600483",
				"SharesShortPriorMonth": "23242369",
				"ShortRatio": "4.51",
				"ShortPercentOutstanding": "0.02",
				"ShortPercentFloat": "0.0243",
				"PercentInsiders": "0.108",
				"PercentInstitutions": "58.555",
				"ForwardAnnualDividendRate": "6.52",
				"ForwardAnnualDividendYield": "0.0525",
				"PayoutRatio": "0.7358",
				"DividendDate": "2020-09-10",
				"ExDividendDate": "2020-08-07",
				"LastSplitFactor": "2:1",
				"LastSplitDate": "1999-05-27"
			}
		`),
	}
	ctx := context.TODO()

	data, err := CompanyProfile(ctx, httpClient, "demo", "IBM")
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Equal(t, "https://www.alphavantage.co/query?function=OVERVIEW&symbol=IBM&apikey=demo", httpClient.Request.URL.String())
	require.Equal(t, "IBM", data.Symbol)
}

func TestBalanceSheet(t *testing.T) {
	httpClient := &fakeHTTPClient{
		StatusCode: http.StatusOK,
		Result: []byte(`
		{
			"symbol": "IBM",
			"annualReports": [
				{
					"fiscalDateEnding": "2019-12-31",
					"reportedCurrency": "USD",
					"totalAssets": "152186000000",
					"intangibleAssets": "15235000000",
					"earningAssets": "None",
					"otherCurrentAssets": "3997000000",
					"totalLiabilities": "131202000000",
					"totalShareholderEquity": "20841000000",
					"deferredLongTermLiabilities": "3851000000",
					"otherCurrentLiabilities": "13406000000",
					"commonStock": "55895000000",
					"retainedEarnings": "162954000000",
					"otherLiabilities": "35519000000",
					"goodwill": "58222000000",
					"otherAssets": "16369000000",
					"cash": "8313000000",
					"totalCurrentLiabilities": "37701000000",
					"shortTermDebt": "8797000000",
					"currentLongTermDebt": "8797000000",
					"otherShareholderEquity": "-198010000000",
					"propertyPlantEquipment": "10010000000",
					"totalCurrentAssets": "38420000000",
					"longTermInvestments": "10786000000",
					"netTangibleAssets": "-52617000000",
					"shortTermInvestments": "696000000",
					"netReceivables": "23795000000",
					"longTermDebt": "54102000000",
					"inventory": "1619000000",
					"accountsPayable": "15498000000",
					"totalPermanentEquity": "None",
					"additionalPaidInCapital": "None",
					"commonStockTotalEquity": "55895000000",
					"preferredStockTotalEquity": "None",
					"retainedEarningsTotalEquity": "162954000000",
					"treasuryStock": "-169413000000",
					"accumulatedAmortization": "None",
					"otherNonCurrrentAssets": "14333000000",
					"deferredLongTermAssetCharges": "None",
					"totalNonCurrentAssets": "5182000000",
					"capitalLeaseObligations": "5259000000",
					"totalLongTermDebt": "54102000000",
					"otherNonCurrentLiabilities": "35547000000",
					"totalNonCurrentLiabilities": "93500000000",
					"negativeGoodwill": "None",
					"warrants": "None",
					"preferredStockRedeemable": "None",
					"capitalSurplus": "55447410000",
					"liabilitiesAndShareholderEquity": "152186000000",
					"cashAndShortTermInvestments": "9009000000",
					"accumulatedDepreciation": "-22018000000",
					"commonStockSharesOutstanding": "887110455"
				}
			],
			"quarterlyReports": [
				{
					"fiscalDateEnding": "2015-09-30",
					"reportedCurrency": "USD",
					"totalAssets": "108649000000",
					"intangibleAssets": "2775000000",
					"earningAssets": "None",
					"otherCurrentAssets": "4158000000",
					"totalLiabilities": "95198000000",
					"totalShareholderEquity": "13294000000",
					"deferredLongTermLiabilities": "3593000000",
					"otherCurrentLiabilities": "14360000000",
					"commonStock": "53220000000",
					"retainedEarnings": "142898000000",
					"otherLiabilities": "29344000000",
					"goodwill": "30275000000",
					"otherAssets": "55875000000",
					"cash": "9480000000",
					"totalCurrentLiabilities": "33732000000",
					"shortTermDebt": "7538000000",
					"currentLongTermDebt": "None",
					"otherShareholderEquity": "-28155000000",
					"propertyPlantEquipment": "10661000000",
					"totalCurrentAssets": "42112000000",
					"longTermInvestments": "9517000000",
					"netTangibleAssets": "-19756000000",
					"shortTermInvestments": "88000000",
					"netReceivables": "26773000000",
					"longTermDebt": "32122000000",
					"inventory": "1613000000",
					"accountsPayable": "11834000000",
					"totalPermanentEquity": "0",
					"additionalPaidInCapital": "0",
					"commonStockTotalEquity": "53220000000",
					"preferredStockTotalEquity": "0",
					"retainedEarningsTotalEquity": "142898000000",
					"treasuryStock": "-154669000000",
					"accumulatedAmortization": "None",
					"otherNonCurrrentAssets": "9619000000",
					"deferredLongTermAssetCharges": "3690000000",
					"totalNonCurrentAssets": "66537000000",
					"capitalLeaseObligations": "None",
					"totalLongTermDebt": "32122000000",
					"otherNonCurrentLiabilities": "25751000000",
					"totalNonCurrentLiabilities": "61466000000",
					"negativeGoodwill": "None",
					"warrants": "None",
					"preferredStockRedeemable": "None",
					"capitalSurplus": "None",
					"liabilitiesAndShareholderEquity": "108649000000",
					"cashAndShortTermInvestments": "9568000000",
					"accumulatedDepreciation": "-18568000000",
					"commonStockSharesOutstanding": "979000000"
				}
			]
		}
		`),
	}
	ctx := context.TODO()

	data, err := BalanceSheets(ctx, httpClient, "demo", "IBM")
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Equal(t, "https://www.alphavantage.co/query?function=BALANCE_SHEET&symbol=IBM&apikey=demo", httpClient.Request.URL.String())
	require.Equal(t, 2, len(data))
}

func TestIncomeStatement(t *testing.T) {
	httpClient := &fakeHTTPClient{
		StatusCode: http.StatusOK,
		Result: []byte(`
		{
			"symbol": "IBM",
			"annualReports": [
				{
					"fiscalDateEnding": "2019-12-31",
					"reportedCurrency": "USD",
					"totalRevenue": "77147000000",
					"totalOperatingExpense": "25945000000",
					"costOfRevenue": "40659000000",
					"grossProfit": "36488000000",
					"ebit": "11511000000",
					"netIncome": "9431000000",
					"researchAndDevelopment": "5989000000",
					"effectOfAccountingCharges": "None",
					"incomeBeforeTax": "10166000000",
					"minorityInterest": "144000000",
					"sellingGeneralAdministrative": "19956000000",
					"otherNonOperatingIncome": "968000000",
					"operatingIncome": "10543000000",
					"otherOperatingExpense": "-614000000",
					"interestExpense": "1344000000",
					"taxProvision": "731000000",
					"interestIncome": "349000000",
					"netInterestIncome": "-995000000",
					"extraordinaryItems": "-150000000",
					"nonRecurring": "None",
					"otherItems": "None",
					"incomeTaxExpense": "731000000",
					"totalOtherIncomeExpense": "529000000",
					"discontinuedOperations": "-4000000",
					"netIncomeFromContinuingOperations": "9435000000",
					"netIncomeApplicableToCommonShares": "9431000000",
					"preferredStockAndOtherAdjustments": "None"
				}
			],
			"quarterlyReports": [
				{
					"fiscalDateEnding": "2020-06-30",
					"reportedCurrency": "USD",
					"totalRevenue": "18123000000",
					"totalOperatingExpense": "6627000000",
					"costOfRevenue": "9423000000",
					"grossProfit": "8700000000",
					"ebit": "1894000000",
					"netIncome": "1361000000",
					"researchAndDevelopment": "1582000000",
					"effectOfAccountingCharges": "None",
					"incomeBeforeTax": "1571000000",
					"minorityInterest": "137000000",
					"sellingGeneralAdministrative": "5045000000",
					"otherNonOperatingIncome": "-179000000",
					"operatingIncome": "2073000000",
					"otherOperatingExpense": "-202000000",
					"interestExpense": "323000000",
					"taxProvision": "209000000",
					"interestIncome": "23000000",
					"netInterestIncome": "-300000000",
					"extraordinaryItems": "-1000000",
					"nonRecurring": "None",
					"otherItems": "None",
					"incomeTaxExpense": "209000000",
					"totalOtherIncomeExpense": "-200000000",
					"discontinuedOperations": "-1000000",
					"netIncomeFromContinuingOperations": "1362000000",
					"netIncomeApplicableToCommonShares": "1361000000",
					"preferredStockAndOtherAdjustments": "None"
				}
			]
		}
		`),
	}
	ctx := context.TODO()

	data, err := IncomeStatements(ctx, httpClient, "demo", "IBM")
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Equal(t, "https://www.alphavantage.co/query?function=INCOME_STATEMENT&symbol=IBM&apikey=demo", httpClient.Request.URL.String())
	require.Equal(t, 2, len(data))
}

func TestCashFlow(t *testing.T) {
	httpClient := &fakeHTTPClient{
		StatusCode: http.StatusOK,
		Result: []byte(`
		{
			"symbol": "IBM",
			"annualReports": [
				{
					"fiscalDateEnding": "2015-12-31",
					"reportedCurrency": "USD",
					"investments": "-629000000",
					"changeInLiabilities": "81000000",
					"cashflowFromInvestment": "-8159000000",
					"otherCashflowFromInvestment": "-3952000000",
					"netBorrowings": "19000000",
					"cashflowFromFinancing": "-9166000000",
					"otherCashflowFromFinancing": "322000000",
					"changeInOperatingActivities": "-3222000000",
					"netIncome": "13190000000",
					"changeInCash": "-317000000",
					"operatingCashflow": "17008000000",
					"otherOperatingCashflow": "-3470000000",
					"depreciation": "3855000000",
					"dividendPayout": "-4897000000",
					"stockSaleAndPurchase": "-4287000000",
					"changeInInventory": "133000000",
					"changeInAccountReceivables": "812000000",
					"changeInNetIncome": "2407000000",
					"capitalExpenditures": "3579000000",
					"changeInReceivables": "812000000",
					"changeInExchangeRate": "-473000000",
					"changeInCashAndCashEquivalents": "-790000000"
				}
			],
			"quarterlyReports": [
				{
					"fiscalDateEnding": "2015-09-30",
					"reportedCurrency": "USD",
					"investments": "272000000",
					"changeInLiabilities": "None",
					"cashflowFromInvestment": "-1343000000",
					"otherCashflowFromInvestment": "-667000000",
					"netBorrowings": "913000000",
					"cashflowFromFinancing": "-1848000000",
					"otherCashflowFromFinancing": "None",
					"changeInOperatingActivities": "None",
					"netIncome": "2950000000",
					"changeInCash": "1044000000",
					"operatingCashflow": "4235000000",
					"otherOperatingCashflow": "-205000000",
					"depreciation": "936000000",
					"dividendPayout": "-1270000000",
					"stockSaleAndPurchase": "-1493000000",
					"changeInInventory": "0",
					"changeInAccountReceivables": "0",
					"changeInNetIncome": "555000000",
					"capitalExpenditures": "948000000",
					"changeInReceivables": "None",
					"changeInExchangeRate": "42000000",
					"changeInCashAndCashEquivalents": "1087000000"
				}
			]
		}
		`),
	}
	ctx := context.TODO()

	data, err := CashFlows(ctx, httpClient, "demo", "IBM")
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Equal(t, "https://www.alphavantage.co/query?function=CASH_FLOW&symbol=IBM&apikey=demo", httpClient.Request.URL.String())
	require.Equal(t, 2, len(data))
}
