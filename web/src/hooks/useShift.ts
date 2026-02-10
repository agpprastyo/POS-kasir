import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { axiosInstance } from "@/lib/api/client";
import {
    POSKasirInternalDtoStartShiftRequest,
    POSKasirInternalDtoEndShiftRequest,
    POSKasirInternalDtoCashTransactionRequest,
    POSKasirInternalDtoShiftResponse,
    POSKasirInternalDtoCashTransactionResponse
} from "@/lib/api/generated";
import { toast } from "sonner";

export const useGetCurrentShift = (options?: { enabled?: boolean }) => {
    return useQuery({
        queryKey: ["shift", "current"],
        queryFn: async () => {
            try {
                const { data } = await axiosInstance.get<{ data: POSKasirInternalDtoShiftResponse }>("/shifts/current");
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
        mutationFn: async (req: POSKasirInternalDtoStartShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: POSKasirInternalDtoShiftResponse }>("/shifts/start", req);
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
        mutationFn: async (req: POSKasirInternalDtoEndShiftRequest) => {
            const { data } = await axiosInstance.post<{ data: POSKasirInternalDtoShiftResponse }>("/shifts/end", req);
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
        mutationFn: async (req: POSKasirInternalDtoCashTransactionRequest) => {
            const { data } = await axiosInstance.post<{ data: POSKasirInternalDtoCashTransactionResponse }>("/shifts/cash-transaction", req);
            return data.data;
        },
        onSuccess: () => {
            toast.success("Transaction recorded");
        },
    });
};
