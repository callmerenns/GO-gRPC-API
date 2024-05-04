package models

type Product struct {
	Id          uint64
	Name        string
	Stock       uint32
	Price       float64
	Category_id uint32
}
