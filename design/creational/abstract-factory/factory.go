package abstract_factory

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
type BaseProduct struct {
	style string
}

func (base BaseProduct) getPrice() int {
	return 10
}

//ProductA 产品A组合了基础产品
type ProductA struct {
	BaseProduct
}

//Creator 工厂接口
type Creator interface {
	CreateProduct(pt ProductType) (Product, error)
}

type CreatorA struct{}

func (ca CreatorA) CreateProduct(pt ProductType) (Product, error) {
	style := "A style"
	switch pt {
	case A:
		return &ProductA{BaseProduct: BaseProduct{style: style}}, nil
	default:
		return nil, errors.New("产品不存在")
	}
}
