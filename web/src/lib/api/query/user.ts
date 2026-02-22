import { useMutation, useQuery, useQueryClient, queryOptions, keepPreviousData } from '@tanstack/react-query'
import { usersApi } from "../../api/client.ts"
import {
    InternalUserCreateUserRequest,
    InternalUserUpdateUserRequest,
    InternalUserUsersResponse,
    InternalUserProfileResponse,
    UsersGetRoleEnum,
    UsersGetSortByEnum,
    UsersGetSortOrderEnum,
    UsersGetStatusEnum,
    POSKasirInternalCommonErrorResponse,
} from "../generated"
import { toast } from "sonner";
import { AxiosError } from "axios";
import { useRBAC } from "@/lib/auth/rbac";
import i18n from '@/lib/i18n';



export type UsersListParams = {
    page?: number
    limit?: number
    search?: string
    role?: UsersGetRoleEnum
    isActive?: boolean
    status?: UsersGetStatusEnum
    sortBy?: UsersGetSortByEnum
    sortOrder?: UsersGetSortOrderEnum
}

// --- QUERY: Get All Users (/api/v1/users) ---
export const usersListQueryOptions = (params?: UsersListParams) =>
    queryOptions<
        InternalUserUsersResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['users', 'list', params],
        queryFn: async () => {
            const res = await usersApi.usersGet(
                params?.page,
                params?.limit,
                params?.search,
                params?.role,
                params?.isActive,
                params?.status,
                params?.sortBy,
                params?.sortOrder
            )
            return (res.data as any).data;
        },
        placeholderData: keepPreviousData,
    })

export const useUsersListQuery = (params?: UsersListParams) => useQuery(usersListQueryOptions(params))


// --- QUERY: Get User By ID (/api/v1/users/{id}) ---
export const userDetailQueryOptions = (id: string) =>
    queryOptions<
        InternalUserProfileResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['users', 'detail', id],
        queryFn: async () => {
            const res = await usersApi.usersIdGet(id)
            return (res.data as any).data;
        },
        enabled: !!id,
    })

export const useUserDetailQuery = (id: string) => useQuery(userDetailQueryOptions(id))



export const useCreateUserMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/users')

    return {
        ...useMutation({
            mutationKey: ['users', 'create'],
            mutationFn: async (body: InternalUserCreateUserRequest) => {
                const res = await usersApi.usersPost(body)
                return (res.data as any).data as InternalUserProfileResponse;
            },
            onSuccess: async () => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                toast.success(i18n.t('users.messages.create_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(i18n.t('users.messages.create_failed', { message: msg }))
            }
        }), isAllowed
    }
}

export const useUpdateUserMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/users/{id}')

    return {
        ...useMutation({
            mutationKey: ['users', 'update'],
            mutationFn: async ({ id, body }: { id: string; body: InternalUserUpdateUserRequest }) => {
                const res = await usersApi.usersIdPut(id, body)
                return (res.data as any).data as InternalUserProfileResponse;
            },
            onSuccess: async (data) => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                await qc.invalidateQueries({ queryKey: ['users', 'detail', data.id] })
                toast.success(i18n.t('users.messages.update_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(i18n.t('users.messages.update_failed', { message: msg }))
            }
        }), isAllowed
    }
}

export const useDeleteUserMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('DELETE', '/users/{id}')

    return {
        ...useMutation({
            mutationKey: ['users', 'delete'],
            mutationFn: async (id: string) => {
                const res = await usersApi.usersIdDelete(id)
                return (res.data as any).data;
            },
            onSuccess: async () => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                toast.success(i18n.t('users.messages.delete_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error || "Unknown error";
                toast.error(i18n.t('users.messages.delete_failed', { message: msg }))
            }
        }), isAllowed
    }
}

export const useToggleUserStatusMutation = () => {
    const qc = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/users/{id}/toggle-status')

    return {
        ...useMutation({
            mutationKey: ['users', 'toggle'],
            mutationFn: async (id: string) => {
                const res = await usersApi.usersIdToggleStatusPost(id)
                return (res.data as any).data;
            },
            onSuccess: async (_, id) => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                await qc.invalidateQueries({ queryKey: ['users', 'detail', id] })
                toast.success(i18n.t('users.messages.update_status_success'))
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                console.error("Masuk ke onError Mutation:", error)

                const errorMessage = error.response?.data?.error || "Unknown error";
                toast.error(i18n.t('users.messages.update_status_failed', { message: errorMessage }))
            }
        }), isAllowed
    }
}