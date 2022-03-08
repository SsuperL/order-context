package aggregate

import (
	"order-service/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetInvoiceID(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "", common.UnInvoiced, 0)
	newID := common.RandomString(10)
	invoiceAg.SetID(newID)
	require.Equal(t, invoiceAg.Invoice.ID, newID)
}

func TestGetInvoiceID(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "", common.UnInvoiced, 0)
	id := invoiceAg.GetID()
	require.Equal(t, id, invoiceAg.Invoice.ID)
}

func TestGetInvoiceStatus(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "", common.UnInvoiced, 0)
	status := invoiceAg.GetStatus()
	require.Equal(t, status, invoiceAg.Invoice.Status)
}

func TestSetInvoiceStatus(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "", common.UnInvoiced, 0)
	newStatus := common.Invoiced
	invoiceAg.SetStatus(newStatus)
	require.Equal(t, newStatus, invoiceAg.Invoice.Status)
}

func TestGetPath(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "/", common.UnInvoiced, 0)
	path := invoiceAg.GetPath()
	require.Equal(t, path, invoiceAg.Invoice.Path)
}

func TestSetInvoicePath(t *testing.T) {
	invoiceAg := randomInvoiceAggregate("invoice1", "", "", common.UnInvoiced, 0)
	newPath := "/"
	invoiceAg.SetPath(newPath)
	require.Equal(t, newPath, invoiceAg.Invoice.Path)
}

func randomInvoiceAggregate(name, code, path string, status common.InvoiceStatusType, price float64) *InvoiceAggregate {
	return NewInvoiceAggregate(common.RandomString(10), common.RandomString(10), name, code, path, status, price)
}
