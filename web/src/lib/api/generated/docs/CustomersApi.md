# CustomersApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**customersGet**](#customersget) | **GET** /customers | List customers|
|[**customersIdDelete**](#customersiddelete) | **DELETE** /customers/{id} | Delete a customer|
|[**customersIdGet**](#customersidget) | **GET** /customers/{id} | Get a customer|
|[**customersIdPut**](#customersidput) | **PUT** /customers/{id} | Update a customer|
|[**customersPost**](#customerspost) | **POST** /customers | Create a customer|

# **customersGet**
> CustomersGet200Response customersGet()

List customers with pagination and search

### Example

```typescript
import {
    CustomersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CustomersApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Page size limit (optional) (default to undefined)
let search: string; //Search by name, phone, or email (optional) (default to undefined)

const { status, data } = await apiInstance.customersGet(
    page,
    limit,
    search
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Page size limit | (optional) defaults to undefined|
| **search** | [**string**] | Search by name, phone, or email | (optional) defaults to undefined|


### Return type

**CustomersGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Customers retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **customersIdDelete**
> POSKasirInternalCommonSuccessResponse customersIdDelete()

Delete customer by ID

### Example

```typescript
import {
    CustomersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CustomersApi(configuration);

let id: string; //Customer ID (default to undefined)

const { status, data } = await apiInstance.customersIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Customer ID | defaults to undefined|


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
|**200** | Customer deleted successfully |  -  |
|**400** | Invalid ID format |  -  |
|**404** | Customer not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **customersIdGet**
> CustomersPost201Response customersIdGet()

Get customer by ID

### Example

```typescript
import {
    CustomersApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CustomersApi(configuration);

let id: string; //Customer ID (default to undefined)

const { status, data } = await apiInstance.customersIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Customer ID | defaults to undefined|


### Return type

**CustomersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Customer retrieved successfully |  -  |
|**400** | Invalid ID format |  -  |
|**404** | Customer not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **customersIdPut**
> CustomersPost201Response customersIdPut(request)

Update customer by ID

### Example

```typescript
import {
    CustomersApi,
    Configuration,
    InternalCustomersUpdateCustomerRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CustomersApi(configuration);

let id: string; //Customer ID (default to undefined)
let request: InternalCustomersUpdateCustomerRequest; //Customer details to update

const { status, data } = await apiInstance.customersIdPut(
    id,
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalCustomersUpdateCustomerRequest**| Customer details to update | |
| **id** | [**string**] | Customer ID | defaults to undefined|


### Return type

**CustomersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Customer updated successfully |  -  |
|**400** | Invalid request |  -  |
|**404** | Customer not found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **customersPost**
> CustomersPost201Response customersPost(request)

Create a new customer (Roles: admin, manager, cashier)

### Example

```typescript
import {
    CustomersApi,
    Configuration,
    InternalCustomersCreateCustomerRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CustomersApi(configuration);

let request: InternalCustomersCreateCustomerRequest; //Customer details

const { status, data } = await apiInstance.customersPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalCustomersCreateCustomerRequest**| Customer details | |


### Return type

**CustomersPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Customer created successfully |  -  |
|**400** | Invalid request |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

