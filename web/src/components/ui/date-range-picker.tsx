import * as React from "react"
import { format, subDays, startOfDay, endOfDay } from "date-fns"
import { id as localeId, enUS as localeEn } from "date-fns/locale"
import { Calendar as CalendarIcon, Check } from "lucide-react"
import { DateRange } from "react-day-picker"

import { cn, formatDate } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useTranslation } from "react-i18next"

interface DateRangePickerProps extends React.HTMLAttributes<HTMLDivElement> {
  date: { from: string; to: string }
  onDateChange: (date: { from: string; to: string }) => void
}

export function DateRangePicker({
  className,
  date,
  onDateChange,
}: DateRangePickerProps) {
  const { t, i18n } = useTranslation()
  const currentLocale = i18n.language === 'id' ? localeId : localeEn

  const [isOpen, setIsOpen] = React.useState(false)

  const [internalDate, setInternalDate] = React.useState<DateRange | undefined>(
    {
      from: date.from ? new Date(date.from) : undefined,
      to: date.to ? new Date(date.to) : undefined,
    }
  )

  React.useEffect(() => {
    if (!isOpen) {
      setInternalDate({
        from: date.from ? new Date(date.from) : undefined,
        to: date.to ? new Date(date.to) : undefined,
      })
    }
  }, [date, isOpen])

  const formatDateString = (d?: Date) => {
    if (!d) return ""
    // Format to YYYY-MM-DD for consistency
    return format(d, "yyyy-MM-dd")
  }

  const parseDateString = (d?: Date) => {
    if (!d) return ""
    return formatDate(d)
  }

  const handleApply = (newRange?: DateRange) => {
    if (newRange?.from && newRange?.to) {
      onDateChange({
        from: formatDateString(newRange.from),
        to: formatDateString(newRange.to),
      })
      setIsOpen(false)
    }
  }

  const presets = [
    {
      label: t("reports.date_ranges.today", "Today"),
      getValue: () => {
        const today = new Date()
        return { from: startOfDay(today), to: endOfDay(today) }
      },
    },
    {
      label: t("reports.date_ranges.yesterday", "Yesterday"),
      getValue: () => {
        const yesterday = subDays(new Date(), 1)
        return { from: startOfDay(yesterday), to: endOfDay(yesterday) }
      },
    },
    {
      label: t("reports.date_ranges.last_7_days", "Last 7 Days"),
      getValue: () => {
        const today = new Date()
        return { from: startOfDay(subDays(today, 6)), to: endOfDay(today) }
      },
    },
    {
      label: t("reports.date_ranges.last_30_days", "Last 30 Days"),
      getValue: () => {
        const today = new Date()
        return { from: startOfDay(subDays(today, 29)), to: endOfDay(today) }
      },
    },
  ]

  return (
    <div className={cn("grid gap-2", className)}>
      <Popover open={isOpen} onOpenChange={setIsOpen}>
        <PopoverTrigger asChild>
          <Button
            id="date"
            variant={"outline"}
            className={cn(
              "w-[260px] justify-start text-left font-normal",
              !date && "text-muted-foreground"
            )}
          >
            <CalendarIcon className="mr-2 h-4 w-4" />
            {date?.from ? (
              date.to ? (
                <>
                  {parseDateString(new Date(date.from))} -{" "}
                  {parseDateString(new Date(date.to))}
                </>
              ) : (
                parseDateString(new Date(date.from))
              )
            ) : (
              <span>{t("common.pick_date_range", "Pick a date range")}</span>
            )}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto p-0" align="end">
          <div className="flex flex-col md:flex-row">
            <div className="flex flex-col gap-1 border-r p-3 pr-4">
              <span className="mb-2 text-xs font-medium text-muted-foreground uppercase px-2">
                {t("reports.date_ranges.quick_select", "Quick Select")}
              </span>
              {presets.map((preset) => (
                <Button
                  key={preset.label}
                  variant="ghost"
                  size="sm"
                  className="justify-start text-sm font-normal"
                  onClick={() => {
                    const value = preset.getValue()
                    setInternalDate(value)
                    handleApply(value)
                  }}
                >
                  {preset.label}
                </Button>
              ))}
            </div>
            <div className="p-3">
              <Calendar
                initialFocus
                mode="range"
                defaultMonth={internalDate?.from}
                selected={internalDate}
                onSelect={(range) => setInternalDate(range)}
                numberOfMonths={2}
                locale={currentLocale}
              />
              <div className="mt-4 flex justify-end gap-2 border-t pt-4">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => {
                    setIsOpen(false)
                    setInternalDate({
                      from: date.from ? new Date(date.from) : undefined,
                      to: date.to ? new Date(date.to) : undefined,
                    })
                  }}
                >
                  {t("common.cancel", "Cancel")}
                </Button>
                <Button
                  size="sm"
                  onClick={() => handleApply(internalDate)}
                  disabled={!internalDate?.from || !internalDate?.to}
                >
                  <Check className="mr-2 h-4 w-4" />
                  {t("common.apply", "Apply")}
                </Button>
              </div>
            </div>
          </div>
        </PopoverContent>
      </Popover>
    </div>
  )
}
