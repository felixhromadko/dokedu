import type * as Types from '../../schema';

import gql from 'graphql-tag';
import { ChatMessageFragmentDoc } from '../../fragments/chatMessage';
import * as Urql from '@urql/vue';
export type Omit<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;
export type ChatQueryVariables = Types.Exact<{
  id: Types.Scalars['ID']['input'];
}>;


export type ChatQuery = { __typename?: 'Query', chat: { __typename?: 'Chat', id: string, name?: string | null, messages: Array<{ __typename?: 'ChatMessage', id: string, message: string, isEdited: boolean, createdAt: never, user: { __typename?: 'User', id: string, firstName: string, lastName: string, email?: string | null } }> } };


export const ChatDocument = gql`
    query chat($id: ID!) {
  chat(id: $id) {
    id
    name
    messages {
      ...chatMessage
    }
  }
}
    ${ChatMessageFragmentDoc}`;

export function useChatQuery(options: Omit<Urql.UseQueryArgs<never, ChatQueryVariables>, 'query'>) {
  return Urql.useQuery<ChatQuery, ChatQueryVariables>({ query: ChatDocument, ...options });
};