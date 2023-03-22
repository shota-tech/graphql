import { useEffect } from 'react'
import { useRouter } from 'next/router'
import { useAuth0 } from '@auth0/auth0-react'
import { Center, Spinner } from '@chakra-ui/react'
import { useFetchUserQuery, useCreateUserMutation } from '@/graphql/generated'

const Callback = () => {
  const router = useRouter()
  const { user } = useAuth0()
  const [fetchUserResult] = useFetchUserQuery({ pause: !user })
  const [, createUser] = useCreateUserMutation()
  const { data } = fetchUserResult

  useEffect(() => {
    if (!user || !data) return
    if (!data.fetchUser) {
      createUser({ name: user.name ?? '' }).then((result) => {
        if (result.error) {
          console.error('failed to update user: ', result.error)
        }
      })
    }
    router.push('/')
  }, [user, data])

  return (
    <Center h='xl'>
      <Spinner color='teal.700' size='xl' />
    </Center>
  )
}

export default Callback
