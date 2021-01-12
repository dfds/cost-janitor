package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/coreos/go-oidc"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var LISTEN_ADDRESS = os.Getenv("COST_JANITOR_LISTEN_ADDRESS")
var BASIC_VALUE = os.Getenv("COST_JANITOR_BASIC_VALUE")
var redis_ctx = context.Background()
var rdb *redis.Client

func main() {
	fmt.Println("Launching cost-janitor")

	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	provider, err := oidc.NewProvider(context.Background(), "https://login.microsoftonline.com/73a99466-ad05-4221-9f90-e7142aa2f6c1/v2.0")
	if err != nil {
		log.Fatal(err)
	}

	authMiddleware := authenticationMiddleware{
		ClientID: "24420be9-46e5-4584-acd7-64850d2f2a03",
		Provider: provider,
	}

	r := mux.NewRouter()
	r.Handle("/get-monthly-total-cost/{accountid}", authMiddleware.Middleware(http.HandlerFunc(GetMonthlyTotalCost)))
	r.Handle("/basic/get-monthly-total-cost/{accountid}", BasicAuthMiddleware(http.HandlerFunc(GetMonthlyTotalCost)))

	addr := fmt.Sprintf("%s:8080", LISTEN_ADDRESS)
	fmt.Printf("HTTP server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, handlers.CompressHandler(r))); err != nil {
		log.Fatal(err)
	}
}

func getCurrentFullMonthDateRange() (string, string) {
	now := time.Now()
	year, month, _ := now.Date()
	endOfThisMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, now.Location())
	monthNumerical := fmt.Sprintf("%v", int(now.Month()))
	if int(now.Month()) <= 9 {
		monthNumerical = fmt.Sprintf("0%v", int(now.Month()))
	}

	return fmt.Sprintf("%v-%v-%v", year, monthNumerical, "01"), fmt.Sprintf("%v-%v-%v", year, monthNumerical, endOfThisMonth.Day())
}

func GetMonthlyTotalCost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	redisKey := fmt.Sprintf("currentmonth.acc.%s", vars["accountid"])


	val, err := rdb.Get(redis_ctx, redisKey).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("No cached result, querying AWS")
	case err != nil:
		log.Fatal("Get failed: ", err)
	}

	if val != "" {
		fmt.Println("Cached entry found, using for response.")
		w.WriteHeader(200)
		w.Write([]byte(val))
		return
	}

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
	startOfMonth, endOfMonth := getCurrentFullMonthDateRange()
	dateInterval.SetStart(startOfMonth)
	dateInterval.SetEnd(endOfMonth)

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

	err = rdb.Set(redis_ctx, redisKey, resp.String(), time.Hour).Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response queried for account %s is now cached for the next hour\n", vars["accountid"])

	w.WriteHeader(200)
	w.Write([]byte(resp.String()))
}


type authenticationMiddleware struct {
	ClientID string
	Provider *oidc.Provider
}

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Basic")
		if len(splitToken) != 2 {
			http.Error(w, "Basic value doesn't seem right", http.StatusUnauthorized)
			return
		}

		basicValue := strings.TrimSpace(splitToken[1])
		decoded, err := base64.URLEncoding.DecodeString(basicValue)
		if err != nil {
			http.Error(w, "Unable to decode basic value", http.StatusUnauthorized)
			return
		}
		cred := string(decoded)

		if cred != BASIC_VALUE {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
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