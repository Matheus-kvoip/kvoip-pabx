import Image from 'next/image';
import { appConfig } from '@/lib/config';

type BrandLogoProps = {
  priority?: boolean;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
};

const sizes = {
  sm: { width: 120, height: 36 },
  md: { width: 160, height: 48 },
  lg: { width: 240, height: 72 },
};

export function BrandLogo({
  priority = false,
  className = '',
  size = 'md',
}: BrandLogoProps) {
  const dim = sizes[size];
  return (
    <Image
      src="/logo-kvoip.png"
      alt={appConfig.appName}
      width={dim.width}
      height={dim.height}
      priority={priority}
      className={`brand-logo ${className}`.trim()}
    />
  );
}
