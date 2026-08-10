import { Badge } from '@/components/ui/badge';
import { useAuthHealth, useBooksHealth } from '@/api/health';

export function ServiceStatus() {
  const authHealth = useAuthHealth();
  const booksHealth = useBooksHealth();

  const getStatus = (isLoading: boolean, isError: boolean, data?: { ready: boolean }) => {
    if (isLoading) return 'loading';
    if (isError) return 'error';
    return data?.ready ? 'ok' : 'error';
  };

  const authStatus = getStatus(authHealth.isLoading, authHealth.isError, authHealth.data);
  const booksStatus = getStatus(booksHealth.isLoading, booksHealth.isError, booksHealth.data);

  return (
    <div className="flex items-center gap-2">
      <span className="text-xs text-muted-foreground">Services:</span>
      <Badge 
        variant={authStatus === 'ok' ? 'success' : authStatus === 'loading' ? 'outline' : 'destructive'}
        className="text-xs"
      >
        Auth {authStatus === 'ok' ? '●' : authStatus === 'loading' ? '○' : '✕'}
      </Badge>
      <Badge 
        variant={booksStatus === 'ok' ? 'success' : booksStatus === 'loading' ? 'outline' : 'destructive'}
        className="text-xs"
      >
        Books {booksStatus === 'ok' ? '●' : booksStatus === 'loading' ? '○' : '✕'}
      </Badge>
    </div>
  );
}

