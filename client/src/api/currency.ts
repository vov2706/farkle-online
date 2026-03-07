import fetchApi from "@/packages/fetchApi.ts";

export interface Currency {
  id: number
  slug: string
  name: string
}

export enum CurrencyType {
  Bronze = "bronze",
  Silver = "silver",
  Gold = "gold"
}

export const getCurrencies = async (): Promise<Currency[]> => {
  const {data} = await fetchApi.get('/currencies');

  return data.data;
}
