package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/coreos/go-oidc"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	fmt.Println("Launching cost-janitor")

	provider, err := oidc.NewProvider(context.Background(), "https://login.microsoftonline.com/73a99466-ad05-4221-9f90-e7142aa2f6c1/v2.0")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := authenticationMiddleware{
		ClientID: "24420be9-46e5-4584-acd7-64850d2f2a03",
		Provider: provider,
	}

	r := mux.NewRouter()
	r.HandleFunc("/get-monthly-total-cost/{accountid}", GetMonthlyTotalCost)

	println("HTTP server listening on :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(authMiddleware.Middleware(r)))); err != nil {
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


type authenticationMiddleware struct {
	ClientID string
	Provider *oidc.Provider
}

func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	var verifier = amw.Provider.Verifier(&oidc.Config{ClientID: amw.ClientID})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqToken := r.Header.Get("Authorization") //Authorization: Bearer a7ydfs87afasd8f990
		splitToken := strings.Split(reqToken, "Bearer")
		if len(splitToken) != 2 {
			http.Error(w, "Token doesn't seem right", http.StatusUnauthorized)
			return
		}

		reqToken = strings.TrimSpace(splitToken[1])

		idToken, err := verifier.Verify(r.Context(), reqToken)
		if err != nil {
			http.Error(w, "Unable to verify token", http.StatusUnauthorized)
			return
		}

		var claims struct {
			Emails []string `json:"emails"`
		}
		if err := idToken.Claims(&claims); err != nil {
			fmt.Println(err)
			http.Error(w, "Unable to retrieve claims", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}