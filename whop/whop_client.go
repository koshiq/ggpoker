package whop

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	WHOP_API_BASE = "https://api.whop.com/api/v2"
)

type WhopClient struct {
	APIKey     string
	HTTPClient *http.Client
}

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	ProfilePicURL string    `json:"profile_pic_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Subscription struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	ProductID   string     `json:"product_id"`
	Status      string     `json:"status"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     time.Time  `json:"end_date"`
	RenewalDate time.Time  `json:"renewal_date"`
	CanceledAt  *time.Time `json:"canceled_at"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Interval    string  `json:"interval"`
}

type AccessCheck struct {
	HasAccess bool   `json:"has_access"`
	Reason    string `json:"reason,omitempty"`
}

func NewWhopClient(apiKey string) *WhopClient {
	return &WhopClient{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// GetUser retrieves user information from WHOP
func (wc *WhopClient) GetUser(userID string) (*User, error) {
	url := fmt.Sprintf("%s/users/%s", WHOP_API_BASE, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+wc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (wc *WhopClient) GetUserSubscriptions(userID string) ([]Subscription, error) {
	url := fmt.Sprintf("%s/users/%s/subscriptions", WHOP_API_BASE, userID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+wc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var subscriptions []Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return subscriptions, nil
}

// CheckAccess verifies if a user has access to a specific product
func (wc *WhopClient) CheckAccess(userID, productID string) (*AccessCheck, error) {
	url := fmt.Sprintf("%s/users/%s/access/%s", WHOP_API_BASE, userID, productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+wc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var accessCheck AccessCheck
	if err := json.NewDecoder(resp.Body).Decode(&accessCheck); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &accessCheck, nil
}

// GetProduct retrieves product information
func (wc *WhopClient) GetProduct(productID string) (*Product, error) {
	url := fmt.Sprintf("%s/products/%s", WHOP_API_BASE, productID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+wc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &product, nil
}

// ValidateUserToken validates a user's access token and returns user info
func (wc *WhopClient) ValidateUserToken(token string) (*User, error) {
	url := fmt.Sprintf("%s/me", WHOP_API_BASE)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wc.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &user, nil
}

// IsUserSubscribed checks if a user has an active subscription to a product
func (wc *WhopClient) IsUserSubscribed(userID, productID string) (bool, error) {
	subscriptions, err := wc.GetUserSubscriptions(userID)
	if err != nil {
		return false, err
	}

	for _, sub := range subscriptions {
		if sub.ProductID == productID && sub.Status == "active" {
			return true, nil
		}
	}

	return false, nil
}
