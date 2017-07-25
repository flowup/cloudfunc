package main

import (
	"github.com/flowup/cloudfunc/api"
	"math"
)

type Config struct {
	Amount     float64 `json:"amount"`
	TermMonths int64 `json:"termMonths"`
	Rate       float64 `json:"rate"`
	Payment    float64 `json:"payment"`
}

func (c *Config) MonthlyRate() float64 {
	return c.Rate / 12
}

type AmortizationTable struct {
	Config *Config `json:"config"`
	Rows   []AmortizationRow `json:"rows"`
}

type AmortizationRow struct {
	Payment       float64 `json:"payment"`
	Interest      float64 `json:"interest"`
	Principal     float64 `json:"principal"`
	EndingBalance float64 `json:"endingBalance"`
}

// PMT calculates monthly payments for the mortgage
func CalculatePayment(amount float64, rate float64, months int64) float64 {
	return (amount * (rate / 12)) / (1 - math.Pow(1+(rate/12), float64(-1*months)))
}

// ConstructAmortizationTable creates an amortization table from the configuration
func ConstructAmortizationTable(config *Config) *AmortizationTable {
	table := &AmortizationTable{
		Config: config,
		Rows:   make([]AmortizationRow, config.TermMonths),
	}

	interest := config.Amount * config.MonthlyRate()
	principal := config.Payment - interest
	endingBalance := config.Amount - principal

	table.Rows[0].Payment = config.Payment
	table.Rows[0].Interest = interest
	table.Rows[0].Principal = principal
	table.Rows[0].EndingBalance = endingBalance

	for i := int64(1); i < config.TermMonths; i++ {
		interest = table.Rows[i-1].EndingBalance * config.MonthlyRate()
		principal = config.Payment - interest
		endingBalance = table.Rows[i-1].EndingBalance - principal

		table.Rows[i] = AmortizationRow{
			Payment: config.Payment,
			Interest: interest,
			Principal: principal,
			EndingBalance: endingBalance,
		}
	}

	return table
}

func main() {
	// create configuration and get user input
	config := &Config{}
	api.GetInput(&config)
	config.Payment = CalculatePayment(config.Amount, config.Rate, config.TermMonths)

	// send back the result
	api.Send(ConstructAmortizationTable(config))
}
