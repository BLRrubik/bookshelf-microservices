import { useState, useCallback } from 'react';
import { ImagePlus, Loader2, Upload } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useUploadCover } from '@/api/covers';

interface CoverUploadProps {
  bookId: string;
  onSuccess?: () => void;
}

const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const MAX_SIZE = 5 * 1024 * 1024; // 5MB

export function CoverUpload({ bookId, onSuccess }: CoverUploadProps) {
  const [isDragging, setIsDragging] = useState(false);
  const upload = useUploadCover(bookId);

  const handleFile = useCallback((file: File) => {
    if (!ACCEPTED_TYPES.includes(file.type)) {
      return;
    }
    if (file.size > MAX_SIZE) {
      return;
    }
    
    upload.mutate(file, {
      onSuccess: () => {
        onSuccess?.();
      },
    });
  }, [upload, onSuccess]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    const file = e.dataTransfer.files[0];
    if (file) {
      handleFile(file);
    }
  }, [handleFile]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleClick = useCallback(() => {
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = ACCEPTED_TYPES.join(',');
    input.onchange = (e) => {
      const file = (e.target as HTMLInputElement).files?.[0];
      if (file) {
        handleFile(file);
      }
    };
    input.click();
  }, [handleFile]);

  return (
    <div
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onClick={handleClick}
      className={cn(
        "border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors",
        isDragging 
          ? "border-primary bg-primary/5" 
          : "border-border hover:border-primary/50 hover:bg-muted/50",
        upload.isPending && "pointer-events-none opacity-50"
      )}
    >
      {upload.isPending ? (
        <div className="flex flex-col items-center gap-2">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Загрузка...</p>
        </div>
      ) : (
        <div className="flex flex-col items-center gap-2">
          <div className="p-3 bg-muted rounded-full">
            {isDragging ? (
              <Upload className="h-6 w-6 text-primary" />
            ) : (
              <ImagePlus className="h-6 w-6 text-muted-foreground" />
            )}
          </div>
          <p className="font-medium">
            {isDragging ? 'Отпустите для загрузки' : 'Перетащите изображение'}
          </p>
          <p className="text-sm text-muted-foreground">
            или нажмите для выбора
          </p>
          <p className="text-xs text-muted-foreground mt-2">
            JPEG, PNG, WebP (до 5MB)
          </p>
        </div>
      )}
    </div>
  );
}

