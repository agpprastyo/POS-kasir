import React, { createContext, useContext, useState } from "react";
import { useGetCurrentShift } from "@/hooks/useShift";
import { POSKasirInternalDtoShiftResponse } from "@/lib/api/generated";
import { useAuth } from "@/context/AuthContext";

interface ShiftContextType {
    shift: POSKasirInternalDtoShiftResponse | null | undefined;
    isLoading: boolean;
    isShiftOpen: boolean;
    openShiftModalOpen: boolean;
    setOpenShiftModalOpen: (open: boolean) => void;
    closeShiftModalOpen: boolean;
    setCloseShiftModalOpen: (open: boolean) => void;
    refreshShift: () => void;
}

const ShiftContext = createContext<ShiftContextType | undefined>(undefined);

export const ShiftProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const { isAuthenticated } = useAuth();
    const { data: shift, isLoading, refetch } = useGetCurrentShift({ enabled: isAuthenticated });

    const [openShiftModalOpen, setOpenShiftModalOpen] = useState(false);
    const [closeShiftModalOpen, setCloseShiftModalOpen] = useState(false);

    const isShiftOpen = !!shift;

    return (
        <ShiftContext.Provider value={{
            shift,
            isLoading,
            isShiftOpen,
            openShiftModalOpen,
            setOpenShiftModalOpen,
            closeShiftModalOpen,
            setCloseShiftModalOpen,
            refreshShift: refetch
        }}>
            {children}
        </ShiftContext.Provider>
    );
};

export const useShiftContext = () => {
    const context = useContext(ShiftContext);
    if (!context) {
        throw new Error("useShiftContext must be used within a ShiftProvider");
    }
    return context;
};
