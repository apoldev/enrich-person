package person

import "errors"

var (
	ErrNotFound = errors.New("not found")
)

// Person описывает сущность и DTO
type Person struct {
	ID         int     `json:"id,omitempty" db:"id" `
	Name       string  `json:"name,omitempty" db:"name" binding:"required"`
	Surname    string  `json:"surname,omitempty" db:"surname" binding:"required"`
	Patronymic *string `json:"patronymic,omitempty" db:"patronymic"`

	Nationality string `json:"nationality,omitempty" db:"nationality"`
	Age         int    `json:"age,omitempty" db:"age"`
	Gender      string `json:"gender,omitempty" db:"gender"`
}
