import { type ButtonHTMLAttributes, forwardRef } from 'react'
import { cn } from '../../lib/utils'

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(
          'inline-flex items-center justify-center rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:pointer-events-none disabled:opacity-50',
          {
            primary: 'bg-blue-600 text-white hover:bg-blue-700',
            secondary: 'bg-zinc-800 text-zinc-100 hover:bg-zinc-700',
            ghost: 'text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800',
            danger: 'bg-red-600 text-white hover:bg-red-700',
          }[variant],
          {
            sm: 'h-8 px-3 text-sm',
            md: 'h-10 px-4 text-sm',
            lg: 'h-12 px-6 text-base',
          }[size],
          className
        )}
        {...props}
      />
    )
  }
)
Button.displayName = 'Button'
