import {createFileRoute, redirect, RegisteredRouter,} from '@tanstack/react-router'
import { meQueryOptions } from '@/lib/api/query/auth'
import { queryClient } from '@/lib/queryClient'
import { POSKasirInternalRepositoryUserRole } from '@/lib/api/generated/models/poskasir-internal-repository-user-role'
import {FileRouteByToPath} from "@tanstack/router-core/src/routeInfo.ts";

export const Route = createFileRoute('/_dashboard/reports' as FileRouteByToPath<any, any>)({
    beforeLoad: async () => {

        const user = await queryClient.ensureQueryData(meQueryOptions())


        if (user.role !== POSKasirInternalRepositoryUserRole.UserRoleAdmin) {
            console.log('Unauthorized access attempt by:', user.role)

            throw redirect({
                to: '/',
                search: {
                    error: 'Unauthorized'
                }
            } as RegisteredRouter )
        }
    },
    component: ReportsPage,
})

function ReportsPage() {
    return <div>Halaman Laporan Rahasia (Admin Only)</div>
}

