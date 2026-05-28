import ReactDOM from 'react-dom/client'
import { RouterProvider } from '@tanstack/react-router'
import { createRouter } from './router'
import { initWebVitals } from './lib/web-vitals'

import './styles.css'

const router = createRouter()

ReactDOM.createRoot(document.getElementById('root')!).render(
    <RouterProvider router={router} />
)

// Initialize Web Vitals reporting (sends metrics to backend)
initWebVitals()
