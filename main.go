package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"log"
	"os"
)

func main() {
	fmt.Println("Launching cost-janitor")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
//		LogLevel: aws.LogLevel(aws.LogDebug),
	})
	if err != nil {
		log.Println("NewSession")
		log.Fatal(err)
	}

	ce := costexplorer.New(sess)

	dateInterval := &costexplorer.DateInterval{}
	dateInterval.SetStart("2021-01-01")
	dateInterval.SetEnd("2021-01-31")

	filter := &costexplorer.Expression{}
	filter.SetDimensions(&costexplorer.DimensionValues{
		Key:          aws.String(costexplorer.DimensionLinkedAccount),
		Values:       []*string{aws.String(os.Getenv("AWS_ACC_ID"))},
	})

	resp, err := ce.GetCostAndUsage(&costexplorer.GetCostAndUsageInput{
		Metrics:       []*string{aws.String(costexplorer.MetricBlendedCost)},
		TimePeriod:    dateInterval,
		Granularity: aws.String(costexplorer.GranularityMonthly),
		Filter: filter,
	})
	if err != nil {
		log.Println("GetCostAndUsage")
		log.Fatal(err)
	}

	fmt.Println(resp.String())
}