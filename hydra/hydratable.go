package hydra

import (
	"context"
	"fmt"
	"reflect"
)

const tagName = "hydra"

type IHydratable interface {
	Init(interface{}) interface{}
}

type Hydratable struct {
	XDBTypeOverride    string      // The database type to override the default detection, e.g. sqlite, mssql, mariadb, oracle, mysql
	XTableNameOverride string      // Optional explicit table name override
	name               string      // The name of the object to be hydrated
	isInitialized      bool        // Flag to check if the object has been initialized
	self               interface{} // The object that is to be hydrated
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
			field := t.Field(i)
			tag := field.Tag.Get(tagName)

			if tag != "" {
				p(fmt.Sprintf("%d. Field: %v (%v), tag: '%v' - Will be hydrated\n", i+1, field.Name, field.Type, tag))
			} else {
				p(fmt.Sprintf("%d. Field: %v (%v), tag: '%v' - Skipped\n", i+1, field.Name, field.Type, tag))
			}
		}
	}

	h.self = o
	h.isInitialized = true
	h.name = t.Name()
	_ = v

	return o
}

func (h *Hydratable) resolveTableName() (string, error) {
	if !h.isInitialized {
		return "", ErrNotInitialized
	}

	if h.XTableNameOverride != "" {
		if err := validateIdentifier(h.XTableNameOverride); err != nil {
			return "", err
		}
		return h.XTableNameOverride, nil
	}

	if provider, ok := h.self.(tableNameProvider); ok {
		name := provider.HydraTableName()
		if name != "" {
			if err := validateIdentifier(name); err != nil {
				return "", err
			}
			return name, nil
		}
	}

	name := lowerASCII(h.name)
	if err := validateIdentifier(name); err != nil {
		return "", err
	}
	return name, nil
}

func (h *Hydratable) Hydrate(db any, whereClauses map[string]interface{}) error {
	return h.HydrateContext(context.Background(), db, whereClauses)
}

func (h *Hydratable) HydrateContext(ctx context.Context, db any, whereClauses map[string]interface{}) error {
	return h.hydrate(ctx, db, whereClauses)
}

func (h *Hydratable) HydrateByPrimaryKey(db any, value any) error {
	return h.HydrateByPrimaryKeyContext(context.Background(), db, value)
}

func (h *Hydratable) HydrateByPrimaryKeyContext(ctx context.Context, db any, value any) error {
	whereClauses, err := h.buildPrimaryKeyWhereClauses(value)
	if err != nil {
		return err
	}
	return h.HydrateContext(ctx, db, whereClauses)
}

func (h *Hydratable) HydrateByLookup(db any) error {
	return h.HydrateByLookupContext(context.Background(), db)
}

func (h *Hydratable) HydrateByLookupContext(ctx context.Context, db any) error {
	whereClauses, err := h.buildLookupWhereClauses()
	if err != nil {
		return err
	}
	return h.HydrateContext(ctx, db, whereClauses)
}

func (h *Hydratable) buildPrimaryKeyWhereClauses(value any) (map[string]interface{}, error) {
	if !h.isInitialized {
		return nil, ErrNotInitialized
	}

	meta, err := getTypeMetadata(h.self)
	if err != nil {
		return nil, err
	}

	if len(meta.PrimaryKeys) == 0 {
		return nil, fmt.Errorf("no primary-key fields configured on %s", meta.Type.Name())
	}

	if len(meta.PrimaryKeys) == 1 {
		return map[string]interface{}{meta.PrimaryKeys[0].Column: value}, nil
	}

	m, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("composite primary keys require map[string]interface{} input")
	}

	where := make(map[string]interface{}, len(meta.PrimaryKeys))
	for _, field := range meta.PrimaryKeys {
		v, ok := m[field.Column]
		if !ok {
			return nil, fmt.Errorf("missing primary-key value for column %q", field.Column)
		}
		where[field.Column] = v
	}

	return where, nil
}

func (h *Hydratable) buildLookupWhereClauses() (map[string]interface{}, error) {
	if !h.isInitialized {
		return nil, ErrNotInitialized
	}

	meta, err := getTypeMetadata(h.self)
	if err != nil {
		return nil, err
	}

	if len(meta.LookupFields) == 0 {
		return nil, fmt.Errorf("no lookup fields configured on %s", meta.Type.Name())
	}

	v := reflect.ValueOf(h.self)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	where := make(map[string]interface{}, len(meta.LookupFields))
	for _, field := range meta.LookupFields {
		fv := v.Field(field.Index)
		value, ok := reflectValueToInterface(fv)
		if !ok {
			return nil, fmt.Errorf("lookup field %s is nil", field.Name)
		}
		where[field.Column] = value
	}

	return where, nil
}

func reflectValueToInterface(v reflect.Value) (any, bool) {
	if !v.IsValid() {
		return nil, false
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, false
		}
		return reflectValueToInterface(v.Elem())
	}
	return v.Interface(), true
}
