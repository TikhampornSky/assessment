package repos

import "database/sql"

type Err struct {
	Message string `json:"message"`
}

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"` //transfer Name --> name (ตอนรับข้อมูล)
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type handler struct {
	DB *sql.DB
}