import fetchApi from "@/packages/fetchApi.ts";
import type {Currency} from "@/api/currency.ts";
import type {User} from "@/api/user.ts";

export enum JoinType {
  ANYONE = 'anyone',
  FRIENDS = 'friends',
  LINK = 'link'
}

export interface CreateGameInput {
  currency_id: number
  bet: number
  winning_points: number
  join_type: JoinType
}

export interface Game {
  code: string
  bet: number
  creator_id: number
  winning_points: number
  join_type: JoinType
  currency?: Currency
  creator?: User
  players: [Player]
  players_count: number
}

export interface Player {
  id: number
  username: string
  is_host: boolean
}

export const createGame = async (input: CreateGameInput): Promise<Game> => {
  const data = await fetchApi.post('/games', input)

  return data.data;
}

export const getLastCreatedGame = async (): Promise<Game | null> => {
  const {data} = await fetchApi.get('/games/current');

  return data
}

export const getGame = async (code: string): Promise<Game | null> => {
  const {data} = await fetchApi.get(`/games/${code}`);

  return data
}

export const leaveGame = async () => {
  return fetchApi.delete(`/games/leave`)
}

export interface GetGamesInput {
  page: number
  search?: string
}

export const getGames = async (params: GetGamesInput) => {
  const {data} = await fetchApi.get('/games', {
    params: params
  });

  return data;
}

export const joinGame = async (code: string) => {
  return fetchApi.post(`/games/join/${code}`)
}
