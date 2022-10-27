package dolar

import (
	"encoding/json"
	"fmt"
	"github.com/Big-Sh4rk/Balanz-Project/internal/model"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
)

const (
	ars          = "ARS"
	usd          = "USD"
	ext          = "EXT"
	firstEndUrl  = "https://test-algobalanz.herokuapp.com/api/v1/prices/security_id"
	secondEndUrl = "https://test-algobalanz.herokuapp.com/api/v1/prices"
	thirdEndUrl  = "https://test-algobalanz.herokuapp.com/api/v1/prices/security_id/"
	socket       = "wss://test-algobalanz.herokuapp.com/ws/"
)

func ConsumeAPI() error {

	//Primer Endpoint de la API
	errFirst := firstEndPoint()
	if isErrorNotNil(errFirst) {
		return errFirst
	}

	//Segundo Endpoint de la API
	errSecond := secondEndPoint()
	if isErrorNotNil(errSecond) {
		return errSecond
	}

	//El tercer endpoint es utilizado con los valores obtenidos en el primer endpoint

	return nil
}

func firstEndPoint() error {

	response, err := http.Get(firstEndUrl)
	defer response.Body.Close()

	if isErrorNotNil(err) {
		return err
	}

	body, errRead := ioutil.ReadAll(response.Body)
	if isErrorNotNil(errRead) {
		return errRead
	}
	var ids model.SecurityIDs
	errJson := json.Unmarshal(body, &ids)
	if isErrorNotNil(errJson) {
		return errJson
	}

	thirdEndPoint(ids.Response)

	return nil
}

func secondEndPoint() error {

	response, err := http.Get(secondEndUrl)
	defer response.Body.Close()
	if isErrorNotNil(err) {
		return err
	}

	instruments, errConv := mapToFinancialInstrument(response)
	if isErrorNotNil(errConv) {
		return errConv
	}

	fmt.Println("Calculando Dollars con el segundo endpoint")
	fmt.Println("")
	calcularDolar(instruments)

	return nil
}

//En este endpoint usaremos la lista de security_id conseguidos en el primer metodo, calcularemos los respectivos dolares
func thirdEndPoint(ids []string) error {

	var instruments []model.FinancialInstrument
	for _, id := range ids {
		response, err := http.Get(thirdEndUrl + id)
		if isErrorNotNil(err) {
			return err
		}
		body, errRead := ioutil.ReadAll(response.Body)
		if isErrorNotNil(errRead) {
			return errRead
		}
		var in model.NewInstrument
		errJson := json.Unmarshal(body, &in)
		if isErrorNotNil(errJson) {
			return errJson
		}
		instruments = append(instruments, in.Response)
		defer response.Body.Close()
	}

	fmt.Println("Calculando Dollars con el tercer endpoint")
	fmt.Println("")
	calcularDolar(instruments)

	return nil
}

func calcularDolar(instruments []model.FinancialInstrument) {

	//Ordenamos los instrumentos
	fmt.Println("ORDENANDO LISTA")
	sort.Slice(instruments, func(i, j int) bool {
		return instruments[i].SecurityID < instruments[j].SecurityID
	})
	/*
		OPCIONAL SOLO APLICA A LA API:
		Imprimo los SecurityIDs ordenados para que se vea en consola que los valores respectivos(Instrumento y settlementType)
		son distintos o la Currency como tal es la misma, por ende no llego a calcular ningun tipo de dolar. Para probarlo pueden descomentar
		el bucle de abajo.

		for _, ins := range instruments {
			fmt.Println(ins.SecurityID)
		}*/
	fmt.Println("---------------")

	max := len(instruments)
	for i := 0; i < max; i++ {
		for j := i + 1; j < max; j++ {
			if sameInstrument(instruments[i].Symbol, instruments[j].Symbol) && sameST(instruments[i].SettlementType, instruments[j].SettlementType) {
				if isCurrencyMEP(instruments[i].Currency, instruments[j].Currency) {
					dolarMep(instruments[i], instruments[j])
				} else if isCurrencyCABLE(instruments[i].Currency, instruments[j].Currency) {
					dolarCable(instruments[i], instruments[j])
				}
			}
		}
	}
	fmt.Println("SE TERMINO DE RECORRER LA LISTA")
}

func dolarMep(instrumentARS, instrumentUSD model.FinancialInstrument) {
	valueARS := instrumentARS.Last.Price
	valueUSD := instrumentUSD.Last.Price
	mep := valueARS / valueUSD
	if math.IsNaN(mep) {
		fmt.Println("Resultado negativo")
	} else {
		fmt.Printf("Dolar Mep del instrumento %s, valor %f\n", instrumentARS.Symbol, mep)
	}
}

func dolarCable(instrumentARS, instrumentEXT model.FinancialInstrument) {
	valueARS := instrumentARS.Last.Price
	valueEXT := instrumentEXT.Last.Price
	cable := valueARS / valueEXT
	if math.IsNaN(cable) {
		fmt.Println("Resultado negativo")
	} else {
		fmt.Printf("Dolar Cable del instrumento %s, valor %f\n", instrumentARS.Symbol, cable)
	}
}

//En ambas funciones comparo strings, lo hice de esta forma para una mejor lectura de codigo
func sameInstrument(fInstrument, sInstrument string) bool {
	return fInstrument == sInstrument
}

func sameST(firstST, secondST string) bool {
	return firstST == secondST
}

func isCurrencyMEP(firstCur, secondCur string) bool {
	return firstCur == ars && secondCur == usd
}

func isCurrencyCABLE(firstCur, secondCur string) bool {
	return firstCur == ars && secondCur == ext
}

func isErrorNotNil(err error) bool {
	if err != nil {
		fmt.Print(err.Error())
		return true
	}
	return false
}

//En esta funcion trabajo con Unstructured Data, para luego pasarlo a una lista de una struct ya predefinida
func mapToFinancialInstrument(response *http.Response) ([]model.FinancialInstrument, error) {

	var newInstruments []model.FinancialInstrument
	var dataMap map[string]model.FinancialInstrument

	body, errRead := ioutil.ReadAll(response.Body)
	if isErrorNotNil(errRead) {
		return nil, errRead
	}

	errJson := json.Unmarshal(body, &dataMap)
	if isErrorNotNil(errJson) {
		return nil, errJson
	}

	for _, value := range dataMap {
		newInstruments = append(newInstruments, value)
	}

	return newInstruments, nil
}
