package factory

import "errors"

type ProductType int

const (
	A ProductType = iota
	B
)

//Product 产品接口
type Product interface {
	getPrice() int
}

//BaseProduct 基础产品
type BaseProduct struct{}

func (base BaseProduct) getPrice() int {
	return 10
}

//ProductA 产品A组合了基础产品
type ProductA struct {
	BaseProduct
}

type Creator struct{}

func (c Creator) CreateProduct(pt ProductType) (Product, error) {
	switch pt {
	case A:
		return &ProductA{BaseProduct: BaseProduct{}}, nil
	default:
		return nil, errors.New("产品不存在")
	}
}
