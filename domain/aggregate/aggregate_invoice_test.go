package aggregate

import (
	"order-context/utils/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomInvoiceAggregate(rootID, id string, options ...InvoiceOptions) *InvoiceAggregate {
	return NewInvoiceAggregate(rootID, id, options...)
}

func TestSetInvoiceID(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("", "invoice1")
	newID := common.RandomString(10)
	invoiceAg.SetID(newID)
	require.Equal(t, invoiceAg.Invoice.ID, newID)
}

func TestGetInvoiceID(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("", "invoice1")
	id := invoiceAg.GetID()
	require.Equal(t, id, invoiceAg.Invoice.ID)
}

func TestGetInvoiceStatus(t *testing.T) {
	invoiceOption := WithInvoiceOption(common.UnInvoiced, "", "")
	invoiceAg := randomInvoiceAggregate("", "invoice1", invoiceOption)
	status := invoiceAg.GetStatus()
	require.Equal(t, status, invoiceAg.Invoice.Status)
}

func TestSetInvoiceStatus(t *testing.T) {
	invoiceOption := WithInvoiceOption(common.UnInvoiced, "", "")
	invoiceAg := randomInvoiceAggregate("", "invoice1", invoiceOption)
	newStatus := common.Invoiced
	invoiceAg.SetStatus(newStatus)
	require.Equal(t, newStatus, invoiceAg.Invoice.Status)
}

func TestGetPath(t *testing.T) {
	invoiceOption := WithInvoiceOption(common.UnInvoiced, "/path", "")
	invoiceAg := randomInvoiceAggregate("", "invoice1", invoiceOption)
	path := invoiceAg.GetPath()
	require.Equal(t, path, invoiceAg.Invoice.Path)
}

func TestSetInvoicePath(t *testing.T) {
	invoiceOption := WithInvoiceOption(common.UnInvoiced, "/path", "")
	invoiceAg := randomInvoiceAggregate("", "invoice1", invoiceOption)
	newPath := "/"
	invoiceAg.SetPath(newPath)
	require.Equal(t, newPath, invoiceAg.Invoice.Path)
}
