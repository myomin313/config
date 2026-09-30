package Services

import (
	enums "config-service/Enums"
	"config-service/Repositories"
	"config-service/Resources"
	"config-service/config"
	"config-service/helper"
	"config-service/models"
	"config-service/pb/commonpb"
	"config-service/requests"
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type BannerService struct {
	Repo  *Repositories.BannerRepository
	Audit *helper.AuditService
}

func NewBannerService(repo *Repositories.BannerRepository) *BannerService {
	return &BannerService{Repo: repo}
}

func mapToBannerResource(p *models.Banner) Resources.BannerResource {

	conn := config.GrpcClient()
	defer conn.Close()

	client := commonpb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call GetUser
	resp, err := client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.CreatedByID)})
	if err != nil {
		log.Println("Error calling GetUser:", err)
	}
	// Declare pointer for UpdatedBy
	var respupdate *commonpb.GetUserResponse
	respupdate = resp

	if p.CreatedByID != p.UpdatedByID {
		respupdate, err = client.GetUser(ctx, &commonpb.GetUserRequest{Id: int64(p.UpdatedByID)})
		if err != nil {
			log.Println("Error calling GetUser:", err)
		}
	}

	var bannerImages []Resources.BannerImageResource
	if len(p.BannerImages) > 0 {
		for _, image := range p.BannerImages {
			bannerImages = append(bannerImages, Resources.BannerImageResource{
				ID:           image.ID,
				ImagePath:    os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + *image.ImagePath,
				LangCode:     image.LanguageCode,
				DefaultImage: image.DefaultImage,
			})
		}
	}

	endDate := ""
	if p.EndDate != nil {
		endDate = p.EndDate.Format("2006-01-02 15:04:05")
	}
	return Resources.BannerResource{
		ID:              p.ID,
		ImageName:       p.ImageName,
		DisplayPlatform: p.DisplayPlatform.String(),
		DisplayDuration: p.DisplayDuration.String(),
		Status:          p.Status.String(),
		StartDate:       p.StartDate.Format("2006-01-02 15:04:05"),
		EndDate:         endDate,
		UpdatedBy:       respupdate.Name,
		CreatedBy:       resp.Name,
		CreatedAt:       p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       p.UpdatedAt.Format("2006-01-02 15:04:05"),
		BannerImages:    bannerImages,
	}

}

func mapToBannerListResource(p *models.Banner) Resources.BannerListResource {

	var bannerImages []Resources.BannerImageResource
	if len(p.BannerImages) > 0 {
		for _, image := range p.BannerImages {
			bannerImages = append(bannerImages, Resources.BannerImageResource{
				ID:           image.ID,
				ImagePath:    os.Getenv("MINIO_PUBLIC_ENDPOINT") + "/" + *image.ImagePath,
				LangCode:     image.LanguageCode,
				DefaultImage: image.DefaultImage,
			})
		}
	}

	return Resources.BannerListResource{
		ImageName:    p.ImageName,
		BannerImages: bannerImages,
	}

}

// Create a new Banner
func (s *BannerService) Create(req *requests.BannerCreateRequest, id uint64) (Resources.BannerResource, error) {

	banner := &models.Banner{
		ImageName:       req.ImageName,
		DisplayPlatform: enums.ParseDisplayPlatform(req.DisplayPlatform),
		DisplayDuration: enums.ParseDisplayDuration(req.DisplayDuration),
		StartDate:       req.StartDate.ToTime(),
		CreatedByID:     id,
		UpdatedByID:     id,
	}
	// ✅ Only set EndDate when NOT display_duration = 2
	if req.DisplayDuration != "permanent" {
		banner.EndDate = req.EndDate.ToTime()
	}
	var images []models.BannerImage
	for _, img := range req.BannerImage {

		images = append(images, models.BannerImage{
			ImagePath:    img.ImagePath,
			LanguageCode: img.LangCode,
			DefaultImage: img.DefaultImage,
		})
	}
	if err := s.Repo.Create(banner, images); err != nil {
		return Resources.BannerResource{}, err
	}
	return mapToBannerResource(banner), nil
}

