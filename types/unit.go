package types

import "sync"

type Storage []*Unit

type Unit struct {
	Expression       string     `json:"expression"`
	Frequency        int        `json:"frequency"`
	Translates       []string   `json:"translates"`
	TranslatesNative []string   `json:"translatesNative"`
	Usages           []string   `json:"usages"`
	Mutex            sync.Mutex `json:"-"`
}
