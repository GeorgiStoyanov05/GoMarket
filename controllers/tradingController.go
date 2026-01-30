package controllers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/GeorgiStoyanov05/GoMarket2/middlewares"
	"github.com/GeorgiStoyanov05/GoMarket2/models"
	"github.com/GeorgiStoyanov05/GoMarket2/services"
	"github.com/gin-gonic/gin"
)

type PortfolioGroup struct {
	Symbol       string
	Key          string
	Qty          int64
	AvgCost      float64
	CurrentPrice float64
	PnL          float64
	PnLPct       float64
}

func safeKey(sym string) string {
	r := strings.NewReplacer(".", "_", ":", "_", "/", "_", " ", "_", "-", "_")
	return r.Replace(sym)
}

func PostMarketBuy(c *gin.Context) {
	symbol := c.Param("symbol")

	uVal, ok := c.Get("user")
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}
	user := uVal.(models.User)

	qtyStr := strings.TrimSpace(c.PostForm("qty"))
	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if qtyStr == "" || err != nil {
		c.String(http.StatusOK, `<div class="text-danger">Enter a valid quantity.</div>`)
		return
	}

	res, errs := services.MarketBuy(user.ID, symbol, qty)
	if len(errs) > 0 {
		// show first useful error
		if v, ok := errs["balance"]; ok {
			c.String(http.StatusOK, `<div class="text-danger">`+v+`</div>`)
			return
		}
		if v, ok := errs["qty"]; ok {
			c.String(http.StatusOK, `<div class="text-danger">`+v+`</div>`)
			return
		}
		if v, ok := errs["_form"]; ok {
			c.String(http.StatusOK, `<div class="text-danger">`+v+`</div>`)
			return
		}
		c.String(http.StatusOK, `<div class="text-danger">Could not buy.</div>`)
		return
	}

	// refresh position panel on the page + portfolio later
	c.Header("HX-Trigger", "positionUpdated")

	c.String(http.StatusOK,
		`<div class="text-success">Bought `+strconv.FormatInt(res.Qty, 10)+` `+res.Symbol+
			` @ `+format2(res.FillPrice)+
			` (Cost: `+format2(res.Cost)+
			`, New balance: `+format2(res.NewBalance)+`)</div>`)
}

func GetPositionPanel(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))

	uVal, ok := c.Get("user")
	if !ok {
		c.HTML(http.StatusOK, "positionPanel", middlewares.WithAuth(c, gin.H{
			"Symbol":      symbol,
			"HasPosition": false,
		}))
		return
	}
	user := uVal.(models.User)

	pos, err := services.GetUserPosition(user.ID, symbol)
	if err != nil || pos == nil || pos.Qty <= 0 {
		c.HTML(http.StatusOK, "positionPanel", middlewares.WithAuth(c, gin.H{
			"Symbol":      symbol,
			"HasPosition": false,
		}))
		return
	}

	price, err := services.FetchCurrentPrice(symbol)
	if err != nil || price <= 0 {
		price = pos.AvgCost
	}
	price = math.Round(price*100) / 100

	pnl := (price - pos.AvgCost) * float64(pos.Qty)
	pnl = math.Round(pnl*100) / 100

	pct := 0.0
	if pos.AvgCost > 0 {
		pct = (price - pos.AvgCost) / pos.AvgCost * 100.0
		pct = math.Round(pct*100) / 100
	}

	c.HTML(http.StatusOK, "positionPanel", middlewares.WithAuth(c, gin.H{
		"Symbol":       symbol,
		"HasPosition":  true,
		"Position":     pos,
		"CurrentPrice": price,
		"PnL":          pnl,
		"PnLPct":       pct,
	}))
}

func format2(v float64) string {
	return strconv.FormatFloat(v, 'f', 2, 64)
}

func PostMarketSell(c *gin.Context) {
	symbol := c.Param("symbol")

	uVal, ok := c.Get("user")
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}
	user := uVal.(models.User)

	qtyStr := strings.TrimSpace(c.PostForm("qty"))
	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if qtyStr == "" || err != nil {
		c.String(http.StatusOK, `<div class="text-danger">Enter a valid quantity.</div>`)
		return
	}

	res, errs := services.MarketSell(user.ID, symbol, qty)
	if len(errs) > 0 {
		if v, ok := errs["qty"]; ok {
			c.String(http.StatusOK, `<div class="text-danger">`+v+`</div>`)
			return
		}
		if v, ok := errs["_form"]; ok {
			c.String(http.StatusOK, `<div class="text-danger">`+v+`</div>`)
			return
		}
		c.String(http.StatusOK, `<div class="text-danger">Could not sell.</div>`)
		return
	}

	c.Header("HX-Trigger", "positionUpdated")

	c.String(http.StatusOK,
		`<div class="text-success">Sold `+strconv.FormatInt(res.Qty, 10)+` `+res.Symbol+
			` @ `+format2(res.FillPrice)+
			` (Proceeds: `+format2(res.Proceeds)+
			`, New balance: `+format2(res.NewBalance)+`)</div>`)
}

// GET /portfolio (SSR page)
func GetPortfolioPage(c *gin.Context) {
	if c.GetHeader("HX-Request") == "true" {
		c.HTML(200, "portfolio", middlewares.WithAuth(c, gin.H{}))
		return
	}
	c.HTML(200, "index.html", middlewares.WithAuth(c, gin.H{
		"InitialPath": "/portfolio",
	}))
}

// GET /portfolio/positions (HTMX partial)
func GetPortfolioPositions(c *gin.Context) {
	uVal, ok := c.Get("user")
	if !ok {
		c.HTML(http.StatusOK, "portfolioPositions", middlewares.WithAuth(c, gin.H{"Groups": []PortfolioGroup{}}))
		return
	}
	user := uVal.(models.User)

	positions, err := services.ListUserPositions(user.ID)
	if err != nil {
		positions = []models.Position{}
	}

	groups := make([]PortfolioGroup, 0, len(positions))
	for _, p := range positions {
		price, err := services.FetchCurrentPrice(p.Symbol)
		if err != nil || price <= 0 {
			price = p.AvgCost
		}
		price = math.Round(price*100) / 100

		pnl := (price - p.AvgCost) * float64(p.Qty)
		pnl = math.Round(pnl*100) / 100

		pct := 0.0
		if p.AvgCost > 0 {
			pct = (price - p.AvgCost) / p.AvgCost * 100.0
			pct = math.Round(pct*100) / 100
		}

		groups = append(groups, PortfolioGroup{
			Symbol:       p.Symbol,
			Key:          safeKey(p.Symbol),
			Qty:          p.Qty,
			AvgCost:      p.AvgCost,
			CurrentPrice: price,
			PnL:          pnl,
			PnLPct:       pct,
		})
	}

	c.HTML(http.StatusOK, "portfolioPositions", middlewares.WithAuth(c, gin.H{
		"Groups": groups,
	}))
}
