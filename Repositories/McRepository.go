package Repositories

import (
	"config-service/models"

	"gorm.io/gorm"
)

type MaintenanceRepository struct {
	DB *gorm.DB
}

func (r *MaintenanceRepository) Get() (*models.Maintenance, error) {
	var maintenance models.Maintenance

	err := r.DB.First(&maintenance, 1).Error
	if err != nil {
		return nil, err
	}

	return &maintenance, nil
}

func (r *MaintenanceRepository) Update(m *models.Maintenance) error {
	return r.DB.Save(m).Error
}
