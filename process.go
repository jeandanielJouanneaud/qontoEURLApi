package main

import (
	"qonto/api"
	"strings"
	"fmt"
)

type Result struct {
	amount float32
	vat float32
	capital float32
	Rémunération float32
	CotisationsAPayer float32
	TVA_Ventes float32
	TVA_Achat float32
	TVA_A_Rendre float32
	TVA_Perdue float32
	Restant float32
	Récupérable float32
}

func GenerateResult(ps []api.Product, printDetails bool) Result {

	var vat float32 = 0
	var vatSales float32 = 0
	var vatPurchases float32 = 0
	var amount float32 = 0
	var vatToSend float32 = 0
	var vatLost float32 = 0
	var capital float32 = 1000
	var rémunération float32 = 0

	for _, p := range ps {
		if p.Vat < 0 {
			vatPurchases = vatPurchases + p.Vat
		} else {
			vatSales = vatSales + p.Vat
		}

		vat = vat + p.Vat
		amount = amount + p.Price

		if strings.Contains(strings.ToLower(p.Name), "remuneration") {
			rémunération = rémunération + p.Price
		}
	}

	vatPurchases = vatPurchases * -1

	if vatPurchases < vatSales {
		vatToSend = vatSales - vatPurchases
	}else
	{
		vatToSend = vatSales
		vatLost = vatPurchases - vatSales
	}

	rémunération = 	rémunération * -1
	cotisations := rémunération * 0.41

	var finalResult = amount - capital - vatToSend - cotisations
	récupérable := finalResult / 1.41


	if printDetails {
		fmt.Println("Restant = ", "amount - capital - TVA_A_Rendre - cotisations ", amount, "-", capital, "-", vatToSend,"-", cotisations )
		fmt.Println("Recuperable = ", "Restant / 1.41 ", finalResult, "/", 1.41 )
	}

	return Result{amount: amount, vat: vat, capital: capital, Rémunération: rémunération, CotisationsAPayer: cotisations, TVA_Achat: vatPurchases, TVA_Ventes: vatSales, TVA_A_Rendre: vatToSend, TVA_Perdue: vatLost, Restant: finalResult, Récupérable: récupérable}
}