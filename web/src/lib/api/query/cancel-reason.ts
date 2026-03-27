import { queryOptions, useQuery } from '@tanstack/react-query'

import {
    InternalCancellationReasonsCancellationReasonResponse,
    POSKasirInternalCommonErrorResponse
} from "../generated"
import { AxiosError } from "axios"
import { cancellationReasonsApi } from "@/lib/api/client.ts";
import { useRBAC } from '@/lib/auth/rbac';

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

export const useCancellationReasonsListQuery = () => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/cancellation-reasons');
    const query = useQuery({
        ...cancellationReasonsListQueryOptions(),
        enabled: isAllowed
    });
    return { ...query, isAllowed };
}