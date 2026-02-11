import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { settingsApi, printerApi } from '@/lib/api/client'
import { POSKasirInternalDtoUpdateBrandingRequest, POSKasirInternalDtoUpdatePrinterSettingsRequest } from '../generated'
import { useRBAC } from '@/lib/auth/rbac'

export const brandingSettingsQueryKey = ['branding-settings']
export const printerSettingsQueryKey = ['printer-settings']

export function useBrandingSettingsQuery() {
    return useQuery({
        queryKey: brandingSettingsQueryKey,
        queryFn: async () => {
            const res = await settingsApi.settingsBrandingGet()
            // The generated client expects the DTO directly, but the API returns a wrapper { message, data }
            // So we cast to any to access the data property
            return (res.data as any).data
        },
    })
}

export function useUpdateBrandingSettingsMutation() {
    const queryClient = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/settings/branding')

    const mutation = useMutation({
        mutationFn: async (data: POSKasirInternalDtoUpdateBrandingRequest) => {
            const res = await settingsApi.settingsBrandingPut(data)
            return (res.data as any).data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: brandingSettingsQueryKey })
        },
    })

    return { ...mutation, isAllowed }
}

export function useUpdateLogoMutation() {
    const queryClient = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/settings/branding/logo')

    const mutation = useMutation({
        mutationFn: async (file: File) => {
            // settingsBrandingLogoPost takes a File object directly as per generated code
            const res = await settingsApi.settingsBrandingLogoPost(file)
            return (res.data as any).data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: brandingSettingsQueryKey })
        },
    })

    return { ...mutation, isAllowed }
}

export function usePrinterSettingsQuery() {
    return useQuery({
        queryKey: printerSettingsQueryKey,
        queryFn: async () => {
            const res = await settingsApi.settingsPrinterGet()
            return (res.data as any).data
        },
    })
}

export function useUpdatePrinterSettingsMutation() {
    const queryClient = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/settings/printer')

    const mutation = useMutation({
        mutationFn: async (data: POSKasirInternalDtoUpdatePrinterSettingsRequest) => {
            const res = await settingsApi.settingsPrinterPut(data)
            return (res.data as any).data
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: printerSettingsQueryKey })
        },
    })

    return { ...mutation, isAllowed }
}

export function useTestPrintMutation() {
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('POST', '/settings/printer/test')

    const mutation = useMutation({
        mutationFn: async () => {
            const res = await printerApi.settingsPrinterTestPost()
            return (res.data as any).data
        },
    })

    return { ...mutation, isAllowed }
}
