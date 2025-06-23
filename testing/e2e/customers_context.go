package e2e

import (
	"github.com/cucumber/godog"
	"github.com/go-openapi/strfmt"
	"github.com/stackus/errors"

	"eda-in-golang/modules/customers/customersclient"
	"eda-in-golang/modules/customers/customersclient/customer"
	"eda-in-golang/modules/customers/customersclient/models"
)

type customersContext struct {
	*suiteContext
	client           *customersclient.Customers
	fetchedCustomers bool
}

var _ featureContext = (*customersContext)(nil)

func newCustomersContext(sc *suiteContext) featureContext {
	return &customersContext{
		suiteContext: sc,
		client:       customersclient.New(sc.transport, strfmt.Default),
	}
}

func (c *customersContext) register(ctx *godog.ScenarioContext) {
	ctx.Step(`^I register a new customer as "([^"]*)"$`, c.iRegisterANewCustomerAs)
	ctx.Step(`^(?:ensure |expect )?a customer named "([^"]*)" (?:to )?exists?$`, c.expectACustomerNamedToExist)
	ctx.Step(`^(?:ensure |expect )?no customer named "([^"]*)" (?:to )?exists?$`, c.expectNoCustomerNamedToExist)

}

func (c *customersContext) reset() error {
	c.customers = make(map[string]string)
	c.fetchedCustomers = false
	if err := c.truncate("customers.customers"); err != nil {
		return err
	}
	if err := c.truncate("customers.inbox"); err != nil {
		return err
	}
	if err := c.truncate("customers.outbox"); err != nil {
		return err
	}
	return nil
}

func (c *customersContext) iRegisterANewCustomerAs(name string) {
	resp, err := c.client.Customer.CreateCustomer(customer.NewCreateCustomerParams().WithBody(&models.CustomerspbRegisterCustomerRequest{
		Name:      name,
		SmsNumber: "555-555-1212",
	}))
	if err != nil {
		c.lastErr = err
		return
	}

	c.customers[name] = resp.Payload.ID
}

func (c *customersContext) expectACustomerNamedToExist(name string) error {
	if !c.fetchedCustomers {
		err := c.fetchCustomers()
		if err != nil {
			return err
		}
	}

	if _, exists := c.customers[name]; !exists {
		return errors.ErrNotFound.Msgf("the customer `%s` does not exist", name)
	}
	return nil
}

func (c *customersContext) expectNoCustomerNamedToExist(name string) error {
	if !c.fetchedCustomers {
		err := c.fetchCustomers()
		if err != nil {
			return err
		}
	}

	if _, exists := c.customers[name]; exists {
		return errors.ErrNotFound.Msgf("the customer `%s` does exist", name)
	}
	return nil
}

func (c *customersContext) fetchCustomers() error {
	c.fetchedCustomers = true
	rows, err := c.db.Query("SELECT id, name FROM customers.customers")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		c.customers[name] = id
	}

	return nil
}
