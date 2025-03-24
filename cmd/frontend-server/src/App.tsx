import "./App.css";
import { useMutation, useQuery } from "@tanstack/react-query";
import { StockCanvas } from "./StockCanvas";
import {
  SolutionDayPrice,
  SolutionRequest,
  SolutionResponse,
  Stock,
} from "./model/stock";
import { useCallback, useMemo, useState } from "react";
import { xPixelToDay, yPixelToPrice } from "./logic/canvasLogic";

async function getStocks(): Promise<Stock[]> {
  const data = await fetch("http://localhost:8080/stocks");
  return data.json();
}

async function postUserDayPrice(
  requestData: SolutionRequest
): Promise<SolutionResponse> {
  const data = await fetch(`http://localhost:8080/solution`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestData),
  });

  return data.json();
}

function App() {
  const futureDays = 10;
  const width = 800;
  const height = 600;
  const [response, setResponse] = useState<SolutionResponse | undefined>(
    undefined
  );
  const { isPending, isError, data, error, refetch } = useQuery({
    queryKey: ["stocks"],
    queryFn: getStocks,
  });

  const postSolution = useMutation({
    mutationFn: (dayPrice: SolutionDayPrice[]) => {
      const lastEntry = data![data!.length - 1];
      if (lastEntry == undefined) {
        throw new Error("No data");
      }
      return postUserDayPrice({
        symbol: lastEntry.symbol,
        afterDate: lastEntry.date,
        estimatedDayPrices: dayPrice,
      });
    },
  });

  // Find the min and max to determine the vertical space needed
  const minPrice = useMemo(
    () => (data == undefined ? 0 : Math.min(...data.map((s) => s.low)) * 0.85),
    [data]
  );
  const maxPrice = useMemo(
    () => (data == undefined ? 0 : Math.max(...data.map((s) => s.high)) * 1.15),
    [data]
  );
  const totalDays = useMemo(() => (data?.length ?? 0) + futureDays, [data]);

  const [userDrawnPoints, setUserDrawnPoints] = useState<
    { x: number; y: number }[]
  >([]);

  const getUserDrawnPrices = useCallback(() => {
    if (data == undefined) {
      return [];
    }
    let lastDay = -1;
    const onPricePerDay: SolutionDayPrice[] = [];
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
  }, [data, maxPrice, minPrice, userDrawnPoints]);

  const submitSolution = useCallback(async () => {
    const userDayPrices = getUserDrawnPrices();
    console.log(userDayPrices);
    const response = await postSolution.mutateAsync(userDayPrices);
    setResponse(response);
  }, [getUserDrawnPrices, postSolution]);

  const clear = () => {
    setUserDrawnPoints([]);
    setResponse(undefined);
  };

  if (isError) {
    return <div>Error: {error.message}</div>;
  }

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
              response={response}
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
          <button
            disabled={response !== undefined}
            onClick={() => submitSolution()}
          >
            Submit
          </button>
        </div>
        <div id="solution">
          {response == undefined ? (
            ""
          ) : (
            <div>
              <h2>
                {response.symbol}
                <span id="date">{response.stocks[0].date}</span>
              </h2>
              <p>Score: {response.score}</p>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default App;
