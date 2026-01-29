# ReportsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**reportsCancellationsGet**](#reportscancellationsget) | **GET** /reports/cancellations | Get cancellation reports|
|[**reportsCashierPerformanceGet**](#reportscashierperformanceget) | **GET** /reports/cashier-performance | Get cashier performance|
|[**reportsDashboardSummaryGet**](#reportsdashboardsummaryget) | **GET** /reports/dashboard-summary | Get dashboard summary|
|[**reportsPaymentMethodsGet**](#reportspaymentmethodsget) | **GET** /reports/payment-methods | Get payment method performance|
|[**reportsProductsGet**](#reportsproductsget) | **GET** /reports/products | Get product performance|
|[**reportsSalesGet**](#reportssalesget) | **GET** /reports/sales | Get sales reports|

# **reportsCancellationsGet**
> ReportsCancellationsGet200Response reportsCancellationsGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let startDate: string; //Start Date (YYYY-MM-DD) (default to undefined)
let endDate: string; //End Date (YYYY-MM-DD) (default to undefined)

const { status, data } = await apiInstance.reportsCancellationsGet(
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|


### Return type

**ReportsCancellationsGet200Response**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsCashierPerformanceGet**
> ReportsCashierPerformanceGet200Response reportsCashierPerformanceGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let startDate: string; //Start Date (YYYY-MM-DD) (default to undefined)
let endDate: string; //End Date (YYYY-MM-DD) (default to undefined)

const { status, data } = await apiInstance.reportsCashierPerformanceGet(
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|


### Return type

**ReportsCashierPerformanceGet200Response**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsDashboardSummaryGet**
> ReportsDashboardSummaryGet200Response reportsDashboardSummaryGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

const { status, data } = await apiInstance.reportsDashboardSummaryGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ReportsDashboardSummaryGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsPaymentMethodsGet**
> ReportsPaymentMethodsGet200Response reportsPaymentMethodsGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let startDate: string; //Start Date (YYYY-MM-DD) (default to undefined)
let endDate: string; //End Date (YYYY-MM-DD) (default to undefined)

const { status, data } = await apiInstance.reportsPaymentMethodsGet(
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|


### Return type

**ReportsPaymentMethodsGet200Response**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsProductsGet**
> ReportsProductsGet200Response reportsProductsGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let startDate: string; //Start Date (YYYY-MM-DD) (default to undefined)
let endDate: string; //End Date (YYYY-MM-DD) (default to undefined)

const { status, data } = await apiInstance.reportsProductsGet(
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|


### Return type

**ReportsProductsGet200Response**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsSalesGet**
> ReportsSalesGet200Response reportsSalesGet()


### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let startDate: string; //Start Date (YYYY-MM-DD) (default to undefined)
let endDate: string; //End Date (YYYY-MM-DD) (default to undefined)

const { status, data } = await apiInstance.reportsSalesGet(
    startDate,
    endDate
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|


### Return type

**ReportsSalesGet200Response**

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
|**500** | Internal Server Error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

