import { AlertCircle, Clock } from 'lucide-react';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Progress } from '@/components/ui/progress';
import { useRateLimitStore } from '@/stores/rateLimit';

export function RateLimitWarning() {
  const { isWarning, remaining, limit, isBlocked, retryAfter } = useRateLimitStore();

  if (isBlocked) {
    return (
      <Alert variant="destructive" className="fixed bottom-4 left-4 right-4 max-w-md z-50">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Превышен лимит запросов</AlertTitle>
        <AlertDescription className="space-y-2">
          <p>
            Подождите <strong>{retryAfter}</strong> секунд перед следующим запросом.
          </p>
          <div className="flex items-center gap-2 text-xs">
            <Clock className="h-3 w-3" />
            <span>Автоматически разблокируется</span>
          </div>
        </AlertDescription>
      </Alert>
    );
  }

  if (isWarning && remaining < 10) {
    const usagePercent = ((limit - remaining) / limit) * 100;
    
    return (
      <Alert variant="warning" className="fixed bottom-4 left-4 right-4 max-w-md z-50">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Предупреждение</AlertTitle>
        <AlertDescription className="space-y-2">
          <p>
            Осталось <strong>{remaining}</strong> из {limit} запросов.
          </p>
          <Progress value={usagePercent} className="h-2" />
        </AlertDescription>
      </Alert>
    );
  }

  return null;
}





