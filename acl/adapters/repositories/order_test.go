package repositories

import (
	"database/sql/driver"
	"fmt"
	"order-context/common"
	"order-context/domain/aggregate"
	"order-context/domain/entity"
	"order-context/domain/vo"
	ohs_pl "order-context/ohs/local/pl"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// AnyTime mock time.Time
type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

// AnyString mock string
type AnyString struct{}

func (s AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}

func NewMockDB() (sqlmock.Sqlmock, *gorm.DB) {
	sqlDB, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true,
		Conn:                 sqlDB,
	}), &gorm.Config{})
	return mock, gormDB
}

func TestCreateOrder(t *testing.T) {
	mock, gormDB := NewMockDB()
	orderAdapter := OrderAdapter{db: gormDB}
	siteCode := "001001"
	root := &aggregate.AggregateRoot{
		Order: &entity.Order{
			ID:     common.RandomString(10),
			Status: common.Unpaid,
			Price:  float32(200),
		},
		Package: vo.Package{
			Version: "v1",
			Price:   float32(100),
		},
		Space: &entity.Space{
			ID: common.RandomString(10),
		},
	}

	mock.ExpectBegin()
	execSQL := `INSERT INTO "orders" ("id","status","number","space_id","pay_id","price","package_version","package_price","site_code","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	mock.ExpectExec(regexp.QuoteMeta(execSQL)).WithArgs(root.Order.ID, root.GetStatus(), AnyString{}, root.Space.ID, root.Payment.Voucher,
		root.Order.Price, root.Package.Version, root.Package.Price, siteCode,
		AnyTime{}, AnyTime{}).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err := orderAdapter.CreateOrder(root, siteCode)
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

func TestCreateOrderWithError(t *testing.T) {
	mock, gormDB := NewMockDB()
	orderAdapter := OrderAdapter{db: gormDB}
	siteCode := "001001"
	root := &aggregate.AggregateRoot{
		Order: &entity.Order{
			ID:     common.RandomString(10),
			Status: common.Unpaid,
			Price:  float32(200),
		},
		Package: vo.Package{
			Version: "v1",
			Price:   float32(100),
		},
		Space: &entity.Space{
			ID: common.RandomString(10),
		},
	}

	mock.ExpectBegin()
	execSQL := `INSERT INTO "orders" ("id","status","number","space_id","pay_id","price","package_version","package_price","site_code","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	mock.ExpectExec(regexp.QuoteMeta(execSQL)).WithArgs(root.Order.ID, root.GetStatus(), AnyString{}, root.Space.ID, root.Payment.Voucher,
		root.Order.Price, root.Package.Version, root.Package.Price, siteCode,
		AnyTime{}, AnyTime{}).WillReturnError(fmt.Errorf("failed to create order"))

	mock.ExpectRollback()

	err := orderAdapter.CreateOrder(root, siteCode)
	require.Error(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

// func TestGetOrderDetail(t *testing.T) {
// 	mock, gormDB := NewMockDB()
// 	orderAdapter := OrderAdapter{db: gormDB}
// 	siteCode := "001001"
// 	orderID := "testId"

// 	query := `SELECT * FROM  "orders" WHERE id = $1 AND site_code = $2`
// 	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(sqlmock.NewRows([]string{"id", "status", "number", "space_id",
// 		"pay_id", "price", "package_version", "package_price", "site_code", "created_at", "updated_at"}).AddRow(
// 		orderID, common.Unpaid, common.GenerateNumber(), common.RandomString(10), common.RandomString(6), 100, "v1", 100, siteCode, time.Now(), time.Now()))

// 	order, err := orderAdapter.GetOrderDetail(orderID, siteCode)
// 	require.NoError(t, err)
// 	require.Equal(t, order.ID, orderID)
// 	require.Equal(t, order.SiteCode, siteCode)
// }

func TestUpdateOrder(t *testing.T) {
	mock, gormDB := NewMockDB()
	orderAdapter := OrderAdapter{db: gormDB}
	orderID, siteCode := common.RandomString(10), "001001"
	status := common.Paid

	mock.ExpectBegin()
	sql := `UPDATE "orders" SET "status"=$1,"updated_at"=$2 WHERE id = $3 AND site_code = $4`
	mock.ExpectExec(regexp.QuoteMeta(sql)).WithArgs(status, AnyTime{}, orderID, siteCode).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()
	err := orderAdapter.UpdateOrderStatus(orderID, siteCode, status)
	require.NoError(t, err)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("Failed to meet expectations,got err:%v", err)
	}
}

func TestCheckOrderExists(t *testing.T) {
	mock, gormDB := NewMockDB()
	orderAdapter := OrderAdapter{db: gormDB}
	orderID, siteCode := common.RandomString(10), "001001"
	query := `SELECT "id" FROM "orders" WHERE id = $1 AND site_code = $2`

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(orderID, siteCode).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(orderID))

	ok, err := orderAdapter.CheckOrderExists(orderID, siteCode)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestGetOrderList(t *testing.T) {
	mock, gormDB := NewMockDB()
	orderAdapter := OrderAdapter{db: gormDB}
	params := ohs_pl.ListOrderParams{
		SpaceID: common.RandomString(10),
		Status:  common.Unpaid,
		Offset:  0,
		Limit:   1,
	}

	countQuery := `SELECT count(*) FROM "orders" WHERE space_id = $1 AND status = $2`
	mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	query := `SELECT * FROM  "orders" WHERE space_id = $1 AND status = $2`
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(params.SpaceID, params.Status).
		WillReturnRows(sqlmock.NewRows([]string{"id", "status", "number", "space_id",
			"pay_id", "price", "package_version", "package_price", "site_code", "created_at", "updated_at"}).AddRow(
			common.RandomString(10), common.Unpaid, common.GenerateNumber(), common.RandomString(10), common.RandomString(6), 100, "v1", 100, "001001", time.Now(), time.Now()))

	orders, total, err := orderAdapter.GetOrderList(params)
	require.Equal(t, len(orders), total)
	require.NoError(t, err)
	// require.Len(t, orders, total)
	require.NotEmpty(t, orders)
}
