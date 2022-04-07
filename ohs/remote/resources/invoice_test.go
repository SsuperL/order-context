package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	a_pl "order-context/acl/adapters/pl"
	"order-context/common"
	"order-context/common/db"
	"order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"os"
	"reflect"
	"testing"

	err "errors"

	"bou.ke/monkey"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func setUpInvoiceSuite(tb testing.TB) (string, func(tb testing.TB)) {
	fmt.Println("setup suite-------")
	os.Setenv("DRIVER", "sqlite")
	dbInstance, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database :%v", err)
	}

	db.InitTables(dbInstance)
	order := &a_pl.Order{ID: common.RandomString(10), SiteCode: "001001"}
	dbInstance.Create(&order)

	return order.ID, func(tb testing.TB) {
		dbInstance.Migrator().DropTable("orders")
		dbInstance.Migrator().DropTable("invoices")
		fmt.Println("teardown suite-------")
	}
}

// TBD table driven tests
func TestCreateInvoice(t *testing.T) {
	orderID, setupSuite := setUpInvoiceSuite(t)
	defer setupSuite(t)

	testCases := []struct {
		name      string
		params    *pl.CreateInvoiceRequest
		expectErr error
		patch     func(string)
		code      codes.Code
	}{
		{
			name: "create Invoice successfully",
			params: &pl.CreateInvoiceRequest{
				Status:  pl.InvoiceStatus(common.UnInvoiced),
				Price:   100,
				OrderId: orderID,
			},
			patch: func(mockID string) {
				body := map[string]string{"id": mockID}
				jsonBytes, _ := json.Marshal(body)
				r := ioutil.NopCloser(bytes.NewReader(jsonBytes))
				monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       r,
					}, nil
				})
			},
			expectErr: nil,
			code:      codes.OK,
		},
		{
			name: "invalid argument",
			params: &pl.CreateInvoiceRequest{
				Status:  22,
				Price:   100,
				OrderId: "test-invalid-argument",
			},
			patch:     func(mockID string) {},
			expectErr: errors.BadRequest("invalid argument"),
			code:      codes.InvalidArgument,
		},
		{
			name: "get uuid failed",
			params: &pl.CreateInvoiceRequest{
				Status:  pl.InvoiceStatus(common.UnInvoiced),
				Price:   100,
				OrderId: orderID,
			},
			patch: func(mockID string) {
				monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
					return &http.Response{
						StatusCode: 500,
					}, err.New("failed to generate uuid")
				})
			},
			expectErr: errors.InternalServerError("get uuid failed"),
			code:      codes.Internal,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestInvoiceServer(t)

			mockID := "1234"
			tc.patch(mockID)
			ctx := newContext()
			client := newTestInvoiceClient(t, serverAddress)
			res, err := client.CreateInvoice(ctx, tc.params)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, res.Id, mockID)
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}

}

func TestGetInvoiceDetail(t *testing.T) {
	orderID, setupSuite := setUpInvoiceSuite(t)
	defer setupSuite(t)

	mockID := "1234"
	body := map[string]string{"id": mockID}
	jsonBytes, _ := json.Marshal(body)
	r := ioutil.NopCloser(bytes.NewReader(jsonBytes))
	monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	})

	md := metadata.Pairs("site-code", "001001")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	resource := NewInvoiceResource()
	res, err := resource.CreateInvoice(ctx, &pl.CreateInvoiceRequest{
		Status:  pl.InvoiceStatus(common.UnInvoiced),
		Price:   100,
		OrderId: orderID,
	})

	testCases := []struct {
		name      string
		params    *pl.GetInvoiceDetailRequest
		expectErr error
		code      codes.Code
	}{
		{
			name: "get Invoice details success",
			params: &pl.GetInvoiceDetailRequest{
				Id: mockID,
			},
			expectErr: nil,
			code:      codes.OK,
		},
		{
			name: "Invoice not found",
			params: &pl.GetInvoiceDetailRequest{
				Id: "test",
			},
			expectErr: errors.InvoiceNotFound("Invoice not found"),
			code:      codes.NotFound,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestInvoiceServer(t)

			ctx := newContext()
			client := newTestInvoiceClient(t, serverAddress)
			res, err := client.GetInvoiceDetail(ctx, tc.params)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotEmpty(t, res)
				require.Equal(t, res.Result.Id, tc.params.Id)
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}

	require.NoError(t, err)
	require.Equal(t, res.Id, mockID)
}

