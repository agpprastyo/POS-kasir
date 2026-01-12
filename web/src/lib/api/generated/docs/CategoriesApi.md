# CategoriesApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**categoriesCountGet**](#categoriescountget) | **GET** /categories/count | Get total number of categories|
|[**categoriesGet**](#categoriesget) | **GET** /categories | Get all categories|
|[**categoriesIdDelete**](#categoriesiddelete) | **DELETE** /categories/{id} | Delete category by ID|
|[**categoriesIdGet**](#categoriesidget) | **GET** /categories/{id} | Get category by ID|
|[**categoriesIdPut**](#categoriesidput) | **PUT** /categories/{id} | Update category by ID|
|[**categoriesPost**](#categoriespost) | **POST** /categories | Create a new category|

# **categoriesCountGet**
> CategoriesCountGet200Response categoriesCountGet()

Get total number of categories

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

const { status, data } = await apiInstance.categoriesCountGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CategoriesCountGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Category count retrieved successfully |  -  |
|**500** | Failed to retrieve category count |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **categoriesGet**
> CategoriesGet200Response categoriesGet()

Get all categories

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let limit: number; //Number of categories to return (optional) (default to undefined)
let offset: number; //Offset for pagination (optional) (default to undefined)

const { status, data } = await apiInstance.categoriesGet(
    limit,
    offset
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **limit** | [**number**] | Number of categories to return | (optional) defaults to undefined|
| **offset** | [**number**] | Offset for pagination | (optional) defaults to undefined|


### Return type

**CategoriesGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Categories retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Failed to retrieve categories |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **categoriesIdDelete**
> POSKasirInternalCommonSuccessResponse categoriesIdDelete()

Delete category by ID

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let id: number; //Category ID (default to undefined)

const { status, data } = await apiInstance.categoriesIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Category ID | defaults to undefined|


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
|**200** | Category deleted successfully |  -  |
|**400** | Invalid category ID format |  -  |
|**404** | Category not found |  -  |
|**409** | Category cannot be deleted because it is in use |  -  |
|**500** | Failed to delete category |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **categoriesIdGet**
> CategoriesPost201Response categoriesIdGet()

Get category by ID

### Example

```typescript
import {
    CategoriesApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let id: number; //Category ID (default to undefined)

const { status, data } = await apiInstance.categoriesIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**number**] | Category ID | defaults to undefined|


### Return type

**CategoriesPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Category retrieved successfully |  -  |
|**400** | Invalid category ID format |  -  |
|**404** | Category not found |  -  |
|**500** | Failed to retrieve category |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **categoriesIdPut**
> POSKasirInternalCommonSuccessResponse categoriesIdPut(category)

Update category by ID

### Example

```typescript
import {
    CategoriesApi,
    Configuration,
    POSKasirInternalDtoCreateCategoryRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let id: number; //Category ID (default to undefined)
let category: POSKasirInternalDtoCreateCategoryRequest; //Category details

const { status, data } = await apiInstance.categoriesIdPut(
    id,
    category
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **category** | **POSKasirInternalDtoCreateCategoryRequest**| Category details | |
| **id** | [**number**] | Category ID | defaults to undefined|


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
|**200** | Category deleted successfully |  -  |
|**400** | Invalid request body |  -  |
|**404** | Category not found |  -  |
|**409** | Category with this name already exists |  -  |
|**500** | Failed to update category |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **categoriesPost**
> CategoriesPost201Response categoriesPost(category)

Create a new category

### Example

```typescript
import {
    CategoriesApi,
    Configuration,
    POSKasirInternalDtoCreateCategoryRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CategoriesApi(configuration);

let category: POSKasirInternalDtoCreateCategoryRequest; //Category details

const { status, data } = await apiInstance.categoriesPost(
    category
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **category** | **POSKasirInternalDtoCreateCategoryRequest**| Category details | |


### Return type

**CategoriesPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Category created successfully |  -  |
|**400** | Invalid request body |  -  |
|**409** | Category with this name already exists |  -  |
|**500** | Failed to create category |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

