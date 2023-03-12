import { FC } from 'react'
import { Text, VStack } from '@chakra-ui/react'
import { useDraggable } from '@dnd-kit/core'
import { CSS } from '@dnd-kit/utilities'
import { Task } from '@/graphql/generated'

type Props = {
  parent: string
  task: Task
}

export const Card: FC<Props> = ({ parent, task }) => {
  const { attributes, listeners, setNodeRef, transform } = useDraggable({
    id: task.id,
    data: {
      parent: parent,
      task: task,
    },
  })

  return (
    <VStack
      spacing='2'
      bg='gray.50'
      w='xs'
      p='4'
      rounded='md'
      shadow='md'
      ref={setNodeRef}
      style={{ transform: CSS.Translate.toString(transform) }}
      {...listeners}
      {...attributes}
    >
      <Text>{task.text}</Text>
    </VStack>
  )
}
