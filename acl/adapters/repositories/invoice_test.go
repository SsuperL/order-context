package repositories

import (
	"fmt"
	"order-service/common"
	"order-service/domain/aggregate"
	"order-service/domain/entity"
	"order-service/domain/vo"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestCreateInvoice(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	siteCode := "001001"
	invoiceAg := &aggregate.InvoiceAggregate{
		RootID: common.RandomString(10),
		Invoice: &entity.Invoice{
			ID:     common.RandomString(10),
			Status: common.Invoiced,
			Path:   common.RandomString(10),
			Detail: vo.InvoiceDetail{
				Name: common.RandomString(6),
			},
		},
		Order: vo.Order{Price: float64(200)},
	}

	mock.ExpectBegin()
	execSQL := `INSERT INTO "invoices" ("id","order_id","price","name","code","status","path","site_code","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	mock.ExpectExec(regexp.QuoteMeta(execSQL)).WithArgs(invoiceAg.Invoice.ID, invoiceAg.RootID, invoiceAg.Order.Price, invoiceAg.Invoice.Detail.Name, AnyString{},
		invoiceAg.Invoice.Status, invoiceAg.Invoice.Path, siteCode, AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := invoiceAdapter.CreateInvoice(invoiceAg, siteCode)
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

func TestCreateInvoiceWithError(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	siteCode := "001001"
	invoiceAg := &aggregate.InvoiceAggregate{
		RootID: common.RandomString(10),
		Invoice: &entity.Invoice{
			ID:     common.RandomString(10),
			Status: common.UnInvoiced,
			Path:   common.RandomString(10),
			Detail: vo.InvoiceDetail{
				Name: common.RandomString(6),
			},
		},
		Order: vo.Order{Price: float64(200)},
	}

	mock.ExpectBegin()
	execSQL := `INSERT INTO "invoices" ("id","order_id","price","name","code","status","path","site_code","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	mock.ExpectExec(regexp.QuoteMeta(execSQL)).WithArgs(invoiceAg.Invoice.ID, invoiceAg.RootID, invoiceAg.Order.Price, invoiceAg.Invoice.Detail.Name, AnyString{},
		invoiceAg.Invoice.Status, invoiceAg.Invoice.Path, siteCode, AnyTime{}, AnyTime{}).WillReturnError(fmt.Errorf("create invoice failed"))
	mock.ExpectRollback()

	err := invoiceAdapter.CreateInvoice(invoiceAg, siteCode)
	require.Error(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

func TestGetInvoiceDetail(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	siteCode := "001001"
	invoiceID := "testId"

	query := `SELECT * FROM  "invoices" WHERE id = $1 AND site_code = $2`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"order_id",
			"price",
			"name",
			"code",
			"status",
			"path",
			"site_code",
			"created_at",
			"updated_at"}).
			AddRow(
				invoiceID,
				common.RandomString(10),
				100,
				common.RandomString(10),
				common.GenerateNumber(),
				common.Invoiced,
				common.RandomString(12),
				siteCode,
				time.Now(),
				time.Now()))

	order, err := invoiceAdapter.GetInvoiceDetail(invoiceID, siteCode)
	require.NoError(t, err)
	require.Equal(t, order.ID, invoiceID)
	require.Equal(t, order.SiteCode, siteCode)
}

func TestUpdateInvoice(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	invoiceID, siteCode := common.RandomString(10), "001001"
	status := common.Invoiced
	param := common.UpdateInvoiceParams{
		Status: status,
		Path:   common.RandomString(12),
	}

	mock.ExpectBegin()
	sql := `UPDATE "invoices" SET "path"=$1,"status"=$2,"updated_at"=$3 WHERE id = $4 AND site_code = $5`
	mock.ExpectExec(regexp.QuoteMeta(sql)).
		WithArgs(param.Path, param.Status, AnyTime{}, invoiceID, siteCode).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()
	err := invoiceAdapter.UpdateInvoice(invoiceID, siteCode, param)
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

func TestCheckInvoiceExists(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	invoiceID, siteCode := common.RandomString(10), "001001"
	query := `SELECT "id" FROM "invoices" WHERE id = $1 AND site_code = $2`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(invoiceID, siteCode).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(invoiceID))

	err := invoiceAdapter.CheckInvoiceExists(invoiceID, siteCode)
	require.NoError(t, err)
}

func TestGetInvoiceList(t *testing.T) {
	mock, gormDB := NewMockDB()
	invoiceAdapter := InvoiceAdapter{db: gormDB}
	params := common.ListInvoiceParams{
		OrderID: common.RandomString(10),
		Status:  common.UnInvoiced,
		Offset:  0,
		Limit:   1,
	}

	countQuery := `SELECT count(*) FROM "invoices" WHERE order_id = $1 AND status = $2`
	mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(sqlmock.NewRows([]string{}))
	query := `SELECT * FROM  "invoices" WHERE order_id = $1 AND status = $2`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(params.OrderID, params.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "price", "name", "code", "status", "path", "site_code", "created_at", "updated_at"}).
			AddRow(common.RandomString(10), common.RandomString(10), 100.00, common.RandomString(6), common.GenerateNumber(), common.Invoiced, common.RandomString(12), "001001", time.Now(), time.Now()))

	invoices, err := invoiceAdapter.GetInvoiceList(params)
	require.NoError(t, err)
	require.Len(t, invoices, 1)
	require.NotEmpty(t, invoices)
}
