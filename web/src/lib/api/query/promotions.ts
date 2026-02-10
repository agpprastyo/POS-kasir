import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { promotionsApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    POSKasirInternalDtoCreatePromotionRequest,
    POSKasirInternalDtoUpdatePromotionRequest,
    POSKasirInternalRepositoryDiscountType,
    POSKasirInternalRepositoryPromotionRuleType,
    POSKasirInternalRepositoryPromotionScope,
    POSKasirInternalRepositoryPromotionTargetType
} from "../generated"
import { toast } from "sonner"
import { AxiosError } from "axios"
import { useRBAC } from '@/lib/auth/rbac'

export interface PromotionRuleResponse {
    id: string
    rule_type: POSKasirInternalRepositoryPromotionRuleType
    rule_value: string
    description?: string
}

export interface PromotionTargetResponse {
    id: string
    target_type: POSKasirInternalRepositoryPromotionTargetType
    target_id: string
}

export interface Promotion {
    id: string
    name: string
    description?: string
    scope: POSKasirInternalRepositoryPromotionScope
    discount_type: POSKasirInternalRepositoryDiscountType
    discount_value: number
    max_discount_amount?: number
    start_date: string
    end_date: string
    is_active: boolean
    created_at: string
    updated_at: string
    deleted_at?: string
    rules: PromotionRuleResponse[]
    targets: PromotionTargetResponse[]
}

export interface PagedPromotionResponse {
    promotions: Promotion[]
    pagination: {
        current_page: number
        total_page: number
        total_data: number
        per_page: number
    }
}

export type PromotionsListParams = {
    limit?: number
    page?: number
    trash?: boolean
}

export const promotionsListQueryOptions = (params?: PromotionsListParams) =>
    queryOptions<
        PagedPromotionResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['promotions', 'list', params],
        queryFn: async () => {
            const res = await promotionsApi.promotionsGet(
                params?.page ? params.page : undefined,
                params?.limit ? params.limit : undefined,
                params?.trash ? params.trash : undefined
            )
            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const usePromotionsListQuery = (params?: PromotionsListParams) =>
    useQuery(promotionsListQueryOptions(params))

// ... (existing helper hooks for detail/create/update/delete)

export const useRestorePromotionMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/promotions/{id}/restore')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        string
    >({
        mutationKey: ['promotions', 'restore'],
        mutationFn: async (id) => {
            const res = await promotionsApi.promotionsIdRestorePost(id)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['promotions', 'list'] })
            toast.success("Promosi berhasil dipulihkan")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memulihkan promosi"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const promotionDetailQueryOptions = (id: string) =>
    queryOptions<
        Promotion,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['promotions', 'detail', id],
        queryFn: async () => {
            const res = await promotionsApi.promotionsIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const usePromotionDetailQuery = (id: string) =>
    useQuery(promotionDetailQueryOptions(id))

export const useCreatePromotionMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/promotions')

    const mutation = useMutation<
        Promotion,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        POSKasirInternalDtoCreatePromotionRequest
    >({
        mutationKey: ['promotions', 'create'],
        mutationFn: async (body) => {
            const res = await promotionsApi.promotionsPost(body)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['promotions', 'list'] })
            toast.success("Promosi berhasil dibuat")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal membuat promosi"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useUpdatePromotionMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/promotions/{id}')

    const mutation = useMutation<
        Promotion,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoUpdatePromotionRequest }
    >({
        mutationKey: ['promotions', 'update'],
        mutationFn: async ({ id, body }) => {
            const res = await promotionsApi.promotionsIdPut(id, body)
            return (res.data as any).data;
        },
        onSuccess: (data) => {
            qc.invalidateQueries({ queryKey: ['promotions', 'list'] })
            qc.invalidateQueries({ queryKey: ['promotions', 'detail', data.id] })
            toast.success("Promosi berhasil diperbarui")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memperbarui promosi"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}

export const useDeletePromotionMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('DELETE', '/promotions/{id}')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        string
    >({
        mutationKey: ['promotions', 'delete'],
        mutationFn: async (id) => {
            const res = await promotionsApi.promotionsIdDelete(id)
            return (res.data as any).data;
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['promotions', 'list'] })
            toast.success("Promosi berhasil dihapus")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal menghapus promosi"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}
