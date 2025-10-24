/**
 * Analytics JavaScript SDK
 * Lightweight client-side analytics tracking library
 */

(function(global, factory) {
  typeof exports === 'object' && typeof module !== 'undefined'
    ? module.exports = factory()
    : typeof define === 'function' && define.amd
      ? define(factory)
      : (global.Analytics = factory());
}(this, (function() {
  'use strict';

  class Analytics {
    constructor(config = {}) {
      this.config = {
        apiUrl: config.apiUrl || 'http://localhost:8080',
        autoTrack: config.autoTrack !== false,
        bufferSize: config.bufferSize || 10,
        flushInterval: config.flushInterval || 30000, // 30 seconds
        debug: config.debug || false,
        ...config
      };

      this.buffer = [];
      this.sessionId = this._getSessionId();
      this.userId = this._getUserId();
      this.flushTimer = null;

      if (this.config.autoTrack) {
        this._setupAutoTracking();
      }

      if (this.config.flushInterval > 0) {
        this._startAutoFlush();
      }

      // Flush on page unload
      window.addEventListener('beforeunload', () => {
        this.flush(true);
      });

      this._log('Analytics initialized', this.config);
    }

    /**
     * Track a custom event
     * @param {string} eventName - Name of the event
     * @param {object} properties - Additional properties
     */
    track(eventName, properties = {}) {
      const event = {
        event_name: eventName,
        user_id: this.userId,
        session_id: this.sessionId,
        url: window.location.href,
        referrer: document.referrer,
        user_agent: navigator.userAgent,
        timestamp: new Date().toISOString(),
        browser: this._getBrowser(),
        os: this._getOS(),
        device: this._getDevice(),
        properties: JSON.stringify(properties)
      };

      this._log('Tracking event:', event);
      this._addToBuffer(event);
    }

    /**
     * Track a page view
     * @param {string} url - Optional URL override
     */
    pageView(url = null) {
      this.track('page_view', {
        path: url || window.location.pathname,
        title: document.title,
        search: window.location.search,
        hash: window.location.hash
      });
    }

    /**
     * Identify a user
     * @param {string} userId - User identifier
     * @param {object} traits - User traits/properties
     */
    identify(userId, traits = {}) {
      this.userId = userId;
      localStorage.setItem('analytics_user_id', userId);
      this.track('identify', traits);
    }

    /**
     * Track a click event
     * @param {string} elementId - ID or selector of the element
     * @param {object} properties - Additional properties
     */
    trackClick(elementId, properties = {}) {
      this.track('click', {
        element: elementId,
        ...properties
      });
    }

    /**
     * Track form submission
     * @param {string} formId - Form identifier
     * @param {object} properties - Additional properties
     */
    trackForm(formId, properties = {}) {
      this.track('form_submit', {
        form_id: formId,
        ...properties
      });
    }

    /**
     * Track an error
     * @param {Error|string} error - Error object or message
     * @param {object} context - Additional context
     */
    trackError(error, context = {}) {
      this.track('error', {
        message: error.message || error,
        stack: error.stack,
        ...context
      });
    }

    /**
     * Set user properties
     * @param {object} properties - User properties to set
     */
    setUserProperties(properties) {
      this.track('user_properties_updated', properties);
    }

    /**
     * Add event to buffer
     * @private
     */
    _addToBuffer(event) {
      this.buffer.push(event);

      if (this.buffer.length >= this.config.bufferSize) {
        this.flush();
      }
    }

    /**
     * Flush buffered events to the server
     * @param {boolean} async - Whether to use sendBeacon (for page unload)
     */
    flush(async = false) {
      if (this.buffer.length === 0) return;

      const events = [...this.buffer];
      this.buffer = [];

      this._log('Flushing events:', events);

      if (async && navigator.sendBeacon) {
        // Use sendBeacon for reliable delivery on page unload
        events.forEach(event => {
          const blob = new Blob([JSON.stringify(event)], { type: 'application/json' });
          navigator.sendBeacon(`${this.config.apiUrl}/api/track`, blob);
        });
      } else {
        // Use fetch for normal requests
        events.forEach(event => {
          fetch(`${this.config.apiUrl}/api/track`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(event),
            keepalive: true
          }).catch(err => {
            this._log('Error sending event:', err);
          });
        });
      }
    }

    /**
     * Setup automatic tracking
     * @private
     */
    _setupAutoTracking() {
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
        const link = e.target.closest('a');
        if (link) {
          this.track('link_clicked', {
            url: link.href,
            text: link.textContent.trim(),
            external: link.hostname !== window.location.hostname
          });
        }
      });

      // Track form submissions
      document.addEventListener('submit', (e) => {
        const form = e.target;
        if (form.id) {
          this.trackForm(form.id);
        }
      });

      // Track errors
      window.addEventListener('error', (e) => {
        this.trackError(e.error || e.message, {
          filename: e.filename,
          lineno: e.lineno,
          colno: e.colno
        });
      });

      // Track unhandled promise rejections
      window.addEventListener('unhandledrejection', (e) => {
        this.trackError(e.reason, {
          type: 'unhandled_promise_rejection'
        });
      });
    }

    /**
     * Start auto-flush timer
     * @private
     */
    _startAutoFlush() {
      this.flushTimer = setInterval(() => {
        this.flush();
      }, this.config.flushInterval);
    }

    /**
     * Get or create session ID
     * @private
     */
    _getSessionId() {
      let sessionId = sessionStorage.getItem('analytics_session_id');
      if (!sessionId) {
        sessionId = this._generateId();
        sessionStorage.setItem('analytics_session_id', sessionId);
      }
      return sessionId;
    }

    /**
     * Get or create user ID
     * @private
     */
    _getUserId() {
      let userId = localStorage.getItem('analytics_user_id');
      if (!userId) {
        userId = this._generateId();
        localStorage.setItem('analytics_user_id', userId);
      }
      return userId;
    }

    /**
     * Generate a unique ID
     * @private
     */
    _generateId() {
      return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        const r = Math.random() * 16 | 0;
        const v = c === 'x' ? r : (r & 0x3 | 0x8);
        return v.toString(16);
      });
    }

    /**
     * Get browser name
     * @private
     */
    _getBrowser() {
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
     * @private
     */
    _getOS() {
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
     * @private
     */
    _getDevice() {
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
     * @private
     */
    _log(...args) {
      if (this.config.debug) {
        console.log('[Analytics]', ...args);
      }
    }

    /**
     * Destroy the analytics instance
     */
    destroy() {
      if (this.flushTimer) {
        clearInterval(this.flushTimer);
      }
      this.flush(true);
    }
  }

  return Analytics;
})));
