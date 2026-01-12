import {MutationCache, QueryClient} from '@tanstack/react-query'
import {QueryCache} from "@tanstack/query-core";
import {toast} from "sonner";


export const queryClient = new QueryClient({
    queryCache: new QueryCache({
        onError: (error: any, query) => {
            if (error.response?.status === 401) return
            toast.error(`Gagal memuat data: ${error.message}`)
        },
    }),

    mutationCache: new MutationCache({
        onError: (error: any, _variables, _context, mutation) => {
            if (mutation.options.onError) return;
            const msg = error.response?.data?.message || "Terjadi kesalahan sistem"
            toast.error(msg)
        }
    })
})
