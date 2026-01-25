# POSKasirInternalDtoOrderListResponse


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**created_at** | **string** |  | [optional] [default to undefined]
**id** | **string** |  | [optional] [default to undefined]
**items** | [**Array&lt;POSKasirInternalDtoOrderItemResponse&gt;**](POSKasirInternalDtoOrderItemResponse.md) |  | [optional] [default to undefined]
**net_total** | **number** |  | [optional] [default to undefined]
**queue_number** | **string** |  | [optional] [default to undefined]
**status** | [**POSKasirInternalRepositoryOrderStatus**](POSKasirInternalRepositoryOrderStatus.md) |  | [optional] [default to undefined]
**type** | [**POSKasirInternalRepositoryOrderType**](POSKasirInternalRepositoryOrderType.md) |  | [optional] [default to undefined]
**user_id** | **string** |  | [optional] [default to undefined]

## Example

```typescript
import { POSKasirInternalDtoOrderListResponse } from 'restClient';

const instance: POSKasirInternalDtoOrderListResponse = {
    created_at,
    id,
    items,
    net_total,
    queue_number,
    status,
    type,
    user_id,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
