package dto




type GetPvzInfoListResp struct {
	Items []GetPvzInfoResp
}

type GetPvzInfoResp struct {
	Pvz  Pvz
	Receptions []struct{
		Reception Reception
		Products []Product
	}
}