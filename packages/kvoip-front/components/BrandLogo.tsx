import Image from 'next/image';
import { appConfig } from '@/lib/config';

type BrandLogoProps = {
  priority?: boolean;
  className?: string;
  size?: 'sm' | 'md' | 'lg';
  showName?: boolean;
};

const sizes = {
  sm: 40,
  md: 64,
  lg: 128,
};

export function BrandLogo({
  priority = false,
  className = '',
  size = 'md',
  showName = false,
}: BrandLogoProps) {
  const px = sizes[size];
  return (
    <span className={`brand-mark ${className}`.trim()}>
      <Image
        src="/logo-kvoip.png"
        alt={appConfig.appName}
        width={px}
        height={px}
        priority={priority}
        className="brand-logo"
      />
      {showName ? (
        <span className={`brand-name brand-name-${size}`}>{appConfig.appName}</span>
      ) : null}
    </span>
  );
}
