package hydra

import (
	"fmt"
	"reflect"
	"strings"
)

func (h *Hydratable) Hydrate(db any, whereClauses map[string]interface{}) error {
	if !h.isInitialized {
		p("Not initialized, cannot hydrate.")
		return fmt.Errorf("not initialized, cannot hydrate")
	}

	// Reflect on the actual value of `self` (which holds the parent struct)
	v := reflect.ValueOf(h.self)
	if v.Kind() == reflect.Ptr {
		v = v.Elem() // Get the actual struct if it's a pointer
	}
	t := v.Type()
	p("Hydrating struct:", t.Name())

	// Get the table name from the struct's type name
	tableName := strings.ToLower(t.Name())
	p("Table name:", tableName)

	// Fetch data from the database using the table name and primary key (id assumed here)
	data, err := h.Fetch(db, tableName, whereClauses) // Pass the table name and where clauses
	p("Data:", data)
	if err != nil {
		p("Error fetching data:", err)
		return fmt.Errorf("error fetching data: %v", err)
	}

	// Loop through the fields of the struct and hydrate based on hydra tags
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)

		// If the field has a hydra tag, try to hydrate it
		if tag != "" {
			if value, ok := data[tag]; ok {
				fieldValue := v.FieldByName(field.Name)

				// Ensure the field is settable
				if fieldValue.CanSet() {
					// Set the field's value from the database data
					switch fieldValue.Kind() {
					case reflect.String:
						if bytes, ok := value.([]byte); ok {
							fieldValue.SetString(string(bytes))
						}
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						fieldValue.SetInt(reflect.ValueOf(value).Int())
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
						fieldValue.SetUint(reflect.ValueOf(value).Uint())
					case reflect.Float32, reflect.Float64:
						fieldValue.SetFloat(reflect.ValueOf(value).Float())
					case reflect.Bool:
						fieldValue.SetBool(value.(bool))
					case reflect.Slice:
						fieldValue.Set(reflect.ValueOf(value)) // Assuming the value is a slice and directly assignable
					case reflect.Struct:
						fieldValue.Set(reflect.ValueOf(value)) // Assuming the value is a struct
					case reflect.Ptr:
						fieldValue.Set(reflect.ValueOf(value)) // Handle pointers by setting the pointer to the value
					default:
						fieldValue.Set(reflect.ValueOf(value)) // Hail mary
						p(fmt.Sprintf("unhandled type: %v\n", fieldValue.Kind()))
					}
				}
			}
		}
	}

	return nil // Hydration success
}
