# AuthApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authLoginPost**](#authloginpost) | **POST** /auth/login | Login|
|[**authLogoutPost**](#authlogoutpost) | **POST** /auth/logout | Logout|
|[**authMeAvatarPut**](#authmeavatarput) | **PUT** /auth/me/avatar | Update avatar|
|[**authMeGet**](#authmeget) | **GET** /auth/me | Get profile|
|[**authMeUpdatePasswordPost**](#authmeupdatepasswordpost) | **POST** /auth/me/update-password | Update password|
|[**authRefreshPost**](#authrefreshpost) | **POST** /auth/refresh | Refresh token|

# **authLoginPost**
> AuthLoginPost200Response authLoginPost(request)

Login

### Example

```typescript
import {
    AuthApi,
    Configuration,
    POSKasirInternalDtoLoginRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: POSKasirInternalDtoLoginRequest; //Login request

const { status, data } = await apiInstance.authLoginPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoLoginRequest**| Login request | |


### Return type

**AuthLoginPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLogoutPost**
> AuthLogoutPost200Response authLogoutPost()

Logout

### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authLogoutPost();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**AuthLogoutPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMeAvatarPut**
> AuthMeGet200Response authMeAvatarPut()

Update avatar

### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let avatar: File; //Avatar file (default to undefined)

const { status, data } = await apiInstance.authMeAvatarPut(
    avatar
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **avatar** | [**File**] | Avatar file | defaults to undefined|


### Return type

**AuthMeGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMeGet**
> AuthMeGet200Response authMeGet()

Get profile

### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authMeGet();
```

### Parameters
This endpoint does not have any parameters.


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
|**200** | Success |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMeUpdatePasswordPost**
> AuthLogoutPost200Response authMeUpdatePasswordPost(request)

Update password

### Example

```typescript
import {
    AuthApi,
    Configuration,
    POSKasirInternalDtoUpdatePasswordRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: POSKasirInternalDtoUpdatePasswordRequest; //Update password request

const { status, data } = await apiInstance.authMeUpdatePasswordPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoUpdatePasswordRequest**| Update password request | |


### Return type

**AuthLogoutPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**400** | Bad Request |  -  |
|**401** | Unauthorized |  -  |
|**404** | Not Found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRefreshPost**
> AuthLoginPost200Response authRefreshPost()

Refresh token

### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.authRefreshPost();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**AuthLoginPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Success |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

