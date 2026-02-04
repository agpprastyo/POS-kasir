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

# **ordersGet**
> OrdersGet200Response ordersGet()


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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdApplyPromotionPost**
> OrdersPost201Response ordersIdApplyPromotionPost(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirInternalDtoApplyPromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: POSKasirInternalDtoApplyPromotionRequest; //Promotion details

const { status, data } = await apiInstance.ordersIdApplyPromotionPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoApplyPromotionRequest**| Promotion details | |
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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdCancelPost**
> POSKasirInternalCommonSuccessResponse ordersIdCancelPost(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirInternalDtoCancelOrderRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: POSKasirInternalDtoCancelOrderRequest; //Cancel order details

const { status, data } = await apiInstance.ordersIdCancelPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoCancelOrderRequest**| Cancel order details | |
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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdGet**
> OrdersPost201Response ordersIdGet()


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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdItemsPut**
> OrdersPost201Response ordersIdItemsPut(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: Array<POSKasirInternalDtoUpdateOrderItemRequest>; //Update order items

const { status, data } = await apiInstance.ordersIdItemsPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **Array<POSKasirInternalDtoUpdateOrderItemRequest>**| Update order items | |
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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdPayManualPost**
> OrdersPost201Response ordersIdPayManualPost(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirInternalDtoConfirmManualPaymentRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: POSKasirInternalDtoConfirmManualPaymentRequest; //Manual payment details

const { status, data } = await apiInstance.ordersIdPayManualPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoConfirmManualPaymentRequest**| Manual payment details | |
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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdPayMidtransPost**
> OrdersIdPayMidtransPost200Response ordersIdPayMidtransPost()


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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersIdUpdateStatusPost**
> OrdersPost201Response ordersIdUpdateStatusPost(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirInternalDtoUpdateOrderStatusRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let id: string; //Order ID (default to undefined)
let request: POSKasirInternalDtoUpdateOrderStatusRequest; //Order status details

const { status, data } = await apiInstance.ordersIdUpdateStatusPost(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoUpdateOrderStatusRequest**| Order status details | |
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
|**200** | OK |  -  |
|**400** | Bad Request |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ordersPost**
> OrdersPost201Response ordersPost(request)


### Example

```typescript
import {
    OrdersApi,
    Configuration,
    POSKasirInternalDtoCreateOrderRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new OrdersApi(configuration);

let request: POSKasirInternalDtoCreateOrderRequest; //Create order details

const { status, data } = await apiInstance.ordersPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoCreateOrderRequest**| Create order details | |


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
|**201** | Created |  -  |
|**400** | Bad Request |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

