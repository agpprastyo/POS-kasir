# ProductsApi

All URIs are relative to *http://localhost:8080/api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**productsGet**](#productsget) | **GET** /products | List products|
|[**productsIdDelete**](#productsiddelete) | **DELETE** /products/{id} | Delete a product|
|[**productsIdGet**](#productsidget) | **GET** /products/{id} | Get a product by ID|
|[**productsIdImagePost**](#productsidimagepost) | **POST** /products/{id}/image | Upload an image for a product|
|[**productsIdPatch**](#productsidpatch) | **PATCH** /products/{id} | Update a product|
|[**productsIdStockHistoryGet**](#productsidstockhistoryget) | **GET** /products/{id}/stock-history | Get stock history for a product|
|[**productsPost**](#productspost) | **POST** /products | Create a new product|
|[**productsProductIdOptionsOptionIdImagePost**](#productsproductidoptionsoptionidimagepost) | **POST** /products/{product_id}/options/{option_id}/image | Upload product option image|
|[**productsProductIdOptionsOptionIdPatch**](#productsproductidoptionsoptionidpatch) | **PATCH** /products/{product_id}/options/{option_id} | Update a product option|
|[**productsProductIdOptionsPost**](#productsproductidoptionspost) | **POST** /products/{product_id}/options | Create a product option|
|[**productsTrashGet**](#productstrashget) | **GET** /products/trash | List deleted products|
|[**productsTrashIdGet**](#productstrashidget) | **GET** /products/trash/{id} | Get a deleted product|
|[**productsTrashIdRestorePost**](#productstrashidrestorepost) | **POST** /products/trash/{id}/restore | Restore a deleted product|
|[**productsTrashRestoreBulkPost**](#productstrashrestorebulkpost) | **POST** /products/trash/restore-bulk | Bulk restore deleted products|

# **productsGet**
> ProductsGet200Response productsGet()

List products based on query parameters

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Limit the number of products returned (optional) (default to undefined)
let search: string; //Search products by name (optional) (default to undefined)
let categoryId: number; //Search products by category ID (optional) (default to undefined)

const { status, data } = await apiInstance.productsGet(
    page,
    limit,
    search,
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Limit the number of products returned | (optional) defaults to undefined|
| **search** | [**string**] | Search products by name | (optional) defaults to undefined|
| **categoryId** | [**number**] | Search products by category ID | (optional) defaults to undefined|


### Return type

**ProductsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Products retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Failed to retrieve products |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsIdDelete**
> POSKasirInternalCommonSuccessResponse productsIdDelete()

Delete a product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)

const { status, data } = await apiInstance.productsIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|


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
|**200** | Product deleted successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to delete product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsIdGet**
> ProductsPost201Response productsIdGet()

Get a product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)

const { status, data } = await apiInstance.productsIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|


### Return type

**ProductsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product retrieved successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to retrieve product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsIdImagePost**
> ProductsPost201Response productsIdImagePost()

Upload an image for a product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)
let image: File; //Image file (default to undefined)

const { status, data } = await apiInstance.productsIdImagePost(
    id,
    image
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|
| **image** | [**File**] | Image file | defaults to undefined|


### Return type

**ProductsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product image uploaded successfully |  -  |
|**400** | Invalid product ID format or image file is missing |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to upload image |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsIdPatch**
> ProductsPost201Response productsIdPatch(body)

Update a product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration,
    POSKasirInternalDtoUpdateProductRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)
let body: POSKasirInternalDtoUpdateProductRequest; //Product update request

const { status, data } = await apiInstance.productsIdPatch(
    id,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **POSKasirInternalDtoUpdateProductRequest**| Product update request | |
| **id** | [**string**] | Product ID | defaults to undefined|


### Return type

**ProductsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product updated successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to update product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsIdStockHistoryGet**
> ProductsIdStockHistoryGet200Response productsIdStockHistoryGet()

Get stock history for a product by ID with pagination

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)
let page: number; //Page number (optional) (default to undefined)
let limit: number; //Limit (optional) (default to undefined)

const { status, data } = await apiInstance.productsIdStockHistoryGet(
    id,
    page,
    limit
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Limit | (optional) defaults to undefined|


### Return type

**ProductsIdStockHistoryGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Stock history retrieved successfully |  -  |
|**400** | Invalid product ID or query parameters |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to retrieve stock history |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsPost**
> ProductsPost201Response productsPost(body)

Create a new product

### Example

```typescript
import {
    ProductsApi,
    Configuration,
    POSKasirInternalDtoCreateProductRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let body: POSKasirInternalDtoCreateProductRequest; //Product create request

const { status, data } = await apiInstance.productsPost(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **POSKasirInternalDtoCreateProductRequest**| Product create request | |


### Return type

**ProductsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Product created successfully |  -  |
|**400** | Invalid request body |  -  |
|**409** | Product with same name already exists |  -  |
|**500** | Failed to create product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsProductIdOptionsOptionIdImagePost**
> ProductsProductIdOptionsPost201Response productsProductIdOptionsOptionIdImagePost()

Upload image for a specific product option

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let productId: string; //Product ID (default to undefined)
let optionId: string; //Option ID (default to undefined)
let image: File; //Product image (default to undefined)

const { status, data } = await apiInstance.productsProductIdOptionsOptionIdImagePost(
    productId,
    optionId,
    image
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **productId** | [**string**] | Product ID | defaults to undefined|
| **optionId** | [**string**] | Option ID | defaults to undefined|
| **image** | [**File**] | Product image | defaults to undefined|


### Return type

**ProductsProductIdOptionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product option image uploaded successfully |  -  |
|**400** | Invalid ID format or missing file |  -  |
|**404** | Product or option not found |  -  |
|**500** | Failed to upload product image |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsProductIdOptionsOptionIdPatch**
> ProductsProductIdOptionsPost201Response productsProductIdOptionsOptionIdPatch(body)

Update a product option by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration,
    POSKasirInternalDtoUpdateProductOptionRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let productId: string; //Product ID (default to undefined)
let optionId: string; //Option ID (default to undefined)
let body: POSKasirInternalDtoUpdateProductOptionRequest; //Product option update request

const { status, data } = await apiInstance.productsProductIdOptionsOptionIdPatch(
    productId,
    optionId,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **POSKasirInternalDtoUpdateProductOptionRequest**| Product option update request | |
| **productId** | [**string**] | Product ID | defaults to undefined|
| **optionId** | [**string**] | Option ID | defaults to undefined|


### Return type

**ProductsProductIdOptionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Product option updated successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product or option not found |  -  |
|**500** | Failed to update product option |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsProductIdOptionsPost**
> ProductsProductIdOptionsPost201Response productsProductIdOptionsPost(body)

Create a product option for a parent product

### Example

```typescript
import {
    ProductsApi,
    Configuration,
    POSKasirInternalDtoCreateProductOptionRequestStandalone
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let productId: string; //Product ID (default to undefined)
let body: POSKasirInternalDtoCreateProductOptionRequestStandalone; //Product option create request

const { status, data } = await apiInstance.productsProductIdOptionsPost(
    productId,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **POSKasirInternalDtoCreateProductOptionRequestStandalone**| Product option create request | |
| **productId** | [**string**] | Product ID | defaults to undefined|


### Return type

**ProductsProductIdOptionsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**201** | Product option created successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Parent product not found |  -  |
|**500** | Failed to create product option |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsTrashGet**
> ProductsGet200Response productsTrashGet()

List deleted products with pagination and filtering

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let page: number; //Page number (optional) (default to undefined)
let limit: number; //Limit the number of products returned (optional) (default to undefined)
let search: string; //Search products by name (optional) (default to undefined)
let categoryId: number; //Search products by category ID (optional) (default to undefined)

const { status, data } = await apiInstance.productsTrashGet(
    page,
    limit,
    search,
    categoryId
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **page** | [**number**] | Page number | (optional) defaults to undefined|
| **limit** | [**number**] | Limit the number of products returned | (optional) defaults to undefined|
| **search** | [**string**] | Search products by name | (optional) defaults to undefined|
| **categoryId** | [**number**] | Search products by category ID | (optional) defaults to undefined|


### Return type

**ProductsGet200Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Deleted products retrieved successfully |  -  |
|**400** | Invalid query parameters |  -  |
|**500** | Failed to retrieve deleted products |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsTrashIdGet**
> ProductsPost201Response productsTrashIdGet()

Get a deleted product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)

const { status, data } = await apiInstance.productsTrashIdGet(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|


### Return type

**ProductsPost201Response**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | Deleted product retrieved successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to retrieve deleted product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsTrashIdRestorePost**
> POSKasirInternalCommonSuccessResponse productsTrashIdRestorePost()

Restore a deleted product by ID

### Example

```typescript
import {
    ProductsApi,
    Configuration
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let id: string; //Product ID (default to undefined)

const { status, data } = await apiInstance.productsTrashIdRestorePost(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | Product ID | defaults to undefined|


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
|**200** | Product restored successfully |  -  |
|**400** | Invalid product ID format |  -  |
|**404** | Product not found |  -  |
|**500** | Failed to restore product |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **productsTrashRestoreBulkPost**
> POSKasirInternalCommonSuccessResponse productsTrashRestoreBulkPost(body)

Restore multiple deleted products by IDs

### Example

```typescript
import {
    ProductsApi,
    Configuration,
    POSKasirInternalDtoRestoreBulkRequest
} from 'restClient';

const configuration = new Configuration();
const apiInstance = new ProductsApi(configuration);

let body: POSKasirInternalDtoRestoreBulkRequest; //Bulk restore request

const { status, data } = await apiInstance.productsTrashRestoreBulkPost(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **POSKasirInternalDtoRestoreBulkRequest**| Bulk restore request | |


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
|**200** | Products restored successfully |  -  |
|**400** | Invalid request body |  -  |
|**500** | Failed to restore products |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

