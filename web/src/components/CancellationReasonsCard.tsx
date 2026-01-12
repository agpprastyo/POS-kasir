import {useSuspenseQuery} from "@tanstack/react-query";
import {cancellationReasonsListQueryOptions} from "@/lib/api/query/cancel-reason.ts";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Ban} from "lucide-react";
import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table.tsx";
import {Badge} from "@/components/ui/badge.tsx";


export function CancellationReasonsCard() {

    const {data: reasons} = useSuspenseQuery(cancellationReasonsListQueryOptions())
    const reasonsList = Array.isArray(reasons) ? reasons : (reasons as any)?.data || []

    return (
        <Card>
            <CardHeader>
                <CardTitle className="flex items-center gap-2">
                    <Ban className="h-5 w-5"/> Cancellation Reasons
                </CardTitle>
                <CardDescription>
                    Manage standard reasons for transaction cancellations.
                </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="rounded-md border">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Reason</TableHead>
                                <TableHead>Description</TableHead>
                                <TableHead className="w-[100px]">Status</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {reasonsList.length === 0 ? (
                                <TableRow>
                                    <TableCell colSpan={3} className="h-24 text-center text-muted-foreground">
                                        No cancellation reasons found.
                                    </TableCell>
                                </TableRow>
                            ) : (
                                reasonsList.map((item) => (
                                    <TableRow key={item.id}>
                                        <TableCell className="font-medium">{item.reason}</TableCell>
                                        <TableCell>{item.description || '-'}</TableCell>
                                        <TableCell>
                                            <Badge
                                                variant={item.is_active ? 'default' : 'secondary'}
                                                className={item.is_active ? 'bg-green-500 hover:bg-green-600' : ''}
                                            >
                                                {item.is_active ? 'Active' : 'Inactive'}
                                            </Badge>
                                        </TableCell>
                                    </TableRow>
                                ))
                            )}
                        </TableBody>
                    </Table>
                </div>
            </CardContent>

            {/*<CardFooter className="justify-end border-t bg-muted/20 px-6 py-4">*/}
            {/*     <Button><Plus className="mr-2 h-4 w-4"/> Add Reason</Button>*/}
            {/*</CardFooter>*/}
        </Card>
    )
}
