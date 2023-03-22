import type { AppProps } from 'next/app'
import { Auth0Provider } from '@auth0/auth0-react'
import { ChakraProvider, Box, Container } from '@chakra-ui/react'
import { BASE_URL, AUTH0_DOMAIN, AUTH0_CLIENT_ID, API_URL } from '@/config/constants'
import { AuthorizedUrqlProvider, Header } from '@/components'

const App = ({ Component, pageProps }: AppProps) => {
  return (
    <Auth0Provider
      domain={AUTH0_DOMAIN}
      clientId={AUTH0_CLIENT_ID}
      authorizationParams={{
        redirect_uri: BASE_URL + '/callback',
        audience: API_URL,
      }}
    >
      <ChakraProvider>
        <AuthorizedUrqlProvider>
          <Box minH='100vh' bg='gray.50'>
            <Header />
            <Box as='main'>
              <Container maxW='container.xl'>
                <Component {...pageProps} />
              </Container>
            </Box>
          </Box>
        </AuthorizedUrqlProvider>
      </ChakraProvider>
    </Auth0Provider>
  )
}

export default App
