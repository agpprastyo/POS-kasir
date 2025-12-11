export const API_BASE = (import.meta.env.VITE_API_BASE as string) ?? ''

export function buildUrl(path: string) {
    const base = API_BASE.replace(/\/$/, '')
    return base ? `${base}${path.startsWith('/') ? path : `/${path}`}` : path
}
