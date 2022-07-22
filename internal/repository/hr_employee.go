package repository

import (
	"fmt"
	"time"
	"tirelease/commons/configs"
	"tirelease/internal/entity"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ATTENTION: there must not be any create/update/insert operation in this file.
var DB *gorm.DB

func SelectAllHrEmployee() ([]entity.HrEmployee, error) {
	var hrEmployees []entity.HrEmployee
	result := DB.Find(&hrEmployees)
	if result.Error != nil {
		return nil, result.Error
	}

	return hrEmployees, nil
}

func InitHrEmployeeDB() {
	test := configs.Config
	fmt.Printf("%+v\n", test)
	config := configs.Config.EmployeeMysql

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=%s",
		config.UserName,
		config.PassWord,
		config.Host,
		config.Port,
		config.DataBase,
		config.CharSet,
		config.TimeZone,
	)

	// Connect
	conn, err := gorm.Open(mysql.Open(url), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
	sqlDB, err := conn.DB()
	if err != nil {
		panic(err.Error())
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Second * 600)

	DB = conn
}
