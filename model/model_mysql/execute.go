package model_mysql

import (
	"database/sql"
	"fmt"
	"maoguo/henan/misc/page"
	"maoguo/henan/misc/utils/parse"

	"github.com/wonderivan/logger"
)

type Myrows struct {
	sql.Rows
}

//原生sql查询,并返回map[string]string
func Query(sql string, args ...interface{}) map[string]interface{} {
	rows, err := Db.Raw(sql, args...).Rows()
	defer rows.Close()
	if err != nil {
		return nil
	}
	result := ScanMap(rows)
	return result
}

//原生sql查询，并返回[]map[string]string
func QueryList(sql string, args ...interface{}) []map[string]interface{} {
	rows, err := Db.Raw(sql, args...).Rows()
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	result := ScanMaps(rows)
	return result
}

func QueryPage(sql string, pg page.Page, args ...interface{}) page.Page {
	sql = sql + " limit " + parse.IntToString((pg.PageNo-1)*pg.PageSize) + "," + parse.IntToString(pg.PageSize)
	rows, err := Db.Raw(sql, args...).Rows()
	defer rows.Close()
	if err != nil {
		logger.Info("queryPage failed", err)
		return pg
	}
	result := ScanMaps(rows)
	if result == nil {
		result = make([]map[string]interface{}, 0)
	}
	pg.List = result
	return pg
}

func QueryColumn(sql string, args ...interface{}) []interface{} {
	rows, err := Db.Raw(sql, args...).Rows()
	if err != nil {
		return nil
	}
	result := ScanStringList(rows)
	return result
}

//执行原生sql(增删改)成功返回true否则返回false
func Exec(sql string, args ...interface{}) bool {
	err := Db.Debug().Exec(sql, args...).Error
	if err != nil {
		return false
	}
	return true
}

//执行原生sql，返回成功的条数
func ExecInt(sql string, args ...interface{}) int {
	success := Db.Debug().Exec(sql, args...).RowsAffected
	return parse.Int64ToInt(success)
}

//执行原生sql，返回执行出错信息
func Execer(sql string, args ...interface{}) error {
	err := Db.Debug().Exec(sql, args...).Error
	return err
}

// type Myrows sql.Rows
func ScanMap(rows *sql.Rows) map[string]interface{} {
	result := ScanMaps(rows)
	if len(result) > 0 {
		return result[0]
	}
	return nil
}

func ScanMaps(rows *sql.Rows) []map[string]interface{} {
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var result []map[string]interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		var value interface{}
		obj := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			obj[columns[i]] = value
		}
		result = append(result, obj)
	}
	return result
}

func ScanStringList(rows *sql.Rows) []interface{} {
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var result []interface{}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		var value interface{}
		for _, col := range values {
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			result = append(result, value)
		}
	}
	return result
}
