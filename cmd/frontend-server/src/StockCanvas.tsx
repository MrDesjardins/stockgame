import { Stock } from "./model/stock";

export interface StockCanvasProps {
  data: Stock[];
}
export function StockCanvas(props: StockCanvasProps) {
  return (
    <div>
      <h1>Stocks</h1>
      <table>
        <thead>
          <tr>
            <th>Symbol</th>
            <th>Date</th>
            <th>Open</th>
            <th>High</th>
            <th>Low</th>
            <th>Close</th>
            <th>Adj Close</th>
            <th>Volume</th>
          </tr>
        </thead>
        <tbody>
          {props.data.map((stock) => (
            <tr key={stock.id}>
              <td>{stock.symbol}</td>
              <td>{stock.date}</td>
              <td>{stock.open}</td>
              <td>{stock.high}</td>
              <td>{stock.low}</td>
              <td>{stock.close}</td>
              <td>{stock.adj_close}</td>
              <td>{stock.volume}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
