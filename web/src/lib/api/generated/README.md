## restClient@1.0

This generator creates TypeScript/JavaScript client that utilizes [axios](https://github.com/axios/axios). The generated Node module can be used in the following environments:

Environment
* Node.js
* Webpack
* Browserify

Language level
* ES5 - you must have a Promises/A+ library installed
* ES6

Module system
* CommonJS
* ES6 module system

It can be used in both TypeScript and JavaScript. In TypeScript, the definition will be automatically resolved via `package.json`. ([Reference](https://www.typescriptlang.org/docs/handbook/declaration-files/consumption.html))

### Building

To build and compile the typescript sources to javascript use:
```
npm install
npm run build
```

### Publishing

First build the package then run `npm publish`

### Consuming

navigate to the folder of your consuming project and run one of the following commands.

_published:_

```
npm install restClient@1.0 --save
```

_unPublished (not recommended):_

```
npm install PATH_TO_GENERATED_PACKAGE --save
```

### Documentation for API Endpoints

All URIs are relative to *http://localhost:8080/api/v1*

Class | Method | HTTP request | Description
------------ | ------------- | ------------- | -------------
*ActivityLogsApi* | [**activityLogsGet**](docs/ActivityLogsApi.md#activitylogsget) | **GET** /activity-logs | Get activity logs
*AuthApi* | [**authAddPost**](docs/AuthApi.md#authaddpost) | **POST** /auth/add | Add new user
*AuthApi* | [**authLoginPost**](docs/AuthApi.md#authloginpost) | **POST** /auth/login | Login
*AuthApi* | [**authLogoutPost**](docs/AuthApi.md#authlogoutpost) | **POST** /auth/logout | Logout
*AuthApi* | [**authMeAvatarPut**](docs/AuthApi.md#authmeavatarput) | **PUT** /auth/me/avatar | Update avatar
*AuthApi* | [**authMeGet**](docs/AuthApi.md#authmeget) | **GET** /auth/me | Get current profile
*AuthApi* | [**authMePasswordPut**](docs/AuthApi.md#authmepasswordput) | **PUT** /auth/me/password | Update password
*AuthApi* | [**authRefreshPost**](docs/AuthApi.md#authrefreshpost) | **POST** /auth/refresh | Refresh token
*CancellationReasonsApi* | [**cancellationReasonsGet**](docs/CancellationReasonsApi.md#cancellationreasonsget) | **GET** /cancellation-reasons | List cancellation reasons
*CategoriesApi* | [**categoriesCountGet**](docs/CategoriesApi.md#categoriescountget) | **GET** /categories/count | Get total number of categories
*CategoriesApi* | [**categoriesGet**](docs/CategoriesApi.md#categoriesget) | **GET** /categories | Get all categories
*CategoriesApi* | [**categoriesIdDelete**](docs/CategoriesApi.md#categoriesiddelete) | **DELETE** /categories/{id} | Delete category by ID
*CategoriesApi* | [**categoriesIdGet**](docs/CategoriesApi.md#categoriesidget) | **GET** /categories/{id} | Get category by ID
*CategoriesApi* | [**categoriesIdPut**](docs/CategoriesApi.md#categoriesidput) | **PUT** /categories/{id} | Update category by ID
*CategoriesApi* | [**categoriesPost**](docs/CategoriesApi.md#categoriespost) | **POST** /categories | Create a new category
*OrdersApi* | [**ordersGet**](docs/OrdersApi.md#ordersget) | **GET** /orders | List orders
*OrdersApi* | [**ordersIdApplyPromotionPost**](docs/OrdersApi.md#ordersidapplypromotionpost) | **POST** /orders/{id}/apply-promotion | Apply promotion to an order
*OrdersApi* | [**ordersIdCancelPost**](docs/OrdersApi.md#ordersidcancelpost) | **POST** /orders/{id}/cancel | Cancel an order
*OrdersApi* | [**ordersIdGet**](docs/OrdersApi.md#ordersidget) | **GET** /orders/{id} | Get an order by ID
*OrdersApi* | [**ordersIdItemsPut**](docs/OrdersApi.md#ordersiditemsput) | **PUT** /orders/{id}/items | Update items in an order
*OrdersApi* | [**ordersIdPayManualPost**](docs/OrdersApi.md#ordersidpaymanualpost) | **POST** /orders/{id}/pay/manual | Confirm manual payment for an order
*OrdersApi* | [**ordersIdPayMidtransPost**](docs/OrdersApi.md#ordersidpaymidtranspost) | **POST** /orders/{id}/pay/midtrans | Initiate Midtrans payment for an order
*OrdersApi* | [**ordersIdUpdateStatusPost**](docs/OrdersApi.md#ordersidupdatestatuspost) | **POST** /orders/{id}/update-status | Update order operational status
*OrdersApi* | [**ordersPost**](docs/OrdersApi.md#orderspost) | **POST** /orders | Create an order
*OrdersApi* | [**ordersWebhookMidtransPost**](docs/OrdersApi.md#orderswebhookmidtranspost) | **POST** /orders/webhook/midtrans | Midtrans Payment Notification Callback
*PaymentMethodsApi* | [**paymentMethodsGet**](docs/PaymentMethodsApi.md#paymentmethodsget) | **GET** /payment-methods | List payment methods
*PrinterApi* | [**ordersIdPrintDataGet**](docs/PrinterApi.md#ordersidprintdataget) | **GET** /orders/{id}/print-data | Get invoice print data
*PrinterApi* | [**ordersIdPrintPost**](docs/PrinterApi.md#ordersidprintpost) | **POST** /orders/{id}/print | Print invoice for an order
*PrinterApi* | [**settingsPrinterTestPost**](docs/PrinterApi.md#settingsprintertestpost) | **POST** /settings/printer/test | Test printer connection
*ProductsApi* | [**productsGet**](docs/ProductsApi.md#productsget) | **GET** /products | List products
*ProductsApi* | [**productsIdDelete**](docs/ProductsApi.md#productsiddelete) | **DELETE** /products/{id} | Delete a product
*ProductsApi* | [**productsIdGet**](docs/ProductsApi.md#productsidget) | **GET** /products/{id} | Get a product by ID
*ProductsApi* | [**productsIdImagePost**](docs/ProductsApi.md#productsidimagepost) | **POST** /products/{id}/image | Upload an image for a product
*ProductsApi* | [**productsIdPatch**](docs/ProductsApi.md#productsidpatch) | **PATCH** /products/{id} | Update a product
*ProductsApi* | [**productsIdStockHistoryGet**](docs/ProductsApi.md#productsidstockhistoryget) | **GET** /products/{id}/stock-history | Get stock history for a product
*ProductsApi* | [**productsPost**](docs/ProductsApi.md#productspost) | **POST** /products | Create a new product
*ProductsApi* | [**productsProductIdOptionsOptionIdImagePost**](docs/ProductsApi.md#productsproductidoptionsoptionidimagepost) | **POST** /products/{product_id}/options/{option_id}/image | Upload product option image
*ProductsApi* | [**productsProductIdOptionsOptionIdPatch**](docs/ProductsApi.md#productsproductidoptionsoptionidpatch) | **PATCH** /products/{product_id}/options/{option_id} | Update a product option
*ProductsApi* | [**productsProductIdOptionsPost**](docs/ProductsApi.md#productsproductidoptionspost) | **POST** /products/{product_id}/options | Create a product option
*ProductsApi* | [**productsTrashGet**](docs/ProductsApi.md#productstrashget) | **GET** /products/trash | List deleted products
*ProductsApi* | [**productsTrashIdGet**](docs/ProductsApi.md#productstrashidget) | **GET** /products/trash/{id} | Get a deleted product
*ProductsApi* | [**productsTrashIdRestorePost**](docs/ProductsApi.md#productstrashidrestorepost) | **POST** /products/trash/{id}/restore | Restore a deleted product
*ProductsApi* | [**productsTrashRestoreBulkPost**](docs/ProductsApi.md#productstrashrestorebulkpost) | **POST** /products/trash/restore-bulk | Bulk restore deleted products
*PromotionsApi* | [**promotionsGet**](docs/PromotionsApi.md#promotionsget) | **GET** /promotions | List all promotions
*PromotionsApi* | [**promotionsIdDelete**](docs/PromotionsApi.md#promotionsiddelete) | **DELETE** /promotions/{id} | Delete (deactivate) a promotion
*PromotionsApi* | [**promotionsIdGet**](docs/PromotionsApi.md#promotionsidget) | **GET** /promotions/{id} | Get a promotion by ID
*PromotionsApi* | [**promotionsIdPut**](docs/PromotionsApi.md#promotionsidput) | **PUT** /promotions/{id} | Update a promotion
*PromotionsApi* | [**promotionsIdRestorePost**](docs/PromotionsApi.md#promotionsidrestorepost) | **POST** /promotions/{id}/restore | Restore a deleted promotion
*PromotionsApi* | [**promotionsPost**](docs/PromotionsApi.md#promotionspost) | **POST** /promotions | Create a new promotion
*ReportsApi* | [**reportsCancellationsGet**](docs/ReportsApi.md#reportscancellationsget) | **GET** /reports/cancellations | Get cancellation reports
*ReportsApi* | [**reportsCashierPerformanceGet**](docs/ReportsApi.md#reportscashierperformanceget) | **GET** /reports/cashier-performance | Get cashier performance
*ReportsApi* | [**reportsDashboardSummaryGet**](docs/ReportsApi.md#reportsdashboardsummaryget) | **GET** /reports/dashboard-summary | Get dashboard summary
*ReportsApi* | [**reportsPaymentMethodsGet**](docs/ReportsApi.md#reportspaymentmethodsget) | **GET** /reports/payment-methods | Get payment method performance
*ReportsApi* | [**reportsProductsGet**](docs/ReportsApi.md#reportsproductsget) | **GET** /reports/products | Get product performance
*ReportsApi* | [**reportsProfitProductsGet**](docs/ReportsApi.md#reportsprofitproductsget) | **GET** /reports/profit-products | Get product profit reports
*ReportsApi* | [**reportsProfitSummaryGet**](docs/ReportsApi.md#reportsprofitsummaryget) | **GET** /reports/profit-summary | Get profit summary
*ReportsApi* | [**reportsSalesGet**](docs/ReportsApi.md#reportssalesget) | **GET** /reports/sales | Get sales reports
*SettingsApi* | [**settingsBrandingGet**](docs/SettingsApi.md#settingsbrandingget) | **GET** /settings/branding | Get branding settings
*SettingsApi* | [**settingsBrandingLogoPost**](docs/SettingsApi.md#settingsbrandinglogopost) | **POST** /settings/branding/logo | Update app logo
*SettingsApi* | [**settingsBrandingPut**](docs/SettingsApi.md#settingsbrandingput) | **PUT** /settings/branding | Update branding settings
*SettingsApi* | [**settingsPrinterGet**](docs/SettingsApi.md#settingsprinterget) | **GET** /settings/printer | Get printer settings
*SettingsApi* | [**settingsPrinterPut**](docs/SettingsApi.md#settingsprinterput) | **PUT** /settings/printer | Update printer settings
*ShiftsApi* | [**shiftsCashTransactionPost**](docs/ShiftsApi.md#shiftscashtransactionpost) | **POST** /shifts/cash-transaction | Create a cash transaction (Drop/Expense/In)
*ShiftsApi* | [**shiftsCurrentGet**](docs/ShiftsApi.md#shiftscurrentget) | **GET** /shifts/current | Get current open shift
*ShiftsApi* | [**shiftsEndPost**](docs/ShiftsApi.md#shiftsendpost) | **POST** /shifts/end | End current shift
*ShiftsApi* | [**shiftsStartPost**](docs/ShiftsApi.md#shiftsstartpost) | **POST** /shifts/start | Start a new shift
*UsersApi* | [**usersGet**](docs/UsersApi.md#usersget) | **GET** /users | Get all users
*UsersApi* | [**usersIdDelete**](docs/UsersApi.md#usersiddelete) | **DELETE** /users/{id} | Delete user
*UsersApi* | [**usersIdGet**](docs/UsersApi.md#usersidget) | **GET** /users/{id} | Get user by ID
*UsersApi* | [**usersIdPut**](docs/UsersApi.md#usersidput) | **PUT** /users/{id} | Update user
*UsersApi* | [**usersIdToggleStatusPost**](docs/UsersApi.md#usersidtogglestatuspost) | **POST** /users/{id}/toggle-status | Toggle user status
*UsersApi* | [**usersPost**](docs/UsersApi.md#userspost) | **POST** /users | Create user


### Documentation For Models

 - [ActivityLogsGet200Response](docs/ActivityLogsGet200Response.md)
 - [AuthAddPost200Response](docs/AuthAddPost200Response.md)
 - [AuthLoginPost200Response](docs/AuthLoginPost200Response.md)
 - [AuthRefreshPost200Response](docs/AuthRefreshPost200Response.md)
 - [CancellationReasonsGet200Response](docs/CancellationReasonsGet200Response.md)
 - [CategoriesGet200Response](docs/CategoriesGet200Response.md)
 - [CategoriesPost201Response](docs/CategoriesPost201Response.md)
 - [InternalActivitylogActivityLogListResponse](docs/InternalActivitylogActivityLogListResponse.md)
 - [InternalActivitylogActivityLogResponse](docs/InternalActivitylogActivityLogResponse.md)
 - [InternalCancellationReasonsCancellationReasonResponse](docs/InternalCancellationReasonsCancellationReasonResponse.md)
 - [InternalCategoriesCategoryResponse](docs/InternalCategoriesCategoryResponse.md)
 - [InternalCategoriesCategoryWithCountResponse](docs/InternalCategoriesCategoryWithCountResponse.md)
 - [InternalCategoriesCreateCategoryRequest](docs/InternalCategoriesCreateCategoryRequest.md)
 - [InternalOrdersApplyPromotionRequest](docs/InternalOrdersApplyPromotionRequest.md)
 - [InternalOrdersCancelOrderRequest](docs/InternalOrdersCancelOrderRequest.md)
 - [InternalOrdersConfirmManualPaymentRequest](docs/InternalOrdersConfirmManualPaymentRequest.md)
 - [InternalOrdersCreateOrderItemOptionRequest](docs/InternalOrdersCreateOrderItemOptionRequest.md)
 - [InternalOrdersCreateOrderItemRequest](docs/InternalOrdersCreateOrderItemRequest.md)
 - [InternalOrdersCreateOrderRequest](docs/InternalOrdersCreateOrderRequest.md)
 - [InternalOrdersMidtransPaymentResponse](docs/InternalOrdersMidtransPaymentResponse.md)
 - [InternalOrdersOrderDetailResponse](docs/InternalOrdersOrderDetailResponse.md)
 - [InternalOrdersOrderItemOptionResponse](docs/InternalOrdersOrderItemOptionResponse.md)
 - [InternalOrdersOrderItemResponse](docs/InternalOrdersOrderItemResponse.md)
 - [InternalOrdersOrderListResponse](docs/InternalOrdersOrderListResponse.md)
 - [InternalOrdersPagedOrderResponse](docs/InternalOrdersPagedOrderResponse.md)
 - [InternalOrdersPaymentAction](docs/InternalOrdersPaymentAction.md)
 - [InternalOrdersUpdateOrderItemRequest](docs/InternalOrdersUpdateOrderItemRequest.md)
 - [InternalOrdersUpdateOrderStatusRequest](docs/InternalOrdersUpdateOrderStatusRequest.md)
 - [InternalPaymentMethodsPaymentMethodResponse](docs/InternalPaymentMethodsPaymentMethodResponse.md)
 - [InternalProductsCreateProductOptionRequest](docs/InternalProductsCreateProductOptionRequest.md)
 - [InternalProductsCreateProductOptionRequestStandalone](docs/InternalProductsCreateProductOptionRequestStandalone.md)
 - [InternalProductsCreateProductRequest](docs/InternalProductsCreateProductRequest.md)
 - [InternalProductsListProductsResponse](docs/InternalProductsListProductsResponse.md)
 - [InternalProductsPagedStockHistoryResponse](docs/InternalProductsPagedStockHistoryResponse.md)
 - [InternalProductsProductListResponse](docs/InternalProductsProductListResponse.md)
 - [InternalProductsProductOptionResponse](docs/InternalProductsProductOptionResponse.md)
 - [InternalProductsProductResponse](docs/InternalProductsProductResponse.md)
 - [InternalProductsRestoreBulkRequest](docs/InternalProductsRestoreBulkRequest.md)
 - [InternalProductsStockHistoryResponse](docs/InternalProductsStockHistoryResponse.md)
 - [InternalProductsUpdateProductOptionRequest](docs/InternalProductsUpdateProductOptionRequest.md)
 - [InternalProductsUpdateProductRequest](docs/InternalProductsUpdateProductRequest.md)
 - [InternalPromotionsCreatePromotionRequest](docs/InternalPromotionsCreatePromotionRequest.md)
 - [InternalPromotionsCreatePromotionRuleRequest](docs/InternalPromotionsCreatePromotionRuleRequest.md)
 - [InternalPromotionsCreatePromotionTargetRequest](docs/InternalPromotionsCreatePromotionTargetRequest.md)
 - [InternalPromotionsPagedPromotionResponse](docs/InternalPromotionsPagedPromotionResponse.md)
 - [InternalPromotionsPromotionResponse](docs/InternalPromotionsPromotionResponse.md)
 - [InternalPromotionsPromotionRuleResponse](docs/InternalPromotionsPromotionRuleResponse.md)
 - [InternalPromotionsPromotionTargetResponse](docs/InternalPromotionsPromotionTargetResponse.md)
 - [InternalPromotionsUpdatePromotionRequest](docs/InternalPromotionsUpdatePromotionRequest.md)
 - [InternalReportCancellationReportResponse](docs/InternalReportCancellationReportResponse.md)
 - [InternalReportCashierPerformanceResponse](docs/InternalReportCashierPerformanceResponse.md)
 - [InternalReportDashboardSummaryResponse](docs/InternalReportDashboardSummaryResponse.md)
 - [InternalReportPaymentMethodPerformanceResponse](docs/InternalReportPaymentMethodPerformanceResponse.md)
 - [InternalReportProductPerformanceResponse](docs/InternalReportProductPerformanceResponse.md)
 - [InternalReportProductProfitResponse](docs/InternalReportProductProfitResponse.md)
 - [InternalReportProfitSummaryResponse](docs/InternalReportProfitSummaryResponse.md)
 - [InternalReportSalesReport](docs/InternalReportSalesReport.md)
 - [InternalSettingsBrandingSettingsResponse](docs/InternalSettingsBrandingSettingsResponse.md)
 - [InternalSettingsPrinterSettingsResponse](docs/InternalSettingsPrinterSettingsResponse.md)
 - [InternalSettingsUpdateBrandingRequest](docs/InternalSettingsUpdateBrandingRequest.md)
 - [InternalSettingsUpdatePrinterSettingsRequest](docs/InternalSettingsUpdatePrinterSettingsRequest.md)
 - [InternalShiftCashTransactionRequest](docs/InternalShiftCashTransactionRequest.md)
 - [InternalShiftCashTransactionResponse](docs/InternalShiftCashTransactionResponse.md)
 - [InternalShiftEndShiftRequest](docs/InternalShiftEndShiftRequest.md)
 - [InternalShiftShiftResponse](docs/InternalShiftShiftResponse.md)
 - [InternalShiftStartShiftRequest](docs/InternalShiftStartShiftRequest.md)
 - [InternalUserCreateUserRequest](docs/InternalUserCreateUserRequest.md)
 - [InternalUserLoginRequest](docs/InternalUserLoginRequest.md)
 - [InternalUserLoginResponse](docs/InternalUserLoginResponse.md)
 - [InternalUserProfileResponse](docs/InternalUserProfileResponse.md)
 - [InternalUserRegisterRequest](docs/InternalUserRegisterRequest.md)
 - [InternalUserUpdatePasswordRequest](docs/InternalUserUpdatePasswordRequest.md)
 - [InternalUserUpdateUserRequest](docs/InternalUserUpdateUserRequest.md)
 - [InternalUserUsersResponse](docs/InternalUserUsersResponse.md)
 - [OrdersGet200Response](docs/OrdersGet200Response.md)
 - [OrdersIdPayMidtransPost200Response](docs/OrdersIdPayMidtransPost200Response.md)
 - [OrdersPost201Response](docs/OrdersPost201Response.md)
 - [POSKasirInternalActivitylogRepositoryLogActionType](docs/POSKasirInternalActivitylogRepositoryLogActionType.md)
 - [POSKasirInternalActivitylogRepositoryLogEntityType](docs/POSKasirInternalActivitylogRepositoryLogEntityType.md)
 - [POSKasirInternalCommonErrorResponse](docs/POSKasirInternalCommonErrorResponse.md)
 - [POSKasirInternalCommonPaginationPagination](docs/POSKasirInternalCommonPaginationPagination.md)
 - [POSKasirInternalCommonSuccessResponse](docs/POSKasirInternalCommonSuccessResponse.md)
 - [POSKasirInternalOrdersRepositoryOrderStatus](docs/POSKasirInternalOrdersRepositoryOrderStatus.md)
 - [POSKasirInternalOrdersRepositoryOrderType](docs/POSKasirInternalOrdersRepositoryOrderType.md)
 - [POSKasirInternalPromotionsRepositoryDiscountType](docs/POSKasirInternalPromotionsRepositoryDiscountType.md)
 - [POSKasirInternalPromotionsRepositoryPromotionRuleType](docs/POSKasirInternalPromotionsRepositoryPromotionRuleType.md)
 - [POSKasirInternalPromotionsRepositoryPromotionScope](docs/POSKasirInternalPromotionsRepositoryPromotionScope.md)
 - [POSKasirInternalPromotionsRepositoryPromotionTargetType](docs/POSKasirInternalPromotionsRepositoryPromotionTargetType.md)
 - [POSKasirInternalShiftRepositoryCashTransactionType](docs/POSKasirInternalShiftRepositoryCashTransactionType.md)
 - [POSKasirInternalShiftRepositoryShiftStatus](docs/POSKasirInternalShiftRepositoryShiftStatus.md)
 - [POSKasirInternalUserRepositoryUserRole](docs/POSKasirInternalUserRepositoryUserRole.md)
 - [POSKasirPkgPaymentMidtransNotificationPayload](docs/POSKasirPkgPaymentMidtransNotificationPayload.md)
 - [PaymentMethodsGet200Response](docs/PaymentMethodsGet200Response.md)
 - [ProductsGet200Response](docs/ProductsGet200Response.md)
 - [ProductsIdStockHistoryGet200Response](docs/ProductsIdStockHistoryGet200Response.md)
 - [ProductsPost201Response](docs/ProductsPost201Response.md)
 - [ProductsProductIdOptionsPost201Response](docs/ProductsProductIdOptionsPost201Response.md)
 - [PromotionsGet200Response](docs/PromotionsGet200Response.md)
 - [PromotionsPost201Response](docs/PromotionsPost201Response.md)
 - [ReportsCancellationsGet200Response](docs/ReportsCancellationsGet200Response.md)
 - [ReportsCashierPerformanceGet200Response](docs/ReportsCashierPerformanceGet200Response.md)
 - [ReportsDashboardSummaryGet200Response](docs/ReportsDashboardSummaryGet200Response.md)
 - [ReportsPaymentMethodsGet200Response](docs/ReportsPaymentMethodsGet200Response.md)
 - [ReportsProductsGet200Response](docs/ReportsProductsGet200Response.md)
 - [ReportsProfitProductsGet200Response](docs/ReportsProfitProductsGet200Response.md)
 - [ReportsProfitSummaryGet200Response](docs/ReportsProfitSummaryGet200Response.md)
 - [ReportsSalesGet200Response](docs/ReportsSalesGet200Response.md)
 - [SettingsBrandingGet200Response](docs/SettingsBrandingGet200Response.md)
 - [SettingsBrandingLogoPost200Response](docs/SettingsBrandingLogoPost200Response.md)
 - [SettingsPrinterGet200Response](docs/SettingsPrinterGet200Response.md)
 - [ShiftsCashTransactionPost201Response](docs/ShiftsCashTransactionPost201Response.md)
 - [ShiftsCurrentGet200Response](docs/ShiftsCurrentGet200Response.md)
 - [UsersGet200Response](docs/UsersGet200Response.md)


<a id="documentation-for-authorization"></a>
## Documentation For Authorization

Endpoints do not require authorization.

