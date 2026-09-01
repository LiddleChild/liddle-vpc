import { cn } from "@/lib/utils"

type SkeletonProps = React.HTMLAttributes<HTMLSpanElement>

export function Skeleton({ className, ...props }: SkeletonProps) {
  return <span aria-hidden="true" className={cn("skeleton", className)} {...props} />
}
