export function xPixelToDay(x: number, width: number, totalDays: number): number {
  return Math.floor((x / width) * totalDays);
}

export function yPixelToPrice(y: number, height: number, minPrice: number, maxPrice: number): number {
  return maxPrice - (y / height) * (maxPrice - minPrice);
}

export function priceToYPixel(price: number, height: number, minPrice: number, maxPrice: number): number {
  return height - ((price - minPrice) / (maxPrice - minPrice)) * height;
}

export function candelPixelWidth(width: number, totalDays: number): number {
  return width / totalDays;
}