import type { AppProps } from 'next/app'
import { ChakraProvider, Box, Container, Heading, Flex, Spacer } from '@chakra-ui/react'

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ChakraProvider>
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
        <Container maxW='container.xl' minH='100vh'>
          <Component {...pageProps} />
        </Container>
      </Box>
    </ChakraProvider>
  )
}
