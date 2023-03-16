import type { AppProps } from 'next/app'
import { Auth0Provider } from '@auth0/auth0-react'
import { ChakraProvider, Box, Container } from '@chakra-ui/react'
import { createClient, Provider } from 'urql'
import { BASE_URL, AUTH0_DOMAIN, AUTH0_CLIENT_ID } from '@/config/constants'
import { Header } from '@/components'

const client = createClient({
  url: 'http://localhost:8080/graphql',
})

const App = ({ Component, pageProps }: AppProps) => {
  return (
    <Auth0Provider
      domain={AUTH0_DOMAIN}
      clientId={AUTH0_CLIENT_ID}
      authorizationParams={{
        redirect_uri: BASE_URL,
      }}
    >
      <ChakraProvider>
        <Provider value={client}>
          <Box minH='100vh' bg='gray.50'>
            <Header />
            <Box as='main'>
              <Container maxW='container.xl'>
                <Component {...pageProps} />
              </Container>
            </Box>
          </Box>
        </Provider>
      </ChakraProvider>
    </Auth0Provider>
  )
}

export default App
