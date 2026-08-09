import { useState, useCallback, useEffect } from 'react';
import { ImagePlus, Upload, X } from 'lucide-react';
import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';

interface CoverPickerProps {
  file: File | null;
  onFileChange: (file: File | null) => void;
  disabled?: boolean;
}

const ACCEPTED_TYPES = ['image/jpeg', 'image/png', 'image/webp'];
const MAX_SIZE = 5 * 1024 * 1024; // 5MB

export function CoverPicker({ file, onFileChange, disabled }: CoverPickerProps) {
  const [isDragging, setIsDragging] = useState(false);
  const [preview, setPreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Generate preview URL when file changes
  useEffect(() => {
    if (file) {
      const url = URL.createObjectURL(file);
      setPreview(url);
      return () => URL.revokeObjectURL(url);
    } else {
      setPreview(null);
    }
  }, [file]);

  const validateAndSetFile = useCallback((selectedFile: File) => {
    setError(null);
    
    if (!ACCEPTED_TYPES.includes(selectedFile.type)) {
      setError('Только JPEG, PNG или WebP');
      return;
    }
    if (selectedFile.size > MAX_SIZE) {
      setError('Файл слишком большой (макс. 5MB)');
      return;
    }
    
    onFileChange(selectedFile);
  }, [onFileChange]);

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    
    if (disabled) return;
    
    const droppedFile = e.dataTransfer.files[0];
    if (droppedFile) {
      validateAndSetFile(droppedFile);
    }
  }, [disabled, validateAndSetFile]);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    if (!disabled) {
      setIsDragging(true);
    }
  }, [disabled]);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  }, []);

  const handleClick = useCallback(() => {
    if (disabled) return;
    
    const input = document.createElement('input');
    input.type = 'file';
    input.accept = ACCEPTED_TYPES.join(',');
    input.onchange = (e) => {
      const selectedFile = (e.target as HTMLInputElement).files?.[0];
      if (selectedFile) {
        validateAndSetFile(selectedFile);
      }
    };
    input.click();
  }, [disabled, validateAndSetFile]);

  const handleRemove = useCallback((e: React.MouseEvent) => {
    e.stopPropagation();
    onFileChange(null);
    setError(null);
  }, [onFileChange]);

  // Show preview if file is selected
  if (preview && file) {
    return (
      <div className="relative">
        <div className="aspect-[3/4] rounded-lg overflow-hidden bg-muted">
          <img 
            src={preview} 
            alt="Предпросмотр обложки"
            className="w-full h-full object-cover"
          />
        </div>
        {!disabled && (
          <Button
            type="button"
            variant="destructive"
            size="icon"
            className="absolute top-2 right-2 h-8 w-8"
            onClick={handleRemove}
          >
            <X className="h-4 w-4" />
          </Button>
        )}
        <p className="text-xs text-muted-foreground mt-2 text-center truncate">
          {file.name}
        </p>
      </div>
    );
  }

  // Show drop zone
  return (
    <div
      onDrop={handleDrop}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onClick={handleClick}
      className={cn(
        "aspect-[3/4] border-2 border-dashed rounded-lg flex flex-col items-center justify-center cursor-pointer transition-colors",
        isDragging 
          ? "border-primary bg-primary/5" 
          : "border-border hover:border-primary/50 hover:bg-muted/50",
        disabled && "pointer-events-none opacity-50",
        error && "border-destructive"
      )}
    >
      <div className="flex flex-col items-center gap-2 p-4 text-center">
        <div className="p-3 bg-muted rounded-full">
          {isDragging ? (
            <Upload className="h-6 w-6 text-primary" />
          ) : (
            <ImagePlus className="h-6 w-6 text-muted-foreground" />
          )}
        </div>
        <p className="text-sm font-medium">
          {isDragging ? 'Отпустите' : 'Обложка'}
        </p>
        <p className="text-xs text-muted-foreground">
          Перетащите или нажмите
        </p>
        <p className="text-xs text-muted-foreground">
          JPEG, PNG, WebP (до 5MB)
        </p>
        {error && (
          <p className="text-xs text-destructive mt-1">{error}</p>
        )}
      </div>
    </div>
  );
}

