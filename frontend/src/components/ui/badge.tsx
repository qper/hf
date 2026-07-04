import * as React from 'react'

import { cn } from '@/lib/utils'

export interface BadgeProps extends React.HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'secondary' | 'outline'
}

function Badge({ className, variant = 'default', ...props }: BadgeProps) {
  const variantClassName = {
    default: 'border-border bg-surface text-foreground',
    secondary: 'border-accent/30 bg-accent/10 text-accent',
    outline: 'border-border bg-transparent text-foreground',
  }[variant]

  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-medium',
        variantClassName,
        className,
      )}
      {...props}
    />
  )
}

export { Badge }
