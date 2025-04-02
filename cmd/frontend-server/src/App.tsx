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
import { APP_CONSTANTS } from "./model/app";
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
  const [response, setResponse] = useState<SolutionResponse | undefined>(
    undefined
  );
  const [responseCounter, setResponseCounter] = useState(0);
  const { isPending, isError, data, error, refetch, isFetching } = useQuery({
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
    () => (data == undefined ? 0 : Math.min(...data.map((s) => s.low)) * APP_CONSTANTS.min_price_percent),
    [data]
  );
  const maxPrice = useMemo(
    () => (data == undefined ? 0 : Math.max(...data.map((s) => s.high)) * APP_CONSTANTS.max_price_percent),
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
      const day = xPixelToDay(
        dayPrice.x,
        APP_CONSTANTS.canvas_width,
        data.length + futureDays
      );
      const price = yPixelToPrice(
        dayPrice.y,
        APP_CONSTANTS.canvas_height,
        minPrice,
        maxPrice
      );
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
        <h1>Guess the stock price</h1>
        <div
          id="data"
          style={{
            width: APP_CONSTANTS.canvas_width,
            height: APP_CONSTANTS.canvas_height,
          }}
        >
          {isPending ? (
            "Loading"
          ) : (
            <StockCanvas
              width={APP_CONSTANTS.canvas_width}
              height={APP_CONSTANTS.canvas_height}
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
            disabled={isFetching}
            onClick={() => {
              clear();
              refetch();
            }}
          >
            New
          </button>
          <button
            disabled={isFetching || response !== undefined}
            onClick={() => submitSolution()}
          >
            Submit
          </button>
        </div>
        <div id="solution">
          {response == undefined ? (
            ""
          ) : (
            <div id="solution-with-data">
              <div id="solution-with-data-details">
                <h2>Details </h2>
                <div className="caption-value">
                  <span className="caption">Name:</span>
                  <span className="value">{response.name.slice(0, 50)}</span>
                </div>
                <div className="caption-value">
                  <span className="caption">Symbol:</span>
                  <span className="value">{response.symbol}</span>
                </div>
                <div className="caption-value">
                  <span className="caption">First Guessing Date:</span>
                  <span className="value">
                    {response.stocks[0].date.substring(0, 10)}
                  </span>
                </div>
              </div>
              <div id="solution-with-data-score">
                <h2>Score</h2>
                <div className="caption-value">
                  <span className="caption">Total Score:</span>
                  <span className="value">{response.score.total}</span>
                </div>
                <h2>Detail Score</h2>
                <div className="caption-value">
                  <span className="caption">In direction:</span>
                  <span className="value">{response.score.inDirection}</span>
                </div>
                <div className="caption-value">
                  <span className="caption">In Low-High:</span>
                  <span className="value">{response.score.inLowHigh}</span>
                </div>
                <div className="caption-value">
                  <span className="caption">In Open-Close:</span>
                  <span className="value">{response.score.inOpenClose}</span>
                </div>
                <div className="caption-value">
                  <span className="caption">In Bollinger Band:</span>
                  <span className="value">{response.score.inBollinger}</span>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </>
  );
}

export default App;
