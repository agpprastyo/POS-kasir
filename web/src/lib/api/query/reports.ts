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
    InternalReportProductProfitResponse,
    InternalReportLowStockProductResponse,
    InternalReportShiftSummaryResponse,
    InternalPromotionsPromotionResponse
} from "../generated"
import { AxiosError } from "axios"
import { useRBAC } from '@/lib/auth/rbac'

export const dashboardSummaryQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportDashboardSummaryResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'dashboard-summary', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsDashboardSummaryGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useDashboardSummaryQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/dashboard-summary');
    const defaultOptions = dashboardSummaryQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


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

export const useSalesReportQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/sales');
    const defaultOptions = salesReportQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


export const productPerformanceQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportProductPerformanceResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'products', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsProductsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useProductPerformanceQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/products');
    const defaultOptions = productPerformanceQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


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

export const usePaymentMethodPerformanceQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/payment-methods');
    const defaultOptions = paymentMethodPerformanceQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


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

export const useCashierPerformanceQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/cashier-performance');
    const defaultOptions = cashierPerformanceQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


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

export const useCancellationReportsQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/cancellations');
    const defaultOptions = cancellationReportsQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}

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

export const useProfitSummaryQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/profit-summary');
    const defaultOptions = profitSummaryQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


export const productProfitReportsQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportProductProfitResponse,
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'profit-products', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsProfitProductsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useProductProfitReportsQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/profit-products');
    const defaultOptions = productProfitReportsQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


export const lowStockReportQueryOptions = (threshold?: number) =>
    queryOptions<
        InternalReportLowStockProductResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'low-stock', threshold],
        queryFn: async () => {
            const res = await reportsApi.reportsLowStockGet(threshold);
            return (res.data as any).data;
        },
    })

export const useLowStockReportQuery = (threshold?: number) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/low-stock');
    const defaultOptions = lowStockReportQueryOptions(threshold);
    const query = useQuery({
        ...defaultOptions,
        enabled: isAllowed
    });
    return { ...query, isAllowed };
}


export const promotionsReportQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalPromotionsPromotionResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'promotions', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsPromotionsGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const usePromotionsReportQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/promotions');
    const defaultOptions = promotionsReportQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}


export const shiftSummaryReportQueryOptions = (startDate: string, endDate: string) =>
    queryOptions<
        InternalReportShiftSummaryResponse[],
        AxiosError<POSKasirInternalCommonErrorResponse>
    >({
        queryKey: ['reports', 'shift-summary', startDate, endDate],
        queryFn: async () => {
            const res = await reportsApi.reportsShiftSummaryGet(startDate, endDate);
            return (res.data as any).data;
        },
        enabled: !!startDate && !!endDate,
    })

export const useShiftSummaryReportQuery = (startDate: string, endDate: string) => {
    const { canAccessApi } = useRBAC();
    const isAllowed = canAccessApi('GET', '/reports/shift-summary');
    const defaultOptions = shiftSummaryReportQueryOptions(startDate, endDate);
    const query = useQuery({ ...defaultOptions, enabled: defaultOptions.enabled !== false ? isAllowed : false });
    return { ...query, isAllowed };
}
