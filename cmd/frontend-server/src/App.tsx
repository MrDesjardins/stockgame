import "./App.css";
import { useMutation, useQuery } from "@tanstack/react-query";
import { StockCanvas } from "./StockCanvas";
import {
  SolutionDayPrice,
  SolutionRequest,
  SolutionResponse,
  StockPublic,
} from "./model/stock";
import { useCallback, useMemo, useState } from "react";
import { xPixelToDay, yPixelToPrice } from "./logic/canvasLogic";
import { User_stock_to_guess } from "./dynamicConstants";
import { getApiUrl } from "./logic/apiLogics";
console.log(import.meta.env.VITE_API_URL);
async function getStocks(): Promise<StockPublic[]> {
  const data = await fetch(`${getApiUrl()}/stocks`);
  return data.json();
}

async function postUserDayPrice(
  requestData: SolutionRequest
): Promise<SolutionResponse> {
  const data = await fetch(`${getApiUrl()}/solution`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestData),
  });

  return data.json();
}

function App() {
  const futureDays = User_stock_to_guess;
  const width = 800;
  const height = 600;
  const [response, setResponse] = useState<SolutionResponse | undefined>(
    undefined
  );
  const [responseCounter, setResponseCounter] = useState(0);
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
        symbolUUID: lastEntry.symbol_uuid,
        afterDate: lastEntry.date,
        estimatedDayPrices: dayPrice.slice(0, futureDays),
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
  const totalDays = useMemo(
    () => (data?.length ?? 0) + futureDays,
    [data?.length, futureDays]
  );

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
        // First point or same day
        lastDay = day;
        sumDay += price;
        countDay += 1;
      } else {
        // New day
        onPricePerDay.push({ day: lastDay, price: sumDay / countDay });
        lastDay = day;
        sumDay = price;
        countDay = 1;
      }
    }
    // Add the last day
    if (lastDay !== -1) {
      onPricePerDay.push({ day: lastDay, price: sumDay / countDay });
    }
    return onPricePerDay;
  }, [data, futureDays, maxPrice, minPrice, userDrawnPoints]);

  const submitSolution = useCallback(async () => {
    const userDayPrices = getUserDrawnPrices();
    console.log(userDayPrices);
    const response = await postSolution.mutateAsync(userDayPrices);
    setResponse(response);
    setResponseCounter(() => responseCounter + 1);
  }, [getUserDrawnPrices, postSolution, responseCounter]);

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
              responseCounter={responseCounter}
              numberDaysUserNeedToGuess={User_stock_to_guess}
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
                {response.name} ({response.symbol})
                <span id="date">{response.stocks[0].date}</span>
              </h2>
              <div>Score: {response.score.total}</div>
              <div>In direction: {response.score.inDirection}</div>
              <div>In Low-High: {response.score.inLowHigh}</div>
              <div>In Open-Close: {response.score.inOpenClose}</div>
              <div>In Bollinger Band: {response.score.inBollinger}</div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default App;
