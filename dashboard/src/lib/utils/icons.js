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
		'Israel': '🇮🇱',
		'IL': '🇮🇱',
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
	
	return 'unknown';
}

/**
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

