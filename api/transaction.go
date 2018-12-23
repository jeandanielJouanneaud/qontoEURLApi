package api

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"os"
	"net/url"
	"strconv"
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
 	Settled_at       string 
	Emitted_at       string
	Status           string
	Note             string
	Label            string
	Reference        string
	Vat_amount       float32
	Vat_amountCents  float32
	Vat_rate         float32
}

type transactions struct {
	Transactions []transaction
}

func getTransactionUrlRequest(slug string, iban string, current_page int) string {
	return "https://thirdparty.qonto.eu/v2/transactions?slug="+slug+"&iban="+iban+"&per_page=100&current_page="+ strconv.Itoa(current_page)+"&status[]=pending&status[]=completed"
}

func getTransaction(current_page int) *transactions  {

	slug := ""
	iban := ""
	secret := ""

	request, err := http.NewRequest("GET", getTransactionUrlRequest(slug, iban, current_page), nil)

	if err != nil {
		fmt.Println("ERROR while creating Request ", err)
		os.Exit(1)
	}

	request.Header.Add("Authorization", slug+":"+secret)
	client := &http.Client{}

	resp , err := client.Do(request)

	if err != nil {
		fmt.Println("ERROR retrieving http response ", err)
	}

	defer resp.Body.Close()

	b, _ := ioutil.ReadAll(resp.Body)
 
	response := new(transactions)
	json.Unmarshal(b, response)
	return response
}

func getTransactions()  *transactions {
	var trans []transaction
	currentPage := 1
	for {

		response := getTransaction(currentPage)

		if len(response.Transactions) == 0 {
			break
		}

		trans = append(trans, response.Transactions...)
		currentPage = currentPage + 1
	}
	return &transactions{Transactions: trans}
}

func RetrieveTransactions(useProxy bool, proxy string) *transactions {

	if !useProxy  {
		proxy = os.Getenv("http_proxy")
	}

	println(proxy)

	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			fmt.Println("ERROR, Failed to parse proxy ", err)
			os.Exit(1)
		}

		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	}

	response := getTransactions()

		if len(response.Transactions) == 0 {
			fmt.Println("ERREUR, aucune transaction récupérée")
			fmt.Println("essayer avec curl")
	        fmt.Println("curl -H \"Authorization: [slug]:[secret]\" \"https://thirdparty.qonto.eu/v2/transactions?slug=[slug]&iban=[iban]\"")
			os.Exit(1)
		}

	return response
}
