/**
 * DevToolsPanel.tsx
 *
 * Wrapper untuk TanStack DevTools — di-lazy load hanya di development.
 * Tidak akan masuk production bundle karena __root.tsx hanya me-lazy import
 * file ini ketika import.meta.env.DEV === true.
 */
import { TanStackDevtools } from '@tanstack/react-devtools'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'

export default function DevToolsPanel() {
    return (
        <TanStackDevtools
            config={{ position: 'bottom-right' }}
            plugins={[
                {
                    name: 'Tanstack Router',
                    render: () => <TanStackRouterDevtoolsPanel />,
                },
            ]}
        />
    )
}
