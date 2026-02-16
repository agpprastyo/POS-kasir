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

Record a manual cash entry or exit within the active shift (Roles: admin, manager, cashier)

### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    InternalShiftCashTransactionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: InternalShiftCashTransactionRequest; //Cash Transaction Request

const { status, data } = await apiInstance.shiftsCashTransactionPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalShiftCashTransactionRequest**| Cash Transaction Request | |


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
|**201** | Cash transaction created successfully |  -  |
|**400** | Invalid request body or validation failure |  -  |
|**404** | No open shift found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsCurrentGet**
> ShiftsCurrentGet200Response shiftsCurrentGet()

Check for and retrieve the details of an active shift session for the authenticated user (Roles: admin, manager, cashier)

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
|**200** | Open shift retrieved successfully |  -  |
|**404** | No open shift found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsEndPost**
> ShiftsCurrentGet200Response shiftsEndPost(request)

Close the active shift session for the authenticated user (Roles: admin, manager, cashier)

### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    InternalShiftEndShiftRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: InternalShiftEndShiftRequest; //End Shift Request

const { status, data } = await apiInstance.shiftsEndPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalShiftEndShiftRequest**| End Shift Request | |


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
|**200** | Shift ended successfully |  -  |
|**400** | Invalid request body or validation failure |  -  |
|**404** | No open shift found |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **shiftsStartPost**
> ShiftsCurrentGet200Response shiftsStartPost(request)

Create a new shift session for the authenticated user (Roles: admin, manager, cashier)

### Example

```typescript
import {
    ShiftsApi,
    Configuration,
    InternalShiftStartShiftRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ShiftsApi(configuration);

let request: InternalShiftStartShiftRequest; //Start Shift Request

const { status, data } = await apiInstance.shiftsStartPost(
    request
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **request** | **InternalShiftStartShiftRequest**| Start Shift Request | |


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
|**201** | Shift started successfully |  -  |
|**400** | Invalid request body or validation failure |  -  |
|**409** | User already has an open shift |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

