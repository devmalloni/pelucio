# \WalletAPI

All URIs are relative to *http://localhost:8091/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**V1AdminWalletIdBurnPost**](WalletAPI.md#V1AdminWalletIdBurnPost) | **Post** /v1/admin/wallet/{id}/burn | Burn Transaction
[**V1AdminWalletIdLockPost**](WalletAPI.md#V1AdminWalletIdLockPost) | **Post** /v1/admin/wallet/{id}/lock | Lock Transaction
[**V1AdminWalletIdMintPost**](WalletAPI.md#V1AdminWalletIdMintPost) | **Post** /v1/admin/wallet/{id}/mint | Mint Transaction
[**V1AdminWalletIdMintandlockPost**](WalletAPI.md#V1AdminWalletIdMintandlockPost) | **Post** /v1/admin/wallet/{id}/mintandlock | Unlock and burn Transaction
[**V1AdminWalletIdUnlockPost**](WalletAPI.md#V1AdminWalletIdUnlockPost) | **Post** /v1/admin/wallet/{id}/unlock | Unlock Transaction
[**V1AdminWalletPost**](WalletAPI.md#V1AdminWalletPost) | **Post** /v1/admin/wallet | Create a wallet
[**V1AdminWalletsGet**](WalletAPI.md#V1AdminWalletsGet) | **Get** /v1/admin/wallets | Get Wallets
[**V1OpenWalletExternalIdGet**](WalletAPI.md#V1OpenWalletExternalIdGet) | **Get** /v1/open/wallet/external/{id} | Get Wallet by externalID
[**V1OpenWalletIdGet**](WalletAPI.md#V1OpenWalletIdGet) | **Get** /v1/open/wallet/{id} | Get Wallet
[**V1OpenWalletIdRecordsGet**](WalletAPI.md#V1OpenWalletIdRecordsGet) | **Get** /v1/open/wallet/{id}/records | Get Wallet records
[**V1OpenWalletTransferPost**](WalletAPI.md#V1OpenWalletTransferPost) | **Post** /v1/open/wallet/transfer | Transfer Transaction
[**V1OpenWalletUnlockandtransferPost**](WalletAPI.md#V1OpenWalletUnlockandtransferPost) | **Post** /v1/open/wallet/unlockandtransfer | UnlockAndTransfer Transaction



## V1AdminWalletIdBurnPost

> V1AdminWalletIdBurnPost(ctx, id).Model(model).Execute()

Burn Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id
	model := *openapiclient.NewWalletAmountModel() // WalletAmountModel | Amount model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1AdminWalletIdBurnPost(context.Background(), id).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletIdBurnPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletIdBurnPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **model** | [**WalletAmountModel**](WalletAmountModel.md) | Amount model | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletIdLockPost

> V1AdminWalletIdLockPost(ctx, id).Model(model).Execute()

Lock Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id
	model := *openapiclient.NewWalletAmountModel() // WalletAmountModel | Amount model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1AdminWalletIdLockPost(context.Background(), id).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletIdLockPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletIdLockPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **model** | [**WalletAmountModel**](WalletAmountModel.md) | Amount model | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletIdMintPost

> V1AdminWalletIdMintPost(ctx, id).Model(model).Execute()

Mint Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id
	model := *openapiclient.NewWalletAmountModel() // WalletAmountModel | Amount model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1AdminWalletIdMintPost(context.Background(), id).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletIdMintPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletIdMintPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **model** | [**WalletAmountModel**](WalletAmountModel.md) | Amount model | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletIdMintandlockPost

> V1AdminWalletIdMintandlockPost(ctx, id).Model(model).Execute()

Unlock and burn Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id
	model := *openapiclient.NewWalletAmountModel() // WalletAmountModel | Amount model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1AdminWalletIdMintandlockPost(context.Background(), id).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletIdMintandlockPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletIdMintandlockPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **model** | [**WalletAmountModel**](WalletAmountModel.md) | Amount model | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletIdUnlockPost

> V1AdminWalletIdUnlockPost(ctx, id).Model(model).Execute()

Unlock Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id
	model := *openapiclient.NewWalletAmountModel() // WalletAmountModel | Amount model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1AdminWalletIdUnlockPost(context.Background(), id).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletIdUnlockPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletIdUnlockPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **model** | [**WalletAmountModel**](WalletAmountModel.md) | Amount model | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletPost

> WalletWalletResponse V1AdminWalletPost(ctx).Model(model).Execute()

Create a wallet



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	model := *openapiclient.NewWalletCreateWalletModel() // WalletCreateWalletModel | Create wallet model

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WalletAPI.V1AdminWalletPost(context.Background()).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1AdminWalletPost`: WalletWalletResponse
	fmt.Fprintf(os.Stdout, "Response from `WalletAPI.V1AdminWalletPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **model** | [**WalletCreateWalletModel**](WalletCreateWalletModel.md) | Create wallet model | 

### Return type

[**WalletWalletResponse**](WalletWalletResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1AdminWalletsGet

> []WalletWalletResponse V1AdminWalletsGet(ctx).Execute()

Get Wallets



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WalletAPI.V1AdminWalletsGet(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1AdminWalletsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1AdminWalletsGet`: []WalletWalletResponse
	fmt.Fprintf(os.Stdout, "Response from `WalletAPI.V1AdminWalletsGet`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiV1AdminWalletsGetRequest struct via the builder pattern


### Return type

[**[]WalletWalletResponse**](WalletWalletResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1OpenWalletExternalIdGet

> WalletWalletResponse V1OpenWalletExternalIdGet(ctx, id).Execute()

Get Wallet by externalID



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | External id

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WalletAPI.V1OpenWalletExternalIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1OpenWalletExternalIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1OpenWalletExternalIdGet`: WalletWalletResponse
	fmt.Fprintf(os.Stdout, "Response from `WalletAPI.V1OpenWalletExternalIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | External id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1OpenWalletExternalIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**WalletWalletResponse**](WalletWalletResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1OpenWalletIdGet

> WalletWalletResponse V1OpenWalletIdGet(ctx, id).Execute()

Get Wallet



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WalletAPI.V1OpenWalletIdGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1OpenWalletIdGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1OpenWalletIdGet`: WalletWalletResponse
	fmt.Fprintf(os.Stdout, "Response from `WalletAPI.V1OpenWalletIdGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1OpenWalletIdGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**WalletWalletResponse**](WalletWalletResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1OpenWalletIdRecordsGet

> []WalletWalletRecord V1OpenWalletIdRecordsGet(ctx, id).Execute()

Get Wallet records



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	id := "id_example" // string | Wallet id

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.WalletAPI.V1OpenWalletIdRecordsGet(context.Background(), id).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1OpenWalletIdRecordsGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `V1OpenWalletIdRecordsGet`: []WalletWalletRecord
	fmt.Fprintf(os.Stdout, "Response from `WalletAPI.V1OpenWalletIdRecordsGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**id** | **string** | Wallet id | 

### Other Parameters

Other parameters are passed through a pointer to a apiV1OpenWalletIdRecordsGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

[**[]WalletWalletRecord**](WalletWalletRecord.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1OpenWalletTransferPost

> V1OpenWalletTransferPost(ctx).Model(model).Execute()

Transfer Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	model := *openapiclient.NewWalletTransferModel() // WalletTransferModel | Transfer data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1OpenWalletTransferPost(context.Background()).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1OpenWalletTransferPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiV1OpenWalletTransferPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **model** | [**WalletTransferModel**](WalletTransferModel.md) | Transfer data | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## V1OpenWalletUnlockandtransferPost

> V1OpenWalletUnlockandtransferPost(ctx).Model(model).Execute()

UnlockAndTransfer Transaction



### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func main() {
	model := *openapiclient.NewWalletTransferModel() // WalletTransferModel | Transfer data

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.WalletAPI.V1OpenWalletUnlockandtransferPost(context.Background()).Model(model).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletAPI.V1OpenWalletUnlockandtransferPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiV1OpenWalletUnlockandtransferPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **model** | [**WalletTransferModel**](WalletTransferModel.md) | Transfer data | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

