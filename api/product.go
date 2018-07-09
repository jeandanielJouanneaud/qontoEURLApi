package api

import (
	"strings"
	"strconv"
	"os"
)

type Product struct {
	Name string
	Price float32
	Vat float32
}

type Products []Product

func (p Product) ToString() string {
	return p.Name + " " + strconv.FormatFloat(float64(p.Price), 'f', 2, 32) + " " + strconv.FormatFloat(float64(p.Vat), 'f', 2, 32)
}

func RetrieveProducts(useProxy bool, proxy string) Products {
	return RetrieveTransactions(useProxy, proxy).ToProducts()
}

func (t transaction) ToProduct() Product {
	signedAmount := t.Amount
	var signedVat float32 = 0.0

	splittedNote := strings.Split(t.Note,"TVA:")

	if len(splittedNote) > 1 {
		signedVat64, err := strconv.ParseFloat(splittedNote[1], 32)
		signedVat = float32(signedVat64)
		if err != nil {
			println("ERROR", err)
			os.Exit(1)
		}
	}

	if t.Side == "debit" {
		signedAmount = signedAmount * -1
		signedVat = signedVat * -1
	}

	if strings.ToLower(t.Label) == "qonto" {
		signedVat = -1.8
	}

	return Product{Name: t.Label, Price: signedAmount, Vat : signedVat}
}

func (ts transactions) ToProducts() Products {
	var ps Products
	for _, t := range ts.Transactions {
		ps = append(ps, t.ToProduct())
	}
	return ps
}