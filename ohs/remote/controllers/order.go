package controllers

import (
	"fmt"
	"net/http"
	"order-context/domain/aggregate"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/ohs/local/services/ohttp"
	"order-context/utils/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

// OrderGroup ...
func OrderGroup(r *gin.RouterGroup) {
	r.GET("/orders", GetOrders)
	r.POST("/orders", CreateOrder)
}

// GetOrders ...
func GetOrders(c *gin.Context) {
	offset, _ := strconv.Atoi(c.Query("offset"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	params := pl.ListOrderParams{
		SpaceID: c.Query("space_id"),
		Offset:  offset,
		Limit:   limit,
	}
	datas, total, err := ohttp.GetOrderListAppService(params)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err)
	}

	c.JSON(http.StatusOK, gin.H{
		"datas": datas,
		"total": total,
	})
}

// CreateOrder ...
func CreateOrder(c *gin.Context) {
	var params pl.CreateOrderParams
	if err := c.ShouldBindJSON(&params); err != nil {
		errResponse(c, http.StatusBadRequest, errors.BadRequest(fmt.Sprintf("invalid body params, %v", err)))
		return
	}

	siteCode := c.Request.Header.Get("site-code")
	orderOption := aggregate.WithOrderOption(common.StatusType(params.Status), params.Price)
	spaceOption := aggregate.WithSpaceOption(params.SpaceID)
	packageOption := aggregate.WithPackageOption(params.PackageVersion, params.PackagePrice)

	res, err := ohttp.CreateOrderAppService("", siteCode, orderOption, spaceOption, packageOption)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusCreated, res)

}

func errResponse(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}
