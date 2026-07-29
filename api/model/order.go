package model

// Order は注文1件分の情報を表す。
type Order struct {
	ID        string `json:"id"`
	Customer  string `json:"customer"`
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
	OrderedAt string `json:"orderedAt"`
}
