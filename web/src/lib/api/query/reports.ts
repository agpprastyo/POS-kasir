import { queryOptions, useQuery } from '@tanstack/react-query'
import { reportsApi } from "../../api/client"
import {
    POSKasirInternalCommonErrorResponse,
    InternalReportCancellationReportResponse,
    InternalReportCashierPerformanceResponse,
    InternalReportDashboardSummaryResponse,
    InternalReportPaymentMethodPerformanceResponse,
    InternalReportProductPerformanceResponse,
    InternalReportSalesReport,
    InternalReportProfitSummaryResponse,
    InternalReportProductProfitResponse
} from "../generated"
import { AxiosError } from "axios"

export const dashboardSummaryQueryOptions = () =>
    queryOptions<
        InternalReportDashboardSummaryResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'dashboard-summary'],
        queryFn: async () => {
            const res = await reportsApi.reportsDashboardSummaryGet();
            return (res.data as any).data;
        },
    })

export const useDashboardSummaryQuery = () =>
    useQuery(dashboardSummaryQueryOptions())


export const salesReportQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportSalesReport[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'sales', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsSalesGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useSalesReportQuery = (startDate: string, endDate: string) =>
    useQuery(salesReportQueryOptions(startDate, endDate))


export const productPerformanceQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportProductPerformanceResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'products', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsProductsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useProductPerformanceQuery = (startDate: string, endDate: string) =>
    useQuery(productPerformanceQueryOptions(startDate, endDate))


export const paymentMethodPerformanceQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportPaymentMethodPerformanceResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'payment-methods', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsPaymentMethodsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const usePaymentMethodPerformanceQuery = (startDate: string, endDate: string) =>
    useQuery(paymentMethodPerformanceQueryOptions(startDate, endDate))


export const cashierPerformanceQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportCashierPerformanceResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'cashier-performance', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsCashierPerformanceGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useCashierPerformanceQuery = (startDate: string, endDate: string) =>
    useQuery(cashierPerformanceQueryOptions(startDate, endDate))


export const cancellationReportsQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportCancellationReportResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'cancellations', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsCancellationsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useCancellationReportsQuery = (startDate: string, endDate: string) =>
    useQuery(cancellationReportsQueryOptions(startDate, endDate))

export const profitSummaryQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportProfitSummaryResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'profit-summary', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsProfitSummaryGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useProfitSummaryQuery = (startDate: string, endDate: string) =>
    useQuery(profitSummaryQueryOptions(startDate, endDate))


export const productProfitReportsQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportProductProfitResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'profit-products', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsProfitProductsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useProductProfitReportsQuery = (startDate: string, endDate: string) =>
    useQuery(productProfitReportsQueryOptions(startDate, endDate))
