import React, { useRef, useEffect, useCallback, useMemo } from "react";
import { Stock } from "./model/stock";
import {
  candelPixelWidth,
  priceToYPixel,
  yPixelToPrice,
} from "./logic/canvasLogic";

export interface StockCanvasProps {
  data: Stock[];
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
      for (let i = 1; i <= 10; i++) {
        const y = (i / 10) * height;
        ctx.beginPath();
        ctx.moveTo(0, y);
        ctx.lineTo(width, y);
        ctx.stroke();

        const price = yPixelToPrice(y, height, props.minPrice, props.maxPrice);
        ctx.fillStyle = "black";
        const yText = y - 4;
        ctx.fillText(`$${price.toFixed(2)}`, 2, yText);
        ctx.fillText(`$${price.toFixed(2)}`, props.width - 30, yText);
      }
    },
    [props.maxPrice, props.minPrice, props.width]
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
        ctx.fillStyle = "#C0C0C0";
        ctx.fill();
      }
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
      });

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
    props.userDrawnPrices,
    stockIndexToX,
  ]);

  const handleMouseDown = (e: React.MouseEvent) => {
    props.clearUserDrawnPrices();
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
