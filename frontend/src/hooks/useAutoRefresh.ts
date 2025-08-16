import { useEffect, useRef, useCallback } from 'react';

interface UseAutoRefreshOptions {
  interval?: number; // milliseconds
  enabled?: boolean;
  onRefresh: () => void | Promise<void>;
  dependencies?: unknown[];
}

export const useAutoRefresh = ({
  interval = 30000, // 30 seconds default
  enabled = true,
  onRefresh,
  dependencies = []
}: UseAutoRefreshOptions) => {
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);
  const isActiveRef = useRef(false);
  const isRunningRef = useRef(false);

  const startInterval = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
    }
    
    intervalRef.current = setInterval(async () => {
      if (!isActiveRef.current || !enabled || isRunningRef.current) return;
      isRunningRef.current = true;
      try {
        await onRefresh();
      } catch (error) {
        console.error('Auto-refresh error:', error);
      } finally {
        isRunningRef.current = false;
      }
    }, interval);
  }, [interval, enabled, onRefresh]);

  const stopInterval = useCallback(() => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  }, []);

  const pause = useCallback(() => {
    isActiveRef.current = false;
  }, []);

  const resume = useCallback(() => {
    isActiveRef.current = true;
  }, []);

  // Start/stop interval when dependencies change
  useEffect(() => {
    if (enabled) {
      isActiveRef.current = true;
      startInterval();
    } else {
      isActiveRef.current = false;
      stopInterval();
    }

    return () => {
      stopInterval();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [enabled, startInterval, stopInterval, JSON.stringify(dependencies)]);

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      stopInterval();
    };
  }, [stopInterval]);

  return {
    pause,
    resume,
    isActive: isActiveRef.current
  };
};
