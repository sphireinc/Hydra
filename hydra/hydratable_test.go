package hydra

import (
	"context"
	"errors"
	"testing"
)

type initPerson struct {
	Hydratable
	ID int `hydra:"id,pk"`
}

func TestInitTableDriven(t *testing.T) {
	tests := []struct {
		name              string
		build             func() (input interface{}, hyd *Hydratable)
		expectInitialized bool
		expectSelfStored  bool
	}{
		{
			name: "pointer struct",
			build: func() (interface{}, *Hydratable) {
				p := &initPerson{}
				return p, &p.Hydratable
			},
			expectInitialized: true,
			expectSelfStored:  true,
		},
		{
			name: "non pointer struct",
			build: func() (interface{}, *Hydratable) {
				p := initPerson{}
				return p, &p.Hydratable
			},
			expectInitialized: true,
			expectSelfStored:  true,
		},
		{
			name: "non struct input",
			build: func() (interface{}, *Hydratable) {
				h := &Hydratable{}
				return 123, h
			},
			expectInitialized: false,
			expectSelfStored:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, hyd := tt.build()
			hyd.Init(input)

			if hyd.isInitialized != tt.expectInitialized {
				t.Fatalf("unexpected initialized state: %v", hyd.isInitialized)
			}
			if (hyd.self != nil) != tt.expectSelfStored {
				t.Fatalf("unexpected self state: %#v", hyd.self)
			}
		})
	}
}

func TestResolveTableNameRequiresInit(t *testing.T) {
	h := &Hydratable{}
	_, err := h.resolveTableName()
	if !errors.Is(err, ErrNotInitialized) {
		t.Fatalf("expected ErrNotInitialized, got: %v", err)
	}
}

func TestHydrateByPrimaryKeyRequiresTag(t *testing.T) {
	type noPK struct {
		Hydratable
		ID int `hydra:"id"`
	}

	obj := &noPK{}
	obj.Init(obj)

	if err := obj.HydrateByPrimaryKey(nil, 1); err == nil {
		t.Fatal("expected error for missing pk tag")
	}
}

func TestHydrateByLookupRequiresTag(t *testing.T) {
	type noLookup struct {
		Hydratable
		Email string `hydra:"email"`
	}

	obj := &noLookup{}
	obj.Init(obj)

	if err := obj.HydrateByLookupContext(context.Background(), nil); err == nil {
		t.Fatal("expected error for missing lookup tag")
	}
}
