# PromotionsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**promotionsGet**](#promotionsget) | **GET** /promotions | List all promotions|
|[**promotionsIdDelete**](#promotionsiddelete) | **DELETE** /promotions/{id} | Delete (deactivate) a promotion|
|[**promotionsIdGet**](#promotionsidget) | **GET** /promotions/{id} | Get a promotion by ID|
|[**promotionsIdPut**](#promotionsidput) | **PUT** /promotions/{id} | Update a promotion|
|[**promotionsIdRestorePost**](#promotionsidrestorepost) | **POST** /promotions/{id}/restore | Restore a deleted promotion|
|[**promotionsPost**](#promotionspost) | **POST** /promotions | Create a new promotion|

# **promotionsGet**
> POSKasirInternalCommonSuccessResponse promotionsGet()


### Example

```typescript
import {
    PromotionsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Items per page (optional) (default to undefined)
let trash: boolean; //Show trash items (optional) (default to undefined)

const { status, data } = await apiInstance.promotionsGet(
    page,
    limit,
    trash
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Items per page | (optional) defaults to undefined|
| **trash** | [**boolean**] | Show trash items | (optional) defaults to undefined|


### Return type

**POSKasirInternalCommonSuccessResponse**

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

# **promotionsIdDelete**
> POSKasirInternalCommonSuccessResponse promotionsIdDelete()


### Example

```typescript
import {
    PromotionsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let id: string; //Promotion ID (default to undefined)

const { status, data } = await apiInstance.promotionsIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Promotion ID | defaults to undefined|


### Return type

**POSKasirInternalCommonSuccessResponse**

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

# **promotionsIdGet**
> POSKasirInternalCommonSuccessResponse promotionsIdGet()


### Example

```typescript
import {
    PromotionsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let id: string; //Promotion ID (default to undefined)

const { status, data } = await apiInstance.promotionsIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Promotion ID | defaults to undefined|


### Return type

**POSKasirInternalCommonSuccessResponse**

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

# **promotionsIdPut**
> POSKasirInternalCommonSuccessResponse promotionsIdPut(request)


### Example

```typescript
import {
    PromotionsApi,
    Configuration,
    POSKasirInternalDtoUpdatePromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let id: string; //Promotion ID (default to undefined)
let request: POSKasirInternalDtoUpdatePromotionRequest; //Promotion details

const { status, data } = await apiInstance.promotionsIdPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoUpdatePromotionRequest**| Promotion details | |
| **id** | [**string**] | Promotion ID | defaults to undefined|


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

# **promotionsIdRestorePost**
> POSKasirInternalCommonSuccessResponse promotionsIdRestorePost()


### Example

```typescript
import {
    PromotionsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let id: string; //Promotion ID (default to undefined)

const { status, data } = await apiInstance.promotionsIdRestorePost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Promotion ID | defaults to undefined|


### Return type

**POSKasirInternalCommonSuccessResponse**

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

# **promotionsPost**
> POSKasirInternalCommonSuccessResponse promotionsPost(request)


### Example

```typescript
import {
    PromotionsApi,
    Configuration,
    POSKasirInternalDtoCreatePromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let request: POSKasirInternalDtoCreatePromotionRequest; //Promotion details

const { status, data } = await apiInstance.promotionsPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoCreatePromotionRequest**| Promotion details | |


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
|**201** | Created |  -  |
|**400** | Bad Request |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

