import { cookies } from 'next/headers';
import { getAuthCookieName } from './auth';

export async function getServerToken(): Promise<string | undefined> {
  const store = await cookies();
  return store.get(getAuthCookieName())?.value;
}
