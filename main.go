package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"qontoEURLApi/api"
)

func main() {
	argsWithoutProg := os.Args[1:]

	argMap := make(map[string]int)
	for i := 0; i < len(argsWithoutProg); i +=1 {
		if argsWithoutProg[i][0] == '-' {
			argMap[argsWithoutProg[i]] = i + 1
		}
	}

	if argMap["-h"] > 0 {
		printUsage()
		os.Exit(0)
	}

	response := api.RetrieveTransactions()

	products := response.ToProducts()

	if argMap["-p"] > 0 {
		addTestProducts(&products, argsWithoutProg, argMap)
	}

	if argMap["-r"] > 0 {
		addRemuneration(&products, argsWithoutProg, argMap)
	}

	printDetails := argMap["-d"] > 0
	result := GenerateResult(products, printDetails)

	if argMap["-v"] > 0 {
		if argMap["-e"] > 0 {
			for _, p := range products {
				fmt.Println(p.ToString())
			}
		}else{
			for _, p := range products {
				fmt.Println(p)
			}
		}
	}

	fmt.Printf("%+v\n",result)
}

func addRemuneration(p *api.Products, argsWithoutProg []string, argMap map[string]int) {
	remuneration, err := strconv.ParseFloat(argsWithoutProg[argMap["-r"]], 64)
	if err != nil {
		fmt.Println("ERREUR !!", "impossible de convertir " + argsWithoutProg[argMap["-r"]] + " en float32")
		printUsage()
		os.Exit(1)
	}

	*p = append(*p, api.Product{Name: "RemunerationTest", Price: -1 * float32(remuneration), Vat: 0})
}

func addTestProducts(p *api.Products, argsWithoutProg []string, argMap map[string]int) {
	productsInfo := argsWithoutProg[argMap["-p"]]
	splittedProductsInfo := strings.Split(productsInfo, ":")

	for _, splittedProductInfoGroup := range splittedProductsInfo {
		splittedProductInfo := strings.Split(splittedProductInfoGroup, ",")
		price, err := strconv.ParseFloat(splittedProductInfo[0], 32)

		if err != nil {
			fmt.Println("ERREUR !!", "impossible de convertir " + splittedProductInfo[0] + " en float32")
			printUsage()
			os.Exit(1)
		}

		vat, err := strconv.ParseFloat(splittedProductInfo[1], 32)

		if err != nil {
			fmt.Println("ERREUR !!", "impossible de convertir " + splittedProductInfo[1] + " en float32")
			printUsage()
			os.Exit(1)
		}

		testProduct := api.Product{Name:"testProduct", Price: float32(price) * -1, Vat: float32(vat) * -1}

		*p = append(*p, testProduct)
	}
}

func printUsage() {
	fmt.Println("Usage : ")
	fmt.Println("-h : help => see this help ")
	fmt.Println("-v : Verbose => print all details ")
	fmt.Println("-e : export mode => print all details in export mode (without braces) ")
	fmt.Println("-p amount,vat[;amount2,vat2] : add some test products")
	fmt.Println("-r amount : add remuneration")
	fmt.Println("-d : print process details")
}