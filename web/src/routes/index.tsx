import { createFileRoute, redirect } from '@tanstack/react-router'


console.log("API BASE =", import.meta.env.VITE_API_BASE)
export const Route = createFileRoute('/')({


  beforeLoad: () => {
    throw redirect({
      to: '/$locale',
      params: { locale: 'id' },
    })
  },
})
