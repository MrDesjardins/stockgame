import React, { useRef, useEffect, useCallback, useMemo } from "react";
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
}
const NUMBER_PRICE_SHOW = 10;
export function StockCanvas(props: StockCanvasProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);

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
  // Draw horizontal grid lines
  const drawGrid = useCallback(
    (ctx: CanvasRenderingContext2D, width: number, height: number) => {
      ctx.strokeStyle = "#DDD";
      ctx.lineWidth = 1;
      for (let i = 1; i <= NUMBER_PRICE_SHOW; i++) {
        const y = (i / NUMBER_PRICE_SHOW) * height;
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y);
        ctx.stroke();

        const price = yPixelToPrice(y, height, props.minPrice, props.maxPrice);
        ctx.fillStyle = "black";
        const yText = y - 4;
        ctx.fillText(`$${price.toFixed(2)}`, 2, yText);
        if (i < NUMBER_PRICE_SHOW) {
          ctx.fillText(`$${price.toFixed(2)}`, props.width - 30, yText);
        }
      }

      for (let i = 0; i < props.futureDays; i++) {
        const x = props.data.length * candleWidth + i * candleWidth;
        ctx.beginPath();
        ctx.moveTo(x, 0);
        ctx.lineTo(x, height);
        ctx.stroke();
        ctx.fillText(`+${i + 1}`, x + 2, height - 2);
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
          props.response === undefined ? 1 : 0.3
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
    const drawChart = () => {
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw the area the user can predict the price
      ctx.fillStyle = "#F0F0F0";
      ctx.fillRect(
        props.data.length * candleWidth,
        0,
        props.futureDays * candleWidth,
        canvas.height
      );

      // Draw each day's candle
      props.data.forEach((stock, index) => {
        drawSingleStock(ctx, stock, index);
      });
      if (props.response !== undefined) {
        for (let i = 0; i < props.response.stocks.length; i++) {
          drawSingleStock(ctx, props.response.stocks[i], props.data.length + i);
        }
        drawBB20(ctx, props.response.stocks, props.response.bb20);
      }

      drawGrid(ctx, canvas.width, canvas.height);
      drawUserPoints(ctx);
    };

    drawChart();
  }, [
    candleWidth,
    drawGrid,
    props.data,
    props.futureDays,
    props.maxPrice,
    props.minPrice,
    props.response,
    props.userDrawnPrices,
    stockIndexToX,
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
