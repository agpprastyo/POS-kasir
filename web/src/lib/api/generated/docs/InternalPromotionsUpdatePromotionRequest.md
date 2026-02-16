# InternalPromotionsUpdatePromotionRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**description** | **string** |  | [optional] [default to undefined]
**discount_type** | [**POSKasirInternalPromotionsRepositoryDiscountType**](POSKasirInternalPromotionsRepositoryDiscountType.md) |  | [default to undefined]
**discount_value** | **number** |  | [default to undefined]
**end_date** | **string** |  | [default to undefined]
**is_active** | **boolean** |  | [optional] [default to undefined]
**max_discount_amount** | **number** |  | [optional] [default to undefined]
**name** | **string** |  | [default to undefined]
**rules** | [**Array&lt;InternalPromotionsCreatePromotionRuleRequest&gt;**](InternalPromotionsCreatePromotionRuleRequest.md) |  | [optional] [default to undefined]
**scope** | [**POSKasirInternalPromotionsRepositoryPromotionScope**](POSKasirInternalPromotionsRepositoryPromotionScope.md) |  | [default to undefined]
**start_date** | **string** |  | [default to undefined]
**targets** | [**Array&lt;InternalPromotionsCreatePromotionTargetRequest&gt;**](InternalPromotionsCreatePromotionTargetRequest.md) |  | [optional] [default to undefined]

## Example

```typescript
import { InternalPromotionsUpdatePromotionRequest } from 'restClient';

const instance: InternalPromotionsUpdatePromotionRequest = {
    description,
    discount_type,
    discount_value,
    end_date,
    is_active,
    max_discount_amount,
    name,
    rules,
    scope,
    start_date,
    targets,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
