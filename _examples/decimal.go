package main

import (
	"fmt"
	"github.com/earncef/decimal"
)

func main() {
	price, err := decimal.NewSErr("136.02")
	if err != nil {
		panic(err)
	}

	quantity := decimal.NewF(3)

	fee, _ := decimal.NewSErr(".035")
	taxRate, _ := decimal.NewSErr(".08875")

	subtotal := price.Mul(quantity)

	preTax := subtotal.Mul(fee.Add(decimal.NewF(1)))

	total := preTax.Mul(taxRate.Add(decimal.NewF(1)))

	fmt.Println("Subtotal:", subtotal)                      // Subtotal: 408.06
	fmt.Println("Pre-tax:", preTax)                         // Pre-tax: 422.3421
	fmt.Println("Taxes:", total.Sub(preTax))                // Taxes: 37.482861375
	fmt.Println("Total:", total)                            // Total: 459.824961375
	fmt.Println("Tax rate:", total.Sub(preTax).Div(preTax)) // Tax rate: 0.08875
}
