import { useAuth0 } from '@auth0/auth0-react'
import { Button } from '@chakra-ui/react'

export const LoginButton = () => {
  const { loginWithRedirect } = useAuth0()
  return (
    <Button
      color='gray.50'
      bg='transparent'
      _hover={{ bg: 'teal.700' }}
      onClick={() => loginWithRedirect()}
    >
      LOGIN
    </Button>
  )
}
