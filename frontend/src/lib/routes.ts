import Home from '../pages/Home.svelte';
import Query from '../pages/Query.svelte';
import Cameras from '../pages/Cameras.svelte';
import Alerts from '../pages/Alerts.svelte';
import Evidence from '../pages/Evidence.svelte';

export const routes = {
	'/': Home,
	'/query': Query,
	'/cameras': Cameras,
	'/alerts': Alerts,
	'/evidence/:id': Evidence,
	// catch-all
	'*': Home
};
