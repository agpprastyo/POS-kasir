import {
    createContext,
    useContext,
    type ReactNode,
    useMemo,
} from 'react'
import { InternalUserLoginRequest, InternalUserProfileResponse } from "@/lib/api/generated";
import { useLoginMutation, useLogoutMutation, useMeQuery } from "@/lib/api/query/auth.ts";
import { useQueryClient } from "@tanstack/react-query";

type UserProfile = InternalUserProfileResponse

type AuthContextValue = {
    user: UserProfile | null
    isAuthenticated: boolean
    isLoading: boolean
    login: (payload: InternalUserLoginRequest) => Promise<void>
    logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {

    const {
        data: userProfile,
        isLoading: isMeLoading,
        isError: isMeError,
    } = useMeQuery()

    const loginMutation = useLoginMutation()
    const logoutMutation = useLogoutMutation()
    const queryClient = useQueryClient()

    const value: AuthContextValue = useMemo(() => {

        const profile = userProfile ?? null
        const isAuthenticated = !!profile && !isMeError

        const login = async (payload: InternalUserLoginRequest) => {
            await loginMutation.mutateAsync(payload)
        }

        const logout = async () => {
            try {
                await logoutMutation.mutateAsync(undefined)
            } catch (e) {
                console.error("Logout error (server side):", e)
            } finally {
                queryClient.clear()

            }
        }

        return {
            user: profile,
            isAuthenticated,
            isLoading: isMeLoading || loginMutation.isPending || logoutMutation.isPending,
            login,
            logout,
        }
    }, [
        userProfile,
        isMeLoading,
        isMeError,
        loginMutation,
        logoutMutation,
    ])

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const ctx = useContext(AuthContext)
    if (!ctx) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return ctx
}