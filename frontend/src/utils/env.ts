/**
 * Copyright 2022 Redpanda Data, Inc.
 *
 * Use of this software is governed by the Business Source License
 * included in the file https://github.com/redpanda-data/redpanda/blob/dev/licenses/bsl.md
 *
 * As of the Change Date specified in that file, in accordance with
 * the Business Source License, use of this software will be governed
 * by the Apache License, Version 2.0
 */

import { toJson } from "./jsonUtils";
import './extensions';

// const assertions and mapped types are awesome!
// type
// { [key: string]: string }

const envNames = [
    'NODE_ENV',

    'REACT_APP_CONSOLE_GIT_SHA',
    'REACT_APP_CONSOLE_GIT_REF', // 'master' or 'v1.2.3'
    'REACT_APP_BUILD_TIMESTAMP',

    'REACT_APP_CONSOLE_BUSINESS_GIT_SHA',
    'REACT_APP_CONSOLE_BUSINESS_GIT_REF',

    'REACT_APP_BUILT_FROM_PUSH', // was built by 'image-on-push'?
    'REACT_APP_BUSINESS', // is business version?

    'REACT_APP_DEV_HINT', // for debugging, since we can't override NODE_ENV
] as const;

type Environment = { [key in typeof envNames[number]]: string };

const env = {} as Environment;
for (const k of envNames)
    env[k] = process.env[k] || '';

export default env;


//
// Helpers
const isDev = (!process.env.NODE_ENV || process.env.NODE_ENV === 'development') || process.env.REACT_APP_DEV_HINT;
export const IsProd = !isDev; // todo: rename IsProd IsDev
export const IsDev = isDev;

export const IsBusiness = Boolean(env.REACT_APP_BUSINESS);
export const AppName = IsBusiness ? 'Redpanda Console Business' : 'Redpanda Console';

const basePathRaw: string = (window as any)["BASE_URL"];
const basePath = (typeof basePathRaw === 'string' && !basePathRaw.startsWith('__BASE_PATH'))
    ? basePathRaw
    : '';
export const basePathNo = basePath ? basePath.removePrefix('/').removeSuffix('/') : '';
export const basePathS = basePathNo ? '/' + basePathNo : '';
export const basePathE = basePathNo ? basePathNo + '/' : '';


export function getBuildDate(): Date | undefined {
    const timestamp = +env.REACT_APP_BUILD_TIMESTAMP;
    if (timestamp == 0) return undefined;
    return new Date(timestamp * 1000);
}


//
// Print all env vars to console
const envVarDebugObj = {} as any;
const envVarDebugAr: { name: string, value: string }[] = [];

const addProp = (key: string, value: any) => {
    if (value === undefined || value === null || value === "") return;
    key = key.removePrefix("REACT_APP_CONSOLE_").removePrefix("REACT_APP_");
    envVarDebugObj[key] = value;
    envVarDebugAr.push({ name: key, value: value });
}
// - add env vars
for (const k in env) addProp(k, (env as any)[k]);

// - add custom
addProp("appName", AppName);

// - print
console.log(toJson(envVarDebugObj));

export { envVarDebugObj, envVarDebugAr };