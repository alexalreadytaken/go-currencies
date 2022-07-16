package rest

type BtcUsdtHistoryPage struct {
	Total   int                  `json:"total"`
	History []BtcUsdtCourseSlice `json:"history"`
}

type FiatToAnyHistoryPage struct {
	Total   int                    `json:"total"`
	History []FiatToAnyCourseSlice `json:"history"`
}

type FiatPaginationRequest struct {
	Limit    uint     `json:"limit"`
	Offset   uint     `json:"offset"`
	FromDate DateOnly `json:"from_date"`
}

//todo cut limit and offset
type BtcUsdtPaginationRequest struct {
	Limit    uint   `json:"limit"`
	Offset   uint   `json:"offset"`
	FromTime uint64 `json:"from_timestamp"`
}
