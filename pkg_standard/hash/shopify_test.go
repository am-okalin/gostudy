package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func Test4(t *testing.T) {
	h := sha256.New()
	h.Write([]byte("abcde123"))
	b := h.Sum(nil)
	s := hex.EncodeToString(b)
	t.Log(s)
}

func Test3(t *testing.T) {
	secret := "your-secret"
	body := `{"id":7507228688576,"title":"biz_test","body_html":"业务系统测试用, 任何人都可修改..","vendor":"REDMAGIC-Dev","product_type":"","created_at":"2023-08-04T04:04:14-04:00","handle":"biz_test","updated_at":"2023-08-04T07:55:53-04:00","published_at":"2023-08-04T04:04:14-04:00","template_suffix":"","status":"active","published_scope":"web","tags":"","admin_graphql_api_id":"gid://shopify/Product/7507228688576","variants":[{"id":43770570506432,"product_id":7507228688576,"title":"Default Title","price":"100.00","sku":"123","position":1,"inventory_policy":"deny","compare_at_price":"109.00","fulfillment_service":"manual","inventory_management":null,"option1":"Default Title","option2":null,"option3":null,"created_at":"2023-08-04T04:04:14-04:00","updated_at":"2023-08-04T04:04:14-04:00","taxable":true,"barcode":"","grams":0,"image_id":null,"weight":0,"weight_unit":"kg","inventory_item_id":45873500356800,"inventory_quantity":0,"old_inventory_quantity":0,"tax_code":"1","requires_shipping":false,"admin_graphql_api_id":"gid://shopify/ProductVariant/43770570506432"}],"options":[{"id":9635892035776,"product_id":7507228688576,"name":"Title","position":1,"values":["Default Title"]}],"images":[],"image":null}`
	str, err := VerifyWebhook(secret, []byte(body))
	t.Log(str, err)
	// dWTr5BZQOlxrikm9+Q8gJeBbZ163crpffJy6KVGyuYY=
}

func VerifyWebhook(secret string, body []byte) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(body)
	if err != nil {
		return "", err
	}
	str := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	//str := hex.EncodeToString(mac.Sum(nil))
	return str, nil
}
