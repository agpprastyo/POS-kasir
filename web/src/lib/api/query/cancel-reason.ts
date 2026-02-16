import { queryOptions, useQuery } from '@tanstack/react-query'

import {
    InternalCancellationReasonsCancellationReasonResponse,
    POSKasirInternalCommonErrorResponse
} from "../generated"
import { AxiosError } from "axios"
import { cancellationReasonsApi } from "@/lib/api/client.ts";

// --- QUERY: List Cancellation Reasons (/api/v1/cancellation-reasons) ---
export const cancellationReasonsListQueryOptions = () =>
    queryOptions<
        InternalCancellationReasonsCancellationReasonResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['cancellation-reasons', 'list'],
        queryFn: async () => {
            const res = await cancellationReasonsApi.cancellationReasonsGet()

            return (res.data as any).data;
        },

        staleTime: 1000 * 60 * 30,
    })

export const useCancellationReasonsListQuery = () => useQuery(cancellationReasonsListQueryOptions())