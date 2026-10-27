import Home from '../pages/Home.svelte';
import Query from '../pages/Query.svelte';
import Cameras from '../pages/Cameras.svelte';
import Alerts from '../pages/Alerts.svelte';
import Evidence from '../pages/Evidence.svelte';
import Map from '../pages/Map.svelte';

export const routes = {
	'/': Home,
	'/query': Query,
	'/cameras': Cameras,
	'/alerts': Alerts,
	'/evidence/:id': Evidence,
	'/map': Map,
	// catch-all
	'*': Home
};