// GetByID fetches a single banner by its ID
func (s *BannerService) GetByID(id uint64) (Resources.BannerResource, error) {
	banner, err := s.Repo.FindBannerWithImagesById(id)
	if err != nil {
		return Resources.BannerResource{}, err
	}

	return mapToBannerResource(banner), nil
}

// GetAll returns all banners
func (s *BannerService) GetAll(c *fiber.Ctx, req *requests.BannerSearchRequest) (map[string]interface{}, error) {
	pagination := helper.GetPagination(c)
	// 2. Fetch from repo
	banners, total, err := s.Repo.All(pagination.Offset, pagination.Limit, req)
	if err != nil {
		return nil, err
	}
	bannerResources := make([]Resources.BannerResource, 0, total)
	for i := range banners {
		bannerResources = append(bannerResources, mapToBannerResource(&banners[i]))
	}

	totalPage := (total + int64(pagination.PageSize) - 1) / int64(pagination.PageSize)

	// 4. Return formatted paginated response
	return map[string]interface{}{
		"records":    bannerResources,
		"page_no":    pagination.PageNo,
		"page_size":  pagination.PageSize,
		"total":      total,
		"total_page": totalPage,
	}, nil

}

// Update a banner by request
func (s *BannerService) Update(id uint64, req *requests.BannerUpdateRequest, userId uint64) (*models.Banner, error) {
	banner, err := s.Repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ImageName != nil {
		banner.ImageName = *req.ImageName
	}

	if req.DisplayPlatform != nil {
		banner.DisplayPlatform = enums.ParseDisplayPlatform(*req.DisplayPlatform)
	}

	if req.DisplayDuration != nil {
		banner.DisplayDuration = enums.ParseDisplayDuration(*req.DisplayDuration)
	}

	if req.StartDate != nil {
		banner.StartDate = req.StartDate.ToTime()
	}

	if req.EndDate != nil {
		banner.EndDate = req.EndDate.ToTime()
	}

	banner.UpdatedByID = userId
	banner.UpdatedAt = time.Now()

	var images []models.BannerImage
	for _, img := range req.BannerImage {
		var ImageID uint64 = 0
		if img.ID != nil {
			ImageID = *img.ID
		}

		images = append(images, models.BannerImage{
			ID:           uint64(ImageID),
			ImagePath:    img.ImagePath,
			LanguageCode: img.LangCode,
			DefaultImage: img.DefaultImage,
		})
	}
	log.Println("images:", images)
	if err := s.Repo.Update(banner, images); err != nil {
		return nil, err
	}
	return banner, nil
}

// Delete removes a banner by ID
func (s *BannerService) Delete(id uint64) error {
	_, err := s.Repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.Repo.Delete(id)
}

func (s *BannerService) MassUpdateStatus(
	tableIDs []uint,
	status string,
) (int64, error) {
	return s.Repo.MassUpdateStatus(tableIDs, enums.ParseStatus(status))
}

func (s *BannerService) MassDelete(
	tableIDs []uint,
) (int64, error) {
	return s.Repo.MassDelete(tableIDs)
}

func (s *BannerService) GetActiveBanners(
	displayPlatform string, lang string,
) ([]Resources.BannerListResource, error) {

	banners, err := s.Repo.GetActiveBanners(enums.ParseDisplayPlatform(displayPlatform), lang)
	if err != nil {
		return nil, err
	}

	bannerResources := make([]Resources.BannerListResource, 0, len(banners))

	for _, banner := range banners {
		bannerResources = append(
			bannerResources,
			mapToBannerListResource(&banner),
		)
	}

	return bannerResources, nil
}
func (s *BannerService) AuditAction(
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
func (s *BannerService) GetByIDs(ids []uint) ([]Resources.BannerResource, error) {
	banners, err := s.Repo.FindBannersWithImagesByIDs(ids)
	if err != nil {
		return nil, err
	}

	resources := make([]Resources.BannerResource, 0, len(banners))
	for _, banner := range banners {
		b := banner // avoid taking the address of the range variable
		resources = append(resources, mapToBannerResource(&b))
	}

	return resources, nil
}
