import type { AnalyticsConfig, EventData, SessionInfo, DeviceInfo } from './types';

/**
 * Analytics JavaScript SDK
 * Lightweight client-side analytics tracking library
 */
class AnalyticsCore {
    private config: Required<AnalyticsConfig>;
    private buffer: EventData[] = [];
    private sessionId: string;
    private userId: string;
    private flushTimer: number | null = null;
    private sessionStartTime: number;

    constructor(config: AnalyticsConfig = {}) {
        this.config = {
            apiUrl: config.apiUrl || 'http://localhost:8080',
            projectId: config.projectId || 'default',
            autoTrack: config.autoTrack !== false,
            bufferSize: config.bufferSize || 10,
            flushInterval: config.flushInterval || 30000, // 30 seconds
            debug: config.debug || false,
            timeout: config.timeout || 10000,
            maxRetries: config.maxRetries || 3,
        };

        this.sessionId = this.getSessionId();
        this.userId = this.getUserId();
        this.sessionStartTime = Date.now();

        if (typeof window !== 'undefined') {
            if (this.config.autoTrack) {
                this.setupAutoTracking();
            }

            if (this.config.flushInterval > 0) {
                this.startAutoFlush();
            }

            // Flush on page unload
            window.addEventListener('beforeunload', () => {
                this.flush(true);
            });
        }

        this.log('Analytics initialized', this.config);
    }

    /**
     * Track a custom event
     */
    track(eventName: string, properties: Record<string, any> = {}): void {
        if (typeof window === 'undefined') return;

        const event: EventData = {
            event_name: eventName,
            user_id: this.userId,
            session_id: this.sessionId,
            session_duration: Math.floor((Date.now() - this.sessionStartTime) / 1000), // Duration in seconds
            url: window.location.href,
            referrer: document.referrer,
            user_agent: navigator.userAgent,
            timestamp: new Date().toISOString(),
            browser: this.getBrowser(),
            os: this.getOS(),
            device: this.getDevice(),
            project_id: this.config.projectId,
        };

        this.log('Tracking event:', event);
        this.addToBuffer(event);
    }

    /**
     * Track a page view
     */
    pageView(url: string | null = null): void {
        if (typeof window === 'undefined') return;

        this.track('page_view', {
            path: url || window.location.pathname,
            title: document.title,
            search: window.location.search,
            hash: window.location.hash,
        });
        this.flush();
    }

    /**
     * Identify a user
     */
    identify(userId: string, traits: Record<string, any> = {}): void {
        this.userId = userId;
        if (typeof window !== 'undefined') {
            localStorage.setItem('analytics_user_id', userId);
        }
        this.track('identify', traits);
    }

    /**
     * Track a click event
     */
    trackClick(elementId: string, properties: Record<string, any> = {}): void {
        this.track('click', {
            element: elementId,
            ...properties,
        });
    }

    /**
     * Track form submission
     */
    trackForm(formId: string, properties: Record<string, any> = {}): void {
        this.track('form_submit', {
            form_id: formId,
            ...properties,
        });
    }

    /**
     * Track an error
     */
    trackError(error: Error | string, context: Record<string, any> = {}): void {
        const errorData = typeof error === 'string'
            ? { message: error }
            : { message: error.message, stack: error.stack };

        this.track('error', {
            ...errorData,
            ...context,
        });
    }

    /**
     * Set user properties
     */
    setUserProperties(properties: Record<string, any>): void {
        this.track('user_properties_updated', properties);
    }

    /**
     * Add event to buffer
     */
    private addToBuffer(event: EventData): void {
        this.buffer.push(event);

        if (this.buffer.length >= this.config.bufferSize) {
            this.flush();
        }
    }

    /**
     * Flush buffered events to the server
     */
    flush(async = false): Promise<any> {
        if (this.buffer.length === 0) return Promise.resolve();
        if (typeof window === 'undefined') return Promise.resolve();

        const events = [...this.buffer];
        this.buffer = [];

        this.log('Flushing events:', events);

        if (async && navigator.sendBeacon) {
            // Use sendBeacon for reliable delivery on page unload
            return Promise.all(
                events.map((event) => {
                    const blob = new Blob([JSON.stringify(event)], { type: 'application/json' });
                    return navigator.sendBeacon(`${this.config.apiUrl}/api/track`, blob);
                })
            );
        } else {
            // Use fetch for normal requests with better error handling and retries
            return Promise.all(events.map((event) => this.sendEvent(event)));
        }
    }

