package repository


import (
		"fmt"
)


func ConstructPatchQuery(tableName, keyName string, fields []string) string {
		query := fmt.Sprintf("UPDATE %s SET", tableName)	
		for i, v := range fields {
			query += fmt.Sprintf(" %v", v)
			if i != len(fields)-1 {
				query += " = ?," 	
			} else {
				query += " = ?\n"
			}
		} 
		query += fmt.Sprintf("WHERE %s = ?;", keyName)
		return query
}
