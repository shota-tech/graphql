import type { AppProps } from 'next/app'
import { ChakraProvider, Box, Container, Heading, Flex, Spacer } from '@chakra-ui/react'
import { createClient, Provider } from 'urql'

const client = createClient({
  url: 'http://localhost:8080/query',
})

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider>
      <Provider value={client}>
        <Box minH='100vh'>
          <Box as='header' h='12' bg='teal'>
            <Flex direction='row' h='full' px='8' align='center' justify='center'>
              <Heading color='gray.50' size='md'>
                APP NAME
              </Heading>
              <Spacer />
              <Heading color='gray.50' size='md'>
                MENUE
              </Heading>
            </Flex>
          </Box>
          <Box as='main' bg='gray.50'>
            <Container maxW='container.xl'>
              <Component {...pageProps} />
            </Container>
          </Box>
        </Box>
      </Provider>
    </ChakraProvider>
  )
}
