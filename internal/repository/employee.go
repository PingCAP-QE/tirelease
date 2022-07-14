package repository

import (
	"fmt"
	"tirelease/commons/database"
	"tirelease/internal/entity"

	"github.com/pkg/errors"
	"gorm.io/gorm/clause"
)

func BatchCreateOrUpdateEmployees(employees []entity.Employee) error {
	if err := database.DBConn.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&employees).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("create or update employees: %+v failed", employees))
	}
	return nil
}
