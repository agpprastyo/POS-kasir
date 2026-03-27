import { useMutation, useQuery, useQueryClient, queryOptions, keepPreviousData } from '@tanstack/react-query'
import { useTranslation } from "react-i18next";
import { customersApi } from "../../api/client.ts"
import {
    InternalCustomersCreateCustomerRequest,
    InternalCustomersUpdateCustomerRequest,
    InternalCustomersCustomerResponse,
    InternalCustomersPagedCustomerResponse,
    POSKasirInternalCommonErrorResponse,
} from "../generated"
import { toast } from "sonner";
import { AxiosError } from "axios";
import { useRBAC } from "@/lib/auth/rbac";

export type CustomersListParams = {
    page?: number
    limit?: number
    search?: string
}

// --- QUERY: Get All Customers (/api/v1/customers) ---
export const customersListQueryOptions = (params?: CustomersListParams) =>
    queryOptions<
        InternalCustomersPagedCustomerResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['customers', 'list', params],
        queryFn: async () => {
            const res = await customersApi.customersGet(
                params?.page,
                params?.limit,
                params?.search
            )
            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useCustomersListQuery = (params?: CustomersListParams) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/customers');
    const query = useQuery({
        ...customersListQueryOptions(params),
        enabled: isAllowed
    });
    return { ...query, isAllowed };
}

// --- QUERY: Get Customer By ID (/api/v1/customers/{id}) ---
export const customerDetailQueryOptions = (id: string) =>
    queryOptions<
        InternalCustomersCustomerResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['customers', 'detail', id],
        queryFn: async () => {
            const res = await customersApi.customersIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const useCustomerDetailQuery = (id: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/customers/{id}');
    const options = customerDetailQueryOptions(id);
    const query = useQuery({
        ...options,
        enabled: options.enabled !== false ? isAllowed : false
    });
    return { ...query, isAllowed };
}

export const useCreateCustomerMutation = () => {
    const { t } = useTranslation()
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/customers')

    return {
        ...useMutation({
            mutationKey: ['customers', 'create'],
            mutationFn: async (body: InternalCustomersCreateCustomerRequest) => {
                const res = await customersApi.customersPost(body)
                return (res.data as any).data as InternalCustomersCustomerResponse;
            },
            onSuccess: async () => {
                await qc.invalidateQueries({ queryKey: ['customers', 'list'] })
                toast.success(t('customers.messages.create_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(t('customers.messages.create_failed', { message: msg }))
            }
        }), isAllowed
    }
}

export const useUpdateCustomerMutation = () => {
    const { t } = useTranslation()
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/customers/{id}')

    return {
        ...useMutation({
            mutationKey: ['customers', 'update'],
            mutationFn: async ({ id, body }: { id: string; body: InternalCustomersUpdateCustomerRequest }) => {
                const res = await customersApi.customersIdPut(id, body)
                return (res.data as any).data as InternalCustomersCustomerResponse;
            },
            onSuccess: async (variables) => {
                await qc.invalidateQueries({ queryKey: ['customers', 'list'] })
                await qc.invalidateQueries({ queryKey: ['customers', 'detail', variables.id] })
                toast.success(t('customers.messages.update_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(t('customers.messages.update_failed', { message: msg }))
            }
        }), isAllowed
    }
}

export const useDeleteCustomerMutation = () => {
    const { t } = useTranslation()
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('DELETE', '/customers/{id}')

    return {
        ...useMutation({
            mutationKey: ['customers', 'delete'],
            mutationFn: async (id: string) => {
                const res = await customersApi.customersIdDelete(id)
                return (res.data as any).data;
            },
            onSuccess: async () => {
                await qc.invalidateQueries({ queryKey: ['customers', 'list'] })
                toast.success(t('customers.messages.delete_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(t('customers.messages.delete_failed', { message: msg }))
            }
        }), isAllowed
    }
}
