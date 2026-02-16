# AuthApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**authAddPost**](#authaddpost) | **POST** /auth/add | Add new user|
|[**authLoginPost**](#authloginpost) | **POST** /auth/login | Login|
|[**authLogoutPost**](#authlogoutpost) | **POST** /auth/logout | Logout|
|[**authMeAvatarPut**](#authmeavatarput) | **PUT** /auth/me/avatar | Update avatar|
|[**authMeGet**](#authmeget) | **GET** /auth/me | Get current profile|
|[**authMePasswordPut**](#authmepasswordput) | **PUT** /auth/me/password | Update password|
|[**authRefreshPost**](#authrefreshpost) | **POST** /auth/refresh | Refresh token|

# **authAddPost**
> AuthAddPost200Response authAddPost(request)

Register a new user with a specific role (Roles: admin)

### Example

```typescript
import {
    AuthApi,
    Configuration,
    InternalUserRegisterRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: InternalUserRegisterRequest; //New user details

const { status, data } = await apiInstance.authAddPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalUserRegisterRequest**| New user details | |


### Return type

**AuthAddPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | User added successfully |  -  |
|**400** | Invalid request body or validation failed |  -  |
|**403** | Forbidden - higher role assignment attempt |  -  |
|**409** | User, username, or email already exists |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLoginPost**
> AuthLoginPost200Response authLoginPost(request)

Authenticate user and return access/refresh tokens via cookies and response body (Roles: public)

### Example

```typescript
import {
    AuthApi,
    Configuration,
    InternalUserLoginRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: InternalUserLoginRequest; //Login credentials

const { status, data } = await apiInstance.authLoginPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalUserLoginRequest**| Login credentials | |


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
|**200** | Authenticated successfully |  -  |
|**400** | Invalid request body or validation failed |  -  |
|**401** | Invalid username or password |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authLogoutPost**
> POSKasirInternalCommonSuccessResponse authLogoutPost()

Clear access and refresh token cookies (Roles: authenticated)

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

**POSKasirInternalCommonSuccessResponse**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Successfully logged out |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMeAvatarPut**
> AuthAddPost200Response authMeAvatarPut()

Upload and update the profile picture for the current user (Roles: authenticated)

### Example

```typescript
import {
    AuthApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let avatar: File; //Avatar image file (default to undefined)

const { status, data } = await apiInstance.authMeAvatarPut(
    avatar
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **avatar** | [**File**] | Avatar image file | defaults to undefined|


### Return type

**AuthAddPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Avatar updated successfully |  -  |
|**400** | Invalid file or dimensions |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMeGet**
> AuthAddPost200Response authMeGet()

Get detailed profile information for the authenticated user session (Roles: authenticated)

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

**AuthAddPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Profile retrieved successfully |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authMePasswordPut**
> POSKasirInternalCommonSuccessResponse authMePasswordPut(request)

Update the password for the current user session (Roles: authenticated)

### Example

```typescript
import {
    AuthApi,
    Configuration,
    InternalUserUpdatePasswordRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let request: InternalUserUpdatePasswordRequest; //Password update details

const { status, data } = await apiInstance.authMePasswordPut(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalUserUpdatePasswordRequest**| Password update details | |


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
|**200** | Password updated successfully |  -  |
|**400** | Invalid request body or validation failed |  -  |
|**401** | Unauthorized |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **authRefreshPost**
> AuthRefreshPost200Response authRefreshPost()

Issue a new access token using a valid refresh token cookie (Roles: public/authenticated)

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

**AuthRefreshPost200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Token refreshed successfully |  -  |
|**401** | Invalid or expired session |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

