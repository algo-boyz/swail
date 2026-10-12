<script lang="ts">
	import { locale, locales, setLocale, _ } from './i18n';
	import type { Locale } from './i18n';

	function onLangChange(e: Event) {
		const code = (e.currentTarget as HTMLSelectElement).value as Locale;
		setLocale(code);
	}
</script>

<label class="lang-switch" title={$_('lang.label')}>
	<span class="lang-label">{$_('lang.label')}</span>
	<select
		class="lang-select"
		value={$locale}
		onchange={onLangChange}
		aria-label={$_('lang.label')}
	>
		{#each locales as l}
			<option value={l.code}>{l.label}</option>
		{/each}
	</select>
</label>

<style>
	.lang-switch {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		font-size: 12px;
		color: var(--muted);
		flex-shrink: 0;
		--wails-draggable: no-drag;
	}

	.lang-label {
		font-weight: 700;
		letter-spacing: 0.06em;
		text-transform: uppercase;
		font-size: 10px;
	}

	.lang-select {
		background: var(--glass);
		border: 1px solid var(--glass-border);
		border-radius: 8px;
		color: var(--text);
		padding: 5px 10px;
		font-size: 12.5px;
		font-weight: 600;
		cursor: pointer;
		outline: none;
	}

	.lang-select:focus {
		border-color: rgba(255, 77, 77, 0.55);
	}

	@media (max-width: 640px) {
		.lang-label {
			display: none;
		}
	}
</style>
