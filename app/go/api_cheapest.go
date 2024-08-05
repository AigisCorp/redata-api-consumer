/*
 * PVPC schema
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package redata_api_consumer

import (
	"encoding/json"
	"sort"

	"github.com/AigisCorp/redata-api-consumer/app/redata"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type CheapestAPI struct {
}

// Get /api/v1/cheapest
// Get cheapest hour
func (api *CheapestAPI) GetCheapest(c *gin.Context) {
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

	// Your handler implementation
	c.JSON(200, gin.H{"Price": price, "Datetime": datetime})
}
