package process

import (
	"qontoEURLApi/api"
	"strings"
	"fmt"
    "time"
    "math"
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
	Immobilisation float32
	IS float32
	Restant float32
	Récupérable float32
}

func GenerateResult(ps []api.Product, printDetails bool) Result {

	var vat float32 = 0
	var vatSales float32 = 0
	var vatPurchases float32 = 0
	var amount float32 = 0
	var vatToSendBack float32 = 0
	var vatLost float32 = 0
	var capital float32 = 1000
	var rémunération float32 = 0
    var immobilisation float32 = 0
    var ImpotSociete float32 = 0
    
	for _, p := range ps {
        var signed float32 = 1
        
		if p.Vat < 0 {
			vatPurchases = vatPurchases + p.Vat
			signed = -1
		} else {
			vatSales = vatSales + p.Vat
		}

		vat = vat + p.Vat
		amount = amount + p.Price

		rémunération = rémunération - p.Remuneration

		if strings.Contains(strings.ToLower(p.Name), "remuneration") {
			rémunération = rémunération + p.Price
		}
		
		if p.Immobilisation {
            var priceHT float32 = (p.Price - p.Vat) * signed
            var prodDate, _ = time.Parse("2006-01-02T15:04:05.000Z", p.ProductDate)
            var yearDiff = (time.Now().Sub(prodDate).Hours() / 24 / 365)
            if(yearDiff < 3) {
                roundYearDiff := float32(math.Max(math.Round(yearDiff), 1.0))
                immobilisation = immobilisation + (priceHT - priceHT / roundYearDiff / 3)
            }
        }
	}

	vatPurchases = vatPurchases * -1

	if vatPurchases < vatSales {
		vatToSendBack = vatSales - vatPurchases
	}else
	{
		vatToSendBack = 0
		vatLost = vatPurchases - vatSales
	}

	rémunération = 	rémunération * -1
	cotisations := rémunération * 0.41

		 
	if(immobilisation < 38120){
      ImpotSociete = ImpotSociete + (immobilisation * 0.15)
    } else{
        ImpotSociete = ImpotSociete + (38120 * 0.15) + ((immobilisation - 38120) * 0.28)
    }
    
	var finalResult = amount - capital - vatToSendBack - cotisations - ImpotSociete
	récupérable := finalResult / 1.41

	if printDetails {
		fmt.Println("Restant = ", "amount - capital - TVA_A_Rendre - cotisations - ImpotSociete ", amount, "-", capital, "-", vatToSendBack,"-", cotisations,"-",ImpotSociete)
		fmt.Println("Recuperable = ", "Restant / 1.41 ", finalResult, "/", 1.41 )
	}
    
	return Result{amount: amount, vat: vat, capital: capital, Rémunération: rémunération, CotisationsAPayer: cotisations, TVA_Achat: vatPurchases, TVA_Ventes: vatSales, TVA_A_Rendre: vatToSendBack, TVA_Perdue: vatLost,Immobilisation: immobilisation,IS: ImpotSociete, Restant: finalResult, Récupérable: récupérable}
}
