package Models

import "time"

type User struct {
	ID        int       `json:"id"`
	Nickname  string    `json:"nickname,omitempty"`
	UserName  string    `json:"user_name,omitempty"`
	Password  string    `json:"password,omitempty"`
	Sex       string    `json:"sex,omitempty"`
	Age       string    `json:"age,omitempty"`
	Address   string    `json:"address,omitempty"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
