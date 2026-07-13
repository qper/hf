export function Skeleton({
  className = '',
}: {
  className?: string
}) {
  return (
    <div
      className={`bg-zinc-800 animate-pulse rounded ${className}`}
    />
  )
}
