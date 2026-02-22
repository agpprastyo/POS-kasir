import { MutationCache, QueryClient } from '@tanstack/react-query'
import { QueryCache } from "@tanstack/query-core";
import { toast } from "sonner";
import i18n from "@/lib/i18n";


export const queryClient = new QueryClient({
    queryCache: new QueryCache({
        onError: (error: any) => {
            if (error.response?.status === 401) return
            toast.error(i18n.t('common.error.load_failed', { message: error.message }))
        },
    }),

    mutationCache: new MutationCache({
        onError: (error: any, _variables, _context, mutation) => {
            if (mutation.options.onError) return;
            const msg = error.response?.data?.message || i18n.t("common.error.default_message")
            toast.error(msg)
        }
    })
})
