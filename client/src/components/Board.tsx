import { FC } from 'react'
import { HStack } from '@chakra-ui/react'
import { DndContext, rectIntersection, DragEndEvent } from '@dnd-kit/core'
import { Todo } from '@/graphql/generated'
import { List } from '.'

type Props = {
  onDragEnd: (e: DragEndEvent) => void
  todoTasks: Todo[]
  inProgressTasks: Todo[]
  doneTasks: Todo[]
}

export const Board: FC<Props> = ({ onDragEnd, todoTasks, inProgressTasks, doneTasks }) => {
  return (
    <DndContext collisionDetection={rectIntersection} onDragEnd={onDragEnd}>
      <HStack align='start' spacing='4'>
        <List title='TODO' tasks={todoTasks} />
        <List title='IN PROGRESS' tasks={inProgressTasks} />
        <List title='DONE' tasks={doneTasks} />
      </HStack>
    </DndContext>
  )
}
