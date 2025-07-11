package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Rice struct {
	BatchID     string `json:"batchID"`
	Variety     string `json:"variety"`
	HarvestDate string `json:"harvestDate"`
	Quantity    string `json:"quantity"`    // Keep as string to match chaincode
	FarmerName  string `json:"farmerName"`
}

func main() {
	router := gin.Default()
	router.Static("/public", "./public")
	router.LoadHTMLGlob("templates/*")

	// Load home page
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Create Rice Batch (Org1 - Farmer)
	router.POST("/api/rice", func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("❌ PANIC CURSED:", r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Chaincode transaction failed",
					"error":   fmt.Sprint(r),
				})
			}
		}()

		var req Rice
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request", "error": err.Error()})
			return
		}

		fmt.Println("➡️ Request Data:", req)

		result := submitTxnFn("org1", "mychannel", "rice", "RiceContract", "invoke",
			map[string][]byte{}, "CreateRiceBatch",
			req.BatchID, req.Variety, req.HarvestDate, req.Quantity, req.FarmerName)

		c.JSON(http.StatusOK, gin.H{
			"message": "Batch created",
			"data":    result,
		})
	})

	// Read Rice Batch (by Key)
	router.GET("/api/rice/:id", func(ctx *gin.Context) {
		// Prevent panic on unexpected error
		defer func() {
			if r := recover(); r != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to query the rice batch",
					"error":   fmt.Sprint(r),
				})
			}
		}()

		batchID := ctx.Param("id")
		if batchID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Batch ID is required"})
			return
		}

		// Call Chaincode Read
		result := submitTxnFn("org1", "mychannel", "rice", "RiceContract", "query",
			make(map[string][]byte), "ReadRiceBatch", batchID)

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Fetched rice batch: %s", batchID),
			"data":    result,
		})
	})

	router.GET("/api/rice/all", func(c *gin.Context) {
		result := submitTxnFn("org1", "mychannel", "rice", "RiceContract", "query",
			map[string][]byte{}, "GetAllRiceBatches")
		c.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/rice/range", func(c *gin.Context) {
		start := c.Query("start")
		end := c.Query("end")
		result := submitTxnFn("org1", "mychannel", "rice", "RiceContract", "query", map[string][]byte{}, "GetRiceBatchByRange", start, end)
		c.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.GET("/api/rice/history/:id", func(c *gin.Context) {
		batchID := c.Param("id")
		result := submitTxnFn("org1", "mychannel", "rice", "RiceContract", "query", map[string][]byte{}, "GetRiceBatchHistory", batchID)
		c.JSON(http.StatusOK, gin.H{"data": result})
	})

	router.POST("/api/orders", func(c *gin.Context) {
		type ProcessOrder struct {
			Variety     string `json:"variety"`
			Quantity    string `json:"quantityInKg"`
			MillerName  string `json:"millerName"`
			OrderID     string `json:"orderID"`
		}
		var req ProcessOrder
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		transient := map[string][]byte{
			"variety":     []byte(req.Variety),
			"quantityInKg": []byte(req.Quantity),
			"millerName":  []byte(req.MillerName),
		}
		result := submitTxnFn("org2", "mychannel", "rice", "RiceContract", "private", transient, "CreateProcessingOrder", req.OrderID)
		c.JSON(http.StatusOK, gin.H{"message": result})
	})

	router.POST("/api/orders/match", func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("‼️ MatchProcessingOrder PANIC:", r)
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"message": "Failed to match processing order",
					"error":   fmt.Sprint(r),
				})
			}
		}()

		// Struct to hold expected JSON input
		type MatchInfo struct {
			BatchID string `json:"batchID"`
			OrderID string `json:"orderID"`
		}

		var data MatchInfo

		// Bind incoming JSON to struct
		if err := ctx.BindJSON(&data); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad match input", "details": err.Error()})
			return
		}

		fmt.Println("Matching Batch:", data.BatchID)
		fmt.Println("With Order    :", data.OrderID)

		// Call chaincode
		result := submitTxnFn("org2", "mychannel", "rice", "RiceContract",
			"invoke", map[string][]byte{}, "MatchProcessingOrder",
			data.BatchID, data.OrderID)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Matched batch to order",
			"result":  result,
		})
	})

	router.POST("/api/rice/dispatch", func(c *gin.Context) {
		type Dispatch struct {
			BatchID     string `json:"batchID"`
			RetailerName string `json:"retailerName"`
		}
		var req Dispatch
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}
		result := submitTxnFn("org3", "mychannel", "rice", "RiceContract", "invoke", map[string][]byte{}, "DispatchToRetailer", req.BatchID, req.RetailerName)
		c.JSON(http.StatusOK, gin.H{"message": result})
	})

	


	// Start the server
	router.Run("localhost:3001")
}