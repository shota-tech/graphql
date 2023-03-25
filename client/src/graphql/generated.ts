import gql from 'graphql-tag';
import * as Urql from 'urql';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: string;
  String: string;
  Boolean: boolean;
  Int: number;
  Float: number;
};

export type CreateTaskInput = {
  text: Scalars['String'];
};

export type CreateTodoInput = {
  taskID: Scalars['String'];
  text: Scalars['String'];
};

export type CreateUserInput = {
  name: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createTask: Task;
  createTodo: Todo;
  createUser: User;
  updateTask: Task;
  updateTodo: Todo;
};


export type MutationCreateTaskArgs = {
  input: CreateTaskInput;
};


export type MutationCreateTodoArgs = {
  input: CreateTodoInput;
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


export type MutationUpdateTaskArgs = {
  input: UpdateTaskInput;
};


export type MutationUpdateTodoArgs = {
  input: UpdateTodoInput;
};

export type Query = {
  __typename?: 'Query';
  fetchTasks: Array<Task>;
  fetchUser?: Maybe<User>;
};

export enum Status {
  Done = 'DONE',
  InProgress = 'IN_PROGRESS',
  Todo = 'TODO'
}

export type Task = {
  __typename?: 'Task';
  id: Scalars['ID'];
  status: Status;
  text: Scalars['String'];
  todos: Array<Todo>;
  user: User;
};

export type Todo = {
  __typename?: 'Todo';
  done: Scalars['Boolean'];
  id: Scalars['ID'];
  task: Task;
  text: Scalars['String'];
};

export type UpdateTaskInput = {
  id: Scalars['ID'];
  status?: InputMaybe<Status>;
  text?: InputMaybe<Scalars['String']>;
};

export type UpdateTodoInput = {
  done?: InputMaybe<Scalars['Boolean']>;
  id: Scalars['ID'];
  text?: InputMaybe<Scalars['String']>;
};

export type User = {
  __typename?: 'User';
  id: Scalars['ID'];
  name: Scalars['String'];
  tasks: Array<Task>;
};

export type CreateTaskMutationVariables = Exact<{
  text: Scalars['String'];
}>;


export type CreateTaskMutation = { __typename?: 'Mutation', createTask: { __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string } } };

export type CreateUserMutationVariables = Exact<{
  name: Scalars['String'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: { __typename?: 'User', id: string, name: string } };

export type FetchTasksQueryVariables = Exact<{ [key: string]: never; }>;


export type FetchTasksQuery = { __typename?: 'Query', fetchTasks: Array<{ __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string, name: string } }> };

export type FetchUserQueryVariables = Exact<{ [key: string]: never; }>;


export type FetchUserQuery = { __typename?: 'Query', fetchUser?: { __typename?: 'User', id: string, name: string } | null };

export type UpdateTaskMutationVariables = Exact<{
  id: Scalars['ID'];
  text?: InputMaybe<Scalars['String']>;
  status?: InputMaybe<Status>;
}>;


export type UpdateTaskMutation = { __typename?: 'Mutation', updateTask: { __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string } } };


export const CreateTaskDocument = gql`
    mutation createTask($text: String!) {
  createTask(input: {text: $text}) {
    id
    text
    status
    user {
      id
    }
  }
}
    `;

export function useCreateTaskMutation() {
  return Urql.useMutation<CreateTaskMutation, CreateTaskMutationVariables>(CreateTaskDocument);
};
export const CreateUserDocument = gql`
    mutation createUser($name: String!) {
  createUser(input: {name: $name}) {
    id
    name
  }
}
    `;

export function useCreateUserMutation() {
  return Urql.useMutation<CreateUserMutation, CreateUserMutationVariables>(CreateUserDocument);
};
export const FetchTasksDocument = gql`
    query fetchTasks {
  fetchTasks {
    id
    text
    status
    user {
      id
      name
    }
  }
}
    `;

export function useFetchTasksQuery(options?: Omit<Urql.UseQueryArgs<FetchTasksQueryVariables>, 'query'>) {
  return Urql.useQuery<FetchTasksQuery, FetchTasksQueryVariables>({ query: FetchTasksDocument, ...options });
};
export const FetchUserDocument = gql`
    query fetchUser {
  fetchUser {
    id
    name
  }
}
    `;

export function useFetchUserQuery(options?: Omit<Urql.UseQueryArgs<FetchUserQueryVariables>, 'query'>) {
  return Urql.useQuery<FetchUserQuery, FetchUserQueryVariables>({ query: FetchUserDocument, ...options });
};
export const UpdateTaskDocument = gql`
    mutation updateTask($id: ID!, $text: String, $status: Status) {
  updateTask(input: {id: $id, text: $text, status: $status}) {
    id
    text
    status
    user {
      id
    }
  }
}
    `;

export function useUpdateTaskMutation() {
  return Urql.useMutation<UpdateTaskMutation, UpdateTaskMutationVariables>(UpdateTaskDocument);
};