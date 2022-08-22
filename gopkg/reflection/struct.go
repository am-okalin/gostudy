package reflection

import (
	"github.com/hasura/go-graphql-client"
)

type SellingPlanPricingPolicyAdjustmentType graphql.String

type SellingPlanPricingPolicyInput struct { //最多两个定价政策
	Fixed struct {
		Id             graphql.ID
		AdjustmentType SellingPlanPricingPolicyAdjustmentType
	}
	AdjustmentType SellingPlanPricingPolicyAdjustmentType
	AfterCycle     graphql.Int
}

func mockSpi() interface{} {
	return SellingPlanPricingPolicyInput{
		Fixed: struct {
			Id             graphql.ID
			AdjustmentType SellingPlanPricingPolicyAdjustmentType
		}{
			Id:             1,
			AdjustmentType: "test",
		},
		AdjustmentType: "test",
		AfterCycle:     1,
	}
}

//func toMap() error {
//	m := make(map[string]interface{})
//	spi := mockSpi()
//
//	v := reflect.ValueOf(spi)
//
//	if v.Kind() != reflect.Struct {
//		return fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
//	}
//
//	t := v.Type()
//	for i:=0; i<v.NumField(); i++ {
//
//		fi := t.Field(i)
//
//		if tagValue := fi.Tag.Get(tag)
//	}
//	return nil
//}
