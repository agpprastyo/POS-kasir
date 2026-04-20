# InternalOrdersCheckoutOrderRequest


## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**cash_received** | **number** |  | [optional] [default to undefined]
**customer_id** | **string** |  | [optional] [default to undefined]
**items** | [**Array&lt;InternalOrdersCreateOrderItemRequest&gt;**](InternalOrdersCreateOrderItemRequest.md) |  | [default to undefined]
**payment_method_id** | **number** |  | [optional] [default to undefined]
**promotion_id** | **string** |  | [optional] [default to undefined]
**type** | [**POSKasirInternalOrdersRepositoryOrderType**](POSKasirInternalOrdersRepositoryOrderType.md) |  | [default to undefined]

## Example

```typescript
import { InternalOrdersCheckoutOrderRequest } from 'restClient';

const instance: InternalOrdersCheckoutOrderRequest = {
    cash_received,
    customer_id,
    items,
    payment_method_id,
    promotion_id,
    type,
};
```

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)
