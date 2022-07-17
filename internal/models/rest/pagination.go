package rest

type BtcUsdtHistoryPage struct {
	Total   int                  `json:"total"`
	History []BtcUsdtCourseSlice `json:"history"`
}

type AnyToFiatHistoryPage struct {
	Total   int                    `json:"total"`
	History []AnyToFiatCourseSlice `json:"history"`
}

type AnyToFiatPaginationRequest struct {
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
