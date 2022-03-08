package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order-service/acl/adapters/pl"
	"order-service/acl/ports/clients"
	"order-service/common"
)

// UUIDAdapter UUID适配器
type UUIDAdapter struct {
	HTTPClient *http.Client
}

var _ clients.UUIDClient = (*UUIDAdapter)(nil)

// NewUUIDAdapter UUID适配器构造函数
func NewUUIDAdapter() clients.UUIDClient {
	return &UUIDAdapter{
		HTTPClient: &http.Client{},
	}
}

// GetUUID 获取uuid
func (u *UUIDAdapter) GetUUID(limit int) (pl.UUIDRes, error) {
	url := common.UUIDurl + "/uuid/generate/"
	resp, err := u.HTTPClient.Get(url)
	if err != nil {
		log.Fatalf("ERROR: uuid-service: Fail to get uuid:%v\n", err)
		return pl.UUIDRes{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("ERROR: uuid-service: Response code = %v", resp.StatusCode)
		err = fmt.Errorf("ERROR: uuid-service: Response code = %v", resp.StatusCode)
		return pl.UUIDRes{}, err
	}

	defer resp.Body.Close()
	var res pl.UUIDRes
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ERROR: Read http response failed: %v", err)
		return pl.UUIDRes{}, err
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Fatalf("ERROR: Fail to unmarshal body: %v", err)
		return pl.UUIDRes{}, err
	}
	fmt.Println(res)

	return res, nil
}
