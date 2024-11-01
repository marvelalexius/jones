package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marvelalexius/jones/model"
	"github.com/marvelalexius/jones/utils"
	"github.com/marvelalexius/jones/utils/logger"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

func (h *HTTPService) Subscribe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.Errorln(c, "failed to get user id from context")
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when subscribing",
		})

		return
	}

	var req model.SubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Errorln(c, "failed to bind json", err)
		ve := utils.ValidationResponse(err)
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrorRes{
			Message: "something went wrong when validating the requests",
			Errors:  ve,
		})

		return
	}

	subscription, err := h.SubscriptionService.Subscribe(c, userID.(string), req)
	if err != nil {
		logger.Errorln(c, "failed to subscribe", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when subscribing",
			Errors:  err,
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "success",
		Data:    subscription,
	})
}

func (h *HTTPService) ManageSubscription(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		logger.Errorln(c, "failed to get user id from context")
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when managing subscription",
		})

		return
	}

	portal, err := h.SubscriptionService.CustomerPortal(c, userID.(string))
	if err != nil {
		logger.Errorln(c, "failed to manage subscription", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when managing subscription",
			Errors:  err,
		})

		return
	}

	utils.SuccessResponse(c, http.StatusOK, utils.SuccessRes{
		Message: "success",
		Data:    portal,
	})
}

func (h *HTTPService) HandleCallback(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Errorln(c, "failed to read request body", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when reading request body",
			Errors:  err,
		})

		return
	}

	event, err := webhook.ConstructEvent(body, c.Request.Header.Get("Stripe-Signature"), h.Conf.Stripe.WebhookSecret)
	if err != nil {
		logger.Errorln(c, "failed to construct event", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
			Message: "something went wrong when constructing event",
			Errors:  err,
		})

		return
	}

	switch event.Type {
	case "customer.subscription.deleted":
		var subscription stripe.Subscription

		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			logger.Errorln(c, "failed to unmarshal subscription", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
				Message: "something went wrong when unmarshaling subscription",
				Errors:  err,
			})

			return
		}

		h.SubscriptionService.HandleSubscriptionDeleted(c, &subscription)
	case "invoice.paid":
		var invoice stripe.Invoice

		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			logger.Errorln(c, "failed to unmarshal invoice", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
				Message: "something went wrong when unmarshaling invoice",
				Errors:  err,
			})

			return
		}

		h.SubscriptionService.HandleInvoicePaid(c, &invoice)
	case "invoice.payment_failed":
		var invoice stripe.Invoice

		if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
			logger.Errorln(c, "failed to unmarshal invoice", err)
			utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrorRes{
				Message: "something went wrong when unmarshaling invoice",
				Errors:  err,
			})

			return
		}

		h.SubscriptionService.HandleInvoicePaymentFailed(c, invoice.CustomerEmail)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
