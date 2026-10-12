import Home from '../pages/Home.svelte';
import Marktplaats from '../pages/Marktplaats.svelte';
import Catawiki from '../pages/Catawiki.svelte';

export const routes = {
	'/': Home,
	'/marktplaats': Marktplaats,
	'/catawiki': Catawiki,
	// catch-all
	'*': Home
};