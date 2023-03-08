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

export type CreateTodoInput = {
  text: Scalars['String'];
  userId: Scalars['String'];
};

export type CreateUserInput = {
  name: Scalars['String'];
};

export type Mutation = {
  __typename?: 'Mutation';
  createTodo: Todo;
  createUser: User;
};


export type MutationCreateTodoArgs = {
  input: CreateTodoInput;
};


export type MutationCreateUserArgs = {
  input: CreateUserInput;
};

export type Query = {
  __typename?: 'Query';
  todos: Array<Todo>;
  user?: Maybe<User>;
};


export type QueryUserArgs = {
  id: Scalars['ID'];
};

export type Todo = {
  __typename?: 'Todo';
  done: Scalars['Boolean'];
  id: Scalars['ID'];
  text: Scalars['String'];
  user: User;
};

export type User = {
  __typename?: 'User';
  id: Scalars['ID'];
  name: Scalars['String'];
};

export type CreateTodoMutationVariables = Exact<{
  text: Scalars['String'];
  userId: Scalars['String'];
}>;


export type CreateTodoMutation = { __typename?: 'Mutation', createTodo: { __typename?: 'Todo', id: string, text: string, done: boolean, user: { __typename?: 'User', id: string } } };

export type CreateUserMutationVariables = Exact<{
  name: Scalars['String'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: { __typename?: 'User', id: string, name: string } };

export type FetchTodosQueryVariables = Exact<{ [key: string]: never; }>;


export type FetchTodosQuery = { __typename?: 'Query', todos: Array<{ __typename?: 'Todo', id: string, text: string, done: boolean, user: { __typename?: 'User', id: string, name: string } }> };

export type FetchUserQueryVariables = Exact<{
  id: Scalars['ID'];
}>;


export type FetchUserQuery = { __typename?: 'Query', user?: { __typename?: 'User', id: string, name: string } | null };


export const CreateTodoDocument = gql`
    mutation createTodo($text: String!, $userId: String!) {
  createTodo(input: {text: $text, userId: $userId}) {
    id
    text
    done
    user {
      id
    }
  }
}
    `;

export function useCreateTodoMutation() {
  return Urql.useMutation<CreateTodoMutation, CreateTodoMutationVariables>(CreateTodoDocument);
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
export const FetchTodosDocument = gql`
    query fetchTodos {
  todos {
    id
    text
    done
    user {
      id
      name
    }
  }
}
    `;

export function useFetchTodosQuery(options?: Omit<Urql.UseQueryArgs<FetchTodosQueryVariables>, 'query'>) {
  return Urql.useQuery<FetchTodosQuery, FetchTodosQueryVariables>({ query: FetchTodosDocument, ...options });
};
export const FetchUserDocument = gql`
    query fetchUser($id: ID!) {
  user(id: $id) {
    id
    name
  }
}
    `;

export function useFetchUserQuery(options: Omit<Urql.UseQueryArgs<FetchUserQueryVariables>, 'query'>) {
  return Urql.useQuery<FetchUserQuery, FetchUserQueryVariables>({ query: FetchUserDocument, ...options });
};