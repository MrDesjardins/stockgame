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
export interface SolutionRequest {
  afterDate: string;
  symbol: string;
  estimatedDayPrices: SolutionDayPrice[];
}

export interface BB20Payload{
  date: string;
  upperBand: number;
  lowerBand: number;
}
export interface SolutionResponse {
  symbol: string;
  stocks: Stock[];
  score: number;
  bb20: Record<string, BB20Payload>
}
