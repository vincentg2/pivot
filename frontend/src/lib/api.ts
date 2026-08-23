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
  const response = await fetch(`/api/v1${path}`, {
    credentials: 'include',
    ...options,
    headers: { 'Content-Type': 'application/json', ...options.headers },
  })
  if (!response.ok) {
    const problem = (await response
      .json()
      .catch(() => ({
        type: 'about:blank',
        title: 'Request failed',
        status: response.status,
      }))) as Problem
    throw new ApiError(problem)
  }
  if (response.status === 204) return undefined as T
  return response.json() as Promise<T>
}
