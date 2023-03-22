import { FC, ReactNode } from 'react'
import { useAuth0 } from '@auth0/auth0-react'
import {
  createClient,
  Provider,
  dedupExchange,
  cacheExchange,
  fetchExchange,
  Exchange,
  Operation,
} from 'urql'
import { pipe, map, mergeMap, fromPromise, fromValue } from 'wonka'
import { API_URL } from '@/config/constants'

export const AuthorizedUrqlProvider: FC<{ children: ReactNode }> = ({ children }) => {
  const { getAccessTokenSilently } = useAuth0()

  const fetchOptionsExchange =
    (fn: any): Exchange =>
    ({ forward }) =>
    (ops$) => {
      return pipe(
        ops$,
        mergeMap((operation: Operation) => {
          const result = fn(operation.context.fetchOptions)
          return pipe(
            (typeof result.then === 'function' ? fromPromise(result) : fromValue(result)) as any,
            map((fetchOptions: RequestInit | (() => RequestInit)) => ({
              ...operation,
              context: { ...operation.context, fetchOptions },
            })),
          )
        }),
        forward,
      )
    }

  const client = createClient({
    url: API_URL,
    exchanges: [
      dedupExchange,
      cacheExchange,
      fetchOptionsExchange(async (fetchOptions: any) => {
        const token = await getAccessTokenSilently({
          authorizationParams: {
            audience: API_URL,
            scope: 'read:tasks write:tasks read:user write:user',
          },
        })

        return Promise.resolve({
          ...fetchOptions,
          headers: {
            Authorization: token ? `Bearer ${token}` : '',
          },
        })
      }),
      fetchExchange,
    ],
    requestPolicy: 'network-only',
  })

  return <Provider value={client}>{children}</Provider>
}
