package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Launching cost-janitor")

	r := mux.NewRouter()
	r.HandleFunc("/get-monthly-total-cost/{accountid}", GetMonthlyTotalCost)

	println("HTTP server listening on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))); err != nil {
		log.Fatal(err)
	}
}

func GetMonthlyTotalCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
		//		LogLevel: aws.LogLevel(aws.LogDebug),
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	ce := costexplorer.New(sess)

	dateInterval := &costexplorer.DateInterval{}
	dateInterval.SetStart("2021-01-01")
	dateInterval.SetEnd("2021-01-31")

	filter := &costexplorer.Expression{}
	filter.SetDimensions(&costexplorer.DimensionValues{
		Key:          aws.String(costexplorer.DimensionLinkedAccount),
		Values:       []*string{aws.String(vars["accountid"])},
	})

	resp, err := ce.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		Metrics:       []*string{aws.String(costexplorer.MetricBlendedCost)},
		TimePeriod:    dateInterval,
		Granularity: aws.String(costexplorer.GranularityMonthly),
		Filter: filter,
	})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		log.Println(err)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(resp.String()))
}