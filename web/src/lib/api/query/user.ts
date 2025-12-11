import { useMutation, useQuery, useQueryClient, queryOptions, keepPreviousData } from '@tanstack/react-query'
import { usersApi } from "../../api/client.ts"
import {
    POSKasirInternalDtoCreateUserRequest,
    POSKasirInternalDtoUpdateUserRequest,
    ApiV1UsersGetRoleEnum,
    ApiV1UsersGetStatusEnum,
    ApiV1UsersGetSortByEnum,
    ApiV1UsersGetSortOrderEnum
} from "@/lib/api/generated"

// Tipe helper untuk parameter list users agar tidak terlalu panjang di function arguments
export type UsersListParams = {
    page?: number
    limit?: number
    search?: string
    role?: ApiV1UsersGetRoleEnum
    isActive?: boolean
    status?: ApiV1UsersGetStatusEnum
    sortBy?: ApiV1UsersGetSortByEnum
    sortOrder?: ApiV1UsersGetSortOrderEnum
}

// --- QUERY: Get All Users (/api/v1/users) ---
export const usersListQueryOptions = (params?: UsersListParams) =>
    queryOptions({
        queryKey: ['users', 'list', params],
        queryFn: async () => {
            const res = await usersApi.apiV1UsersGet(
                params?.page,
                params?.limit,
                params?.search,
                params?.role,
                params?.isActive,
                params?.status,
                params?.sortBy,
                params?.sortOrder
            )
            // Asumsi response wrapper: res.data.data berisi list dan pagination
            return res.data.data
        },
        // keepPreviousData berguna untuk pagination agar UI tidak flicker saat pindah halaman
        placeholderData: keepPreviousData,
    })

export const useUsersListQuery = (params?: UsersListParams) => useQuery(usersListQueryOptions(params))


// --- QUERY: Get User By ID (/api/v1/users/{id}) ---
export const userDetailQueryOptions = (id: string) =>
    queryOptions({
        queryKey: ['users', 'detail', id],
        queryFn: async () => {
            const res = await usersApi.apiV1UsersIdGet(id)
            return res.data.data
        },
        enabled: !!id, // Query tidak jalan jika ID kosong
    })

export const useUserDetailQuery = (id: string) => useQuery(userDetailQueryOptions(id))


// --- MUTATION: Create User (/api/v1/users) ---
export const useCreateUserMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['users', 'create'],
        mutationFn: async (body: POSKasirInternalDtoCreateUserRequest) => {
            const res = await usersApi.apiV1UsersPost(body)
            return res.data.data
        },
        onSuccess: async () => {
            // Refresh list user setelah create sukses
            await qc.invalidateQueries({ queryKey: ['users', 'list'] })
        },
    })
}


// --- MUTATION: Update User (/api/v1/users/{id}) ---
export const useUpdateUserMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['users', 'update'],
        mutationFn: async ({ id, body }: { id: string; body: POSKasirInternalDtoUpdateUserRequest }) => {
            const res = await usersApi.apiV1UsersIdPut(id, body)
            return res.data.data
        },
        onSuccess: async (_, variables) => {
            // Refresh list dan detail user spesifik yang diupdate
            await qc.invalidateQueries({ queryKey: ['users', 'list'] })
            await qc.invalidateQueries({ queryKey: ['users', 'detail', variables.id] })

            // Jika user mengupdate dirinya sendiri, mungkin perlu invalidate 'auth/me' juga
            // await qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
    })
}


// --- MUTATION: Delete User (/api/v1/users/{id}) ---
export const useDeleteUserMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['users', 'delete'],
        mutationFn: async (id: string) => {
            const res = await usersApi.apiV1UsersIdDelete(id)
            return res.data
        },
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: ['users', 'list'] })
        },
    })
}


// --- MUTATION: Toggle User Status (/api/v1/users/{id}/toggle) ---
export const useToggleUserStatusMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['users', 'toggle'],
        mutationFn: async (id: string) => {
            const res = await usersApi.apiV1UsersIdTogglePut(id)
            return res.data
        },
        onSuccess: async (_, id) => {
            await qc.invalidateQueries({ queryKey: ['users', 'list'] })
            await qc.invalidateQueries({ queryKey: ['users', 'detail', id] })
        },
    })
}