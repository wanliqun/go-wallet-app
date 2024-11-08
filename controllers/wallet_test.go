package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wanliqun/go-wallet-app/controllers"
	"github.com/wanliqun/go-wallet-app/middlewares"
	"github.com/wanliqun/go-wallet-app/mocks"
	"github.com/wanliqun/go-wallet-app/models"
	"github.com/wanliqun/go-wallet-app/services"
)

var (
	userGenerator models.FakeUserGenerator
)

func setupTestRouter(walletService *mocks.MockWalletService, userService *mocks.MockUserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(middlewares.AuthMiddleware(userService))
	router.Use(middlewares.CorsMiddleware())

	walletController := controllers.NewWalletController(walletService, userService)
	walletRouter := router.Group("/")
	{
		walletRouter.POST("/deposit", walletController.Deposit)
		walletRouter.POST("/withdraw", walletController.Withdraw)
		walletRouter.POST("/transfer", walletController.Transfer)
		walletRouter.GET("/balances", walletController.GetBalances)
		walletRouter.GET("/transactions", walletController.GetTransactionHistory)
	}

	return router
}

func TestWalletController_Deposit(t *testing.T) {
	mockWalletService := new(mocks.MockWalletService)
	mockUserService := new(mocks.MockUserService)

	router := setupTestRouter(mockWalletService, mockUserService)

	testUser := userGenerator.Generate()
	currency := "USDT"

	t.Run("should deposit successfully", func(t *testing.T) {
		amount := decimal.NewFromFloat(100.0)

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("Deposit", testUser.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		})).Return(nil)

		body, _ := json.Marshal(map[string]interface{}{
			"currency": currency,
			"amount":   amount.String(),
		})
		req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response Body: %s", w.Body.String())
	})

	t.Run("should return error for invalid amount", func(t *testing.T) {
		amount := decimal.NewFromFloat(-100.0)

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("Deposit", testUser.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		})).Return(services.ErrInvalidAmount)

		body, _ := json.Marshal(map[string]interface{}{
			"currency": currency,
			"amount":   amount.String(),
		})
		req, _ := http.NewRequest("POST", "/deposit", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["message"], "Field validation for 'Amount' failed on the 'positive_decimal' tag")
	})
}

func TestWalletController_Withdraw(t *testing.T) {
	mockWalletService := new(mocks.MockWalletService)
	mockUserService := new(mocks.MockUserService)
	router := setupTestRouter(mockWalletService, mockUserService)

	testUser := userGenerator.Generate()
	currency := "USDT"

	t.Run("should withdraw successfully", func(t *testing.T) {
		amount := decimal.NewFromFloat(100.0)

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("Withdraw", testUser.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		})).Return(nil)

		body, _ := json.Marshal(map[string]interface{}{
			"currency": currency,
			"amount":   amount.String(),
		})
		req, _ := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response Body: %s", w.Body.String())
	})

	t.Run("should return error for insufficient balance", func(t *testing.T) {
		amount := decimal.NewFromFloat(200.0)

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("Withdraw", testUser.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		})).Return(services.ErrInsufficientBalance)

		body, _ := json.Marshal(map[string]interface{}{
			"currency": currency,
			"amount":   amount.String(),
		})
		req, _ := http.NewRequest("POST", "/withdraw", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["message"], services.ErrInsufficientBalance.Error())
	})
}

func TestWalletController_Transfer(t *testing.T) {
	mockWalletService := new(mocks.MockWalletService)
	mockUserService := new(mocks.MockUserService)
	router := setupTestRouter(mockWalletService, mockUserService)

	sender := userGenerator.Generate()
	recipient := userGenerator.Generate()
	currency := "USDT"
	memo := "test transfer"

	t.Run("should transfer successfully", func(t *testing.T) {
		amount := decimal.NewFromFloat(30.0)

		mockUserService.On("GetUserByName", sender.Name).Return(sender, true, nil)
		mockUserService.On("GetUserByName", recipient.Name).Return(recipient, true, nil)
		mockWalletService.On("Transfer", sender.ID, recipient.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		}), memo).Return(nil)

		reqBody, _ := json.Marshal(map[string]interface{}{
			"recipient": recipient.Name, "currency": currency, "amount": amount.String(), "memo": memo,
		})
		req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+sender.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response Body: %s", w.Body.String())
	})

	t.Run("should return error for insufficient balance", func(t *testing.T) {
		amount := decimal.NewFromFloat(100.0)

		mockUserService.On("GetUserByName", sender.Name).Return(sender, true, nil)
		mockUserService.On("GetUserByName", recipient.Name).Return(recipient, true, nil)
		mockWalletService.On("Transfer", sender.ID, recipient.ID, currency, mock.MatchedBy(func(a decimal.Decimal) bool {
			return a.Equal(amount)
		}), memo).Return(services.ErrInsufficientBalance)

		reqBody, _ := json.Marshal(map[string]interface{}{
			"recipient": recipient.Name, "currency": currency, "amount": amount.String(), "memo": memo,
		})
		req, _ := http.NewRequest("POST", "/transfer", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+sender.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Contains(t, resp["message"], services.ErrInsufficientBalance.Error())
	})
}

func TestWalletController_Transfer_GetBalances(t *testing.T) {
	mockWalletService := new(mocks.MockWalletService)
	mockUserService := new(mocks.MockUserService)
	router := setupTestRouter(mockWalletService, mockUserService)

	testUser := userGenerator.Generate()
	currencies := []string{"USDT", "BTC"}

	t.Run("should get balances successfully", func(t *testing.T) {
		expectedVaults := []models.Vault{
			{UserID: testUser.ID, Currency: "USDT", Amount: decimal.NewFromFloat(100.0)},
			{UserID: testUser.ID, Currency: "BTC", Amount: decimal.NewFromFloat(1)},
		}

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("GetBalances", testUser.ID, currencies).Return(expectedVaults, nil)

		req, _ := http.NewRequest("GET", "/balances?currency=USDT&currency=BTC", nil)
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		var resp map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, len(expectedVaults), len(resp["data"].([]interface{})))
	})
}

func TestWalletController_GetTransactions(t *testing.T) {
	mockWalletService := new(mocks.MockWalletService)
	mockUserService := new(mocks.MockUserService)
	router := setupTestRouter(mockWalletService, mockUserService)

	testUser := userGenerator.Generate()

	t.Run("should get transactions successfully", func(t *testing.T) {
		txnType := models.Deposit
		cursor := "cursor"
		expectedCursor := "next_cursor"
		expectedTransactions := []models.Transaction{
			{UserID: testUser.ID, Type: "deposit", Currency: "USDT", Amount: decimal.NewFromFloat(100.0)},
			{UserID: testUser.ID, Type: "withdraw", Currency: "BTC", Amount: decimal.NewFromFloat(1)},
		}

		mockUserService.On("GetUserByName", testUser.Name).Return(testUser, true, nil)
		mockWalletService.On("GetTransactionHistory", testUser.ID, txnType, cursor, services.SortOrderDesc, 10).
			Return(expectedTransactions, expectedCursor, nil)

		req, _ := http.NewRequest("GET", "/transactions?type=deposit&cursor=cursor&order=desc&limit=10", nil)
		req.Header.Set("Authorization", "Bearer "+testUser.Name)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		var resp struct {
			Code    int
			Message string
			Data    controllers.GetTransactionHistoryResponse
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, expectedCursor, resp.Data.NextCursor)
		assert.Equal(t, len(expectedTransactions), len(resp.Data.Transactions))
	})
}
