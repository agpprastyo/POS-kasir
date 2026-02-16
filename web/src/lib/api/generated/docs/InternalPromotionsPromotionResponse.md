# InternalPromotionsPromotionResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** |  | [optional] [default to undefined]
**deleted_at** | **string** |  | [optional] [default to undefined]
**description** | **string** |  | [optional] [default to undefined]
**discount_type** | [**POSKasirInternalPromotionsRepositoryDiscountType**](POSKasirInternalPromotionsRepositoryDiscountType.md) |  | [optional] [default to undefined]
**discount_value** | **number** |  | [optional] [default to undefined]
**end_date** | **string** |  | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**is_active** | **boolean** |  | [optional] [default to undefined]
**max_discount_amount** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [optional] [default to undefined]
**rules** | [**Array&lt;InternalPromotionsPromotionRuleResponse&gt;**](InternalPromotionsPromotionRuleResponse.md) |  | [optional] [default to undefined]
**scope** | [**POSKasirInternalPromotionsRepositoryPromotionScope**](POSKasirInternalPromotionsRepositoryPromotionScope.md) |  | [optional] [default to undefined]
**start_date** | **string** |  | [optional] [default to undefined]
**targets** | [**Array&lt;InternalPromotionsPromotionTargetResponse&gt;**](InternalPromotionsPromotionTargetResponse.md) |  | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { InternalPromotionsPromotionResponse } from 'restClient';

const instance: InternalPromotionsPromotionResponse = {
    created_at,
    deleted_at,
    description,
    discount_type,
    discount_value,
    end_date,
    id,
    is_active,
    max_discount_amount,
    name,
    rules,
    scope,
    start_date,
    targets,
    updated_at,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
