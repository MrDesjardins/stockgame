export function xPixelToDay(
  x: number,
  width: number,
  totalDays: number
): number {
  return Math.floor((x / width) * totalDays);
}

export function yPixelToPrice(
  y: number,
  height: number,
  minPrice: number,
  maxPrice: number
): number {
  return maxPrice - (y / height) * (maxPrice - minPrice);
}

/**
 *
 * @param price Price of the stock to be represented by a pixel return
 * @param height  The canvas height
 * @param minPrice  The min (bottom of the canvas) price of the stock
 * @param maxPrice  The max (top of the canvas) price of the stock
 * @returns The y pixel value of the price (from the top-left corner of the canvas)
 */
export function priceToYPixel(
  price: number,
  height: number,
  minPrice: number,
  maxPrice: number
): number {
  return height - ((price - minPrice) / (maxPrice - minPrice)) * height;
}

/**
 *
 * @param width The Canvas width
 * @param totalDays The total amount of days to be represented in the width
 * @returns The total pixel width of a single day
 */
export function candelPixelWidth(width: number, totalDays: number): number {
  return width / totalDays;
}
