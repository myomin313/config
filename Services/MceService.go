package Services

import (
	"config-service/Repositories"
	"config-service/helper"
	"config-service/models"
	"config-service/requests"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type MaintenanceService struct {
	Repo  *Repositories.MaintenanceRepository
	Audit *helper.AuditService
}

func (s *MaintenanceService) Get() (*models.Maintenance, error) {
	return s.Repo.Get()
}

func (s *MaintenanceService) Update(req *requests.MaintenanceUpdateRequest, userID uint) (*models.Maintenance, error) {
	m, err := s.Repo.Get()
	log.Printf("MaintenanceService Update: %v", m)
	if err != nil {
		return nil, err
	}

	if req.ApplicationMaintenance != nil {
		m.ApplicationMaintenance = *req.ApplicationMaintenance
	}

	if req.GameMaintenance != nil {
		m.GameMaintenance = *req.GameMaintenance
	}

	m.UpdatedByID = &userID
	if m.CreatedByID == nil {
		m.CreatedByID = &userID
	}
	log.Printf("MaintenanceService after patch: %v", m)
	if err := s.Repo.Update(m); err != nil {
		return nil, err
	}

	return m, nil
}
func (s *MaintenanceService) AuditAction(
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
