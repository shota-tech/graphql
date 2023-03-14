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
  userId: Scalars['ID'];
};

export type CreateUserInput = {
  name: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createTask: Task;
  createUser: User;
  updateTask: Task;
};


export type MutationCreateTaskArgs = {
  input: CreateTaskInput;
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};


export type MutationUpdateTaskArgs = {
  input?: InputMaybe<UpdateTaskInput>;
};

export type Query = {
  __typename?: 'Query';
  fetchTasks: Array<Task>;
  fetchUser?: Maybe<User>;
};


export type QueryFetchTasksArgs = {
  userId: Scalars['ID'];
};


export type QueryFetchUserArgs = {
  id: Scalars['ID'];
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
  user: User;
};

export type UpdateTaskInput = {
  id: Scalars['ID'];
  status?: InputMaybe<Status>;
  text?: InputMaybe<Scalars['String']>;
};

export type User = {
  __typename?: 'User';
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type CreateTaskMutationVariables = Exact<{
  text: Scalars['String'];
  userId: Scalars['ID'];
}>;


export type CreateTaskMutation = { __typename?: 'Mutation', createTask: { __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string } } };

export type CreateUserMutationVariables = Exact<{
  name: Scalars['String'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: { __typename?: 'User', id: string, name: string } };

export type FetchTasksQueryVariables = Exact<{
  userId: Scalars['ID'];
}>;


export type FetchTasksQuery = { __typename?: 'Query', fetchTasks: Array<{ __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string, name: string } }> };

export type FetchUserQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type FetchUserQuery = { __typename?: 'Query', fetchUser?: { __typename?: 'User', id: string, name: string } | null };

export type UpdateTaskMutationVariables = Exact<{
  id: Scalars['ID'];
  text?: InputMaybe<Scalars['String']>;
  status?: InputMaybe<Status>;
}>;


export type UpdateTaskMutation = { __typename?: 'Mutation', updateTask: { __typename?: 'Task', id: string, text: string, status: Status, user: { __typename?: 'User', id: string } } };


export const CreateTaskDocument = gql`
    mutation createTask($text: String!, $userId: ID!) {
  createTask(input: {text: $text, userId: $userId}) {
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
    query fetchTasks($userId: ID!) {
  fetchTasks(userId: $userId) {
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

export function useFetchTasksQuery(options: Omit<Urql.UseQueryArgs<FetchTasksQueryVariables>, 'query'>) {
  return Urql.useQuery<FetchTasksQuery, FetchTasksQueryVariables>({ query: FetchTasksDocument, ...options });
};
export const FetchUserDocument = gql`
    query fetchUser($id: ID!) {
  fetchUser(id: $id) {
    id
    name
  }
}
    `;

export function useFetchUserQuery(options: Omit<Urql.UseQueryArgs<FetchUserQueryVariables>, 'query'>) {
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