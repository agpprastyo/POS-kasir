# POSKasirInternalDtoOrderDetailResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**applied_promotion_id** | **string** |  | [optional] [default to undefined]
**created_at** | **string** |  | [optional] [default to undefined]
**discount_amount** | **number** |  | [optional] [default to undefined]
**gross_total** | **number** |  | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**items** | [**Array&lt;POSKasirInternalDtoOrderItemResponse&gt;**](POSKasirInternalDtoOrderItemResponse.md) |  | [optional] [default to undefined]
**net_total** | **number** |  | [optional] [default to undefined]
**payment_gateway_reference** | **string** |  | [optional] [default to undefined]
**payment_method_id** | **number** |  | [optional] [default to undefined]
**status** | [**POSKasirInternalRepositoryOrderStatus**](POSKasirInternalRepositoryOrderStatus.md) |  | [optional] [default to undefined]
**type** | [**POSKasirInternalRepositoryOrderType**](POSKasirInternalRepositoryOrderType.md) |  | [optional] [default to undefined]
**updated_at** | **string** |  | [optional] [default to undefined]
**user_id** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { POSKasirInternalDtoOrderDetailResponse } from 'restClient';

const instance: POSKasirInternalDtoOrderDetailResponse = {
    applied_promotion_id,
    created_at,
    discount_amount,
    gross_total,
    id,
    items,
    net_total,
    payment_gateway_reference,
    payment_method_id,
    status,
    type,
    updated_at,
    user_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
