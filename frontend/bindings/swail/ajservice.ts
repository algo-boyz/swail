// Hand-written bindings for main.AJService (Call.ByName — no ID hash required).
// Regenerate with `wails3 generate bindings -ts -names` when convenient.

import { Call as $Call, CancellablePromise as $CancellablePromise } from '@wailsio/runtime';

export interface QueryArgs {
	tool: string;
	filter: string;
	time_window_sec: number;
	limit: number;
	camera_id?: string;
}

export interface QueryResultRow {
	title: string;
	meta: string;
	confidence: number;
	object_id?: string;
	kind: string;
	raw?: string;
}

export interface MockAlert {
	id: string;
	kind: string;
	camera_id: string;
	message: string;
	confidence: number;
	timestamp: string;
	object_id?: string;
}

export interface MockCamera {
	id: string;
	label: string;
	status: string;
	scene: string;
}

export function Health(): $CancellablePromise<Record<string, unknown>> {
	return $Call.ByName('main.AJService.Health') as $CancellablePromise<Record<string, unknown>>;
}

export function Query(args: QueryArgs): $CancellablePromise<QueryResultRow[] | null> {
	return $Call.ByName('main.AJService.Query', args) as $CancellablePromise<QueryResultRow[] | null>;
}

export function GetRecentAlerts(): $CancellablePromise<MockAlert[] | null> {
	return $Call.ByName('main.AJService.GetRecentAlerts') as $CancellablePromise<MockAlert[] | null>;
}

export function GetCameras(): $CancellablePromise<MockCamera[] | null> {
	return $Call.ByName('main.AJService.GetCameras') as $CancellablePromise<MockCamera[] | null>;
}
