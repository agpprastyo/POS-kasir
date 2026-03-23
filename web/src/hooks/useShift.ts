import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { axiosInstance } from "@/lib/api/client";
import {
    InternalShiftStartShiftRequest,
    InternalShiftEndShiftRequest,
    InternalShiftCashTransactionRequest,
    InternalShiftShiftResponse,
    InternalShiftCashTransactionResponse
} from "@/lib/api/generated";
import { toast } from "sonner";

export const useGetCurrentShift = (options?: { enabled?: boolean }) => {
    return useQuery({
        queryKey: ["shift", "current"],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: InternalShiftShiftResponse }>("/shifts/current");
                return data.data;
            } catch (error: any) {
                if (error.response?.status === 404) {
                    return null;
                }
                throw error;
            }
        },
        retry: false,
        enabled: options?.enabled,
    });
};

export const useStartShift = () => {
    const { t } = useTranslation();
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (req: InternalShiftStartShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftShiftResponse }>("/shifts/start", req);
            return data.data;
        },
        onSuccess: (data) => {
            toast.success(t("shift.messages.start_success"));
            queryClient.setQueryData(["shift", "current"], data);
        },
    });
};

export const useEndShift = () => {
    const { t } = useTranslation();
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (req: InternalShiftEndShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftShiftResponse }>("/shifts/end", req);
            return data.data;
        },
        onSuccess: () => {
            toast.success(t("shift.messages.end_success"));
            queryClient.setQueryData(["shift", "current"], null);
            queryClient.invalidateQueries({ queryKey: ["shift"] });
        },
    });
};

export const useCreateCashTransaction = () => {
    const { t } = useTranslation();
    return useMutation({
        mutationFn: async (req: InternalShiftCashTransactionRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftCashTransactionResponse }>("/shifts/cash-transaction", req);
            return data.data;
        },
        onSuccess: () => {
            toast.success(t("shift.messages.transaction_recorded"));
        },
    });
};
