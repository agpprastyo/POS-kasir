import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { categoriesApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    InternalCategoriesCategoryResponse,
    InternalCategoriesCreateCategoryRequest
} from "../generated"
import { toast } from "sonner"
import { AxiosError } from "axios"
import { useRBAC } from '@/lib/auth/rbac'
import i18n from '@/lib/i18n'

// Gunakan tipe dari generated, atau definisikan ulang tanpa description
export type Category = InternalCategoriesCategoryResponse

export type CategoriesListParams = {
    limit?: number
    offset?: number
}

// --- QUERY: Get All Categories ---
export const categoriesListQueryOptions = (params?: CategoriesListParams) =>
    queryOptions<
        Category[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['categories', 'list', params],
        queryFn: async () => {
            const res = await categoriesApi.categoriesGet(
                params?.limit,
                params?.offset
            )
            // Unwrapping data
            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useCategoriesListQuery = (params?: CategoriesListParams) =>
    useQuery(categoriesListQueryOptions(params))


// --- QUERY: Get Category By ID ---
export const categoryDetailQueryOptions = (id: number) =>
    queryOptions<
        Category,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['categories', 'detail', id],
        queryFn: async () => {
            const res = await categoriesApi.categoriesIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })


// --- MUTATION: Create Category ---
export const useCreateCategoryMutation = () => {
    const qc = useQueryClient()

    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/categories')

    const mutation = useMutation<
        Category,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalCategoriesCreateCategoryRequest
    >({
        mutationKey: ['categories', 'create'],
        mutationFn: async (body) => {
            const res = await categoriesApi.categoriesPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['categories', 'list'] })
            toast.success(i18n.t('category.create_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal membuat kategori"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}


// --- MUTATION: Update Category ---
export const useUpdateCategoryMutation = () => {
    const qc = useQueryClient()

    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/categories/{id}')

    const mutation = useMutation<
        Category,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: number; body: InternalCategoriesCreateCategoryRequest }
    >({
        mutationKey: ['categories', 'update'],
        mutationFn: async ({ id, body }) => {
            const res = await categoriesApi.categoriesIdPut(id, body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['categories', 'list'] })
            toast.success(i18n.t('category.update_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal update kategori"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}


// --- MUTATION: Delete Category ---
export const useDeleteCategoryMutation = () => {
    const qc = useQueryClient()

    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('DELETE', '/categories/{id}')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        number
    >({
        mutationKey: ['categories', 'delete'],
        mutationFn: async (id) => {
            const res = await categoriesApi.categoriesIdDelete(id)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['categories', 'list'] })
            toast.success(i18n.t('category.delete_success'))
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal menghapus kategori"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}