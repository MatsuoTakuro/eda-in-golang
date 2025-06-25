package main

import (
	"eda-in-golang/internal/system"
)

type monolith struct {
	svc     system.Service
	modules []system.Module
}

func (m *monolith) startupModules() error {
	for _, module := range m.modules {
		if err := module.Startup(m.svc.Runner().Context(), m.svc); err != nil {
			return err
		}
	}

	return nil
}
