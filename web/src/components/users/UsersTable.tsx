import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { InternalUserProfileResponse } from '@/lib/api/generated'

interface UsersTableProps {
    users: InternalUserProfileResponse[]
    t: any
    renderActions: (user: InternalUserProfileResponse) => React.ReactNode
}

export function UsersTable({ users, t, renderActions }: UsersTableProps) {
    return (
        <div className="rounded-md border bg-card">
            <Table>
                <TableHeader>
                    <TableRow>
                        <TableHead className="w-20">{t('users.table.avatar')}</TableHead>
                        <TableHead>{t('users.table.user')}</TableHead>
                        <TableHead>{t('users.table.role')}</TableHead>
                        <TableHead>{t('users.table.status')}</TableHead>
                        <TableHead className="text-right">{t('users.table.actions')}</TableHead>
                    </TableRow>
                </TableHeader>
                <TableBody>
                    {users.length === 0 ? (
                        <TableRow>
                            <TableCell colSpan={5} className="h-24 text-center">{t('users.table.empty')}</TableCell>
                        </TableRow>
                    ) : (
                        users.map((user) => (
                            <TableRow key={user.id}>
                                <TableCell>
                                    <Avatar className="h-9 w-9">
                                        <AvatarImage src={user.avatar || undefined} alt={user.username} />
                                        <AvatarFallback>{user.username?.slice(0, 2).toUpperCase()}</AvatarFallback>
                                    </Avatar>
                                </TableCell>
                                <TableCell>
                                    <div className="flex flex-col">
                                        <span className="font-medium">{user.username}</span>
                                        <span className="text-xs text-muted-foreground">{user.email}</span>
                                    </div>
                                </TableCell>
                                <TableCell>
                                    <Badge variant="outline" className="capitalize">{t(`users.roles.${user.role}`)}</Badge>
                                </TableCell>
                                <TableCell>
                                    <Badge
                                        variant={user.is_active ? 'default' : 'destructive'}
                                    >
                                        {user.is_active ? t('users.status.active') : t('users.status.inactive')}
                                    </Badge>
                                </TableCell>
                                <TableCell className="text-right">
                                    {renderActions(user)}
                                </TableCell>
                            </TableRow>
                        ))
                    )}
                </TableBody>
            </Table>
        </div>
    )
}
