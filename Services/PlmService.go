package Services

import (
	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"errors"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PlatformService struct {
	Repo  *Repositories.PlatformRepository
	Audit *helper.AuditService
}

// Create a new Platform
func (s *PlatformService) Create(
	name string,
	email string,
	status string,
	ipAllowlist []string,
	apiKey string,
	secretKey string,
) (*models.Platform, error) {

	existing := &models.Platform{}

	err := s.Repo.DB.
		Where("email = ?", email).
		First(existing).Error

	if err == nil {
		return nil, errors.New("email already exists")
	}
	if ipAllowlist == nil {
		ipAllowlist = []string{}
	}
	platform := &models.Platform{
		Name:        name,
		Email:       email,
		Status:      status,
		IPAllowlist: ipAllowlist,
	}

	api := &models.PlatformApiKey{
		Key:    apiKey,
		Secret: secretKey,
		Status: "active",
	}

	if err := s.Repo.CreatePlatformWithApiKey(platform, api); err != nil {
		if strings.Contains(err.Error(), "uni_platform_api_keys_secret") {
			return nil, errors.New("secret key already exists")
		}
		if strings.Contains(err.Error(), "uni_platform_api_keys") {
			return nil, errors.New("api key already exists")
		}
		return nil, err
	}

	return platform, nil
}

// All returns all Platforms
func (s *PlatformService) All() ([]models.Platform, error) {
	var platforms []models.Platform
	if err := s.Repo.FindAll(&platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}

// FindByID returns a Platform by ID
func (s *PlatformService) FindByID(id uint64) (*models.Platform, error) {
	var platform models.Platform
	if err := s.Repo.FindByID(id, &platform); err != nil {
		return nil, errors.New("platform not found")
	}
	return &platform, nil
}

// Update a Platform by ID
// func (s *PlatformService) Update(id uint64, name string, status string, ipAllowlist []string, apiKey string,
// 	secretKey string) (*models.Platform, error) {
// 	platform, err := s.FindByID(id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if name != "" {
// 		platform.Name = name
// 	}
// 	if status != "" {
// 		platform.Status = status
// 	}
// 	// if ipAllowlist != "" {
// 	// 	platform.IPAllowlist = ipAllowlist
// 	// }
// 	if len(ipAllowlist) > 0 {
// 		platform.IPAllowlist = ipAllowlist
// 	}

// 	platform.UpdatedAt = time.Now()

// 	if err := s.Repo.Update(platform); err != nil {
// 		return nil, err
// 	}

//		return platform, nil
//	}
func (s *PlatformService) Update(
	id uint64,
	name string,
	email string,
	status string,
	ipAllowlist []string,
	apiKey string,
	secretKey string,
) (*models.Platform, error) {

	platform, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		platform.Name = name
	}
	if email != "" {
		platform.Email = email
	}

	if status != "" {
		platform.Status = status
	}

	if len(ipAllowlist) > 0 {
		platform.IPAllowlist = ipAllowlist
	}

	if err := s.Repo.Update(platform); err != nil {

		if strings.Contains(err.Error(), "uni_platforms_email") {
			return nil, errors.New("email already exists")
		}

		return nil, err
	}
	log.Println("secret_key_outside")
	// Update API Key table
	if apiKey != "" || secretKey != "" {

		platformApiKey, err := s.Repo.GetApiKeyByPlatformID(id)
		if err != nil {
			return nil, err
		}

		if apiKey != "" {
			platformApiKey.Key = apiKey
		}

		if secretKey != "" {
			platformApiKey.Secret = secretKey
		}
		if platform.Status == "disabled" {
			platformApiKey.Status = "disabled"
		} else {
			platformApiKey.Status = "active"
		}

		if err := s.Repo.UpdateApiKey(platformApiKey); err != nil {
			if strings.Contains(err.Error(), "uni_platform_api_keys_secret") {
				log.Println("secret key")
				return nil, errors.New("secret key already exists")
			}
			if strings.Contains(err.Error(), "uni_platform_api_keys") {
				return nil, errors.New("api key already exists")
			}
			return nil, err
		}
	}

	return platform, nil
}

// Delete a Platform by ID
func (s *PlatformService) Delete(id uint64) error {
	platform, err := s.FindByID(id)
	if err != nil {
		return err
	}
	return s.Repo.Delete(platform)
}
func (s *PlatformService) GetMerchantCurrency() ([]models.MerchantCurrency, error) {
	//return s.Repo.GetMerchantCurrency()
	var currencies []models.MerchantCurrency
	if err := s.Repo.GetMerchantCurrency(&currencies); err != nil {
		return nil, err
	}
	return currencies, nil
}
func (s *PlatformService) GetMerchantList(platformID uint64) ([]models.MerchantCurrency, error) {
	//return s.Repo.GetMerchantCurrency()
	var currencies []models.MerchantCurrency
	if err := s.Repo.GetMerchantCurrencyByPlatformID(platformID, &currencies); err != nil {
		return nil, err
	}
	return currencies, nil
}

func (s *PlatformService) CreateMerchantCurrency(
	platformID uint,
	merchantCode string,
	currency string,
	status string,
) (*models.MerchantCurrency, error) {

	exists, err := s.Repo.Exists(merchantCode)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New("merchant code already exists")
	}

	merchantCurrency := &models.MerchantCurrency{
		MerchantCode: merchantCode,
		Currency:     currency,
		PlatformID:   platformID,
		Status:       status,
	}

	err = s.Repo.CreateMerchantCurrency(merchantCurrency)
	if err != nil {
		return nil, err
	}

	return merchantCurrency, nil

}
func (s *PlatformService) UpdateMerchantCurrency(
	id uint,
	merchantCode string,
	currency string,
	status string,
) (*models.MerchantCurrency, error) {

	merchantCurrency, err := s.Repo.FindMerchantCurrencyByID(id)
	if err != nil {
		return nil, err
	}

	exists, err := s.Repo.ExistsMerchantCurrency(
		merchantCode,
		id,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("merchant code already exists")
	}

	merchantCurrency.MerchantCode = merchantCode
	merchantCurrency.Currency = currency
	merchantCurrency.Status = status
	if err := s.Repo.UpdateMerchantCurrency(merchantCurrency); err != nil {
		return nil, err
	}

	return merchantCurrency, nil
}

// All returns all Platforms
func (s *PlatformService) AllEnable() ([]models.Platform, error) {
	var platforms []models.Platform
	if err := s.Repo.FindAllEnable(&platforms); err != nil {
		return nil, err
	}
	return platforms, nil
}
func (s *PlatformService) AuditAction(
	c *fiber.Ctx,
	action string,
	model string,
	user jwt.MapClaims,
	oldData any,
	newData any,
) {

	helper.PublishAudit(
		s.Audit,
		c,
		action,
		model,
		user,
		oldData,
		newData,
	)
}
func (s *PlatformService) FindMerchantCurrencyByID(id uint) (*models.MerchantCurrency, error) {
	merchantCurrency, err := s.Repo.FindMerchantCurrencyByID(id)
	if err != nil {
		return nil, err
	}
	return merchantCurrency, nil
}
