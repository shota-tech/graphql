import { useAuth0 } from '@auth0/auth0-react'
import { Button } from '@chakra-ui/react'

export const LogoutButton = () => {
  const { logout } = useAuth0()
  return (
    <Button
      color='gray.50'
      bg='transparent'
      _hover={{ bg: 'teal.700' }}
      onClick={() => logout({ logoutParams: { returnTo: window.location.origin } })}
    >
      LOGOUT
    </Button>
  )
}
