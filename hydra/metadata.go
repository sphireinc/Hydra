package hydra

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type HydraFieldConverter func(src any) (any, error)

type HydraConverterProvider interface {
	HydraConverters() map[string]HydraFieldConverter
}

type HydraValueConverter interface {
	HydraConvert(src any) error
}

type tableNameProvider interface {
	HydraTableName() string
}

type fieldMetadata struct {
	Index    int
	Name     string
	Column   string
	IsPK     bool
	IsLookup bool
}

type typeMetadata struct {
	Type         reflect.Type
	Fields       []fieldMetadata
	Columns      []string
	PrimaryKeys  []fieldMetadata
	LookupFields []fieldMetadata
}

var (
	identifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
	metadataCache     sync.Map
)

func validateIdentifier(identifier string) error {
	if !identifierPattern.MatchString(identifier) {
		return fmt.Errorf("invalid identifier %q", identifier)
	}
	return nil
}

func resolveStructType(self any) (reflect.Type, error) {
	if self == nil {
		return nil, fmt.Errorf("hydrated object is nil")
	}
	t := reflect.TypeOf(self)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("hydrated object must be a struct, got %T", self)
	}
	return t, nil
}

func getTypeMetadata(self any) (*typeMetadata, error) {
	t, err := resolveStructType(self)
	if err != nil {
		return nil, err
	}

	if cached, ok := metadataCache.Load(t); ok {
		return cached.(*typeMetadata), nil
	}

	meta := &typeMetadata{Type: t}
	seenColumns := make(map[string]struct{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		rawTag := strings.TrimSpace(field.Tag.Get(tagName))
		if rawTag == "" {
			continue
		}

		parts := strings.Split(rawTag, ",")
		column := strings.TrimSpace(parts[0])
		if column == "" {
			return nil, fmt.Errorf("field %s has empty hydra column tag", field.Name)
		}
		if err := validateIdentifier(column); err != nil {
			return nil, fmt.Errorf("field %s: %w", field.Name, err)
		}

		fm := fieldMetadata{
			Index:  i,
			Name:   field.Name,
			Column: column,
		}

		for _, part := range parts[1:] {
			switch strings.TrimSpace(strings.ToLower(part)) {
			case "pk", "primary", "primary_key":
				fm.IsPK = true
			case "lookup":
				fm.IsLookup = true
			case "":
			default:
				return nil, fmt.Errorf("field %s: unsupported hydra tag option %q", field.Name, strings.TrimSpace(part))
			}
		}

		meta.Fields = append(meta.Fields, fm)

		if _, ok := seenColumns[column]; !ok {
			seenColumns[column] = struct{}{}
			meta.Columns = append(meta.Columns, column)
		}

		if fm.IsPK {
			meta.PrimaryKeys = append(meta.PrimaryKeys, fm)
		}
		if fm.IsLookup {
			meta.LookupFields = append(meta.LookupFields, fm)
		}
	}

	if len(meta.Columns) == 0 {
		return nil, fmt.Errorf("no hydra-tagged fields found on %s", t.Name())
	}

	sort.Strings(meta.Columns)
	metadataCache.Store(t, meta)
	return meta, nil
}

func lowerASCII(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}
