package pharmacy_orders

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc gin.HandlerFunc
}

func NewRouter(handleFunctions ApiHandleFunctions) *gin.Engine {
	return NewRouterWithGinEngine(gin.Default(), handleFunctions)
}

func NewRouterWithGinEngine(router *gin.Engine, handleFunctions ApiHandleFunctions) *gin.Engine {
	for _, route := range getRoutes(handleFunctions) {
		if route.HandlerFunc == nil {
			route.HandlerFunc = DefaultHandleFunc
		}
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodPatch:
			router.PATCH(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	return router
}

func DefaultHandleFunc(c *gin.Context) {
	c.String(http.StatusNotImplemented, "501 not implemented")
}

type ApiHandleFunctions struct {
	PharmacyOrdersAPI      PharmacyOrdersAPI
	PharmacyMedicinesAPI   PharmacyMedicinesAPI
	PharmacyDispensingsAPI PharmacyDispensingsAPI
}

func getRoutes(handleFunctions ApiHandleFunctions) []Route {
	return []Route{
		{
			"GetMedicines",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/medicines",
			handleFunctions.PharmacyMedicinesAPI.GetMedicines,
		},
		{
			"CreateMedicine",
			http.MethodPost,
			"/api/pharmacy/:pharmacyId/medicines",
			handleFunctions.PharmacyMedicinesAPI.CreateMedicine,
		},
		{
			"GetMedicine",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/medicines/:medicineId",
			handleFunctions.PharmacyMedicinesAPI.GetMedicine,
		},
		{
			"UpdateMedicine",
			http.MethodPut,
			"/api/pharmacy/:pharmacyId/medicines/:medicineId",
			handleFunctions.PharmacyMedicinesAPI.UpdateMedicine,
		},
		{
			"DeleteMedicine",
			http.MethodDelete,
			"/api/pharmacy/:pharmacyId/medicines/:medicineId",
			handleFunctions.PharmacyMedicinesAPI.DeleteMedicine,
		},
		{
			"CreateOrder",
			http.MethodPost,
			"/api/pharmacy/:pharmacyId/orders",
			handleFunctions.PharmacyOrdersAPI.CreateOrder,
		},
		{
			"DeleteOrder",
			http.MethodDelete,
			"/api/pharmacy/:pharmacyId/orders/:orderId",
			handleFunctions.PharmacyOrdersAPI.DeleteOrder,
		},
		{
			"GetOrder",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/orders/:orderId",
			handleFunctions.PharmacyOrdersAPI.GetOrder,
		},
		{
			"GetOrders",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/orders",
			handleFunctions.PharmacyOrdersAPI.GetOrders,
		},
		{
			"UpdateOrder",
			http.MethodPut,
			"/api/pharmacy/:pharmacyId/orders/:orderId",
			handleFunctions.PharmacyOrdersAPI.UpdateOrder,
		},
		{
			"UpdateOrderStatus",
			http.MethodPatch,
			"/api/pharmacy/:pharmacyId/orders/:orderId/status",
			handleFunctions.PharmacyOrdersAPI.UpdateOrderStatus,
		},
		{
			"GetDispensings",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/dispensings",
			handleFunctions.PharmacyDispensingsAPI.GetDispensings,
		},
		{
			"CreateDispensing",
			http.MethodPost,
			"/api/pharmacy/:pharmacyId/dispensings",
			handleFunctions.PharmacyDispensingsAPI.CreateDispensing,
		},
		{
			"GetOrderReceipt",
			http.MethodGet,
			"/api/pharmacy/:pharmacyId/orders/:orderId/receipt",
			handleFunctions.PharmacyOrdersAPI.GetOrderReceipt,
		},
		{
			"ReceiveOrder",
			http.MethodPost,
			"/api/pharmacy/:pharmacyId/orders/:orderId/receive",
			handleFunctions.PharmacyOrdersAPI.ReceiveOrder,
		},
	}
}
