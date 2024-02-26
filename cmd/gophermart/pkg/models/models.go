package models

const (
	New        = "NEW"
	Processing = "PROCESSING"
	Invalid    = "INVALID"
	Processed  = "PROCESSED"
	Registered = "REGISTERED"
)

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Order struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float32 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}

type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float32 `json:"accrual,omitempty"`
}

type UserBalance struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type WithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Withdraw struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
