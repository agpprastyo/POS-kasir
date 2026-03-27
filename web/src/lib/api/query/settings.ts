import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { settingsApi, printerApi } from '@/lib/api/client'
import { InternalSettingsUpdateBrandingRequest, InternalSettingsUpdatePrinterSettingsRequest } from '../generated'
import { useRBAC } from '@/lib/auth/rbac'

export const brandingSettingsQueryKey = ['branding-settings']
export const printerSettingsQueryKey = ['printer-settings']

export function useBrandingSettingsQuery() {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/settings/branding');
    const query = useQuery({
        queryKey: brandingSettingsQueryKey,
        queryFn: async () => {
            const res = await settingsApi.settingsBrandingGet()
            return (res.data as any).data
        },
        staleTime: Infinity,
        gcTime: Infinity,
        enabled: isAllowed,
    })
    return { ...query, isAllowed }
}

export function useUpdateBrandingSettingsMutation() {
    const queryClient = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/settings/branding')

    const mutation = useMutation({
        mutationFn: async (data: InternalSettingsUpdateBrandingRequest) => {
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
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/settings/printer');
    const query = useQuery({
        queryKey: printerSettingsQueryKey,
        queryFn: async () => {
            const res = await settingsApi.settingsPrinterGet()
            return (res.data as any).data
        },
        enabled: isAllowed,
    })
    return { ...query, isAllowed }
}

export function useUpdatePrinterSettingsMutation() {
    const queryClient = useQueryClient()
    const { canAccessApi } = useRBAC()
    const isAllowed = canAccessApi('PUT', '/settings/printer')

    const mutation = useMutation({
        mutationFn: async (data: InternalSettingsUpdatePrinterSettingsRequest) => {
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

export function useDiscoverPrintersQuery(enabled: boolean = false) {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/settings/printer/discover');
    const query = useQuery({
        queryKey: ['discovered-printers'],
        queryFn: async () => {
            const res = await printerApi.settingsPrinterDiscoverGet()
            return (res.data as any).data as { ip: string, port: number, name: string }[]
        },
        enabled: enabled && isAllowed,
        staleTime: 0,
    })
    return { ...query, isAllowed }
}
