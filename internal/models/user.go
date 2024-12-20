package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// Вместо обычный типов go, были использованны типы pgtype
// потому что почти все поля таблици - nullable
type User struct {
	ID          uint64      `json:"id"`
	Age         pgtype.Int8 `json:"age"`
	Username    pgtype.Text `json:"username"`
	Email       pgtype.Text `json:"email"`
	City        pgtype.Text `json:"city"`
	DateBirth   pgtype.Date `json:"date_birth"`
	Phone       pgtype.Text `json:"phone"`
	Description pgtype.Text `json:"description"`
}
