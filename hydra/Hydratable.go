package hydra

import (
	"reflect"
)

const tagName = "hydra"

type IHydratable interface {
	Init(interface{}) interface{}
}

type Hydratable struct {
	XDBTypeOverride string      // The database type to override the default detection, e.g. sqlite, mssql, mariadb, oracle, mysql
	name            string      // The name of the object to be hydrated
	isInitialized   bool        // Flag to check if the object has been initialized
	self            interface{} // The object that is to be hydrated
}

func (h *Hydratable) Init(o interface{}) interface{} {
	v := reflect.ValueOf(o)
	t := reflect.TypeOf(o)

	// If it's a pointer, get the underlying element
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		p("Passed value is not a struct")
		return o
	}

	if debug {
		p("Initializing type:", t.Name())
		p("Number of fields:", t.NumField())

		for i := 0; i < t.NumField(); i++ {
			// Get the field metadata
			field := t.Field(i)

			// Get the field tag value for "hydra"
			tag := field.Tag.Get(tagName)

			// If the tag is not empty, add it to the hydrate list
			if tag != "" {
				p("%d. Field: %v (%v), tag: '%v' - Will be hydrated\n", i+1, field.Name, field.Type, tag)
			} else {
				p("%d. Field: %v (%v), tag: '%v' - Skipped\n", i+1, field.Name, field.Type, tag)
			}
		}
	}

	// Save a reference to the passed object in 'self'
	h.self = o
	h.isInitialized = true
	h.name = t.Name()

	return o
}
