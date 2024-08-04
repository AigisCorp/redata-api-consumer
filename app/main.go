package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"sort"
	"sync"
	"time"
	_ "time/tzdata"

	log "github.com/sirupsen/logrus"

	api "aigiscorp.dev/redata-api-consumer/ogen-redata/api"
	redata "aigiscorp.dev/redata-api-consumer/redata"
)

//go:embed static
var static embed.FS

type apiService struct {
	price  map[int32]api.Price
	charge api.Charge
	mux    sync.Mutex
}

func (p *apiService) Charge(ctx context.Context) (api.ChargeRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	responseBody := redata.GetMercadosPreciosTempoReal()
	data := redata.MercadosPreciosTempoReal{}
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		log.Fatal(err)
	}

	prices := data.Included[0].Attributes.Values
	var total float64
	for v := range prices {
		total = total + data.Included[0].Attributes.Values[v].Value
	}
	mean := total / 24

	log.Info(mean)

	var charge bool
	t := time.Now()
	currentHour := t.Hour()
	priceCurrentHour := data.Included[0].Attributes.Values[currentHour].Value
	if priceCurrentHour < mean {
		charge = true
		log.Info(fmt.Sprintf("Current price %.2f, is under the mean price %.2f", priceCurrentHour, mean))
	} else {
		charge = false
		log.Info(fmt.Sprintf("Current price %.2f, is over the mean price %.2f", priceCurrentHour, mean))
	}

	return &api.Charge{Charge: charge}, nil
}

func (p *apiService) GetCheap(ctx context.Context, params api.GetCheapParams) (api.GetCheapRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	queryHours := params.Hours
	if queryHours > 24 {
		msg := map[string]interface{}{
			"Hours": "A day only has 24 hours!",
		}
		errorMessage, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		return &api.Prices{"Error": errorMessage}, nil
	}
	responseBody := redata.GetMercadosPreciosTempoReal()
	data := redata.MercadosPreciosTempoReal{}
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		log.Fatal(err)
	}

	prices := data.Included[0].Attributes.Values
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Value < prices[j].Value
	})

	type price struct {
		Price    float64 `json:"price"`
		Datetime string  `json:"datetime"`
	}
	var hours []price
	var i uint8
	for i = 0; i < queryHours; i++ {
		hours = append(hours,
			price{Price: data.Included[0].Attributes.Values[i].Value, Datetime: data.Included[0].Attributes.Values[i].Datetime.String()})
	}

	result, err := json.Marshal(hours)
	if err != nil {
		log.Fatal(err)
	}

	return &api.Prices{"result": result}, nil
}

func (p *apiService) GetCheapest(ctx context.Context) (api.GetCheapestRes, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	responseBody := redata.GetMercadosPreciosTempoReal()
	data := redata.MercadosPreciosTempoReal{}
	err := json.Unmarshal(responseBody, &data)
	if err != nil {
		log.Fatal(err)
	}

	prices := data.Included[0].Attributes.Values
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Value < prices[j].Value
	})

	price := data.Included[0].Attributes.Values[0].Value
	datetime := data.Included[0].Attributes.Values[0].Datetime.String()

	return &api.Price{Price: price, Datetime: datetime}, nil
}

// func main() {
// 	log.SetOutput(os.Stdout)
// 	log.SetLevel(log.DebugLevel)
// 	log.SetFormatter(&log.JSONFormatter{})
// 	log.Info("Starting...")

// 	// Create service instance.
// 	service := &apiService{
// 		price:  map[int32]api.Price{},
// 		charge: api.Charge{},
// 	}
// 	// Create generated server.
// 	srv, err := api.NewServer(service)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if err := http.ListenAndServe(":8080", srv); err != nil {
// 		log.Fatal(err)
// 	}
// }

func main() {
	location, err := time.LoadLocation(os.Getenv("location"))
	if err != nil {
		log.Fatal("Error loading location: %v", err)
	}
	time.Local = location

	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Starting...")

	// Create service instance.
	service := &apiService{
		price:  map[int32]api.Price{},
		charge: api.Charge{},
	}
	srv, err := api.NewServer(service)
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", srv))
	// Register static files.
	mux.Handle("/static/", http.FileServer(http.FS(static)))
	{
		// Register pprof handlers.
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	http.ListenAndServe(":8080", mux)
}
