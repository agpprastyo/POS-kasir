import { useTranslation } from 'react-i18next';
import { useRouter, useLocation } from '@tanstack/react-router';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Globe } from 'lucide-react';

export function LanguageSwitcher() {
    const { i18n } = useTranslation();
    const router = useRouter();
    const location = useLocation();

    const handleLanguageChange = (value: string) => {
        const currentPath = location.pathname;
        const segments = currentPath.split('/');

        if (segments.length >= 2) {
            segments[1] = value;
        } else {
            segments.push(value);
        }
        const newPath = segments.join('/');
        router.navigate({
            to: newPath,
            search: (old) => old,
            replace: true
        });
    };

    return (
        <Select value={i18n.language} onValueChange={handleLanguageChange}>
            <SelectTrigger className="w-full gap-2">
                <Globe className="h-4 w-4 text-muted-foreground" />
                <SelectValue placeholder="Language" />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="id">Bahasa (ID)</SelectItem>
                <SelectItem value="en">English (EN)</SelectItem>
            </SelectContent>
        </Select>
    );
}
