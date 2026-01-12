import { createFileRoute } from '@tanstack/react-router'
import { useSuspenseQuery } from '@tanstack/react-query'
import { meQueryOptions } from '@/lib/api/query/auth'

import { ShieldCheck, User } from 'lucide-react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"

import { UpdateAvatarCard } from "@/components/UpdateAvatarCard.tsx";
import { UpdatePasswordCard } from "@/components/UpdatePasswordCard.tsx";



export const Route = createFileRoute("/_dashboard/account")({
  component: AccountPage,
})

function AccountPage() {
  const { data: user } = useSuspenseQuery(meQueryOptions())
  const profile = (user as any)?.data ?? user

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Account</h1>
        <p className="text-muted-foreground">
          Manage your account settings and preferences.
        </p>
      </div>

      <Tabs defaultValue="account" className="space-y-4">
        <TabsList>
          <TabsTrigger value="account" className="flex items-center gap-2">
            <User className="h-4 w-4" />
            Profile
          </TabsTrigger>
          <TabsTrigger value="security" className="flex items-center gap-2">
            <ShieldCheck className="h-4 w-4" />
            Security
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
