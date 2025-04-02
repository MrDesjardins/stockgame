import { useEffect, useRef } from "react";
import { animate } from "./animate";
import { usePrefersColorScheme } from "./usePrefersColorScheme";
import { APP_CONSTANTS } from "./model/app";

export interface LoadingCanvasProps {
  width: number;
  height: number;
}
export function LoadingCanvas(props: LoadingCanvasProps) {
  const theme = usePrefersColorScheme();
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animation = useRef(
    animate(1000, 50, (frame: number, _timeFromBeginningMs: number) => {
      // Draw the candles on the left-side first
      // Draw only up to the animated candle count
      const ctx = canvasRef.current?.getContext("2d");
      if (ctx === null || ctx === undefined) {
        return;
      }
      renderLoadingText(ctx, frame);
    })
  );
  const renderLoadingText = (ctx: CanvasRenderingContext2D, frame: number) => {
    ctx.fillStyle = theme === "dark" ? "white" : "black";
    ctx.font = "20px Arial";
    const fullText = "Loading...";
    const textToShow = fullText.slice(0, frame % fullText.length);
    ctx.fillText(
      textToShow,
      APP_CONSTANTS.canvas_width / 2 - fullText.length * 2,
      APP_CONSTANTS.canvas_height / 2 - 10
    );
  };
  useEffect(() => {
    animation.current.start();
  }, []);
  return (
    <canvas ref={canvasRef} width={props.width} height={props.height}></canvas>
  );
}
