package dolar

import (
	"encoding/json"
	"fmt"
	"github.com/Big-Sh4rk/Balanz-Project/internal/model"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	ars   = "ARS"
	usd   = "USD"
	ext   = "EXT"
	mep   = "MEP"
	cable = "CABLE"
)

var urls = [...]string{"https://test-algobalanz.herokuapp.com/api/v1/prices/security_id", "https://test-algobalanz.herokuapp.com/api/v1/prices", "https://test-algobalanz.herokuapp.com/api/v1/prices/security_id/"}

func ConsumeAPI() error {
	err := endpoints()
	return checkError(err)
}

//Este metodo debo convertirlo en un pointer reciver para poder trabajar con formatos json no declarados como struct
func endpoints() error {

	errFirst := firstEndPoint()
	if errFirst != nil {
		fmt.Print(errFirst.Error())
		return errFirst
	}
	//Primer Endpoint de la API

	//Segundo Endpoint de la API

	/*
		response, errFirst := http.Get("https://test-algobalanz.herokuapp.com/api/v1/prices/")
		return checkError(errFirst)

		responseData, errRead := ioutil.ReadAll(response.Body)
		return checkError(errRead)
	*/

	//Tercer Endpoint de la API

	return nil
}

func firstEndPoint() error {

	dataMap := make(map[string]string)
	var dataArray []string
	response, err := http.Get("https://test-algobalanz.herokuapp.com/api/v1/prices/security_id")
	defer response.Body.Close()

	if err != nil {
		fmt.Print(err.Error())
	}

	data, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		fmt.Println(err)
	}

	for _, value := range dataMap {
		dataArray = append(dataArray, value)
	}

	for securityId := range dataArray {
		fmt.Println(securityId)
	}

	return nil
}

func secondEndPoint() error {

	/*
		response, err := http.Get("https://test-algobalanz.herokuapp.com/api/v1/prices/")

		if err != nil {
			return err
			os.Exit(1)
		}

		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		var responseObject Response
		err = json.Unmarshal(responseData, &responseObject)
		if err != nil {
			return err
		}
		fmt.Println(responseObject.AL30C0001CCTEXT)
	*/
	return nil
}

func thirdEndPoint() error {
	return nil
}

func calcularDolar(instrument1, instrument2 model.FinancialInstrument, dolarType string) {

	// Validamos que sean del mismo instrumento
	if sameStringValue(instrument1.Symbol, instrument2.Symbol) {
		// Validamos el plazo y el tipo de concurrencia de cada uno
		if sameStringValue(instrument1.SettlementType, instrument2.SettlementType) && isCurrencyCorrect(instrument1.Currency, instrument2.Currency, dolarType) {
			pesos := instrument1.Last.Price
			if dolarType == mep {
				dolares := instrument2.Last.Price
				mep := pesos / dolares
				fmt.Printf("Dolar Mep del instrumento %s, valor %f\n", instrument1.Symbol, mep)
			} else {
				exts := instrument2.Last.Price
				cable := pesos / exts
				fmt.Printf("Dolar Cable del instrumento %s, valor %f\n", instrument1.Symbol, cable)
			}

		}
	} else {
		fmt.Println("Invalid Instrument")
	}

}

func sameStringValue(str1, str2 string) bool {
	return str1 == str2
}

func isCurrencyCorrect(currency1, currency2, dolar string) bool {

	switch dolar {
	case mep:
		return currency1 == ars && currency2 == usd
	case cable:
		return currency1 == ars && currency2 == ext
	default:
		return false
	}
}

func checkError(err error) error {
	if err != nil {
		return err
		os.Exit(1)
	}
	return nil
}
