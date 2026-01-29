import { LanguageSwitcher } from '@/components/LanguageSwitcher'
import { ThemeToggle } from '@/components/ThemeToggle'

export function SettingsPanel() {
  return (
    <div className="flex items-center gap-3">
      <LanguageSwitcher />
      <ThemeToggle />
    </div>
  )
}
