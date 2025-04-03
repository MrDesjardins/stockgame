export function animate(
  id: string,
  animation: (
    frame: number,
    timeFromBeginningMs: number,
    fps: number,
    maxFps: number
  ) => void,
  animationExpectedDuration?: number,
  totalAnimationFramesInitial?: number,
  autoReset?: boolean
) {
  let startTime = 0;
  let animationFrameDone = 0;
  let timePerFrame = 0;
  let totalAnimationFrames = totalAnimationFramesInitial;
  let isDone = false;
  const isAutoReset = autoReset ?? false;
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
      if (
        animationExpectedDuration !== undefined &&
        totalAnimationFrames !== undefined
      ) {
        timePerFrame = animationExpectedDuration / totalAnimationFrames;
      }
    },
    tick: (timestamp: number, fps: number, maxFps: number) => {
      if (startTime === 0) {
        return false;
      }
      const elapsedTime = timestamp - startTime;
      if (
        !isDone &&
        totalAnimationFrames !== undefined &&
        animationFrameDone >= totalAnimationFrames
      ) {
        console.log(id, "Ended at after", elapsedTime, "ms");
        if (isAutoReset) {
          animationFrameDone = 0;
          startTime = timestamp;
          isDone = false;
        } else {
          isDone = true;
        }
      }
      animation(animationFrameDone, elapsedTime, fps, maxFps);
      if (
        !isDone &&
        (totalAnimationFrames === undefined ||
          (totalAnimationFrames !== undefined &&
            animationFrameDone < totalAnimationFrames &&
            elapsedTime >= timePerFrame * animationFrameDone))
      ) {
        animationFrameDone++;
      }

      return true;
    },
    getFrame: () => {
      return animationFrameDone;
    },
    getId: () => {
      return id;
    },
  };
}

export function animationEngine(
  target_fps: number,
  initAnimations?: ReturnType<typeof animate>[]
) {
  let animationFrameRef: number | null = null;
  const animations: ReturnType<typeof animate>[] = [];
  let frameCount = 0;
  let frameRenderedCount = 0;
  let fpsStartTime = performance.now();
  let maxFps: number = target_fps; // Target but will be adjusted for the user's machine
  let actualFps = target_fps; // Target but will be adjusted for the user's machine
  if (initAnimations !== undefined) {
    animations.push(...initAnimations);
  }
  const loopFunc = (timestamp: number) => {
    if (timestamp - fpsStartTime >= 1000) {
      maxFps = frameCount;
      actualFps = frameRenderedCount;
      frameCount = 0;
      frameRenderedCount = 0;
      fpsStartTime = timestamp;
    }

    if (frameCount % Math.floor(maxFps / target_fps) === 0) {
      frameRenderedCount++;

      for (const animation of animations) {
        animation.tick(timestamp, actualFps, maxFps);
      }
    }
    frameCount++;
    animationFrameRef = requestAnimationFrame(loopFunc);
  };

  return {
    addAnimate: (animation: ReturnType<typeof animate>) => {
      animations.push(animation);
    },
    startAll: () => {
      if (animationFrameRef === null) {
        for (const animation of animations) {
          animation.start();
        }
        animationFrameRef = requestAnimationFrame(loopFunc);
      }
    },
    start: (animationId: string) => {
      const animation = animations.find((a) => a.getId() === animationId);
      if (animation) {
        animation.start();
      }
    },
    reset: (animationId: string, value: number) => {
      const animation = animations.find((a) => a.getId() === animationId);
      if (animation) {
        animation.reset(value);
      }
    },
    stopAll: () => {
      if (animationFrameRef !== null) {
        cancelAnimationFrame(animationFrameRef);
        animationFrameRef = null;
      }
    },
  };
}
