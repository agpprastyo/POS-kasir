import { createFileRoute } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { meQueryOptions } from '@/lib/api/query/auth'
import { ShieldCheck, User } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { UpdateAvatarCard } from "@/components/UpdateAvatarCard.tsx";
import { UpdatePasswordCard } from "@/components/UpdatePasswordCard.tsx";
import { useTranslation } from 'react-i18next'
import { z } from 'zod'

const searchSchema = z.object({
  tab: z.string().optional().default('account'),
})

export const Route = createFileRoute("/$locale/_dashboard/account")({
  validateSearch: (search) => searchSchema.parse(search),
  component: AccountPage,
})

function AccountPage() {
  const { t } = useTranslation()
  const { tab } = Route.useSearch()
  const navigate = Route.useNavigate()

  const { data: user } = useSuspenseQuery(meQueryOptions())
  const profile = (user as any)?.data ?? user

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">{t('account.title')}</h1>
        <p className="text-muted-foreground">
          {t('account.subtitle')}
        </p>
      </div>

      <Tabs
        value={tab}
        onValueChange={(val) => navigate({ search: (old) => ({ ...old, tab: val }), replace: true })}
        className="space-y-4"
      >
        <TabsList>
          <TabsTrigger value="account" className="flex items-center gap-2">
            <User className="h-4 w-4" />
            {t('account.tabs.profile')}
          </TabsTrigger>
          <TabsTrigger value="security" className="flex items-center gap-2">
            <ShieldCheck className="h-4 w-4" />
            {t('account.tabs.security')}
          </TabsTrigger>
        </TabsList>

        <TabsContent value="account">
          <div className="grid gap-6">
            <UpdateAvatarCard
              currentAvatar={profile?.avatar}
              username={profile?.username}
            />
          </div>
        </TabsContent>
        <TabsContent value="security">
          <div className="grid gap-6">
            <UpdatePasswordCard />
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
