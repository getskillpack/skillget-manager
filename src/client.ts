import { registryBaseUrl } from "./config.js";

export async function fetchJson<T>(path: string, init?: RequestInit): Promise<T> {
  const url = `${registryBaseUrl()}${path.startsWith("/") ? path : `/${path}`}`;
  const res = await fetch(url, {
    ...init,
    headers: {
      Accept: "application/json",
      ...init?.headers,
    },
  });
  if (!res.ok) {
    const text = await res.text().catch(() => "");
    throw new Error(`Registry ${res.status} ${res.statusText}: ${text || url}`);
  }
  return (await res.json()) as T;
}
