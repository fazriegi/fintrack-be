package pkg

import (
	"errors"
	"reflect"

	"github.com/jmoiron/sqlx"
)

// returns elements that appear in both slice1 and slice2.
func Intersection[T comparable](slice1, slice2 []T) []T {
	set := make(map[T]struct{})
	for _, v := range slice1 {
		set[v] = struct{}{}
	}

	var matches []T
	for _, v := range slice2 {
		if _, ok := set[v]; ok {
			matches = append(matches, v)
		}
	}

	return matches
}

func ScanRowsIntoStructs(rows *sqlx.Rows, destSlice interface{}) error {
	destVal := reflect.ValueOf(destSlice)
	if destVal.Kind() != reflect.Ptr || destVal.Elem().Kind() != reflect.Slice {
		return errors.New("destSlice must be a pointer to a slice")
	}

	sliceType := destVal.Elem().Type().Elem()
	if sliceType.Kind() != reflect.Struct {
		return errors.New("slice elements must be structs")
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}

	for rows.Next() {
		elem := reflect.New(sliceType).Elem()
		// Create a DB column → struct field mapping based on the db tag
		fieldMap := make(map[string]reflect.Value)
		for i := 0; i < elem.NumField(); i++ {
			field := sliceType.Field(i)
			tag := field.Tag.Get("db")
			if tag == "" {
				continue
			}
			fieldMap[tag] = elem.Field(i)
		}

		// Prepare a place to scan query results
		scanArgs := make([]interface{}, len(columns))
		for i, col := range columns {
			if field, exists := fieldMap[col]; exists {
				// Handle interface{} fields by scanning into appropriate types
				if field.Kind() == reflect.Interface {
					switch columnTypes[i].DatabaseTypeName() {
					case "INT", "INTEGER", "BIGINT":
						var v int64
						scanArgs[i] = &v
					case "FLOAT", "REAL", "DOUBLE":
						var v float64
						scanArgs[i] = &v
					case "DECIMAL", "NUMERIC":
						var v float64
						scanArgs[i] = &v
					case "VARCHAR", "TEXT", "STRING":
						var v string
						scanArgs[i] = &v
					default:
						// Fallback to interface{} for unknown types
						scanArgs[i] = field.Addr().Interface()
					}
				} else {
					scanArgs[i] = field.Addr().Interface()
				}
			} else {
				// Ignore unused columns
				var dummy interface{}
				scanArgs[i] = &dummy
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}

		// Convert scanned values to interface{} for target fields
		for i, col := range columns {
			if field, exists := fieldMap[col]; exists && field.Kind() == reflect.Interface {
				val := reflect.ValueOf(scanArgs[i]).Elem().Interface()
				field.Set(reflect.ValueOf(val))
			}
		}

		destVal.Elem().Set(reflect.Append(destVal.Elem(), elem))
	}

	return rows.Err()
}
