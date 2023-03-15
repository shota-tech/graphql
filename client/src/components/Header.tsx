import { useAuth0 } from '@auth0/auth0-react'
import { Box, Heading, Flex, Spacer } from '@chakra-ui/react'
import { LoginButton, LogoutButton } from '@/components'

export const Header = () => {
  const { isAuthenticated } = useAuth0()

  return (
    <Box as='header' h='12' bg='teal'>
      <Flex direction='row' h='full' px='8' align='center' justify='center'>
        <Heading color='gray.50' size='md'>
          KANBAN BOARD
        </Heading>
        <Spacer />
        {isAuthenticated ? <LogoutButton /> : <LoginButton />}
      </Flex>
    </Box>
  )
}
