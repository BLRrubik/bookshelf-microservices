import { useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { Loader2, CheckCircle, XCircle, Clock } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Progress } from '@/components/ui/progress';
import { useCoverStatus } from '@/api/covers';
import type { CoverStatus as CoverStatusType } from '@/types/api';

interface CoverStatusProps {
  bookId: string;
  initialStatus?: CoverStatusType;
}

export function CoverStatus({ bookId, initialStatus = 'none' }: CoverStatusProps) {
  const queryClient = useQueryClient();
  const { data, isLoading } = useCoverStatus(bookId, initialStatus);
  const prevStatusRef = useRef<string | undefined>(initialStatus);
  
  const status = data?.status || initialStatus;

  // When cover processing completes, invalidate book query to get updated cover_url
  useEffect(() => {
    const prevStatus = prevStatusRef.current;
    
    // If status changed from 'processing' to 'ready' or 'failed', refresh book data
    if (prevStatus === 'processing' && (status === 'ready' || status === 'failed')) {
      queryClient.invalidateQueries({ queryKey: ['books', bookId] });
    }
    
    prevStatusRef.current = status;
  }, [status, bookId, queryClient]);

  if (status === 'none') {
    return null;
  }

  if (isLoading && initialStatus === 'none') {
    return null;
  }

  return (
    <div className="space-y-2">
      {status === 'processing' && (
        <>
          <div className="flex items-center gap-2">
            <Loader2 className="h-4 w-4 animate-spin text-primary" />
            <span className="text-sm text-muted-foreground">
              Обработка обложки...
            </span>
          </div>
          <Progress value={undefined} className="h-2" />
        </>
      )}

      {status === 'ready' && (
        <Badge variant="success" className="flex items-center gap-1 w-fit">
          <CheckCircle className="h-3 w-3" />
          Обложка готова
        </Badge>
      )}

      {status === 'failed' && (
        <div className="space-y-1">
          <Badge variant="destructive" className="flex items-center gap-1 w-fit">
            <XCircle className="h-3 w-3" />
            Ошибка обработки
          </Badge>
          {data?.error && (
            <p className="text-xs text-destructive">{data.error}</p>
          )}
        </div>
      )}

      {data?.started_at && status === 'processing' && (
        <p className="text-xs text-muted-foreground flex items-center gap-1">
          <Clock className="h-3 w-3" />
          Начато: {new Date(data.started_at).toLocaleTimeString('ru-RU')}
        </p>
      )}
    </div>
  );
}




