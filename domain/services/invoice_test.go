package services

import (
	"fmt"
	"order-context/acl/adapters/pl"
	"order-context/domain/aggregate"
	msg "order-context/ohs/local/pl"
	"order-context/utils/common"
	"order-context/utils/mock"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createInvoiceAggregate(rootID, id string, options ...aggregate.InvoiceOptions) *aggregate.InvoiceAggregate {
	return aggregate.NewInvoiceAggregate(rootID, id, options...)
}

func TestCreateInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrderRepo := mock.NewMockOrderRepository(ctrl)
	mockUUIDClient := mock.NewMockUUIDClient(ctrl)
	mockInvoiceRepo := mock.NewMockInvoiceRepository(ctrl)
	mockID := "1234"

	invoiceOption := aggregate.WithInvoiceOption(common.UnInvoiced, "", "")
	invoiceAg := createInvoiceAggregate("root_id", "", invoiceOption)
	isv := WithPortAndClientParams(invoiceAg, mockOrderRepo, mockInvoiceRepo, mockUUIDClient)
	siteCode := "001001"

	// order exists
	mockOrderRepo.EXPECT().CheckOrderExists(invoiceAg.RootID, siteCode).Return(true, nil)
	mockUUIDClient.EXPECT().GetUUID(1).Return(pl.UUIDRes{ID: mockID}, nil)
	mockInvoiceRepo.EXPECT().CreateInvoice(isv.InvoiceAg, siteCode).Return(nil)

	invoiceID, err := isv.CreateInvoice(siteCode)
	require.NoError(t, err)
	require.Equal(t, invoiceID, mockID)

	// order not exists
	mockOrderRepo.EXPECT().CheckOrderExists(invoiceAg.RootID, siteCode).Return(false, fmt.Errorf("order not exists"))
	invoiceID, err = isv.CreateInvoice(siteCode)
	require.Error(t, err)
	require.Empty(t, invoiceID)
}

func TestUpdateInvoice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvoiceRepo := mock.NewMockInvoiceRepository(ctrl)

	invoiceOption := aggregate.WithInvoiceOption(common.UnInvoiced, "", "")
	invoiceAg := createInvoiceAggregate("root_id", "", invoiceOption)
	isv := WithPortAndClientParams(invoiceAg, nil, mockInvoiceRepo, nil)
	siteCode := "001001"

	params := msg.UpdateInvoiceParams{
		Status: common.Invoiced,
	}

	// invoice exists
	mockInvoiceRepo.EXPECT().CheckInvoiceExists(invoiceAg.GetID(), siteCode).Return(nil)
	mockInvoiceRepo.EXPECT().UpdateInvoice(invoiceAg.GetID(), siteCode, params).Return(nil)

	err := isv.UpdateInvoice(siteCode, params)
	require.NoError(t, err)

	// invoice not found
	mockInvoiceRepo.EXPECT().CheckInvoiceExists(invoiceAg.GetID(), siteCode).Return(gorm.ErrRecordNotFound)
	err = isv.UpdateInvoice(siteCode, params)
	require.Error(t, err)
}

// func TestGetInvoiceDetail(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockInvoiceRepo := mock.NewMockInvoiceRepository(ctrl)

// 	invoiceAg := createInvoiceAggregate("test", "", "", common.UnInvoiced, 100)
// 	isv := WithPortAndClientParams(invoiceAg, nil, mockInvoiceRepo, nil)
// 	siteCode := "001001"

// 	mockInvoiceRepo.EXPECT().GetInvoiceDetail(isv.InvoiceAg.GetID(), siteCode).Return(pl.Invoice{ID: isv.InvoiceAg.GetID()}, nil)
// 	invoice, err := isv.GetInvoiceDetail(siteCode)
// 	require.NoError(t, err)
// 	require.Equal(t, invoice.ID, isv.InvoiceAg.GetID())

// 	// order not found
// 	mockInvoiceRepo.EXPECT().GetInvoiceDetail(isv.InvoiceAg.GetID(), siteCode).Return(pl.Invoice{}, gorm.ErrRecordNotFound)
// 	invoice, err = isv.GetInvoiceDetail(siteCode)
// 	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
// 	require.Empty(t, invoice)
// }

// func TestGetInvoiceList(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockInvoiceRepo := mock.NewMockInvoiceRepository(ctrl)

// 	invoiceAg := createInvoiceAggregate("test", "", "", common.UnInvoiced, 100)
// 	isv := WithPortAndClientParams(invoiceAg, nil, mockInvoiceRepo, nil)
// 	siteCode := "001001"

// 	args := ohs_pl.ListInvoiceParams{
// 		Status: common.UnInvoiced,
// 	}
// 	mockInvoices := []pl.Invoice{{ID: isv.InvoiceAg.GetID(), SiteCode: siteCode}}
// 	mockInvoiceRepo.EXPECT().GetInvoiceList(args).Return(mockInvoices, 1, nil)

// 	invoices, total, err := isv.GetInvoiceList(args)
// 	require.NoError(t, err)
// 	require.Equal(t, invoices, mockInvoices)
// 	require.Len(t, invoices, total)
// }
