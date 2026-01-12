# CancellationReasonsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**cancellationReasonsGet**](#cancellationreasonsget) | **GET** /cancellation-reasons | List cancellation reasons|

# **cancellationReasonsGet**
> CancellationReasonsGet200Response cancellationReasonsGet()


### Example

```typescript
import {
    CancellationReasonsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new CancellationReasonsApi(configuration);

const { status, data } = await apiInstance.cancellationReasonsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**CancellationReasonsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | List of cancellation reasons |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

