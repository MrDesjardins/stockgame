import "./App.css";
import { useQuery } from "@tanstack/react-query";

interface Stock {
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
async function getStocks(): Promise<Stock[]> {
  const data = await fetch("http://localhost:8080/stocks");
  return data.json();
}

function App() {
  const { isPending, isError, data, error, refetch } = useQuery({
    queryKey: ["stocks"],
    queryFn: getStocks,
  });

  if (isError) {
    return <div>Error: {error.message}</div>;
  }

  return (
    <>
      <div>
        <div id="data">{isPending ? "Loading" : JSON.stringify(data)}</div>
        <div id="controls">
          <button onClick={() => refetch()}>Fetch</button>
        </div>
      </div>
    </>
  );
}

export default App;
