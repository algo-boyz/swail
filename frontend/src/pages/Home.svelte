<script lang="ts">
	import { push } from 'svelte-spa-router';
	import { _ } from '../lib/i18n';
	import LangSwitcher from '../lib/LangSwitcher.svelte';

	const cards = $derived([
		{ path: '/query', title: $_('nav.query'), desc: $_('nav.query_desc') },
		{ path: '/cameras', title: $_('nav.cameras'), desc: $_('nav.cameras_desc') },
		{ path: '/alerts', title: $_('nav.alerts'), desc: $_('nav.alerts_desc') },
		{ path: '/evidence/demo', title: $_('nav.evidence'), desc: $_('nav.evidence_desc') }
	]);
</script>

<div class="home">
	<header class="top-bar">
		<span class="brand-mark" data-wml-openURL="https://www.agentjetson.ai">
			<img src="/wails.png" class="brand-logo" alt="AgentJetson" />
		</span>
		<LangSwitcher />
	</header>

	<span class="badge">{$_('home.badge')}</span>

	<h1 class="title">
		<span class="title-accent">AgentJetson</span>
	</h1>
	<p class="subtitle">
		{$_('home.subtitle')}
	</p>

	<nav class="platforms" aria-label="Main sections">
		{#each cards as c}
			<button class="platform-card" onclick={() => push(c.path)} aria-label={c.title}>
				<span>
					<span class="platform-name">{c.title}</span>
					<span class="card-desc">{c.desc}</span>
				</span>
				<span class="platform-arrow">
					<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2"
						stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
						<line x1="5" y1="12" x2="19" y2="12"></line>
						<polyline points="12 5 19 12 12 19"></polyline>
					</svg>
				</span>
			</button>
		{/each}
	</nav>
</div>

<style>
	.home {
		width: 100%;
	}

	.top-bar {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 0.75rem;
		margin-bottom: var(--s-2);
		flex-shrink: 0;
	}

	.top-bar :global(.brand-logo) {
		height: calc(var(--s-4) * 0.85);
	}

	.card-desc {
		display: block;
		font-size: 12.5px;
		font-weight: 500;
		color: var(--muted);
		margin-top: 3px;
	}
	.platform-card {
		align-items: flex-start;
	}
</style>
