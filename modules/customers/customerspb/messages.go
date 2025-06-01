package customerspb

import (
	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/registry/registrar"
)

const (
	CustomerAggregateChannel = "mallbots.customers.events.Customer"

	CustomerRegisteredEvent = "customersapi.CustomerRegistered"
	CustomerSmsChangedEvent = "customersapi.CustomerSmsChanged"
	CustomerEnabledEvent    = "customersapi.CustomerEnabled"
	CustomerDisabledEvent   = "customersapi.CustomerDisabled"

	CommandChannel = "mallbots.customers.commands"

	AuthorizeCustomerCommand = "customersapi.AuthorizeCustomer"
)

func RegisterMessages(reg registry.Registry) error {
	regtr := registrar.NewProtoRegistrar(reg)

	// Customer events
	if err := regtr.Register(&CustomerRegistered{}); err != nil {
		return err
	}
	if err := regtr.Register(&CustomerSmsChanged{}); err != nil {
		return err
	}
	if err := regtr.Register(&CustomerEnabled{}); err != nil {
		return err
	}
	if err := regtr.Register(&CustomerDisabled{}); err != nil {
		return err
	}

	// commands
	if err := regtr.Register(&AuthorizeCustomer{}); err != nil {
		return err
	}
	return nil
}

var (
	_ registry.Registrable = (*CustomerRegistered)(nil)
	_ registry.Registrable = (*CustomerSmsChanged)(nil)
	_ registry.Registrable = (*CustomerEnabled)(nil)
	_ registry.Registrable = (*CustomerDisabled)(nil)
)

func (*CustomerRegistered) Key() string { return CustomerRegisteredEvent }
func (*CustomerSmsChanged) Key() string { return CustomerSmsChangedEvent }
func (*CustomerEnabled) Key() string    { return CustomerEnabledEvent }
func (*CustomerDisabled) Key() string   { return CustomerDisabledEvent }

var (
	_ registry.Registrable = (*AuthorizeCustomer)(nil)
)

func (*AuthorizeCustomer) Key() string { return AuthorizeCustomerCommand }
