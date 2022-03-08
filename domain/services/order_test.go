package services

import (
	"order-service/acl/adapters/pl"
	"order-service/common"
	"order-service/domain/aggregate"
	"order-service/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createAggregateRoot(voucher, source, version, currency string, status common.StatusType, price, total float64) *aggregate.AggregateRoot {
	return aggregate.NewOrderAggregateRoot(common.RandomString(10), common.RandomString(10), voucher, source, version, currency, status, price, total)
}

func TestCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	// 断言被调用
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockUUIDClient := mock.NewMockUUIDClient(ctrl)
	mockID := "1234"
	// 不存在未支付订单
	mockOrderRepo.EXPECT().GetOrderList(gomock.Any()).Return([]pl.Order{}, nil)
	mockUUIDClient.EXPECT().GetUUID(1).Return(pl.UUIDRes{ID: mockID}, nil)

	root := createAggregateRoot("", "", "", "USD", common.Unpaid, 100, 100)
	siteCode := "001001"
	osv := OrderService{
		Port:       mockOrderRepo,
		Order:      root,
		UUIDClient: mockUUIDClient,
	}

	mockOrderRepo.EXPECT().CreateOrder(osv.Order, siteCode).Return(nil)

	orderID, err := osv.CreateOrder(siteCode)
	require.NoError(t, err)
	require.Equal(t, orderID, mockID)

	d, _ := time.ParseDuration("-25h")
	// 存在未支付订单
	// expired
	mockOrderRepo.EXPECT().GetOrderList(gomock.Any()).
		Return([]pl.Order{{
			ID:        mockID,
			SpaceID:   root.Space.ID,
			CreatedAt: time.Now().Add(d),
		}}, nil)

	mockOrderRepo.EXPECT().CheckOrderExists(root.Order.ID, siteCode).Times(1).Return(true, nil)
	mockOrderRepo.EXPECT().UpdateOrderStatus(osv.Order.GetID(), siteCode, common.Failed).Times(1).Return(nil)
	mockUUIDClient.EXPECT().GetUUID(1).Return(pl.UUIDRes{ID: mockID}, nil)
	mockOrderRepo.EXPECT().CreateOrder(osv.Order, siteCode).Return(nil)

	orderID, err = osv.CreateOrder(siteCode)
	require.NoError(t, err)
	require.Equal(t, orderID, mockID)

	// 存在未支付订单，未过期
	d, _ = time.ParseDuration("-2h")
	// 存在未支付订单
	// 已过期
	mockOrderRepo.EXPECT().GetOrderList(gomock.Any()).
		Return([]pl.Order{{
			ID:        mockID,
			SpaceID:   root.Space.ID,
			CreatedAt: time.Now().Add(d),
		}}, nil)

	orderID, err = osv.CreateOrder(siteCode)
	// TODO: 具体错误
	require.Error(t, err)
	require.Empty(t, orderID)
}

func TestGetOrderDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)

	root := createAggregateRoot("", "", "", "USD", common.Unpaid, 100, 100)
	siteCode := "001001"
	osv := OrderService{
		Port:  mockOrderRepo,
		Order: root,
	}

	mockOrderRepo.EXPECT().GetOrderDetail(osv.Order.GetID(), siteCode).Return(pl.Order{ID: osv.Order.GetID()}, nil)
	order, err := osv.GetOrderDetail(siteCode)
	require.NoError(t, err)
	require.Equal(t, order.ID, osv.Order.GetID())

	// order not found
	mockOrderRepo.EXPECT().GetOrderDetail(osv.Order.GetID(), siteCode).Return(pl.Order{}, gorm.ErrRecordNotFound)
	order, err = osv.GetOrderDetail(siteCode)
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	require.Empty(t, order)
}

func TestGetOrderList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	root := createAggregateRoot("", "", "", "USD", common.Unpaid, 100, 100)
	osv := OrderService{
		Port:  mockOrderRepo,
		Order: root,
	}

	// TODO: 参数位置调整
	args := common.ListOrderParams{
		SpaceID: root.Space.ID,
	}
	mockOrders := []pl.Order{{ID: root.Order.ID, SpaceID: root.Space.ID}}
	mockOrderRepo.EXPECT().GetOrderList(args).Return(mockOrders, nil)

	orders, err := osv.GetOrderList(args)
	require.NoError(t, err)
	require.Equal(t, orders, mockOrders)
	require.Len(t, orders, 1)
}

func TestUpdateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	root := createAggregateRoot("", "", "", "USD", common.Unpaid, 100, 100)
	osv := OrderService{
		Port:  mockOrderRepo,
		Order: root,
	}
	siteCode := "001001"
	status := common.Paid

	mockOrderRepo.EXPECT().CheckOrderExists(osv.Order.GetID(), siteCode).Return(true, nil)
	mockOrderRepo.EXPECT().UpdateOrderStatus(osv.Order.GetID(), siteCode, status).Return(nil)

	err := osv.UpdateOrderStatus(siteCode, status)
	require.NoError(t, err)

	// order not found
	mockOrderRepo.EXPECT().CheckOrderExists(osv.Order.GetID(), siteCode).Return(false, gorm.ErrRecordNotFound)
	err = osv.UpdateOrderStatus(siteCode, status)
	require.Error(t, err)
}
