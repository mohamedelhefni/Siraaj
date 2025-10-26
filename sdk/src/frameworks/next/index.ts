import { useEffect, useRef } from 'react';
import { useRouter, usePathname, useSearchParams } from 'next/navigation';
import { analytics } from '../../core/analytics';
import type { AnalyticsConfig } from '../../core/types';

export { analytics, useAnalytics, AnalyticsProvider } from '../react';

/**
 * Next.js App Router hook for automatic page tracking
 * Tracks route changes in Next.js 13+ App Router
 */
export function useNextAnalytics() {
    const pathname = usePathname();
    const searchParams = useSearchParams();
    const initialized = useRef(false);

    useEffect(() => {
        if (initialized.current) {
            analytics.pageView(pathname + (searchParams?.toString() ? `?${searchParams.toString()}` : ''));
        } else {
            initialized.current = true;
        }
    }, [pathname, searchParams]);
}

/**
 * Next.js Pages Router hook for automatic page tracking
 * Tracks route changes in Next.js Pages Router (legacy)
 */
export function useNextPagesAnalytics() {
    const router = useRouter();
    const initialized = useRef(false);

    useEffect(() => {
        const handleRouteChange = (url: string) => {
            analytics.pageView(url);
        };

        if (!initialized.current) {
            // Track initial page
            analytics.pageView(router.pathname);
            initialized.current = true;
        }

        router.events?.on('routeChangeComplete', handleRouteChange);

        return () => {
            router.events?.off('routeChangeComplete', handleRouteChange);
        };
    }, [router]);
}

/**
 * Initialize analytics in Next.js App Router layout
 */
export function initNextAnalytics(config: AnalyticsConfig) {
    if (typeof window !== 'undefined') {
        Object.assign((analytics as any).config, config);
    }
}

/**
 * Track events with automatic client-side only execution
 */
export function trackEvent(eventName: string, properties?: Record<string, any>) {
    if (typeof window !== 'undefined') {
        analytics.track(eventName, properties);
    }
}
