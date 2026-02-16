# InternalOrdersOrderListResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** |  | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**is_paid** | **boolean** |  | [optional] [default to undefined]
**items** | [**Array&lt;InternalOrdersOrderItemResponse&gt;**](InternalOrdersOrderItemResponse.md) |  | [optional] [default to undefined]
**net_total** | **number** |  | [optional] [default to undefined]
**queue_number** | **string** |  | [optional] [default to undefined]
**status** | [**POSKasirInternalOrdersRepositoryOrderStatus**](POSKasirInternalOrdersRepositoryOrderStatus.md) |  | [optional] [default to undefined]
**type** | [**POSKasirInternalOrdersRepositoryOrderType**](POSKasirInternalOrdersRepositoryOrderType.md) |  | [optional] [default to undefined]
**user_id** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { InternalOrdersOrderListResponse } from 'restClient';

const instance: InternalOrdersOrderListResponse = {
    created_at,
    id,
    is_paid,
    items,
    net_total,
    queue_number,
    status,
    type,
    user_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
