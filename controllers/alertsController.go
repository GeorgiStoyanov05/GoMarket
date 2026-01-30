package controllers

import (
	"strconv"
	"net/http"
    "strings"
    "sort"
    "github.com/GeorgiStoyanov05/GoMarket/middlewares"
	"github.com/GeorgiStoyanov05/GoMarket/models"
	"github.com/GeorgiStoyanov05/GoMarket/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gin-gonic/gin"
)

type AlertGroup struct {
	Symbol string
	Alerts []models.PriceAlert
}

func PostCreateAlert(c *gin.Context) {
	symbol := c.Param("symbol")
	cond := strings.TrimSpace(c.PostForm("condition"))
	targetStr := strings.TrimSpace(c.PostForm("targetPrice"))

	uVal, ok := c.Get("user")
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}
	user, ok := uVal.(models.User)
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}

	target, err := strconv.ParseFloat(targetStr, 64)
	if targetStr == "" || err != nil {
		c.String(http.StatusOK, `<div class="text-danger">Please enter a valid target price.</div>`)
		return
	}

	_, errs := services.CreatePriceAlert(user.ID, symbol, cond, target)
	if len(errs) > 0 {
		msgs := make([]string, 0, len(errs))
		if v, ok := errs["_form"]; ok && v != "" {
			msgs = append(msgs, v)
		}
		for k, v := range errs {
			if k == "_form" {
				continue
			}
			if v != "" {
				msgs = append(msgs, v)
			}
		}
		c.String(http.StatusOK, `<div class="text-danger">`+strings.Join(msgs, "<br>")+`</div>`)
		return
	}

	c.Header("HX-Trigger", "alertsUpdated")

	c.String(http.StatusOK, `<div class="text-success">Alert created ✅</div>`)
}

func GetAlertsList(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))

	uVal, ok := c.Get("user")
	if !ok {
		c.HTML(http.StatusOK, "alertsList", middlewares.WithAuth(c, gin.H{
			"Symbol": symbol,
			"Alerts": []models.PriceAlert{},
		}))
		return
	}

	user := uVal.(models.User)

	alerts, err := services.ListPriceAlerts(user.ID, symbol)
	if err != nil {
		alerts = []models.PriceAlert{}
	}

	c.HTML(http.StatusOK, "alertsList", middlewares.WithAuth(c, gin.H{
		"Symbol": symbol,
		"Alerts": alerts,
	}))
}

// POST /alerts/:symbol/:id/delete
func PostDeleteAlert(c *gin.Context) {
	symbol := strings.ToUpper(strings.TrimSpace(c.Param("symbol")))
	idStr := strings.TrimSpace(c.Param("id"))

	uVal, ok := c.Get("user")
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}
	user, ok := uVal.(models.User)
	if !ok {
		c.String(http.StatusUnauthorized, `<div class="text-danger">Unauthorized</div>`)
		return
	}

	oid, err := primitive.ObjectIDFromHex(idStr)
	if err == nil {
		_ = services.DeletePriceAlert(user.ID, oid)
	}

	alerts, err := services.ListPriceAlerts(user.ID, symbol)
	if err != nil {
		alerts = []models.PriceAlert{}
	}

	c.HTML(http.StatusOK, "alertsList", middlewares.WithAuth(c, gin.H{
		"Symbol": symbol,
		"Alerts": alerts,
	}))
}

// GET /alerts/list
func GetWatchlistAlerts(c *gin.Context) {
	uVal, ok := c.Get("user")
	if !ok {
		c.HTML(http.StatusOK, "watchlistAlerts", middlewares.WithAuth(c, gin.H{
			"Groups": []AlertGroup{},
		}))
		return
	}
	user := uVal.(models.User)

	alerts, err := services.ListAllUserAlerts(user.ID)
	if err != nil {
		alerts = []models.PriceAlert{}
	}

	// group by symbol
	m := map[string][]models.PriceAlert{}
	for _, a := range alerts {
		s := strings.ToUpper(strings.TrimSpace(a.Symbol))
		m[s] = append(m[s], a)
	}

	groups := make([]AlertGroup, 0, len(m))
	for sym, list := range m {
		groups = append(groups, AlertGroup{Symbol: sym, Alerts: list})
	}

	// stable order
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Symbol < groups[j].Symbol
	})

	c.HTML(http.StatusOK, "watchlistAlerts", middlewares.WithAuth(c, gin.H{
		"Groups": groups,
	}))
}

func PostDeleteAlertGlobal(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))

	uVal, ok := c.Get("user")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}
	user := uVal.(models.User)

	oid, err := primitive.ObjectIDFromHex(idStr)
	if err == nil {
		_ = services.DeletePriceAlert(user.ID, oid)
	}

	// Make both Details + Watchlist refresh wherever they are
	c.Header("HX-Trigger", "alertsUpdated")

	// nothing to swap (watchlist uses hx-swap="none")
	c.Status(http.StatusNoContent)
}
