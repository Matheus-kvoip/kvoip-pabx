export const appConfig = {
  apiUrl: process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:3001/api',
  appName: process.env.NEXT_PUBLIC_APP_NAME ?? 'KVOIP',
  appTagline: process.env.NEXT_PUBLIC_APP_TAGLINE ?? 'PABX Virtual',
  apiHostLabel:
    process.env.NEXT_PUBLIC_API_HOST_LABEL ?? 'localhost:3001',
  authCookie:
    process.env.NEXT_PUBLIC_AUTH_COOKIE ??
    process.env.AUTH_COOKIE ??
    'kvoip_token',
};
