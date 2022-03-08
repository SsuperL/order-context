package aggregate

import (
	"order-service/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetID(t *testing.T) {
	root := createAggregateRoot("", "", "", "", common.Unpaid, 20.5, 42)
	newID := common.RandomString(10)
	root.SetID(newID)
	require.Equal(t, root.Order.ID, newID)
}

func TestGetStatus(t *testing.T) {
	root := createAggregateRoot("", "", "", "", common.Unpaid, 0, 0)
	status := root.GetStatus()
	require.Equal(t, root.Order.Status, status)
}

func TestSetStatus(t *testing.T) {
	root := createAggregateRoot("", "", "", "", common.Unpaid, 0, 0)
	newStatus := common.Paid
	root.SetStatus(newStatus)
	require.Equal(t, root.Order.Status, newStatus)
}

func TestGetID(t *testing.T) {
	root := createAggregateRoot("", "", "", "", common.Unpaid, 20.5, 42)
	id := root.GetID()
	require.Equal(t, id, root.Order.ID)
}

func createAggregateRoot(voucher, source, version, currency string, status common.StatusType, price, total float64) *AggregateRoot {
	return NewOrderAggregateRoot(common.RandomString(10), common.RandomString(10), voucher, source, version, currency, status, price, total)
}
