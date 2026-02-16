import { keepPreviousData, queryOptions, useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ordersApi, printerApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    InternalOrdersApplyPromotionRequest,
    InternalOrdersCancelOrderRequest,
    InternalOrdersConfirmManualPaymentRequest,
    InternalOrdersCreateOrderRequest,
    InternalOrdersUpdateOrderStatusRequest,
    InternalOrdersUpdateOrderItemRequest,
    InternalOrdersOrderDetailResponse,
    InternalOrdersMidtransPaymentResponse,
} from "../generated"
import { toast } from "sonner"
import { AxiosError } from "axios"
import { useRBAC } from "@/lib/auth/rbac"

export type OrdersListParams = {
    limit?: number
    page?: number
    statuses?: string[]
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
                params?.statuses ? (params.statuses as any) : undefined,
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

import { UseQueryOptions } from '@tanstack/react-query'

export const useOrderDetailQuery = (id: string, options?: Omit<UseQueryOptions<any, AxiosError<POSKasirInternalCommonErrorResponse>>, 'queryKey' | 'queryFn'>) =>
    useQuery({
        ...orderDetailQueryOptions(id),
        ...options
    })

export const useApplyPromotionMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/apply-promotion')

    const mutation = useMutation<
        InternalOrdersOrderDetailResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalOrdersApplyPromotionRequest }
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

    return { ...mutation, isAllowed }
}

export const useCancelOrderMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/cancel')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalOrdersCancelOrderRequest }
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

    return { ...mutation, isAllowed }
}

export const useConfirmManualPaymentMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/pay-manual')

    const mutation = useMutation<
        InternalOrdersOrderDetailResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalOrdersConfirmManualPaymentRequest }
    >({
        mutationKey: ['orders', 'complete-manual-payment'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdPayManualPost(id, body)
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

    return { ...mutation, isAllowed }
}

export const useInitiateMidtransPaymentMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/pay-midtrans')

    const mutation = useMutation<
        InternalOrdersMidtransPaymentResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string }
    >({
        mutationKey: ['orders', 'process-payment'],
        mutationFn: async ({ id }) => {
            const res = await ordersApi.ordersIdPayMidtransPost(id)
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

    return { ...mutation, isAllowed }
}

export const useUpdateOrderStatusMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/update-status')

    const mutation = useMutation<
        InternalOrdersOrderDetailResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalOrdersUpdateOrderStatusRequest }
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

    return { ...mutation, isAllowed }
}

export const useCreateOrderMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders')

    const mutation = useMutation<
        InternalOrdersOrderDetailResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalOrdersCreateOrderRequest
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

    return { ...mutation, isAllowed }
}

export const useUpdateOrderItemsMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/orders/{id}/items')

    const mutation = useMutation<
        InternalOrdersOrderDetailResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string; body: InternalOrdersUpdateOrderItemRequest[] }
    >({
        mutationKey: ['orders', 'update-items'],
        mutationFn: async ({ id, body }) => {
            const res = await ordersApi.ordersIdItemsPut(id, body)
            return (res.data as any).data
        },
        onSuccess: (_, variables) => {
            qc.invalidateQueries({ queryKey: ['orders', 'list'] })
            qc.invalidateQueries({ queryKey: ['orders', 'detail', variables.id] })
            toast.success("Item order berhasil diperbarui")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Gagal memperbarui item order"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}


export const usePrintInvoiceMutation = () => {
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/orders/{id}/print')

    const mutation = useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        { id: string }
    >({
        mutationKey: ['orders', 'print'],
        mutationFn: async ({ id }) => {
            const res = await printerApi.ordersIdPrintPost(id)
            return (res.data as any).data
        },
        onSuccess: () => {
            toast.success("Printing invoice...")
        },
        onError: (error) => {
            const msg = error.response?.data?.message || "Failed to print invoice"
            toast.error(msg)
        }
    })

    return { ...mutation, isAllowed }
}
