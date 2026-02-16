# OrdersApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**ordersGet**](#ordersget) | **GET** /orders | List orders|
|[**ordersIdApplyPromotionPost**](#ordersidapplypromotionpost) | **POST** /orders/{id}/apply-promotion | Apply promotion to an order|
|[**ordersIdCancelPost**](#ordersidcancelpost) | **POST** /orders/{id}/cancel | Cancel an order|
|[**ordersIdGet**](#ordersidget) | **GET** /orders/{id} | Get an order by ID|
|[**ordersIdItemsPut**](#ordersiditemsput) | **PUT** /orders/{id}/items | Update items in an order|
|[**ordersIdPayManualPost**](#ordersidpaymanualpost) | **POST** /orders/{id}/pay/manual | Confirm manual payment for an order|
|[**ordersIdPayMidtransPost**](#ordersidpaymidtranspost) | **POST** /orders/{id}/pay/midtrans | Initiate Midtrans payment for an order|
|[**ordersIdUpdateStatusPost**](#ordersidupdatestatuspost) | **POST** /orders/{id}/update-status | Update order operational status|
|[**ordersPost**](#orderspost) | **POST** /orders | Create an order|
|[**ordersWebhookMidtransPost**](#orderswebhookmidtranspost) | **POST** /orders/webhook/midtrans | Midtrans Payment Notification Callback|

# **ordersGet**
> OrdersGet200Response ordersGet()

Get a list of orders with filtering by status and user (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Number of orders per page (optional) (default to undefined)
let statuses: Array<'open' | 'in_progress' | 'served' | 'paid' | 'cancelled'>; //Order statuses (optional) (default to undefined)
let userId: string; //Filter by User ID (optional) (default to undefined)

const { status, data } = await apiInstance.ordersGet(
    page,
    limit,
    statuses,
    userId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Number of orders per page | (optional) defaults to undefined|
| **statuses** | **Array<&#39;open&#39; &#124; &#39;in_progress&#39; &#124; &#39;served&#39; &#124; &#39;paid&#39; &#124; &#39;cancelled&#39;>** | Order statuses | (optional) defaults to undefined|
| **userId** | [**string**] | Filter by User ID | (optional) defaults to undefined|


### Return type

**OrdersGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Orders retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Failed to retrieve orders |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdApplyPromotionPost**
> OrdersPost201Response ordersIdApplyPromotionPost(request)

Apply a specific promotion to an existing order by its ID (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    InternalOrdersApplyPromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: InternalOrdersApplyPromotionRequest; //Promotion details

const { status, data } = await apiInstance.ordersIdApplyPromotionPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalOrdersApplyPromotionRequest**| Promotion details | |
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Promotion applied successfully |  -  |
|**400** | Invalid order ID format or request body |  -  |
|**404** | Order or Promotion not found |  -  |
|**500** | Failed to apply promotion |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdCancelPost**
> POSKasirInternalCommonSuccessResponse ordersIdCancelPost(request)

Cancel an existing order with a reason (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    InternalOrdersCancelOrderRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: InternalOrdersCancelOrderRequest; //Cancel order details

const { status, data } = await apiInstance.ordersIdCancelPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalOrdersCancelOrderRequest**| Cancel order details | |
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**POSKasirInternalCommonSuccessResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order cancelled successfully |  -  |
|**400** | Invalid order ID format or request body |  -  |
|**404** | Order not found |  -  |
|**409** | Order cannot be cancelled |  -  |
|**500** | Failed to cancel order |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdGet**
> OrdersPost201Response ordersIdGet()

Retrieve detailed information of a specific order by its ID (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)

const { status, data } = await apiInstance.ordersIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order retrieved successfully |  -  |
|**400** | Invalid order ID format |  -  |
|**404** | Order not found |  -  |
|**500** | Failed to retrieve order |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdItemsPut**
> OrdersPost201Response ordersIdItemsPut(request)

Update, add, or remove items in an existing open order (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: Array<InternalOrdersUpdateOrderItemRequest>; //Update order items

const { status, data } = await apiInstance.ordersIdItemsPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **Array<InternalOrdersUpdateOrderItemRequest>**| Update order items | |
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order items updated successfully |  -  |
|**400** | Invalid order ID format or request body |  -  |
|**404** | Order not found |  -  |
|**500** | Failed to update order items |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdPayManualPost**
> OrdersPost201Response ordersIdPayManualPost(request)

Process a manual payment (Cash) and finalize an order (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    InternalOrdersConfirmManualPaymentRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: InternalOrdersConfirmManualPaymentRequest; //Manual payment details

const { status, data } = await apiInstance.ordersIdPayManualPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalOrdersConfirmManualPaymentRequest**| Manual payment details | |
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Payment completed successfully |  -  |
|**400** | Invalid order ID format or request body |  -  |
|**404** | Order not found |  -  |
|**409** | Order might have been paid or cancelled |  -  |
|**500** | Failed to complete payment |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdPayMidtransPost**
> OrdersIdPayMidtransPost200Response ordersIdPayMidtransPost()

Create a QRIS/Gopay payment session via Midtrans for an existing order (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)

const { status, data } = await apiInstance.ordersIdPayMidtransPost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersIdPayMidtransPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | QRIS payment initiated successfully |  -  |
|**400** | Invalid order ID format |  -  |
|**404** | Order not found |  -  |
|**500** | Failed to process payment |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdUpdateStatusPost**
> OrdersPost201Response ordersIdUpdateStatusPost(request)

Update the status of an existing order (e.g., to in_progress, served) (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    InternalOrdersUpdateOrderStatusRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: InternalOrdersUpdateOrderStatusRequest; //Order status details

const { status, data } = await apiInstance.ordersIdUpdateStatusPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalOrdersUpdateOrderStatusRequest**| Order status details | |
| **id** | [**string**] | Order ID | defaults to undefined|


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Order status updated successfully |  -  |
|**400** | Invalid order ID format or request body |  -  |
|**404** | Order not found |  -  |
|**409** | Invalid status transition |  -  |
|**500** | Failed to update order status |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersPost**
> OrdersPost201Response ordersPost(request)

Create a new order with multiple items (Roles: admin, manager, cashier)

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    InternalOrdersCreateOrderRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let request: InternalOrdersCreateOrderRequest; //Create order details

const { status, data } = await apiInstance.ordersPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalOrdersCreateOrderRequest**| Create order details | |


### Return type

**OrdersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Order created successfully |  -  |
|**400** | Invalid request body |  -  |
|**500** | Failed to create order |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersWebhookMidtransPost**
> POSKasirInternalCommonSuccessResponse ordersWebhookMidtransPost(payload)

Webhook for Midtrans to notify order payment status updates

### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirPkgPaymentMidtransNotificationPayload
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let payload: POSKasirPkgPaymentMidtransNotificationPayload; //Midtrans Notification Payload

const { status, data } = await apiInstance.ordersWebhookMidtransPost(
    payload
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **payload** | **POSKasirPkgPaymentMidtransNotificationPayload**| Midtrans Notification Payload | |


### Return type

**POSKasirInternalCommonSuccessResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Notification received successfully |  -  |
|**400** | Invalid notification format |  -  |
|**500** | Failed to handle notification |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

