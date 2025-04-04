export interface AppConstants {
  canvas_width: number;
  canvas_height: number;
  min_price_percent: number;
  max_price_percent: number;
  fps: number;
}

export const APP_CONSTANTS: AppConstants = {
  canvas_width: 1000,
  canvas_height: 600,
  min_price_percent: 0.8,
  max_price_percent: 1.2,
  fps: 60,
};
