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

export interface StockPublic {
  date: string;
  open: number;
  high: number;
  low: number;
  close: number;
  adj_close: number;
  volume: number;
  symbol_uuid: string;
}

export interface SolutionDayPrice {
  day: number;
  price: number;
}
export interface SolutionRequest {
  afterDate: string;
  symbolUUID: string;
  estimatedDayPrices: SolutionDayPrice[];
}

export interface BB20Payload {
  date: string;
  upperBand: number;
  lowerBand: number;
}

export interface SolutionScore {
  total: number;
  inLowHigh: number;
  inOpenClose: number;
  inBollinger: number;
  inDirection: number;
}
export interface SolutionResponse {
  symbol: string;
  name: string;
  stocks: StockPublic[];
  score: SolutionScore;
  bb20: Record<string, BB20Payload>;
}
