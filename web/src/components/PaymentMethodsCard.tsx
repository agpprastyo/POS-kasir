import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import { usePaymentMethodsListQuery } from "@/lib/api/query/payment-methods"
import { Loader2, CreditCard } from "lucide-react"

export function PaymentMethodsCard() {
    const { data: paymentMethods, isLoading } = usePaymentMethodsListQuery()

    return (
        <Card>
            <CardHeader>
                <CardTitle>Payment Methods</CardTitle>
                <CardDescription>
                    Available payment methods for POS transactions.
                </CardDescription>
            </CardHeader>
            <CardContent>
                {isLoading ? (
                    <div className="flex h-40 items-center justify-center">
                        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                    </div>
                ) : (
                    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                        {paymentMethods && paymentMethods.length > 0 ? (
                            paymentMethods.map((method: any) => (
                                <div
                                    key={method.id}
                                    className="flex flex-col gap-2 rounded-lg border p-4 "
                                >
                                    <div className="flex items-center gap-2">
                                        <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
                                            <CreditCard className="h-5 w-5 text-primary" />
                                        </div>
                                        <div className="flex flex-col">
                                            <span className="font-semibold">{method.name}</span>
                                        </div>
                                    </div>
                                   
                                </div>
                            ))
                        ) : (
                            <div className="flex h-20 items-center justify-center rounded-lg border border-dashed text-muted-foreground">
                                No payment methods found.
                            </div>
                        )}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}
