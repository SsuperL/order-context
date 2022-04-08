package clients

import (
	"net/http"
	"order-context/utils/common"
	"testing"

	"github.com/nbio/st"
	"github.com/stretchr/testify/require"
	"gopkg.in/h2non/gock.v1"
)

func TestGetUUID(t *testing.T) {
	defer gock.Off()

	mockID := common.RandomString(10)
	gock.New(common.FileConfig.UUIDSrv.HOST).
		Get("/uuid/generate/").
		Reply(200).
		JSON(map[string]string{"id": mockID})

	// mock http client
	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)

	uuidClient := UUIDAdapter{HTTPClient: client}
	res, err := uuidClient.GetUUID(1)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Equal(t, res.ID, mockID)
	require.Len(t, res.IDList, 0)
	st.Expect(t, gock.IsDone(), true)

	// limit > 1
	gock.New(common.FileConfig.UUIDSrv.HOST).
		Get("/uuid/generate/").
		Reply(200).
		JSON(map[string][]string{"id_list": {mockID}})

	res, err = uuidClient.GetUUID(20)
	require.NoError(t, err)
	require.NotEmpty(t, res)
	require.Empty(t, res.ID)
	require.Len(t, res.IDList, 1)
}
