import React, {
  useRef,
  useEffect,
  useCallback,
  useMemo,
  useState,
} from "react";
import { SolutionResponse, Stock } from "./model/stock";
import {
  candelPixelWidth,
  priceToYPixel,
  yPixelToPrice,
} from "./logic/canvasLogic";

export interface StockCanvasProps {
  data: Stock[];
  response: SolutionResponse | undefined;
  minPrice: number;
  maxPrice: number;
  totalDays: number;
  width: number;
  height: number;
  futureDays: number;
  userDrawnPrices: { x: number; y: number }[];
  addUserDrawnPrice: (x: number, y: number) => void;
  clearUserDrawnPrices: () => void;
  responseCounter: number;
}
const NUMBER_PRICE_SHOW = 10;
const ANIMATION_SPEED = 20; // ms per candle
export function StockCanvas(props: StockCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  // Animation of the data candles
  const [animatedDataCandleCount, setAnimatedCandleCount] = useState(0);
  const [isAnimatingData, setIsAnimating] = useState(false);
  const animationDataRef = useRef<number | null>(null);
  const prevDataLengthRef = useRef<number>(0);
  const prevDataFirstDateRef = useRef<string>("");
  // Animation of the response candles
  const [animatedResponseCandleCount, setAnimatedResponseCandleCount] =
    useState(0);
  const [isAnimatingResponse, setIsAnimatingResponse] = useState(false);
  const animationResponseRef = useRef<number | null>(null);
  const prevResponseCounterRef = useRef<number>(0);
  const candleWidth = useMemo(() => {
    return candelPixelWidth(props.width, props.totalDays);
  }, [props.width, props.totalDays]);

  // Give the most left side of the candle
  const stockIndexToX = useCallback(
    (dayIndex: number): number => {
      return dayIndex * candleWidth;
    },
    [candleWidth]
  );

  // Reset and start animation when new data is received
  useEffect(() => {
    // Don't animate if there's no data
    if (props.data.length === 0) {
      return;
    }

    // Check if we received new data by comparing length and first item date
    // Lenght because starts with nothing then with about 40 candles, symbol because will be different each roll
    const currentDataLength = props.data.length;
    const currentSymbol = props.data[0].symbol;

    // Only trigger animation if the data has actually changed (length or symbol)
    if (
      currentDataLength !== prevDataLengthRef.current ||
      currentSymbol !== prevDataFirstDateRef.current
    ) {
      // Update refs for next comparison
      prevDataLengthRef.current = currentDataLength;
      prevDataFirstDateRef.current = currentSymbol;

      // Reset animation state and start
      setAnimatedCandleCount(0);
      setIsAnimating(true);
    }
  }, [props.data]);

  // Animation effect
  useEffect(() => {
    if (!isAnimatingData) {
      return;
    }

    let lastTimestamp = 0;
    const animate = (timestamp: number) => {
      if (!isAnimatingData) {
        return;
      }

      if (lastTimestamp === 0 || timestamp - lastTimestamp > ANIMATION_SPEED) {
        lastTimestamp = timestamp;
        setAnimatedCandleCount((prev) => {
          if (prev >= props.data.length) {
            setIsAnimating(false);
            return props.data.length;
          }
          return prev + 1;
        });
      }

      animationDataRef.current = requestAnimationFrame(animate);
    };

    animationDataRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationDataRef.current) {
        cancelAnimationFrame(animationDataRef.current);
      }
    };
  }, [isAnimatingData, props.data.length]);

  useEffect(() => {
    // Only trigger the animation once
    if (prevResponseCounterRef.current !== props.responseCounter) {
      setAnimatedResponseCandleCount(0);
      setIsAnimatingResponse(true);
    }
  }, [props.responseCounter]);

  useEffect(() => {
    if (!isAnimatingResponse) {
      return;
    }
    let lastTimestamp = 0;
    const animate = (timestamp: number) => {
      if (!isAnimatingResponse) {
        return;
      }

      if (lastTimestamp === 0 || timestamp - lastTimestamp > ANIMATION_SPEED) {
        lastTimestamp = timestamp;
        setAnimatedResponseCandleCount((prev) => {
          if (prev >= props.data.length) {
            setIsAnimatingResponse(false);
            return props.data.length - 1;
          }
          return prev + 1;
        });
      }

      animationResponseRef.current = requestAnimationFrame(animate);
    };

    animationResponseRef.current = requestAnimationFrame(animate);

    return () => {
      if (animationResponseRef.current) {
        cancelAnimationFrame(animationResponseRef.current);
      }
    };
  }, [isAnimatingData, isAnimatingResponse, props.data.length]);

  // Draw horizontal/vertical grid lines in the background
  const drawBackgroundGridLines = useCallback(
    (ctx: CanvasRenderingContext2D, width: number, height: number) => {
      ctx.strokeStyle = "#DDD";
      ctx.lineWidth = 1;
      // Horizontal lines
      for (let i = 1; i <= NUMBER_PRICE_SHOW; i++) {
        const y = (i / NUMBER_PRICE_SHOW) * height;
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y);
        ctx.stroke();

        // Add price text
        const price = yPixelToPrice(y, height, props.minPrice, props.maxPrice);
        ctx.fillStyle = "black";
        const yText = y - 4;
        ctx.fillText(`$${price.toFixed(2)}`, 2, yText);
        if (i < NUMBER_PRICE_SHOW) {
          ctx.fillText(`$${price.toFixed(2)}`, props.width - 30, yText);
        }
      }

      // Vertical lines
      for (let i = 0; i < props.futureDays; i++) {
        const x = props.data.length * candleWidth + i * candleWidth;
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, height);
        ctx.stroke();
        // Text of the days to guess (+1, +2, ...)
        ctx.fillText(`+${i + 1}`, x - 2, height - 2);
      }
    },
    [
      candleWidth,
      props.data.length,
      props.futureDays,
      props.maxPrice,
      props.minPrice,
      props.width,
    ]
  );

  useEffect(() => {
    if (!canvasRef.current) {
      // Not yet ready, the canvas will be initialized in the next render
      return;
    }
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d")!;

    const drawUserPoints = (ctx: CanvasRenderingContext2D) => {
      for (const point of props.userDrawnPrices) {
        ctx.beginPath();
        ctx.arc(point.x, point.y, 5, 0, 2 * Math.PI);
        ctx.fillStyle = `rgba(192, 192, 192, ${
          props.response === undefined ? 1 : 0.2
        })`;
        ctx.fill();
      }
    };

    const drawSingleStock = (
      ctx: CanvasRenderingContext2D,
      stock: Stock,
      index: number
    ) => {
      const x = stockIndexToX(index);
      const openY = priceToYPixel(
        stock.open,
        canvas.height,
        props.minPrice,
        props.maxPrice
      );
      const closeY = priceToYPixel(
        stock.close,
        canvas.height,
        props.minPrice,
        props.maxPrice
      );
      const highY = priceToYPixel(
        stock.high,
        canvas.height,
        props.minPrice,
        props.maxPrice
      );
      const lowY = priceToYPixel(
        stock.low,
        canvas.height,
        props.minPrice,
        props.maxPrice
      );
      const isBullish = stock.close > stock.open;

      ctx.strokeStyle = "black";
      ctx.fillStyle = isBullish ? "green" : "red";

      // Draw the candle high and low
      ctx.beginPath();
      ctx.moveTo(x + candleWidth / 2, highY);
      ctx.lineTo(x + candleWidth / 2, lowY);
      ctx.stroke();

      // Draw the candle body
      ctx.fillRect(
        x,
        Math.min(openY, closeY),
        candleWidth * 0.8,
        Math.abs(closeY - openY)
      );
      // Draw the candle border
      ctx.strokeRect(
        x,
        Math.min(openY, closeY),
        candleWidth * 0.8,
        Math.abs(closeY - openY)
      );
    };

    const drawBB20 = (
      ctx: CanvasRenderingContext2D,
      stocks: Stock[],
      bb20: Record<string, { upperBand: number; lowerBand: number }>
    ) => {
      ctx.strokeStyle = "blue";
      ctx.lineWidth = 1;
      const firstStock = stocks[0];
      let x = stockIndexToX(props.data.length - 1);
      let y = priceToYPixel(
        bb20[firstStock.date].upperBand,
        canvas.height,
        props.minPrice,
        props.maxPrice
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
          canvas.height,
          props.minPrice,
          props.maxPrice
        );
        x = stockIndexToX(props.data.length - 1 + i);
        ctx.lineTo(x + candleWidth / 2, y);
      }
      ctx.stroke();
      ctx.strokeStyle = "red";
      ctx.lineWidth = 1;
      ctx.beginPath();
      x = stockIndexToX(props.data.length - 1);
      y = priceToYPixel(
        bb20[firstStock.date].lowerBand,
        canvas.height,
        props.minPrice,
        props.maxPrice
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
          canvas.height,
          props.minPrice,
          props.maxPrice
        );
        x = stockIndexToX(props.data.length - 1 + i);
        ctx.lineTo(x + candleWidth / 2, y);
      }
      ctx.stroke();
    };

    const drawFullChart = () => {
      // Reset the canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw the area the user can predict the price (the user can draw, we make the background light gray)
      ctx.fillStyle = "#F0F0F0";
      ctx.fillRect(
        props.data.length * candleWidth,
        0,
        props.futureDays * candleWidth,
        canvas.height
      );

      // Draw the candles on the left-side first
      // Draw only up to the animated candle count
      const displayCount = isAnimatingData
        ? animatedDataCandleCount
        : props.data.length;
      props.data.slice(0, displayCount).forEach((stock, index) => {
        drawSingleStock(ctx, stock, index);
      });

      // Draw the response candles (solution on the right side)
      if (props.response !== undefined) {
        // Draw only up to the animated candle count
        const displayResponseCount = isAnimatingResponse
          ? animatedResponseCandleCount
          : props.response.stocks.length;

        props.response.stocks
          .slice(0, displayResponseCount)
          .forEach((stock, index) => {
            drawSingleStock(ctx, stock, props.data.length + index);
          });

        // Add the Bollinger Bands
        drawBB20(ctx, props.response.stocks, props.response.bb20);
      }

      drawBackgroundGridLines(ctx, canvas.width, canvas.height);
      drawUserPoints(ctx);
    };

    drawFullChart();
  }, [
    candleWidth,
    drawBackgroundGridLines,
    props.data,
    props.futureDays,
    props.maxPrice,
    props.minPrice,
    props.response,
    props.userDrawnPrices,
    stockIndexToX,
    animatedDataCandleCount,
    isAnimatingData,
    isAnimatingResponse,
    animatedResponseCandleCount,
  ]);

  const handleMouseDown = (e: React.MouseEvent) => {
    if (props.response === undefined) {
      props.clearUserDrawnPrices();
    }
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    // Support left click mouse only
    if (e.buttons !== 1) {
      return;
    }
    // Must have the canvas ready
    if (!canvasRef.current) {
      return;
    }

    if (props.response !== undefined) {
      return;
    }
    const canvas = canvasRef.current;
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    if (
      x < stockIndexToX(props.data.length - 1) + candleWidth ||
      x > canvas.width
    ) {
      return;
    }
    props.addUserDrawnPrice(x, y);
  };

  return (
    <div>
      <canvas
        ref={canvasRef}
        width={props.width}
        height={props.height}
        style={{ border: "1px solid black" }}
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
      ></canvas>
    </div>
  );
}
