package aggregate

import (
	"order-context/common"
	"testing"

	"github.com/stretchr/testify/require"
)

func createAggregateRoot(id string, options ...RootOptions) *AggregateRoot {
	return NewOrderAggregateRoot(id, options...)
}

func TestSetID(t *testing.T) {
	root := createAggregateRoot("")
	newID := common.RandomString(10)
	root.SetID(newID)
	require.Equal(t, root.Order.ID, newID)
}

func TestGetStatus(t *testing.T) {
	orderOption := WithOrderOption(common.Unpaid, 0)
	root := createAggregateRoot("", orderOption)
	status := root.GetStatus()
	require.Equal(t, root.Order.Status, status)
}

func TestSetStatus(t *testing.T) {
	orderOption := WithOrderOption(common.Unpaid, 0)
	root := createAggregateRoot("", orderOption)
	newStatus := common.Paid
	root.SetStatus(newStatus)
	require.Equal(t, root.Order.Status, newStatus)
}

func TestGetID(t *testing.T) {
	orderOption := WithOrderOption(common.Unpaid, 0)
	root := createAggregateRoot("", orderOption)
	id := root.GetID()
	require.Equal(t, id, root.Order.ID)
}
