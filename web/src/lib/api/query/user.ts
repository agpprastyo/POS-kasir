import { useMutation, useQuery, useQueryClient, queryOptions, keepPreviousData } from '@tanstack/react-query'
import { usersApi } from "../../api/client.ts"
import {
    POSKasirInternalDtoCreateUserRequest,
    POSKasirInternalDtoUpdateUserRequest,
    POSKasirInternalDtoUsersResponse,
    POSKasirInternalDtoProfileResponse,
    UsersGetRoleEnum,
    UsersGetSortByEnum,
    UsersGetSortOrderEnum,
    UsersGetStatusEnum, POSKasirInternalCommonErrorResponse,
} from "../generated"
import { toast } from "sonner";
import { AxiosError } from "axios";
import { useRBAC } from "@/lib/auth/rbac";



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
        POSKasirInternalDtoUsersResponse,
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
        POSKasirInternalDtoProfileResponse,
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
            mutationFn: async (body: POSKasirInternalDtoCreateUserRequest) => {
                const res = await usersApi.usersPost(body)
                return (res.data as any).data as POSKasirInternalDtoProfileResponse;
            },
            onSuccess: async () => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                toast.success("User created successfully")
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error
                toast.error("Gagal membuat user: " + msg)
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
            mutationFn: async ({ id, body }: { id: string; body: POSKasirInternalDtoUpdateUserRequest }) => {
                const res = await usersApi.usersIdPut(id, body)
                return (res.data as any).data as POSKasirInternalDtoProfileResponse;
            },
            onSuccess: async (data) => {
                await qc.invalidateQueries({ queryKey: ['users', 'list'] })
                await qc.invalidateQueries({ queryKey: ['users', 'detail', data.id] })
                toast.success("User updated successfully")
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error
                toast.error("Gagal memperbarui user: " + msg)
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
                toast.success("User deleted successfully")
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                const msg = error.response?.data?.error
                toast.error("Gagal menghapus user: " + msg)
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
                toast.success("Status user berhasil diubah")
            },
            onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
                console.error("Masuk ke onError Mutation:", error)

                const errorMessage = error.response?.data?.error
                toast.error("Gagal mengubah status: " + errorMessage)
            }
        }), isAllowed
    }
}