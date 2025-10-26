import { watch, ref } from 'vue';
import { useRoute } from 'vue-router';
import { analytics } from '../../core/analytics';
import type { AnalyticsConfig } from '../../core/types';

export { analytics, useAnalytics, AnalyticsPlugin } from '../vue';

/**
 * Nuxt 3 composable for automatic page tracking
 * Tracks route changes in Nuxt 3
 */
export function useNuxtAnalytics() {
    const route = useRoute();
    const isFirstRender = ref(true);

    watch(
        () => route.fullPath,
        (newPath) => {
            if (!isFirstRender.value) {
                analytics.pageView(newPath);
            } else {
                isFirstRender.value = false;
            }
        },
        { immediate: true }
    );
}

/**
 * Initialize analytics in Nuxt plugin
 */
export function initNuxtAnalytics(config: AnalyticsConfig) {
    if (process.client) {
        Object.assign((analytics as any).config, config);
    }
}

/**
 * Nuxt 3 plugin definition
 */
export default defineNuxtPlugin(() => {
    // Analytics will be initialized via plugin options
    // This provides the analytics instance globally
    return {
        provide: {
            analytics: {
                track: (eventName: string, properties?: Record<string, any>) => {
                    if (process.client) {
                        analytics.track(eventName, properties);
                    }
                },
                pageView: (url?: string | null) => {
                    if (process.client) {
                        analytics.pageView(url);
                    }
                },
                identify: (userId: string, traits?: Record<string, any>) => {
                    if (process.client) {
                        analytics.identify(userId, traits);
                    }
                },
                trackClick: (elementId: string, properties?: Record<string, any>) => {
                    if (process.client) {
                        analytics.trackClick(elementId, properties);
                    }
                },
                trackForm: (formId: string, properties?: Record<string, any>) => {
                    if (process.client) {
                        analytics.trackForm(formId, properties);
                    }
                },
                trackError: (error: Error | string, context?: Record<string, any>) => {
                    if (process.client) {
                        analytics.trackError(error, context);
                    }
                },
            },
        },
    };
});

/**
 * Track events with automatic client-side only execution
 */
export function trackEvent(eventName: string, properties?: Record<string, any>) {
    if (process.client) {
        analytics.track(eventName, properties);
    }
}

// Type augmentation for Nuxt
declare module '#app' {
    interface NuxtApp {
        $analytics: {
            track: (eventName: string, properties?: Record<string, any>) => void;
            pageView: (url?: string | null) => void;
            identify: (userId: string, traits?: Record<string, any>) => void;
            trackClick: (elementId: string, properties?: Record<string, any>) => void;
            trackForm: (formId: string, properties?: Record<string, any>) => void;
            trackError: (error: Error | string, context?: Record<string, any>) => void;
        };
    }
}

declare module 'vue' {
    interface ComponentCustomProperties {
        $analytics: {
            track: (eventName: string, properties?: Record<string, any>) => void;
            pageView: (url?: string | null) => void;
            identify: (userId: string, traits?: Record<string, any>) => void;
            trackClick: (elementId: string, properties?: Record<string, any>) => void;
            trackForm: (formId: string, properties?: Record<string, any>) => void;
            trackError: (error: Error | string, context?: Record<string, any>) => void;
        };
    }
}
