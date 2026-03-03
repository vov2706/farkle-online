import fetchApi from '../packages/fetchApi.ts';
import type {Currency} from "@/api/common.ts";
import type {Game} from "@/api/game.ts";

export interface Balance {
  amount: number
  currency: Currency
}

export interface User {
  id: number
  username: string
  balance: Balance
  current_game?: Game | null
}

export const getProfile = async () => {
  const { data } = await fetchApi.get('/profile');

  return data.data as User;
}
