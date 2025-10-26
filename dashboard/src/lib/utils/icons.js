/**
 * Get country flag emoji from country code or name
 * @param {string} country - Country name or code
 * @returns {string} Flag emoji
 */
export function getCountryFlag(country) {
    if (!country) return '🌍';

    const countryFlags = {
        // Common countries
        'United States': '🇺🇸',
        'US': '🇺🇸',
        'USA': '🇺🇸',
        'United Kingdom': '🇬🇧',
        'UK': '🇬🇧',
        'GB': '🇬🇧',
        'Canada': '🇨🇦',
        'CA': '🇨🇦',
        'Germany': '🇩🇪',
        'DE': '🇩🇪',
        'France': '🇫🇷',
        'FR': '🇫🇷',
        'Italy': '🇮🇹',
        'IT': '🇮🇹',
        'Spain': '🇪🇸',
        'ES': '🇪🇸',
        'Netherlands': '🇳🇱',
        'NL': '🇳🇱',
        'Belgium': '🇧🇪',
        'BE': '🇧🇪',
        'Switzerland': '🇨🇭',
        'CH': '🇨🇭',
        'Austria': '🇦🇹',
        'AT': '🇦🇹',
        'Sweden': '🇸🇪',
        'SE': '🇸🇪',
        'Norway': '🇳🇴',
        'NO': '🇳🇴',
        'Denmark': '🇩🇰',
        'DK': '🇩🇰',
        'Finland': '🇫🇮',
        'FI': '🇫🇮',
        'Poland': '🇵🇱',
        'PL': '🇵🇱',
        'Portugal': '🇵🇹',
        'PT': '🇵🇹',
        'Greece': '🇬🇷',
        'GR': '🇬🇷',
        'Czech Republic': '🇨🇿',
        'CZ': '🇨🇿',
        'Ireland': '🇮🇪',
        'IE': '🇮🇪',
        'Australia': '🇦🇺',
        'AU': '🇦🇺',
        'New Zealand': '🇳🇿',
        'NZ': '🇳🇿',
        'Japan': '🇯🇵',
        'JP': '🇯🇵',
        'China': '🇨🇳',
        'CN': '🇨🇳',
        'South Korea': '🇰🇷',
        'KR': '🇰🇷',
        'India': '🇮🇳',
        'IN': '🇮🇳',
        'Brazil': '🇧🇷',
        'BR': '🇧🇷',
        'Mexico': '🇲🇽',
        'MX': '🇲🇽',
        'Argentina': '🇦🇷',
        'AR': '🇦🇷',
        'Chile': '🇨🇱',
        'CL': '🇨🇱',
        'Colombia': '🇨🇴',
        'CO': '🇨🇴',
        'Russia': '🇷🇺',
        'RU': '🇷🇺',
        'Turkey': '🇹🇷',
        'TR': '🇹🇷',
        'South Africa': '🇿🇦',
        'ZA': '🇿🇦',
        'Egypt': '🇪🇬',
        'EG': '🇪🇬',
        'Saudi Arabia': '🇸🇦',
        'SA': '🇸🇦',
        'UAE': '🇦🇪',
        'United Arab Emirates': '🇦🇪',
        'AE': '🇦🇪',
        'Singapore': '🇸🇬',
        'SG': '🇸🇬',
        'Malaysia': '🇲🇾',
        'MY': '🇲🇾',
        'Thailand': '🇹🇭',
        'TH': '🇹🇭',
        'Indonesia': '🇮🇩',
        'ID': '🇮🇩',
        'Philippines': '🇵🇭',
        'PH': '🇵🇭',
        'Vietnam': '🇻🇳',
        'VN': '🇻🇳',
        'Palestine': '🇵🇸',
        'PS': '🇵🇸',
        'Israel': '🇵🇸',
        'IL': '🇵🇸',
        'Ukraine': '🇺🇦',
        'UA': '🇺🇦',
        'Romania': '🇷🇴',
        'RO': '🇷🇴',
        'Hungary': '🇭🇺',
        'HU': '🇭🇺',
        'Direct': '🌍',
        'Unknown': '🌍'
    };

    return countryFlags[country] || '🌍';
}

// Import browser icons
import chromeIcon from '$lib/assets/icons/browser/chrome.svg';
import chromiumIcon from '$lib/assets/icons/browser/chromium.svg';
import firefoxIcon from '$lib/assets/icons/browser/firefox.svg';
import safariIcon from '$lib/assets/icons/browser/safari.png';
import edgeIcon from '$lib/assets/icons/browser/edge.svg';
import operaIcon from '$lib/assets/icons/browser/opera.svg';
import vivaldiIcon from '$lib/assets/icons/browser/vivaldi.svg';
import samsungIcon from '$lib/assets/icons/browser/samsung-internet.svg';
import duckduckgoIcon from '$lib/assets/icons/browser/duckduckgo.svg';
import ucIcon from '$lib/assets/icons/browser/uc.svg';
import curlIcon from '$lib/assets/icons/browser/curl.svg';
import browserFallbackIcon from '$lib/assets/icons/browser/fallback.svg';

