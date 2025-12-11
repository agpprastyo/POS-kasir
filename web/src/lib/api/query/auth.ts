import { useMutation, useQuery, useQueryClient, queryOptions } from '@tanstack/react-query'
import {authApi} from "../../api/client.ts";
import {POSKasirInternalDtoLoginRequest, POSKasirInternalDtoUpdatePasswordRequest} from "@/lib/api/generated";


// --- QUERY: current user (/auth/profile) ---
export const meQueryOptions = () =>
    queryOptions({
        queryKey: ['auth', 'me'],
        queryFn: async () => {
            const res = await authApi.authMeGet()
            return res.data.data
        },
        retry: false,
    })

export const useMeQuery = () => useQuery(meQueryOptions())

// --- MUTATION: login (/auth/login) ---
export const useLoginMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['auth', 'login'],
        mutationFn: async (body: POSKasirInternalDtoLoginRequest) => {
            const res = await authApi.authLoginPost(body)
            return res.data.data
        },
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
    })
}

// --- MUTATION: logout (/auth/logout) ---
export const useLogoutMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['auth', 'logout'],
        mutationFn: async () => {
            const res = await authApi.authLogoutPost()
            return res.data.data
        },
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
    })
}


// --- MUTATION: Update Avatar (/auth/avatar) ---
export const useUpdateAvatarMutation = () => {
    const qc = useQueryClient()

    return useMutation({
        mutationKey: ['auth', 'update-avatar'],
        mutationFn: async (file: File) => {
            const res = await authApi.authMeAvatarPut(file)
            return res.data
        },
        onSuccess: () => {
            return qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
    })
}

// --- MUTATION: Update Password (/auth/update-password) ---
export const useUpdatePasswordMutation = () => {
    return useMutation({
        mutationKey: ['auth', 'update-password'],
        mutationFn: async (payload: POSKasirInternalDtoUpdatePasswordRequest) => {
            const res = await authApi.authUpdatePasswordPost(payload)
            return res.data
        },
    })
}