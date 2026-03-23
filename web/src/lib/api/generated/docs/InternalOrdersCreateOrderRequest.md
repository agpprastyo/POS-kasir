# InternalOrdersCreateOrderRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**customer_id** | **string** |  | [optional] [default to undefined]
**items** | [**Array&lt;InternalOrdersCreateOrderItemRequest&gt;**](InternalOrdersCreateOrderItemRequest.md) |  | [default to undefined]
**type** | [**POSKasirInternalOrdersRepositoryOrderType**](POSKasirInternalOrdersRepositoryOrderType.md) |  | [default to undefined]

## Example

```typescript
import { InternalOrdersCreateOrderRequest } from 'restClient';

const instance: InternalOrdersCreateOrderRequest = {
    customer_id,
    items,
    type,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