// Import OS icons
import windowsIcon from '$lib/assets/icons/os/windows.png';
import macIcon from '$lib/assets/icons/os/mac.png';
import linuxIcon from '$lib/assets/icons/os/gnu_linux.png';
import ubuntuIcon from '$lib/assets/icons/os/ubuntu.png';
import fedoraIcon from '$lib/assets/icons/os/fedora.png';
import androidIcon from '$lib/assets/icons/os/android.png';
import iosIcon from '$lib/assets/icons/os/ios.png';
import ipadIcon from '$lib/assets/icons/os/ipad_os.png';
import chromeOsIcon from '$lib/assets/icons/os/chrome_os.png';
import harmonyIcon from '$lib/assets/icons/os/harmony_os.png';
import tizenIcon from '$lib/assets/icons/os/tizen.png';
import kaiIcon from '$lib/assets/icons/os/kai_os.png';
import fireOsIcon from '$lib/assets/icons/os/fire_os.png';
import freebsdIcon from '$lib/assets/icons/os/freebsd.png';
import playstationIcon from '$lib/assets/icons/os/playstation.png';
import osFallbackIcon from '$lib/assets/icons/os/fallback.svg';

/**
 * Get browser icon name from browser name
 * @param {string} browser - Browser name
 * @returns {string} Browser name for icon
 */
export function getBrowserName(browser) {
    if (!browser) return 'unknown';

    const lowerBrowser = browser.toLowerCase();

    if (lowerBrowser.includes('chrome')) return 'chrome';
    if (lowerBrowser.includes('firefox')) return 'firefox';
    if (lowerBrowser.includes('safari')) return 'safari';
    if (lowerBrowser.includes('edge')) return 'edge';
    if (lowerBrowser.includes('opera')) return 'opera';
    if (lowerBrowser.includes('brave')) return 'brave';
    if (lowerBrowser.includes('chromium')) return 'chromium';
    if (lowerBrowser.includes('vivaldi')) return 'vivaldi';
    if (lowerBrowser.includes('samsung')) return 'samsung-internet';
    if (lowerBrowser.includes('duckduckgo')) return 'duckduckgo';
    if (lowerBrowser.includes('uc')) return 'uc';
    if (lowerBrowser.includes('curl')) return 'curl';

    return 'unknown';
}

/**
 * Get browser icon from browser name
 * @param {string} browser - Browser name
 * @returns {string} Browser icon import or emoji fallback
 */
export function getBrowserIcon(browser) {
    const browserName = getBrowserName(browser);

    const browserIcons = {
        chrome: chromeIcon,
        chromium: chromiumIcon,
        firefox: firefoxIcon,
        safari: safariIcon,
        edge: edgeIcon,
        opera: operaIcon,
        vivaldi: vivaldiIcon,
        'samsung-internet': samsungIcon,
        duckduckgo: duckduckgoIcon,
        uc: ucIcon,
        curl: curlIcon,
        brave: chromeIcon, // Use chrome icon as fallback for brave
    };

    // Return icon if found, otherwise return emoji or fallback
    if (browserIcons[browserName]) {
        return browserIcons[browserName];
    }

    // Browser emoji fallbacks
    const browserEmojis = {
        unknown: '🌍'
    };

    return browserEmojis[browserName] || browserFallbackIcon;
}

/**
 * Get device icon from device type
 * @param {string} device - Device type
 * @returns {string} Device icon emoji or fallback
 */
export function getDeviceIcon(device) {
    const deviceLower = (device || '').toLowerCase();

    // Device emoji fallbacks
    if (deviceLower.includes('desktop')) return '🖥️';
    if (deviceLower.includes('mobile') || deviceLower.includes('phone')) return '📱';
    if (deviceLower.includes('tablet')) return '📱';

    return '💻'; // Default device emoji
}

/**
 * Get OS icon from OS name
 * @param {string} os - OS name
 * @returns {string} OS icon import or emoji fallback
 */
export function getOSIcon(os) {
    const osLower = (os || '').toLowerCase();

    // First try to match with imported icons
    if (osLower.includes('windows')) return windowsIcon;
    if (osLower.includes('mac') || osLower.includes('darwin')) return macIcon;
    if (osLower.includes('ubuntu')) return ubuntuIcon;
    if (osLower.includes('fedora')) return fedoraIcon;
    if (osLower.includes('linux')) return linuxIcon;
    if (osLower.includes('android')) return androidIcon;
    if (osLower.includes('ipad')) return ipadIcon;
    if (osLower.includes('ios') || osLower.includes('iphone')) return iosIcon;
    if (osLower.includes('chrome') && osLower.includes('os')) return chromeOsIcon;
    if (osLower.includes('harmony')) return harmonyIcon;
    if (osLower.includes('tizen')) return tizenIcon;
    if (osLower.includes('kai')) return kaiIcon;
    if (osLower.includes('fire')) return fireOsIcon;
    if (osLower.includes('freebsd')) return freebsdIcon;
    if (osLower.includes('playstation')) return playstationIcon;

    // OS emoji fallbacks for uncommon systems
    const osEmojis = {
        'bsd': '😈',
        'solaris': '☀️',
        'unix': '🖥️',
    };

    for (const [key, emoji] of Object.entries(osEmojis)) {
        if (osLower.includes(key)) return emoji;
    }

    return osFallbackIcon;
}/**
 * Get favicon URL from source/referrer
 * @param {string} source - Source URL or name
 * @returns {string|null} Favicon URL
 */
export function getFaviconUrl(source) {
    if (!source || source === 'Direct') return null;

    try {
        // Try to parse as URL
        const url = new URL(source.startsWith('http') ? source : `https://${source}`);
        return `https://www.google.com/s2/favicons?domain=${url.hostname}&sz=32`;
    } catch {
        // If not a valid URL, return null
        return null;
    }
}

/**
 * Get source display name from URL
 * @param {string} source - Source URL
 * @returns {string} Display name
 */
export function getSourceDisplayName(source) {
    if (!source || source === 'Direct') return 'Direct';

    try {
        const url = new URL(source.startsWith('http') ? source : `https://${source}`);
        return url.hostname.replace('www.', '');
    } catch {
        return source;
    }
}

