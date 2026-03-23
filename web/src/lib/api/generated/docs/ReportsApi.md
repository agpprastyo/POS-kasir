# ReportsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**reportsCancellationsGet**](#reportscancellationsget) | **GET** /reports/cancellations | Get cancellation reports|
|[**reportsCashierPerformanceGet**](#reportscashierperformanceget) | **GET** /reports/cashier-performance | Get cashier performance|
|[**reportsDashboardSummaryGet**](#reportsdashboardsummaryget) | **GET** /reports/dashboard-summary | Get dashboard summary|
|[**reportsLowStockGet**](#reportslowstockget) | **GET** /reports/low-stock | Get low stock products|
|[**reportsPaymentMethodsGet**](#reportspaymentmethodsget) | **GET** /reports/payment-methods | Get payment method performance|
|[**reportsProductsGet**](#reportsproductsget) | **GET** /reports/products | Get product performance|
|[**reportsProfitProductsGet**](#reportsprofitproductsget) | **GET** /reports/profit-products | Get product profit reports|
|[**reportsProfitSummaryGet**](#reportsprofitsummaryget) | **GET** /reports/profit-summary | Get profit summary|
|[**reportsPromotionsGet**](#reportspromotionsget) | **GET** /reports/promotions | Get promotion performance|
|[**reportsSalesGet**](#reportssalesget) | **GET** /reports/sales | Get sales reports|
|[**reportsShiftSummaryGet**](#reportsshiftsummaryget) | **GET** /reports/shift-summary | Get shift summary records|

# **reportsCancellationsGet**
> ReportsCancellationsGet200Response reportsCancellationsGet()

Get statistics on order cancellations grouped by reason (Roles: admin, manager, cashier)

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
|**200** | Cancellation reports retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsCashierPerformanceGet**
> ReportsCashierPerformanceGet200Response reportsCashierPerformanceGet()

Get order counts and sales totals handled by each cashier (Roles: admin, manager, cashier)

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
|**200** | Cashier performance data retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsDashboardSummaryGet**
> ReportsDashboardSummaryGet200Response reportsDashboardSummaryGet()

Get high-level summary metrics (totals) for the dashboard (Roles: admin, manager, cashier)

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

const { status, data } = await apiInstance.reportsDashboardSummaryGet(
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

**ReportsDashboardSummaryGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Dashboard summary retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsLowStockGet**
> ReportsLowStockGet200Response reportsLowStockGet()

Get products with stock below threshold

### Example

```typescript
import {
    ReportsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ReportsApi(configuration);

let threshold: number; //Threshold (default: 5) (optional) (default to undefined)
let _export: string; //Export format (csv) (optional) (default to undefined)

const { status, data } = await apiInstance.reportsLowStockGet(
    threshold,
    _export
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **threshold** | [**number**] | Threshold (default: 5) | (optional) defaults to undefined|
| **_export** | [**string**] | Export format (csv) | (optional) defaults to undefined|


### Return type

**ReportsLowStockGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Low stock products retrieved successfully |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsPaymentMethodsGet**
> ReportsPaymentMethodsGet200Response reportsPaymentMethodsGet()

Get usage counts and totals for each payment method (Roles: admin, manager, cashier)

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
|**200** | Payment method performance data retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsProductsGet**
> ReportsProductsGet200Response reportsProductsGet()

Get sales performance metrics for each product (Roles: admin, manager, cashier)

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
|**200** | Product performance data retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsProfitProductsGet**
> ReportsProfitProductsGet200Response reportsProfitProductsGet()

Get profitability metrics for each product sold (Roles: admin, manager, cashier)

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

const { status, data } = await apiInstance.reportsProfitProductsGet(
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

**ReportsProfitProductsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product profit reports retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsProfitSummaryGet**
> ReportsProfitSummaryGet200Response reportsProfitSummaryGet()

Get gross profit analytics grouped by date (Roles: admin, manager, cashier)

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

const { status, data } = await apiInstance.reportsProfitSummaryGet(
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

**ReportsProfitSummaryGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Profit summary retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsPromotionsGet**
> ReportsPromotionsGet200Response reportsPromotionsGet()

Get metrics of promotions usage

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
let _export: string; //Export format (csv) (optional) (default to undefined)

const { status, data } = await apiInstance.reportsPromotionsGet(
    startDate,
    endDate,
    _export
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|
| **_export** | [**string**] | Export format (csv) | (optional) defaults to undefined|


### Return type

**ReportsPromotionsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Promotion performance retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsSalesGet**
> ReportsSalesGet200Response reportsSalesGet()

Get aggregated sales data grouped by date within a specified range (Roles: admin, manager, cashier)

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
|**200** | Sales reports retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **reportsShiftSummaryGet**
> ReportsShiftSummaryGet200Response reportsShiftSummaryGet()

Get historical shifts and their cash differences

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
let _export: string; //Export format (csv) (optional) (default to undefined)

const { status, data } = await apiInstance.reportsShiftSummaryGet(
    startDate,
    endDate,
    _export
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **startDate** | [**string**] | Start Date (YYYY-MM-DD) | defaults to undefined|
| **endDate** | [**string**] | End Date (YYYY-MM-DD) | defaults to undefined|
| **_export** | [**string**] | Export format (csv) | (optional) defaults to undefined|


### Return type

**ReportsShiftSummaryGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Shift summary retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Internal server error |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

