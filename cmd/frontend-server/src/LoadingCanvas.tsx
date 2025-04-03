import { useEffect, useRef } from "react";
import { animate, animationEngine } from "./animate";
import { usePrefersColorScheme } from "./usePrefersColorScheme";
import { APP_CONSTANTS } from "./model/app";

export interface LoadingCanvasProps {
  width: number;
  height: number;
}
const text = "Loading...";
export function LoadingCanvas(props: LoadingCanvasProps) {
  const theme = usePrefersColorScheme();
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const animationEngineRef = useRef(
    animationEngine(APP_CONSTANTS.fps, [
      animate(
        "LoadingAnimation",
        (frame: number, timeFromBeginningMs: number) => {
          //console.log("LoadingAnimation", frame, timeFromBeginningMs);
          // Draw the candles on the left-side first
          // Draw only up to the animated candle count
          const ctx = canvasRef.current?.getContext("2d");
          if (ctx === null || ctx === undefined) {
            return;
          }
          renderLoadingText(ctx, frame, text);
        },
        2000,
        text.length,
        true
      ),
    ])
  );
  const renderLoadingText = (
    ctx: CanvasRenderingContext2D,
    frame: number,
    text: string
  ) => {
    ctx.fillStyle = theme === "dark" ? "white" : "black";
    ctx.font = "28px Arial";
    const textToShow = text.slice(0, frame % text.length);
    ctx.clearRect(
      0,
      0,
      APP_CONSTANTS.canvas_width,
      APP_CONSTANTS.canvas_height
    );
    ctx.fillText(
      textToShow,
      APP_CONSTANTS.canvas_width / 2 - text.length * 5,
      APP_CONSTANTS.canvas_height / 2 - 10
    );
  };
  useEffect(() => {
    animationEngineRef.current.startAll();
  }, []);
  return (
    <canvas ref={canvasRef} width={props.width} height={props.height}></canvas>
  );
}
