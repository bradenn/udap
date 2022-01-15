// Copyright (c) 2021 Braden Nicholson

package models

// Core attributes
var (
	On = &Attribute{
		Key: "on",
	}
	Dim = &Attribute{
		Key: "dim",
	}
	Hue = &Attribute{
		Key: "hue",
	}
	Cct = &Attribute{
		Key: "cct",
	}
)

func NewMediaEntity(name string, module string) *Entity {
	e := Entity{
		Name:       name,
		Type:       "media",
		Attributes: []*Attribute{},
		Module:     module,
	}
	return &e
}

func NewEntity(name string, module string) *Entity {
	e := Entity{
		Name:       name,
		Type:       "switch",
		Module:     module,
		Attributes: []*Attribute{},
	}
	return &e
}

func NewSpectrum(name string, module string) *Entity {
	e := Entity{
		Name:       name,
		Type:       "spectrum",
		Module:     module,
		Attributes: []*Attribute{On, Dim, Hue, Cct},
	}
	return &e
}

func NewSwitch(name string, module string) *Entity {
	e := Entity{
		Name:       name,
		Type:       "switch",
		Module:     module,
		Attributes: []*Attribute{On},
	}
	return &e
}
