import { MutationCache, QueryClient } from '@tanstack/react-query'
import { QueryCache } from "@tanstack/query-core";
import { toast } from "sonner";
import i18n from "@/lib/i18n";


export const makeQueryClient = () => {
    return new QueryClient({
        queryCache: new QueryCache({
            onError: (error: any) => {
                if (typeof window === 'undefined') return;
                if (error.response?.status === 401) return;
                toast.error(i18n.t('common.error.load_failed', { message: error.message }))
            },
        }),
        mutationCache: new MutationCache({
            onError: (error: any, _variables, _context, mutation) => {
                if (typeof window === 'undefined') return;
                if (mutation.options.onError) return;
                const msg = error.response?.data?.message || i18n.t("common.error.default_message")
                toast.error(msg)
            }
        }),
        defaultOptions: {
            queries: {
                staleTime: 1000 * 60,
                refetchOnWindowFocus: false,
                retry: false,
                // On the server we might want to prevent throwing to error boundaries or suspending indefinitely for failed API requests
                throwOnError: typeof window === 'undefined' ? false : undefined,
            }
        }
    })
}
