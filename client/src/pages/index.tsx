import { ChangeEvent, useState, useEffect } from 'react'
import { withAuthenticationRequired } from '@auth0/auth0-react'
import { IconButton, Input, HStack, VStack } from '@chakra-ui/react'
import { AddIcon } from '@chakra-ui/icons'
import { DragEndEvent } from '@dnd-kit/core'
import {
  Task,
  Status,
  useFetchTasksQuery,
  useCreateTaskMutation,
  useUpdateTaskMutation,
} from '@/graphql/generated'
import { Board } from '@/components'

const Home = () => {
  const [text, setText] = useState('')
  const [todoTasks, setTodoTasks] = useState<Task[]>([])
  const [inProgressTasks, setInProgressTasks] = useState<Task[]>([])
  const [doneTasks, setDoneTasks] = useState<Task[]>([])
  const [fetchTasksResult] = useFetchTasksQuery()
  const [, createTask] = useCreateTaskMutation()
  const [, updateTask] = useUpdateTaskMutation()
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
    createTask({ text: text }).then((result) => {
      if (result.error) {
        console.error('Oh no!', result.error)
      }
    })
    setText('')
  }

  const handleDragEnd = (e: DragEndEvent) => {
    const newParent = e.over?.id ?? ''
    const oldParent = e.active.data.current?.parent ?? ''
    const task = e.active.data.current?.task ?? null
    if (newParent === oldParent) {
      return
    }
    if (oldParent == 'TODO') {
      setTodoTasks(todoTasks.filter((value) => value.id !== task.id))
    } else if (oldParent === 'IN PROGRESS') {
      setInProgressTasks(inProgressTasks.filter((value) => value.id !== task.id))
    } else if (oldParent === 'DONE') {
      setDoneTasks(doneTasks.filter((value) => value.id !== task.id))
    }
    let newStatus = Status.Todo
    if (newParent === 'TODO') {
      setTodoTasks([...todoTasks, task])
    } else if (newParent === 'IN PROGRESS') {
      newStatus = Status.InProgress
      setInProgressTasks([...inProgressTasks, task])
    } else if (newParent === 'DONE') {
      newStatus = Status.Done
      setDoneTasks([...doneTasks, task])
    }
    updateTask({ id: task.id, status: newStatus }).then((result) => {
      if (result.error) {
        console.error('Oh no!', result.error)
      }
    })
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

export default withAuthenticationRequired(Home)
