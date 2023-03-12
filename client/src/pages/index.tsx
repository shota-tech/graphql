import { ChangeEvent, useState, useEffect } from 'react'
import { IconButton, Input, HStack, VStack } from '@chakra-ui/react'
import { AddIcon } from '@chakra-ui/icons'
import { DragEndEvent } from '@dnd-kit/core'
import { Task, Status, useFetchTasksQuery, useCreateTaskMutation } from '@/graphql/generated'
import { Board } from '@/components'

const userId = 'cg1ltn51nm6u7l352ma0'

export default function Home() {
  const [text, setText] = useState('')
  const [todoTasks, setTodoTasks] = useState<Task[]>([])
  const [inProgressTasks, setInProgressTasks] = useState<Task[]>([])
  const [doneTasks, setDoneTasks] = useState<Task[]>([])
  const [fetchTasksResult] = useFetchTasksQuery()
  const [_, createTask] = useCreateTaskMutation()
  const { data, fetching, error } = fetchTasksResult

  useEffect(() => {
    if (!data) {
      return
    }
    setTodoTasks(data.fetchTasks.filter((task) => task.status === Status.Todo))
    setInProgressTasks(data.fetchTasks.filter((task) => task.status === Status.InProgress))
    setDoneTasks(data.fetchTasks.filter((task) => task.status === Status.Done))
  }, [data])

  if (fetching) return <p>Loading...</p>
  if (error) return <p>Oh no... {error.message}</p>

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setText(e.target.value)
  }

  const handleClick = () => {
    if (!text) {
      return
    }
    createTask({ text: text, userId: userId }).then((result) => {
      if (result.error) {
        console.error('Oh no!', result.error)
      }
    })
    setText('')
  }

  const handleDragEnd = (e: DragEndEvent) => {
    const container = e.over?.id
    const parent = e.active.data.current?.parent ?? ''
    const task = e.active.data.current?.task ?? null
    if (container === parent) {
      return
    }
    if (parent == 'TODO') {
      setTodoTasks(todoTasks.filter((value) => value.id !== task.id))
    } else if (parent === 'IN PROGRESS') {
      setInProgressTasks(inProgressTasks.filter((value) => value.id !== task.id))
    } else if (parent === 'DONE') {
      setDoneTasks(doneTasks.filter((value) => value.id !== task.id))
    }
    if (container === 'TODO') {
      setTodoTasks([...todoTasks, task])
    } else if (container === 'IN PROGRESS') {
      setInProgressTasks([...inProgressTasks, task])
    } else if (container === 'DONE') {
      setDoneTasks([...doneTasks, task])
    }
  }

  return (
    <VStack align='center' spacing='8' pt='8'>
      <HStack bg='gray.200' w='sm' p='4' rounded='md' shadow='md'>
        <Input bg='gray.50' value={text} onChange={handleChange} />
        <IconButton
          colorScheme='teal'
          aria-label='add task'
          onClick={handleClick}
          icon={<AddIcon />}
        />
      </HStack>
      <Board
        onDragEnd={handleDragEnd}
        todoTasks={todoTasks}
        inProgressTasks={inProgressTasks}
        doneTasks={doneTasks}
      />
    </VStack>
  )
}
