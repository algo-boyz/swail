import { writable, derived, get } from 'svelte/store';

import en from './locales/en.json';
import ar from './locales/ar.json';
import fr from './locales/fr.json';
import es from './locales/es.json';
import de from './locales/de.json';
import nl from './locales/nl.json';

export type Locale = 'en' | 'ar' | 'fr' | 'es' | 'de' | 'nl';

export const locales: { code: Locale; label: string; dir: 'ltr' | 'rtl' }[] = [
	{ code: 'en', label: 'English', dir: 'ltr' },
	{ code: 'ar', label: '\u0627\u0644\u0639\u0631\u0628\u064a\u0629', dir: 'rtl' },
	{ code: 'fr', label: 'Fran\u00e7ais', dir: 'ltr' },
	{ code: 'es', label: 'Espa\u00f1ol', dir: 'ltr' },
	{ code: 'de', label: 'Deutsch', dir: 'ltr' },
	{ code: 'nl', label: 'Nederlands', dir: 'ltr' }
];

const catalogs: Record<Locale, Record<string, string>> = {
	en,
	ar,
	fr,
	es,
	de,
	nl
};

const STORAGE_KEY = 'swail-locale';

function detectInitial(): Locale {
	if (typeof localStorage !== 'undefined') {
		const saved = localStorage.getItem(STORAGE_KEY) as Locale | null;
		if (saved && catalogs[saved]) return saved;
	}
	if (typeof navigator !== 'undefined') {
		const nav = (navigator.language || 'en').slice(0, 2).toLowerCase();
		if (nav in catalogs) return nav as Locale;
	}
	return 'en';
}

export const locale = writable<Locale>(detectInitial());

locale.subscribe((code) => {
	if (typeof document !== 'undefined') {
		const meta = locales.find((l) => l.code === code);
		document.documentElement.lang = code;
		document.documentElement.dir = meta?.dir ?? 'ltr';
	}
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem(STORAGE_KEY, code);
	}
});

/** Translate a key. Supports simple {param} interpolation. */
export function t(key: string, params?: Record<string, string | number>): string {
	const code = get(locale);
	const catalog = catalogs[code] ?? catalogs.en;
	let str = catalog[key] ?? catalogs.en[key] ?? key;
	if (params) {
		for (const [k, v] of Object.entries(params)) {
			str = str.replace(new RegExp(`\\{${k}\\}`, 'g'), String(v));
		}
	}
	return str;
}

/** Reactive translation — re-runs when locale changes. */
export const _ = derived(locale, () => t);

export function setLocale(code: Locale) {
	if (catalogs[code]) locale.set(code);
}
