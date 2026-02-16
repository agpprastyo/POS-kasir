import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
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
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (req: InternalShiftStartShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftShiftResponse }>("/shifts/start", req);
            return data.data;
        },
        onSuccess: (data) => {
            toast.success("Shift started successfully");
            queryClient.setQueryData(["shift", "current"], data);
        },
    });
};

export const useEndShift = () => {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (req: InternalShiftEndShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftShiftResponse }>("/shifts/end", req);
            return data.data;
        },
        onSuccess: () => {
            toast.success("Shift ended successfully");
            queryClient.setQueryData(["shift", "current"], null);
            queryClient.invalidateQueries({ queryKey: ["shift"] });
        },
    });
};

export const useCreateCashTransaction = () => {
    return useMutation({
        mutationFn: async (req: InternalShiftCashTransactionRequest) => {
            const { data } = await axiosInstance.post<{ data: InternalShiftCashTransactionResponse }>("/shifts/cash-transaction", req);
            return data.data;
        },
        onSuccess: () => {
            toast.success("Transaction recorded");
        },
    });
};
