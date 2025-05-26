package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kopeisec/prometheus-pushgateway-json/conf"
)

// DoRequest sends a request to the Prometheus Pushgateway to push metrics
func DoRequest(endpoint string, header http.Header, jobName string, key string, value string) error {
	// Construct the PushGateway URL
	pushURL := fmt.Sprintf("%s/metrics/job/%s", endpoint, jobName)

	// Create a Prometheus-formatted metric
	metric := fmt.Sprintf("%s %s\n", key, value)

	// Prepare the request body
	body := []byte(metric)

	// Send the POST request to Pushgateway
	req, err := http.NewRequest("POST", pushURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set header
	req.Header = header

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request to Pushgateway: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to push metrics to Pushgateway, status: %s", resp.Status)
	}

	// Success
	return nil
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.POST("/metrics/job/:job_name", func(c *gin.Context) {
		jobName := c.Param("job_name")
		m := make(map[string]string)

		if err := c.ShouldBindJSON(&m); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var errList []error
		for key, value := range m {
			if err := DoRequest(conf.PushGatewayEndpoint(), c.Request.Header, jobName, key, value); err != nil {
				errList = append(errList, err)
			}
		}
		if len(errList) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": errors.Join(errList...).Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg": "OK",
			"job": jobName,
		})
	})
	if err := eng.Run(conf.BindAddr()); err != nil {
		panic(err)
	}
}
