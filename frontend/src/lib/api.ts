export interface Problem {
  type: string
  title: string
  status: number
  detail?: string
}

export class ApiError extends Error {
  constructor(public problem: Problem) {
    super(problem.detail || problem.title)
  }
}

export async function api<T>(path: string, options: RequestInit = {}): Promise<T> {
  const baseURL = import.meta.env.VITE_API_BASE_URL || '/api/v1'
  const response = await fetch(`${baseURL}${path}`, {
    credentials: 'include',
    ...options,
    headers: { 'Content-Type': 'application/json', ...options.headers },
  })
  if (!response.ok) {
    const problem = (await response.json().catch(() => ({
      type: 'about:blank',
      title: 'Request failed',
      status: response.status,
    }))) as Problem
    throw new ApiError(problem)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}
