import { useMutation, useQuery, useQueryClient, queryOptions } from '@tanstack/react-query'
import { authApi } from "../../api/client.ts";
import {
    InternalUserLoginRequest,
    InternalUserUpdatePasswordRequest,
    InternalUserProfileResponse,
    InternalUserLoginResponse,
    POSKasirInternalCommonErrorResponse
} from "@/lib/api/generated";
import { AxiosError } from "axios";
import { toast } from "sonner";
import i18n from "@/lib/i18n";

// --- QUERY: current user (/auth/me) ---
export const meQueryOptions = () =>
    queryOptions<
        InternalUserProfileResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['auth', 'me'],
        queryFn: async () => {
            const res = await authApi.authMeGet()
            return (res.data as any).data;
        },
        retry: false,
        staleTime: 1000 * 60 * 5,
    })

export const useMeQuery = () => useQuery(meQueryOptions())


// --- MUTATION: login (/auth/login) ---
export const useLoginMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        InternalUserLoginResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalUserLoginRequest
    >({
        mutationKey: ['auth', 'login'],
        mutationFn: async (body) => {
            const res = await authApi.authLoginPost(body)

            return (res.data as any).data;
        },
        onSuccess: async () => {
            await qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
    })
}


// --- MUTATION: logout (/auth/logout) ---
export const useLogoutMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        void
    >({
        mutationKey: ['auth', 'logout'],
        mutationFn: async () => {
            const res = await authApi.authLogoutPost()
            return (res.data as any).data;
        },
        onSuccess: async () => {
            qc.setQueryData(['auth', 'me'], null);
            await qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
        onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
            const msg = error.response?.data?.error || "Unknown error";
            toast.error(i18n.t('auth.logout_failed', { message: msg }));
        }
    })
}


// --- MUTATION: Update Avatar (/auth/me/avatar) ---
export const useUpdateAvatarMutation = () => {
    const qc = useQueryClient()

    return useMutation<
        InternalUserProfileResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        File
    >({
        mutationKey: ['auth', 'update-avatar'],
        mutationFn: async (file) => {
            const res = await authApi.authMeAvatarPut(file)
            return (res.data as any).data;
        },
        onSuccess: (newData) => {
            qc.setQueryData(['auth', 'me'], newData);
            return qc.invalidateQueries({ queryKey: ['auth', 'me'] })
        },
        onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
            const msg = error.response?.data?.error || "Unknown error";
            toast.error(i18n.t('auth.update_avatar_failed', { message: msg }));
        }
    })
}


// --- MUTATION: Update Password (/auth/update-password) ---
export const useUpdatePasswordMutation = () => {
    return useMutation<
        any,
        AxiosError<POSKasirInternalCommonErrorResponse>,
        InternalUserUpdatePasswordRequest
    >({
        mutationKey: ['auth', 'update-password'],
        mutationFn: async (payload) => {
            const res = await authApi.authMePasswordPut(payload)
            return (res.data as any).data;
        },
        onSuccess: async () => {
            toast.success(i18n.t('auth.update_password_success'))
        },
        onError: (error: AxiosError<POSKasirInternalCommonErrorResponse>) => {
            const msg = error.response?.data?.error || "Unknown error";
            toast.error(i18n.t('auth.update_password_failed', { message: msg }));
        }
    })
}