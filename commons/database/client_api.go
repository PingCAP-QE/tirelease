package database

import (
	"strings"

	"gorm.io/gorm"
)

func (conn MysqlInfo) RawWrapper(sql string, values ...interface{}) (tx *gorm.DB) {
	if strings.Contains(sql, "@") || strings.Contains(sql, "?") {
		return conn.DB.Raw(sql, values...)
	} else {
		return conn.DB.Raw(sql)
	}
}
