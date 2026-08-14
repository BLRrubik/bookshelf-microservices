import { BookOpen } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { useRateLimitStore } from "@/stores/rateLimit";

export function Footer() {
  const { remaining, limit } = useRateLimitStore();

  return (
    <footer className="border-t py-6 mt-auto">
      <div className="container flex flex-col sm:flex-row items-center justify-between gap-4">
        <div className="flex items-center gap-2 text-muted-foreground">
          <BookOpen className="h-4 w-4" />
          <span className="text-sm">Bookshelf by Praxis © 2026</span>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-xs text-muted-foreground">Rate Limit:</span>
          <Badge
            variant={
              remaining > 10
                ? "outline"
                : remaining > 0
                ? "secondary"
                : "destructive"
            }
            className="text-xs"
          >
            {remaining}/{limit}
          </Badge>
        </div>

        <p className="text-sm text-muted-foreground">Проект 4: API Gateway</p>
      </div>
    </footer>
  );
}
