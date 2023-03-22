import { FC } from 'react'
import { HStack } from '@chakra-ui/react'
import { DndContext, rectIntersection, DragEndEvent } from '@dnd-kit/core'
import { Task } from '@/components'
import { List } from '.'

type Props = {
  onDragEnd: (e: DragEndEvent) => void
  todoTasks: Task[]
  inProgressTasks: Task[]
  doneTasks: Task[]
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