func TestUpdateInvoice(t *testing.T) {
	orderID, setupSuite := setUpInvoiceSuite(t)
	defer setupSuite(t)

	mockID := "12345"
	body := map[string]string{"id": mockID}
	jsonBytes, _ := json.Marshal(body)
	r := ioutil.NopCloser(bytes.NewReader(jsonBytes))
	monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	})

	md := metadata.Pairs("site-code", "001001")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	resource := NewInvoiceResource()
	res, err := resource.CreateInvoice(ctx, &pl.CreateInvoiceRequest{
		Status:  pl.InvoiceStatus(common.UnInvoiced),
		Price:   100,
		OrderId: orderID,
	})

	testCases := []struct {
		name      string
		params    *pl.UpdateInvoiceRequest
		expectErr error
		code      codes.Code
	}{
		{
			name: "update Invoice successfully",
			params: &pl.UpdateInvoiceRequest{
				Id:     mockID,
				Status: pl.InvoiceStatus(common.Paid),
			},
			expectErr: nil,
			code:      codes.OK,
		},
		{
			name: "Invoice not found",
			params: &pl.UpdateInvoiceRequest{
				Id: "test",
			},
			expectErr: errors.InvoiceNotFound("Invoice not found"),
			code:      codes.NotFound,
		},
		{
			name: "invalid argument",
			params: &pl.UpdateInvoiceRequest{
				Id:     mockID,
				Status: 111,
			},
			expectErr: errors.BadRequest("invalid argument"),
			code:      codes.InvalidArgument,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestInvoiceServer(t)

			ctx := newContext()
			client := newTestInvoiceClient(t, serverAddress)
			res, err := client.UpdateInvoice(ctx, tc.params)
			if tc.code == codes.OK {
				require.NotEmpty(t, res)
				require.True(t, res.Success)
				require.Nil(t, err)
				res, err := client.GetInvoiceDetail(ctx, &pl.GetInvoiceDetailRequest{Id: tc.params.Id})
				require.Nil(t, err)
				require.Equal(t, res.Result.Status, tc.params.Status)
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}

	require.NoError(t, err)
	require.Equal(t, res.Id, mockID)
}

func TestGetInvoiceList(t *testing.T) {
	orderID, setupSuite := setUpInvoiceSuite(t)
	defer setupSuite(t)

	mockID := "1234"
	body := map[string]string{"id": mockID}
	jsonBytes, _ := json.Marshal(body)
	r := ioutil.NopCloser(bytes.NewReader(jsonBytes))
	monkey.PatchInstanceMethod(reflect.TypeOf(http.DefaultClient), "Get", func(c *http.Client, url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	})

	md := metadata.Pairs("site-code", "001001")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	resource := NewInvoiceResource()
	res, err := resource.CreateInvoice(ctx, &pl.CreateInvoiceRequest{
		Status:  pl.InvoiceStatus(common.UnInvoiced),
		Price:   100,
		OrderId: orderID,
	})

	testCases := []struct {
		name      string
		params    *pl.GetInvoiceListRequest
		total     int
		expectErr error
		code      codes.Code
	}{
		{
			name: "get Invoice list successfully",
			params: &pl.GetInvoiceListRequest{
				OrderId: orderID,
			},
			total:     1,
			expectErr: nil,
			code:      codes.OK,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestInvoiceServer(t)

			ctx := newContext()
			client := newTestInvoiceClient(t, serverAddress)
			res, err := client.GetInvoiceList(ctx, tc.params)
			if tc.code == codes.OK {
				require.NotEmpty(t, res)

				require.Nil(t, err)
				res, err := client.GetInvoiceList(ctx, tc.params)
				require.Nil(t, err)
				require.Equal(t, tc.total, int(res.Total))
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}

	require.NoError(t, err)
	require.Equal(t, res.Id, mockID)
}

func startTestInvoiceServer(t *testing.T) string {
	InvoiceServer := NewInvoiceResource()
	grpcServer := grpc.NewServer()
	pl.RegisterInvoiceServiceServer(grpcServer, InvoiceServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestInvoiceClient(t *testing.T, serverAddress string) pl.InvoiceServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pl.NewInvoiceServiceClient(conn)
}
