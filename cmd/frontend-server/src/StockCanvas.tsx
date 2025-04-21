import React, { useRef, useEffect } from "react";
import { SolutionResponse, StockPublic } from "./model/stock";
import {
  candelPixelWidth,
  priceToYPixel,
  yPixelToPrice,
} from "./logic/canvasLogic";
import { animate, animationEngine } from "./animate";
import { usePrefersColorScheme } from "./usePrefersColorScheme";
import { APP_CONSTANTS } from "./model/app";

export interface StockCanvasProps {
  data: StockPublic[];
  response: SolutionResponse | undefined;
  minPrice: number;
  maxPrice: number;
  totalDays: number;
  width: number;
  height: number;
  userDrawnPrices: { x: number; y: number }[];
  addUserDrawnPrice: (x: number, y: number) => void;
  clearUserDrawnPrices: () => void;
  responseCounter: number;
  numberDaysUserNeedToGuess: number;
}

export function StockCanvas(props: StockCanvasProps) {
  const theme = usePrefersColorScheme();
  const canvasRef = useRef<HTMLCanvasElement>(null);

  // Animation of the data candles
  const dataRef = useRef<StockPublic[]>([]);
  const responseDataRef = useRef<SolutionResponse | undefined>(undefined);
  const minPriceRef = useRef(props.minPrice);
  const maxPriceRef = useRef(props.maxPrice);
  const userDrawnPricesRef = useRef<{ x: number; y: number }[]>([]);

  const animationEngineRef = useRef(
    animationEngine(APP_CONSTANTS.fps, [
      animate(
        "Background",
        (
          _frame: number,
          _timeFromBeginning: number,
          fps: number,
          maxFps: number
        ) => {
          const ctx = canvasRef.current?.getContext("2d");
          if (ctx !== null && ctx !== undefined) {
            ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);

            renderFPS(ctx, maxFps, fps);
            renderFullChart(ctx);
          }
        }
      ),
      animate(
        "StockSingleStocks",
        (frame: number, _timeFromBeginningMs: number) => {
          // Draw the candles on the left-side first
          // Draw only up to the animated candle count
          const ctx = canvasRef.current?.getContext("2d");
          if (ctx === null || ctx === undefined) {
            return;
          }

          dataRef.current.slice(0, frame).forEach((stock, index) => {
            renderSingleStock(ctx, stock, index);
          });
        },
        1000,
        () => dataRef.current.length
      ),
      animate(
        "StockResponsesStocks",
        (frame: number, _timeFromBeginningMs: number) => {
          const ctx = canvasRef.current?.getContext("2d");
          if (ctx === null || ctx === undefined) {
            return;
          }
          if (responseDataRef.current !== undefined) {
            // Draw only up to the animated candle count
            responseDataRef.current.stocks
              .slice(0, frame)
              .forEach((stock, index) => {
                renderSingleStock(ctx, stock, dataRef.current.length + index);
              });
          }
        },
        500,
        () => responseDataRef.current?.stocks.length ?? 0
      ),
    ])
  );

  const responseChangedCounterId = useRef(0);
  const candleWidth = useRef(0);
  useEffect(() => {
    candleWidth.current = candelPixelWidth(props.width, props.totalDays);
  }, [props.width, props.totalDays]);

  // Give the most left side of the candle
  const stockIndexToX = (dayIndex: number): number => {
    return dayIndex * candleWidth.current;
  };

  // New data
  useEffect(() => {
    // Don't animate if there's no data
    if (props.data.length === 0) {
      return;
    }
    // Reset animation state and start
    animationEngineRef.current.start("StockSingleStocks");
    animationEngineRef.current.reset("StockResponsesStocks");

    responseDataRef.current = undefined;
    // Set the data into a ref for the animation loop
    dataRef.current = props.data;
  }, [props.data]);

  useEffect(() => {
    // Reset animation state and start

    if (
      props.responseCounter > 0 &&
      props.responseCounter !== responseChangedCounterId.current
    ) {
      responseDataRef.current = props.response;
      responseChangedCounterId.current = props.responseCounter;

      animationEngineRef.current.reset("StockResponsesStocks");
      animationEngineRef.current.start("StockResponsesStocks");
    }
  }, [props.response, props.responseCounter]);

  useEffect(() => {
    userDrawnPricesRef.current = props.userDrawnPrices;
  }, [props.userDrawnPrices]);

  useEffect(() => {
    minPriceRef.current = props.minPrice;
    maxPriceRef.current = props.maxPrice;
  }, [props.minPrice, props.maxPrice]);
  // Animation effect
  useEffect(() => {
    animationEngineRef.current.startAll();
  }, []);
  const renderFPS = (
    ctx: CanvasRenderingContext2D,
    fpsMax: number,
    fps: number
  ) => {
    ctx.fillStyle = theme === "dark" ? "white" : "black";
    ctx.font = "16px Arial";
    ctx.fillText(`Max FPS: ${fpsMax}`, 1, 20);
    ctx.fillText(`FPS: ${fps}`, 1, 40);
  };
  // Draw horizontal/vertical grid lines in the background
  const renderBackgroundGridLines = (
    ctx: CanvasRenderingContext2D,
    width: number,
    height: number
  ) => {
    ctx.strokeStyle = "#DDD";
    ctx.lineWidth = 1;
    // Horizontal lines
    for (let i = 1; i <= props.numberDaysUserNeedToGuess; i++) {
      const y = (i / props.numberDaysUserNeedToGuess) * height;
      ctx.beginPath();
      ctx.moveTo(0, y);
      ctx.lineTo(width, y);
      ctx.stroke();

      // Add price text
      const price = yPixelToPrice(
        y,
        height,
        minPriceRef.current,
        maxPriceRef.current
      );
      ctx.fillStyle = theme === "dark" ? "white" : "black";
      const yText = y - 4;
      ctx.fillText(`$${price.toFixed(2)}`, 2, yText);

      // The right side price
      if (i < props.numberDaysUserNeedToGuess) {
        const priceStr = price.toFixed(2);
        ctx.fillText(
          `$${price.toFixed(2)}`,
          props.width - (priceStr.length >= 6 ? 75 : 50),
          yText
        );
      }
    }

    // Vertical lines
    for (let i = 0; i < props.numberDaysUserNeedToGuess; i++) {
      const x =
        dataRef.current.length * candleWidth.current + i * candleWidth.current;
      ctx.beginPath();
      ctx.moveTo(x, 0);
      ctx.lineTo(x, height);
      ctx.stroke();
      // Text of the days to guess (+1, +2, ...)
      ctx.font = "10px Arial";
      ctx.fillText(`+${i + 1}`, x - 2, height - 2);
    }
  };

  const renderUserPoints = (ctx: CanvasRenderingContext2D) => {
    for (const point of userDrawnPricesRef.current) {
      ctx.beginPath();
      ctx.arc(point.x, point.y, 5, 0, 2 * Math.PI);
      ctx.fillStyle = `rgba(192, 192, 192, ${
        responseDataRef.current?.stocks ? 0.5 : 1
      })`;
      ctx.fill();
    }
  };

  const renderStockVolume = (ctx: CanvasRenderingContext2D) => {
    const volumeMax = Math.max(
      ...dataRef.current.map((stock) => stock.volume),
      0
    );
    const volumeHeight = ctx.canvas.height / 4;
    const padding = 2;

    for (let i = 0; i < dataRef.current.length; i++) {
      const stock = dataRef.current[i];
      const x = stockIndexToX(i);
      const y = priceToYPixel(stock.volume, volumeHeight, 0, volumeMax);
      ctx.fillStyle = `rgba(192, 192, 192, 0.4)`;
      ctx.fillRect(
        x + padding,
        ctx.canvas.height - y,
        candleWidth.current - 2 * padding,
        y
      );
    }
  };
  const renderSingleStock = (
    ctx: CanvasRenderingContext2D,
    stock: StockPublic,
    index: number
  ) => {
    const x = stockIndexToX(index);
    const openY = priceToYPixel(
      stock.open,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    const closeY = priceToYPixel(
      stock.close,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    const highY = priceToYPixel(
      stock.high,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    const lowY = priceToYPixel(
      stock.low,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    const isBullish = stock.close > stock.open;

    ctx.strokeStyle = "black";
    ctx.fillStyle = isBullish ? "green" : "red";

    // Draw the candle high and low
    ctx.beginPath();
    ctx.moveTo(x + candleWidth.current / 2, highY);
    ctx.lineTo(x + candleWidth.current / 2, lowY);
    ctx.stroke();

    // Draw the candle body
    ctx.fillRect(
      x,
      Math.min(openY, closeY),
      candleWidth.current * 0.8,
      Math.abs(closeY - openY)
    );
    // Draw the candle border
    ctx.strokeRect(
      x,
      Math.min(openY, closeY),
      candleWidth.current * 0.8,
      Math.abs(closeY - openY)
    );
  };

  const renderBB20 = (
    ctx: CanvasRenderingContext2D,
    stocks: StockPublic[],
    bb20: Record<string, { upperBand: number; lowerBand: number }>
  ) => {
    ctx.strokeStyle = "blue";
    ctx.lineWidth = 1;
    const firstStock = stocks[0];
    if (Object.keys(bb20).length === 0) {
      return;
    }
    let x = stockIndexToX(dataRef.current.length - 1);
    let y = priceToYPixel(
      bb20[firstStock.date].upperBand,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    ctx.beginPath();
    ctx.moveTo(x, y);
    for (let i = 1; i < stocks.length; i++) {
      const stock = stocks[i];
      const band = bb20[stock.date];
      if (band === undefined) {
        continue;
      }
      y = priceToYPixel(
        band.upperBand,
        ctx.canvas.height,
        minPriceRef.current,
        maxPriceRef.current
      );
      x = stockIndexToX(dataRef.current.length - 1 + i);
      ctx.lineTo(x + candleWidth.current / 2, y);
    }
    ctx.stroke();
    ctx.strokeStyle = "red";
    ctx.lineWidth = 1;
    ctx.beginPath();
    x = stockIndexToX(dataRef.current.length - 1);
    y = priceToYPixel(
      bb20[firstStock.date].lowerBand,
      ctx.canvas.height,
      minPriceRef.current,
      maxPriceRef.current
    );
    ctx.moveTo(x, y);
    for (let i = 1; i < stocks.length; i++) {
      const stock = stocks[i];
      const band = bb20[stock.date];
      if (band === undefined) {
        continue;
      }
      y = priceToYPixel(
        band.lowerBand,
        ctx.canvas.height,
        minPriceRef.current,
        maxPriceRef.current
      );
      x = stockIndexToX(dataRef.current.length - 1 + i);
      ctx.lineTo(x + candleWidth.current / 2, y);
    }
    ctx.stroke();
  };

  const renderFullChart = (ctx: CanvasRenderingContext2D) => {
    // Draw the area the user can predict the price (the user can draw, we make the background light gray)
    ctx.fillStyle = theme === "dark" ? "#3b3b3b" : "#F0F0F0";
    ctx.fillRect(
      dataRef.current.length * candleWidth.current,
      0,
      props.numberDaysUserNeedToGuess * candleWidth.current,
      ctx.canvas.height
    );

    renderBackgroundGridLines(ctx, ctx.canvas.width, ctx.canvas.height);
    renderUserPoints(ctx);
    renderStockVolume(ctx);
    // Add the Bollinger Bands
    if (responseDataRef.current !== undefined) {
      renderBB20(
        ctx,
        responseDataRef.current.stocks,
        responseDataRef.current.bb20
      );
    }
  };

  const handleMouseDown = (_e: React.MouseEvent) => {
    if (responseDataRef.current === undefined) {
      props.clearUserDrawnPrices();
    }
  };
  const handleTouchStart = (_e: React.TouchEvent) => {
    if (responseDataRef.current === undefined) {
      props.clearUserDrawnPrices();
    }
  };
  const handleTouchMove = (e: React.TouchEvent) => {
    // Support left click mouse only
    if (e.touches.length !== 1) {
      return;
    }

    handleMove(e.touches[0].clientX, e.touches[0].clientY);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    // Support left click mouse only
    if (e.buttons !== 1) {
      return;
    }

    handleMove(e.clientX, e.clientY);
  };

  const handleMove = (xClick: number, yClick: number) => {
    // Must have the canvas ready
    if (!canvasRef.current) {
      return;
    }
    if (responseDataRef.current !== undefined) {
      return;
    }

    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const x = xClick - rect.left;
    const y = yClick - rect.top;
    if (
      x < stockIndexToX(dataRef.current.length - 1) + candleWidth.current ||
      x > canvas.width
    ) {
      return;
    }
    props.addUserDrawnPrice(x, y);
  };

  return (
    <canvas
      ref={canvasRef}
      width={props.width}
      height={props.height}
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
    ></canvas>
  );
}
