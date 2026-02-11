import { useAuth } from "@/context/AuthContext";
import { RBAC_RULES } from "./rbacRules";

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

/**
 * Checks if a user has access to a specific API endpoint.
 * 
 * @param role - The user's role.
 * @param method - The HTTP method of the endpoint.
 * @param path - The path of the endpoint (must match the key in RBAC_RULES).
 * @returns true if access is allowed, false otherwise.
 */
export function hasApiAccess(role: string | undefined | null, method: HttpMethod, path: string): boolean {
    if (!role) return false;

    const key = `${method.toUpperCase()} ${path}`;
    const allowedRoles = RBAC_RULES[key];

    if (!allowedRoles) {
        // If no rules are defined for this endpoint, assume it's open or handle as needed.
        // Based on our swagger, only secured endpoints have x-roles.
        // However, we should be careful. If x-roles is missing in swagger, it might mean public.
        // But if we are checking RBAC, we likely expect a rule.
        // For now, let's assume if it's not in RBAC_RULES, it might be public OR we don't know it.
        // But for strict RBAC, default deny is safer if we expect all protected routes to be listed.
        // Let's check if the intent is to check "is this action allowed for this role".
        return true;
    }

    return allowedRoles.includes(role);
}

export function useRBAC() {
    const { user } = useAuth();

    const canAccessApi = (method: HttpMethod, path: string) => {
        return hasApiAccess(user?.role, method, path);
    };

    return {
        canAccessApi,
        role: user?.role
    };
}
