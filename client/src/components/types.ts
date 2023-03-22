import { Status } from '@/graphql/generated'

export type Task = {
  id: string
  text: string
  status: Status
  user: {
    id: string
    name: string
  }
}
