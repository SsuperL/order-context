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
	"order-context/ohs/local/pl"
	ohs_pl "order-context/ohs/local/pl"
	"order-context/ohs/local/pl/errors"
	"order-context/utils/common"
	"order-context/utils/common/db"
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

func setUpSuite(tb testing.TB) func(tb testing.TB) {
	fmt.Println("setup suite-------")
	os.Setenv("DRIVER", "sqlite")
	dbInstance, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to database :%v", err)
	}

	db.InitTables(dbInstance)

	return func(tb testing.TB) {
		dbInstance.Migrator().DropTable("orders")
		dbInstance.Migrator().DropTable("invoices")
		fmt.Println("teardown suite-------")
	}
}

// TBD table driven tests
func TestCreateOrder(t *testing.T) {
	setupSuite := setUpSuite(t)
	defer setupSuite(t)

	testCases := []struct {
		name      string
		params    *ohs_pl.CreateOrderRequest
		expectErr error
		patch     func(string)
		code      codes.Code
	}{
		{
			name: "create order successfully",
			params: &ohs_pl.CreateOrderRequest{
				Status:  ohs_pl.OrderStatus(common.Unpaid),
				Price:   100,
				SpaceId: "test-space",
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
			name: "unpaid order exists",
			params: &ohs_pl.CreateOrderRequest{
				Status:  ohs_pl.OrderStatus(common.Unpaid),
				Price:   100,
				SpaceId: "test-space",
			},
			patch:     func(mockID string) {},
			expectErr: errors.UnpaidOrderExists("unpaid order exists"),
			code:      codes.AlreadyExists,
		},
		{
			name: "invalid argument",
			params: &ohs_pl.CreateOrderRequest{
				Status:  22,
				Price:   100,
				SpaceId: "test-invalid-argument",
			},
			patch:     func(mockID string) {},
			expectErr: errors.BadRequest("invalid argument"),
			code:      codes.InvalidArgument,
		},
		{
			name: "get uuid failed",
			params: &ohs_pl.CreateOrderRequest{
				Status:  ohs_pl.OrderStatus(common.Unpaid),
				Price:   100,
				SpaceId: "test-space1",
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
			serverAddress := startTestOrderServer(t)

			mockID := "1234"
			tc.patch(mockID)
			ctx := newContext()
			client := newTestOrderClient(t, serverAddress)
			res, err := client.CreateOrder(ctx, tc.params)
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

func TestGetOrderDetail(t *testing.T) {
	setupSuite := setUpSuite(t)
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
	resource := NewOrderResource()
	res, err := resource.CreateOrder(ctx, &ohs_pl.CreateOrderRequest{
		Status:  ohs_pl.OrderStatus(common.Unpaid),
		Price:   100,
		SpaceId: "test-space",
	})

	testCases := []struct {
		name      string
		params    *ohs_pl.GetOrderDetailRequest
		expectErr error
		code      codes.Code
	}{
		{
			name: "get order details success",
			params: &ohs_pl.GetOrderDetailRequest{
				Id: mockID,
			},
			expectErr: nil,
			code:      codes.OK,
		},
		{
			name: "order not found",
			params: &ohs_pl.GetOrderDetailRequest{
				Id: "test",
			},
			expectErr: errors.OrderNotFound("order not found"),
			code:      codes.NotFound,
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestOrderServer(t)

			ctx := newContext()
			client := newTestOrderClient(t, serverAddress)
			res, err := client.GetOrderDetail(ctx, tc.params)
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

func TestUpdateOrder(t *testing.T) {
	setupSuite := setUpSuite(t)
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
	resource := NewOrderResource()
	res, err := resource.CreateOrder(ctx, &ohs_pl.CreateOrderRequest{
		Status:  ohs_pl.OrderStatus(common.Unpaid),
		Price:   100,
		SpaceId: "test-space",
	})

	testCases := []struct {
		name      string
		params    *ohs_pl.UpdateOrderRequest
		expectErr error
		code      codes.Code
	}{
		{
			name: "update order successfully",
			params: &ohs_pl.UpdateOrderRequest{
				Id:     mockID,
				Status: ohs_pl.OrderStatus(common.Paid),
			},
			expectErr: nil,
			code:      codes.OK,
		},
		{
			name: "order not found",
			params: &ohs_pl.UpdateOrderRequest{
				Id: "test",
			},
			expectErr: errors.OrderNotFound("order not found"),
			code:      codes.NotFound,
		},
		{
			name: "invalid argument",
			params: &ohs_pl.UpdateOrderRequest{
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
			serverAddress := startTestOrderServer(t)

			ctx := newContext()
			client := newTestOrderClient(t, serverAddress)
			res, err := client.UpdateOrder(ctx, tc.params)
			if tc.code == codes.OK {
				require.NotEmpty(t, res)
				require.True(t, res.Success)
				require.Nil(t, err)
				res, err := client.GetOrderDetail(ctx, &pl.GetOrderDetailRequest{Id: tc.params.Id})
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

func TestGetOrderList(t *testing.T) {
	setupSuite := setUpSuite(t)
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
	resource := NewOrderResource()
	res, err := resource.CreateOrder(ctx, &ohs_pl.CreateOrderRequest{
		Status:  ohs_pl.OrderStatus(common.Unpaid),
		Price:   100,
		SpaceId: "test-space",
	})

	testCases := []struct {
		name      string
		params    *ohs_pl.GetOrderListRequest
		total     int
		expectErr error
		code      codes.Code
	}{
		{
			name: "get order list successfully",
			params: &ohs_pl.GetOrderListRequest{
				SpaceId: "test-space",
			},
			total:     1,
			expectErr: nil,
			code:      codes.OK,
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			serverAddress := startTestOrderServer(t)

			ctx := newContext()
			client := newTestOrderClient(t, serverAddress)
			res, err := client.GetOrderList(ctx, tc.params)
			if tc.code == codes.OK {
				require.NotEmpty(t, res)

				require.Nil(t, err)
				res, err := client.GetOrderList(ctx, tc.params)
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

func startTestOrderServer(t *testing.T) string {
	orderServer := NewOrderResource()
	grpcServer := grpc.NewServer()
	ohs_pl.RegisterOrderServiceServer(grpcServer, orderServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
}

func newTestOrderClient(t *testing.T, serverAddress string) ohs_pl.OrderServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return ohs_pl.NewOrderServiceClient(conn)
}

func newContext() context.Context {
	siteCode := "001001"
	md := metadata.Pairs("site-code", siteCode)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}
