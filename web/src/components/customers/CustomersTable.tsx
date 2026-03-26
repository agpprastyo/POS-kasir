import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { InternalCustomersCustomerResponse } from '@/lib/api/generated'

interface CustomersTableProps {
    customers: InternalCustomersCustomerResponse[]
    t: any
    renderActions: (customer: InternalCustomersCustomerResponse) => React.ReactNode
}

export function CustomersTable({ customers, t, renderActions }: CustomersTableProps) {
    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead>{t('customers.table.name', 'Name')}</TableHead>
                        <TableHead>{t('customers.table.phone', 'Phone')}</TableHead>
                        <TableHead>{t('customers.table.email', 'Email')}</TableHead>
                        <TableHead>{t('customers.table.address', 'Address')}</TableHead>
                        <TableHead className="text-right">{t('customers.table.actions', 'Actions')}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {customers.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={5} className="h-24 text-center">{t('customers.table.empty', 'No customers found.')}</TableCell>
                        </TableRow>
                    ) : (
                        customers.map((customer) => (
                            <TableRow key={customer.id}>
                                <TableCell className="font-medium">{customer.name}</TableCell>
                                <TableCell>{customer.phone || '-'}</TableCell>
                                <TableCell>{customer.email || '-'}</TableCell>
                                <TableCell>{customer.address || '-'}</TableCell>
                                <TableCell className="text-right">
                                    {renderActions(customer)}
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    )
}
