package dolar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Big-Sh4rk/Balanz-Project/internal/model"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

const (
	ars          = "ARS"
	usd          = "USD"
	ext          = "EXT"
	firstEndUrl  = "https://test-algobalanz.herokuapp.com/api/v1/prices/security_id"
	secondEndUrl = "https://test-algobalanz.herokuapp.com/api/v1/prices"
	thirdEndUrl  = "https://test-algobalanz.herokuapp.com/api/v1/prices/security_id/"
	socketUri    = "wss://test-algobalanz.herokuapp.com/ws/example"
)

func ConsumeAPI() {

	fmt.Println("Consumiendo API...")
	fmt.Println("")
	//Primer Endpoint de la API
	errFirst := firstEndPoint()
	isErrorNotNil(errFirst)

	//Segundo Endpoint de la API
	errSecond := secondEndPoint()
	isErrorNotNil(errSecond)

	//El tercer endpoint es utilizado con los valores obtenidos en el primer endpoint

}

func ConsumeSocket() error {

	c, _, err := websocket.DefaultDialer.Dial(socketUri, nil)
	isErrorNotNil(err)
	defer c.Close()

	var instruments []model.FinancialInstrument

	// Recivimos mensajes
	go func() {
		log.Println("")
		log.Println("Consumiendo Socket...")
		log.Println("")
		for {
			_, message, _ := c.ReadMessage()
			//Se transforma el valor al instrumento deseado
			instrument := byteToInstrument(message)
			//Lo agregamos a la lista
			instruments = append(instruments, instrument)
			//Eliminamos valores viejos con la intencion de mejorar significativamente el rendimiento
			if len(instruments) > 2 {
				instruments = removeOldValues(instruments)
				calcularDolar(instruments)
			}
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		_ = c.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		log.Printf("Scan: %s", scanner.Text())
	}

	return nil
}

func firstEndPoint() error {
	fmt.Println("Extrayendo los Security_id")
	response, err := http.Get(firstEndUrl)
	defer response.Body.Close()

	isErrorNotNil(err)

	body, errRead := ioutil.ReadAll(response.Body)
	isErrorNotNil(errRead)

	var ids model.SecurityIDs
	errJson := json.Unmarshal(body, &ids)
	isErrorNotNil(errJson)

	thirdEndPoint(ids.Response)

	return nil
}

func secondEndPoint() error {
	fmt.Println("Obteniendo todos los instrumentos disponibles en el 2do endpoint")
	response, err := http.Get(secondEndUrl)
	defer response.Body.Close()
	isErrorNotNil(err)

	instruments := mapToFinancialInstrument(response)
	calcularDolar(instruments)

	return nil
}

//En este endpoint usaremos la lista de security_id conseguidos en el primer metodo, calcularemos los respectivos dolares
func thirdEndPoint(ids []string) error {
	fmt.Println("Tercer Endpoint")
	fmt.Println("Obteniendo los Instrumentos con los security_ids obtenidos en el primer endpoint")
	var instruments []model.FinancialInstrument
	for _, id := range ids {
		response, err := http.Get(thirdEndUrl + id)
		isErrorNotNil(err)

		body, errRead := ioutil.ReadAll(response.Body)
		isErrorNotNil(errRead)

		var in model.NewInstrument
		errJson := json.Unmarshal(body, &in)
		isErrorNotNil(errJson)

		instruments = append(instruments, in.Response)
		defer response.Body.Close()
	}

	calcularDolar(instruments)
	return nil
}

func calcularDolar(instruments []model.FinancialInstrument) {

	instruments = sortInstruments(instruments)

	max := len(instruments)
	for i := 0; i < max; i++ {
		fIns := instruments[i]
		for j := i + 1; j < max; j++ {
			sIns := instruments[j]
			if !sameInstrument(fIns.Symbol, sIns.Symbol) || !sameST(fIns.SettlementType, sIns.SettlementType) || !diferentCurrency(fIns.Currency, sIns.Currency) {
				continue
			}

			values := mapWithData(fIns.Currency, fIns.Last.Price, sIns.Currency, sIns.Last.Price)
			checkingTypeDolar(values, substractIns(fIns.Symbol))

		}
	}
}

func mapWithData(firstCur string, firstPrice float64, secondCur string, secondPrice float64) map[string]float64 {
	data := make(map[string]float64)
	data[firstCur] = firstPrice
	data[secondCur] = secondPrice
	return data
}

func checkingTypeDolar(data map[string]float64, instrument string) {
	peso, isArs := data[ars]
	dolarU, isUsd := data[usd]
	dolarE, isExt := data[ext]
	//La validacion de los numeros es para evitar que me de un resultado que tienda a infinito.
	if isArs && isUsd && peso != 0.0 && dolarU != 0.0 {
		dolarMep(instrument, peso, dolarU)
	}
	if isArs && isExt && peso != 0.0 && dolarE != 0.0 {
		dolarCable(instrument, peso, dolarE)
	}
}

func dolarMep(instrument string, peso, dolar float64) {

	result := peso / dolar
	fmt.Printf("Dolar MEP %s: %.2f\n", instrument, result)
}

func dolarCable(instrument string, peso, dolar float64) {

	result := peso / dolar
	fmt.Printf("Dolar Cable %s: %.2f\n", instrument, result)
}

func sameInstrument(fInstrument, sInstrument string) bool {
	//Hago estas comparaciones ya que existen instrumentos que tienen un caracter extra
	return substractIns(fInstrument) == substractIns(sInstrument)
}

func substractIns(instrument string) string {
	return instrument[0:4]
}

func sameST(firstST, secondST string) bool {
	return firstST == secondST
}

func diferentCurrency(firstCur, secondCur string) bool {
	return firstCur != secondCur
}

func isErrorNotNil(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//En esta funcion trabajo con Unstructured Data, para luego pasarlo a una lista de una struct ya predefinida
func mapToFinancialInstrument(response *http.Response) []model.FinancialInstrument {

	var newInstruments []model.FinancialInstrument
	dataMap := make(map[string]model.FinancialInstrument)

	body, errRead := ioutil.ReadAll(response.Body)
	isErrorNotNil(errRead)

	errJson := json.Unmarshal(body, &dataMap)
	isErrorNotNil(errJson)

	for _, value := range dataMap {
		newInstruments = append(newInstruments, value)
	}

	return newInstruments
}

func byteToInstrument(response []byte) model.FinancialInstrument {

	socketIns := model.SocketInstrument{}
	err := json.Unmarshal(response, &socketIns)
	isErrorNotNil(err)

	return socketIns.Msg
}

func removeOldValues(instruments []model.FinancialInstrument) []model.FinancialInstrument {

	instruments = sortInstruments(instruments)

	uniqPointer := 0

	for i := 1; i < len(instruments); i++ {
		if instruments[uniqPointer].SecurityID != instruments[i].SecurityID {
			uniqPointer++
			instruments[uniqPointer] = instruments[i]
		}
	}

	return instruments[:uniqPointer+1]
}

func sortInstruments(instruments []model.FinancialInstrument) []model.FinancialInstrument {
	//Ordenamos los instrumentos
	sort.Slice(instruments, func(i, j int) bool {
		return instruments[i].SecurityID < instruments[j].SecurityID
	})
	return instruments
}
