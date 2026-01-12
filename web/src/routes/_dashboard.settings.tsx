import { createFileRoute } from '@tanstack/react-router'


import { Ban, ShieldCheck, CreditCard } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

import { CancellationReasonsCard } from "@/components/CancellationReasonsCard.tsx";
import { CategoriesCard } from "@/components/CategoriesCard.tsx";
import { PaymentMethodsCard } from "@/components/PaymentMethodsCard.tsx";



export const Route = createFileRoute("/_dashboard/settings")({
    component: SettingsPage,
})


function SettingsPage() {
    return (
        <div className="flex flex-col gap-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">Settings</h1>
                <p className="text-muted-foreground">
                    Manage application settings and configurations.
                </p>
            </div>

            <Tabs defaultValue="cancellation" className="space-y-4">
                <TabsList>
                    <TabsTrigger value="cancellation" className="flex items-center gap-2">
                        <Ban className="h-4 w-4" />
                        Cancellation Reasons
                    </TabsTrigger>

                    <TabsTrigger value="category" className="flex items-center gap-2">
                        <ShieldCheck className="h-4 w-4" />
                        Category
                    </TabsTrigger>

                    <TabsTrigger value="payment-methods" className="flex items-center gap-2">
                        <CreditCard className="h-4 w-4" />
                        Payment Methods
                    </TabsTrigger>
                </TabsList>

                <TabsContent value="cancellation">
                    <div className="grid gap-6">
                        <CancellationReasonsCard />
                    </div>
                </TabsContent>
                <TabsContent value="category">
                    <div className="grid gap-6">
                        <CategoriesCard />
                    </div>
                </TabsContent>
                <TabsContent value="payment-methods">
                    <div className="grid gap-6">
                        <PaymentMethodsCard />
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}




