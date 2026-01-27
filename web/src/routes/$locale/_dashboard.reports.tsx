import { createFileRoute, redirect, RegisteredRouter, } from '@tanstack/react-router'
import { useTranslation } from 'react-i18next'
import { meQueryOptions } from '@/lib/api/query/auth'
import { queryClient } from '@/lib/queryClient'
import { POSKasirInternalRepositoryUserRole } from '@/lib/api/generated/models/poskasir-internal-repository-user-role'
import { FileRouteByToPath } from "@tanstack/router-core/src/routeInfo.ts";

export const Route = createFileRoute('/$locale/_dashboard/reports' as FileRouteByToPath<any, any>)({
    beforeLoad: async () => {
        const user = await queryClient.ensureQueryData(meQueryOptions())

        const allowedRoles = [
            POSKasirInternalRepositoryUserRole.UserRoleAdmin,
            POSKasirInternalRepositoryUserRole.UserRoleManager
        ]

        if (!user.role || !allowedRoles.includes(user.role)) {
            console.log('Unauthorized access attempt by:', user.role)

            throw redirect({
                to: '/',
                search: {
                    error: 'Unauthorized'
                }
            } as RegisteredRouter)
        }
    },
    component: ReportsPage,
})

function ReportsPage() {
    const { t } = useTranslation()
    return <div>{t('dashboard.secret_reports')}</div>
}

