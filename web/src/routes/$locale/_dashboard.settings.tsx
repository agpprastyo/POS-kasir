import { createFileRoute } from '@tanstack/react-router'
import { Ban, ShieldCheck, CreditCard, Palette, Printer } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { CancellationReasonsCard } from "@/components/CancellationReasonsCard.tsx";
import { CategoriesCard } from "@/components/CategoriesCard.tsx";
import { PaymentMethodsCard } from "@/components/PaymentMethodsCard.tsx";
import { BrandingSettingsCard } from "@/components/BrandingSettingsCard.tsx";
import { PrinterSettingsCard } from "@/components/PrinterSettingsCard.tsx";
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

const searchSchema = z.object({
    tab: z.string().optional().default('cancellation'),
})

export const Route = createFileRoute("/$locale/_dashboard/settings")({
    validateSearch: (search) => searchSchema.parse(search),
    component: SettingsPage,
})

function SettingsPage() {
    const { t } = useTranslation()
    const { tab } = Route.useSearch()
    const navigate = Route.useNavigate()

    return (
        <div className="flex flex-col gap-6">
            <div>
                <h1 className="text-3xl font-bold tracking-tight">{t('settings.title')}</h1>
                <p className="text-muted-foreground">
                    {t('settings.description')}
                </p>
            </div>

            <Tabs
                value={tab}
                onValueChange={(val) => navigate({ search: (old) => ({ ...old, tab: val }), replace: true })}
                className="space-y-4"
            >
                <TabsList>
                    <TabsTrigger value="cancellation" className="flex items-center gap-2">
                        <Ban className="h-4 w-4" />
                        {t('settings.tabs.cancellation')}
                    </TabsTrigger>

                    <TabsTrigger value="category" className="flex items-center gap-2">
                        <ShieldCheck className="h-4 w-4" />
                        {t('settings.tabs.category')}
                    </TabsTrigger>

                    <TabsTrigger value="payment-methods" className="flex items-center gap-2">
                        <CreditCard className="h-4 w-4" />
                        {t('settings.tabs.payment_methods')}
                    </TabsTrigger>

                    <TabsTrigger value="branding" className="flex items-center gap-2">
                        <Palette className="h-4 w-4" />
                        {t('settings.tabs.branding')}
                    </TabsTrigger>

                    <TabsTrigger value="printer" className="flex items-center gap-2">
                        <Printer className="h-4 w-4" />
                        {t('settings.tabs.printer')}
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
                <TabsContent value="branding">
                    <div className="grid gap-6">
                        <BrandingSettingsCard />
                    </div>
                </TabsContent>
                <TabsContent value="printer">
                    <div className="grid gap-6">
                        <PrinterSettingsCard />
                    </div>
                </TabsContent>
            </Tabs>
        </div >
    )
}




