package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/wildanasyrof/backend-topup/internal/domain/dto"
	"github.com/wildanasyrof/backend-topup/internal/domain/entity"
	"github.com/wildanasyrof/backend-topup/internal/repository"
	"github.com/wildanasyrof/backend-topup/pkg/logger"
	"gorm.io/gorm"
)

// Define constants for external service configuration
const (
	DFBaseURL  = "http://localhost:3000" // Should ideally be configured externally (e.g., environment variable)
	DFEndpoint = "/v1/price-list?brand=Mobile%20Legends"
	DFUsername = "demo_buyer"
	DFAPIKey   = "secret_apikey_demo" // Should *never* be hardcoded in production code! Use a secret manager.
	DFCommand  = "prepaid"
)

// ExternalService defines the interface for external API interactions.
type ExternalService interface {
	DFGetProductList(ctx context.Context) ([]dto.DFProductListRes, error)
	DFSaveProductList(ctx context.Context) ([]dto.DFProductListRes, error)
}

// externalService holds the dependencies for interacting with the external API.
type externalService struct {
	httpClient  *http.Client
	logger      logger.Logger
	productRepo repository.ProductRepository
	// Consider adding base URL/credentials here if they vary or need to be configured at startup.
}

// NewExternalService is the constructor for externalService.
func NewExternalService(httpClient *http.Client, logger logger.Logger, productRepo repository.ProductRepository) ExternalService {
	return &externalService{httpClient: httpClient, logger: logger, productRepo: productRepo}
}

// makeSign generates the MD5 hash for API authentication.
func makeSign(username, apiKey string) string {
	// Concatenate the required components for the hash
	dataToHash := username + apiKey + "pricelist"
	hash := md5.Sum([]byte(dataToHash))
	return hex.EncodeToString(hash[:])
}

// DFProductList fetches the product list from the external service.
func (e *externalService) DFGetProductList(ctx context.Context) ([]dto.DFProductListRes, error) {
	// 1. Prepare Request Body
	reqBody := dto.DFProductListReq{
		Cmd:      DFCommand,
		Username: DFUsername,
		Sign:     makeSign(DFUsername, DFAPIKey),
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request error: %w", err)
	}

	// 2. Create HTTP Request
	url := DFBaseURL + DFEndpoint
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// 3. Execute HTTP Request
	res, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http execution error: %w", err)
	}
	defer func() {
		// Ensure the response body is closed to prevent resource leaks
		if closeErr := res.Body.Close(); closeErr != nil {
			e.logger.Warn(fmt.Sprintf("failed to close response body: %v", closeErr))
		}
	}()

	// 4. Read Response Body
	rawBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}

	// 5. Handle Non-2xx Status Codes
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		// Return a specific error indicating a bad response
		return nil, fmt.Errorf("external service bad response: status %d, body: %s", res.StatusCode, rawBody)
	}

	// 6. Unmarshal Response
	var payload dto.DFBaseRes
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal response error: %w", err)
	}

	return payload.Data, nil
}

// DFGetProductList implements ExternalService.
func (e *externalService) DFSaveProductList(ctx context.Context) ([]dto.DFProductListRes, error) {
	products, err := e.DFGetProductList(ctx)

	if err != nil {
		return nil, err
	}

	for _, it := range products {

		product, err := e.productRepo.FindByCode(ctx, it.BuyerSkuCode)
		e.logger.Error(err, "wtf")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			product := &entity.Product{
				Name:        it.ProductName,
				SkuCode:     it.BuyerSkuCode,
				SellerName:  it.SellerName,
				CategoryID:  12,
				ProviderID:  5,
				Status:      StatusMapper(it.BuyerProductStatus),
				Stock:       StockMapper(it.Stock, it.UnlimitedStock),
				BasePrice:   float64(it.Price),
				Description: it.Desc,
				ImgUrl:      "",
				StartOff:    it.StartCutOff,
				EndOff:      it.EndCutOff,
			}

			if err := e.productRepo.Create(ctx, product); err != nil {
				return nil, err
			}
		}

		product.BasePrice = float64(it.Price)
		product.Stock = StockMapper(it.Stock, it.UnlimitedStock)
		product.Status = StatusMapper(it.BuyerProductStatus)
		product.StartOff = it.StartCutOff
		product.EndOff = it.EndCutOff

		if err := e.productRepo.Update(ctx, product); err != nil {
			return nil, err
		}
	}

	return products, err
}

func StatusMapper(status bool) entity.CatStatus {
	if status {
		return entity.CatActive
	}

	return entity.CatInactive
}

func StockMapper(stock int32, isUnlimited bool) int64 {
	if isUnlimited {
		return 9999999
	}

	return int64(stock)
}
