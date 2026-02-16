# PrinterApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**ordersIdPrintPost**](#ordersidprintpost) | **POST** /orders/{id}/print | Print invoice for an order|
|[**settingsPrinterTestPost**](#settingsprintertestpost) | **POST** /settings/printer/test | Test printer connection|

# **ordersIdPrintPost**
> POSKasirInternalCommonSuccessResponse ordersIdPrintPost()

Trigger printing of invoice for a specific order (Roles: admin, manager, cashier)

### Example

```typescript
import {
    PrinterApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PrinterApi(configuration);

let id: string; //Order ID (default to undefined)

const { status, data } = await apiInstance.ordersIdPrintPost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Order ID | defaults to undefined|


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
|**200** | Invoice sent to printer |  -  |
|**400** | Invalid order ID |  -  |
|**500** | Failed to print invoice |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **settingsPrinterTestPost**
> POSKasirInternalCommonSuccessResponse settingsPrinterTestPost()

Send a test print command to the configured printer (Roles: admin)

### Example

```typescript
import {
    PrinterApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PrinterApi(configuration);

const { status, data } = await apiInstance.settingsPrinterTestPost();
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
|**200** | Test print command sent associated with configured printer |  -  |
|**500** | Failed to send test print |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

