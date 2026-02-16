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
> PromotionsGet200Response promotionsGet()

Get a list of promotions with pagination and optional trash filter (Roles: admin, manager, cashier)

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

**PromotionsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Promotions retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **promotionsIdDelete**
> POSKasirInternalCommonSuccessResponse promotionsIdDelete()

Soft delete a promotion by its ID (Roles: admin, manager)

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
|**200** | Promotion deleted successfully |  -  |
|**400** | Invalid promotion ID format |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **promotionsIdGet**
> PromotionsPost201Response promotionsIdGet()

Retrieve details of a specific promotion by its ID (Roles: admin, manager, cashier)

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

**PromotionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Promotion retrieved successfully |  -  |
|**400** | Invalid promotion ID format |  -  |
|**404** | Promotion not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **promotionsIdPut**
> PromotionsPost201Response promotionsIdPut(request)

Update details of an existing promotion by its ID (Roles: admin, manager)

### Example

```typescript
import {
    PromotionsApi,
    Configuration,
    InternalPromotionsUpdatePromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let id: string; //Promotion ID (default to undefined)
let request: InternalPromotionsUpdatePromotionRequest; //Promotion details

const { status, data } = await apiInstance.promotionsIdPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalPromotionsUpdatePromotionRequest**| Promotion details | |
| **id** | [**string**] | Promotion ID | defaults to undefined|


### Return type

**PromotionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Promotion updated successfully |  -  |
|**400** | Invalid project ID format or request body |  -  |
|**404** | Promotion not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **promotionsIdRestorePost**
> POSKasirInternalCommonSuccessResponse promotionsIdRestorePost()

Restore a soft-deleted promotion by its ID (Roles: admin, manager)

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
|**200** | Promotion restored successfully |  -  |
|**400** | Invalid promotion ID format |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **promotionsPost**
> PromotionsPost201Response promotionsPost(request)

Create a new promotion with rules and targets (Roles: admin, manager)

### Example

```typescript
import {
    PromotionsApi,
    Configuration,
    InternalPromotionsCreatePromotionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PromotionsApi(configuration);

let request: InternalPromotionsCreatePromotionRequest; //Promotion details

const { status, data } = await apiInstance.promotionsPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalPromotionsCreatePromotionRequest**| Promotion details | |


### Return type

**PromotionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Promotion created successfully |  -  |
|**400** | Invalid request body or validation failed |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

