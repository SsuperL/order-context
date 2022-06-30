package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"order-context/acl/adapters/pl"
	"order-context/acl/ports/clients"
	"order-context/utils/common"
	"time"
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
	url := common.FileConfig.UUIDSrv.HOST + "/uuid/generate/"

	// 设置超时
	req, err := http.NewRequest("GET", url, nil)
	ctx, cancel := context.WithTimeout(req.Context(), time.Millisecond*100)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := u.HTTPClient.Do(req)

	if err != nil {
		log.Printf("ERROR: uuid-service: Failed to get uuid: %v", err)
		return pl.UUIDRes{}, err
	}

	defer resp.Body.Close()

	// resp, err := u.HTTPClient.Get(url)
	// if err != nil {
	// 	log.Printf("ERROR: uuid-service: Fail to get uuid:%v\n", err)
	// 	return pl.UUIDRes{}, err
	// }
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("ERROR: uuid-service: Response code = %v", resp.StatusCode)
		err = fmt.Errorf("ERROR: uuid-service: Response code = %v", resp.StatusCode)
		return pl.UUIDRes{}, err
	}

	defer resp.Body.Close()
	var res pl.UUIDRes
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ERROR: Read http response failed: %v", err)
		return pl.UUIDRes{}, err
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		log.Printf("ERROR: Fail to unmarshal body: %v", err)
		return pl.UUIDRes{}, err
	}
	// fmt.Println(res)

	return res, nil
}
