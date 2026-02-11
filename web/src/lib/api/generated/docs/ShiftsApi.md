# ShiftsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**shiftsCashTransactionPost**](#shiftscashtransactionpost) | **POST** /shifts/cash-transaction | Create a cash transaction (Drop/Expense/In)|
|[**shiftsCurrentGet**](#shiftscurrentget) | **GET** /shifts/current | Get current open shift|
|[**shiftsEndPost**](#shiftsendpost) | **POST** /shifts/end | End current shift|
|[**shiftsStartPost**](#shiftsstartpost) | **POST** /shifts/start | Start a new shift|

# **shiftsCashTransactionPost**
> ShiftsCashTransactionPost201Response shiftsCashTransactionPost(request)


### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    POSKasirInternalDtoCashTransactionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: POSKasirInternalDtoCashTransactionRequest; //Cash Transaction Request

const { status, data } = await apiInstance.shiftsCashTransactionPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoCashTransactionRequest**| Cash Transaction Request | |


### Return type

**ShiftsCashTransactionPost201Response**

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
|**404** | No open shift found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsCurrentGet**
> ShiftsCurrentGet200Response shiftsCurrentGet()


### Example

```typescript
import {
    ShiftsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

const { status, data } = await apiInstance.shiftsCurrentGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ShiftsCurrentGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**404** | No open shift found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsEndPost**
> ShiftsCurrentGet200Response shiftsEndPost(request)


### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    POSKasirInternalDtoEndShiftRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: POSKasirInternalDtoEndShiftRequest; //End Shift Request

const { status, data } = await apiInstance.shiftsEndPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoEndShiftRequest**| End Shift Request | |


### Return type

**ShiftsCurrentGet200Response**

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
|**404** | No open shift found |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsStartPost**
> ShiftsCurrentGet200Response shiftsStartPost(request)


### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    POSKasirInternalDtoStartShiftRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: POSKasirInternalDtoStartShiftRequest; //Start Shift Request

const { status, data } = await apiInstance.shiftsStartPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **POSKasirInternalDtoStartShiftRequest**| Start Shift Request | |


### Return type

**ShiftsCurrentGet200Response**

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
|**409** | User already has an open shift |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

