# UsersApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**usersGet**](#usersget) | **GET** /users | Get all users|
|[**usersIdDelete**](#usersiddelete) | **DELETE** /users/{id} | Delete user|
|[**usersIdGet**](#usersidget) | **GET** /users/{id} | Get user by ID|
|[**usersIdPut**](#usersidput) | **PUT** /users/{id} | Update user|
|[**usersIdToggleStatusPost**](#usersidtogglestatuspost) | **POST** /users/{id}/toggle-status | Toggle user status|
|[**usersPost**](#userspost) | **POST** /users | Create user|

# **usersGet**
> UsersGet200Response usersGet()

Retrieve a list of users with pagination, filtering, and sorting

### Example

```typescript
import {
    UsersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let page: number; //Page number (default 1) (optional) (default to 1)
let limit: number; //Items per page (default 10) (optional) (default to 10)
let search: string; //Search by username or email (optional) (default to undefined)
let role: 'admin' | 'cashier' | 'manager'; //Filter by User Role (optional) (default to undefined)
let isActive: boolean; //Filter by Active Status (optional) (default to undefined)
let status: 'active' | 'deleted' | 'all'; //Filter by Account Status (optional) (default to undefined)
let sortBy: 'created_at' | 'username'; //Sort by column (optional) (default to undefined)
let sortOrder: 'asc' | 'desc'; //Sort direction (optional) (default to undefined)

const { status, data } = await apiInstance.usersGet(
    page,
    limit,
    search,
    role,
    isActive,
    status,
    sortBy,
    sortOrder
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number (default 1) | (optional) defaults to 1|
| **limit** | [**number**] | Items per page (default 10) | (optional) defaults to 10|
| **search** | [**string**] | Search by username or email | (optional) defaults to undefined|
| **role** | [**&#39;admin&#39; | &#39;cashier&#39; | &#39;manager&#39;**]**Array<&#39;admin&#39; &#124; &#39;cashier&#39; &#124; &#39;manager&#39;>** | Filter by User Role | (optional) defaults to undefined|
| **isActive** | [**boolean**] | Filter by Active Status | (optional) defaults to undefined|
| **status** | [**&#39;active&#39; | &#39;deleted&#39; | &#39;all&#39;**]**Array<&#39;active&#39; &#124; &#39;deleted&#39; &#124; &#39;all&#39;>** | Filter by Account Status | (optional) defaults to undefined|
| **sortBy** | [**&#39;created_at&#39; | &#39;username&#39;**]**Array<&#39;created_at&#39; &#124; &#39;username&#39;>** | Sort by column | (optional) defaults to undefined|
| **sortOrder** | [**&#39;asc&#39; | &#39;desc&#39;**]**Array<&#39;asc&#39; &#124; &#39;desc&#39;>** | Sort direction | (optional) defaults to undefined|


### Return type

**UsersGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Invalid query parameters |  -  |
|**404** | No users found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **usersIdDelete**
> POSKasirInternalCommonSuccessResponse usersIdDelete()


### Example

```typescript
import {
    UsersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let id: string; //User ID (default to undefined)

const { status, data } = await apiInstance.usersIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | User ID | defaults to undefined|


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
|**404** | User not found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **usersIdGet**
> AuthMeGet200Response usersIdGet()


### Example

```typescript
import {
    UsersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let id: string; //User ID (default to undefined)

const { status, data } = await apiInstance.usersIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | User ID | defaults to undefined|


### Return type

**AuthMeGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | User ID is required |  -  |
|**404** | User not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **usersIdPut**
> AuthMeGet200Response usersIdPut(user)


### Example

```typescript
import {
    UsersApi,
    Configuration,
    POSKasirInternalDtoUpdateUserRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let id: string; //User ID (default to undefined)
let user: POSKasirInternalDtoUpdateUserRequest; //User details

const { status, data } = await apiInstance.usersIdPut(
    id,
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **POSKasirInternalDtoUpdateUserRequest**| User details | |
| **id** | [**string**] | User ID | defaults to undefined|


### Return type

**AuthMeGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**400** | Invalid request body |  -  |
|**401** | Unauthorized |  -  |
|**403** | You are not authorized to change user roles |  -  |
|**404** | User not found |  -  |
|**409** | Username already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **usersIdToggleStatusPost**
> POSKasirInternalCommonSuccessResponse usersIdToggleStatusPost()


### Example

```typescript
import {
    UsersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let id: string; //User ID (default to undefined)

const { status, data } = await apiInstance.usersIdToggleStatusPost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | User ID | defaults to undefined|


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
|**400** | User ID is required |  -  |
|**404** | User not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **usersPost**
> AuthMeGet200Response usersPost(user)


### Example

```typescript
import {
    UsersApi,
    Configuration,
    POSKasirInternalDtoCreateUserRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new UsersApi(configuration);

let user: POSKasirInternalDtoCreateUserRequest; //User details

const { status, data } = await apiInstance.usersPost(
    user
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **user** | **POSKasirInternalDtoCreateUserRequest**| User details | |


### Return type

**AuthMeGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Created |  -  |
|**400** | Invalid request body |  -  |
|**409** | User already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

