# PaymentMethodsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**paymentMethodsGet**](#paymentmethodsget) | **GET** /payment-methods | List payment methods|

# **paymentMethodsGet**
> PaymentMethodsGet200Response paymentMethodsGet()

Get a list of all active payment methods (e.g., Cash, QRIS)

### Example

```typescript
import {
    PaymentMethodsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new PaymentMethodsApi(configuration);

const { status, data } = await apiInstance.paymentMethodsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**PaymentMethodsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of payment methods retrieved successfully |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

