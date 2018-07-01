package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
)

type transaction struct {
	ID               string
	Amount           float32
	AmountCents      uint64
	LocalAmount      float32
	LocalAmountCents uint64
	Side             string
	OperationType    string
	Currency         string
	LocalCurrency    string
	Status           string
	Note             string
	Label            string
}

type transactions struct {
	Transactions []transaction
}

func RetrieveTransactions() *transactions {
	slug := ""
	iban := ""
	secret := ""
	request, _ := http.NewRequest("GET", "https://thirdparty.qonto.eu/v2/transactions?slug="+slug+"&iban="+iban+"&status[]=pending&status[]=completed", nil)
	request.Header.Add("Authorization", slug+":"+secret)
	client := &http.Client{}

	resp , _ := client.Do(request)
	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)

	response := new(transactions)
	json.Unmarshal(b, response)

	if len(response.Transactions) == 0 {
		fmt.Println("ERREUR, aucune transaction récupérée")
		fmt.Println("essayer avec curl")
		fmt.Println("curl -H \"Authorization: [slug]:[secret]\" \"https://thirdparty.qonto.eu/v2/transactions?slug=[slug]&iban=[iban]\"")
		fmt.Println("curl -H \"Authorization: "+slug+":"+secret+"\" \"https://thirdparty.qonto.eu/v2/transactions?slug="+slug+"&iban="+iban+"\"")
		os.Exit(1)
	}

	return response
}