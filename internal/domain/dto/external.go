package dto

type DFProductListReq struct {
	Cmd      string `json:"cmd"`      // "prepaid"
	Username string `json:"username"` // from cfg
	Sign     string `json:"sign"`     // md5(username+apiKey+"pricelist")
}

type DFBaseRes struct {
	Data []DFProductListRes `json:"data"`
}

type DFProductListRes struct {
	ProductName         string `json:"product_name"`
	Category            string `json:"category"`
	Brand               string `json:"brand"`
	Type                string `json:"type"`
	SellerName          string `json:"seller_name"`
	Price               int32  `json:"price"`
	BuyerSkuCode        string `json:"buyer_sku_code"`
	BuyerProductStatus  bool   `json:"buyer_product_status"`
	SellerProductStatus bool   `json:"seller_product_status"`
	UnlimitedStock      bool   `json:"unlimited_stock"`
	Stock               int32  `json:"stock"`
	Multi               bool   `json:"multi"`
	StartCutOff         string `json:"start_cut_off"`
	EndCutOff           string `json:"end_cut_off"`
	Desc                string `json:"desc"`
}
