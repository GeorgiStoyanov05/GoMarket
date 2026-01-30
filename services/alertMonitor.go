package services

import (
	"context"
	"log"
	"time"

	"github.com/GeorgiStoyanov05/GoMarket2/models"
)

func StartPriceAlertMonitor(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runAlertTick()
			}
		}
	}()
}

func runAlertTick() {
	alerts, err := ListActiveAlerts()
	if err != nil || len(alerts) == 0 {
		return
	}

	// Group by symbol so we fetch 1 quote per symbol per tick.
	bySymbol := map[string][]models.PriceAlert{}
	for _, a := range alerts {
		bySymbol[a.Symbol] = append(bySymbol[a.Symbol], a)
	}

	for sym, group := range bySymbol {
		price, err := FetchCurrentPrice(sym)
		if err != nil {
			continue
		}

		for _, a := range group {
			hit := (a.Condition == "above" && price >= a.TargetPrice) ||
				(a.Condition == "below" && price <= a.TargetPrice)

			if !hit {
				continue
			}

			if err := MarkAlertTriggered(a.ID, price); err != nil {
				log.Println("alert monitor: mark triggered:", err)
			}
		}
	}
}
