import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ordersApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    POSKasirInternalDtoApplyPromotionRequest,
    POSKasirInternalDtoCancelOrderRequest,
    POSKasirInternalDtoCompleteManualPaymentRequest,
    POSKasirInternalDtoCreateOrderRequest,
    POSKasirInternalDtoUpdateOrderStatusRequest,
    OrdersGetStatusEnum,
} from "../generated"
import { toast } from "sonner"
import { AxiosError } from "axios"

export type OrdersListParams = {
    limit?: number
    page?: number
    status?: string
    userId?: string
}

export const ordersListQueryOptions = (params?: OrdersListParams) =>
    queryOptions<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['orders', 'list', params],
        queryFn: async () => {
            const res = await ordersApi.ordersGet(
                params?.page ? params.page : undefined,
                params?.limit ? params.limit : undefined,
                params?.status ? (params.status as OrdersGetStatusEnum) : undefined,
                params?.userId ? params.userId : undefined,
            )

            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useOrdersListQuery = (params?: OrdersListParams) =>
    useQuery(ordersListQueryOptions(params))

export const orderDetailQueryOptions = (id: string) =>
    queryOptions<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['orders', 'detail', id],
        queryFn: async () => {
            const res = await ordersApi.ordersIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const useOrderDetailQuery = (id: string) =>
    useQuery(orderDetailQueryOptions(id))

export const useApplyPromotionMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoApplyPromotionRequest }
    >({
        mutationKey: ['orders', 'apply-promotion'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdApplyPromotionPost(id, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Promo berhasil diterapkan")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal menerapkan promo"
            toast.error(msg)
        }
    })
}

export const useCancelOrderMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoCancelOrderRequest }
    >({
        mutationKey: ['orders', 'cancel'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdCancelPost(id, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Order berhasil dibatalkan")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal membatalkan order"
            toast.error(msg)
        }
    })
}

export const useCompleteManualPaymentMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoCompleteManualPaymentRequest }
    >({
        mutationKey: ['orders', 'complete-manual-payment'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdCompleteManualPaymentPost(id, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Pembayaran manual berhasil dikonfirmasi")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal konfirmasi pembayaran manual"
            toast.error(msg)
        }
    })
}

export const useProcessPaymentMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string }
    >({
        mutationKey: ['orders', 'process-payment'],
        mutationFn: async ({ id }) => {
            const res = await ordersApi.ordersIdProcessPaymentPost(id)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Pembayaran diproses")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memproses pembayaran"
            toast.error(msg)
        }
    })
}

export const useUpdateOrderStatusMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: POSKasirInternalDtoUpdateOrderStatusRequest }
    >({
        mutationKey: ['orders', 'update-status'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdUpdateStatusPost(id, body)
            return (res.data as any).data;
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Status order diperbarui")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memperbarui status order"
            toast.error(msg)
        }
    })
}

export const useCreateOrderMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        POSKasirInternalDtoCreateOrderRequest
    >({
        mutationKey: ['orders', 'create'],
        mutationFn: async (body) => {
            const res = await ordersApi.ordersPost(body)
            return (res.data as any).data
        },
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            toast.success("Order berhasil dibuat")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal membuat order"
            toast.error(msg)
        }
    })
}