    /**
     * Send individual event with retry logic
     */
    private sendEvent(event: EventData, retries: number = this.config.maxRetries): Promise<any> {
        if (typeof window === 'undefined') return Promise.resolve();

        return fetch(`${this.config.apiUrl}/api/track`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(event),
            keepalive: true,
            signal: AbortSignal.timeout(this.config.timeout),
        })
            .then((response) => {
                if (!response.ok) {
                    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
                }
                this.log('Event sent successfully:', event.event_name);
                return response;
            })
            .catch((err) => {
                this.log('Error sending event:', err, event);

                // Retry on network errors
                if (retries > 0 && (err.name === 'NetworkError' || err.name === 'TypeError')) {
                    this.log(`Retrying event send, ${retries} attempts left`);
                    return new Promise((resolve) => {
                        setTimeout(() => {
                            resolve(this.sendEvent(event, retries - 1));
                        }, 1000 * (this.config.maxRetries + 1 - retries)); // Exponential backoff
                    });
                }

                // If all retries failed, add back to buffer for next flush
                if (retries === 0) {
                    this.buffer.unshift(event);
                    this.log('Event added back to buffer after failed retries');
                }

                throw err;
            });
    }

    /**
     * Setup automatic tracking
     */
    private setupAutoTracking(): void {
        if (typeof window === 'undefined') return;

        // Track initial page view
        if (document.readyState === 'complete') {
            this.pageView();
        } else {
            window.addEventListener('load', () => this.pageView());
        }

        // Track page visibility changes
        document.addEventListener('visibilitychange', () => {
            this.track(document.hidden ? 'page_hidden' : 'page_visible');
        });

        // Track clicks on links
        document.addEventListener('click', (e) => {
            const link = (e.target as HTMLElement).closest('a');
            if (link) {
                this.track('link_clicked', {
                    url: link.href,
                    text: link.textContent?.trim(),
                    external: link.hostname !== window.location.hostname,
                });
            }
        });

        // Track form submissions
        document.addEventListener('submit', (e) => {
            const form = e.target as HTMLFormElement;
            if (form.id) {
                this.trackForm(form.id);
            }
        });

        // Track errors
        window.addEventListener('error', (e) => {
            this.trackError(e.error || e.message, {
                filename: e.filename,
                lineno: e.lineno,
                colno: e.colno,
            });
        });

        // Track unhandled promise rejections
        window.addEventListener('unhandledrejection', (e) => {
            this.trackError(e.reason, {
                type: 'unhandled_promise_rejection',
            });
        });
    }

    /**
     * Start auto-flush timer
     */
    private startAutoFlush(): void {
        if (typeof window === 'undefined') return;

        this.flushTimer = window.setInterval(() => {
            this.flush();
        }, this.config.flushInterval);
    }

    /**
     * Get or create session ID
     */
    private getSessionId(): string {
        if (typeof window === 'undefined') return this.generateId();

        let sessionId = sessionStorage.getItem('analytics_session_id');
        const sessionStart = sessionStorage.getItem('analytics_session_start');

        // Create new session if doesn't exist or if it's been more than 30 minutes
        if (!sessionId || !sessionStart || (Date.now() - parseInt(sessionStart)) > 30 * 60 * 1000) {
            sessionId = this.generateId();
            sessionStorage.setItem('analytics_session_id', sessionId);
            sessionStorage.setItem('analytics_session_start', Date.now().toString());
            this.sessionStartTime = Date.now();
        } else {
            this.sessionStartTime = parseInt(sessionStart);
        }

        return sessionId;
    }

    /**
     * Get or create user ID
     */
    private getUserId(): string {
        if (typeof window === 'undefined') return this.generateId();

        let userId = localStorage.getItem('analytics_user_id');
        if (!userId) {
            userId = this.generateId();
            localStorage.setItem('analytics_user_id', userId);
        }
        return userId;
    }

    /**
     * Generate a unique ID
     */
    private generateId(): string {
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
            const r = (Math.random() * 16) | 0;
            const v = c === 'x' ? r : (r & 0x3) | 0x8;
            return v.toString(16);
        });
    }

    /**
     * Get browser name
     */
    private getBrowser(): string {
        if (typeof window === 'undefined') return 'Unknown';

        const ua = navigator.userAgent;
        if (ua.indexOf('Firefox') > -1) return 'Firefox';
        if (ua.indexOf('Chrome') > -1) return 'Chrome';
        if (ua.indexOf('Safari') > -1) return 'Safari';
        if (ua.indexOf('Edge') > -1) return 'Edge';
        if (ua.indexOf('Opera') > -1 || ua.indexOf('OPR') > -1) return 'Opera';
        if (ua.indexOf('MSIE') > -1 || ua.indexOf('Trident') > -1) return 'IE';
        return 'Unknown';
    }

    /**
     * Get operating system
     */
    private getOS(): string {
        if (typeof window === 'undefined') return 'Unknown';

        const ua = navigator.userAgent;
        if (ua.indexOf('Win') > -1) return 'Windows';
        if (ua.indexOf('Mac') > -1) return 'MacOS';
        if (ua.indexOf('Linux') > -1) return 'Linux';
        if (ua.indexOf('Android') > -1) return 'Android';
        if (ua.indexOf('iOS') > -1 || ua.indexOf('iPhone') > -1 || ua.indexOf('iPad') > -1) return 'iOS';
        return 'Unknown';
    }

    /**
     * Get device type
     */
    private getDevice(): string {
        if (typeof window === 'undefined') return 'Desktop';

        const ua = navigator.userAgent;
        if (/(tablet|ipad|playbook|silk)|(android(?!.*mobi))/i.test(ua)) {
            return 'Tablet';
        }
        if (/Mobile|Android|iP(hone|od)|IEMobile|BlackBerry|Kindle|Silk-Accelerated|(hpw|web)OS|Opera M(obi|ini)/.test(ua)) {
            return 'Mobile';
        }
        return 'Desktop';
    }

    /**
     * Debug logging
     */
    private log(...args: any[]): void {
        if (this.config.debug) {
            console.log('[Analytics]', ...args);
        }
    }

    /**
     * Destroy the analytics instance
     */
    destroy(): void {
        if (this.flushTimer) {
            clearInterval(this.flushTimer);
        }
        this.flush(true);
    }
}

export const analytics = new AnalyticsCore();
export type { AnalyticsConfig, EventData, SessionInfo, DeviceInfo } from './types';
