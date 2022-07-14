package entity

import "time"

// Employee entity which the TiRelease logic should use.
type Employee struct {
	ID           int64      `gorm:"column:id"`
	HrEmployeeID string     `gorm:"column:hr_employee_id"`
	Name         string     `gorm:"column:name"`
	Email        string     `gorm:"column:email"`
	GithubId     string     `gorm:"column:github_id"`
	Country      string     `gorm:"column:country"`
	CreateDate   *time.Time `gorm:"column:create_date"`
	DeleteDate   *time.Time `gorm:"column:delete_date"`
	IsActive     bool       `gorm:"column:active"`
	GhEmail      string     `gorm:"column:gh_email"`
	GhId         string     `gorm:"column:gh_id"`
	GhName       string     `gorm:"column:gh_name"`
	JobNumber    string     `gorm:"column:job_number"`
}

// Employee entity from hr_employee table
// Which is derived ouside TiRelease
// Hence it's ought not to be modified.
type HrEmployee struct {
	ID         string
	Name       string
	Email      string
	GithubId   string
	Country    string
	CreateDate string
	DeleteDate string
	IsActive   bool `gorm:"column:active"`
	GhEmail    string
	GhId       string
	GhName     string
	JobNumber  string
}

func (employee HrEmployee) Trans() Employee {
	var createDate *time.Time
	if employee.CreateDate == "" {
		createDate = nil
	} else {
		tmp, _ := time.Parse("2006-01-02", employee.CreateDate)
		createDate = &tmp
	}

	var deleteDate *time.Time
	if employee.DeleteDate == "" {
		deleteDate = nil
	} else {
		tmp, _ := time.Parse("2006-01-02", employee.DeleteDate)
		deleteDate = &tmp
	}
	return Employee{
		HrEmployeeID: employee.ID,
		Name:         employee.Name,
		Email:        employee.Email,
		GithubId:     employee.GithubId,
		Country:      employee.Country,
		CreateDate:   createDate,
		DeleteDate:   deleteDate,
		IsActive:     employee.IsActive,
		GhEmail:      employee.GhEmail,
		GhId:         employee.GhId,
		GhName:       employee.GhName,
		JobNumber:    employee.JobNumber,
	}
}
