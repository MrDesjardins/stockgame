/**
 * Represent a single animation
 * @param id - A unique name to identify the animation. Useful during debugging.
 * @param animation - The function to be called for each frame of the animation.
 * @param animationExpectedDuration  - The time between frame
 * @param totalAnimationFramesInitial - The total number of frames for the animation before stopping. Undefined mean infinite animation.
 * @param autoReset - When true, the animation will reset when it reaches the end. When false, the animation will stop at the end.
 * @returns - An object with commands
 */
export function animate(
  id: string,
  animation: (
    frame: number,
    timeFromBeginningMs: number,
    fps: number,
    maxFps: number
  ) => void,
  animationExpectedDuration?: number,
  totalAnimationFramesInitial?: () => number,
  autoReset?: boolean
) {
  let startTime = 0;
  let animationFrameDone = 0;

  let isDone = false;
  const isAutoReset = autoReset ?? false;
  return {
    reset: () => {
      animationFrameDone = 0;
      startTime = 0;
      isDone = true;
    },
    start: () => {
      startTime = performance.now();
      isDone = false;
      animationFrameDone = 0;
    },
    tick: (timestamp: number, fps: number, maxFps: number) => {
      if (startTime === 0) {
        return false;
      }
      const elapsedTime = timestamp - startTime;
      if (
        !isDone &&
        totalAnimationFramesInitial !== undefined &&
        animationFrameDone >= totalAnimationFramesInitial()
      ) {
        if (isAutoReset) {
          animationFrameDone = 0;
          startTime = timestamp;
          isDone = false;
        } else {
          isDone = true;
        }
      }
      animation(animationFrameDone, elapsedTime, fps, maxFps);
      if (!isDone) {
        let timePerFrame = 0;
        if (
          animationExpectedDuration !== undefined &&
          totalAnimationFramesInitial !== undefined
        ) {
          timePerFrame =
            animationExpectedDuration / totalAnimationFramesInitial();
        }
        if (
          totalAnimationFramesInitial === undefined ||
          (totalAnimationFramesInitial !== undefined &&
            animationFrameDone < totalAnimationFramesInitial() &&
            elapsedTime >= timePerFrame * animationFrameDone)
        ) {
          animationFrameDone++;
        }
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

/**
 * The engine is the single loop for all animations. It controls the frame per second and ensure that the target fps is
 * respected before calling individual animations.
 * @param target_fps - The number of time we repaint the canvas per second
 * @param initAnimations - The list of animation to run every frame.
 * @returns - Commands on the engine
 */
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
    reset: (animationId: string) => {
      const animation = animations.find((a) => a.getId() === animationId);
      if (animation) {
        animation.reset();
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
