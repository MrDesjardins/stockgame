import "./App.css";
import { useMutation, useQuery } from "@tanstack/react-query";
import { StockCanvas } from "./StockCanvas";
import {
  SolutionDayPrice,
  SolutionRequest,
  SolutionResponse,
  StockPublic,
} from "./model/stock";
import { useCallback, useEffect, useMemo, useState } from "react";
import { xPixelToDay, yPixelToPrice } from "./logic/canvasLogic";
import {
  Number_initial_stock_shown,
  User_stock_to_guess,
} from "./dynamicConstants";
import { getApiUrl } from "./logic/apiLogics";
import { APP_CONSTANTS } from "./model/app";
import { LoadingCanvas } from "./LoadingCanvas";
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
  const [response, setResponse] = useState<SolutionResponse | undefined>(
    undefined
  );
  const [responseCounter, setResponseCounter] = useState(0);
  const [isNarrowView, setIsNarrowView] = useState(false);
  const { isPending, isError, data, error, refetch, isFetching } = useQuery({
    queryKey: ["stocks"],
    queryFn: getStocks,
  });
  const [dataToShow, setDataToShow] = useState<StockPublic[]>([]);

  useEffect(() => {
    const onResize = () => {
      setIsNarrowView(window.innerWidth < APP_CONSTANTS.canvas_width);
    };
    window.addEventListener("resize", onResize);
    onResize();
    return () => {
      window.removeEventListener("resize", onResize);
    };
  }, []);

  useEffect(() => {
    setDataToShow(data ?? []);
  }, [data]);
  const postSolution = useMutation({
    mutationFn: (dayPrice: SolutionDayPrice[]) => {
      if (dataToShow === undefined) {
        throw new Error("No data");
      }
      const lastEntry = dataToShow[dataToShow.length - 1];
      if (lastEntry == undefined) {
        throw new Error("No data");
      }
      return postUserDayPrice({
        symbolUUID: lastEntry.symbol_uuid,
        afterDate: lastEntry.date,
        estimatedDayPrices: dayPrice.slice(0, User_stock_to_guess),
      });
    },
  });

  // Find the min and max to determine the vertical space needed
  const minPrice = useMemo(
    () =>
      dataToShow.length === 0
        ? 0
        : Math.min(...dataToShow.map((s) => s.low)) *
          APP_CONSTANTS.min_price_percent,
    [dataToShow]
  );
  const maxPrice = useMemo(
    () =>
      dataToShow.length === 0
        ? 0
        : Math.max(...dataToShow.map((s) => s.high)) *
          APP_CONSTANTS.max_price_percent,
    [dataToShow]
  );

  const [userDrawnPoints, setUserDrawnPoints] = useState<
    { x: number; y: number }[]
  >([]);

  const getUserDrawnPrices = useCallback(() => {
    if (dataToShow.length === 0) {
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
        Number_initial_stock_shown + User_stock_to_guess
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
  }, [dataToShow, maxPrice, minPrice, userDrawnPoints]);

  const submitSolution = useCallback(async () => {
    const userDayPrices = getUserDrawnPrices();
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
  const isNewDisabled = isFetching || response === undefined;
  const isSubmitDisabled = isFetching || response !== undefined;
  return (
    <div className={"app-container" + (isNarrowView ? " narrow" : "")}>
      <h1>Guess the stock price</h1>
      <div id="instruction">
        Click the gray area and drag from left to right where you think the
        price will move in the next 10 days. You can redraw and when ready click
        submit.
      </div>
      <div
        id="data"
        style={{
          width: APP_CONSTANTS.canvas_width,
          height: APP_CONSTANTS.canvas_height,
        }}
      >
        {isPending ? (
          <LoadingCanvas
            width={APP_CONSTANTS.canvas_width}
            height={APP_CONSTANTS.canvas_height}
          />
        ) : (
          <StockCanvas
            width={APP_CONSTANTS.canvas_width}
            height={APP_CONSTANTS.canvas_height}
            data={dataToShow}
            response={response}
            minPrice={minPrice}
            maxPrice={maxPrice}
            totalDays={Number_initial_stock_shown + User_stock_to_guess}
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
          className={!isNewDisabled ? "active" : ""}
          disabled={isNewDisabled}
          onClick={() => {
            clear();
            refetch();
            setDataToShow([]);
          }}
        >
          New
        </button>
        <button
          className={!isSubmitDisabled ? "active" : ""}
          disabled={isSubmitDisabled}
          onClick={() => submitSolution()}
        >
          Submit
        </button>
      </div>
      <div id="solution">
        {response === undefined ? (
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
  );
}

export default App;
