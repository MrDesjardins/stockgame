export function animate(
  animationExpectedDuration: number,
  totalAnimationFramesInitial: number,
  animation: (frame: number, timeFromBeginningMs: number) => void
) {
  let startTime = 0;
  let animationFrameDone = 0;
  let timePerFrame = 0;
  let totalAnimationFrames = totalAnimationFramesInitial;
  let isDone = false;
  return {
    reset: (totalAnimationFramesInitial: number) => {
      animationFrameDone = 0;
      startTime = 0;
      isDone = false;
      timePerFrame = 0;
      totalAnimationFrames = totalAnimationFramesInitial;
    },
    start: () => {
      startTime = performance.now();
      isDone = false;
      animationFrameDone = 0;
      timePerFrame = animationExpectedDuration / totalAnimationFrames;
    },
    tick: (timestamp: number) => {
      if (startTime === 0) {
        return false;
      }
      const elapsedTime = timestamp - startTime;
      if (!isDone && animationFrameDone >= totalAnimationFrames) {
        console.log("Ended at after", elapsedTime, "ms");
        isDone = true;
      }

      if (
        !isDone &&
        animationFrameDone < totalAnimationFrames &&
        elapsedTime >= timePerFrame * animationFrameDone
      ) {
        animationFrameDone++;
        console.log(animationFrameDone, elapsedTime);
      }
      animation(animationFrameDone, elapsedTime);
      return true;
    },
    getFrame: () => {
      return animationFrameDone;
    },
  };
}
