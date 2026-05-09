package cron

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fazriegi/fintrack-be/internal/repository"
	"github.com/go-co-op/gocron"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type YahooFinanceResponse struct {
	OptionChain YahooFinanceOptionChain `json:"optionChain"`
}

type YahooFinanceOptionChain struct {
	Result []YahooFinanceResult `json:"result"`
}

type YahooFinanceResult struct {
	UnderlyingSymbol string            `json:"underlyingSymbol"`
	Quote            YahooFinanceQuote `json:"quote"`
}

type YahooFinanceQuote struct {
	RegularMarketPrice         decimal.Decimal `json:"regularMarketPrice"`
	RegularMarketChange        decimal.Decimal `json:"regularMarketChange"`
	RegularMarketChangePercent decimal.Decimal `json:"regularMarketChangePercent"`
}

func UpdateStockPrice(db *sqlx.DB, assetRepo repository.AssetRepository, appLogger *log.Logger) {
	s := gocron.NewScheduler(time.Local)

	_, err := s.Every(1).Day().At("17:00").Do(func() {
		appLogger.Println("Starting scheduled stock price update...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		tickers, err := assetRepo.GetTickers(ctx, db)
		if err != nil {
			appLogger.Printf("ERROR: Failed to get tickers: %v", err)
			return
		}

		var wg sync.WaitGroup
		errChan := make(chan error, len(*tickers))

		for _, ticker := range *tickers {
			wg.Add(1)

			go func() {
				defer wg.Done()
				url := fmt.Sprintf("https://yahoo-finance-real-time1.p.rapidapi.com/stock/get-options?symbol=%s.JK&lang=en-US", ticker)
				method := "GET"

				client := &http.Client{}
				req, err := http.NewRequest(method, url, nil)

				if err != nil {
					appLogger.Printf("ERROR: Failed to create stock price request: %v", err)
					errChan <- err
					return
				}
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("x-rapidapi-host", "yahoo-finance-real-time1.p.rapidapi.com")
				req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))

				res, err := client.Do(req)
				if err != nil {
					appLogger.Printf("ERROR: Failed to get stock price: %v", err)
					errChan <- err
					return
				}
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				if err != nil {
					appLogger.Printf("ERROR: Failed to read stock price response: %v", err)
					errChan <- err
					return
				}
				if res.StatusCode != 200 {
					errChan <- fmt.Errorf("stock price API returned status code %d", res.StatusCode)
					return
				}

				var yahooFinanceResp YahooFinanceResponse
				err = json.Unmarshal(body, &yahooFinanceResp)
				if err != nil {
					appLogger.Printf("ERROR: Failed to unmarshal stock price response: %v", err)
					errChan <- err
					return
				}

				if len(yahooFinanceResp.OptionChain.Result) == 0 {
					appLogger.Printf("WARNING: No result found for ticker %s", ticker)
					return
				}

				result := yahooFinanceResp.OptionChain.Result[0]
				if result.Quote.RegularMarketPrice.IsZero() {
					appLogger.Printf("WARNING: No market price found for ticker %s", ticker)
					return
				}

				err = assetRepo.UpdateStockPrice(ctx, db, ticker, result.Quote.RegularMarketPrice)
				if err != nil {
					appLogger.Printf("ERROR: Failed to update stock price: %v", err)
					errChan <- err
					return
				}
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				appLogger.Printf("ERROR: Failed to update stock prices: %v", err)
				return
			}
		}

		appLogger.Printf("SUCCESS: Stock prices updated at %s", time.Now().Format("2006-01-02 15:04:05"))
	})

	if err != nil {
		appLogger.Fatalf("Failed to schedule job: %v", err)
	}

	s.StartAsync()

	appLogger.Println("Stock price update scheduler is active.")
}
