import { FC } from 'react'
import { Heading, VStack } from '@chakra-ui/react'
import { useDroppable } from '@dnd-kit/core'
import { Task } from '@/components'
import { Card } from '.'

type Props = {
  title: string
  tasks: Task[]
}

export const List: FC<Props> = ({ title, tasks }) => {
  const { setNodeRef } = useDroppable({
    id: title,
  })

  return (
    <VStack
      spacing='4'
      bg='gray.200'
      align='center'
      w='sm'
      p='4'
      rounded='md'
      shadow='md'
      ref={setNodeRef}
    >
      <Heading size='md'>{title}</Heading>
      <VStack direction='column' spacing='4' minH='60vh'>
        {tasks?.map((task) => (
          <Card key={task.id} parent={title} task={task} />
        ))}
      </VStack>
    </VStack>
  )
}
