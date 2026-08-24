import { api } from '@/lib/api'

let required: boolean | undefined

export async function setupRequired(): Promise<boolean> {
  if (required === undefined) {
    required = (await api<{ setupRequired: boolean }>('/setup/status')).setupRequired
  }
  return required
}

export function markInstalled() {
  required = false
}
