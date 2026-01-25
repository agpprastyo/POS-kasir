import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { productsApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    POSKasirInternalDtoCreateProductOptionRequestStandalone,
    POSKasirInternalDtoCreateProductRequest,
    POSKasirInternalDtoListProductsResponse,
    POSKasirInternalDtoProductResponse,
    POSKasirInternalDtoUpdateProductOptionRequest,
    POSKasirInternalDtoUpdateProductRequest,
    POSKasirInternalDtoRestoreBulkRequest
} from "../generated"
import { toast } from "sonner"
import { AxiosError } from "axios"


export type Product = POSKasirInternalDtoProductResponse
export type ProductListResponse = POSKasirInternalDtoListProductsResponse

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
            console.info("Product detail :", res.data)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const useProductDetailQuery = (id: string) =>
    useQuery(productDetailQueryOptions(id))

export const useCreateProductMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        POSKasirInternalDtoCreateProductRequest
    >({
        mutationKey: ['products', 'create'],
        mutationFn: async (body) => {
            const res = await productsApi.productsPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            toast.success("Produk berhasil dibuat")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal membuat produk"
            toast.error(msg)
        }
    })
}

export const useUpdateProductMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        Product,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoUpdateProductRequest }
    >({
        mutationKey: ['products', 'update'],
        mutationFn: async ({ id, body }) => {
            const res = await productsApi.productsIdPatch(id, body)
            return (res.data as any).data;
        },
        onSuccess: (data) => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'detail', data.id] })
            toast.success("Produk berhasil diperbarui")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memperbarui produk"
            toast.error(msg)
        }
    })
}


export const useDeleteProductMutation = () => {
    const qc = useQueryClient()

    return useMutation<
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
            toast.success("Produk berhasil dihapus")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal menghapus produk"
            toast.error(msg)
        }
    })
}

export const useUploadProductImageMutation = () => {
    const qc = useQueryClient()

    return useMutation<
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
            toast.success("Gambar produk berhasil diupload")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal upload gambar"
            toast.error(msg)
        }
    })
}

export const useCreateProductOptionMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { productId: string; body: POSKasirInternalDtoCreateProductOptionRequestStandalone }
    >({
        mutationKey: ['products', 'create-option'],
        mutationFn: async ({ productId, body }) => {
            const res = await productsApi.productsProductIdOptionsPost(productId, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['products', 'detail', variables.productId] })
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            toast.success("Varian produk berhasil dibuat")
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || "Gagal membuat varian")
        }
    })
}

export const useUpdateProductOptionMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { productId: string; optionId: string; body: POSKasirInternalDtoUpdateProductOptionRequest }
    >({
        mutationKey: ['products', 'update-option'],
        mutationFn: async ({ productId, optionId, body }) => {
            const res = await productsApi.productsProductIdOptionsOptionIdPatch(productId, optionId, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['products', 'detail', variables.productId] })
            toast.success("Varian berhasil diperbarui")
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || "Gagal update varian")
        }
    })
}

export const useUploadProductOptionImageMutation = () => {
    const qc = useQueryClient()

    return useMutation<
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
            toast.success("Gambar varian berhasil diupload")
        },
        onError: (error) => {
            toast.error(error.response?.data?.message || "Gagal upload gambar varian")
        }
    })
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

    return useMutation<
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
            toast.success("Produk berhasil dipulihkan")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memulihkan produk"
            toast.error(msg)
        }
    })
}

export const useRestoreBulkProductMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        POSKasirInternalDtoRestoreBulkRequest
    >({
        mutationKey: ['products', 'restore-bulk'],
        mutationFn: async (body) => {
            const res = await productsApi.productsTrashRestoreBulkPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['products', 'list'] })
            qc.invalidateQueries({ queryKey: ['products', 'trash'] })
            toast.success("Produk berhasil dipulihkan")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memulihkan produk"
            toast.error(msg)
        }
    })
}