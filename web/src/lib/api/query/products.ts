import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { productsApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    InternalProductsCreateProductOptionRequestStandalone,
    InternalProductsCreateProductRequest,
    InternalProductsListProductsResponse,
    InternalProductsProductResponse,
    InternalProductsUpdateProductOptionRequest,
    InternalProductsUpdateProductRequest,
    InternalProductsRestoreBulkRequest,
    InternalProductsPagedStockHistoryResponse
} from "../generated"
import { toast } from "sonner"
import { useTranslation } from "react-i18next"
import { AxiosError } from "axios"
import { useRBAC } from '@/lib/auth/rbac'


export type Product = InternalProductsProductResponse
export type ProductListResponse = InternalProductsListProductsResponse

export type ProductsListParams = {
    limit?: number
    page?: number
    search?: string
    category?: number

}

export const productsListQueryOptions = (params?: ProductsListParams) =>
    queryOptions<
        ProductListResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['products', 'list', params],
        queryFn: async () => {
            const res = await productsApi.productsGet(
                params?.page ? params.page : undefined,
                params?.limit ? params.limit : undefined,
                params?.search ? params.search : undefined,
                params?.category ? params.category : undefined,
            )

            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useProductsListQuery = (params?: ProductsListParams) =>
    useQuery(productsListQueryOptions(params))

export const productDetailQueryOptions = (id: string) =>
    queryOptions<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['products', 'detail', id],
        queryFn: async () => {
            const res = await productsApi.productsIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const useProductDetailQuery = (id: string) =>
    useQuery(productDetailQueryOptions(id))

export const useCreateProductMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products')

    const mutation = useMutation<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalProductsCreateProductRequest
    >({
        mutationKey: ['products', 'create'],
        mutationFn: async (body) => {
            const res = await productsApi.productsPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            toast.success(t('products.messages.create_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.create_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useUpdateProductMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PATCH', '/products/{id}')

    const mutation = useMutation<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalProductsUpdateProductRequest }
    >({
        mutationKey: ['products', 'update'],
        mutationFn: async ({ id, body }) => {
            const res = await productsApi.productsIdPatch(id, body)
            return (res.data as any).data;
        },
        onSuccess: (data) => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'detail', data.id] })
            toast.success(t('products.messages.update_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.update_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}


export const useDeleteProductMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('DELETE', '/products/{id}')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        string
    >({
        mutationKey: ['products', 'delete'],
        mutationFn: async (id) => {
            const res = await productsApi.productsIdDelete(id)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            toast.success(t('products.messages.delete_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.delete_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useUploadProductImageMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products/{id}/image')

    const mutation = useMutation<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; file: File }
    >({
        mutationKey: ['products', 'upload-image'],
        mutationFn: async ({ id, file }) => {
            const res = await productsApi.productsIdImagePost(id, file)
            return (res.data as any).data;
        },
        onSuccess: (data) => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'detail', data.id] })
            toast.success(t('products.messages.upload_image_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.upload_image_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useCreateProductOptionMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products/{product_id}/options')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { productId: string; body: InternalProductsCreateProductOptionRequestStandalone }
    >({
        mutationKey: ['products', 'create-option'],
        mutationFn: async ({ productId, body }) => {
            const res = await productsApi.productsProductIdOptionsPost(productId, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['products', 'detail', variables.productId] })
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            toast.success(t('products.messages.create_variant_success'))
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || t('products.messages.create_variant_error'))
        }
    })

    return { ...mutation, isAllowed }
}

export const useUpdateProductOptionMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PATCH', '/products/{product_id}/options/{option_id}')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { productId: string; optionId: string; body: InternalProductsUpdateProductOptionRequest }
    >({
        mutationKey: ['products', 'update-option'],
        mutationFn: async ({ productId, optionId, body }) => {
            const res = await productsApi.productsProductIdOptionsOptionIdPatch(productId, optionId, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['products', 'detail', variables.productId] })
            toast.success(t('products.messages.update_variant_success'))
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || t('products.messages.update_variant_error'))
        }
    })

    return { ...mutation, isAllowed }
}

export const useUploadProductOptionImageMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products/{product_id}/options/{option_id}/image')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { productId: string; optionId: string; file: File }
    >({
        mutationKey: ['products', 'upload-option-image'],
        mutationFn: async ({ productId, optionId, file }) => {
            const res = await productsApi.productsProductIdOptionsOptionIdImagePost(productId, optionId, file)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['products', 'detail', variables.productId] })
            toast.success(t('products.messages.upload_variant_image_success'))
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || t('products.messages.upload_variant_image_error'))
        }
    })

    return { ...mutation, isAllowed }
}

export const trashProductsListQueryOptions = (params?: ProductsListParams) =>
    queryOptions<
        ProductListResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['products', 'trash', params],
        queryFn: async () => {
            const res = await productsApi.productsTrashGet(
                params?.page ? params.page : undefined,
                params?.limit ? params.limit : undefined,
                params?.search ? params.search : undefined,
                params?.category ? params.category : undefined,
            )

            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useTrashProductsListQuery = (params?: ProductsListParams) =>
    useQuery(trashProductsListQueryOptions(params))

export const useRestoreProductMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products/trash/{id}/restore')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        string
    >({
        mutationKey: ['products', 'restore'],
        mutationFn: async (id) => {
            const res = await productsApi.productsTrashIdRestorePost(id)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'trash'] })
            toast.success(t('products.messages.restore_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.restore_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useRestoreBulkProductMutation = () => {
    const qc = useQueryClient()
    const { t } = useTranslation()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/products/trash/restore-bulk')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalProductsRestoreBulkRequest
    >({
        mutationKey: ['products', 'restore-bulk'],
        mutationFn: async (body) => {
            const res = await productsApi.productsTrashRestoreBulkPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'trash'] })
            toast.success(t('products.messages.restore_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || t('products.messages.restore_error')
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}
export type StockHistoryParams = {
    page?: number
    limit?: number
}

export const stockHistoryQueryOptions = (productId: string, params?: StockHistoryParams) =>
    queryOptions<
        InternalProductsPagedStockHistoryResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['products', 'stock-history', productId, params],
        queryFn: async () => {
            const res = await productsApi.productsIdStockHistoryGet(
                productId,
                params?.page,
                params?.limit
            )
            return (res.data as any).data;
        },
        enabled: !!productId,
        placeholderData: keepPreviousData,
    })

export const useStockHistoryQuery = (productId: string, params?: StockHistoryParams) =>
    useQuery(stockHistoryQueryOptions(productId, params))
