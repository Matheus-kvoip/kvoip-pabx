import { appConfig } from './config';

type RequestOptions = RequestInit & {
  token?: string;
};

export async function apiRequest<T>(
  path: string,
  init?: RequestOptions,
  getToken?: () => Promise<string | undefined> | string | undefined,
): Promise<T> {
  const explicit = init?.token;
  const token =
    explicit ??
    (getToken ? await getToken() : undefined);

  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(init?.headers ?? {}),
  };

  if (token) {
    (headers as Record<string, string>).Authorization = `Bearer ${token}`;
  }

  const { token: _ignored, ...rest } = init ?? {};
  void _ignored;

  const response = await fetch(`${appConfig.apiUrl}${path}`, {
    ...rest,
    headers,
    cache: 'no-store',
  });

  if (!response.ok) {
    let message = `Falha na API (${response.status})`;
    try {
      const data = (await response.json()) as { message?: string | string[] };
      if (Array.isArray(data.message)) {
        message = data.message.join(', ');
      } else if (data.message) {
        message = data.message;
      }
    } catch {
      const text = await response.text().catch(() => '');
      if (text) message = text;
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return response.json() as Promise<T>;
}
