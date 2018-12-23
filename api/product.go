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
	Remuneration float32
	Immobilisation bool
	ProductDate string
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
	
	var signedRem float32 = 0.0
  	splittedRem := strings.Split(t.Note,"REMUNERATION:")

	if len(splittedRem) > 1 {
        splittedFinalRem := strings.Split(splittedRem[1], "\n")[0]
		signedRem64, err := strconv.ParseFloat(splittedFinalRem, 32)
		signedRem = float32(signedRem64)
		if err != nil {
			println("ERROR", err)
			os.Exit(1)
		}
	}
	
	var isImmo bool = len(strings.Split(t.Note,"IMMOBILISATION")) > 1

	var signedVat float32 = t.Vat_amount
    
	if t.Side == "debit" {
		signedAmount = signedAmount * -1
		signedVat = signedVat * -1
	}

	if strings.ToLower(t.Label) == "qonto" {
		signedVat = -1.8
	}
    
	return Product{Name: t.Label, Price: signedAmount, Vat : signedVat, ProductDate: t.Emitted_at, Remuneration: signedRem, Immobilisation: isImmo}
}

func (ts transactions) ToProducts() Products {
	var ps Products
	for _, t := range ts.Transactions {
		ps = append(ps, t.ToProduct())
	}
	return ps
}
