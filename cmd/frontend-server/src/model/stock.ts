export interface Stock {
  id: number;
  symbol: string;
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  adj_close: number;
  volume: number;
}

export interface SolutionDayPrice {
  day: number;
  price: number;
}
export interface Solution {
  afterDate: string;
  symbol: string;
  dayPrice: SolutionDayPrice[];
}