import { appConfig } from './config';

export function getAuthCookieName() {
  return appConfig.authCookie;
}

export function setClientToken(token: string) {
  const maxAge = 60 * 60 * 8;
  document.cookie = `${getAuthCookieName()}=${encodeURIComponent(token)}; Path=/; Max-Age=${maxAge}; SameSite=Lax`;
}

export function clearClientToken() {
  document.cookie = `${getAuthCookieName()}=; Path=/; Max-Age=0; SameSite=Lax`;
}

export function getClientToken(): string | undefined {
  if (typeof document === 'undefined') return undefined;
  const match = document.cookie
    .split('; ')
    .find((row) => row.startsWith(`${getAuthCookieName()}=`));
  if (!match) return undefined;
  return decodeURIComponent(match.split('=').slice(1).join('='));
}
