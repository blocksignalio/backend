package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/blocksignalio/core"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const msgInvalidContract = "Invalid contract address provided."

// (contract, valid).
func extractContract(c *gin.Context) (string, bool) {
	contract := c.Params.ByName("contract")
	// hasSuffix := strings.HasSuffix(contract, ".json")
	// if hasSuffix {
	// 	contract = strings.TrimSuffix(contract, ".json")
	// }
	if ok := core.ValidateAddress(contract); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    codeClientError,
			"message": msgInvalidContract,
			"data":    nil,
		})
		c.Abort()
		return contract, false
	}
	// if hasSuffix {
	// 	value := fmt.Sprintf("attachment; filename=%s.json", contract)
	// 	c.Header("Content-Disposition", value)
	// }
	return contract, true
}

// extractPagination returns the tuple (page, pageSize).
func extractPagination(c *gin.Context) (int, int) {
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "0"
	}

	pageSizeStr := c.Query("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = "100"
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 100
	}
	if pageSize != 50 && pageSize != 100 && pageSize != 200 {
		pageSize = 100
	}

	return page, pageSize
}

func extractTopic(c *gin.Context) string {
	topic := c.Query("topic")
	topic = core.SanitizeHex(topic)
	if len(topic) > 10 {
		topic = topic[:10]
	}
	return topic
}

func registerEndpoints(r *gin.Engine) {
	// Initialize the shared database instance.
	db, err := core.Open()
	if err != nil {
		panic(err)
	}

	// Ping test.
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Events endpoint.
	r.GET("/api/contracts/:contract/events", makeEventsHandler(false))
	r.GET("/api/contracts/:contract/events.json", makeEventsHandler(true))

	// Logs endpoint.
	r.GET("/api/contracts/:contract/logs", makeLogsHandler(db, false))
	r.GET("/api/contracts/:contract/logs.json", makeLogsHandler(db, true))

	// r.GET("/api/contracts/:contract/details", make())
	// r.GET("/api/contracts/:contract/source", make())
}

func respond(c *gin.Context, contract, suffix string, data any, download bool) {
	if download {
		value := fmt.Sprintf("attachment; filename=%s-%s.json", contract, suffix)
		c.Header("Content-Disposition", value)
		c.JSON(http.StatusOK, data)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    codeSuccess,
			"message": "success",
			"data":    data,
		})
	}
}

func makeEventsHandler(download bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		contract, valid := extractContract(c)
		if !valid {
			return
		}

		events, err := core.GetContractEvents(contract)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    codeServerError,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		respond(c, contract, "events", adaptEvents(events), download)
	}
}

// GET parameters:
//   - page
//   - pageSize
func makeLogsHandler(db *gorm.DB, download bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		contract, valid := extractContract(c)
		if !valid {
			return
		}

		if err := core.BackfillLogs(c.Request.Context(), db, contract); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    codeServerError,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		topic := extractTopic(c)
		page, pageSize := extractPagination(c)
		logs, err := core.SelectLogs(db, contract, topic, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    codeServerError,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}

		respond(c, contract, "logs", adaptLogs(logs), download)
	}
}
