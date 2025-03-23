import "./App.css";
import { useQuery } from "@tanstack/react-query";
import { StockCanvas } from "./StockCanvas";
import { Stock } from "./model/stock";
import { useCallback, useMemo, useState } from "react";
import { xPixelToDay, yPixelToPrice } from "./logic/canvasLogic";

async function getStocks(): Promise<Stock[]> {
  const data = await fetch("http://localhost:8080/stocks");
  return data.json();
}

function App() {
  const futureDays = 10;
  const width = 800;
  const height = 400;
  const { isPending, isError, data, error, refetch } = useQuery({
    queryKey: ["stocks"],
    queryFn: getStocks,
  });

  // Find the min and max to determine the vertical space needed
  const minPrice = useMemo(
    () => (data == undefined ? 0 : Math.min(...data.map((s) => s.low)) * 0.95),
    [data]
  );
  const maxPrice = useMemo(
    () => (data == undefined ? 0 : Math.max(...data.map((s) => s.high)) * 1.05),
    [data]
  );
  const totalDays = useMemo(() => (data?.length ?? 0) + futureDays, [data]);

  const [userDrawnPoints, setUserDrawnPoints] = useState<
    { x: number; y: number }[]
  >([]);

  const showSolution = useCallback(() => {
    console.log("Solution");
  }, []);

  const clear = () => {
    setUserDrawnPoints([]);
  };

  if (isError) {
    return <div>Error: {error.message}</div>;
  }

  const getUserDrawnPrices = () => {
    if (data == undefined) {
      return [];
    }
    let lastDay = -1;
    const onPricePerDay: { day: number; price: number }[] = [];
    let sumDay = 0;
    let countDay = 0;
    for (const dayPrice of userDrawnPoints) {
      const day = xPixelToDay(dayPrice.x, width, data.length + futureDays);
      const price = yPixelToPrice(dayPrice.y, height, minPrice, maxPrice);
      if (lastDay === -1 || lastDay === day) {
        lastDay = day;
        sumDay += price;
        countDay += 1;
      } else {
        onPricePerDay.push({ day: lastDay, price: sumDay / countDay });
        lastDay = day;
        sumDay = price;
        countDay = 1;
      }
    }
    onPricePerDay.push({ day: lastDay, price: sumDay / countDay });
    return onPricePerDay;
  };

  return (
    <>
      <div>
        <div id="data">
          {isPending ? (
            "Loading"
          ) : (
            <StockCanvas
              width={width}
              height={height}
              data={data}
              minPrice={minPrice}
              maxPrice={maxPrice}
              totalDays={totalDays}
              futureDays={futureDays}
              userDrawnPrices={userDrawnPoints}
              addUserDrawnPrice={(x, y) => {
                if (
                  userDrawnPoints.length == 0 ||
                  x > userDrawnPoints[userDrawnPoints.length - 1].x
                ) {
                  setUserDrawnPoints((prev) => [...prev, { x, y }]);
                }
              }}
              clearUserDrawnPrices={clear}
            />
          )}
        </div>
        <div id="controls">
          <button
            onClick={() => {
              clear();
              refetch();
            }}
          >
            New
          </button>
          <button onClick={() => console.log(getUserDrawnPrices())}>
            Get Drawn Prices
          </button>
          <button onClick={() => showSolution()}>Solution</button>
        </div>
      </div>
    </>
  );
}

export default App;
