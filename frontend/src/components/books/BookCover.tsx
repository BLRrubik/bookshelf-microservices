import { BookOpen, Loader2 } from 'lucide-react';
import { cn } from '@/lib/utils';
import type { CoverStatus } from '@/types/api';

interface BookCoverProps {
  coverUrl?: string | null;
  coverStatus: CoverStatus;
  title: string;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

const sizeClasses = {
  sm: 'h-24 w-16',
  md: 'h-48 w-32',
  lg: 'h-96 w-64',
};

export function BookCover({ 
  coverUrl, 
  coverStatus, 
  title, 
  size = 'md',
  className 
}: BookCoverProps) {
  const showPlaceholder = !coverUrl && coverStatus !== 'processing';
  const showProcessing = coverStatus === 'processing';

  return (
    <div 
      className={cn(
        "relative bg-muted rounded-lg flex items-center justify-center overflow-hidden",
        sizeClasses[size],
        className
      )}
    >
      {coverUrl && coverStatus === 'ready' ? (
        <img 
          src={coverUrl} 
          alt={`Обложка книги "${title}"`}
          className="w-full h-full object-cover"
          loading="lazy"
        />
      ) : showProcessing ? (
        <div className="flex flex-col items-center gap-2 text-muted-foreground">
          <Loader2 className={cn(
            "animate-spin",
            size === 'sm' ? 'h-4 w-4' : size === 'md' ? 'h-8 w-8' : 'h-12 w-12'
          )} />
          {size !== 'sm' && (
            <span className="text-xs">Обработка...</span>
          )}
        </div>
      ) : showPlaceholder ? (
        <BookOpen className={cn(
          "text-muted-foreground",
          size === 'sm' ? 'h-6 w-6' : size === 'md' ? 'h-12 w-12' : 'h-24 w-24'
        )} />
      ) : null}
    </div>
  );
}

