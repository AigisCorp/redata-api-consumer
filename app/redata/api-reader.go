package redata

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type MercadosPreciosTempoReal struct { // https://mholt.github.io/json-to-go/
	Data struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes struct {
			Title       string    `json:"title"`
			LastUpdate  time.Time `json:"last-update"`
			Description any       `json:"description"`
		} `json:"attributes"`
		Meta struct {
			CacheControl struct {
				Cache    string `json:"cache"`
				ExpireAt string `json:"expireAt"`
			} `json:"cache-control"`
		} `json:"meta"`
	} `json:"data"`
	Included []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		GroupID    any    `json:"groupId"`
		Attributes struct {
			Title       string    `json:"title"`
			Description any       `json:"description"`
			Color       string    `json:"color"`
			Type        any       `json:"type"`
			Magnitude   string    `json:"magnitude"`
			Composite   bool      `json:"composite"`
			LastUpdate  time.Time `json:"last-update"`
			Values      []struct {
				Value      float64   `json:"value"`
				Percentage float64   `json:"percentage"`
				Datetime   time.Time `json:"datetime"`
			} `json:"values"`
		} `json:"attributes"`
	} `json:"included"`
}

func TodayString() string {
	t := time.Now()
	today := fmt.Sprintf("%d-%d-%d", t.Year(), t.Month(), t.Day())
	return today
}

func CreateFile(filename string, data []byte) {
	_ = os.WriteFile(filename+".json", data, 0644)

	log.Info("Created file:" + filename + ".json")
}

func GetMercadosPreciosTempoReal() []byte {
	var responseBody []byte
	today := TodayString()
	filename := today + ".json"
	if _, err := os.Stat(filename); err == nil {
		log.Info("Reading existing file: " + filename)
		file, err := os.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		responseBody = file
	} else {
		log.Info("File: '" + filename + "' doesn't exists. Querying REData API for the first time today...")
		apiUrl := "https://apidatos.ree.es/es/datos/mercados/precios-mercados-tiempo-real?start_date=" + today + "T00:00&end_date=" + today + "T23:59&time_trunc=hour"
		log.Info("Making a request to " + apiUrl + "...")
		request, err := http.NewRequest("GET", apiUrl, nil)
		if err != nil {
			fmt.Println(err)
		}
		request.Header.Set("Content-Type", "application/json; charset=utf-8")
		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			fmt.Println(error)
		}

		responseBody, error = io.ReadAll(response.Body)
		if error != nil {
			fmt.Println(error)
		}

		// Uncomment this to debug responseBody
		// var formattedData bytes.Buffer
		// err = json.Indent(&formattedData, responseBody, "", "\t")
		// if err != nil {
		// 	fmt.Println(error)
		// }
		// fmt.Println("Status: ", response.Status)
		// fmt.Println("Response body: ", formattedData.String())

		CreateFile(today, responseBody)

		// clean up memory after execution
		defer response.Body.Close()
	}
	return responseBody
}
