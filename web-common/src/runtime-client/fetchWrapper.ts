import { fetchTokenFromLocalStorage } from "@rilldata/web-common/runtime-client/utils";

export type FetchWrapperOptions = {
  baseUrl?: string;
  url: string;
  method: string;
  headers?: HeadersInit;
  params?: Record<string, unknown>;
  data?: any;
  signal?: AbortSignal;
};

export interface HTTPError {
  response: {
    status: number;
    data: {
      message: string;
    };
  };
  message: string;
}

export async function fetchWrapper({
  url,
  method,
  headers,
  data,
  params,
  signal,
}: FetchWrapperOptions) {
  if (signal && signal.aborted) return Promise.reject(new Error("Aborted"));

  headers ??= { "Content-Type": "application/json" };
  const TOKEN = fetchTokenFromLocalStorage();
  if(TOKEN) {
    headers["Authorization"] = `Bearer ${TOKEN}`;
  }
  url = encodeURI(url);

  if (params) {
    const paramParts = [];
    for (const p in params) {
      paramParts.push(`${p}=${encodeURIComponent(params[p] as string)}`);
    }
    if (paramParts.length) {
      url = `${url}?${paramParts.join("&")}`;
    }
  }

  const resp = await fetch(url, {
    method,
    ...(data ? { body: serializeBody(data) } : {}),
    headers,
    signal,
  });
  if (!resp.ok) {
    const data = await resp.json();

    // Return runtime errors in the same form as the Axios client had previously
    if (data.code && data.message) {
      return Promise.reject({
        response: {
          status: resp.status,
          data,
        },
      });
    }

    // Fallback error handling
    const err = new Error();
    (err as any).response = await resp.json();
    return Promise.reject(err);
  }
  return resp.json();
}

function serializeBody(body: BodyInit | Record<string, unknown>): BodyInit {
  return body instanceof FormData ? body : JSON.stringify(body);
}
