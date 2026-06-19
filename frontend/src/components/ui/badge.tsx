import type { HTMLAttributes } from 'react'
import { cn } from '../../lib/utils'

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'danger'
}

export function Badge({
  className,
  variant = 'default',
  ...props
}: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium',
        {
          default: 'bg-zinc-800 text-zinc-300',
          success: 'bg-emerald-900/50 text-emerald-400',
          warning: 'bg-amber-900/50 text-amber-400',
          danger: 'bg-red-900/50 text-red-400',
        }[variant],
        className
      )}
      {...props}
    />
  )
}
