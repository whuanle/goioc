package goioc

type ServiceLifetime int

const (
	Transient ServiceLifetime = iota
	Scope
	Singleton
)
