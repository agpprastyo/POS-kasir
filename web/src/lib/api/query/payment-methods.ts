import { queryOptions, useQuery } from '@tanstack/react-query'
import { paymentMethodsApi } from "../client"
import {
    POSKasirInternalCommonErrorResponse,
    POSKasirInternalDtoPaymentMethodResponse,
} from "../generated"
import { AxiosError } from "axios"

export const paymentMethodsListQueryOptions = () =>
    queryOptions<
        POSKasirInternalDtoPaymentMethodResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['payment-methods', 'list'],
        queryFn: async () => {
            const res = await paymentMethodsApi.paymentMethodsGet()
            return (res.data as any).data;
        },
    })

export const usePaymentMethodsListQuery = () =>
    useQuery(paymentMethodsListQueryOptions())
